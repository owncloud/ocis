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
use TestHelpers\HttpRequestHelper;
use Behat\Gherkin\Node\TableNode;

require_once 'bootstrap.php';

/**
 * Context for the TUS-specific steps using the Graph API
 */
class SettingsContext implements Context {
	private FeatureContext $featureContext;
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
		$this->baseUrl = \trim($this->featureContext->getBaseUrl(), "/");
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function getRoles(string $user): ResponseInterface {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "roles-list";
		return HttpRequestHelper::post(
			$fullUrl,
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			null,
			"{}"
		);
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
		$response = $this->getRoles($user);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @param string $user
	 * @param string $userId
	 * @param string $roleId
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function assignRoleToUser(string $user, string $userId, string $roleId): ResponseInterface {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "assignments-add";
		$body = json_encode(["account_uuid" => $userId, "role_id" => $roleId], JSON_THROW_ON_ERROR);
		return HttpRequestHelper::post(
			$fullUrl,
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			null,
			$body
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
	public function getAssignmentsList(string $user, string $userId): ResponseInterface {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "assignments-list";
		$body = json_encode(["account_uuid" => $userId], JSON_THROW_ON_ERROR);
		return HttpRequestHelper::post(
			$fullUrl,
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			null,
			$body
		);
	}

	/**
	 * @Given /^the administrator has given "([^"]*)" the role "([^"]*)" using the settings api$/
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
		$roleId = $this->getRoleIdByRoleName($admin, $role);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id') ?? $user;
		$response = $this->assignRoleToUser($admin, $userId, $roleId);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			201,
			"Expected response status code should be 201",
			$response
		);
	}

	/**
	 * @param string $user
	 * @param string $role
	 *
	 * @return string
	 */
	public function getRoleIdByRoleName(string $user, string $role): string {
		// Sometimes the response body is not complete and results invalid json.
		// So we try again until we get a valid json.
		$retried = 0;
		do {
			$response = $this->getRoles($user);
			$this->featureContext->theHTTPStatusCodeShouldBe(
				201,
				"Expected response status code should be 201",
				$response
			);

			$rawBody =  $response->getBody()->getContents();
			try {
				$decodedBody = \json_decode($rawBody, true, 512, JSON_THROW_ON_ERROR);
				$tryAgain = false;
			} catch (Exception $e) {
				$tryAgain = $retried < HttpRequestHelper::numRetriesOnHttpTooEarly();

				if (!$tryAgain) {
					throw $e;
				}
			}

			if ($tryAgain) {
				$retried += 1;
				echo "Invalid json body, retrying ($retried)...\n";
				// wait 500ms and try again
				\usleep(500 * 1000);
			}
		} while ($tryAgain);

		Assert::assertArrayHasKey(
			'bundles',
			$decodedBody,
			__METHOD__ . " could not find bundles in body"
		);
		$bundles = $decodedBody["bundles"];

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
		$roleId = $this->getRoleIdByRoleName($this->featureContext->getAdminUserName(), $role);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id');
		$response = $this->assignRoleToUser($user, $userId, $roleId);
		$this->featureContext->setResponse($response);
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
		$roleId = $this->getRoleIdByRoleName($this->featureContext->getAdminUserName(), $role);
		$userId = $this->featureContext->getAttributeOfCreatedUser($assignedUser, 'id');
		$response = $this->assignRoleToUser($user, $userId, $roleId);
		$this->featureContext->setResponse($response);
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
		$this->featureContext->setResponse($this->getAssignmentsList($user, $userId));
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
		$response = $this->getAssignmentsList($this->featureContext->getAdminUserName(), $userId);
		$assignmentResponse = $this->featureContext->getJsonDecodedResponseBodyContent($response);
		if (isset($assignmentResponse->assignments[0]->roleId)) {
			$actualRoleId = $assignmentResponse->assignments[0]->roleId;
			Assert::assertEquals($this->getRoleIdByRoleName($this->featureContext->getAdminUserName(), $role), $actualRoleId, "user $user has no role $role");
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
		Assert::assertEquals($this->getRoleIdByRoleName($this->featureContext->getAdminUserName(), $role), $assignmentRoleId, "user has no role $role");
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendRequestGetBundlesList(string $user): ResponseInterface {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "bundles-list";
		return HttpRequestHelper::post(
			$fullUrl,
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			null,
			"{}"
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
		$response = $this->sendRequestGetBundlesList($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			201,
			"Expected response status code should be 201",
			$response
		);

		$body = json_decode((string)$response->getBody(), true, 512, JSON_THROW_ON_ERROR);
		foreach ($body["bundles"] as $value) {
			if ($value["displayName"] === $bundleName) {
				return $value;
			}
		}
		return [];
	}

	/**
	 * @param string $user
	 * @param array $headers
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendRequestGetSettingsValuesList(string $user, array $headers = null): ResponseInterface {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "values-list";
		$body = json_encode(["account_uuid" => "me"], JSON_THROW_ON_ERROR);
		return HttpRequestHelper::post(
			$fullUrl,
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$headers,
			$body
		);
	}

	/**
	 * @param string $user
	 *
	 * @return bool
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function getAutoAcceptSharesSettingValue(string $user): bool {
		$response = $this->sendRequestGetSettingsValuesList($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			201,
			"Expected response status code should be 201",
			$response
		);

		$body = $this->featureContext->getJsonDecodedResponseBodyContent($response);

		if (empty($body)) {
			return true;
		}

		foreach ($body->values as $value) {
			if ($value->identifier->setting === "auto-accept-shares") {
				return $value->value->boolValue;
			}
		}

		return true;
	}

	/**
	 * @When /^user "([^"]*)" lists values-list with headers using the Settings API$/
	 *
	 * @param string $user
	 * @param TableNode $headersTable
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theUserListsAllValuesListWithHeadersUsingSettingsApi(string $user, TableNode $headersTable): void {
		$this->featureContext->verifyTableNodeColumns(
			$headersTable,
			['header', 'value']
		);
		$headers = [];
		foreach ($headersTable as $row) {
			$headers[$row['header']] = $row ['value'];
		}
		$this->featureContext->setResponse($this->sendRequestGetSettingsValuesList($user, $headers));
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
		$response = $this->sendRequestGetSettingsValuesList($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			201,
			"Expected response status code should be 201",
			$response
		);

		$body = json_decode((string)$response->getBody(), true, 512, JSON_THROW_ON_ERROR);

		// if no language is set, the request body is empty return English as the default language
		if (empty($body)) {
			return "en";
		}
		foreach ($body["values"] as $value) {
			if ($value["identifier"]["setting"] === "language") {
				return $value["value"]["listValue"]["values"][0]["stringValue"];
			}
		}
		// if a language setting was still not found, return English
		return "en";
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
		return HttpRequestHelper::post(
			$fullUrl,
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			null,
			$body
		);
	}

	/**
	 * @Given /^user "([^"]*)" has switched the system language to "([^"]*)" using the settings API$/
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

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendRequestToDisableAutoAccepting(string $user): ResponseInterface {
		$fullUrl = $this->baseUrl . $this->settingsUrl . "values-save";
		$body = json_encode(
			[
				"value" => [
					"account_uuid" => "me",
					"bundleId" => "2a506de7-99bd-4f0d-994e-c38e72c28fd9",
					"settingId" => "ec3ed4a3-3946-4efc-8f9f-76d38b12d3a9",
					"resource" => [
						"type" => "TYPE_USER"
					],
					"boolValue" => false
				]
			],
			JSON_THROW_ON_ERROR
		);

		return HttpRequestHelper::post(
			$fullUrl,
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			[],
			$body
		);
	}

	/**
	 * @Given user :user has disabled auto-accepting
	 * @Given user :user has disabled the auto-sync share
	 *
	 * @param string $user
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function theUserHasDisabledAutoAccepting(string $user): void {
		$response = $this->sendRequestToDisableAutoAccepting($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			201,
			"Expected response status code should be 201",
			$response
		);
		$this->featureContext->rememberUserAutoSyncSetting($user, false);
	}
}
