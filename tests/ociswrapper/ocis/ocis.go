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
var stopSignal = false

func Start(envMap map[string]any) {
	stopSignal = false
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
			// -1 exit code means that the process was killed by a signal (syscall.SIGINT)
			if status.ExitStatus() > 0 && !stopSignal {
				waitUntilCompleteShutdown()

				log.Println(fmt.Sprintf("oCIS server exited with code %v", status.ExitStatus()))

				// retry to start oCIS server
				retryCount++
				maxRetry, _ := strconv.Atoi(config.Get("retry"))
				if retryCount <= maxRetry {
					log.Println(fmt.Sprintf("Retry starting oCIS server... (retry %v)", retryCount))
					Start(envMap)
				}
			}
		}
	}
}

func Stop() {
	stopSignal = true

	// SIGINT allows oCIS server to gracefully shutdown
	err := cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "process already finished") {
			log.Fatalln(err)
		}
	}
	cmd.Process.Wait()
	waitUntilCompleteShutdown()
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

func waitUntilCompleteShutdown() {
	timeout := 30 * time.Second
	startTime := time.Now()

	c := exec.Command("sh", "-c", "ps ax | grep 'ocis server' | grep -v grep | awk '{print $1}'")
	output, err := c.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}
	for strings.TrimSpace(string(output)) != "" {
		output, _ = c.CombinedOutput()

		if time.Since(startTime) >= timeout {
			log.Println(fmt.Sprintf("Unable to kill oCIS server after %v seconds", int64(timeout.Seconds())))
			break
		}
	}
}

func Restart(envMap map[string]any) bool {
	log.Println("Restarting oCIS server...")
	Stop()

	common.Wg.Add(1)
	go Start(envMap)

	return WaitForConnection()
}
