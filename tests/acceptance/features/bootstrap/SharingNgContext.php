<?php declare(strict_types=1);
/**
 * @author Viktor Scharf <scharf.vi@gmail.com>
 *
 * @copyright Copyright (c) 2023, ownCloud GmbH
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
use TestHelpers\GraphHelper;
use TestHelpers\WebDavHelper;
use Behat\Gherkin\Node\TableNode;
use PHPUnit\Framework\Assert;

require_once 'bootstrap.php';

/**
 * Acceptance test steps related to testing sharing ng features
 */
class SharingNgContext implements Context {
	private FeatureContext $featureContext;
	private SpacesContext $spacesContext;

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
	}

	/**
	 * @param string $user
	 * @param TableNode|null $body
	 *
	 * @return ResponseInterface
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function createLinkShare(string $user, TableNode $body): ResponseInterface {
		$bodyRows = $body->getRowsHash();
		$space = $bodyRows['space'];
		$resource = $bodyRows['resource'];

		$spaceId = ($this->spacesContext->getSpaceByName($user, $space))["id"];
		$itemId = $this->spacesContext->getResourceId($user, $space, $resource);

		$bodyRows['displayName'] = $bodyRows['displayName'] ?? null;
		$bodyRows['expirationDateTime'] = $bodyRows['expirationDateTime'] ?? null;
		$bodyRows['password'] = $bodyRows['password'] ?? null;
		$body = [
			'type' => $bodyRows['permissionsRole'],
			'displayName' => $bodyRows['displayName'],
			'expirationDateTime' => $bodyRows['expirationDateTime'],
			'password' => $this->featureContext->getActualPassword($bodyRows['password'])
		];

		return GraphHelper::createLinkShare(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$spaceId,
			$itemId,
			\json_encode($body)
		);
	}

	/**
	 * @param string $user
	 * @param string $fileOrFolder   (file|folder)
	 * @param string $space
	 * @param string $resource
	 *
	 * @return ResponseInterface
	 * @throws Exception
	 */
	public function getPermissionsList(string $user, string $fileOrFolder, string $space, ?string $resource = ''):ResponseInterface {
		$spaceId = ($this->spacesContext->getSpaceByName($user, $space))["id"];

		if ($fileOrFolder === 'folder') {
			$itemId = $this->spacesContext->getResourceId($user, $space, $resource);
		} else {
			$itemId = $this->spacesContext->getFileId($user, $space, $resource);
		}

		return GraphHelper::getPermissionsList(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$spaceId,
			$itemId
		);
	}

	/**
	 * @When /^user "([^"]*)" gets permissions list for (folder|file) "([^"]*)" of the space "([^"]*)" using the Graph API$/
	 *
	 * @param string $user
	 * @param string $fileOrFolder   (file|folder)
	 * @param string $resource
	 * @param string $space
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userGetsPermissionsListForResourceOfTheSpaceUsingTheGraphiAPI(string $user, string $fileOrFolder, string $resource, string $space):void {
		$this->featureContext->setResponse(
			$this->getPermissionsList($user, $fileOrFolder, $space, $resource)
		);
	}

	/**
	 * @When /^user "([^"]*)" lists the permissions of space "([^"]*)" using permissions endpoint of the Graph API$/
	 *
	 * @param string $user
	 * @param string $space
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userListsThePermissionsOfSpaceUsingTheGraphApi($user, $space):void {
		$this->featureContext->setResponse(
			$this->getPermissionsList($user, 'folder', $space)
		);
	}

	/**
	 * @When /^user "([^"]*)" tries to list the permissions of space "([^"]*)" owned by "([^"]*)" using permissions endpoint of the Graph API$/
	 *
	 * @param string $user
	 * @param string $space
	 * @param string $spaceOwner
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userTriesToListThePermissionsOfSpaceUsingPermissionsEndpointOfTheGraphApi(string $user, string $space, string $spaceOwner):void {
		$spaceId = ($this->spacesContext->getSpaceByName($spaceOwner, $space))["id"];
		$itemId = $this->spacesContext->getResourceId($spaceOwner, $space, '');

		$this->featureContext->setResponse(
			GraphHelper::getPermissionsList(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$user,
				$this->featureContext->getPasswordForUser($user),
				$spaceId,
				$itemId
			)
		);
	}

	/**
	 * @param string $user
	 * @param TableNode $table
	 * @param string|null $fileId
	 *
	 * @return ResponseInterface
	 *
	 * @throws JsonException
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function sendShareInvitation(string $user, TableNode $table, string $fileId = null): ResponseInterface {
		$rows = $table->getRowsHash();
		if ($rows['space'] === 'Personal' || $rows['space'] === 'Shares') {
			$space = $this->spacesContext->getSpaceByName($user, $rows['space']);
		} else {
			$space = $this->spacesContext->getCreatedSpace($rows['space']);
		}
		$spaceId = $space['id'];

		// $fileId is used for trying to share deleted files
		if ($fileId) {
			$itemId = $fileId;
		} else {
			$resource = $rows['resource'] ?? '';

			// for a disabled and deleted space, resource id is not accessible, so get resource id from the saved response
			if ($resource === '' && !\in_array($rows['space'], ['Personal', 'Shares'])) {
				$itemId = $space['fileId'];
			} else {
				$itemId = $this->spacesContext->getResourceId($user, $rows['space'], $resource);
			}
		}

		if (\array_key_exists('shareeId', $rows)) {
			$shareeIds[] = $rows['shareeId'];
			$shareTypes[] = $rows['shareType'];
		} else {
			$sharees = array_map('trim', explode(',', $rows['sharee']));
			$shareTypes = array_map('trim', explode(',', $rows['shareType']));

			foreach ($sharees as $index => $sharee) {
				$shareType = $shareTypes[$index];
				$shareeId = "";
				if ($shareType === "user") {
					$shareeId = $this->featureContext->getAttributeOfCreatedUser($sharee, 'id');
				} elseif ($shareType === "group") {
					$shareeId = $this->featureContext->getAttributeOfCreatedGroup($sharee, 'id');
				}
				// for non-existing group or user, generate random id
				$shareeIds[] = $shareeId ?: WebDavHelper::generateUUIDv4();
			}
		}

		$permissionsRole = $rows['permissionsRole'] ?? null;
		$permissionsAction = $rows['permissionsAction'] ?? null;
		$expireDate = $rows["expireDate"] ?? null;

		$response = GraphHelper::sendSharingInvitation(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$spaceId,
			$itemId,
			$shareeIds,
			$shareTypes,
			$permissionsRole,
			$permissionsAction,
			$expireDate
		);
		if ($response->getStatusCode() === 200) {
			$this->featureContext->shareNgAddToCreatedUserGroupShares($response);
		}
		return $response;
	}

	/**
	 * @Given /^user "([^"]*)" has sent the following share invitation:$/
	 *
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userHasSentTheFollowingShareInvitation(string $user, TableNode $table): void {
		$response = $this->sendShareInvitation($user, $table);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "", $response);
	}

	/**
	 * @When /^user "([^"]*)" sends the following share invitation using the Graph API:$/
	 * @When /^user "([^"]*)" tries to send the following share invitation using the Graph API:$/
	 * @When user :user sends the following share invitation for space using the Graph API:
	 *
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userSendsTheFollowingShareInvitationUsingTheGraphApi(string $user, TableNode $table): void {
		$this->featureContext->setResponse(
			$this->sendShareInvitation($user, $table)
		);
	}

	/**
	 * @When user :user updates the last share with the following using the Graph API:
	 *
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function userUpdatesTheLastShareWithFollowingUsingGraphApi($user, TableNode $table) {
		$response = $this->featureContext->shareNgGetLastCreatedUserGroupShare();
		$permissionID = json_decode($response->getBody()->getContents())->value[0]->id;
		$this->featureContext->setResponse(
			$this->updateResourceShare(
				$user,
				$table,
				$permissionID
			)
		);
	}

	/**
	 * @param string $user
	 * @param TableNode $body
	 * @param string $permissionID
	 *
	 * @return ResponseInterface
	 */
	public function updateResourceShare(string $user, TableNode  $body, string $permissionID): ResponseInterface {
		$bodyRows = $body->getRowsHash();
		$space = $bodyRows['space'];
		$resource = $bodyRows['resource'];
		$spaceId = ($this->spacesContext->getSpaceByName($user, $space))["id"];
		$itemId = $this->spacesContext->getResourceId($user, $space, $resource);
		$body = [];

		if (\array_key_exists('permissionsRole', $bodyRows)) {
			$body['roles'] = [GraphHelper::getPermissionsRoleIdByName($bodyRows['permissionsRole'])];
		}

		if (\array_key_exists('expirationDateTime', $bodyRows)) {
			$body['expirationDateTime'] = empty($bodyRows['expirationDateTime']) ? null : $bodyRows['expirationDateTime'];
		}

		return GraphHelper::updateShare(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$spaceId,
			$itemId,
			\json_encode($body),
			$permissionID
		);
	}

	/**
	 * @When user :user sends the following share invitation with file-id :fileId using the Graph API:
	 *
	 * @param string $user
	 * @param string $fileId
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userSendsTheFollowingShareInvitationWithFileIdUsingTheGraphApi(string $user, string $fileId, TableNode $table): void {
		$this->featureContext->setResponse(
			$this->sendShareInvitation($user, $table, $fileId)
		);
	}

	/**
	 * @When /^user "([^"]*)" creates the following link share using the Graph API:$/
	 *
	 * @param string $user
	 * @param TableNode|null $body
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userCreatesAPublicLinkShareWithSettings(string $user, TableNode  $body):void {
		$response = $this->createLinkShare($user, $body);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given /^user "([^"]*)" has created the following link share:$/
	 *
	 * @param string $user
	 * @param TableNode|null $body
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userHasCreatedTheFollowingLinkShare(string $user, TableNode  $body): void {
		$response = $this->createLinkShare($user, $body);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "Failed while creating public share link!", $response);
		$this->featureContext->shareNgAddToCreatedLinkShares($response);
	}

	/**
	 * @When /^user "([^"]*)" updates the last public link share using the Graph API with$/
	 *
	 * @param string $user
	 * @param TableNode|null $body
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userUpdatesLastPublicLinkShareUsingTheGraphApiWith(string $user, TableNode  $body):void {
		$this->featureContext->setResponse(
			$this->updateLinkShare(
				$user,
				$body,
				$this->featureContext->shareNgGetLastCreatedLinkShareID()
			)
		);
	}

	/**
	 * @param string $user
	 * @param TableNode $body
	 * @param string $permissionID
	 *
	 * @return ResponseInterface
	 */
	public function updateLinkShare(string $user, TableNode  $body, string $permissionID): ResponseInterface {
		$bodyRows = $body->getRowsHash();
		$space = $bodyRows['space'];
		$resource = $bodyRows['resource'];
		$spaceId = ($this->spacesContext->getSpaceByName($user, $space))['id'];
		$itemId = $this->spacesContext->getResourceId($user, $space, $resource);
		$body = [];

		if (\array_key_exists('permissionsRole', $bodyRows)) {
			$body['link']['type'] = $bodyRows['permissionsRole'];
		}

		if (\array_key_exists('expirationDateTime', $bodyRows)) {
			$body['expirationDateTime'] = empty($bodyRows['expirationDateTime']) ? null : $bodyRows['expirationDateTime'];
		}

		return GraphHelper::updateShare(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$spaceId,
			$itemId,
			\json_encode($body),
			$permissionID
		);
	}

	/**
	 * @When user :user sets the following password for the last link share using the Graph API:
	 *
	 * @param string $user
	 * @param TableNode|null $body
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userSetsOrUpdatesFollowingPasswordForLastLinkShareUsingTheGraphApi(string $user, TableNode  $body):void {
		$bodyRows = $body->getRowsHash();
		$space = $bodyRows['space'];
		$resource = $bodyRows['resource'];
		$spaceId = ($this->spacesContext->getSpaceByName($user, $space))["id"];
		$itemId = $this->spacesContext->getResourceId($user, $space, $resource);

		if (\array_key_exists('password', $bodyRows)) {
			$body = [
				"password" => $this->featureContext->getActualPassword($bodyRows['password']),
			];
		} else {
			throw new Error('Password is missing to set for share link!');
		}

		$response = GraphHelper::setLinkSharePassword(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$spaceId,
			$itemId,
			\json_encode($body),
			$this->featureContext->shareNgGetLastCreatedLinkShareID()
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @param string $sharer
	 * @param string $shareType (user|group|link)
	 * @param string $space
	 * @param string|null $resource
	 * @param string|null $recipient
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function removeAccessToSpaceItem(
		string $sharer,
		string $shareType,
		string $space,
		?string $resource = null,
		?string $recipient = null
	): ResponseInterface {
		$spaceId = ($this->spacesContext->getSpaceByName($sharer, $space))["id"];
		$itemId = (isset($resource)) ? $this->spacesContext->getResourceId($sharer, $space, $resource) : $this->spacesContext->getResourceId($sharer, $space, $space);

		$permId = ($shareType === 'link')
			? $this->featureContext->shareNgGetLastCreatedLinkShareID()
			: $this->featureContext->shareNgGetLastCreatedUserGroupShareID();
		return
			GraphHelper::removeAccessToSpaceItem(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$sharer,
				$this->featureContext->getPasswordForUser($sharer),
				$spaceId,
				$itemId,
				$permId
			);
	}

	/**
	 * @param string $sharer
	 * @param string $shareType (user|group|link)
	 * @param string $space
	 * @param string|null $recipient
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function removeAccessToSpace(
		string $sharer,
		string $shareType,
		string $space,
		?string $recipient = null
	): ResponseInterface {
		$spaceId = ($this->spacesContext->getSpaceByName($sharer, $space))["id"];

		$permId = ($shareType === 'link')
			? $this->featureContext->shareNgGetLastCreatedLinkShareID()
			: $this->featureContext->shareNgGetLastCreatedUserGroupShareID();
		return
			GraphHelper::removeAccessToSpace(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$sharer,
				$this->featureContext->getPasswordForUser($sharer),
				$spaceId,
				$permId
			);
	}

	/**
	 * @When /^user "([^"]*)" removes the access of (user|group) "([^"]*)" from (?:file|folder|resource) "([^"]*)" of space "([^"]*)" using the Graph API$/
	 *
	 * @param string $sharer
	 * @param string $recipientType (user|group)
	 * @param string $recipient can be both user or group
	 * @param string $resource
	 * @param string $space
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userRemovesAccessOfUserOrGroupFromResourceOfSpaceUsingGraphAPI(
		string $sharer,
		string $recipientType,
		string $recipient,
		string $resource,
		string $space
	): void {
		$this->featureContext->setResponse(
			$this->removeAccessToSpaceItem($sharer, $recipientType, $space, $resource)
		);
	}

	/**
	 * @When /^user "([^"]*)" removes the access of (user|group) "([^"]*)" from space "([^"]*)" using permissions endpoint of the Graph API$/
	 *
	 * @param string $sharer
	 * @param string $recipientType (user|group)
	 * @param string $recipient can be both user or group
	 * @param string $space
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userRemovesAccessOfUserOrGroupFromSpaceUsingPermissionsEndpointOfGraphAPI(
		string $sharer,
		string $recipientType,
		string $recipient,
		string $space
	): void {
		$this->featureContext->setResponse(
			$this->removeAccessToSpaceItem($sharer, $recipientType, $space)
		);
	}

	/**
	 * @When /^user "([^"]*)" removes the link of (?:file|folder) "([^"]*)" from space "([^"]*)" using the Graph API$/
	 *
	 * @param string $sharer
	 * @param string $resource
	 * @param string $space
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userRemovesSharePermissionOfAResourceInLinkShareUsingGraphAPI(
		string $sharer,
		string $resource,
		string $space
	):void {
		$this->featureContext->setResponse(
			$this->removeAccessToSpaceItem($sharer, 'link', $space, $resource)
		);
	}

	/**
	 * @When /^user "([^"]*)" removes the access of (user|group) "([^"]*)" from space "([^"]*)" using root endpoint of the Graph API$/
	 *
	 * @param string $sharer
	 * @param string $recipientType (user|group)
	 * @param string $recipient can be both user or group
	 * @param string $space
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userRemovesAccessOfUserOrGroupFromSpaceUsingGraphAPI(
		string $sharer,
		string $recipientType,
		string $recipient,
		string $space
	): void {
		$this->featureContext->setResponse(
			$this->removeAccessToSpace($sharer, $recipientType, $space)
		);
	}

	/**
	 * @When /^user "([^"]*)" removes the link from space "([^"]*)" using root endpoint of the Graph API$/
	 *
	 * @param string $sharer
	 * @param string $space
	 *
	 * @return void
	 * @throws JsonException
	 * @throws GuzzleException
	 */
	public function userRemovesLinkFromSpaceUsingRootEndpointOfGraphAPI(
		string $sharer,
		string $space
	):void {
		$this->featureContext->setResponse(
			$this->removeAccessToSpace($sharer, 'link', $space)
		);
	}

	/**
	 * @Then /^for user "([^"]*)" the space Shares should (not|)\s?contain these (files|entries):$/
	 *
	 * @param string $user
	 * @param string $shouldOrNot
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function forUserTheSpaceSharesShouldContainTheseEntries(string $user, string $shouldOrNot, TableNode $table): void {
		$should = $shouldOrNot !== 'not';
		$rows = $table->getRows();
		$response = GraphHelper::getSharesSharedWithMe(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user)
		);
		$contents = \json_decode($response->getBody()->getContents(), true);

		$fileFound = empty(array_diff(array_map(fn ($row) => trim($row[0], '/'), $rows), array_column($contents['value'], 'name')));

		$assertMessage = $should
			? "Response does not contain the entry."
			: "Response does contain the entry but should not.";

		Assert::assertSame($should, $fileFound, $assertMessage);
	}

	/**
	 * @When user :user disables sync of share :share using the Graph API
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userDisablesSyncOfShareUsingTheGraphApi(string $user):void {
		$shareItemId = $this->featureContext->shareNgGetLastCreatedUserGroupShareID();
		$shareSpaceId = FeatureContext::SHARES_SPACE_ID;
		$itemId = $shareSpaceId . '!' . $shareItemId;
		$response = GraphHelper::disableShareSync(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			$itemId,
			$shareSpaceId,
		);
		$this->featureContext->setResponse($response);
		$this->featureContext->pushToLastStatusCodesArrays();
	}

	/**
	 * @When user :user enables sync of share :share offered by :offeredBy from :space space using the Graph API
	 *
	 * @param string $user
	 * @param string $share
	 * @param string $offeredBy
	 * @param string $space
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userEnablesSyncOfShareUsingTheGraphApi(string $user, string $share, string $offeredBy, string $space):void {
		$share = ltrim($share, '/');
		$itemId = $this->spacesContext->getResourceId($offeredBy, $space, $share);
		$shareSpaceId = FeatureContext::SHARES_SPACE_ID;
		$response =  GraphHelper::enableShareSync(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			$itemId,
			$shareSpaceId
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user tries to enable share sync of a resource :resource using the Graph API
	 *
	 * @param string $user
	 * @param string $resource
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function userTriesToEnableShareSyncOfResourceUsingTheGraphApi(string $user, string $resource):void {
		$shareSpaceId = FeatureContext::SHARES_SPACE_ID;
		$itemId = ($resource === 'nonexistent') ? WebDavHelper::generateUUIDv4() : '';

		$response =  GraphHelper::enableShareSync(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			$itemId,
			$shareSpaceId
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Then /^user "([^"]*)" should have sync (enabled|disabled) for share "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $status
	 * @param string $resource
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userShouldHaveSyncEnabledOrDisabledForShare(string $user, string $status, string $resource):void {
		$response = GraphHelper::getSharesSharedWithMe(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user)
		);
		$responseBody = $this->featureContext->getJsonDecodedResponse($response);
		$expectedValue = $status === "enabled" ? "true" : "false";
		$actualValue = "";
		foreach ($responseBody["value"] as $value) {
			if ($value["remoteItem"]["name"] === $resource) {
				// var_export converts values to their string representations
				// e.g.: true -> 'true'
				$actualValue = var_export($value["@client.synchronize"], true);
				break;
			}
		}
		Assert::assertSame(
			$actualValue,
			$expectedValue,
			"Expected property '@client.synchronize' to be '$expectedValue' but found '$actualValue'"
		);
	}

	/**
	 * @Then user :user should be able to send share invitation with all allowed permission roles
	 *
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userShouldBeAbleToSendShareInvitationWithAllAllowedPermissionRoles(string $user, TableNode $table): void {
		$listPermissionResponse = $this->featureContext->getJsonDecodedResponseBodyContent();
		if (!isset($listPermissionResponse->{'@libre.graph.permissions.roles.allowedValues'})) {
			Assert::fail(
				"The following response does not contain '@libre.graph.permissions.roles.allowedValues' property:\n" . $listPermissionResponse
			);
		}
		Assert::assertNotEmpty(
			$listPermissionResponse->{'@libre.graph.permissions.roles.allowedValues'},
			"'@libre.graph.permissions.roles.allowedValues' should not be empty"
		);
		$allowedPermissionRoles = $listPermissionResponse->{'@libre.graph.permissions.roles.allowedValues'};
		// this info is needed for log to see which roles allowed and which were not when tests fail
		$shareInvitationRequestResult = "From the given allowed role lists from the permissions:\n";
		$areAllSendInvitationSuccessFullForAllowedRoles = true;
		$rows = $table->getRowsHash();
		// when sending share invitation for a project space, the resource to be shared is project space itself. So resource can be put as empty
		$resource = $rows['resource'] ?? '';
		$shareType = $rows['shareType'];
		$space = $rows['space'];
		//this details is needed for result logging purpose to determine whether the resource shared is a resource or a project space
		$resourceDetail = ($resource) ? "resource '" . $resource : "space '" . $space;
		foreach ($allowedPermissionRoles as $role) {
			//we should be able to send share invitation for each of the role allowed for the files/folders which are  listed in permissions (allowed)
			$roleAllowed = GraphHelper::getPermissionNameByPermissionRoleId($role->id);
			$responseSendInvitation = $this->sendShareInvitation($user, new TableNode(array_merge($table->getTable(), [['permissionsRole', $roleAllowed]])));
			$jsonResponseSendInvitation = $this->featureContext->getJsonDecodedResponseBodyContent($responseSendInvitation);
			$httpsStatusCode = $responseSendInvitation->getStatusCode();
			if ($httpsStatusCode === 200 && !empty($jsonResponseSendInvitation->value)) {
				// remove the share so that the same user can be share for the next allowed roles
				$removePermissionsResponse = $this->removeAccessToSpaceItem($user, $shareType, $space, $resource);
				Assert::assertEquals(204, $removePermissionsResponse->getStatusCode());
			} else {
				$areAllSendInvitationSuccessFullForAllowedRoles = false;
				$shareInvitationRequestResult .= "\tShare invitation for " . $resourceDetail . "' with role '" . $roleAllowed . "' failed and was not allowed.\n";
			}
		}
		Assert::assertTrue($areAllSendInvitationSuccessFullForAllowedRoles, $shareInvitationRequestResult);
	}

	/**
	 * @When /^user "([^"]*)" (?:tries to list|lists) the permissions of space "([^"]*)" using root endpoint of the Graph API$/
	 *
	 * @param string $user
	 * @param string $space
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 *
	 */
	public function userListsThePermissionsOfDriveUsingRootEndPointOFTheGraphApi(string $user, string $space):void {
		$spaceId = ($this->spacesContext->getSpaceByName($user, $space))["id"];

		$response = GraphHelper::getDrivePermissionsList(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$spaceId
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^user "([^"]*)" (?:tries to send|sends) the following share invitation using root endpoint of the Graph API:$/
	 *
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userSendsTheFollowingShareInvitationUsingRootEndPointTheGraphApi(string $user, TableNode $table):void {
		$shareeIds = [];
		$rows = $table->getRowsHash();
		if ($rows['space'] === 'Personal' || $rows['space'] === 'Shares') {
			$space = $this->spacesContext->getSpaceByName($user, $rows['space']);
		} else {
			$space = $this->spacesContext->getCreatedSpace($rows['space']);
		}
		$spaceId = $space['id'];

		$sharees = array_map('trim', explode(',', $rows['sharee']));
		$shareTypes = array_map('trim', explode(',', $rows['shareType']));

		foreach ($sharees as $index => $sharee) {
			$shareType = $shareTypes[$index];
			if ($sharee === "") {
				// set empty value to $shareeIds
				$shareeIds[] = "";
				continue;
			}
			$shareeId = "";
			if ($shareType === "user") {
				$shareeId = $this->featureContext->getAttributeOfCreatedUser($sharee, 'id');
			} elseif ($shareType === "group") {
				$shareeId = $this->featureContext->getAttributeOfCreatedGroup($sharee, 'id');
			}
			// for non-existing group or user, generate random id
			$shareeIds[] = $shareeId ?: WebDavHelper::generateUUIDv4();
		}

		$permissionsRole = $rows['permissionsRole'] ?? null;
		$permissionsAction = $rows['permissionsAction'] ?? null;
		$expireDate = $rows["expireDate"] ?? null;

		$response = GraphHelper::sendSharingInvitationForDrive(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$spaceId,
			$shareeIds,
			$shareTypes,
			$permissionsRole,
			$permissionsAction,
			$expireDate
		);

		$this->featureContext->setResponse($response);
	}
}
