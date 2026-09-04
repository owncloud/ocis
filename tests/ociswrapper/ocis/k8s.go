package ocis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"ociswrapper/log"
	"ociswrapper/ocis/config"
	"os/exec"
	"slices"
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

	initialEnvs, err := getInitialEnvs(service)
	if err != nil {
		return false, "error getting existing envs"
	}

	_, ok := K8sOcisInitEnv[service]
	if !ok {
		K8sOcisInitEnv[service] = &ServiceConfig{
			CurrentPod: podName,
			Envs:       initialEnvs,
		}
	} else {
		extraEnvs := diffEnvs(K8sOcisInitEnv[service].Envs, initialEnvs)
		K8sOcisInitEnv[service].Envs = append(K8sOcisInitEnv[service].Envs, extraEnvs...)
		K8sOcisInitEnv[service].CurrentPod = podName
	}

	// envMap may introduce vars that have no prior explicit value on the pod at all (e.g. only
	// a code-level default was in effect) - the tracked baseline, just established/updated
	// above by either branch, has no entry for those, so on its own it is not a rollback target
	// that removes them; kubectl set env is additive and never strips a var it isn't told
	// about. Mark any such brand-new var for removal now, while we still know it is new, or it
	// silently survives every future rollback.
	log.Println(fmt.Sprintf("[%s] DEBUG requested envMap: %s", service, strings.Join(envMap, ", ")))
	log.Println(fmt.Sprintf("[%s] DEBUG baseline before newlyIntroduced check: %s", service, strings.Join(K8sOcisInitEnv[service].Envs, ", ")))
	newlyIntroduced := diffEnvs(K8sOcisInitEnv[service].Envs, envMap)
	log.Println(fmt.Sprintf("[%s] DEBUG newlyIntroduced: %s", service, strings.Join(newlyIntroduced, ", ")))
	K8sOcisInitEnv[service].Envs = append(K8sOcisInitEnv[service].Envs, newlyIntroduced...)
	log.Println(fmt.Sprintf("[%s] DEBUG baseline immediately after append, ptr=%p: %s", service, K8sOcisInitEnv[service], strings.Join(K8sOcisInitEnv[service].Envs, ", ")))

	envSet, skipWaitForService, err := setServiceEnv(service, envMap, "Failed to set env")
	if err != nil {
		return false, "error setting env"
	}

	if !skipWaitForService {
		_, err = waitForService(service, envSet)
		if err != nil {
			return false, "error waiting for service"
		}
	}

	return true, "ok"
}

func diffEnvs(initialEnvsMap []string, currentEnvMap []string) []string {
	extraEnvs := []string{}
	for _, env := range getEnvKeys(currentEnvMap) {
		if !slices.Contains(getEnvKeys(initialEnvsMap), env) {
			extraEnvs = append(extraEnvs, env+"-")
		}
	}
	return extraEnvs
}

func getEnvKeys(envMap []string) []string {
	envKeys := []string{}
	for _, env := range envMap {
		envKey := strings.Split(env, "=")[0]
		envKeys = append(envKeys, envKey)
	}
	return envKeys
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
	// The chart legitimately renders some vars twice in the pod spec (a default value, then a
	// later override for the same name - Go's own os.Environ() takes the last one, which is
	// what the running process actually observes). Deduping here, keeping the last occurrence,
	// ensures this list is safe to replay through a single `kubectl set env` call later (e.g.
	// during rollback): passing the same key twice in one invocation against a spec that
	// already has two entries for it has been observed to drop the variable entirely instead of
	// converging on one value, rather than raising an error.
	return dedupeEnvs(flatEnvVars), nil
}

func dedupeEnvs(envs []string) []string {
	indexByKey := make(map[string]int, len(envs))
	deduped := make([]string, 0, len(envs))
	for _, env := range envs {
		key := strings.SplitN(env, "=", 2)[0]
		if idx, ok := indexByKey[key]; ok {
			deduped[idx] = env
			continue
		}
		indexByKey[key] = len(deduped)
		deduped = append(deduped, env)
	}
	return deduped
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

			output, err := checkServiceHealth(service)
			if err != nil {
				time.Sleep(pollInterval)
				continue
			}

			if strings.Contains(output, "200200") {
				err = checkServiceGrpc(service, podName)
				if err != nil {
					time.Sleep(pollInterval)
					continue
				}

				log.Println(fmt.Sprintf("[%s] Service is healthy and ready. Pod: %s", service, podName))
				return true, nil
			}

			log.Println(fmt.Sprintf("[%s] Waiting for service. Pod: %s. Output: %s", service, podName, output))
			time.Sleep(pollInterval)
		}
	}
}

func setServiceEnv(service string, envMap []string, errMsgPrefix string) (bool, bool, error) {
	// kubectl set env has been observed to behave unreliably (silently dropping the variable
	// entirely, or picking an arbitrary one of the existing values) when the target pod spec
	// already has multiple pre-existing entries for a key being set - which happens here
	// because the chart legitimately renders some vars twice (a default value, then a later
	// vault-mode override, relying on the running process's own last-one-wins env handling,
	// not on kubectl ever reconciling it). Strip any existing entries for every key this call
	// touches first, in its own invocation, so the actual set below always starts from a clean
	// (zero-or-one-entry) state instead of leaving kubectl to reconcile a pre-existing
	// duplicate on its own. Removing a key that isn't set is a no-op, so this is safe to do
	// unconditionally.
	removalArgs := []string{}
	seenKeys := map[string]bool{}
	for _, env := range envMap {
		key := strings.TrimSuffix(strings.SplitN(env, "=", 2)[0], "-")
		if !seenKeys[key] {
			seenKeys[key] = true
			removalArgs = append(removalArgs, key+"-")
		}
	}
	if len(removalArgs) > 0 {
		removeCmdArgs := append([]string{"set", "env", "-n", config.Get("namespace"), "deployment", service}, removalArgs...)
		if _, err := exec.Command("kubectl", removeCmdArgs...).Output(); err != nil {
			errMsg := ""
			if exitErr, ok := err.(*exec.ExitError); ok {
				errMsg = strings.TrimSpace(string(exitErr.Stderr))
			}
			log.Println(fmt.Sprintf("[%s] Failed to pre-remove existing envs before setting them. %s", service, errMsg))
			return false, true, fmt.Errorf("error removing existing env before set")
		}
	}

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
		return false, true, fmt.Errorf("error setting env")
	}
	outString := strings.TrimSpace(string(output))
	if strings.Contains(outString, "env updated") {
		return true, false, nil
	}
	log.Println(fmt.Sprintf("[%s] No change in env. Current pod will be used.", service))
	return true, true, nil
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
		log.Println(fmt.Sprintf("[%s] DEBUG config ptr=%p at top of rollback loop", service, config))
		envs := config.Envs
		log.Println(fmt.Sprintf("[%s] Rolling envs: %s", service, strings.Join(envs, ", ")))
		podName, err := getPodName(service)
		if err != nil {
			return false, "error getting pod name"
		}
		K8sOcisInitEnv[service].CurrentPod = podName
		log.Println(fmt.Sprintf("[%s] Rolling back service. Current Pod: %s", service, podName))

		envSet, skipWaitForService, err := setServiceEnv(service, envs, fmt.Sprintf("Failed to rollback service. Pod: %s", podName))
		if err != nil {
			return false, "failed to rollback"
		}

		if !skipWaitForService {
			_, err = waitForService(service, envSet)			
			if err != nil {
				return false, "error waiting for service"
			}
		}
	}
	return true, "ok"
}
