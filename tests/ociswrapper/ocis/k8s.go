package ocis

import (
	"bytes"
	"fmt"
	"ociswrapper/ocis/config"
	"os/exec"
	"strings"
	"time"
)

var K3dServiceEnvConfigs = make(map[string][]string)

func UpdateEnv(service string, envMap []string) (bool, string) {
	if envMap == nil {
		envMap = []string{}
	}

	cmdArgs := new(bytes.Buffer)
	for _, value := range envMap {
		fmt.Fprintf(cmdArgs, "%s ", value)
	}
	envMap = append([]string{"set", "env", "-n", config.Get("namespace"), "deployment", service}, envMap...)
	cmd = exec.Command("kubectl", envMap...)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return false, "error"
	}
	IsServiceRunning(service)
	return true, "ok"
}

func IsServiceRunning(service string) {
	timeout := time.After(10 * time.Second)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-timeout:
			fmt.Printf("Timeout: %s service did not become ready in time.\n",service)
			return
		case <-tick.C:
			cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl get pods -n %s -A | grep %s | wc -l", config.Get("namespace"),service))
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if strings.TrimSpace(string(stdout)) == "1" {
				for {
					select {
					case <-timeout:
						fmt.Println("Timeout: service did not reach 'Running' state in time.")
						return
					case <-tick.C:
						cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl get pods -n %s -A | grep %s | grep Running | wc -l", config.Get("namespace"),service))
						stdout, err := cmd.Output()
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						if strings.TrimSpace(string(stdout)) == "1" {
							return
						}
					}
				}
			}
		}
	}
}

func Rollback() (bool, string) {
	for service, envs := range K3dServiceEnvConfigs {
		cmdArgs := []string{"set", "env", "-n", config.Get("namespace")}
		cmdArgs = append(cmdArgs, fmt.Sprintf("deployment/%s",service))
		for _, env := range envs {
			cmdArgs = append(cmdArgs, strings.SplitN(env, "=", 2)[0]+"-")
			cmd = exec.Command("kubectl", cmdArgs...)
			_, err := cmd.Output()
			if err != nil {
				return false, "service didnt restart"
			}
			IsServiceRunning(service)
			delete(K3dServiceEnvConfigs, service)
		}
	}
	return true, "ok"
}