package TestHelpers

import (
	"fmt"
	"net/http"
	"strings"
)

func GetUser(
	baseUrl string,
	xRequestId string,
	adminUser string,
	adminPassword string,
	userName string,
) (*http.Response, error) {
	url := getFullUrl(baseUrl, "users/"+userName)
	payload := strings.NewReader(``)
	response, err := get(
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

	if err != nil {
		return nil, err

	} else {
		return response, nil
	}

}

func CreateUser(
	baseUrl string,
	xRequestId string,
	adminUser string,
	adminPassword string,
	userName string,
	password string,
	email string,
	displayName string,
) (*http.Response, error) {
	payload := prepareCreateUserPayload(userName, password, email, displayName)

	url := getFullUrl(baseUrl, "users")

	response, err := post(
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
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
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
