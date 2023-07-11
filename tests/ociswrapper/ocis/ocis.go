package ocis

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"ociswrapper/common"
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

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		m := stderrScanner.Text()
		fmt.Println(m)
	}
	stdoutScanner := bufio.NewScanner(stdout)
	for stdoutScanner.Scan() {
		m := stdoutScanner.Text()
		fmt.Println(m)
		retryCount++

		maxRetry, _ := strconv.Atoi(config.Get("retry"))
		if retryCount <= maxRetry {
			fmt.Println(fmt.Sprintf("Retry starting oCIS server... (retry %v)", retryCount))
			// Stop and start again
			Stop()
			Start(envMap)
		}
	}
}

func Stop() {
	err := cmd.Process.Kill()
	if err != nil {
		log.Panic("Cannot kill oCIS server")
	}
	cmd.Wait()
}

func WaitForConnection() bool {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 5 seconds timeout
	timeoutValue := 5 * time.Second

	client := http.Client{
		Timeout:   timeoutValue,
		Transport: transport,
	}

	timeout := time.After(timeoutValue)

	for {
		select {
		case <-timeout:
			fmt.Println(fmt.Sprintf("%v seconds timeout waiting for oCIS server", int64(timeoutValue.Seconds())))
			return false
		default:
			_, err := client.Get(config.Get("url"))
			if err != nil {
				fmt.Println("Waiting for oCIS server...")
			} else {
				fmt.Println(fmt.Sprintf("oCIS server is ready to accept requests"))
				return true
			}
			// 500 milliseconds poll interval
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func Restart(envMap map[string]any) bool {
	log.Print("Restarting oCIS server...")
	Stop()

	common.Wg.Add(1)
	go Start(envMap)

	return WaitForConnection()
}
