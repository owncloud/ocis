<?php

declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Viktor Scharf <v.scharf@owncloud.com>
 * @copyright Copyright (c) 2022 Viktor Scharf v.scharf@owncloud.com
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use GuzzleHttp\Exception\GuzzleException;
use TestHelpers\HttpRequestHelper;
use Behat\Gherkin\Node\TableNode;
use TestHelpers\WebDavHelper;

require_once 'bootstrap.php';

/**
 * Context for the TUS-specific steps using the Graph API
 */
class SpacesTUSContext implements Context {

	/**
	 * @var FeatureContext
	 */
	private FeatureContext $featureContext;

	/**
	 * @var TUSContext
	 */
	private TUSContext $tusContext;

	/**
	 * @var SpacesContext
	 */
	private SpacesContext $spacesContext;

	/**
	 * @var string
	 */
	private string $baseUrl;

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
		$this->tusContext = $environment->getContext('TUSContext');
		$this->baseUrl = \trim($this->featureContext->getBaseUrl(), "/");
	}

	/**
	 * @Given /^user "([^"]*)" has uploaded a file from "([^"]*)" to "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $source
	 * @param string $destination
	 * @param string $spaceName
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userHasUploadedFileViaTusInSpace(string $user, string $source, string $destination, string $spaceName): void {
		$this->userUploadsAFileViaTusInsideOfTheSpaceUsingTheWebdavApi($user, $source, $destination, $spaceName);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "Expected response status code should be 200");
	}

	/**
	 * @When /^user "([^"]*)" uploads a file from "([^"]*)" to "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $source
	 * @param string $spaceName
	 * @param string $destination
	 * @param array|null $uploadMetadata
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userUploadsAFileViaTusInsideOfTheSpaceUsingTheWebdavApi(
		string $user,
		string $source,
		string $destination,
		string $spaceName,
		?array $uploadMetadata = null
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->userUploadsUsingTusAFileTo($user, $source, $destination);
	}

	/**
	 * @Given user :user has created a new TUS resource for the space :spaceName using the WebDAV API with these headers:
	 *
	 * @param string $user
	 * @param string $spaceName
	 * @param TableNode $headers
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userHasCreatedANewTusResourceForTheSpaceUsingTheWebdavApiWithTheseHeaders(
		string $user,
		string $spaceName,
		TableNode $headers
	): void {
		$this->userCreatesANewTusResourceForTheSpaceUsingTheWebdavApiWithTheseHeaders($user, $spaceName, $headers);
		$this->featureContext->theHTTPStatusCodeShouldBe(201, "Expected response status code should be 201");
	}

	/**
	 * @When user :user creates a new TUS resource for the space :spaceName using the WebDAV API with these headers:
	 *
	 * @param string $user
	 * @param string $spaceName
	 * @param TableNode $headers
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userCreatesANewTusResourceForTheSpaceUsingTheWebdavApiWithTheseHeaders(
		string $user,
		string $spaceName,
		TableNode $headers
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->createNewTUSResourceWithHeaders($user, $headers, '');
	}

	/**
	 * @When /^user "([^"]*)" uploads a file with content "([^"]*)" to "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $content
	 * @param string $resource
	 * @param string $spaceName
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function userUploadsAFileWithContentToViaTusInsideOfTheSpaceUsingTheWebdavApi(
		string $user,
		string $content,
		string $resource,
		string $spaceName
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->userUploadsAFileWithContentToUsingTus($user, $content, $resource);
	}
}
