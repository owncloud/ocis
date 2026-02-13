package ocis

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"ociswrapper/common"
	"ociswrapper/log"
	"ociswrapper/ocis/config"

	"github.com/creack/pty"
)

var cmd *exec.Cmd
var retryCount = 0
var stopSignal = false
var runningServices = make(map[string]int)

// exported variables
var ServiceEnvConfigs = make(map[string][]string)
var OcisServiceName = "ocis"

func Start(envMap []string) {
	log.Println("Starting oCIS service...")
	StartService(OcisServiceName, envMap)
}

func Stop() (bool, string) {
	log.Println("Stopping oCIS server...")
	stopSignal = true

	var stopErrors []string
	for service := range runningServices {
		success, message := StopService(service)
		if !success {
			stopErrors = append(stopErrors, message)
		}
	}
	if len(stopErrors) > 0 {
		log.Println("Errors occurred while stopping services:")
		for _, err := range stopErrors {
			log.Println(err)
		}
	}

	success, message := waitUntilCompleteShutdown()

	cmd = nil
	return success, message
}

func Restart(envMap []string) (bool, string) {
	Stop()

	common.Wg.Add(1)
	go Start(envMap)

	return WaitForConnection()
}

func IsOcisRunning() bool {
	if runningServices[OcisServiceName] == 0 {
		return false
	}

	_, err := os.FindProcess(runningServices[OcisServiceName])
	if err != nil {
		delete(runningServices, OcisServiceName)
		return false
	}
	return true
}

func waitAllServices(startTime time.Time, timeout time.Duration) {
	timeoutS := timeout * time.Second

	c := exec.Command(config.Get("bin"), "list")
	_, err := c.CombinedOutput()
	if err != nil {
		if time.Since(startTime) <= timeoutS {
			time.Sleep(500 * time.Millisecond)
			waitAllServices(startTime, timeout)
		}
		return
	}
}

func WaitForConnection() (bool, string) {
	waitAllServices(time.Now(), 30)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 30 seconds timeout
	timeoutValue := 30 * time.Second

	client := http.Client{
		Timeout:   timeoutValue,
		Transport: transport,
	}

	var req *http.Request
	if config.Get("adminUsername") != "" && config.Get("adminPassword") != "" {
		req, _ = http.NewRequest("GET", config.Get("url")+"/graph/v1.0/me/drives", nil)
		req.SetBasicAuth(config.Get("adminUsername"), config.Get("adminPassword"))
	} else {
		req, _ = http.NewRequest("GET", config.Get("url")+"/ocs/v1.php/cloud/capabilities?format=json", nil)
	}

	timeout := time.After(timeoutValue)

	for {
		select {
		case <-timeout:
			log.Println(fmt.Sprintf("%v seconds timeout waiting for oCIS server", int64(timeoutValue.Seconds())))
			return false, "Timeout waiting for oCIS server to start"
		default:
			req.Header.Set("X-Request-ID", "ociswrapper-"+strconv.Itoa(int(time.Now().UnixMilli())))

			res, err := client.Do(req)
			if err != nil || res.StatusCode != 200 {
				// 500 milliseconds poll interval
				time.Sleep(500 * time.Millisecond)
				continue
			}

			log.Println("oCIS server is ready to accept requests")
			return true, "oCIS server is up and running"
		}
	}
}

func waitUntilCompleteShutdown() (bool, string) {
	timeout := 30 * time.Second
	startTime := time.Now()

	c := exec.Command("sh", "-c", "ps ax | grep 'ocis server' | grep -v grep | awk '{print $1}'")
	output, err := c.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}
	for strings.TrimSpace(string(output)) != "" {
		output, _ = c.CombinedOutput()
		log.Println("Process found. Waiting...")

		if time.Since(startTime) >= timeout {
			log.Println(fmt.Sprintf("Unable to kill oCIS server after %v seconds", int64(timeout.Seconds())))
			return false, "Timeout waiting for oCIS server to stop"
		}
	}
	return true, "oCIS server stopped successfully"
}

