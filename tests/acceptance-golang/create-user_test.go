package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	graphHelper "goAcceptanceTest/test-helpers"
	setupHelper "goAcceptanceTest/test-helpers"

	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
)

var statusCode int
var responseBody string

func adminCreatesUser(username, password, email, displayName string) error {
	response, err := graphHelper.CreateUser(
		setupHelper.GetBaseUrl(),
		"",
		"",
		"",
		username,
		password,
		email,
		displayName,
	)

	if err != nil {
		return err
	}
	statusCode = response.StatusCode
	json.NewDecoder(response.Body).Decode(&responseBody)
	return nil
}

func userExists(username string) bool {
	response, err := graphHelper.GetUser(
		setupHelper.GetBaseUrl(),
		"",
		"",
		"",
		username,
	)
	statusCode = response.StatusCode

	if err != nil {
		return false
	}

	if statusCode >= 400 {
		return false
	}
	return true
}

func theAdministratorCreatesUserUsingTheGraphAPIWithTheFollowingSettings(table *godog.Table) error {

	assist := assistdog.NewDefault()
	userInfo, err := assist.ParseMap(table)
	if err != nil {
		return err
	}
	username := userInfo["userName"]
	password := userInfo["password"]
	email := userInfo["email"]
	displayName := userInfo["displayName"]
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
	if !userExists(username) {
		return fmt.Errorf("User '%s' should exist but does not exist", username)
	}
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
