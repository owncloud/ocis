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
use Behat\Gherkin\Node\TableNode;
use PHPUnit\Framework\Assert;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\HttpRequestHelper;
use TestHelpers\SettingsHelper;
use TestHelpers\BehatHelper;

require_once 'bootstrap.php';

/**
 * Context for the TUS-specific steps using the Graph API
 */
class SettingsContext implements Context {
	private FeatureContext $featureContext;
	private array $autoAcceptSharesSettingIds = [];

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
		$this->featureContext = BehatHelper::getContext($scope, $environment, 'FeatureContext');
	}

	/**
	 * @param string $user
	 *
	 * @return string
	 */
	public function getAutoAcceptShareSettingId(string $user): string {
		if (!empty($this->autoAcceptSharesSettingIds) && \array_key_exists($user, $this->autoAcceptSharesSettingIds)) {
			return $this->autoAcceptSharesSettingIds[$user];

		}
		return '';
	}

	/**
	 * @param string $user
	 * @param string $id
	 *
	 * @return void
	 */
	public function setAutoAcceptShareSettingId(string $user, string $id): void {
		$this->autoAcceptSharesSettingIds[$user] = $id;
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
		return SettingsHelper::getRolesList(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
		);
	}

	/**
	 * @When /^user "([^"]*)" tries to get all existing roles using the settings API$/
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
		return SettingsHelper::assignRoleToUser(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$userId,
			$roleId,
		);
	}

	/**
	 * @param string $user
	 * @param string $userId
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function getAssignmentsList(string $user, string $userId): ResponseInterface {
		return SettingsHelper::getAssignmentsList(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$userId,
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
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id') ?: $user;
		$response = $this->assignRoleToUser($admin, $userId, $roleId);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			201,
			"Expected response status code should be 201",
			$response,
		);
	}

	/**
	 * @When user :assigner assigns the role :role to user :assignee using the settings API
	 *
	 * @param string $assigner
	 * @param string $role
	 * @param string $assignee
	 *
	 * @return void
	 *
	 * @throws Exception
	 */
	public function userAssignsTheRoleToUserUsingTheSettingsApi(
		string $assigner,
		string $role,
		string $assignee,
	): void {
		$response = $this->assignRoleToUser(
			$assigner,
			$this->featureContext->getAttributeOfCreatedUser($assignee, 'id'),
			$this->getRoleIdByRoleName($assigner, $role),
		);
		$this->featureContext->setResponse($response);
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
				$response,
			);

			$rawBody = $response->getBody()->getContents();
			try {
				$decodedBody = \json_decode($rawBody, true, 512, JSON_THROW_ON_ERROR);
				$tryAgain = false;
			} catch (Exception $e) {
				$tryAgain = $retried < HttpRequestHelper::maxHTTPRequestRetries();

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
			__METHOD__ . " could not find bundles in body",
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
	 * @When /^user "([^"]*)" tries to get list of assignment using the settings API$/
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
			Assert::assertEquals(
				$this->getRoleIdByRoleName($this->featureContext->getAdminUserName(), $role),
				$actualRoleId,
				"user $user has no role $role",
			);
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
		$assignmentRoleId = $this->featureContext->getJsonDecodedResponse(
			$this->featureContext->getResponse(),
		)["assignments"][0]["roleId"];
		Assert::assertEquals(
			$this->getRoleIdByRoleName($this->featureContext->getAdminUserName(), $role),
			$assignmentRoleId,
			"user has no role $role",
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
	public function sendRequestGetBundlesList(string $user): ResponseInterface {
		return SettingsHelper::getBundlesList(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
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
	public function getBundleByName(string $user, string $bundleName): array {
		return SettingsHelper::getBundleByName(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$bundleName,
		);
	}

	/**
	 * @param string $user
	 * @param array|null $headers
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendRequestGetSettingsValuesList(string $user, ?array $headers = null): ResponseInterface {
		return SettingsHelper::getValuesList(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$headers,
		);
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
			['header', 'value'],
		);
		$headers = [];
		foreach ($headersTable as $row) {
			$headers[$row['header']] = $row ['value'];
		}
		$this->featureContext->setResponse($this->sendRequestGetSettingsValuesList($user, $headers));
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
		$profileBundlesList = $this->getBundleByName($user, "Profile");
		Assert::assertNotEmpty($profileBundlesList, "bundles list is empty");

		$settingId = '';
		foreach ($profileBundlesList["settings"] as $value) {
			if ($value["name"] === "language") {
				$settingId = $value["id"];
				break;
			}
		}
		Assert::assertNotEmpty($settingId, "settingId is empty");

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
								"stringValue" => $language,
							],
						],
					],
					"resource" => [
						"type" => "TYPE_USER",
					],
					"settingId" => $settingId,
				],
			],
			JSON_THROW_ON_ERROR,
		);
		return SettingsHelper::updateSettings(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$body,
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
			$response,
		);
	}

	/**
	 * @When user :user switches the system language to :language using the settings API
	 *
	 * @param string $user
	 * @param string $language
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userSwitchesTheSystemLanguageUsingTheSettingsApi(string $user, string $language): void {
		$response = $this->sendRequestToSwitchSystemLanguage($user, $language);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @param string $user
	 * @param bool $status
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function toggleAutoAcceptSharesSetting(string $user, bool $status): ResponseInterface {
		$body = [
			"value" => [
				"account_uuid" => "me",
				"bundleId" => SettingsHelper::getBundleId(),
				"settingId" => SettingsHelper::getSettingIdUsingEventName("Auto Accept Shares"),
				"resource" => [
					"type" => "TYPE_USER",
				],
				"boolValue" => $status,
			],
		];
		$autoAcceptSharesSettingId = $this->getAutoAcceptShareSettingId($user);
		if ($autoAcceptSharesSettingId) {
			// use existing id if available
			$body["value"]["id"] = $autoAcceptSharesSettingId;
		}
		$body = json_encode($body, JSON_THROW_ON_ERROR);

		$response = SettingsHelper::updateSettings(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$body,
		);
		if (!$autoAcceptSharesSettingId && $response->getStatusCode() === 201) {
			// save id for future use
			// updating the setting without id will create a new setting entry
			$data = $this->featureContext->getJsonDecodedResponseBodyContent($response);
			$this->setAutoAcceptShareSettingId($user, $data->value->value->id);
		}
		$response->getBody()->rewind();
		return $response;
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
		$response = $this->toggleAutoAcceptSharesSetting($user, false);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			201,
			"Expected response status code should be 201",
			$response,
		);
		$this->featureContext->rememberUserAutoSyncSetting($user, false);
	}

	/**
	 * @When /^user "([^"]*)" (disables|enables) the auto-sync share using the settings API$/
	 *
	 * @param string $user
	 * @param string $status
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userEnablesOrDisablesTheAutoSyncShareUsingTheSettingsApi(string $user, string $status): void {
		$enable = $status === "enables";
		$response = $this->toggleAutoAcceptSharesSetting($user, $enable);
		$this->featureContext->setResponse($response);
		$this->featureContext->rememberUserAutoSyncSetting($user, $enable);
	}

	/**
	 * @When user :user disables email notification using the settings API
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function userDisablesEmailNotificationUsingTheSettingsAPI(string $user): void {
		$body = [
			"value" => [
				"account_uuid" => "me",
				"bundleId" => SettingsHelper::getBundleId(),
				"settingId" => SettingsHelper::getSettingIdUsingEventName("Disable Email Notifications"),
				"resource" => [
					"type" => "TYPE_USER",
				],
				"boolValue" => true,
			],
		];
		$response = SettingsHelper::updateSettings(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			json_encode($body),
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @param string $user
	 * @param string $enableOrDisable
	 * @param array $event
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function enableOrDisableNotification(
		string $user,
		string $enableOrDisable,
		array $event,
	): ResponseInterface {
		$body = [
			"value" => [
				"account_uuid" => "me",
				"bundleId" => SettingsHelper::getBundleId(),
				"settingId" => SettingsHelper::getSettingIdUsingEventName($event['event']),
				"resource" => [
					"type" => "TYPE_USER",
				],
			],
		];
		$notificationTypes = explode(',', $event['notificationTypes']);

		foreach ($notificationTypes as $type) {
			if ($type === "mail") {
				$body["value"]["collectionValue"]["values"][] = [
					"key" => "mail",
					"boolValue" => $enableOrDisable === "enables",
				];
			} elseif ($type === "in-app") {
				$body["value"]["collectionValue"]["values"][] = [
					"key" => "in-app",
					"boolValue" => $enableOrDisable === "enables",
				];
			}
		}

		return SettingsHelper::updateSettings(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			json_encode($body),
		);
	}

	/**
	 * @Given /^user "([^"]*)" has (disabled|enabled) notification for the following event using the settings API:$/
	 *
	 * @param string $user
	 * @param string $enableOrDisable
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function userHasEnabledOrDisabledNotificationForTheFollowingEventUsingTheSettingsApi(
		string $user,
		string $enableOrDisable,
		TableNode $table,
	): void {
		$event = $table->getRowsHash();
		$response = $this->enableOrDisableNotification($user, $enableOrDisable, $event);
		$this->featureContext->theHTTPStatusCodeShouldBe(201, "", $response);
	}

	/**
	 * @When /^user "([^"]*)" (disables|enables) notification for the following event using the settings API:$/
	 *
	 * @param string $user
	 * @param string $enableOrDisable
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function userEnablesOrDisablesNotificationForTheFollowingEventUsingTheSettingsApi(
		string $user,
		string $enableOrDisable,
		TableNode $table,
	): void {
		$event = $table->getRowsHash();
		$response = $this->enableOrDisableNotification($user, $enableOrDisable, $event);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @param string $user
	 * @param string $interval
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function setEmailSendingInterval(string $user, string $interval): ResponseInterface {
		$body = [
			"value" => [
				"account_uuid" => "me",
				"bundleId" => SettingsHelper::getBundleId(),
				"settingId" => SettingsHelper::getSettingIdUsingEventName("Email Sending Interval"),
				"resource" => [
					"type" => "TYPE_USER",
				],
				"stringValue" => $interval,
			],
		];
		return SettingsHelper::updateSettings(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			json_encode($body, JSON_THROW_ON_ERROR),
		);
	}

	/**
	 * @When /^user "([^"]*)" sets the email sending interval to "([^"]*)" using the settings API$/
	 *
	 * @param string $user
	 * @param string $interval
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function userSetsTheEmailSendingIntervalToUsingTheSettingsAPI(
		string $user,
		string $interval,
	): void {
		$response = $this->setEmailSendingInterval($user, $interval);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given /^user "([^"]*)" has set the email sending interval to "([^"]*)" using the settings API$/
	 *
	 * @param string $user
	 * @param string $interval
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function userHasSetTheEmailSendingIntervalToUsingTheSettingsAPI(
		string $user,
		string $interval,
	): void {
		$response = $this->setEmailSendingInterval($user, $interval);
		$this->featureContext->theHTTPStatusCodeShouldBe(201, "", $response);
	}

	/**
	 * @Then for user :user setting :setting should have value :value
	 *
	 * @param string $user
	 * @param string $settingName
	 * @param string $value
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function forUserSettingShouldHaveValue(string $user, string $settingName, string $value): void {
		$response = SettingsHelper::getValuesBySettingID(
			$this->featureContext->getBaseUrl(),
			SettingsHelper::PROFILE_SETTINGS[$settingName],
			$user,
			$this->featureContext->getPasswordForUser($user),
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(201, "", $response);

		$setting = HttpRequestHelper::getJsonDecodedResponseBodyContent($response)->value->value;
		if (\property_exists($setting, 'stringValue')) {
			$settingValue = $setting->stringValue;
		} elseif (\property_exists($setting, 'boolValue')) {
			$settingValue = $setting->boolValue;
			$value = $value === "true" ? true : false;
		} elseif (\property_exists($setting, 'listValue')) {
			$settingValue = $setting->listValue->values[0]->stringValue;
		} else {
			Assert::fail(
				"Setting '$settingName' does not have a stringValue, boolValue or listValue."
				. "\n"
				. json_encode($setting),
			);
		}

		Assert::assertSame(
			$value,
			$settingValue,
			"Expected setting value '$value' but got '$settingValue'",
		);
	}
}