func RunCommand(command string, inputs []string) (int, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// build the command
	cmdArgs := strings.Split(command, " ")
	c := exec.CommandContext(ctx, config.Get("bin"), cmdArgs...)

	// Start the command with a pty (pseudo terminal)
	// This is required to interact with the command
	ptyF, err := pty.Start(c)
	if err != nil {
		log.Fatalln(err)
	}
	defer ptyF.Close()

	var output bytes.Buffer
	done := make(chan error, 1)
	// read concurrently from the pty
	go func() {
		_, err := io.Copy(&output, ptyF)
		done <- err
	}()

	// send inputs to the command
	for _, input := range inputs {
		fmt.Fprintf(ptyF, "%s\n", input)
	}

	var timeoutErr string
	if err := c.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			timeoutErr = "[TimeoutError] Command timed out\n"
		}
	}
	// wait for copy to finish
	<-done

	cmdOutput := output.String()
	if timeoutErr != "" {
		cmdOutput += timeoutErr
	}

	// remove stdins from the output
	cmdOutput = strings.TrimLeft(cmdOutput, strings.Join(inputs, "\r\n"))

	return c.ProcessState.ExitCode(), cmdOutput
}

func StartService(service string, envMap []string) {
	// Initialize command args based on service presence
	cmdArgs := []string{"server"} // Default command args

	if envMap == nil {
		envMap = []string{}
	}

	if service == "" {
		service = OcisServiceName
	} else if service != OcisServiceName {
		cmdArgs = append([]string{service}, cmdArgs...)
	}

	// wait for the log scanner to finish
	var wg sync.WaitGroup
	wg.Add(2)

	stopSignal = false
	if retryCount == 0 {
		defer common.Wg.Done()
	}

	cmd = exec.Command(config.Get("bin"), cmdArgs...)
	cmd.Env = append(os.Environ(), append(ServiceEnvConfigs[service], envMap...)...)

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

	logScanner := bufio.NewScanner(logs)
	outputScanner := bufio.NewScanner(output)
	outChan := make(chan string)

	// save the service process ID
	runningServices[service] = cmd.Process.Pid

	// Read the logs when the 'ocis server' command is running
	go func() {
		defer wg.Done()
		for logScanner.Scan() {
			outChan <- logScanner.Text()
		}
	}()

	go func() {
		defer wg.Done()
		for outputScanner.Scan() {
			outChan <- outputScanner.Text()
		}
	}()

	// Fetch logs from the channel and print them
	go func() {
		for s := range outChan {
			fmt.Println(s)
		}
	}()

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
					wg.Wait()
					close(outChan)
					log.Println(fmt.Sprintf("Retry starting oCIS server... (retry %v)", retryCount))
					// wait 500 milliseconds before retrying
					time.Sleep(500 * time.Millisecond)
					StartService(service, envMap)
					return
				}
			}
		}
	}
	wg.Wait()
	close(outChan)
}

// Stop oCIS service or a specific service by its unique identifier
func StopService(service string) (bool, string) {
	pid, exists := runningServices[service]
	if !exists {
		return false, fmt.Sprintf("Service %s is not running", service)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false, fmt.Sprintf("Failed to find service %s process running with ID %d", service, pid)
	}

	pKillError := process.Signal(syscall.SIGINT)
	if pKillError != nil {
		return false, fmt.Sprintf("Failed to stop service with process id %d", pid)
	}

	success := WaitForServiceStatus(service, false, false)
	if !success {
		StopService(service)
	}

	delete(runningServices, service)
	return true, fmt.Sprintf(fmt.Sprintf("%s service stopped successfully", service))
}

func WaitForServiceStatus(service string, waitForUp bool, isRunningInK8s bool) bool {
	overallTimeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	port := config.GetServiceDebugPort(service)

	for {
		select {
		case <-overallTimeout:
			status := "available"
			if !waitForUp {
				status = "shut down"
			}
			log.Println(fmt.Errorf("Timeout: %s service did not %s within 30 seconds", service, status).Error())
			return false
		case <-ticker.C:
			if !isRunningInK8s {
				if _, exists := runningServices[service]; !exists {
					log.Println(fmt.Sprintf("Service %s not found in running services. Retrying...", service))
					continue
				}
			}

			address := fmt.Sprintf(":%d", port)
			conn, err := net.DialTimeout("tcp", address, 1*time.Second)
			if waitForUp {
				if err == nil {
					_ = conn.Close()
					log.Println(fmt.Sprintf("%s service is ready to listen on port %d", service, port))
					return true
				}
				log.Println(fmt.Sprintf("%s service is not ready on port %d. %v", service, port, err))
			} else {
				if err != nil {
					log.Println(fmt.Sprintf("%s service port %d is no longer reachable", service, port))
					return true
				}
				_ = conn.Close()
				log.Println(fmt.Sprintf("%s service port %d is still active. Retrying...", service, port))
			}
		}
	}
}
