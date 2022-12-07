package TestHelpers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
)

func post(
	url string,
	xRequestId string,
	user string,
	password string,
	headers map[string]string,
	body *strings.Reader,
	config []string,
	cookies []string,
	stream bool,
) *http.Response {
	response := sendRequest(
		url,
		xRequestId,
		"POST",
		user,
		password,
		headers,
		body,
		config,
		cookies,
		stream,
		0,
		nil)
	return response
}

func sendRequest(
	url string,
	xRequestId string,
	method string,
	user string,
	password string,
	headers map[string]string,
	body *strings.Reader,
	config []string,
	cookies []string,
	stream bool,
	timeout int,
	client *http.Client) *http.Response {
	if client == nil {
		client = createClient(
			user,
			password,
			config,
			cookies,
			stream,
			timeout,
		)
	}

	request := createRequest(
		url,
		xRequestId,
		method,
		headers,
		body,
	)

	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	return response
}

func createClient(
	user string,
	password string,
	config []string,
	cookies []string,
	stream bool,
	timeout int,
) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// ClientAuth
	// options := make(map[string]string)
	client := &http.Client{Transport: tr}
	return client
}

func createRequest(
	url string,
	xRequestId string,
	method string,
	headers map[string]string,
	body *strings.Reader,
) *http.Request {
	request, err := http.NewRequest(method, url, body)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	request.SetBasicAuth(GetAdminUsername(), GetAdminPassword())

	request.Header.Add("Content-Type", headers["Content-Type"])
	request.Header.Add("X-Request-Id", xRequestId)
	return request
}
