package ocis

import (
	"fmt"
	"ociswrapper/log"
	"ociswrapper/ocis/config"
	"os/exec"
	"strings"
	"time"
)

var K3dServiceEnvConfigs = make(map[string][]string)

func K8sUpdateEnv(service string, envMap []string) (bool, string) {
	if envMap == nil {
		envMap = []string{}
	}
	K3dServiceEnvConfigs[service] = append(K3dServiceEnvConfigs[service], envMap...)

	cmdArgs := append([]string{"set", "env", "-n", "ocis-server", "deployment", service}, envMap...)
	cmd := exec.Command("kubectl", cmdArgs...)
	_, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
		return false, "error"
	}
	waitForService(service)
	return true, "ok"
}

func waitForService(service string) {
	timeout := time.After(30 * time.Second)
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()

	port := config.GetServiceDebugPort(service)
	healthUrl := fmt.Sprintf("http://%s:%d/healthz", service, port)
	readyUrl := fmt.Sprintf("http://%s:%d/readyz", service, port)
	for {
		select {
		case <-timeout:
			log.Println(fmt.Sprintf("%s seconds timeout waiting for '%s' service.", "30", service))
			return
		case <-tick.C:
			curlCmd := fmt.Sprintf("curl %s -s -o /dev/null -w '%%{http_code}';", healthUrl)
			curlCmd += fmt.Sprintf("curl %s -s -o /dev/null -w '%%{http_code}';echo", readyUrl)
			cmdString := fmt.Sprintf("kubectl run healthcheck -n %s --rm -it --image=curlimages/curl --restart=Never -- sh -c", "ocis-server")
			cmdString += fmt.Sprintf(" \"%s\"", curlCmd)
			cmd := exec.Command("sh", "-c", cmdString)
			stdout, err := cmd.Output()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			output := strings.ReplaceAll(strings.TrimSpace(string(stdout)), "\n", ":")
			if strings.Contains(output, "200200") {
				log.Println(fmt.Sprintf("'%s' service is healthy and ready.", service))
				return
			} else {
				log.Println(fmt.Sprintf("Waiting for '%s' service. Output: %s", service, output))
			}
		}
	}
}

func K8sRollback() (bool, string) {
	for service, envs := range K3dServiceEnvConfigs {
		cmdArgs := []string{"set", "env", "-n", "ocis-server", "deployment", service}
		cmdArgs = append(cmdArgs, envs...)
		cmd := exec.Command("kubectl", cmdArgs...)
		_, err := cmd.Output()
		if err != nil {
			log.Println(fmt.Sprintf("Failed to rollback service '%s'. %s", service, err.Error()))
			return false, "failed to rollback"
		}
		waitForService(service)
		delete(K3dServiceEnvConfigs, service)
	}
	return true, "ok"
}
