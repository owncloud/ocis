package ocis

import (
	"bytes"
	"fmt"
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
	envMap = append([]string{"set", "env", "-n", "ocis", "deployment", service}, envMap...)
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
			return
		case <-tick.C:
			cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl get pods -n ocis -A | grep %s | wc -l", service))
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
						cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl get pods -n ocis -A | grep %s | grep Running | wc -l", service))
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
		cmdArgs := []string{"set", "env", "-n", "ocis"}
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
