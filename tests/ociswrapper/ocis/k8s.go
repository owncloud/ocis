package ocis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ociswrapper/log"
	"ociswrapper/ocis/config"
	"os/exec"
	"slices"
	"strings"
	"time"
)

var K3dServiceEnvConfigs = make(map[string][]string)

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
	log.Println(fmt.Sprintf("Updating environment variables for service '%s'...", service))
	if envMap == nil {
		envMap = []string{}
	}
	initialEnvs, err := getInitialEnvs(service, getEnvKeys(envMap))
	if err != nil {
		return false, "error getting existing envs"
	}
	K3dServiceEnvConfigs[service] = append(K3dServiceEnvConfigs[service], initialEnvs...)

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
		return false, "service is ready"
	}
	return true, "ok"
}

func getInitialEnvs(service string, filterEnvKeys []string) ([]string, error) {
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

	var filteredEnvVars []string
	var allEnvs []EnvVar
	err = json.Unmarshal(output, &allEnvs)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	for _, env := range allEnvs {
		envName := env.Name
		envValue := ""
		if !slices.Contains(filterEnvKeys, envName) {
			continue
		}
		// if env has 'valueFrom' field (used for secrets), use the empty string value
		if env.ValueFrom == nil {
			envValue = env.Value
		}
		filteredEnvVars = append(filteredEnvVars, fmt.Sprintf("%s=%s", envName, envValue))
	}
	for _, envKey := range filterEnvKeys {
		if !slices.Contains(getEnvKeys(filteredEnvVars), envKey) {
			filteredEnvVars = append(filteredEnvVars, fmt.Sprintf("%s=", envKey))
		}
	}

	return filteredEnvVars, nil
}

func getEnvKeys(envs []string) []string {
	var keys []string
	for _, env := range envs {
		parts := strings.SplitN(env, "=", 2)
		keys = append(keys, strings.TrimSpace(parts[0]))
	}
	return keys
}

func waitForService(service string) (bool, error) {
	timeout := time.After(30 * time.Second)
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()

	port := config.GetServiceDebugPort(service)
	healthUrl := fmt.Sprintf("http://%s:%d/healthz", service, port)
	readyUrl := fmt.Sprintf("http://%s:%d/readyz", service, port)
	for {
		select {
		case <-timeout:
			return false, fmt.Errorf("%s seconds timeout waiting for '%s' service.", "30", service)
		case <-tick.C:
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
				log.Println(fmt.Sprintf("'%s' service is healthy and ready.", service))
				return true, nil
			} else {
				log.Println(fmt.Sprintf("Waiting for '%s' service. Output: %s", service, output))
			}
		}
	}
}

func K8sRollback() (bool, string) {
	for service, envs := range K3dServiceEnvConfigs {
		log.Println(fmt.Sprintf("Rolling back '%s' service...", service))
		cmdArgs := []string{"set", "env", "-n", config.Get("namespace"), "deployment", service}
		cmdArgs = append(cmdArgs, envs...)
		cmd := exec.Command("kubectl", cmdArgs...)
		_, err := cmd.Output()
		if err != nil {
			log.Println(fmt.Sprintf("Failed to rollback service '%s'. %s", service, err.Error()))
			return false, "failed to rollback"
		}
		_, err = waitForService(service)
		if err != nil {
			log.Println(err.Error())
			return false, "service is ready"
		}
		delete(K3dServiceEnvConfigs, service)
	}
	return true, "ok"
}
