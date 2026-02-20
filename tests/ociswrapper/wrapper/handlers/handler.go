package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"ociswrapper/common"
	"ociswrapper/ocis"
	"ociswrapper/ocis/config"
	"strconv"
	"strings"
)

type BasicResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type CommandResponse struct {
	*BasicResponse
	ExitCode int `json:"exitCode"`
}

func parseJsonBody(reqBody io.ReadCloser) (map[string]any, error) {
	body, _ := io.ReadAll(reqBody)

	if len(body) == 0 || !json.Valid(body) {
		return nil, errors.New("invalid json data")
	}

	var bodyMap map[string]any
	json.Unmarshal(body, &bodyMap)

	return bodyMap, nil
}

func sendResponse(res http.ResponseWriter, statusCode int, message string) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	var status string
	if statusCode == http.StatusOK {
		status = "OK"
	} else {
		status = "ERROR"
	}

	resBody := BasicResponse{
		Status:  status,
		Message: message,
	}

	jsonResponse, _ := json.Marshal(resBody)
	res.Write(jsonResponse)
}

func sendCmdResponse(res http.ResponseWriter, exitCode int, message string) {
	resBody := CommandResponse{
		BasicResponse: &BasicResponse{
			Message: message,
		},
		ExitCode: exitCode,
	}

	if exitCode == 0 {
		resBody.BasicResponse.Status = "OK"
	} else {
		resBody.BasicResponse.Status = "ERROR"
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")

	jsonResponse, _ := json.Marshal(resBody)
	res.Write(jsonResponse)
}

func SetEnvHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		sendResponse(res, http.StatusMethodNotAllowed, "")
		return
	}

	envBody, err := parseJsonBody(req.Body)
	if err != nil {
		sendResponse(res, http.StatusMethodNotAllowed, "Invalid json body")
		return
	}

	var envMap []string
	for key, value := range envBody {
		envMap = append(envMap, fmt.Sprintf("%s=%v", key, value))
	}
	ocis.ServiceEnvConfigs["ocis"] = append(ocis.ServiceEnvConfigs["ocis"], envMap...)

	var message string

	success, _ := ocis.Restart(ocis.ServiceEnvConfigs["ocis"])
	if success {
		message = "oCIS configured successfully"
		sendResponse(res, http.StatusOK, message)
		return
	}

	message = "Failed to restart oCIS with new configuration"
	sendResponse(res, http.StatusInternalServerError, message)
}

func K8sSetEnvHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		sendResponse(res, http.StatusMethodNotAllowed, "")
		return
	}
	envBody, err := parseJsonBody(req.Body)
	if err != nil {
		sendResponse(res, http.StatusMethodNotAllowed, "Invalid json body")
		return
	}

	var message string
	var envMap []string

	for service, envs := range envBody {
		for env, value := range envs.(map[string]any) {
			envMap = append(envMap, fmt.Sprintf("%s=%v", env, value))
		}
		ocis.K3dServiceEnvConfigs[service] = append(ocis.K3dServiceEnvConfigs[service], envMap...)
		success, _ := ocis.UpdateEnv(service, envMap)
		if !success {
			message = "Failed to restart oCIS with new configuration"
			sendResponse(res, http.StatusInternalServerError, message)
			return
		}
		envMap = []string{}
	}
	message = "oCIS configured successfully"
	sendResponse(res, http.StatusOK, message)
}

func RollbackHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		sendResponse(res, http.StatusMethodNotAllowed, "")
		return
	}

	var message string
	ocis.ServiceEnvConfigs["ocis"] = []string{}
	success, _ := ocis.Restart([]string{})
	if success {
		message = "oCIS configuration rolled back successfully"
		sendResponse(res, http.StatusOK, message)
		return
	}

	message = "Failed to restart oCIS with initial configuration"
	sendResponse(res, http.StatusInternalServerError, message)
}

func K8sRollbackHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		sendResponse(res, http.StatusMethodNotAllowed, "")
		return
	}
	var message string
	success, _ := ocis.Rollback()

	if success {
		message = "oCIS configured successfully"
		sendResponse(res, http.StatusOK, message)
		return
	}
	message = "Failed to restart oCIS with new configuration"
	sendResponse(res, http.StatusInternalServerError, message)
}

func StopOcisHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		sendResponse(res, http.StatusMethodNotAllowed, "")
		return
	}

	success, message := ocis.Stop()
	if success {
		sendResponse(res, http.StatusOK, message)
		return
	}

	sendResponse(res, http.StatusInternalServerError, message)
}

func StartOcisHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		sendResponse(res, http.StatusMethodNotAllowed, "")
		return
	}

	if ocis.IsOcisRunning() {
		sendResponse(res, http.StatusConflict, "oCIS server is already running")
		return
	}

	common.Wg.Add(1)
	go ocis.Start(ocis.ServiceEnvConfigs[ocis.OcisServiceName])

	success, message := ocis.WaitForConnection()
	if success {
		sendResponse(res, http.StatusOK, message)
		return
	}

	sendResponse(res, http.StatusInternalServerError, message)
}

func CommandHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		sendResponse(res, http.StatusMethodNotAllowed, "")
		return
	}

	if req.Body == nil {
		sendResponse(res, http.StatusBadRequest, "Body is missing")
		return
	}

	body, err := parseJsonBody(req.Body)
	if err != nil {
		sendResponse(res, http.StatusBadRequest, "Invalid json body")
		return
	}
	if _, ok := body["command"]; !ok {
		sendResponse(res, http.StatusBadRequest, "Command is missing")
		return
	}

	command := body["command"].(string)

	stdIn := []string{}
	if _, ok := body["inputs"]; ok {
		if inputs, ok := body["inputs"].([]interface{}); ok {
			for _, input := range inputs {
				if _, ok := input.(string); ok {
					stdIn = append(stdIn, input.(string))
				} else {
					sendResponse(res, http.StatusBadRequest, "Invalid input data. Expected string")
					return
				}
			}
		}
	}

	exitCode, output := ocis.RunCommand(command, stdIn)
	sendCmdResponse(res, exitCode, output)
}

func OcisServiceHandler(res http.ResponseWriter, req *http.Request) {
	serviceName := req.PathValue("service")

	var envBody map[string]interface{}

	if req.Body != nil && req.ContentLength > 0 {
		var err error
		envBody, err = parseJsonBody(req.Body)
		if err != nil {
			sendResponse(res, http.StatusBadRequest, "Invalid json body")
			return
		}
	}

	if req.Method == http.MethodPost {
		for key, value := range envBody {
			ocis.ServiceEnvConfigs[serviceName] = append(ocis.ServiceEnvConfigs[serviceName], fmt.Sprintf("%s=%v", key, value))
			if strings.HasSuffix(key, "DEBUG_ADDR") {
				address := strings.Split(value.(string), ":")
				port, _ := strconv.Atoi(address[1])
				config.SetServiceDebugPort(serviceName, port)
			}
		}
		log.Println(fmt.Sprintf("Starting '%s' service...", serviceName))

		common.Wg.Add(1)
		go ocis.StartService(serviceName, ocis.ServiceEnvConfigs[serviceName])

		success := ocis.WaitForServiceStatus(serviceName, true, false)
		if success {
			sendResponse(res, http.StatusOK, fmt.Sprintf("'%s' service started successfully", serviceName))
		} else {
			sendResponse(res, http.StatusInternalServerError, fmt.Sprintf("Failed to start '%s' service.", serviceName))
		}
		return
	} else if req.Method == http.MethodDelete {
		success, message := ocis.StopService(serviceName)
		log.Println(message)
		if success {
			sendResponse(res, http.StatusOK, message)
		} else {
			sendResponse(res, http.StatusInternalServerError, message)
		}
		return
	} else if req.Method == http.MethodPatch {
		success, message := ocis.StopService(serviceName)
		if success {
			var serviceEnvMap []string
			for key, value := range envBody {
				serviceEnvMap = append(serviceEnvMap, fmt.Sprintf("%s=%v", key, value))
				if strings.HasSuffix(key, "DEBUG_ADDR") {
					address := strings.Split(value.(string), ":")
					port, _ := strconv.Atoi(address[1])
					config.SetServiceDebugPort(serviceName, port)
				}
				serviceEnvMap = append(ocis.ServiceEnvConfigs[serviceName], serviceEnvMap...)
			}
			common.Wg.Add(1)
			log.Println(fmt.Sprintf("Restarting '%s' service...", serviceName))
			go ocis.StartService(serviceName, serviceEnvMap)

			success := ocis.WaitForServiceStatus(serviceName, true, false)
			if success {
				sendResponse(res, http.StatusOK, fmt.Sprintf("'%s' service updated successfully", serviceName))
				return
			}
		}
		sendResponse(res, http.StatusInternalServerError, message)
		return
	}
	sendResponse(res, http.StatusMethodNotAllowed, "Method Not Allowed")
}

func RollbackServicesHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		sendResponse(res, http.StatusMethodNotAllowed, "")
		return
	}

	var rollbackFailed bool
	var message string

	for serviceName := range ocis.ServiceEnvConfigs {
		if serviceName != "ocis" {
			if success := ocis.WaitForServiceStatus(serviceName, true, false); success {
				success, message := ocis.StopService(serviceName)
				log.Println(message)
				if !success {
					rollbackFailed = true
					break
				}
				delete(ocis.ServiceEnvConfigs, serviceName)
			}
		}
	}

	if rollbackFailed {
		sendResponse(res, http.StatusInternalServerError, message)
	} else {
		sendResponse(res, http.StatusOK, "All services have been rolled back successfully")
	}
}
