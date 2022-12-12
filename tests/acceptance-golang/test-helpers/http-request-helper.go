package TestHelpers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
)

func get(
	url string,
	xRequestId string,
	user string,
	password string,
	headers map[string]string,
	body *strings.Reader,
	config []string,
	cookies []string,
	stream bool,
) (*http.Response, error) {
	response, err := sendRequest(
		url,
		xRequestId,
		"GET",
		user,
		password,
		headers,
		body,
		config,
		cookies,
		stream,
		0,
		nil)
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

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
) (*http.Response, error) {
	response, err := sendRequest(
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
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
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
	client *http.Client) (*http.Response, error) {
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
		return nil, err
	}

	defer response.Body.Close()
	return response, nil
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
	// options for the client go here
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

	if headers == nil {
		headers = make(map[string]string)
	} else {
		request.Header.Add("Content-Type", headers["Content-Type"])
	}
	if xRequestId != "" {
		request.Header.Add("X-Request-Id", xRequestId)
	}
	return request
}
