package command

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tw "github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/pkg/xattr"
	"github.com/rogpeppe/go-internal/lockedfile"
	"github.com/urfave/cli/v2"
)

// BenchmarkCommand is the entrypoint for the benchmark commands.
func BenchmarkCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "benchmark",
		Usage:       "cli tools to test low and high level performance",
		Category:    "benchmark",
		Subcommands: []*cli.Command{BenchmarkClientCommand(cfg), BenchmarkSyscallsCommand(cfg)},
	}
}

// BenchmarkClientCommand is the entrypoint for the benchmark client command.
func BenchmarkClientCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name: "client",

		Usage: "Start a client that continuously makes web requests and prints stats. The options mimic curl, but URL must be at the end.",
		Flags: []cli.Flag{

			// TODO with v3 'flag.Persistent: true' can be set to make the order of flags no longer relevant \o/
			// flags mimicing curl
			&cli.StringFlag{
				Name:    "request",
				Aliases: []string{"X"},
				Value:   "PROPFIND",
				Usage:   "Specifies a custom request method to use when communicating with the HTTP server.",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Value:   "admin:admin",
				Usage:   "Specify the user name and password to use for server authentication.",
			},
			&cli.BoolFlag{
				Name:    "insecure",
				Aliases: []string{"k"},
				Usage:   "Skip the TLS verification step and proceed without checking.",
			},
			&cli.StringFlag{
				Name:    "data",
				Aliases: []string{"d"},
				Usage:   "Sends the specified data in a request to the HTTP server.",
				// TODE support multiple data flags, support data-binary, data-raw
			},
			&cli.StringSliceFlag{
				Name:    "header",
				Aliases: []string{"H"},
				Usage:   "Extra header to include in information sent.",
			},
			&cli.StringFlag{
				Name: "rate",
				Usage: `Specify the maximum transfer frequency you allow a client to use - in number of transfer starts per time unit (sometimes called request rate).
	The request rate is provided as "N/U" where N is an integer number and U is a time unit. Supported units are 's' (second), 'm' (minute), 'h' (hour) and 'd' /(day, as in a 24 hour unit). The default time unit, if no "/U" is provided, is number of transfers per hour.`,
			},
			/*
				&cli.StringFlag{
					Name:    "oauth2-bearer",
					Usage:   "Specify the Bearer Token for OAUTH 2.0 server authentication.",
				},
				&cli.StringFlag{
					Name:    "user-agent",
					Aliases: []string{"A"},
					Value:   "admin:admin",
					Usage:   "Specify the User-Agent string to send to the HTTP	server.",
				},
			*/
			// other flags
			&cli.StringFlag{
				Name:  "bearer-token-command",
				Usage: "Command to execute for a bearer token, e.g. 'oidc-token OCIS'. When set, disables basic auth.",
			},
			&cli.IntFlag{
				Name:  "every",
				Usage: "Aggregate stats every time this amount of seconds has passed.",
			},
			&cli.IntFlag{
				Name:    "jobs",
				Aliases: []string{"j"},
				Value:   1,
				Usage:   "Number of parallel clients to start.",
			},
		},
		Category: "benchmark",
		Action: func(c *cli.Context) error {
			opt := clientOptions{
				request:  c.String("request"),
				url:      c.Args().First(),
				insecure: c.Bool("insecure"),
				jobs:     c.Int("jobs"),
				headers:  make(map[string]string),
				data:     []byte(c.String("data")),
			}
			if opt.url == "" {
				log.Fatal(errors.New("no URL specified"))
			}

			for _, h := range c.StringSlice("headers") {
				parts := strings.SplitN(h, ":", 2)
				if len(parts) != 2 {
					log.Fatal(errors.New("invalid header '" + h + "'"))
				}
				opt.headers[parts[0]] = strings.TrimSpace(parts[1])
			}

			rate := c.String("rate")
			if rate != "" {
				parts := strings.SplitN(rate, "/", 2)
				num, err := strconv.Atoi(parts[0])
				if err != nil {
					fmt.Println(err)
				}
				unit := time.Hour // default
				if len(parts) == 2 {
					switch parts[1] {
					case "s":
						unit = time.Second
					case "m":
						unit = time.Minute
					case "d":
						unit = time.Hour * 24
					default:
						log.Fatal(errors.New("unsupported rate unit. Use s, m, h or d"))
					}
				}
				opt.rateDelay = unit / time.Duration(num)
			}

			user := c.String("user")
			opt.auth = func() string {
				return "Basic " + base64.StdEncoding.EncodeToString([]byte(user))
			}

			btc := c.String("bearer-token-command")
			if btc != "" {
				parts := strings.SplitN(btc, " ", 2)
				var cmd *exec.Cmd
				opt.auth = func() string {
					if len(parts) > 1 {
						cmd = exec.Command(parts[0], parts[1])
					} else {
						cmd = exec.Command(parts[0])
					}
					output, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Println(err)
					}
					return "Bearer " + string(output)
				}
			}

			every := c.Int("every")
			if every != 0 {
				opt.ticker = time.NewTicker(time.Second * time.Duration(every))
				defer opt.ticker.Stop()
			}

			return client(opt)

		},
	}
}

type clientOptions struct {
	request   string
	url       string
	auth      func() string
	insecure  bool
	headers   map[string]string
	rateDelay time.Duration
	data      []byte
	ticker    *time.Ticker
	jobs      int
}

