package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"ociswrapper/ocis"
)

func parseJsonBody(reqBody io.ReadCloser) (map[string]any, error) {
	body, _ := io.ReadAll(reqBody)

	if len(body) == 0 || !json.Valid(body) {
		return nil, errors.New("Invalid json data")
	}

	var bodyMap map[string]any
	json.Unmarshal(body, &bodyMap)

	return bodyMap, nil
}

func sendResponse(res http.ResponseWriter, ocisStatus bool) {
	resBody := make(map[string]string)

	if ocisStatus {
		res.WriteHeader(http.StatusOK)
		resBody["status"] = "OK"
		resBody["message"] = "oCIS server is running"
	} else {
		res.WriteHeader(http.StatusInternalServerError)
		resBody["status"] = "ERROR"
		resBody["message"] = "oCIS server error"
	}
	res.Header().Set("Content-Type", "application/json")

	jsonResponse, _ := json.Marshal(resBody)
	res.Write(jsonResponse)
}

func SetEnvHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	environments, err := parseJsonBody(req.Body)
	if err != nil {
		http.Error(res, "Bad request", http.StatusBadRequest)
		return
	}

	ocisStatus := ocis.Restart(environments)

	sendResponse(res, ocisStatus)
}

func RollbackHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ocisStatus := ocis.Restart(nil)

	sendResponse(res, ocisStatus)
}
