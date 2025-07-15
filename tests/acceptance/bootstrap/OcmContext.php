<?php declare(strict_types=1);
/**
 * @author Viktor Scharf <scharf.vi@gmail.com>
 *
 * @copyright Copyright (c) 2024, ownCloud GmbH
 * @license AGPL-3.0
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3,
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License, version 3,
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 *
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\OcmHelper;
use TestHelpers\WebDavHelper;
use TestHelpers\BehatHelper;
use TestHelpers\HttpRequestHelper;
use Behat\Gherkin\Node\TableNode;

/**
 * Acceptance test steps related to testing federation share(ocm) features
 */
class OcmContext implements Context {
	private FeatureContext $featureContext;
	private SpacesContext $spacesContext;
	private ArchiverContext $archiverContext;
	private SharingNgContext $sharingNgContext;
	private string $invitationToken;
	private array $acceptedUsers = ["LOCAL" => [], "REMOTE" => []];

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
		$this->spacesContext = BehatHelper::getContext($scope, $environment, 'SpacesContext');
		$this->archiverContext = BehatHelper::getContext($scope, $environment, 'ArchiverContext');
		$this->sharingNgContext = BehatHelper::getContext($scope, $environment, 'SharingNgContext');
	}

	/**
	 * @return string
	 * @throws Exception
	 */
	public function getLastFederatedInvitationToken(): string {
		if (empty($this->invitationToken)) {
			throw new Exception(__METHOD__ . " token not found");
		}
		return $this->invitationToken;
	}

	/**
	 * @param string $user
	 * @param string|null $email
	 * @param string|null $description
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function createInvitation(
		string $user,
		?string $email = null,
		?string $description = null,
	): ResponseInterface {
		$response = OcmHelper::createInvitation(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$email,
			$description,
		);
		$responseData = \json_decode($response->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
		if (isset($responseData["token"])) {
			$this->invitationToken = $responseData["token"];
		}
		return $response;
	}

	/**
	 * @When :user creates the federation share invitation
	 * @When :user creates the federation share invitation with email :email and description :description
	 *
	 * @param string $user
	 * @param string|null $email
	 * @param string|null $description
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userCreatesTheFederationShareInvitation(
		string $user,
		?string $email = null,
		?string $description = null,
	): void {
		$this->featureContext->setResponse($this->createInvitation($user, $email, $description));
	}

	/**
	 * @Given :user has created the federation share invitation
	 * @Given :user has created the federation share invitation with email :email and description :description
	 *
	 * @param string $user
	 * @param string|null $email
	 * @param string|null $description
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userHasCreatedTheFederationShareInvitation(
		string $user,
		?string $email = null,
		?string $description = null,
	): void {
		$response = $this->createInvitation($user, $email, $description);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @param string $user
	 * @param string|null $token
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function acceptInvitation(string $user, ?string $token = null): ResponseInterface {
		$providerDomain = $this->featureContext->getLocalBaseUrlWithoutScheme();
		if ($this->featureContext->getCurrentServer() === "LOCAL") {
			$providerDomain = $this->featureContext->getRemoteBaseUrlWithoutScheme();
		}
		return OcmHelper::acceptInvitation(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$token ? $token : $this->getLastFederatedInvitationToken(),
			$providerDomain,
		);
	}

	/**
	 * @When :user accepts the last federation share invitation
	 * @When :user tries to accept the last federation share invitation
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userAcceptsTheLastFederationShareInvitation(string $user): void {
		$this->featureContext->setResponse($this->acceptInvitation($user));
	}

	/**
	 * @When :user tries to accept the invitation with invalid token
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userTriesToAcceptInvitationWithInvalidToken(string $user): void {
		$this->featureContext->setResponse($this->acceptInvitation($user, WebDavHelper::generateUUIDv4()));
	}

	/**
	 * @Given :user has accepted invitation
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userHasAcceptedTheLastFederationShareInvitation(string $user): void {
		$response = $this->acceptInvitation($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @When :user tries to accept the federation share invitation from same instance
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function triesToAcceptTheFederationShareInvitationFromSameInstance(string $user): void {
		$providerDomain = $this->featureContext->getLocalBaseUrlWithoutScheme();
		$token = $this->getLastFederatedInvitationToken();
		$this->featureContext->setResponse(
			OcmHelper::acceptInvitation(
				$this->featureContext->getBaseUrl(),
				$user,
				$this->featureContext->getPasswordForUser($user),
				$token,
				$providerDomain,
			),
		);
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function findAcceptedUsers(string $user): ResponseInterface {
		$currentServer = $this->featureContext->getCurrentServer();
		$response = OcmHelper::findAcceptedUsers(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
		);
		if ($response->getStatusCode() === 200) {
			$users = $this->featureContext->getJsonDecodedResponse($response);
			$this->acceptedUsers[$currentServer] = \array_merge($this->acceptedUsers[$currentServer], $users);
			$response->getBody()->rewind();
		}
		return $response;
	}

	/**
	 * @When :user searches for accepted users
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userFindsAcceptedUsers(string $user): void {
		$this->featureContext->setResponse($this->findAcceptedUsers($user));
	}

	/**
	 *
	 * @param string $user
	 * @param string $ocmUserName
	 *
	 * @return array
	 * @throws GuzzleException
	 */
	public function getAcceptedUserByName(string $user, string $ocmUserName): array {
		$currentServer = $this->featureContext->getCurrentServer();
		$displayName = $this->featureContext->getUserDisplayName($ocmUserName);
		$acceptedUsers = $this->acceptedUsers[$currentServer];
		foreach ($acceptedUsers as $acceptedUser) {
			if ($acceptedUser["display_name"] === $displayName) {
				return $acceptedUser;
			}
		}
		// fetch the accepted users
		$response = $this->findAcceptedUsers($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			200,
			"failed to list accepted users by '$user'",
			$response,
		);
		$users = ($this->featureContext->getJsonDecodedResponse($response));
		foreach ($users as $acceptedUser) {
			if ($acceptedUser["display_name"] === $displayName) {
				return $acceptedUser;
			}
		}
		throw new \Exception("Could not find user with name '{$ocmUserName}' in the accepted users list.");
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function listInvitations(string $user): ResponseInterface {
		return OcmHelper::listInvite(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
		);
	}

	/**
	 * @When :user lists the created invitations
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userListsCreatedInvitations(string $user): void {
		$this->featureContext->setResponse($this->listInvitations($user));
	}

	/**
	 * @When the user waits :number seconds for the invitation token to expire
	 *
	 * @param int $number
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theUserWaitsForTokenToExpire(int $number): void {
		\sleep($number);
	}

	/**
	 * @When user :user deletes federated connection with user :ocmUser using the Graph API
	 * @When user :user tries to delete federated connection with user :ocmUser using the Graph API
	 *
	 * @param string $user
	 * @param string $ocmUser
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userDeletesFederatedConnectionWithUserUsingTheGraphApi(string $user, string $ocmUser): void {
		$this->featureContext->setResponse($this->deleteConnection($user, $ocmUser));
	}

	/**
	 * @When user :user has deleted federated connection with user :ocmUser
	 *
	 * @param string $user
	 * @param string $ocmUser
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userHasDeletedFederatedConnectionWithUser(string $user, string $ocmUser): void {
		$response = $this->deleteConnection($user, $ocmUser);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			200,
			"failed while deleting connection with user $ocmUser",
			$response,
		);
	}

	/**
	 * @param string $user
	 * @param string $ocmUser
	 * @param string|null $idp
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function deleteConnection(string $user, string $ocmUser, ?string $idp = null): ResponseInterface {
		$ocmUser = $this->getAcceptedUserByName($user, $ocmUser);
		$ocmUser['idp'] = $idp ?? $ocmUser['idp'];
		return OcmHelper::deleteConnection(
			$this->featureContext->getBaseUrl(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$ocmUser['user_id'],
			$ocmUser['idp'],
		);
	}

	/**
	 * @Then user :user should be able to download federated shared file :resource
	 *
	 * @param string $user
	 * @param string $resource
	 *
	 * @return void
	 */
	public function userShouldBeAbleToDownloadFederatedSharedFile(string $user, string $resource): void {
		$remoteItemId = $this->spacesContext->getSharesRemoteItemId($user, $resource);
		$baseUrl = $this->featureContext->getRemoteBaseUrl();
		$davPath = WebDavHelper::getDavPath($this->featureContext->getDavPathVersion());
		$response = HttpRequestHelper::get(
			"$baseUrl/$davPath/$remoteItemId",
			$user,
			$this->featureContext->getPasswordForUser($user),
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "Failed to download resource $resource", $response);
	}

	/**
	 * @Then user :user should be able to download archive of federated shared folder :resource
	 *
	 * @param string $user
	 * @param string $resource
	 *
	 * @return void
	 */
	public function userShouldBeAbleToDownloadArchiveOfFederatedSharedFolder(string $user, string $resource): void {
		$queryString = $this->archiverContext->getArchiverQueryString($user, $resource, 'remoteItemIds');
		$response = HttpRequestHelper::get(
			$this->archiverContext->getArchiverUrl($queryString),
			$user,
			$this->featureContext->getPasswordForUser($user),
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			200,
			"Failed to download archive of resource $resource",
			$response,
		);
	}

	/**
	 * @When user :user sends PROPFIND request to federated share :share with depth :folderDepth using the WebDAV API
	 *
	 * @param string $user
	 * @param string $share
	 * @param string $folderDepth
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userSendsPropfindRequestToFederatedShareWithDepthUsingTheWebdavApi(
		string $user,
		string $share,
		string $folderDepth,
	): void {
		$response = $this->spacesContext->sendPropfindRequestToSpace($user, "", $share, null, $folderDepth, true);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Then user :sharee should have the following federated shares:
	 *
	 * @param string $sharee
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function userShouldHaveTheFollowingFederatedShares(string $sharee, TableNode $table): void {
		$shares = $table->getColumnsHash();
		foreach ($shares as $share) {
			$this->sharingNgContext->checkIfShareExists(
				$share["resource"],
				$sharee,
				$share["sharer"],
				'',
				true,
				true,
				$share["permissionsRole"],
			);
		}
	}

	/**
	 * @When user :user tries to delete federated connection with user :ocmUser and provider :idp using the Graph API
	 *
	 * @param string $user
	 * @param string $ocmUser
	 * @param string $idp
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userDeletesFederatedConnectionWithUserAndProviderUsingTheGraphApi(
		string $user,
		string $ocmUser,
		string $idp,
	): void {
		$this->featureContext->setResponse($this->deleteConnection($user, $ocmUser, $idp));
	}
}