func client(o clientOptions) error {

	type stat struct {
		job      int
		duration time.Duration
		status   int
	}
	stats := make(chan stat)
	for i := 0; i < o.jobs; i++ {
		go func(i int) {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion:         tls.VersionTLS12,
					InsecureSkipVerify: o.insecure,
				},
			}
			client := &http.Client{Transport: tr}

			cookies := map[string]*http.Cookie{}
			for {
				req, err := http.NewRequest(o.request, o.url, bytes.NewReader(o.data))
				if err != nil {
					log.Printf("client %d: could not create request: %s\n", i, err)
					return
				}
				req.Header.Set("Authorization", strings.TrimSpace(o.auth()))
				for k, v := range o.headers {
					req.Header.Set(k, v)
				}
				for _, cookie := range cookies {
					req.AddCookie(cookie)
				}

				start := time.Now()
				res, err := client.Do(req)
				duration := -time.Until(start)
				if err != nil {
					log.Printf("client %d: could not create request: %s\n", i, err)
					time.Sleep(time.Second)
				} else {
					res.Body.Close()
					stats <- stat{
						job:      i,
						duration: duration,
						status:   res.StatusCode,
					}
					for _, c := range res.Cookies() {
						cookies[c.Name] = c
					}
				}
				time.Sleep(o.rateDelay - duration)
			}
		}(i)
	}

	numRequests := 0
	if o.ticker == nil {
		// no ticker, just write every request
		for {
			stat := <-stats
			numRequests++
			fmt.Printf("req %d took %v and returned status %d\n", numRequests, stat.duration, stat.status)
		}
	}

	var duration time.Duration
	for {
		select {
		case stat := <-stats:
			numRequests++
			duration += stat.duration
		case <-o.ticker.C:
			if numRequests > 0 {
				fmt.Printf("%d req at %v/req\n", numRequests, duration/time.Duration(numRequests))
				numRequests = 0
				duration = 0
			}
		}
	}

}

// BenchmarkSyscallsCommand is the entrypoint for the benchmark syscalls command.
func BenchmarkSyscallsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "syscalls",
		Usage: "test the performance of syscalls",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "Path to test",
			},
			&cli.StringFlag{
				Name:  "iterations",
				Value: "100",
				Usage: "Number of iterations to execute",
			},
		},
		Category: "benchmark",
		Action: func(c *cli.Context) error {

			path := c.String("path")
			if path == "" {
				f, err := os.CreateTemp("", "ocis-bench-temp-")
				if err != nil {
					log.Fatal(err)
				}
				path = f.Name()
				f.Close()
				defer os.Remove(path)
			}

			iterations := c.Int("iterations")

			return benchmark(iterations, path)
		},
	}
}

func benchmark(iterations int, path string) error {
	tests := map[string]func() error{
		"lockedfile open(wo,c,t) close": func() error {
			for i := 0; i < iterations; i++ {
				lockedFile, err := lockedfile.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
				if err != nil {
					return err
				}
				lockedFile.Close()
			}
			return nil
		},
		"stat": func() error {
			for i := 0; i < iterations; i++ {
				_, err := os.Stat(path)
				if err != nil {
					return err
				}
			}
			return nil
		},
		"fopen(ro) close": func() error {
			for i := 0; i < iterations; i++ {
				h, err := os.OpenFile(path, os.O_RDONLY, 0600)
				if err != nil {
					return err
				}
				h.Close()
			}
			return nil
		},
		"fopen(wo,t) write close": func() error {
			for i := 0; i < iterations; i++ {
				h, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0600)
				if err != nil {
					return err
				}
				_, err = h.WriteString("1234567890")
				if err != nil {
					h.Close()
					return err
				}
				h.Close()
			}
			return nil
		},
		"fopen(ro) read close": func() error {
			for i := 0; i < iterations; i++ {
				bytes := make([]byte, 0, 10)
				h, err := os.OpenFile(path, os.O_RDONLY, 0600)
				if err != nil {
					return err
				}
				_, err = h.Read(bytes)
				if err != nil {
					h.Close()
					return err
				}
				h.Close()
			}
			return nil
		},
		"xattr-set": func() error {
			for i := 0; i < iterations; i++ {
				err := xattr.Set(path, "user.test", []byte("123456"))
				if err != nil {
					return err
				}
			}
			return nil
		},
		"xattr-get": func() error {
			for i := 0; i < iterations; i++ {
				_, err := xattr.Get(path, "user.test")
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	fmt.Println("Version: " + version.GetString())
	fmt.Printf("Compiled: %s\n", version.Compiled())
	fmt.Printf("Path: %s\n", path)
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Println("")

	table := tw.NewWriter(os.Stdout)
	table.SetHeader([]string{"Test", "Iterations", "dur/it", "total"})
	table.SetAutoFormatHeaders(false)
	table.SetColumnAlignment([]int{tw.ALIGN_LEFT, tw.ALIGN_RIGHT, tw.ALIGN_RIGHT, tw.ALIGN_RIGHT})
	table.SetAutoMergeCellsByColumnIndex([]int{2, 3})
	for _, t := range []string{"lockedfile open(wo,c,t) close", "stat", "fopen(wo,t) write close", "fopen(ro) close", "fopen(ro) read close", "xattr-set", "xattr-get"} {
		start := time.Now()
		err := tests[t]()
		end := time.Now()
		delta := end.Sub(start)
		if err != nil {
			table.Append([]string{t, fmt.Sprintf("%d", iterations), err.Error(), err.Error()})
		} else {
			table.Append([]string{t, fmt.Sprintf("%d", iterations), strconv.Itoa(int(delta.Nanoseconds())/iterations) + "ns", delta.String()})
		}
	}
	table.Render()
	return nil
}

func init() {
	register.AddCommand(BenchmarkCommand)
}
