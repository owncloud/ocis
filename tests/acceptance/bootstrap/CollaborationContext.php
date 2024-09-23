<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Amrita Shrestha <amrita@jankaritech.com>
 * @copyright Copyright (c) 2024 Amrita Shrestha <amrita@jankaritech.com>
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
use TestHelpers\HttpRequestHelper;
use TestHelpers\WebDavHelper;
use TestHelpers\CollaborationHelper;

/**
 * steps needed to re-configure oCIS server
 */
class CollaborationContext implements Context {
	private FeatureContext $featureContext;
	private SpacesContext $spacesContext;
	private string $lastAppOpenData;

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
	 * @param string $data
	 *
	 * @return void
	 */
	public function setLastAppOpenData(string $data): void {
		$this->lastAppOpenData = $data;
	}

	/**
	 * @return string
	 */
	public function getLastAppOpenData(): string {
		return $this->lastAppOpenData;
	}

	/**
	 * @When user :user checks the information of file :file of space :space using office :app
	 * @When user :user checks the information of file :file of space :space using office :app with view mode :view
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $space
	 * @param string $app
	 * @param string|null $viewMode
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 */
	public function userChecksTheInformationOfFileOfSpaceUsingOffice(string $user, string $file, string $space, string $app, string $viewMode = null): void {
		$fileId = $this->spacesContext->getFileId($user, $space, $file);
		$response = \json_decode(
			CollaborationHelper::sendPOSTRequestToAppOpen(
				$fileId,
				$app,
				$this->featureContext->getActualUsername($user),
				$this->featureContext->getPasswordForUser($user),
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$viewMode
			)->getBody()->getContents()
		);

		$accessToken = $response->form_parameters->access_token;

		// Extract the WOPISrc from the app_url
		$parsedUrl = parse_url($response->app_url);
		parse_str($parsedUrl['query'], $queryParams);
		$wopiSrc = $queryParams['WOPISrc'];

		$this->featureContext->setResponse(
			HttpRequestHelper::get(
				$wopiSrc . "?access_token=$accessToken",
				$this->featureContext->getStepLineRef()
			)
		);
	}

