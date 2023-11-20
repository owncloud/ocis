package ocis

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"ociswrapper/common"
	"ociswrapper/log"
	"ociswrapper/ocis/config"
)

var cmd *exec.Cmd
var retryCount = 0

func Start(envMap map[string]any) {
	if retryCount == 0 {
		defer common.Wg.Done()
	}

	cmd = exec.Command(config.Get("bin"), "server")
	cmd.Env = os.Environ()
	var environments []string
	if envMap != nil {
		for key, value := range envMap {
			environments = append(environments, fmt.Sprintf("%s=%v", key, value))
		}
	}
	cmd.Env = append(cmd.Env, environments...)

	logs, err := cmd.StderrPipe()
	if err != nil {
		log.Panic(err)
	}
	output, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Panic(err)
	}

	// Read and print the logs when the 'ocis server' command is running
	logScanner := bufio.NewScanner(logs)
	for logScanner.Scan() {
		m := logScanner.Text()
		fmt.Println(m)
	}
	// Read output when the 'ocis server' command gets exited
	outputScanner := bufio.NewScanner(output)
	for outputScanner.Scan() {
		m := outputScanner.Text()
		fmt.Println(m)
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.Sys().(syscall.WaitStatus)
			// retry only if oCIS server exited with code > 0
			// -1 exit code means that the process was killed by a signal (cmd.Process.Kill())
			if status.ExitStatus() > 0 {
				log.Println(fmt.Sprintf("oCIS server exited with code %v", status.ExitStatus()))

				// retry to start oCIS server
				retryCount++
				maxRetry, _ := strconv.Atoi(config.Get("retry"))
				if retryCount <= maxRetry {
					log.Println(fmt.Sprintf("Retry starting oCIS server... (retry %v)", retryCount))
					// Stop and start again
					Stop()
					Start(envMap)
				}
			}
		}
	}
}

func Stop() {
	err := cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		log.Println("here......")
		log.Println(err.Error())
		if !strings.HasSuffix(err.Error(), "process already finished") {
			log.Fatalln(err)
		}
	}
	cmd.Wait()
}

func WaitForConnection() bool {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 30 seconds timeout
	timeoutValue := 30 * time.Second

	client := http.Client{
		Timeout:   timeoutValue,
		Transport: transport,
	}

	req, _ := http.NewRequest("GET", config.Get("url")+"/graph/v1.0/users/"+config.Get("adminUsername"), nil)
	req.SetBasicAuth(config.Get("adminUsername"), config.Get("adminPassword"))

	timeout := time.After(timeoutValue)

	for {
		select {
		case <-timeout:
			log.Println(fmt.Sprintf("%v seconds timeout waiting for oCIS server", int64(timeoutValue.Seconds())))
			return false
		default:
			req.Header.Set("X-Request-ID", "ociswrapper-"+strconv.Itoa(int(time.Now().UnixMilli())))

			res, err := client.Do(req)
			if err != nil || res.StatusCode != 200 {
				// 500 milliseconds poll interval
				time.Sleep(500 * time.Millisecond)
				continue
			}

			log.Println("oCIS server is ready to accept requests")
			return true
		}
	}
}

func Restart(envMap map[string]any) bool {
	log.Println("Stopping oCIS server...")
	Stop()

	time.Sleep(5 * time.Second)

	common.Wg.Add(1)
	go Start(envMap)
	log.Println("Restarting oCIS server...")

	time.Sleep(5 * time.Second)

	return WaitForConnection()
}
