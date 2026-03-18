package ocis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"ociswrapper/log"
	"ociswrapper/ocis/config"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
)

type ServiceConfig struct {
	CurrentPod string
	Envs       []string
}

var K8sOcisInitEnv = make(map[string]*ServiceConfig)

type EnvVar struct {
	Name      string `json:"name"`
	Value     string `json:"value,omitempty"`
	ValueFrom *struct {
		SecretKeyRef struct {
			Name string `json:"name"`
			Key  string `json:"key"`
		} `json:"secretKeyRef"`
	} `json:"valueFrom,omitempty"`
}

func K8sUpdateEnv(service string, envMap []string) (bool, string) {
	podName, err := getPodName(service)
	if err != nil {
		return false, "error getting pod name"
	}
	log.Println(fmt.Sprintf("[%s] Updating env variables. Current Pod: %s", service, podName))

	if envMap == nil {
		envMap = []string{}
	}

	_, ok := K8sOcisInitEnv[service]
	if !ok {
		initialEnvs, err := getInitialEnvs(service)
		if err != nil {
			return false, "error getting existing envs"
		}
		K8sOcisInitEnv[service] = &ServiceConfig{
			CurrentPod: podName,
			Envs:       initialEnvs,
		}
	}
	K8sOcisInitEnv[service].CurrentPod = podName

	envSet, err := setServiceEnv(service, envMap, "Failed to set env")
	if err != nil {
		return false, "error setting env"
	}

	_, err = waitForService(service, envSet)
	if err != nil {
		return false, "error waiting for service"
	}
	return true, "ok"
}

func getInitialEnvs(service string) ([]string, error) {
	filter := "jsonpath=\"{.spec.template.spec.containers[*].env}\""
	cmdArgs := []string{"get", "-n", config.Get("namespace"), "deployment", service, "-o", filter}
	cmd := exec.Command("kubectl", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		errMsg := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			// stderr from the command
			errMsg = strings.TrimSpace(string(exitErr.Stderr))
		}
		log.Println(fmt.Sprintf("[%s] Failed to get initial envs. %s", service, errMsg))
		return nil, err
	}
	output = bytes.TrimSpace(output)
	output = bytes.Trim(output, "\"")

	var flatEnvVars []string
	var allEnvs []EnvVar
	err = json.Unmarshal(output, &allEnvs)
	if err != nil {
		log.Println(fmt.Sprintf("[%s] Failed to parse envs. %s", service, err.Error()))
		return nil, err
	}

	for _, env := range allEnvs {
		// do not include env vars with valueFrom (includes secrets).
		if env.ValueFrom == nil && env.Value != "" {
			flatEnvVars = append(flatEnvVars, fmt.Sprintf("%s=%s", env.Name, env.Value))
		}
	}
	return flatEnvVars, nil
}

func waitForService(service string, waitDeletion bool) (bool, error) {
	timeoutInSecond := 30
	timeout := time.After(time.Duration(timeoutInSecond) * time.Second)
	pollInterval := 5 * time.Second

	if waitDeletion {
		_, err := waitPodDelete(K8sOcisInitEnv[service].CurrentPod, timeoutInSecond)
		if err != nil {
			return false, fmt.Errorf("[%s] Pod not deleted", service)
		}
		log.Println(fmt.Sprintf("[%s] Old pod '%s' deleted.", service, K8sOcisInitEnv[service].CurrentPod))
	}
	log.Println(fmt.Sprintf("[%s] Waiting for service to be ready...", service))

	for {
		select {
		case <-timeout:
			log.Println(fmt.Sprintf("[%s] %d seconds timeout waiting service.", service, timeoutInSecond))
			return false, fmt.Errorf("timeout waiting for service")
		default:
			_, err := waitPodReady(service, timeoutInSecond)
			if err != nil {
				time.Sleep(pollInterval)
				continue
			}

			podName, err := getPodName(service)
			if err != nil {
				time.Sleep(pollInterval)
				continue
			}

			err = checkServiceGrpc(service, podName)
			if err != nil {
				time.Sleep(pollInterval)
				continue
			}

			output, err := checkServiceHealth(service)
			if err != nil {
				time.Sleep(pollInterval)
				continue
			}

			if strings.Contains(output, "200200") {
				log.Println(fmt.Sprintf("[%s] Service is healthy and ready. Pod: %s", service, podName))
				return true, nil
			}

			log.Println(fmt.Sprintf("[%s] Waiting for service. Pod: %s. Output: %s", service, podName, output))
			time.Sleep(pollInterval)
		}
	}
}

func setServiceEnv(service string, envMap []string, errMsgPrefix string) (bool, error) {
	cmdArgs := append([]string{"set", "env", "-n", config.Get("namespace"), "deployment", service}, envMap...)
	cmd := exec.Command("kubectl", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		errMsg := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			// stderr from the command
			errMsg = strings.TrimSpace(string(exitErr.Stderr))
		}
		log.Println(fmt.Sprintf("[%s] %s. %s", service, errMsgPrefix, errMsg))
		return false, fmt.Errorf("error setting env")
	}
	outString := strings.TrimSpace(string(output))
	if strings.Contains(outString, "env updated") {
		return true, nil
	}
	log.Println(fmt.Sprintf("[%s] No change in env. Current pod will be used.", service))
	return false, nil
}

