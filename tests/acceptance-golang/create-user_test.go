package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	graphHelper "goAcceptanceTest/test-helpers"
	setupHelper "goAcceptanceTest/test-helpers"

	"github.com/cucumber/godog"
)

var statusCode int
var responseBody string

func adminCreatesUser(username, password, email, displayName string) error {
	response := graphHelper.CreateUser(
		setupHelper.GetBaseUrl(),
		"",
		"",
		"",
		username,
		password,
		email,
		displayName,
	)

	statusCode = response.StatusCode
	json.NewDecoder(response.Body).Decode(&responseBody)
	return nil
}

func theAdministratorCreatesUserUsingTheGraphAPIWithTheFollowingSettings() error {
	username := "alice"
	password := "123456"
	email := "alice@example.com"
	displayName := "Alice Hansen"
	adminCreatesUser(username, password, email, displayName)
	return nil
}

func theHTTPStatusCodeShouldBe(expectedStatusCode string) error {
	actualStatusCode := strconv.Itoa(statusCode)
	if expectedStatusCode != actualStatusCode {
		return fmt.Errorf("expected status code %s, got %s", expectedStatusCode, actualStatusCode)
	}
	return nil
}

func userShouldExist(username string) error {
	// TODO Implement
	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fail()
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the administrator creates user using the Graph API with the following settings:$`, theAdministratorCreatesUserUsingTheGraphAPIWithTheFollowingSettings)
	ctx.Step(`^the HTTP status code should be "([^"]*)"$`, theHTTPStatusCodeShouldBe)
	ctx.Step(`^user "([^"]*)" should exist$`, userShouldExist)
}
