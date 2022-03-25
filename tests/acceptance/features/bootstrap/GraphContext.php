<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Kiran Parajuli <kiran@jankaritech.com>
 * @copyright Copyright (c) 2021 Kiran Parajuli kiran@jankaritech.com
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Behat\Gherkin\Node\TableNode;
use TestHelpers\GraphHelper;
use TestHelpers\HttpRequestHelper;
use GuzzleHttp\Exception\GuzzleException;
use PHPUnit\Framework\Assert;

require_once "bootstrap.php";

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
	 * @When /^the administrator sends a user creation request for user "([^"]*)" password "([^"]*)" using the graph API$/
	 *
	 * @param string $user
	 * @param string $password
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function adminSendsUserCreationRequestUsingTheGraphApi(string $user, string $password):void {
		$user = $this->featureContext->getActualUsername($user);
		$password = $this->featureContext->getActualPassword($password);
		$response = GraphHelper::createUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$user,
			$password
		);
		$this->featureContext->setResponse($response);
		$this->featureContext->pushToLastStatusCodesArrays();
		$success = $this->featureContext->theHTTPStatusCodeWasSuccess();
		$this->featureContext->addUserToCreatedUsersList(
			$user,
			$password,
			null,
			null,
			$success
		);
	}

	/**
	 * @When /^the administrator sends a user creation request for the following users with password using the graph API$/
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theAdministratorSendsAUserCreationRequestForTheFollowingUsersWithPasswordUsingTheGraphAPI(TableNode $table) {
		$this->featureContext->verifyTableNodeColumns($table, ["username", "password"]);
		$users = $table->getHash();
		foreach ($users as $user) {
			$this->adminSendsUserCreationRequestUsingTheGraphApi($user["username"], $user["password"]);
		}
	}

	/**
	 * @Then /^the graph API response should return the following error$/
	 *
	 * @param TableNode $body
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theGraphApiResponseShouldReturnTheFollowingError(TableNode $body):void {
		$this->featureContext->verifyTableNodeRows($body, ['code', 'message']);
		$bodyRows = $body->getRowsHash();
		$responseData = $this->featureContext->getJsonDecodedResponse();
		// parse "{space}" to " " from the message
		$bodyRows['message'] = \str_replace('{space}', ' ', $bodyRows['message']);
		Assert::assertEquals(
			$bodyRows['code'],
			$responseData['error']['code'],
			"Status code is not as expected"
		);
		Assert::assertEquals(
			$bodyRows['message'],
			$responseData['error']['message'],
			"Status message is not as expected"
		);
	}

	/**
	 * @When /^the administrator sends a group creation request for group "([^"]*)" using the graph API$/
	 *
	 * @param string $group
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function adminSendsGroupCreationRequestUsingTheGraphAPI(string $group):void {
		$response = GraphHelper::createGroup(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$group
		);
		$this->featureContext->setResponse($response);
		$justCreatedGroup = $this->featureContext->getJsonDecodedResponse();
		$this->featureContext->pushToLastStatusCodesArrays();
		$this->featureContext->addGroupToCreatedGroupsList($group, true, true, $justCreatedGroup['id']);
	}

	/**
	 * @When /^the administrator sends a group creation request for the following groups using the graph API$/
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theAdministratorSendsAGroupCreationRequestForTheFollowingGroupsUsingTheGraphAPI(TableNode $table) {
		$this->featureContext->verifyTableNodeColumns($table, ["group_display_name"], ['comment']);
		$groups = $table->getHash();
		foreach ($groups as $group) {
			$this->adminSendsGroupCreationRequestUsingTheGraphAPI($group["group_display_name"]);
		}
	}
}
