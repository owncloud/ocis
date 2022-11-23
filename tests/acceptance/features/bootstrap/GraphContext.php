<?php

declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Kiran Parajuli <kiran@jankaritech.com>
 * @copyright Copyright (c) 2021 Kiran Parajuli kiran@jankaritech.com
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Behat\Gherkin\Node\TableNode;
use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\GraphHelper;
use TestHelpers\WebDavHelper;
use PHPUnit\Framework\Assert;

require_once 'bootstrap.php';

/**
 * Context for the provisioning specific steps using the Graph API
 */
class GraphContext implements Context {
	/**
	 * @var FeatureContext
	 */
	private FeatureContext $featureContext;

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
	}

	/**
	 * @param string $user
	 * @param string|null $userName
	 * @param string|null $password
	 * @param string|null $email
	 * @param string|null $displayName
	 * @param string|null $requester
	 * @param string|null $requesterPassword
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userHasBeenEditedUsingTheGraphApi(
		string $user,
		?string $userName = null,
		?string $password = null,
		?string $email = null,
		?string $displayName = null,
		?string $requester = null,
		?string $requesterPassword = null
	): void {
		if (!$requester) {
			$requester = $this->featureContext->getAdminUsername();
			$requesterPassword = $this->featureContext->getAdminPassword();
		}
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id');
		$response = GraphHelper::editUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$requester,
			$requesterPassword,
			$userId,
			$userName,
			$password,
			$email,
			$displayName
		);
		$this->featureContext->setResponse($response);
		$this->featureContext->theHttpStatusCodeShouldBe(200);
	}

	/**
	 * @param string $user
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function adminHasRetrievedUserUsingTheGraphApi(string $user): void {
		$user = $this->featureContext->getActualUsername($user);
		try {
			$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		} catch (Exception $e) {
			$userId = $user;
		}
		$result = GraphHelper::getUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$userId
		);
		$this->featureContext->setResponse($result);
		$this->featureContext->thenTheHTTPStatusCodeShouldBe(200);
	}

	/**
	 * @param $requestingUser
	 * @param $targetUser
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userHasRetrievedUserUsingTheGraphApi(
		$requestingUser,
		$targetUser
	): void {
		$requester = $this->featureContext->getActualUsername($requestingUser);
		$requesterPassword = $this->featureContext->getPasswordForUser($requestingUser);
		$user = $this->featureContext->getActualUsername($targetUser);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$response = GraphHelper::getUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$requester,
			$requesterPassword,
			$userId
		);
		$this->featureContext->setResponse($response);
		$this->featureContext->thenTheHTTPStatusCodeShouldBe(200);
	}

	/**
	 * @param string $groupId
	 * @param string|null $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function userDeletesGroupWithGroupId(
		string $groupId,
		?string $user = null
	): ResponseInterface {
		$credentials = $this->getAdminOrUserCredentials($user);

		return GraphHelper::deleteGroup(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials["username"],
			$credentials["password"],
			$groupId
		);
	}

	/**
	 * @param string $groupId
	 * @param bool $checkResult
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function adminDeletesGroupWithGroupId(
		string $groupId,
		bool $checkResult = false
	): void {
		$this->featureContext->setResponse(
			$this->userDeletesGroupWithGroupId($groupId)
		);
		if ($checkResult) {
			$this->featureContext->thenTheHTTPStatusCodeShouldBe(204);
		}
	}

	/**
	 * @param string $group
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function adminDeletesGroupUsingTheGraphApi(
		string $group
	): void {
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		if ($groupId) {
			$this->adminDeletesGroupWithGroupId($groupId);
		} else {
			throw new Exception(
				"Group id does not exist for '$group' in the created list."
				. " Cannot delete group without id when using the Graph API."
			);
		}
	}

	/**
	 * @param string $group
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function adminHasDeletedGroupUsingTheGraphApi(string $group): void {
		$this->adminDeletesGroupUsingTheGraphApi($group);
		$this->featureContext->thenTheHTTPStatusCodeShouldBe(204);
	}

	/**
	 * sends a request to delete a user using the Graph API
	 *
	 * @param string $user username is used as the id
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function adminDeletesUserUsingTheGraphApi(string $user): void {
		$this->featureContext->setResponse(
			GraphHelper::deleteUser(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$this->featureContext->getAdminUsername(),
				$this->featureContext->getAdminPassword(),
				$user
			)
		);
	}

	/**
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function adminHasRemovedUserFromGroupUsingTheGraphApi(string $user, string $group): void {
		$user = $this->featureContext->getActualUsername($user);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		$response = GraphHelper::removeUserFromGroup(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$userId,
			$groupId,
		);
		$this->featureContext->setResponse($response);
		$this->featureContext->thenTheHTTPStatusCodeShouldBe(204);
	}

	/**
	 * check if the provided user is present as a member in the provided group
	 *
	 * @param string $user
	 * @param string $group
	 *
	 * @return bool
	 * @throws JsonException
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function getUserPresenceInGroupUsingTheGraphApi(string $user, string $group): bool {
		$user = $this->featureContext->getActualUsername($user);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$members = $this->theAdminHasRetrievedMembersListOfGroupUsingTheGraphApi($group);
		$found = false;
		foreach ($members as $member) {
			if ($member["id"] === $userId) {
				$found = true;
				break;
			}
		}
		return $found;
	}

	/**
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userShouldNotBeMemberInGroupUsingTheGraphApi(string $user, string $group): void {
		$found = $this->getUserPresenceInGroupUsingTheGraphApi($user, $group);
		Assert::assertFalse($found, __METHOD__ . " User $user is member of group $group");
	}

	/**
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userShouldBeMemberInGroupUsingTheGraphApi(string $user, string $group): void {
		$found = $this->getUserPresenceInGroupUsingTheGraphApi($user, $group);
		Assert::assertTrue($found, __METHOD__ . "User $user is not member of group $group");
	}

	/**
	 * @param string $user
	 * @param string $password
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function adminChangesPasswordOfUserToUsingTheGraphApi(
		string $user,
		string $password
	): void {
		$user = $this->featureContext->getActualUsername($user);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id');
		$response = GraphHelper::editUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$userId,
			null,
			$password
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 *
	 * @param array $groups
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theseGroupsShouldBeInTheResponse(array $groups): void {
		$respondedGroups = $this->getArrayOfGroupsResponded($this->featureContext->getResponse());
		foreach ($groups as $group) {
			$found = false;
			foreach ($respondedGroups as $respondedGroup) {
				if ($respondedGroup["displayName"] === $group) {
					$found = true;
					break;
				}
			}
			Assert::assertTrue($found, "Group '$group' not found in the list");
		}
	}

	/**
	 *
	 * @param array $users
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theseUsersShouldBeInTheResponse(array $users): void {
		$respondedUsers = $this->getArrayOfUsersResponded($this->featureContext->getResponse());
		foreach ($users as $user) {
			$found = false;
			foreach ($respondedUsers as $respondedUser) {
				if ($respondedUser["onPremisesSamAccountName"] === $user) {
					$found = true;
					break;
				}
			}
			Assert::assertTrue($found, "User '$user' not found in the list");
		}
	}

	/**
	 *
	 * @param string|null $user
	 *
	 * @return array
	 */
	public function getAdminOrUserCredentials(?string $user): array {
		$credentials["username"] = $user ? $this->featureContext->getActualUsername($user) : $this->featureContext->getAdminUsername();
		$credentials["password"] = $user ? $this->featureContext->getPasswordForUser($user) : $this->featureContext->getAdminPassword();
		return $credentials;
	}
	/**
	 *
	 * @param string|null $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function listGroups(?string $user = null): ResponseInterface {
		$credentials = $this->getAdminOrUserCredentials($user);

		return GraphHelper::getGroups(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials["username"],
			$credentials["password"]
		);
	}

	/**
	 * returns list of groups
	 *
	 * @param ResponseInterface $response
	 *
	 * @return array
	 * @throws Exception
	 */
	public function getArrayOfGroupsResponded(ResponseInterface $response): array {
		if ($response->getStatusCode() === 200) {
			$jsonResponseBody = $this->featureContext->getJsonDecodedResponse($response);
			return $jsonResponseBody["value"];
		} else {
			$this->throwHttpException($response, "Could not retrieve groups list.");
		}
	}

	/**
	 *
	 * @return array
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function adminHasRetrievedGroupListUsingTheGraphApi(): array {
		return  $this->getArrayOfGroupsResponded($this->listGroups());
	}

	/**
	 *
	 * @param string $group
	 * @param string|null $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function listGroupMembers(string $group, ?string $user = null): ResponseInterface {
		$credentials = $this->getAdminOrUserCredentials($user);

		return GraphHelper::getMembersList(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials["username"],
			$credentials["password"],
			$this->featureContext->getAttributeOfCreatedGroup($group, 'id')
		);
	}

	/**
	 * returns list of users of a group
	 *
	 * @param ResponseInterface $response
	 *
	 * @return array
	 * @throws Exception
	 */
	public function getArrayOfUsersResponded(ResponseInterface $response): array {
		if ($response->getStatusCode() === 200) {
			return $this->featureContext->getJsonDecodedResponse($response);
		} else {
			$this->throwHttpException($response, "Could not retrieve group members list.");
		}
	}

	/**
	 * returns a list of members in group
	 *
	 * @param string $group
	 *
	 * @return array
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function theAdminHasRetrievedMembersListOfGroupUsingTheGraphApi(string $group): array {
		return $this->getArrayOfUsersResponded($this->listGroupMembers($group));
	}

	/**
	 * creates a user with provided data
	 * actor: the administrator
	 *
	 * @param string $user
	 * @param string $password
	 * @param string $email
	 * @param string $displayName
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function theAdminHasCreatedUser(
		string $user,
		string $password,
		string $email,
		string $displayName
	): void {
		$response = GraphHelper::createUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$user,
			$password,
			$email,
			$displayName
		);
		if ($response->getStatusCode() !== 200) {
			$this->throwHttpException($response, "Could not create user $user");
		} else {
			$this->featureContext->setResponse($response);
		}
	}

	/**
	 * @When /^the user "([^"]*)" creates a new user using GraphAPI with the following settings:$/
	 *
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function theUserCreatesNewUser(string $user, TableNode $table): void {
		$rows = $table->getRowsHash();
		$response = GraphHelper::createUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$rows["userName"],
			$rows["password"],
			$rows["email"],
			$rows["displayName"]
		);

		// add created user to list except for the user with an empty name
		// because request /graph/v1.0/users/emptyUserName exits with 200
		// and we cannot check that the user with empty name doesn't exist
		if (!empty($rows["userName"])) {
			$this->featureContext->addUserToCreatedUsersList(
				$rows["userName"],
				$rows["password"],
				$rows["displayName"],
				$rows["email"]
			);
		}
		$this->featureContext->setResponse($response);
	}

	/**
	 * adds a user to a group
	 *
	 * @param string $group
	 * @param string $user
	 * @param string|null $byUser
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function addUserToGroup(string $group, string $user, ?string $byUser = null): ResponseInterface {
		$credentials = $this->getAdminOrUserCredentials($byUser);
		try {
			$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		} catch (Exception $e) {
			$groupId = WebDavHelper::generateUUIDv4();
		}
		try {
			$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		} catch (Exception $e) {
			$userId = WebDavHelper::generateUUIDv4();
		}

		return GraphHelper::addUserToGroup(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials['username'],
			$credentials['password'],
			$userId,
			$groupId
		);
	}

	/**
	 * @Given /^the administrator has added a user "([^"]*)" to the group "([^"]*)" using GraphApi$/
	 *
	 * @param string $user
	 * @param string $group
	 * @param bool $checkResult
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function adminHasAddedUserToGroupUsingTheGraphApi(
		string $user,
		string $group,
		bool $checkResult = true
	): void {
		$result = $this->addUserToGroup($group, $user);
		if ($checkResult && ($result->getStatusCode() !== 204)) {
			$this->throwHttpException($result, "Could not add user '$user' to group '$group'.");
		}
	}

	/**
	 * @When the administrator adds the following users to the following groups using the Graph API
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function theAdministratorAddsTheFollowingUsersToTheFollowingGroupsUsingTheGraphAPI(TableNode $table): void {
		$this->featureContext->verifyTableNodeColumns($table, ['username', 'groupname']);
		$userGroupList = $table->getColumnsHash();

		foreach ($userGroupList as $userGroup) {
			$this->featureContext->setResponse($this->addUserToGroup($userGroup['groupname'], $userGroup['username']));
			$this->featureContext->pushToLastHttpStatusCodesArray();
		}
	}

	/**
	 * @When the administrator tries to add user :user to group :group using the Graph API
	 *
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 */
	public function theAdministratorTriesToAddUserToGroupUsingTheGraphAPI(string $user, string $group): void {
		$this->featureContext->setResponse($this->addUserToGroup($group, $user));
	}
 
	/**
	 * @When user :user tries to add himself/herself to group :group using the Graph API
	 *
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 */
	public function theUserTriesToAddHimselfToGroupUsingTheGraphAPI(string $user, string $group): void {
		$this->featureContext->setResponse($this->addUserToGroup($group, $user, $user));
	}

	/**
	 * @When user :byUser tries to add user :user to group :group using the Graph API
	 *
	 * @param string $byUser
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 */
	public function theUserTriesToAddAnotherUserToGroupUsingTheGraphAPI(string $byUser, string $user, string $group): void {
		$this->featureContext->setResponse($this->addUserToGroup($group, $byUser, $user));
	}

	/**
	 *
	 * @param string $group
	 * @param ?string $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function createGroup(string $group, ?string $user = null): ResponseInterface {
		$credentials = $this->getAdminOrUserCredentials($user);

		return GraphHelper::createGroup(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials["username"],
			$credentials["password"],
			$group,
		);
	}

	/**
	 * @When /^the administrator creates a group "([^"]*)" using the Graph API$/
	 * @When user :user creates a group :group using the Graph API
	 * @When user :user tries to create a group :group using the Graph API
	 *
	 * @param string $group
	 * @param ?string $user
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userCreatesGroupUsingTheGraphApi(string $group, ?string $user = null): void {
		$response = $this->createGroup($group, $user);
		$this->featureContext->setResponse($response);
		$this->featureContext->pushToLastHttpStatusCodesArray((string) $response->getStatusCode());

		if ($response->getStatusCode() === 200) {
			$groupId = $this->featureContext->getJsonDecodedResponse($response)["id"];
			$this->featureContext->addGroupToCreatedGroupsList($group, true, true, $groupId);
		}
	}

	/**
	 * create group with provided data
	 *
	 * @param string $group
	 *
	 * @return array
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function adminHasCreatedGroupUsingTheGraphApi(string $group): array {
		$result = $this->createGroup($group);
		if ($result->getStatusCode() === 200) {
			return $this->featureContext->getJsonDecodedResponse($result);
		} else {
			$this->throwHttpException($result, "Could not create group '$group'.");
		}
	}

	/**
	 * @param ResponseInterface $response
	 * @param string $errorMsg
	 *
	 * @return void
	 * @throws Exception
	 */
	private function throwHttpException(ResponseInterface $response, string $errorMsg): void {
		try {
			$jsonBody = $this->featureContext->getJsonDecodedResponse($response);
			throw new Exception(
				__METHOD__
				. "\n$errorMsg"
				. "\nHTTP status code: " . $response->getStatusCode()
				. "\nError code: " . $jsonBody["error"]["code"]
				. "\nMessage: " . $jsonBody["error"]["message"]
			);
		} catch (TypeError $e) {
			throw new Exception(
				__METHOD__
				. "\n$errorMsg"
				. "\nHTTP status code: " . $response->getStatusCode()
				. "\nResponse body: " . $response->getBody()
			);
		}
	}

	/**
	 * @param string $shouldOrNot (not|)
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theseGroupsShouldNotExist(string $shouldOrNot, TableNode $table): void {
		$should = ($shouldOrNot !== "not");
		$this->featureContext->verifyTableNodeColumns($table, ['groupname']);
		$actualGroupsList = $this->adminHasRetrievedGroupListUsingTheGraphApi();
		$expectedGroups = $table->getColumnsHash();
		// check if every expected group is(not) in the actual groups list
		foreach ($expectedGroups as $expectedGroup) {
			$groupName = $expectedGroup['groupname'];
			$groupExists = false;
			foreach ($actualGroupsList as $actualGroup) {
				if ($actualGroup['displayName'] === $groupName) {
					$groupExists = true;
					break;
				}
			}
			if ($groupExists !== $should) {
				throw new Exception(
					__METHOD__
					. "\nGroup '$groupName' is expected " . ($should ? "" : "not ")
					. "to exist, but it does" . ($should ? " not" : "") . " exist."
				);
			}
		}
	}

	/**
	 * @When /^the user "([^"]*)" changes its own password "([^"]*)" to "([^"]*)" using the Graph API$/
	 *
	 * @param string $user
	 * @param string $currentPassword
	 * @param string $newPassword
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function userChangesOwnPassword(string $user, string $currentPassword, string $newPassword): void {
		$response = GraphHelper::changeOwnPassword(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$currentPassword,
			$newPassword
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user gets all the groups using the Graph API
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function userGetsAllTheGroupsUsingTheGraphApi(string $user): void {
		$this->featureContext->setResponse($this->listGroups($user));
	}

	/**
	 * @When user :user gets all the members of group :group using the Graph API
	 *
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 */
	public function userGetsAllTheMembersOfGroupUsingTheGraphApi($user, $group): void {
		$this->featureContext->setResponse($this->listGroupMembers($group, $user));
	}

	/**
	 * @Then the last response should be an unauthorized response
	 *
	 * @return void
	 */
	public function theLastResponseShouldBeUnauthorizedReponse(): void {
		$response = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse());
		$errorText = $response['error']['message'];

		Assert::assertEquals(
			'Unauthorized',
			$errorText,
			__METHOD__
			. "\nExpected unauthorized message but got '" . $errorText . "'"
		);
	}

	/**
	 * @When user :user deletes group :group using the Graph API
	 * @When the administrator deletes group :group using the Graph API
	 * @When user :user tries to delete group :group using the Graph API
	 *
	 * @param string $group
	 * @param string|null $user
	 *
	 * @return void
	 */
	public function userDeletesGroupUsingTheGraphApi(string $group, ?string $user): void {
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		$response = $this->userDeletesGroupWithGroupId($groupId, $user);
		$this->featureContext->setResponse($response);
	}
  
	/**
	 * @Then the following users should be listed in the following groups
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theFollowingUsersShouldBeListedInFollowingGroups(TableNode $table): void {
		$this->featureContext->verifyTableNodeColumns($table, ['username', 'groupname']);
		$usersGroups = $table->getColumnsHash();
		foreach ($usersGroups as $userGroup) {
			$members = $this->listGroupMembers($userGroup['groupname']);
			$members = $this->featureContext->getJsonDecodedResponse($members);

			$exists = false;
			foreach ($members as $member) {
				if ($member['onPremisesSamAccountName'] === $userGroup['username']) {
					$exists = true;
					break;
				}
			}
			Assert::assertTrue(
				$exists,
				__METHOD__
				. "\nExpected user '" . $userGroup['username'] . "' to be in group '" . $userGroup['groupname'] . "'. But not found."
			);
		}
	}
}
