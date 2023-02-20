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

require_once 'bootstrap.php';

/**
 * Context for the TUS-specific steps using the Graph API
 */
class RoleAssignmentContext implements Context {

	/**
	 * @var FeatureContext
	 */
	private FeatureContext $featureContext;

	/**
	 * @var SpacesContext
	 */
	private SpacesContext $spacesContext;

	/**
	 * @var string
	 */
	private string $baseUrl;

	/**
	 * @var string
	 */
	private string $setttingsUrl = '/api/v0/settings/';

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
		$fullUrl = $this->baseUrl . $this->setttingsUrl . "roles-list";
		$this->featureContext->setResponse(
			$this->spacesContext->sendPostRequestToUrl($fullUrl, $user, $this->featureContext->getPasswordForUser($user), "{}")
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
		$fullUrl = $this->baseUrl . $this->setttingsUrl . "assignments-add";
		$body = json_encode(["account_uuid" => $userId, "role_id" => $roleId], JSON_THROW_ON_ERROR);

		$this->featureContext->setResponse(
			$this->spacesContext->sendPostRequestToUrl($fullUrl, $user, $this->featureContext->getPasswordForUser($user), $body)
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
	public function sendRequestAssignmentsList(string $user, string $userId): void {
		$fullUrl = $this->baseUrl . $this->setttingsUrl . "assignments-list";
		$body = json_encode(["account_uuid" => $userId], JSON_THROW_ON_ERROR);

		$this->featureContext->setResponse(
			$this->spacesContext->sendPostRequestToUrl($fullUrl, $user, $this->featureContext->getPasswordForUser($user), $body)
		);
	}

	/**
	 * @When /^the administrator has given "([^"]*)" the role "([^"]*)" using the settings api$/
	 *
	 * @param string $user
	 * @param string $role
	 *
	 * @return void
	 *
	 * @throws GuzzleException
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
	public function userGetRoleIdByRoleName($user, $role): string {
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
	public function setRoleToUser($user, $userId, $roleId): void {
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
		$this->sendRequestAssignmentsList($user, $userId);
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
		$this->sendRequestAssignmentsList($this->featureContext->getAdminUserName(), $userId);
		$rawBody =  $this->featureContext->getResponse()->getBody()->getContents();
		$assignmentRoleId = \json_decode($rawBody, true, 512, JSON_THROW_ON_ERROR)["assignments"][0]["roleId"];
		Assert::assertEquals($this->userGetRoleIdByRoleName($this->featureContext->getAdminUserName(), $role), $assignmentRoleId, "user $user has no role $role");
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
}