func checkServiceGrpc(service string, podName string) error {
	grpcPort := config.GetServiceGRPCPort(service)
	if grpcPort == 0 {
		return nil
	}

	checkCmd := fmt.Sprintf("-plaintext -max-time 1 %s:%d list", service, grpcPort)
	cmdString := fmt.Sprintf(
		"run grpccheck -n %s --rm --attach --image=fullstorydev/grpcurl --restart=Never -- %s",
		config.Get("namespace"),
		checkCmd,
	)
	cmdArgs := strings.Split(cmdString, " ")
	c := exec.Command("kubectl", cmdArgs...)

	// Start the command with a pty (pseudo terminal)
	// This is required by grpc connection
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

	// wait for copy to finish
	<-done
	cmdOutput := output.String()
	cmdOutput = strings.ReplaceAll(strings.TrimSpace(string(cmdOutput)), "\n", ". ")
	if strings.Contains(cmdOutput, "reflection API") {
		log.Println(fmt.Sprintf("[%s] gRPC service is ready. Pod: %s", service, podName))
		return nil
	}
	log.Println(fmt.Sprintf("[%s] gRPC service is not reachable. Pod: %s. Output: %s", service, podName, cmdOutput))
	return fmt.Errorf("gRPC service not reachable")
}

func checkServiceHealth(service string) (string, error) {
	port := config.GetServiceDebugPort(service)
	if port == 0 {
		log.Println(fmt.Sprintf("[%s] Debug port not found", service))
		return "", fmt.Errorf("invalid debug port")
	}
	healthUrl := fmt.Sprintf("http://%s:%d/healthz", service, port)
	readyUrl := fmt.Sprintf("http://%s:%d/readyz", service, port)

	curlCmd := fmt.Sprintf("curl %s -s -o /dev/null -w '%%{http_code}';", healthUrl)
	curlCmd += fmt.Sprintf("curl %s -s -o /dev/null -w '%%{http_code}';echo", readyUrl)
	cmdString := fmt.Sprintf("kubectl run healthcheck -n %s --rm -it --image=curlimages/curl --restart=Never -- sh -c", config.Get("namespace"))
	cmdString += fmt.Sprintf(" \"%s\"", curlCmd)

	cmd := exec.Command("sh", "-c", cmdString)
	stdout, err := cmd.Output()
	if err != nil {
		errMsg := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			// stderr from the command
			errMsg = strings.TrimSpace(string(exitErr.Stderr))
		}
		log.Println(fmt.Sprintf("[%s] Failed to run health check. %s", service, errMsg))
		return "", err
	}
	output := strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", ". ")
	return output, nil
}

func getPodName(service string) (string, error) {
	cmdString := fmt.Sprintf("kubectl get pods -n %s -l app=%s -o jsonpath=\"{.items[0].metadata.name}\"", config.Get("namespace"), service)
	cmd := exec.Command("sh", "-c", cmdString)
	stdout, err := cmd.Output()
	if err != nil {
		errMsg := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			// stderr from the command
			errMsg = strings.TrimSpace(string(exitErr.Stderr))
		}
		log.Println(fmt.Sprintf("[%s] Failed to get pod name. %s", service, errMsg))

		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

func waitPodReady(service string, timeout int) (string, error) {
	cmdString := fmt.Sprintf("kubectl -n %s wait pod --for=condition=Ready -l app=%s --timeout=%ds", config.Get("namespace"), service, timeout)
	cmd := exec.Command("sh", "-c", cmdString)
	stdout, err := cmd.Output()
	if err != nil {
		errMsg := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			// stderr from the command
			errMsg = strings.TrimSpace(string(exitErr.Stderr))
		}
		log.Println(fmt.Sprintf("[%s] Pod not in ready state. %s", service, errMsg))
		return "", err
	}
	return strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", ". "), nil
}

func waitPodDelete(podName string, timeout int) (string, error) {
	cmdString := fmt.Sprintf("kubectl -n %s wait pod %s --for=delete --timeout=%ds", config.Get("namespace"), podName, timeout)
	cmd := exec.Command("sh", "-c", cmdString)
	stdout, err := cmd.Output()
	if err != nil {
		errMsg := ""
		if exitErr, ok := err.(*exec.ExitError); ok {
			// stderr from the command
			errMsg = strings.TrimSpace(string(exitErr.Stderr))
		}
		log.Println(fmt.Sprintf("Pod '%s' not deleted. %s", podName, errMsg))
		return "", err
	}
	return strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", ". "), nil
}

func K8sRollback() (bool, string) {
	for service, config := range K8sOcisInitEnv {
		envs := config.Envs
		log.Println(fmt.Sprintf("[%s] Rolling envs: %s", service, strings.Join(envs, ", ")))
		podName, err := getPodName(service)
		if err != nil {
			return false, "error getting pod name"
		}
		K8sOcisInitEnv[service].CurrentPod = podName
		log.Println(fmt.Sprintf("[%s] Rolling back service. Current Pod: %s", service, podName))

		envSet, err := setServiceEnv(service, envs, fmt.Sprintf("Failed to rollback service. Pod: %s", podName))
		if err != nil {
			return false, "failed to rollback"
		}

		_, err = waitForService(service, envSet)
		if err != nil {
			return false, "error waiting for service"
		}
	}
	return true, "ok"
}
