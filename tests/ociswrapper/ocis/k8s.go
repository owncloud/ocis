package ocis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ociswrapper/log"
	"ociswrapper/ocis/config"
	"os/exec"
	"strings"
	"time"
)

type ServiceConfig struct {
	CurrentPod string
	Envs       map[string][]string
}

var K8sOcisServices = ServiceConfig{
	Envs: make(map[string][]string),
}

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
		log.Println(err.Error())
		return false, "error getting pod name"
	}
	K8sOcisServices.CurrentPod = podName
	log.Println(fmt.Sprintf("Updating env variables for '%s' service. Current Pod: %s", service, podName))

	if envMap == nil {
		envMap = []string{}
	}
	initialEnvs, err := getInitialEnvs(service)
	if err != nil {
		return false, "error getting existing envs"
	}
	K8sOcisServices.Envs[service] = initialEnvs

	cmdArgs := append([]string{"set", "env", "-n", config.Get("namespace"), "deployment", service}, envMap...)
	cmd := exec.Command("kubectl", cmdArgs...)
	_, err = cmd.Output()
	if err != nil {
		log.Println(err.Error())
		return false, "error"
	}
	_, err = waitForService(service)
	if err != nil {
		log.Println(err.Error())
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
		log.Println(err.Error())
		return nil, err
	}
	output = bytes.TrimSpace(output)
	output = bytes.Trim(output, "\"")

	var flatEnvVars []string
	var allEnvs []EnvVar
	err = json.Unmarshal(output, &allEnvs)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	for _, env := range allEnvs {
		// do not include env vars with valueFrom (includes secrets).
		if env.ValueFrom == nil {
			flatEnvVars = append(flatEnvVars, fmt.Sprintf("%s=%s", env.Name, env.Value))
		}
	}
	return flatEnvVars, nil
}

func waitForService(service string) (bool, error) {
	timeoutInSecond := 30
	timeout := time.After(time.Duration(timeoutInSecond) * time.Second)
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()

	_, err := waitPodDelete(K8sOcisServices.CurrentPod, timeoutInSecond)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(fmt.Sprintf("Old pod '%s' deleted.", K8sOcisServices.CurrentPod))

	port := config.GetServiceDebugPort(service)
	healthUrl := fmt.Sprintf("http://%s:%d/healthz", service, port)
	readyUrl := fmt.Sprintf("http://%s:%d/readyz", service, port)

	log.Println(fmt.Sprintf("Waiting for '%s' service to be ready...", service))

	for {
		select {
		case <-timeout:
			return false, fmt.Errorf("%d seconds timeout waiting for '%s' service.", timeoutInSecond, service)
		case <-tick.C:
			_, err := waitPodReady(service, timeoutInSecond)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			podName, err := getPodName(service)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			curlCmd := fmt.Sprintf("curl %s -s -o /dev/null -w '%%{http_code}';", healthUrl)
			curlCmd += fmt.Sprintf("curl %s -s -o /dev/null -w '%%{http_code}';echo", readyUrl)
			cmdString := fmt.Sprintf("kubectl run healthcheck -n %s --rm -it --image=curlimages/curl --restart=Never -- sh -c", config.Get("namespace"))
			cmdString += fmt.Sprintf(" \"%s\"", curlCmd)
			cmd := exec.Command("sh", "-c", cmdString)
			stdout, err := cmd.Output()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			output := strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", ": ")
			if strings.Contains(output, "200200") {
				log.Println(fmt.Sprintf("'%s' service is healthy and ready. Pod: %s", service, podName))
				return true, nil
			}
			log.Println(fmt.Sprintf("Waiting for '%s' service. Pod: %s. Output: %s", service, podName, output))
		}
	}
}

func getPodName(service string) (string, error) {
	cmdString := fmt.Sprintf("kubectl get pods -n %s -l app=%s -o jsonpath=\"{.items[0].metadata.name}\"", config.Get("namespace"), service)
	cmd := exec.Command("sh", "-c", cmdString)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

func waitPodReady(service string, timeout int) (string, error) {
	cmdString := fmt.Sprintf("kubectl -n %s wait pod --for=condition=Ready -l app=%s --timeout=%ds", config.Get("namespace"), service, timeout)
	cmd := exec.Command("sh", "-c", cmdString)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", ". "), nil
}

func waitPodDelete(podName string, timeout int) (string, error) {
	cmdString := fmt.Sprintf("kubectl -n %s wait pod %s --for=delete --timeout=%ds", config.Get("namespace"), podName, timeout)
	cmd := exec.Command("sh", "-c", cmdString)
	stdout, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", ". "), nil
}

func K8sRollback() (bool, string) {
	for service, envs := range K8sOcisServices.Envs {
		podName, err := getPodName(service)
		if err != nil {
			log.Println(err.Error())
			return false, "error getting pod name"
		}
		K8sOcisServices.CurrentPod = podName
		log.Println(fmt.Sprintf("Rolling back '%s' service. Current Pod: %s", service, podName))

		cmdArgs := []string{"set", "env", "-n", config.Get("namespace"), "deployment", service}
		cmdArgs = append(cmdArgs, envs...)
		cmd := exec.Command("kubectl", cmdArgs...)
		_, err = cmd.Output()
		if err != nil {
			log.Println(fmt.Sprintf("Failed to rollback service '%s'. Pod: %s. %s", service, podName, err.Error()))
			return false, "failed to rollback"
		}
		_, err = waitForService(service)
		if err != nil {
			log.Println(err.Error())
			return false, "error waiting for service"
		}
		delete(K8sOcisServices.Envs, service)
	}
	return true, "ok"
}
