package TestHelpers

import (
	"fmt"
	"net/http"
	"strings"
)

func CreateUser(
	baseUrl string,
	xRequestId string,
	adminUser string,
	adminPassword string,
	userName string,
	password string,
	email string,
	displayName string,
) *http.Response {
	payload := prepareCreateUserPayload(userName, password, email, displayName)

	url := getFullUrl(baseUrl, "users")

	return post(
		url,
		xRequestId,
		adminUser,
		adminPassword,
		getRequestHeaders(),
		payload,
		nil,
		nil,
		false,
	)
}

func prepareCreateUserPayload(userName string, password string, email string, displayName string) *strings.Reader {
	payload := strings.NewReader(fmt.Sprintf(`{
		"onPremisesSamAccountName": "%s",
		"passwordProfile": {"password": "%s"},
		"displayName": "%s",
		"mail": "%s"
	}`, userName, password, displayName, email))

	return payload
}

func getFullUrl(baseUrl string, endpoint string) string {
	fullUrl := baseUrl
	if !strings.HasSuffix(baseUrl, "/") {
		fullUrl += "/"
	}

	fullUrl += "graph/v1.0/" + endpoint

	return fullUrl
}

func getRequestHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}
