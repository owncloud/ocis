<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Kiran Parajuli <kiran@jankaritech.com>
 * @copyright Copyright (c) 2021 Kiran Parajuli kiran@jankaritech.com
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use GuzzleHttp\Exception\GuzzleException;
use TestHelpers\GraphHelper;
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
	public function before(BeforeScenarioScope $scope):void {
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
	 * @return array
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
	): array {
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
		$this->featureContext->theHTTPStatusCodeShouldBeSuccess();
		$response = GraphHelper::getUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$requester,
			$requesterPassword,
			$userId
		);
		$this->featureContext->setResponse($response);
		$this->featureContext->theHTTPStatusCodeShouldBeSuccess();
		return $this->featureContext->getJsonDecodedResponse();
	}

	/**
	 * @param string $user
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function adminHasRetrievedUserUsingTheGraphApi(string $user):void {
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
	):void {
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
	 * @param string $group
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function adminDeletesGroupUsingTheGraphApi(
		string $group
	) {
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		if ($groupId) {
			$this->featureContext->setResponse(
				GraphHelper::deleteGroup(
					$this->featureContext->getBaseUrl(),
					$this->featureContext->getStepLineRef(),
					$this->featureContext->getAdminUsername(),
					$this->featureContext->getAdminPassword(),
					$groupId
				)
			);
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
	public function adminHasDeletedGroupUsingTheGraphApi(string $group):void {
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
	public function adminDeletesUserUsingTheGraphApi(string $user) {
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
	public function adminHasRemovedUserFromGroupUsingTheGraphApi(string $user, string $group):void {
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
	public function userShouldNotBeMemberInGroupUsingTheGraphApi(string $user, string $group):void {
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
	public function userShouldBeMemberInGroupUsingTheGraphApi(string $user, string $group):void {
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
	):void {
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
	 * @return array
	 * @throws Exception
	 */
	public function adminHasRetrievedGroupListUsingTheGraphApi():array {
		$response =  GraphHelper::getGroups(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword()
		);
		if ($response->getStatusCode() === 200) {
			return $this->featureContext->getJsonDecodedResponse($response);
		} else {
			try {
				$jsonBody = $this->featureContext->getJsonDecodedResponse($response);
				throw new Exception(
					__METHOD__
					. "\nCould not retrieve groups list."
					. "\nHTTP status code: " . $response->getStatusCode()
					. "\nError code: " . $jsonBody["error"]["code"]
					. "\nMessage: " . $jsonBody["error"]["message"]
				);
			} catch (TypeError $e) {
				throw new Exception(
					__METHOD__
					. "\nCould not retrieve groups list."
					. "\nHTTP status code: " . $response->getStatusCode()
					. "\nResponse body: " . $response->getBody()
				);
			}
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
	public function theAdminHasRetrievedMembersListOfGroupUsingTheGraphApi(string $group):array {
		$response = GraphHelper::getMembersList(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$this->featureContext->getAttributeOfCreatedGroup($group, 'id')
		);
		if ($response->getStatusCode() === 200) {
			return $this->featureContext->getJsonDecodedResponse($response);
		} else {
			try {
				$jsonBody = $this->featureContext->getJsonDecodedResponse($response);
				throw new Exception(
					__METHOD__
					. "\nCould not retrieve members list for group $group."
					. "\nHTTP status code: " . $response->getStatusCode()
					. "\nError code: " . $jsonBody["error"]["code"]
					. "\nMessage: " . $jsonBody["error"]["message"]
				);
			} catch (TypeError $e) {
				throw new Exception(
					__METHOD__
					. "\nCould not retrieve members list for group $group."
					. "\nHTTP status code: " . $response->getStatusCode()
					. "\nResponse body: " . $response->getBody()
				);
			}
		}
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
	 * @return array
	 * @throws Exception
	 */
	public function theAdminHasCreatedUser(
		string $user,
		string $password,
		string $email,
		string $displayName
	): array {
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
			try {
				$jsonResponseBody = $this->featureContext->getJsonDecodedResponse($response);
				throw new Exception(
					__METHOD__
					. "\nCould not create user $user"
					. "\nError code: {$jsonResponseBody['error']['code']}"
					. "\nError message: {$jsonResponseBody['error']['message']}"
				);
			} catch (TypeError $e) {
				throw new Exception(
					__METHOD__
					. "\nCould not create user $user"
					. "\nHTTP status code: " . $response->getStatusCode()
					. "\nResponse body: " . $response->getBody()
				);
			}
		} else {
			return $this->featureContext->getJsonDecodedResponse($response);
		}
	}

	/**
	 * adds a user to a group
	 *
	 * @param string $user
	 * @param string $group
	 * @param bool $checkResult
	 *
	 * @return void
	 * @throws JsonException
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function adminHasAddedUserToGroupUsingTheGraphApi(
		string $user,
		string $group,
		bool $checkResult = true
	) {
		$groupId = $this->featureContext->getAttributeOfCreatedGroup($group, "id");
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, "id");
		$result = GraphHelper::addUserToGroup(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$userId,
			$groupId
		);
		if ($checkResult && ($result->getStatusCode() !== 204)) {
			throw new Exception(
				__METHOD__
				. "\nCould not add user to group. "
				. "\n HTTP status: " . $result->getStatusCode()
				. "\n Response body: " . $result->getBody()
			);
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
	public function adminHasCreatedGroupUsingTheGraphApi(string $group):array {
		$result = GraphHelper::createGroup(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$group,
		);
		if ($result->getStatusCode() === 200) {
			return $this->featureContext->getJsonDecodedResponse($result);
		} else {
			try {
				$jsonBody = $this->featureContext->getJsonDecodedResponse($result);
				throw new Exception(
					__METHOD__
					. "\nError: failed creating group '$group'"
					. "\nStatus code: " . $jsonBody['error']['code']
					. "\nMessage: " . $jsonBody['error']['message']
				);
			} catch (TypeError $e) {
				throw new Exception(
					__METHOD__
					. "\nError: failed creating group '$group'"
					. "\nHTTP status code: " . $result->getStatusCode()
					. "\nResponse body: " . $result->getBody()
				);
			}
		}
	}
}