	/**
	 * @When user :user creates a file :file inside folder :folder in space :space using wopi endpoint
	 * @When user :user tries to create a file :file inside folder :folder in space :space using wopi endpoint
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $folder
	 * @param string $space
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userCreatesFileInsideFolderInSpaceUsingWopiEndpoint(string $user, string $file, string $folder, string $space): void {
		$parentContainerId = $this->spacesContext->getResourceId($user, $space, $folder);
		$this->featureContext->setResponse(
			CollaborationHelper::createFile(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$user,
				$this->featureContext->getPasswordForUser($user),
				$parentContainerId,
				$file
			)
		);
	}

	/**
	 * @param string $file
	 * @param string $password
	 * @param string $folder
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function createFile(string $file, string $password, string $folder = ""): void {
		$token = $this->featureContext->shareNgGetLastCreatedLinkShareToken();
		$davPath = WebDavHelper::getDavPath($token, null, "public-files-new") . "/$folder";
		$response = HttpRequestHelper::sendRequest(
			$this->featureContext->getBaseUrl() . "/$davPath",
			$this->featureContext->getStepLineRef(),
			"PROPFIND",
			"public",
			$this->featureContext->getActualPassword($password)
		);
		$responseXml = HttpRequestHelper::getResponseXml(
			$response,
			__METHOD__
		);
		$xmlPart = $responseXml->xpath("//d:prop/oc:fileid");
		$parentContainerId = (string) $xmlPart[0];

		$headers = [
			"Public-Token" => $token
		];
		$this->featureContext->setResponse(
			CollaborationHelper::createFile(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				"public",
				$this->featureContext->getActualPassword($password),
				$parentContainerId,
				$file,
				$headers
			)
		);
	}

	/**
	 * @When the public creates a file :file inside the last shared public link folder with password :password using wopi endpoint
	 * @When the public tries to create a file :file inside the last shared public link folder with password :password using wopi endpoint
	 *
	 * @param string $file
	 * @param string $password
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function thePublicCreatesAFileInsideTheLastSharedPublicLinkFolderWithPasswordUsingWopiEndpoint(string $file, string $password): void {
		$this->createFile($file, $password);
	}

	/**
	 * @When the public creates a file :file inside folder :folder in the last shared public link space with password :password using wopi endpoint
	 * @When the public tries to create a file :file inside folder :folder in the last shared public link space with password :password using wopi endpoint
	 *
	 * @param string $file
	 * @param string $folder
	 * @param string $password
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function thePublicCreatesAFileInsideFolderInTheLastSharedPublicLinkSpaceWithPasswordUsingWopiEndpoint(string $file, string $folder, string $password): void {
		$this->createFile($file, $password, $folder);
	}

	/**
	 * @When user :user tries to check the information of file :file of space :space using office :app with invalid file-id
	 *
	 * @param string $user
	 * @param string $file
	 * @param string $space
	 * @param string $app
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userTriesToCheckTheInformationOfFileOfSpaceUsingOfficeWithInvalidFileId(string $user, string $file, string $space, string $app): void {
		$response = \json_decode(
			CollaborationHelper::sendPOSTRequestToAppOpen(
				$this->spacesContext->getFileId($user, $space, $file),
				$app,
				$this->featureContext->getActualUsername($user),
				$this->featureContext->getPasswordForUser($user),
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef()
			)->getBody()->getContents()
		);
		$accessToken = $response->form_parameters->access_token;

		// Extract the WOPISrc from the app_url
		$parsedUrl = parse_url($response->app_url);
		parse_str($parsedUrl['query'], $queryParams);
		$wopiSrc = $queryParams['WOPISrc'];
		$position = strpos($wopiSrc, '/files/') + \strlen('/files/');

		// Extract the base URL up to and including '/files/'
		$fullUrl = substr($wopiSrc, 0, $position) . WebDavHelper::generateUUIDv4();
		$this->featureContext->setResponse(
			HttpRequestHelper::get(
				$fullUrl . "?access_token=$accessToken",
				$this->featureContext->getStepLineRef()
			)
		);
	}

	/**
	 * @When user :user tries to create a file :file inside deleted folder using wopi endpoint
	 *
	 * @param string $user
	 * @param string $file
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userTriesToCreateAFileInsideDeletedFolderUsingWopiEndpoint(string $user, string $file): void {
		$parentContainerId = $this->featureContext->getStoredFileID();
		$this->featureContext->setResponse(
			CollaborationHelper::createFile(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$user,
				$this->featureContext->getPasswordForUser($user),
				$parentContainerId,
				$file
			)
		);
	}

	/**
	 * @Given user :user has sent the following app-open request:
	 *
	 * @param string $user
	 * @param TableNode $properties
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userHasSentTheFollowingAppOpenRequest(string $user, TableNode $properties): void {
		$rows = $properties->getRowsHash();
		$appResponse = CollaborationHelper::sendPOSTRequestToAppOpen(
			$this->spacesContext->getFileId($user, $rows['space'], $rows['resource']),
			$rows['app'],
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef()
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $appResponse);
		$this->setLastAppOpenData($appResponse->getBody()->getContents());
	}

	/**
	 * @When user :user tries to get the information of the last opened file using wopi endpoint
	 * @When user :user gets the information of the last opened file using wopi endpoint
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userGetsTheInformationOfTheLastOpenedFileUsingWopiEndpoint(string $user): void {
		$response = json_decode($this->getLastAppOpenData());
		$accessToken = $response->form_parameters->access_token;

		// Extract the WOPISrc from the app_url
		$parsedUrl = parse_url($response->app_url);
		parse_str($parsedUrl['query'], $queryParams);
		$wopiSrc = $queryParams['WOPISrc'];

		$this->featureContext->setResponse(
			HttpRequestHelper::get(
				$wopiSrc . "?access_token=$accessToken",
				$this->featureContext->getStepLineRef()
			)
		);
	}
}
