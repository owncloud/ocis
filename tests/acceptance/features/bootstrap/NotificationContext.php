<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Viktor Scharf <vscharf@owncloud.com>
 * @copyright Copyright (c) 2023 Viktor Scharf vscharf@owncloud.com
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use TestHelpers\OcsApiHelper;
use Behat\Gherkin\Node\PyStringNode;
use Helmich\JsonAssert\JsonAssertions;

require_once 'bootstrap.php';

/**
 * Defines application features from the specific context.
 */
class NotificationContext implements Context {
	/**
	 * @var FeatureContext
	 */
	private $featureContext;

	/**
	 * @var string
	 */
	private string $notificationEndpointPath = '/apps/notifications/api/v1/notifications?format=json';

	/**
	 * @var array[]
	 */
	private $notificationIds;

	/**
	 * @return array[]
	 */
	public function getNotificationIds():array {
		return $this->notificationIds;
	}

	/**
	 * @return array[]
	 */
	public function getLastNotificationIds():array {
		return \end($this->notificationIds);
	}

	/**
	 * @var string
	 */
	private string $userRecipient;

	/**
	 * @param string $userRecipient
	 *
	 * @return void
	 */
	public function setUserRecipient(string $userRecipient): void {
		$this->userRecipient = $userRecipient;
	}

	/**
	 * @return string
	 */
	public function getUserRecipient(): string {
		return $this->userRecipient;
	}

	/**
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 * @throws Exception
	 */
	public function setUpScenario(BeforeScenarioScope $scope):void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context
		$this->featureContext = $environment->getContext('FeatureContext');
	}

	/**
	 * @Then /^user "([^"]*)" list all notifications$/
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function userListAllNotifications(string $user):void {
		$this->setUserRecipient($user);
		$response = OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			'GET',
			$this->notificationEndpointPath,
			$this->featureContext->getStepLineRef()
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Then /^the JSON response should contain message with type "([^"]*)" and match$/
	 *
	 * @param string $messageType
	 * @param string|null $spaceName
	 * @param PyStringNode $schemaString
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theJsonDataFromLastResponseShouldMatch(
		string $messageType,
		?string $spaceName = null,
		PyStringNode $schemaString
	): void {
		if (isset($this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data)) {
			$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data;
			foreach ($responseBody as $value) {
				if (isset($value->subject) && $value->subject === $messageType) {
					$responseBody = $value;
					// set notificationId
					$this->notificationIds[] = $value->notification_id;
					break;
				}
			}
		} else {
			$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent();
		}

		// substitute the value here
		$schemaString = $schemaString->getRaw();
		$schemaString = $this->featureContext->substituteInLineCodes(
			$schemaString,
			$this->featureContext->getCurrentUser(),
			[],
			[
				[
					"code" => "%space_id%",
					"function" =>
						[$this, "getSpaceIdByNameFromResponse"],
					"parameter" => [$spaceName]
				]
			],
			null,
			$this->getUserRecipient(),
		);
		JsonAssertions::assertJsonDocumentMatchesSchema(
			$responseBody,
			$this->featureContext->getJSONSchema($schemaString)
		);
	}
}
