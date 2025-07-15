<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Artur Neumann <artur@jankaritech.com>
 * @copyright Copyright (c) 2018 Artur Neumann artur@jankaritech.com
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License,
 * as published by the Free Software Foundation;
 * either version 3 of the License, or any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>
 *
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Behat\Gherkin\Node\TableNode;
use GuzzleHttp\Exception\GuzzleException;
use PHPUnit\Framework\Assert;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\HttpRequestHelper;
use TestHelpers\WebDavHelper;
use TestHelpers\BehatHelper;

require_once 'bootstrap.php';

/**
 * context containing API steps needed for the locking mechanism of webdav
 */
class WebDavLockingContext implements Context {
	private FeatureContext $featureContext;
	private PublicWebDavContext $publicWebDavContext;
	private SpacesContext $spacesContext;

	/**
	 *
	 * @var string[][]
	 */
	private array $tokenOfLastLock = [];

	/**
	 *
	 * @param string $user
	 * @param string $file
	 * @param TableNode $properties table with no heading with | property | value |
	 * @param string|null $fullUrl
	 * @param boolean $public if the file is in a public share or not
	 * @param boolean $expectToSucceed
	 * @param string|null $spaceId
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	private function lockFile(
		string $user,
		string $file,
		TableNode $properties,
		?string $fullUrl = null,
		bool $public = false,
		bool $expectToSucceed = true,
		?string $spaceId = null,
	): ResponseInterface {
		$user = $this->featureContext->getActualUsername($user);
		if ($public === true) {
			$type = "public-files";
			$password = $this->featureContext->getActualPassword("%public%");
		} else {
			$type = "files";
			$password = $this->featureContext->getPasswordForUser($user);
		}
		$body
			= "<?xml version='1.0' encoding='UTF-8'?>" .
			"<d:lockinfo xmlns:d='DAV:'> ";
		$headers = [];
		// depth is only 0 or infinity. We don't need to set it more, as there is no lock for the folder
		$this->featureContext->verifyTableNodeRows($properties, [], ['lockscope', 'timeout']);
		$propertiesRows = $properties->getRowsHash();

		foreach ($propertiesRows as $property => $value) {
			if ($property === "timeout") {
				// properties that are set in the header not in the xml
				$headers[$property] = $value;
			} else {
				$body .= "<d:$property><d:$value/></d:$property>";
			}
		}
		$body .= "</d:lockinfo>";

		if (isset($fullUrl)) {
			$response = HttpRequestHelper::sendRequest(
				$fullUrl,
				"LOCK",
				$this->featureContext->getActualUsername($user),
				$this->featureContext->getPasswordForUser($user),
				$headers,
				$body,
			);
		} else {
			$baseUrl = $this->featureContext->getBaseUrl();
			$response = WebDavHelper::makeDavRequest(
				$baseUrl,
				$user,
				$password,
				"LOCK",
				$file,
				$headers,
				$spaceId,
				$body,
				$this->featureContext->getDavPathVersion(),
				$type,
			);
		}

		$responseXmlObject = HttpRequestHelper::getResponseXml($response, __METHOD__);
		$xmlPart = $responseXmlObject->xpath("//d:locktoken/d:href");
		if (isset($xmlPart[0])) {
			$this->tokenOfLastLock[$user][$file] = (string) $xmlPart[0];
		} else {
			if ($expectToSucceed === true) {
				Assert::fail("could not find lock token after trying to lock '$file'");
			}
		}
		return $response;
	}

	/**
	 * @When user :user locks file :file using the WebDAV API setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 */
	public function userLocksFileSettingPropertiesUsingWebDavAPI(
		string $user,
		string $file,
		TableNode $properties,
	): void {
		$spaceId = null;
		if (\str_starts_with($file, "Shares/")
			&& $this->featureContext->getDavPathVersion() === WebDavHelper::DAV_VERSION_SPACES
		) {
			$spaceId = $this->spacesContext->getSpaceIdByName($user, "Shares");
			$file = \str_replace("Shares/", "", $file);
		}
		$response = $this->lockFile($user, $file, $properties, null, false, true, $spaceId);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user tries to lock file/folder :file using the WebDAV API setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 */
	public function userTriesToLockFileSettingPropertiesUsingWebDavAPI(
		string $user,
		string $file,
		TableNode $properties,
	): void {
		$response = $this->lockFile($user, $file, $properties, null, false, false);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user locks file :file inside the space :space using the WebDAV API setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $space
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userLocksFileInProjectSpaceUsingWebDavAPI(
		string $user,
		string $file,
		string $space,
		TableNode $properties,
	): void {
		$this->featureContext->setResponse($this->userLocksFileInProjectSpace($user, $file, $space, $properties));
	}

	/**
	 * @param string $user
	 * @param string $file
	 * @param string $space
	 * @param TableNode $properties
	 *
	 * @return ResponseInterface|null
	 *
	 * @throws GuzzleException
	 */
	public function userLocksFileInProjectSpace(
		string $user,
		string $file,
		string $space,
		TableNode $properties,
	): ?ResponseInterface {
		$spaceId = $this->spacesContext->getSpaceIdByName($user, $space);
		$baseUrl = $this->featureContext->getBaseUrl();
		$davPathVersion = $this->featureContext->getDavPathVersion();
		$suffixPath = $user;
		if ($davPathVersion === WebDavHelper::DAV_VERSION_SPACES) {
			$suffixPath = $spaceId;
		}

		$davPath = WebDavHelper::getDavPath($davPathVersion, $suffixPath);
		$fullUrl = "$baseUrl/$davPath/$file";
		return $this->lockFile($user, $file, $properties, $fullUrl, false, true, $spaceId);
	}

	/**
	 * @Given user :user has locked file :file inside the space :space setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $space
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userHasLockedFileInProjectSpaceUsingWebDavAPI(
		string $user,
		string $file,
		string $space,
		TableNode $properties,
	): void {
		$response = $this->userLocksFileInProjectSpace($user, $file, $space, $properties);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @When user :user tries to lock file :file inside the space :space using the WebDAV API setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $space
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 */
	public function userTriesToLockFileInProjectSpaceUsingWebDavAPI(
		string $user,
		string $file,
		string $space,
		TableNode $properties,
	): void {
		$spaceId = $this->spacesContext->getSpaceIdByName($user, $space);
		$davPathVersion = $this->featureContext->getDavPathVersion();
		$suffixPath = $user;
		if ($davPathVersion === WebDavHelper::DAV_VERSION_SPACES) {
			$suffixPath = $spaceId;
		}

		$davPath = WebdavHelper::getDavPath($davPathVersion, $suffixPath);
		$fullUrl = $this->featureContext->getBaseUrl() . "/$davPath/$file";
		$response = $this->lockFile($user, $file, $properties, $fullUrl, false, false, $spaceId);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user locks file :file using file-id :fileId using the WebDAV API setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $fileId
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 */
	public function userLocksFileUsingFileIdUsingWebDavAPISettingFollowingProperties(
		string $user,
		string $file,
		string $fileId,
		TableNode $properties,
	): void {
		$davPath = WebdavHelper::getDavPath($this->featureContext->getDavPathVersion());
		$davPath = \rtrim($davPath, '/');
		$fullUrl = $this->featureContext->getBaseUrl() . "/$davPath/$fileId";
		$response = $this->lockFile($user, $file, $properties, $fullUrl);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user tries to lock file :file using file-id :fileId using the WebDAV API setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $fileId
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 */
	public function userTriesToLockFileUsingFileIdUsingWebDavAPI(
		string $user,
		string $file,
		string $fileId,
		TableNode $properties,
	): void {
		$davPath = WebdavHelper::getDavPath($this->featureContext->getDavPathVersion());
		$davPath = \rtrim($davPath, '/');
		$fullUrl = $this->featureContext->getBaseUrl() . "/$davPath/$fileId";
		$response = $this->lockFile($user, $file, $properties, $fullUrl, false, false);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given user :user has locked file :file setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 */
	public function userHasLockedFile(string $user, string $file, TableNode $properties): void {
		$response = $this->lockFile($user, $file, $properties);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @Given user :user has locked file :file inside space :spaceName setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $spaceName
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 */
	public function userHasLockedFileInsideSpaceSettingTheFollowingProperties(
		string $user,
		string $file,
		string $spaceName,
		TableNode $properties,
	): void {
		$spaceId = $this->spacesContext->getSpaceIdByName($this->featureContext->getActualUsername($user), $spaceName);
		$response = $this->lockFile($user, $file, $properties, null, false, true, $spaceId);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @Given user :user has locked file :file using file-id :fileId setting the following properties
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $fileId
	 * @param TableNode $properties table with no heading with | property | value |
	 *
	 * @return void
	 */
	public function userHasLockedFileUsingFileId(
		string $user,
		string $file,
		string $fileId,
		TableNode $properties,
	): void {
		$davPath = WebdavHelper::getDavPath($this->featureContext->getDavPathVersion());
		$davPath = \rtrim($davPath, '/');
		$fullUrl = $this->featureContext->getBaseUrl() . "/$davPath/$fileId";
		$response = $this->lockFile($user, $file, $properties, $fullUrl);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @Given the public has locked the last public link shared file/folder setting the following properties
	 *
	 * @param TableNode $properties
	 *
	 * @return void
	 */
	public function publicHasLockedLastSharedFile(TableNode $properties): void {
		$response = $this->lockFile(
			$this->featureContext->getLastCreatedPublicShareToken(),
			"/",
			$properties,
			null,
			true,
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @When the public locks the last public link shared file using the WebDAV API setting the following properties
	 * @When the public tries to lock the last public link shared file using the WebDAV API setting the following properties
	 *
	 * @param TableNode $properties
	 *
	 * @return void
	 */
	public function publicLocksLastSharedFile(TableNode $properties): void {
		$token = ($this->featureContext->isUsingSharingNG())
		? $this->featureContext->shareNgGetLastCreatedLinkShareToken()
		: $this->featureContext->getLastCreatedPublicShareToken();
		$response = $this->lockFile(
			$token,
			"/",
			$properties,
			null,
			true,
			false,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given the public has locked :file in the last public link shared folder setting the following properties
	 *
	 * @param string $file
	 * @param TableNode $properties
	 *
	 * @return void
	 */
	public function publicHasLockedFileLastSharedFolder(
		string $file,
		TableNode $properties,
	): void {
		$token = ($this->featureContext->isUsingSharingNG())
		? $this->featureContext->shareNgGetLastCreatedLinkShareToken()
		: $this->featureContext->getLastCreatedPublicShareToken();
		$response = $this->lockFile(
			$token,
			$file,
			$properties,
			null,
			true,
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @When /^the public locks "([^"]*)" in the last public link shared folder using the public WebDAV API setting the following properties$/
	 * @When /^the public tries to lock "([^"]*)" in the last public link shared folder using the public WebDAV API setting the following properties$/
	 *
	 * @param string $file
	 * @param TableNode $properties
	 *
	 * @return void
	 */
	public function publicLocksFileLastSharedFolder(
		string $file,
		TableNode $properties,
	): void {
		$token = ($this->featureContext->isUsingSharingNG())
		? $this->featureContext->shareNgGetLastCreatedLinkShareToken()
		: $this->featureContext->getLastCreatedPublicShareToken();
		$response = $this->lockFile(
			$token,
			$file,
			$properties,
			null,
			true,
			false,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user unlocks the last created lock of file :file using the WebDAV API
	 *
	 * @param string $user
	 * @param string $file
	 *
	 * @return void
	 */
	public function unlockLastLockUsingWebDavAPI(string $user, string $file): void {
		$response = $this->unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
			$user,
			$file,
			$user,
			$file,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user unlocks the last created lock of file :file inside space :spaceName using the WebDAV API
	 *
	 * @param string $user
	 * @param string $spaceName
	 * @param string $file
	 *
	 * @return void
	 */
	public function userUnlocksTheLastCreatedLockOfFileInsideSpaceUsingTheWebdavApi(
		string $user,
		string $spaceName,
		string $file,
	): void {
		$spaceId = $this->spacesContext->getSpaceIdByName($this->featureContext->getActualUsername($user), $spaceName);
		$response = $this->unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
			$user,
			$file,
			$user,
			$file,
			false,
			null,
			$spaceId,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user unlocks the last created lock of file :itemToUnlock using file-id :fileId using the WebDAV API
	 *
	 * @param string $user
	 * @param string $itemToUnlock
	 * @param string $fileId
	 *
	 * @return void
	 */
	public function userUnlocksTheLastCreatedLockOfFileWithFileIdPathUsingTheWebdavApi(
		string $user,
		string $itemToUnlock,
		string $fileId,
	): void {
		$davPath = WebdavHelper::getDavPath($this->featureContext->getDavPathVersion());
		$davPath = \rtrim($davPath, '/');
		$fullUrl = $this->featureContext->getBaseUrl() . "/$davPath/$fileId";
		$response = $this->unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
			$user,
			$itemToUnlock,
			$user,
			$itemToUnlock,
			false,
			$fullUrl,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user unlocks file :itemToUnlock with the last created lock of file :itemToUseLockOf using the WebDAV API
	 *
	 * @param string $user
	 * @param string $itemToUnlock
	 * @param string $itemToUseLockOf
	 *
	 * @return void
	 */
	public function unlockItemWithLastLockOfOtherItemUsingWebDavAPI(
		string $user,
		string $itemToUnlock,
		string $itemToUseLockOf,
	): void {
		$response = $this->unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
			$user,
			$itemToUnlock,
			$user,
			$itemToUseLockOf,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user unlocks file :itemToUnlock with the last created public lock of file :itemToUseLockOf using the WebDAV API
	 *
	 * @param string $user
	 * @param string $itemToUnlock
	 * @param string $itemToUseLockOf
	 *
	 * @return void
	 */
	public function unlockItemWithLastPublicLockOfOtherItemUsingWebDavAPI(
		string $user,
		string $itemToUnlock,
		string $itemToUseLockOf,
	): void {
		$lockOwner = $this->featureContext->getLastCreatedPublicShareToken();
		$response = $this->unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
			$user,
			$itemToUnlock,
			$lockOwner,
			$itemToUseLockOf,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 *
	 * @param string $user
	 * @param string $itemToUnlock
	 *
	 * @return int
	 *
	 * @throws Exception|GuzzleException
	 */
	private function countLockOfResources(
		string $user,
		string $itemToUnlock,
	): int {
		$user = $this->featureContext->getActualUsername($user);
		$baseUrl = $this->featureContext->getBaseUrl();
		$password = $this->featureContext->getPasswordForUser($user);
		$body
			= "<?xml version='1.0' encoding='UTF-8'?>" .
			"<d:propfind xmlns:d='DAV:'> " .
			"<d:prop><d:lockdiscovery/></d:prop>" .
			"</d:propfind>";
		$response = WebDavHelper::makeDavRequest(
			$baseUrl,
			$user,
			$password,
			"PROPFIND",
			$itemToUnlock,
			null,
			null,
			$body,
			$this->featureContext->getDavPathVersion(),
		);
		$responseXmlObject = HttpRequestHelper::getResponseXml($response, __METHOD__);
		$xmlPart = $responseXmlObject->xpath("//d:response//d:lockdiscovery/d:activelock");
		if (\is_array($xmlPart)) {
			return \count($xmlPart);
		}
		throw new Exception("xmlPart for 'd:activelock' was expected to be array but found: $xmlPart");
	}

	/**
	 * @Given user :user has unlocked file :itemToUnlock with the last created lock of file :itemToUseLockOf of user :lockOwner using the WebDAV API
	 *
	 * @param string $user
	 * @param string $itemToUnlock
	 * @param string $lockOwner
	 * @param string $itemToUseLockOf
	 * @param boolean $public
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function hasUnlockItemWithTheLastCreatedLock(
		string $user,
		string $itemToUnlock,
		string $lockOwner,
		string $itemToUseLockOf,
		bool $public = false,
	): void {
		$lockCount = $this->countLockOfResources($user, $itemToUnlock);

		$response = $this->unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
			$user,
			$itemToUnlock,
			$lockOwner,
			$itemToUseLockOf,
			$public,
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(204, "", $response);

		$lockCountAfterUnlock = $this->countLockOfResources($user, $itemToUnlock);

		Assert::assertEquals(
			$lockCount - 1,
			$lockCountAfterUnlock,
			"Expected $lockCount lock(s) for '$itemToUnlock' but found '$lockCount'",
		);
	}

	/**
	 * @param string $user
	 * @param string $itemToUnlock
	 * @param string $lockOwner
	 * @param string $itemToUseLockOf
	 * @param boolean $public
	 * @param string|null $fullUrl
	 * @param string|null $spaceId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
		string $user,
		string $itemToUnlock,
		string $lockOwner,
		string $itemToUseLockOf,
		bool $public = false,
		?string $fullUrl = null,
		?string $spaceId = null,
	): ResponseInterface {
		$user = $this->featureContext->getActualUsername($user);
		$lockOwner = $this->featureContext->getActualUsername($lockOwner);
		if ($public === true) {
			$type = "public-files";
			$password = $this->featureContext->getActualPassword("%public%");
		} else {
			$type = "files";
			$password = $this->featureContext->getPasswordForUser($user);
		}
		$baseUrl = $this->featureContext->getBaseUrl();
		if (!isset($this->tokenOfLastLock[$lockOwner][$itemToUseLockOf])) {
			Assert::fail(
				"could not find saved token of '$itemToUseLockOf' " .
				"owned by user '$lockOwner'",
			);
		}
		$headers = [
			"Lock-Token" => $this->tokenOfLastLock[$lockOwner][$itemToUseLockOf],
		];
		if (isset($fullUrl)) {
			$response = HttpRequestHelper::sendRequest(
				$fullUrl,
				"UNLOCK",
				$this->featureContext->getActualUsername($user),
				$this->featureContext->getPasswordForUser($user),
				$headers,
			);
		} else {
			$response = WebDavHelper::makeDavRequest(
				$baseUrl,
				$user,
				$password,
				"UNLOCK",
				$itemToUnlock,
				$headers,
				$spaceId,
				null,
				$this->featureContext->getDavPathVersion(),
				$type,
			);
		}
		return $response;
	}

	/**
	 * @When user :user unlocks file :itemToUnlock with the last created lock of file :itemToUseLockOf of user :lockOwner using the WebDAV API
	 *
	 * @param string $user
	 * @param string $itemToUnlock
	 * @param string $lockOwner
	 * @param string $itemToUseLockOf
	 *
	 * @return void
	 */
	public function userUnlocksItemWithLastLockOfUserAndItemUsingWebDavAPI(
		string $user,
		string $itemToUnlock,
		string $lockOwner,
		string $itemToUseLockOf,
	): void {
		$response = $this->unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
			$user,
			$itemToUnlock,
			$lockOwner,
			$itemToUseLockOf,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When the public unlocks file :itemToUnlock with the last created lock of file :itemToUseLockOf of user :lockOwner using the WebDAV API
	 *
	 * @param string $itemToUnlock
	 * @param string $lockOwner
	 * @param string $itemToUseLockOf
	 *
	 * @return void
	 */
	public function unlockItemAsPublicWithLastLockOfUserAndItemUsingWebDavAPI(
		string $itemToUnlock,
		string $lockOwner,
		string $itemToUseLockOf,
	): void {
		$token = ($this->featureContext->isUsingSharingNG())
		? $this->featureContext->shareNgGetLastCreatedLinkShareToken()
		: $this->featureContext->getLastCreatedPublicShareToken();
		$response = $this->unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
			$token,
			$itemToUnlock,
			$lockOwner,
			$itemToUseLockOf,
			true,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When the public unlocks file :itemToUnlock using the WebDAV API
	 *
	 * @param string $itemToUnlock
	 *
	 * @return void
	 */
	public function unlockItemAsPublicUsingWebDavAPI(string $itemToUnlock): void {
		$token = ($this->featureContext->isUsingSharingNG())
		? $this->featureContext->shareNgGetLastCreatedLinkShareToken()
		: $this->featureContext->getLastCreatedPublicShareToken();
		$response = $this->unlockItemWithLastLockOfUserAndItemUsingWebDavAPI(
			$token,
			$itemToUnlock,
			$token,
			$itemToUnlock,
			true,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^user "([^"]*)" moves (?:file|folder|entry) "([^"]*)" to "([^"]*)" sending the locktoken of (?:file|folder|entry) "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $fileSource
	 * @param string $fileDestination
	 * @param string $itemToUseLockOf
	 *
	 * @return void
	 */
	public function moveItemSendingLockToken(
		string $user,
		string $fileSource,
		string $fileDestination,
		string $itemToUseLockOf,
	): void {
		$response = $this->moveItemSendingLockTokenOfUser(
			$user,
			$fileSource,
			$fileDestination,
			$itemToUseLockOf,
			$user,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @param string $user
	 * @param string $fileSource
	 * @param string $fileDestination
	 * @param string $itemToUseLockOf
	 * @param string $lockOwner
	 *
	 * @return ResponseInterface
	 */
	public function moveItemSendingLockTokenOfUser(
		string $user,
		string $fileSource,
		string $fileDestination,
		string $itemToUseLockOf,
		string $lockOwner,
	): ResponseInterface {
		$user = $this->featureContext->getActualUsername($user);
		$lockOwner = $this->featureContext->getActualUsername($lockOwner);
		$destination = $this->featureContext->destinationHeaderValue(
			$user,
			$fileDestination,
		);
		$token = $this->tokenOfLastLock[$lockOwner][$itemToUseLockOf];
		$headers = [
			"Destination" => $destination,
			"If" => "(<$token>)",
		];
		return $this->featureContext->makeDavRequest(
			$user,
			"MOVE",
			$fileSource,
			$headers,
		);
	}

	/**
	 * @When /^user "([^"]*)" moves (?:file|folder|entry) "([^"]*)" to "([^"]*)" sending the locktoken of (?:file|folder|entry) "([^"]*)" of user "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $fileSource
	 * @param string $fileDestination
	 * @param string $itemToUseLockOf
	 * @param string $lockOwner
	 *
	 * @return void
	 */
	public function userMovesItemSendingLockTokenOfUser(
		string $user,
		string $fileSource,
		string $fileDestination,
		string $itemToUseLockOf,
		string $lockOwner,
	): void {
		$response = $this->moveItemSendingLockTokenOfUser(
			$user,
			$fileSource,
			$fileDestination,
			$itemToUseLockOf,
			$lockOwner,
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^user "([^"]*)" uploads file with content "([^"]*)" to "([^"]*)" sending the locktoken of (?:file|folder|entry) "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $content
	 * @param string $destination
	 * @param string $itemToUseLockOf
	 *
	 * @return void
	 */
	public function userUploadsAFileWithContentTo(
		string $user,
		string $content,
		string $destination,
		string $itemToUseLockOf,
	): void {
		$user = $this->featureContext->getActualUsername($user);
		$token = $this->tokenOfLastLock[$user][$itemToUseLockOf];
		$this->featureContext->pauseUploadDelete();
		$response = $this->featureContext->makeDavRequest(
			$user,
			"PUT",
			$destination,
			["If" => "(<$token>)"],
			$content,
		);
		$this->featureContext->setResponse($response);
		$this->featureContext->setLastUploadDeleteTime(\time());
	}

	/**
	 * @Then :count locks should be reported for file :file of user :user by the WebDAV API
	 *
	 * @param int $count
	 * @param string $file
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function numberOfLockShouldBeReported(int $count, string $file, string $user): void {
		$lockCount = $this->countLockOfResources($user, $file);
		Assert::assertEquals(
			$count,
			$lockCount,
			"Expected $count lock(s) for '$file' but found '$lockCount'",
		);
	}

	/**
	 * @When the user waits for :time seconds to expire the lock
	 *
	 * @param int $time
	 *
	 * @return void
	 */
	public function waitForCertainSecondsToExpireTheLock(int $time): void {
		\sleep($time);
	}

	/**
	 * @Then :count locks should be reported for file :file inside the space :space of user :user
	 *
	 * @param int $count
	 * @param string $file
	 * @param string $spaceName
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function numberOfLockShouldBeReportedInProjectSpace(
		int $count,
		string $file,
		string $spaceName,
		string $user,
	): void {
		$response = $this->spacesContext->sendPropfindRequestToSpace($user, $spaceName, $file, null, '0');
		$this->featureContext->theHTTPStatusCodeShouldBe(207, "", $response);
		$responseXmlObject = HttpRequestHelper::getResponseXml($response, __METHOD__);
		$xmlPart = $responseXmlObject->xpath("//d:response//d:lockdiscovery/d:activelock");
		if (\is_array($xmlPart)) {
			$lockCount = \count($xmlPart);
		} else {
			throw new Exception("xmlPart for 'd:activelock' was expected to be array but found: $xmlPart");
		}
		Assert::assertEquals(
			$count,
			$lockCount,
			"Expected $count lock(s) for '$file' inside space '$spaceName' but found '$lockCount'",
		);
	}

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
		// Get all the contexts you need in this context
		$this->featureContext = BehatHelper::getContext($scope, $environment, 'FeatureContext');
		$this->publicWebDavContext = BehatHelper::getContext($scope, $environment, 'PublicWebDavContext');
		$this->spacesContext = BehatHelper::getContext($scope, $environment, 'SpacesContext');
	}
}
