<?php

declare(strict_types=1);

/**
 * ownCloud
 *
 * @author Viktor Scharf <v.scharf@owncloud.com>
 * @copyright Copyright (c) 2022 Viktor Scharf v.scharf@owncloud.com
 */

use Behat\Behat\Context\Context;
use GuzzleHttp\Exception\GuzzleException;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use PHPUnit\Framework\Assert;
use Psr\Http\Message\ResponseInterface;

require_once 'bootstrap.php';

/**
 * Context for the TUS-specific steps using the Graph API
 */
class SettingsContext implements Context {
	private FeatureContext $featureContext;
	private SpacesContext $spacesContext;
	private string $baseUrl;
	private string $settingsUrl = '/api/v0/settings/';

	/**
	 * This will run before EVERY scenario.
	 * It will set the properties for this object.
	 *
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 */
	public function before(BeforeScenarioScope $scope): void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context from here
		$this->featureContext = $environment->getContext('FeatureContext');
		$this->spacesContext = $environment->getContext('SpacesContext');
		$this->baseUrl = \trim($this->featureContext->getBaseUrl(), "/");
	}

	/**
	 * @When /^user "([^"]*)" tries to get all existing roles$/
	 *
	 * @param string $user
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function getAllExistingRoles(string $user): void {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "roles-list";
		$this->featureContext->setResponse(
			$this->spacesContext->sendPostRequestToUrl($fullUrl, $user, $this->featureContext->getPasswordForUser($user), "{}", $this->featureContext->getStepLineRef())
		);
	}

	/**
	 * @param string $user
	 * @param string $userId
	 * @param string $roleId
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendRequestToAssignRoleToUser(string $user, string $userId, string $roleId): void {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "assignments-add";
		$body = json_encode(["account_uuid" => $userId, "role_id" => $roleId], JSON_THROW_ON_ERROR);

		$this->featureContext->setResponse(
			$this->spacesContext->sendPostRequestToUrl($fullUrl, $user, $this->featureContext->getPasswordForUser($user), $body, $this->featureContext->getStepLineRef())
		);
	}

	/**
	 * @param string $user
	 * @param string $userId
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendRequestAssignmentsList(string $user, string $userId): ResponseInterface {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "assignments-list";
		$body = json_encode(["account_uuid" => $userId], JSON_THROW_ON_ERROR);
		return $this->spacesContext->sendPostRequestToUrl($fullUrl, $user, $this->featureContext->getPasswordForUser($user), $body, $this->featureContext->getStepLineRef());
	}

	/**
	 * @When /^the administrator has given "([^"]*)" the role "([^"]*)" using the settings api$/
	 *
	 * @param string $user
	 * @param string $role
	 *
	 * @return void
	 *
	 * @throws Exception
	 */
	public function theAdministratorHasGivenUserTheRole(string $user, string $role): void {
		$admin = $this->featureContext->getAdminUserName();
		$roleId = $this->userGetRoleIdByRoleName($admin, $role);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id') ?? $user;
		$this->setRoleToUser($admin, $userId, $roleId);
	}

	/**
	 * @param string $user
	 * @param string $role
	 *
	 * @return string
	 */
	public function userGetRoleIdByRoleName(string $user, string $role): string {
		$this->getAllExistingRoles($user);

		if ($this->featureContext->getResponse()) {
			$rawBody =  $this->featureContext->getResponse()->getBody()->getContents();
			$decodedBody = \json_decode($rawBody, true, 512, JSON_THROW_ON_ERROR);
			Assert::assertArrayHasKey(
				'bundles',
				$decodedBody,
				__METHOD__ . " could not find bundles in body"
			);
			$bundles = $decodedBody["bundles"];
		} else {
			$bundles = [];
		}

		$roleToAssign = "";
		foreach ($bundles as $value) {
			// find the selected role
			if ($value["displayName"] === $role) {
				$roleToAssign = $value;
				break;
			}
		}
		Assert::assertNotEmpty($roleToAssign, "The selected role $role could not be found");
		return $roleToAssign["id"];
	}

	/**
	 * @param string $user
	 * @param string $userId
	 * @param string $roleId
	 *
	 * @return void
	 * @throws Exception
	 */
	public function setRoleToUser(string $user, string $userId, string $roleId): void {
		$this->sendRequestToAssignRoleToUser($user, $userId, $roleId);

		if ($this->featureContext->getResponse()) {
			$rawBody = $this->featureContext->getResponse()->getBody()->getContents();
			$decodedBody = \json_decode($rawBody, true, 512, JSON_THROW_ON_ERROR);
			Assert::assertArrayHasKey(
				'assignment',
				$decodedBody,
				__METHOD__ . " could not find assignment in body"
			);
			$assignment = $decodedBody["assignment"];
		} else {
			$assignment = [];
		}

		Assert::assertEquals($userId, $assignment["accountUuid"]);
		Assert::assertEquals($roleId, $assignment["roleId"]);
	}

	/**
	 * @When /^user "([^"]*)" changes his own role to "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $role
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userChangeOwnRole(string $user, string $role): void {
		// we assume that the user knows uuid role.
		$roleId = $this->userGetRoleIdByRoleName($this->featureContext->getAdminUserName(), $role);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id');
		$this->sendRequestToAssignRoleToUser($user, $userId, $roleId);
	}

	/**
	 * @When /^user "([^"]*)" changes the role "([^"]*)" for user "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $role
	 * @param string $assignedUser
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userChangeRoleAnotherUser(string $user, string $role, string $assignedUser): void {
		// we assume that the user knows uuid role.
		$roleId = $this->userGetRoleIdByRoleName($this->featureContext->getAdminUserName(), $role);
		$userId = $this->featureContext->getAttributeOfCreatedUser($assignedUser, 'id');
		$this->sendRequestToAssignRoleToUser($user, $userId, $roleId);
	}

	/**
	 * @When /^user "([^"]*)" tries to get list of assignment$/
	 *
	 * @param string $user
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function userGetAssignmentsList(string $user): void {
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id');
		$this->featureContext->setResponse($this->sendRequestAssignmentsList($user, $userId));
	}

	/**
	 * @Then /^user "([^"]*)" should have the role "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $role
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function userShouldHaveRole(string $user, string $role): void {
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id');
		$response = $this->sendRequestAssignmentsList($this->featureContext->getAdminUserName(), $userId);
		$assignmentResponse = $this->featureContext->getJsonDecodedResponseBodyContent($response);
		if (isset($assignmentResponse->assignments[0]->roleId)) {
			$actualRoleId = $assignmentResponse->assignments[0]->roleId;
			Assert::assertEquals($this->userGetRoleIdByRoleName($this->featureContext->getAdminUserName(), $role), $actualRoleId, "user $user has no role $role");
		} else {
			Assert::fail("Response should contain user role but not found.\n" . json_encode($assignmentResponse));
		}
	}

	/**
	 * @Then /^the setting API response should have the role "([^"]*)"$/
	 *
	 * @param string $role
	 *
	 * @return void
	 *
	 * @throws Exception
	 */
	public function theSettingApiResponseShouldHaveTheRole(string $role): void {
		$assignmentRoleId = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse())["assignments"][0]["roleId"];
		Assert::assertEquals($this->userGetRoleIdByRoleName($this->featureContext->getAdminUserName(), $role), $assignmentRoleId, "user has no role $role");
	}

	/**
	 * @param string $user
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendRequestGetBundlesList(string $user): void {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "bundles-list";
		$this->featureContext->setResponse(
			$this->spacesContext->sendPostRequestToUrl($fullUrl, $user, $this->featureContext->getPasswordForUser($user), '{}', $this->featureContext->getStepLineRef())
		);

		$this->featureContext->theHTTPStatusCodeShouldBe(
			201,
			"Expected response status code should be 201"
		);
	}

	/**
	 * @param string $user
	 * @param string $bundleName
	 *
	 * @return array
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function getBundlesList(string $user, string $bundleName): array {
		$this->sendRequestGetBundlesList($user);
		$body = json_decode((string)$this->featureContext->getResponse()->getBody(), true, 512, JSON_THROW_ON_ERROR);
		foreach ($body["bundles"] as $value) {
			if ($value["displayName"] === $bundleName) {
				return $value;
			}
		}
		return [];
	}

	/**
	 * @param string $user
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendRequestGetSettingsValuesList(string $user): void {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "values-list";
		$body = json_encode(["account_uuid" => "me"], JSON_THROW_ON_ERROR);
		$this->featureContext->setResponse(
			$this->spacesContext->sendPostRequestToUrl($fullUrl, $user, $this->featureContext->getPasswordForUser($user), $body, $this->featureContext->getStepLineRef())
		);

		Assert::assertEquals(
			$this->featureContext->getResponse()->getStatusCode(),
			201,
			"Expected response status code should be 201"
		);
	}

	/**
	 * @param string $user
	 *
	 * @return string
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function getSettingLanguageValue(string $user): string {
		$this->sendRequestGetSettingsValuesList($user);
		$body = json_decode((string)$this->featureContext->getResponse()->getBody(), true, 512, JSON_THROW_ON_ERROR);

		// if no language is set, the request body is empty return English as the default language
		if (empty($body)) {
			return "en";
		}
		foreach ($body["values"] as $value) {
			if ($value["identifier"]["setting"] === "language") {
				return $value["value"]["listValue"]["values"][0]["stringValue"];
			}
		}
	}

	/**
	 * @param string $user
	 * @param string $language
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendRequestToSwitchSystemLanguage(string $user, string $language): ResponseInterface {
		$profileBundlesList = $this->getBundlesList($user, "Profile");
		Assert::assertNotEmpty($profileBundlesList, "bundles list is empty");

		$settingId = '';
		foreach ($profileBundlesList["settings"] as $value) {
			if ($value["name"] === "language") {
				$settingId = $value["id"];
				break;
			}
		}
		Assert::assertNotEmpty($settingId, "settingId is empty");

		$fullUrl = $this->baseUrl . $this->settingsUrl . "values-save";
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id');
		$body = json_encode(
			[
			"value" => [
			"account_uuid" => "me",
			"bundleId" => $profileBundlesList["id"],
			"id" => $userId,
			"listValue" => [
			"values" => [
			  [
				"stringValue" => $language
			  ]
			]
			],
			"resource" => [
			"type" => "TYPE_USER"
			],
			"settingId" => $settingId
			]
			],
			JSON_THROW_ON_ERROR
		);
		return $this->spacesContext->sendPostRequestToUrl(
			$fullUrl,
			$user,
			$this->featureContext->getPasswordForUser($user),
			$body,
			$this->featureContext->getStepLineRef()
		);
	}

	/**
	 * @Given /^user "([^"]*)" has switched the system language to "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $language
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function theUserHasSwitchedSystemLanguage(string $user, string $language): void {
		$response = $this->sendRequestToSwitchSystemLanguage($user, $language);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			201,
			"Expected response status code should be 201",
			$response
		);
	}
}
