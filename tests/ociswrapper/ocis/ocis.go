package ocis

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
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
var EnvConfigs = []string{}
var runningServices = make(map[string]int)

func Start(envMap []string) {
	log.Println("Starting oCIS service........")
	StartService("", envMap)
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

	log.Println("Restarting oCIS server...")
	common.Wg.Add(1)
	go Start(envMap)

	return WaitForConnection()
}

func IsOcisRunning() bool {
	if runningServices["ocis"] == 0 {
		return false
	}

	_, err := os.FindProcess(runningServices["ocis"])
	if err != nil {
		delete(runningServices, "ocis")
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
	log.Println("All services are up")
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

	c := exec.Command("sh", "-c", "ps ax | grep '[o]cis server' | grep -v grep | awk '{print $1}'")
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
	logs := new(strings.Builder)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// build the command
	cmdArgs := strings.Split(command, " ")
	c := exec.CommandContext(ctx, config.Get("bin"), cmdArgs...)

	// Start the command with a pty (pseudo terminal)
	// This is required to interact with the command
	ptyF, err := pty.Start(c)
	if err != nil {
		log.Panic(err)
	}
	defer ptyF.Close()

	for _, input := range inputs {
		fmt.Fprintf(ptyF, "%s\n", input)
	}

	var cmdOutput string
	if err := c.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			cmdOutput = "Command timed out:\n"
		}
	}

	// Copy the logs from the pty
	io.Copy(logs, ptyF)
	cmdOutput += logs.String()

	// TODO: find if there is a better way to remove stdins from the output
	cmdOutput = strings.TrimLeft(cmdOutput, strings.Join(inputs, "\r\n"))

	return c.ProcessState.ExitCode(), cmdOutput
}

func StartService(service string, envMap []string) {
	// Initialize command args based on service presence
	cmdArgs := []string{"server"} // Default command args

	if service != "" {
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

	if len(envMap) == 0 {
		cmd.Env = append(os.Environ(), EnvConfigs...)
	} else {
		cmd.Env = append(os.Environ(), envMap...)
	}

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

	if service == "" {
		runningServices["ocis"] = cmd.Process.Pid
	} else {
		runningServices[service] = cmd.Process.Pid
	}

	for listService, pid := range runningServices {
		log.Println(fmt.Sprintf("%s service started with process id %v", listService, pid))
	}

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
	stopWg := new(sync.WaitGroup)
	stopWg.Add(1)

	resultChan := make(chan struct {
		success bool
		message string
	}, 1)

	go func() {
		defer stopWg.Done()

		pid, exists := runningServices[service]
		if !exists {
			resultChan <- struct {
				success bool
				message string
			}{false, fmt.Sprintf("Service %s is not running", service)}
			return
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			resultChan <- struct {
				success bool
				message string
			}{false, fmt.Sprintf("Failed to find service %s process running with ID %d", service, pid)}
			return
		}

		pKillError := process.Signal(syscall.SIGINT)
		if pKillError != nil {
			resultChan <- struct {
				success bool
				message string
			}{false, fmt.Sprintf("Failed to stop service with process id %d", pid)}
			return
		}

		delete(runningServices, service)
		log.Println(fmt.Sprintf("oCIS service %s has been stopped successfully", service))

		resultChan <- struct {
			success bool
			message string
		}{true, fmt.Sprintf("Service %s stopped successfully", service)}
	}()

	// Wait for the goroutine to finish and retrieve the result.
	stopWg.Wait()
	result := <-resultChan
	close(resultChan)

	return result.success, result.message
}

func WaitUntilPortListens(service string) (bool, error) {
	// Define the overall timeout for the function
	overallTimeout := time.After(30 * time.Second)

	for {
		select {
		case <-overallTimeout:
			// Overall timeout reached, return an error
			errMsg := fmt.Sprintf("Service %s not found in runningServices within the timeout period", service)
			return false, fmt.Errorf(errMsg)

		default:
			pid, exists := runningServices[service]
			if exists {
				// Construct the command to check the port for the PID
				netFilePath := fmt.Sprintf("/proc/%d/net/tcp", pid)

				// Inner timeout for checking the port
				timeoutChan := time.After(30 * time.Second)

				for {
					// Read the /proc/pid/net/tcp file to get the port
					fileContent, err := os.ReadFile(netFilePath)
					if err != nil {
						return false, fmt.Errorf("Error reading %s: %v", netFilePath, err)
					}

					// Process each line in the file
					lines := strings.Split(string(fileContent), "\n")
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if line == "" {
							continue
						}

						// Skip the first line (header)
						if strings.HasPrefix(line, "sl") {
							continue
						}

						// Extract the information from each line
						parts := strings.Fields(line)
						if len(parts) < 10 {
							continue
						}

						// Address field is at index 1
						address := parts[1]

						// Split the address into IP and port parts (IPv4:port format)
						addressParts := strings.Split(address, ":")
						if len(addressParts) < 2 {
							continue
						}

						// The port is the last part (after ":")
						portHex := addressParts[len(addressParts)-1]

						// Convert the hexadecimal port to decimal
						port, err := strconv.ParseInt(portHex, 16, 32)
						if err != nil {
							log.Println(fmt.Sprintf("Error converting port from hex: %v", err))
							continue
						}

						// Ensure we have a valid port
						if port == 0 {
							continue
						}

						log.Println(fmt.Sprintf("Port found for PID %d: %d", pid, port))
						return true, nil
					}

				// If no port found, wait a bit and try again
					select {
					case <-timeoutChan:
						errMsg := fmt.Sprintf("Timeout reached, no port found for PID %d", pid)
						return false, fmt.Errorf(errMsg)
					case <-time.After(2 * time.Second):
					}
				}
			}

			// If service not found, wait a bit and retry
			log.Println(fmt.Sprintf("Service %s not found running, retrying to find service...", service))
			time.Sleep(2 * time.Second)
		}
	}
}
