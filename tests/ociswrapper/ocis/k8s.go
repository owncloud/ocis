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
	newlyIntroduced := diffEnvs(K8sOcisInitEnv[service].Envs, envMap)
	K8sOcisInitEnv[service].Envs = append(K8sOcisInitEnv[service].Envs, newlyIntroduced...)

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

func fetchRawEnvVars(service string) ([]EnvVar, error) {
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
		log.Println(fmt.Sprintf("[%s] Failed to get envs. %s", service, errMsg))
		return nil, err
	}
	output = bytes.TrimSpace(output)
	output = bytes.Trim(output, "\"")

	var allEnvs []EnvVar
	err = json.Unmarshal(output, &allEnvs)
	if err != nil {
		log.Println(fmt.Sprintf("[%s] Failed to parse envs. %s", service, err.Error()))
		return nil, err
	}
	return allEnvs, nil
}

func fetchRawEnvs(service string) ([]string, error) {
	allEnvs, err := fetchRawEnvVars(service)
	if err != nil {
		return nil, err
	}

	var flatEnvVars []string
	for _, env := range allEnvs {
		// do not include env vars with valueFrom (includes secrets).
		if env.ValueFrom == nil && env.Value != "" {
			flatEnvVars = append(flatEnvVars, fmt.Sprintf("%s=%s", env.Name, env.Value))
		}
	}
	return flatEnvVars, nil
}

// countEnvNames counts every entry the live pod spec has per env var name, including
// empty-valued and valueFrom entries that fetchRawEnvs filters out for baseline-tracking
// purposes. A name occupying more than one slot is exactly the situation kubectl set env
// cannot reliably converge on its own (see setServiceEnv) - filtering by value here would
// hide a duplicate that happens to have an empty default alongside a real override.
func countEnvNames(service string) (map[string]int, error) {
	allEnvs, err := fetchRawEnvVars(service)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int, len(allEnvs))
	for _, env := range allEnvs {
		counts[env.Name]++
	}
	return counts, nil
}

func getInitialEnvs(service string) ([]string, error) {
	flatEnvVars, err := fetchRawEnvs(service)
	if err != nil {
		return nil, err
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

// duplicatedKeys returns the subset of the given keys that currently occupy more than one
// slot in the live pod spec (per counts, from countEnvNames). Only these need the
// pre-removal step in setServiceEnv - most keys have a single entry and removing+re-setting
// them would trigger two rolling restarts (one from the removal, one from the set) instead
// of one, for no benefit.
func duplicatedKeys(counts map[string]int, keys []string) map[string]bool {
	duplicated := make(map[string]bool)
	for _, key := range keys {
		if counts[key] > 1 {
			duplicated[key] = true
		}
	}
	return duplicated
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
	// kubectl set env can behave unreliably when a key has multiple entries in the pod spec,
	// which can happen when the chart renders a default and a vault-mode override. Remoce only
	// confirmed duplicate keys first to avoid unnecessary rolling restarts. Count all entries,
	// including empty/vauleFrom ones, since kubectl considers them when resolving the key.
	counts, err := countEnvNames(service)
	if err != nil {
		log.Println(fmt.Sprintf("[%s] Could not check for duplicate env entries, proceeding without pre-removal: %s", service, err.Error()))
	}
	touchedKeys := []string{}
	seenKeys := map[string]bool{}
	for _, env := range envMap {
		// envMap entries are either "KEY=VALUE" or a "KEY-" removal marker.
		key := strings.TrimSuffix(strings.SplitN(env, "=", 2)[0], "-")
		if !seenKeys[key] {
			seenKeys[key] = true
			touchedKeys = append(touchedKeys, key)
		}
	}
	duplicated := duplicatedKeys(counts, touchedKeys)

	removalArgs := []string{}
	for key := range duplicated {
		removalArgs = append(removalArgs, key+"-")
	}
	if len(removalArgs) > 0 {
		removeCmdArgs := append([]string{"set", "env", "-n", config.Get("namespace"), "deployment", service}, removalArgs...)
		if _, err := exec.Command("kubectl", removeCmdArgs...).Output(); err != nil {
			errMsg := ""
			if exitErr, ok := err.(*exec.ExitError); ok {
				errMsg = strings.TrimSpace(string(exitErr.Stderr))
			}
			log.Println(fmt.Sprintf("[%s] Pre-remove of duplicated keys before set failed: %s", service, errMsg))
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
		// A var added by a test (not present in the original baseline) would
		// otherwise never get unset: `kubectl set env` only sets the vars it's
		// given, it doesn't remove anything else already on the deployment.
		currentEnvs, err := getInitialEnvs(service)
		if err != nil {
			return false, "error getting current envs"
		}
		extraEnvs := diffEnvs(config.Envs, currentEnvs)
		envs := append(append([]string{}, config.Envs...), extraEnvs...)
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
