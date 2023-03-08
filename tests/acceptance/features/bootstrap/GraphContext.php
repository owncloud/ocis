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
	 * application Entity
	 *
	 * @var array
	 */
	private $appEntity = [];

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
		$this->featureContext->theHttpStatusCodeShouldBe(200); // TODO 204 when prefer=minimal header was sent
	}

	/**
	 * @When /^the user "([^"]*)" changes the email of user "([^"]*)" to "([^"]*)" using the Graph API$/
	 * @When /^the user "([^"]*)" tries to change the email of user "([^"]*)" to "([^"]*)" using the Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 * @param string $email
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theUserChangesTheEmailOfUserToUsingTheGraphApi(string $byUser, string $user, string $email): void {
		$response = $this->editUserUsingTheGraphApi($byUser, $user, null, null, $email);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^the user "([^"]*)" changes the display name of user "([^"]*)" to "([^"]*)" using the Graph API$/
	 * @When /^the user "([^"]*)" tries to change the display name of user "([^"]*)" to "([^"]*)" using the Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 * @param string $displayName
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theUserChangesTheDisplayNameOfUserToUsingTheGraphApi(string $byUser, string $user, string $displayName): void {
		$response = $this->editUserUsingTheGraphApi($byUser, $user, null, null, null, $displayName);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^the user "([^"]*)" disables user "([^"]*)" using the Graph API$/
	 * @When /^the user "([^"]*)" tries to disable user "([^"]*)" using the Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theUserDisablesUserToUsingTheGraphApi(string $byUser, string $user): void {
		$response = $this->editUserUsingTheGraphApi($byUser, $user, null, null, null, null, false);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given /^the user "([^"]*)" has disabled user "([^"]*)" using the Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theUserHasDisabledUserToUsingTheGraphApi(string $byUser, string $user): void {
		$this->theUserDisablesUserToUsingTheGraphApi($byUser, $user);
		$this->featureContext->thenTheHTTPStatusCodeShouldBe(200);
	}

	/**
	 * @When /^the user "([^"]*)" enables user "([^"]*)" using the Graph API$/
	 * @When /^the user "([^"]*)" tries to enable user "([^"]*)" using the Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theUserEnablesUserToUsingTheGraphApi(string $byUser, string $user): void {
		$response = $this->editUserUsingTheGraphApi($byUser, $user, null, null, null, null, true);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Then /^the user "([^"]*)" should have information with these key and value pairs:$/
	 *
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function theUserShouldHaveInformationWithTheseKeyAndValuePairs(string $user, TableNode $table): void {
		$rows = $table->getHash();
		$this->adminHasRetrievedUserUsingTheGraphApi($user);
		foreach ($rows as $row) {
			$key = $row['key'];
			$expectedValue = $row['value'];
			$responseValue = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse())[$key];
			Assert::assertEquals(
				$expectedValue,
				$responseValue
			);
		}
	}

	/**
	 * Edits the user information
	 *
	 * @param string $byUser
	 * @param string $user
	 * @param string|null $userName
	 * @param string|null $password
	 * @param string|null $email
	 * @param string|null $displayName
	 * @param bool|true $accountEnabled
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function editUserUsingTheGraphApi(string $byUser, string $user, string $userName = null, string $password = null, string $email = null, string $displayName = null, bool $accountEnabled = true): ResponseInterface {
		$user = $this->featureContext->getActualUsername($user);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id');
		$userId = $userId ?? $user;
		return GraphHelper::editUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$byUser,
			$this->featureContext->getPasswordForUser($byUser),
			$userId,
			$userName,
			$password,
			$email,
			$displayName,
			$accountEnabled
		);
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
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$userId = $userId ? $userId : $user;
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
	 * @param string|null $byUser
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function adminDeletesUserUsingTheGraphApi(string $user, ?string $byUser = null): void {
		$credentials = $this->getAdminOrUserCredentials($byUser);

		$this->featureContext->setResponse(
			GraphHelper::deleteUser(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$credentials["username"],
				$credentials["password"],
				$user
			)
		);
	}

	/**
	 * remove user from group
	 *
	 * @param string $groupId
	 * @param string $userId
	 * @param string|null $byUser
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function removeUserFromGroup(string $groupId, string $userId, ?string $byUser = null): ResponseInterface {
		$credentials = $this->getAdminOrUserCredentials($byUser);
		return GraphHelper::removeUserFromGroup(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials['username'],
			$credentials['password'],
			$userId,
			$groupId,
		);
	}

	/**
	 * @When /^the user "([^"]*)" deletes a user "([^"]*)" using the Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function theUserDeletesAUserUsingTheGraphAPI(string $byUser, string $user): void {
		$this->adminDeletesUserUsingTheGraphApi($user, $byUser);
	}

	/**
	 * @Given /^the user "([^"]*)" has deleted a user "([^"]*)" using the Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function theUserHasDeletesAUserUsingTheGraphAPI(string $byUser, string $user): void {
		$this->adminDeletesUserUsingTheGraphApi($user, $byUser);
		$this->featureContext->thenTheHTTPStatusCodeShouldBe(204);
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
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$response = $this->removeUserFromGroup($groupId, $userId);
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
	 * @param string|null $byUser
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function adminChangesPasswordOfUserToUsingTheGraphApi(
		string $user,
		string $password,
		?string $byUser = null
	): void {
		$credentials = $this->getAdminOrUserCredentials($byUser);
		$user = $this->featureContext->getActualUsername($user);
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$userId = $userId ?? $user;
		$response = GraphHelper::editUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials["username"],
			$credentials["password"],
			$userId,
			null,
			$password
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^the user "([^"]*)" resets the password of user "([^"]*)" to "([^"]*)" using the Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 * @param string $password
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theUserResetsThePasswordOfUserToUsingTheGraphApi(string $byUser, string $user, string $password) {
		$this->adminChangesPasswordOfUserToUsingTheGraphApi($user, $password, $byUser);
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
	 *
	 * @param string $user
	 * @param string|null $group
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function listSingleOrAllGroupsAlongWithAllMemberInformation(string $user, ?string $group = null): ResponseInterface {
		$credentials = $this->getAdminOrUserCredentials($user);

		return GraphHelper::getSingleOrAllGroupsAlongWithMembers(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials["username"],
			$credentials["password"],
			($group) ? $this->featureContext->getAttributeOfCreatedGroup($group, 'id') : null
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
	 * @Given /^the user "([^"]*)" has created a new user using the Graph API with the following settings:$/
	 *
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function theUserHasCreatedANewUserUsingGraphapiWithTheFollowingSettings(string $user, TableNode $table): void {
		$this->theUserCreatesNewUser(
			$user,
			$table
		);
		$rows = $table->getRowsHash();
		$response = $this->featureContext->getResponse();

		if ($response->getStatusCode() !== 200) {
			$this->throwHttpException($response, "Could not create user '$rows[userName]'");
		}
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
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
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
	 * @When the administrator tries to add user :user to a nonexistent group using the Graph API
	 * @When the user :byUser tries to add user :user to a nonexistent group using the Graph API
	 *
	 * @param string $user
	 * @param string|null $byUser
	 *
	 * @return void
	 *
	 * @throws GuzzleException | Exception
	 */
	public function theAdministratorTriesToAddUserToNonExistentGroupUsingTheGraphAPI(string $user, ?string $byUser = null): void {
		$credentials = $this->getAdminOrUserCredentials($byUser);
		$groupId = WebDavHelper::generateUUIDv4();
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$this->featureContext->setResponse(
			GraphHelper::addUserToGroup(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$credentials['username'],
				$credentials['password'],
				$userId,
				$groupId
			)
		);
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
	 * @Given /^the administrator has created a group "([^"]*)" using the Graph API$/
	 * @Given user :user has created a group :group using the Graph API
	 *
	 * @param string $group
	 * @param ?string $user
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userHasCreatedGroupUsingTheGraphApi(string $group, ?string $user = null): void {
		$response = $this->createGroup($group, $user);

		if ($response->getStatusCode() === 200) {
			$groupId = $this->featureContext->getJsonDecodedResponse($response)["id"];
			$this->featureContext->addGroupToCreatedGroupsList($group, true, true, $groupId);
		} else {
			$this->throwHttpException($response, "Could not create group '$group'.");
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
	 * @When user :user retrieves all groups along with their members using the Graph API
	 * @When user :user gets all the members information of group :group using the Graph API
	 *
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userRetrievesAllMemberInformationOfSingleOrAllGroups(string $user, string $group = ''): void {
		$this->featureContext->setResponse($this->listSingleOrAllGroupsAlongWithAllMemberInformation($user, $group));
	}

	/**
	 * @Then the group :group should have the following member information
	 *
	 * @param string $group
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theGroupShouldHaveTheFollowingMemberInformation(string $group, TableNode $table): void {
		$response = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse());
		$rows = $table->getHash();
		$currentMemberIndex = 0;
		if (isset($response['value'])) {
			$response = $response['value'];
			$groupFoundInResponse = false;
			foreach ($response as $value) {
				if ($value['displayName'] === $group) {
					$groupFoundInResponse = true;
					foreach ($rows as $row) {
						$this->checkUserInformation($row, $value['members'][$currentMemberIndex]);
						$currentMemberIndex++;
					}
					break;
				}
			}
			if (!$groupFoundInResponse) {
				throw new Error(
					'Group ' . $group . " could not be found in the response."
				);
			}
		} else {
			foreach ($rows as $row) {
				$this->checkUserInformation($row, $response['members'][$currentMemberIndex]);
				$currentMemberIndex++;
			}
		}
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
	public function userDeletesGroupUsingTheGraphApi(string $group, ?string $user = null): void {
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

	/**
	 * rename group name
	 *
	 * @param string $oldGroup
	 * @param string $newGroup
	 * @param string $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function renameGroup(string $oldGroup, string $newGroup, ?string $user = null): ResponseInterface {
		$credentials = $this->getAdminOrUserCredentials($user);
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($oldGroup, "id");

		return GraphHelper::updateGroup(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials['username'],
			$credentials['password'],
			$groupId,
			$newGroup
		);
	}

	/**
	 * @When user :user renames group :oldGroup to :newGroup using the Graph API
	 *
	 * @param string $user
	 * @param string $oldGroup
	 * @param string $newGroup
	 *
	 * @return void
	 */
	public function userRenamesGroupUsingTheGraphApi(string $user, string $oldGroup, string $newGroup): void {
		$this->featureContext->setResponse($this->renameGroup($oldGroup, $newGroup, $user));
	}

	/**
	 * @When the administrator removes the following users from the following groups using the Graph API
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function theAdministratorRemovesTheFollowingUsersFromTheFollowingGroupsUsingTheGraphApi(TableNode $table): void {
		$this->featureContext->verifyTableNodeColumns($table, ['username', 'groupname']);
		$usersGroups = $table->getColumnsHash();

		foreach ($usersGroups as $userGroup) {
			$groupId = $this->featureContext->getAttributeOfCreatedGroup($userGroup['groupname'], "id");
			$userId = $this->featureContext->getAttributeOfCreatedUser($userGroup['username'], "id");
			$this->featureContext->setResponse($this->removeUserFromGroup($groupId, $userId));
			$this->featureContext->pushToLastHttpStatusCodesArray();
		}
	}

	/**
	 * @When the administrator removes user :user from group :group using the Graph API
	 *
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 */
	public function theAdministratorTriesToRemoveUserFromGroupUsingTheGraphAPI(string $user, string $group): void {
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$this->featureContext->setResponse($this->removeUserFromGroup($groupId, $userId));
	}

	/**
	 * @When the administrator tries to remove user :user from group :group using the Graph API
	 * @When user :byUser tries to remove user :user from group :group using the Graph API
	 *
	 * @param string $user
	 * @param string $group
	 * @param string|null $byUser
	 *
	 * @return void
	 * @throws Exception | GuzzleException
	 */
	public function theUserTriesToRemoveAnotherUserFromGroupUsingTheGraphAPI(string $user, string $group, ?string $byUser = null): void {
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$this->featureContext->setResponse($this->removeUserFromGroup($groupId, $userId, $byUser));
	}

	/**
	 * @When the administrator tries to remove user :user from a nonexistent group using the Graph API
	 * @When user :byUser tries to remove user :user from a nonexistent group using the Graph API
	 *
	 * @param string $user
	 * @param string|null $byUser
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theUserTriesToRemoveAnotherUserFromNonExistentGroupUsingTheGraphAPI(string $user, ?string $byUser = null): void {
		$groupId = WebDavHelper::generateUUIDv4();
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$this->featureContext->setResponse($this->removeUserFromGroup($groupId, $userId, $byUser));
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function retrieveUserInformationUsingGraphApi(
		string $user
	):ResponseInterface {
		$credentials = $this->getAdminOrUserCredentials($user);
		return GraphHelper::getOwnInformationAndGroupMemberships(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials["username"],
			$credentials["password"],
		);
	}

	/**
	 * @When /^the user "([^"]*)" retrieves (her|his) information using the Graph API$/
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function userRetrievesHisorHerInformationOfUserUsingGraphApi(
		string $user
	):void {
		$response = $this->retrieveUserInformationUsingGraphApi($user);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Then /^the user retrieve API response should contain the following information:$/
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theUserRetrieveApiResponseShouldContainTheFollowingInformation(TableNode $table): void {
		$rows = $table->getHash();
		$apiResponse = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse());
		foreach ($rows as $row) {
			$this->checkUserInformation($row, $apiResponse);
		}
	}

	/**
	 * @param array $expectedValue
	 * @param array $actualValue
	 *
	 * @throws GuzzleException
	 * @return void
	 */
	public function checkUserInformation(array $expectedValue, array  $actualValue):void {
		foreach (array_keys($expectedValue) as $keyName) {
			switch ($keyName) {
				case "memberOf":
					$memberOfFromApiReponse = [];
					$memberOf = preg_split('/\s*,\s*/', trim($expectedValue['memberOf']));
					foreach ($actualValue['memberOf'] as $member) {
						$memberOfFromApiReponse[] = $member['displayName'];
					}
					Assert::assertEqualsCanonicalizing($memberOf, $memberOfFromApiReponse);
					break;
				case "id":
					if ($expectedValue[$keyName] !== '%uuid_v4%') {
						throw new Error(
							'Only UUIDv4 patterned user id can be checked' . ' but got '
							. trim($expectedValue[$keyName], '%')
						);
					}
					Assert::assertTrue(GraphHelper::isUUIDv4($actualValue['id']), __METHOD__ . ' Expected user_id to have UUIDv4 pattern but found: ' . $actualValue['id']);
					break;
				case "accountEnabled":
					if ($expectedValue[$keyName] === 'true') {
						Assert::assertTrue($actualValue[$keyName], ' Expected ' . $keyName . ' is not true ');
					} else {
						Assert::assertFalse($actualValue[$keyName], ' Expected ' . $keyName . ' is not false ');
					}
					break;
				default:
					Assert::assertEquals(
						$expectedValue[$keyName],
						$actualValue[$keyName],
						__METHOD__ .
						' Expected ' . $keyName . ' to have value ' . $expectedValue[$keyName]
						. ' but got ' . $actualValue[$keyName]
					);
					break;
			}
		}
	}

	/**
	 * @When user :byUser tries to get information of user :user using Graph API
	 * @When user :byUser gets information of user :user using Graph API
	 *
	 * @param string $byUser
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userTriesToGetInformationOfUser(string $byUser, string $user): void {
		$credentials = $this->getAdminOrUserCredentials($byUser);
		$response = GraphHelper::getUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials['username'],
			$credentials['password'],
			$user
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user tries to get all users using the Graph API
	 * @When user :user gets all users using the Graph API
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userGetsAllUserUsingTheGraphApi(string $user) {
		$credentials = $this->getAdminOrUserCredentials($user);
		$response = GraphHelper::getUsers(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials['username'],
			$credentials['password'],
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Then /^the API response should (not|)\s?contain following (user|users) with the information:$/
	 *
	 * @param string    $shouldOrNot   (not|)
	 * @param TableNode $table
	 *
	 * @throws Exception
	 * @return void
	 */
	public function theApiResponseShouldContainAllUserWithFollowingInformation(string $shouldOrNot, TableNode $table): void {
		$values = $table->getHash();
		$apiResponse = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse())['value'];
		foreach ($values as $expectedValue) {
			$found = false;
			foreach ($apiResponse as $key => $actualResponseValue) {
				if ($expectedValue["displayName"] === $actualResponseValue["displayName"]) {
					$found = true;
					$this->checkUserInformation($expectedValue, $actualResponseValue);
					unset($apiResponse[$key]);
					break;
				}
			}
			if ($shouldOrNot === 'not') {
				Assert::assertFalse($found, $expectedValue["displayName"] . ' has been found in the response, but should not be.');
			} else {
				Assert::assertTrue($found, $expectedValue["displayName"] . ' could not be found in the response.');
			}
		}
	}

	/**
	 * @param string $byUser
	 * @param string|null $user
	 *
	 * @return ResponseInterface
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function retrieveUserInformationAlongWithDriveUsingGraphApi(
		string $byUser,
		?string $user = null
	):ResponseInterface {
		$user = $user ?? $byUser;
		$credentials = $this->getAdminOrUserCredentials($user);
		return GraphHelper::getUserWithDriveInformation(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials["username"],
			$credentials["password"],
			$user
		);
	}

	/**
	 * @param string $byUser
	 * @param string|null $user
	 *
	 * @return ResponseInterface
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function retrieveUserInformationAlongWithGroupUsingGraphApi(
		string $byUser,
		?string $user = null
	):ResponseInterface {
		$user = $user ?? $byUser;
		$credentials = $this->getAdminOrUserCredentials($user);
		return GraphHelper::getUserWithGroupInformation(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$credentials["username"],
			$credentials["password"],
			$user
		);
	}

	/**
	 * @When /^the user "([^"]*)" gets user "([^"]*)" along with (his|her) drive information using Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 *
	 * @return void
	 */
	public function userTriesToGetInformationOfUserAlongWithHisDriveData(string $byUser, string $user) {
		$response = $this->retrieveUserInformationAlongWithDriveUsingGraphApi($byUser, $user);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^the user "([^"]*)" gets user "([^"]*)" along with (his|her) group information using Graph API$/
	 *
	 * @param string $byUser
	 * @param string $user
	 *
	 * @return void
	 */
	public function userTriesToGetInformationOfUserAlongWithHisGroup(string $byUser, string $user) {
		$response = $this->retrieveUserInformationAlongWithGroupUsingGraphApi($byUser, $user);
		$this->featureContext->setResponse($response);
	}

	/**
	 *
	 * @When /^the user "([^"]*)" gets (his|her) drive information using Graph API$/
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function userTriesToGetOwnDriveInformation(string $user) {
		$response = $this->retrieveUserInformationAlongWithDriveUsingGraphApi($user);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @param array $driveInformation
	 *
	 * @return string
	 */
	public static function getSpaceIdFromActualDriveinformation(array $driveInformation): string {
		return $driveInformation['id'];
	}

	/**
	 * check if single drive information is correct
	 *
	 * @param array $expectedDriveInformation
	 * @param array $actualDriveInformation
	 *
	 * @return void
	 */
	public function checkUserDriveInformation(array $expectedDriveInformation, array  $actualDriveInformation):void {
		foreach (array_keys($expectedDriveInformation) as $keyName) {
			$actualKeyValue = GraphHelper::separateAndGetValueForKey($keyName, $actualDriveInformation);
			switch ($expectedDriveInformation[$keyName]) {
				case '%user_id%':
					Assert::assertTrue(GraphHelper::isUUIDv4($actualKeyValue), __METHOD__ . ' Expected user_id to have UUIDv4 pattern but found: ' . $actualKeyValue);
					break;
				case '%space_id%':
					Assert::assertTrue(GraphHelper::isSpaceId($actualKeyValue), __METHOD__ . ' Expected space_id to have a UUIDv4:UUIDv4 pattern but found: ' . $actualKeyValue);
					break;
				default:
					$expectedDriveInformation[$keyName] = $this->featureContext->substituteInLineCodes(
						$expectedDriveInformation[$keyName],
						$this->featureContext->getCurrentUser(),
						[],
						[
							[
								// the actual space_id is substituted from the actual drive information rather than making an API request and substituting
								"code" => "%space_id%",
								"function" =>
									[$this, "getSpaceIdFromActualDriveinformation"],
								"parameter" => [$actualDriveInformation]
							],
						]
					);
					Assert::assertEquals($expectedDriveInformation[$keyName], $actualKeyValue);
			}
		}
	}

	/**
	 * @param TableNode $table
	 *
	 * @Then the user retrieve API response should contain the following drive information:
	 *
	 * @return void
	 */
	public function theResponseShouldContainTheFollowingDriveInformation(TableNode $table): void {
		$expectedDriveInformation = $table->getRowsHash();
		// array of user drive information (Personal Drive Information Only)
		$actualDriveInformation = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse());
		if (\is_array($actualDriveInformation) && \array_key_exists('drive', $actualDriveInformation)) {
			$this->checkUserDriveInformation($expectedDriveInformation, $actualDriveInformation['drive']);
		} else {
			throw new Error('Response is not an array or the array does not consist key "drive"');
		}
	}

	/**
	 * add multiple users in a group at once
	 *
	 * @param string $user
	 * @param array $userIds
	 * @param string $groupId
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function addMultipleUsersToGroup(string $user, array $userIds, string $groupId, TableNode $table): void {
		$credentials = $this->getAdminOrUserCredentials($user);
		$this->featureContext->verifyTableNodeColumns($table, ['username']);

		$this->featureContext->setResponse(
			GraphHelper::addUsersToGroup(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$credentials["username"],
				$credentials["password"],
				$groupId,
				$userIds
			)
		);
	}

	/**
	 * @When /^the administrator "([^"]*)" adds the following users to a group "([^"]*)" at once using the Graph API$/
	 *
	 * @param string $user
	 * @param string $group
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function theAdministratorAddsTheFollowingUsersToAGroupInASingleRequestUsingTheGraphApi(string $user, string $group, TableNode $table): void {
		$userIds = [];
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		foreach ($table->getHash() as $row) {
			$userIds[] = $this->featureContext->getAttributeOfCreatedUser($row['username'], "id");
		}
		$this->addMultipleUsersToGroup($user, $userIds, $groupId, $table);
	}

	/**
	 * @When /^the administrator "([^"]*)" tries to add the following users to a nonexistent group at once using the Graph API$/
	 *
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theAdministratorTriesToAddsTheFollowingUsersToANonExistingGroupAtOnceUsingTheGraphApi(string $user, TableNode $table): void {
		$userIds = [];
		$groupId = WebDavHelper::generateUUIDv4();
		foreach ($table->getHash() as $row) {
			$userIds[] = $this->featureContext->getAttributeOfCreatedUser($row['username'], "id");
		}
		$this->addMultipleUsersToGroup($user, $userIds, $groupId, $table);
	}

	/**
	 * @When /^the administrator "([^"]*)" tries to add the following nonexistent users to a group "([^"]*)" at once using the Graph API$/
	 *
	 * @param string $user
	 * @param string $group
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theAdministratorTriesToAddTheFollowingNonExistingUsersToAGroupAtOnceUsingTheGraphApi(string $user, string $group, TableNode $table): void {
		$userIds = [];
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		foreach ($table->getHash() as $row) {
			$userIds[] = WebDavHelper::generateUUIDv4();
		}
		$this->addMultipleUsersToGroup($user, $userIds, $groupId, $table);
	}

	/**
	 * @When /^the administrator "([^"]*)" tries to add the following users to a group "([^"]*)" at once using the Graph API$/
	 * @When /^the administrator "([^"]*)" tries to add the following existent and nonexistent users to a group "([^"]*)" at once using the Graph API$/
	 *
	 * @param string $user
	 * @param string $group
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theAdministratorTriesToAddTheFollowingUsersToAGroupAtOnceUsingTheGraphApi(string $user, string $group, TableNode $table): void {
		$userIds = [];
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		foreach ($table->getHash() as $row) {
			$userId = $this->featureContext->getAttributeOfCreatedUser($row['username'], "id");
			$userIds[] = $userId ? $userId : WebDavHelper::generateUUIDv4();
		}
		$this->addMultipleUsersToGroup($user, $userIds, $groupId, $table);
	}

	/**
	 * @When user :user gets all applications using the Graph API
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userGetsAllApplicationsUsingTheGraphApi(string $user) {
		$response = GraphHelper::getApplications(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Then the user retrieve API response should contain the following applications information:
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function theResponseShouldContainTheFollowingApplicationsInformation(TableNode $table): void {
		Assert::assertIsArray($responseArray = ($this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse()))['value'][0]);
		foreach ($table->getHash() as $row) {
			$key = $row["key"];
			if ($key === 'id') {
				Assert::assertTrue(
					GraphHelper::isUUIDv4($responseArray[$key]),
					__METHOD__ . ' Expected id to have UUIDv4 pattern but found: ' . $row["value"]
				);
			} else {
				Assert::assertEquals($responseArray[$key], $row["value"]);
			}
		}
	}

	/**
	 * @Then the user retrieve API response should contain the following app roles:
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function theResponseShouldContainTheFollowingAppRolesInformation(TableNode $table): void {
		Assert::assertIsArray($responseArray = ($this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse()))['value'][0]);
		foreach ($table->getRows() as $row) {
			$foundRoleInResponse = false;
			foreach ($responseArray['appRoles'] as $role) {
				if ($role['displayName'] === $row[0]) {
					$foundRoleInResponse = true;
					break;
				}
			}
			Assert::assertTrue($foundRoleInResponse, "the response does not contain the role $row[0]");
		}
	}

	/**
	 * @When the user :user gets all users of the group :group using the Graph API
	 *
	 * @param string $user
	 * @param string $group
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userGetsAllUsersOfTheGroupUsingTheGraphApi(string $user, string $group) {
		$groupId = $this->featureContext->getGroupIdByGroupName($group);
		$response = GraphHelper::getUsersWithFilterMemberOf(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$groupId
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When the user :user gets all users of two groups :groups using the Graph API
	 *
	 * @param string $user
	 * @param string $groups
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userGetsAllUsersOfTwoGroupsUsingTheGraphApi(string $user, string $groups) {
		$groupsIdArray = [];
		foreach (explode(',', $groups) as $group) {
			array_push($groupsIdArray, $this->featureContext->getGroupIdByGroupName($group));
		}
		$response = GraphHelper::getUsersOfTwoGroups(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$groupsIdArray
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * Get roleId by role name
	 *
	 * @param string $role
	 *
	 * @return string
	 * @throws GuzzleException
	 */
	public function getRoleIdByRoleName(string $role): string {
		$response = GraphHelper::getApplications(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword()
		);
		$responseData = \json_decode($response->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
		if (isset($responseData["value"][0]["appRoles"])) {
			foreach ($responseData["value"][0]["appRoles"] as $value) {
				if ($value["displayName"] === $role) {
					return $value["id"];
				}
			}
			throw new Exception(__METHOD__ . " role with name $role not found");
		}
	}

	/**
	 * @When the user :user gets all users with role :role using the Graph API
	 *
	 * @param string $user
	 * @param string $role
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userGetsAllUsersWithRoleUsingTheGraphApi(string $user, string $role) {
		$response = GraphHelper::getUsersWithFilterRoleAssignment(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$this->getRoleIdByRoleName($role)
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When the user :user gets all users with role :role and member of the group :group using the Graph API
	 *
	 * @param string $user
	 * @param string $role
	 * @param string $group
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userGetsAllUsersWithRoleAndMemberOfGroupUsingTheGraphApi(string $user, string $role, string $group) {
		$response = GraphHelper::getUsersWithFilterRolesAssignmentAndMemberOf(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$this->getRoleIdByRoleName($role),
			$this->featureContext->getGroupIdByGroupName($group)
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given /^the administrator has assigned the role "([^"]*)" to user "([^"]*)" using the Graph API$/
	 *
	 * @param string $role
	 * @param string $user
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function theAdministratorHasGivenTheRoleUsingTheGraphApi(string $role, string $user): void {
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id') ?? $user;

		if (empty($this->appEntity)) {
			$this->setApplicationEntity();
		}

		$response = GraphHelper::assignRole(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$this->appEntity["appRoles"][$role],
			$this->appEntity["id"],
			$userId
		);
		Assert::assertEquals(
			201,
			$response->getStatusCode(),
			__METHOD__
			. "\nExpected status code '200' but got '" . $response->getStatusCode() . "'"
		);
	}

	/**
	 * @When /^the administrator retrieves the assigned role of user "([^"]*)" using the Graph API$/
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userRetrievesAssignedRoleUsingTheGraphApi(string $user): void {
		$admin = $this->featureContext->getAdminUserName();
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id') ?? $user;
		$this->featureContext->setResponse(
			GraphHelper::getAssignedRole(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$admin,
				$this->featureContext->getPasswordForUser($admin),
				$userId
			)
		);
	}

	/**
	 * set application Entity in global variable
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function setApplicationEntity(): void {
		$applicationEntity = (
			$this->featureContext->getJsonDecodedResponse(
				GraphHelper::getApplications(
					$this->featureContext->getBaseUrl(),
					$this->featureContext->getStepLineRef(),
					$this->featureContext->getAdminUsername(),
					$this->featureContext->getAdminPassword(),
				)
			)
		)['value'][0];
		$this->appEntity["id"] = $applicationEntity["id"];
		foreach ($applicationEntity["appRoles"] as $value) {
			$this->appEntity["appRoles"][$value['displayName']] = $value['id'];
		}
	}

	/**
	 * @Then /^the Graph API response should have the role "([^"]*)"$/
	 *
	 * @param string $role
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function theGraphApiResponseShouldHaveTheRole(string $role): void {
		$response = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse())['value'][0];
		if (empty($this->appEntity)) {
			$this->setApplicationEntity();
		}
		Assert::assertEquals(
			$this->appEntity["appRoles"][$role],
			$response['appRoleId'],
			__METHOD__
			. "\nExpected rolId for role '$role'' to be '" . $this->appEntity["appRoles"][$role] . "' but got '" . $response['appRoleId'] . "'"
		);
	}

	/**
	 * @When user :user gets details of the group :groupName using the Graph API
	 *
	 * @param string $user
	 * @param string $groupName
	 *
	 * @return void
	 */
	public function userGetsDetailsOfTheGroupUsingTheGraphApi(string $user, string $groupName): void {
		$credentials = $this->getAdminOrUserCredentials($user);

		$this->featureContext->setResponse(
			GraphHelper::getGroup(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$credentials["username"],
				$credentials["password"],
				$groupName
			)
		);
	}
}
