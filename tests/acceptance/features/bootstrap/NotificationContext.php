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
use TestHelpers\EmailHelper;
use PHPUnit\Framework\Assert;
use TestHelpers\GraphHelper;
use Behat\Gherkin\Node\TableNode;
use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\ResponseInterface;

require_once 'bootstrap.php';

/**
 * Defines application features from the specific context.
 */
class NotificationContext implements Context {
	private FeatureContext $featureContext;
	private SpacesContext $spacesContext;
	private SettingsContext $settingsContext;
	private string $notificationEndpointPath = '/apps/notifications/api/v1/notifications?format=json';

	private array $notificationIds;

	/**
	 * @return array[]
	 */
	public function getNotificationIds():array {
		return $this->notificationIds;
	}

	/**
	 * @return array[]
	 */
	public function getLastNotificationId():array {
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
		$this->spacesContext = $environment->getContext('SpacesContext');
		$this->settingsContext = $environment->getContext('SettingsContext');
	}

	/**
	 * @When /^user "([^"]*)" lists all notifications$/
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function userListAllNotifications(string $user):void {
		$this->setUserRecipient($user);
		$headers = ["accept-language" => $this->settingsContext->getSettingLanguageValue($user)];
		$response = OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			'GET',
			$this->notificationEndpointPath,
			$this->featureContext->getStepLineRef(),
			[],
			2,
			$headers
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user deletes all notifications
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userDeletesAllNotifications(string $user):void {
		$this->userListAllNotifications($user);
		if (isset($this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data)) {
			$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data;
			foreach ($responseBody as $value) {
				// set notificationId
				$this->notificationIds[] = $value->notification_id;
			}
		}
		$this->featureContext->setResponse($this->userDeletesNotification($user));
	}

	/**
	 * @When user :user deletes a notification related to resource :resource with subject :subject
	 *
	 * @param string $user
	 * @param string $resource
	 * @param string $subject
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userDeletesNotificationOfResourceAndSubject(string $user, string $resource, string $subject):void {
		$this->userListAllNotifications($user);
		$this->filterResponseByNotificationSubjectAndResource($subject, $resource);
		$this->featureContext->setResponse($this->userDeletesNotification($user));
	}

	/**
	 * deletes notification
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userDeletesNotification(string $user):ResponseInterface {
		$this->setUserRecipient($user);
		$payload["ids"] = $this->getNotificationIds();
		return OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			'DELETE',
			$this->notificationEndpointPath,
			$this->featureContext->getStepLineRef(),
			\json_encode($payload),
			2
		);
	}

	/**
	 * @Then the notifications should be empty
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theNotificationsShouldBeEmpty(): void {
		$statusCode = $this->featureContext->getResponse()->getStatusCode();
		if ($statusCode !== 200) {
			$response = $this->featureContext->getResponse()->getBody()->getContents();
			throw new \Exception(
				__METHOD__
				. " Failed to get user notification list" . $response
			);
		}
		$notifications = $this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data;
		Assert::assertNull($notifications, "response should not contain any notification");
	}

	/**
	 * @Then user :user should not have any notification
	 *
	 * @param $user
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldNotHaveAnyNotification($user): void {
		$this->userListAllNotifications($user);
		$this->theNotificationsShouldBeEmpty();
	}

	/**
	 * @Then /^user "([^"]*)" should have "([^"]*)" notifications$/
	 *
	 * @param string $user
	 * @param int $numberOfNotification
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldHaveNotifications(string $user, int $numberOfNotification): void {
		if (!isset($this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data)) {
			throw new Exception("Notification is empty");
		}
		$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data;
		$actualNumber = \count($responseBody);
		Assert::assertEquals(
			$numberOfNotification,
			$actualNumber,
			"Expected number of notification was '$numberOfNotification', but got '$actualNumber'"
		);
	}

	/**
	 * @Then /^the JSON response should contain a notification message with the subject "([^"]*)" and the message-details should match$/
	 *
	 * @param string $subject
	 * @param PyStringNode $schemaString
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theJsonDataFromLastResponseShouldMatch(
		string $subject,
		PyStringNode $schemaString
	): void {
		$responseBody = $this->filterResponseAccordingToNotificationSubject($subject);
		// substitute the value here
		$schemaString = $schemaString->getRaw();
		$schemaString = $this->featureContext->substituteInLineCodes(
			$schemaString,
			$this->featureContext->getCurrentUser(),
			[],
			[],
			null,
			$this->getUserRecipient(),
		);
		$this->featureContext->assertJsonDocumentMatchesSchema(
			$responseBody,
			$this->featureContext->getJSONSchema($schemaString)
		);
	}

	/**
	 * filter notification according to subject
	 *
	 * @param string $subject
	 *
	 * @return object
	 */
	public function filterResponseAccordingToNotificationSubject(string $subject): object {
		$responseBody =  null;
		if (isset($this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data)) {
			$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data;
			foreach ($responseBody as $value) {
				if (isset($value->subject) && $value->subject === $subject) {
					$responseBody = $value;
					// set notificationId
					$this->notificationIds[] = $value->notification_id;
					break;
				}
			}
		} else {
			$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent();
		}
		return $responseBody;
	}

	/**
	 * filter notification according to subject and resource
	 *
	 * @param string $subject
	 * @param string $resource
	 *
	 * @return array
	 */
	public function filterResponseByNotificationSubjectAndResource(string $subject, string $resource): array {
		$responseBodyArray = [];
		$statusCode = $this->featureContext->getResponse()->getStatusCode();
		if ($statusCode !== 200) {
			$response = $this->featureContext->getResponse()->getBody()->getContents();
			Assert::fail($response . " Response should contain status code 200");
		}
		if (isset($this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data)) {
			$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data;
			foreach ($responseBody as $value) {
				if (isset($value->subject) && $value->subject === $subject && isset($value->messageRichParameters->resource->name) && $value->messageRichParameters->resource->name === $resource) {
					$this->notificationIds[] = $value->notification_id;
					$responseBodyArray[] = $value;
				}
			}
		} else {
			$responseBodyArray[] = $this->featureContext->getJsonDecodedResponseBodyContent();
			Assert::fail("Response should contain notification but found: $responseBodyArray");
		}
		return $responseBodyArray;
	}

	/**
	 * @Then /^user "([^"]*)" should (?:get|have) a notification with subject "([^"]*)" and message:$/
	 *
	 * @param string $user
	 * @param string $subject
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function userShouldGetANotificationWithMessage(string $user, string $subject, TableNode $table):void {
		$this->userListAllNotifications($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(200);
		// sometimes the test might try to get notification before the notification is created by the server
		// in order to prevent test from failing because of that list the notifications again
		if (!isset($this->filterResponseAccordingToNotificationSubject($subject)->message)) {
			\sleep(1);
			$this->featureContext->setResponse(null);
			$this->userListAllNotifications($user);
			$this->featureContext->theHTTPStatusCodeShouldBe(200);
		}
		$actualMessage = str_replace(["\r", "\n"], " ", $this->filterResponseAccordingToNotificationSubject($subject)->message);
		$expectedMessage = $table->getColumnsHash()[0]['message'];
		Assert::assertSame(
			$expectedMessage,
			$actualMessage,
			__METHOD__ . "expected message to be '$expectedMessage' but found'$actualMessage'"
		);
	}

	/**
	 * @Then user :user should not have a notification related to resource :resource with subject :subject
	 *
	 * @param string $user
	 * @param string $resource
	 * @param string $subject
	 *
	 * @return void
	 */
	public function userShouldNotHaveANotificationRelatedToResourceWithSubject(string $user, string $resource, string $subject):void {
		$this->userListAllNotifications($user);
		$response = $this->filterResponseByNotificationSubjectAndResource($subject, $resource);
		Assert::assertCount(0, $response, "Response should not contain notification related to resource '$resource' with subject '$subject' but found" . print_r($response, true));
	}

	/**
	 * @Then user :user should have received the following email from user :sender about the share of project space :spaceName
	 *
	 * @param string $user
	 * @param string $sender
	 * @param string $spaceName
	 * @param PyStringNode $content
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldHaveReceivedTheFollowingEmailFromUserAboutTheShareOfProjectSpace(string $user, string $sender, string $spaceName, PyStringNode $content):void {
		$rawExpectedEmailBodyContent = \str_replace("\r\n", "\n", $content->getRaw());
		$this->featureContext->setResponse(
			GraphHelper::getMySpaces(
				$this->featureContext->getBaseUrl(),
				$user,
				$this->featureContext->getPasswordForUser($user)
			)
		);
		$expectedEmailBodyContent = $this->featureContext->substituteInLineCodes(
			$rawExpectedEmailBodyContent,
			$sender,
			[],
			[
				[
					"code" => "%space_id%",
					"function" =>
						[$this->spacesContext, "getSpaceIdByName"],
					"parameter" => [$sender, $spaceName]
				],
			],
			null,
			null
		);
		$this->assertEmailContains($user, $expectedEmailBodyContent);
	}

	/**
	 * @Then user :user should have received the following email from user :sender
	 *
	 * @param string $user
	 * @param string $sender
	 * @param PyStringNode $content
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldHaveReceivedTheFollowingEmailFromUser(string $user, string $sender, PyStringNode $content):void {
		$rawExpectedEmailBodyContent = \str_replace("\r\n", "\n", $content->getRaw());
		$expectedEmailBodyContent = $this->featureContext->substituteInLineCodes(
			$rawExpectedEmailBodyContent,
			$sender
		);
		$this->assertEmailContains($user, $expectedEmailBodyContent);
	}

	/***
	 * @param string $user
	 * @param string $expectedEmailBodyContent
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function assertEmailContains(string $user, string $expectedEmailBodyContent):void {
		$address = $this->featureContext->getEmailAddressForUser($user);
		$this->featureContext->pushEmailRecipientAsMailBox($address);
		$actualEmailBodyContent = EmailHelper::getBodyOfLastEmail($address, $this->featureContext->getStepLineRef());
		Assert::assertStringContainsString(
			$expectedEmailBodyContent,
			$actualEmailBodyContent,
			"The email address '$address' should have received an email with the body containing $expectedEmailBodyContent
			but the received email is $actualEmailBodyContent"
		);
	}

	/**
	 * Delete all the inbucket emails
	 *
	 * @AfterScenario @email
	 *
	 * @return void
	 */
	public function clearInbucketMessages():void {
		try {
			if (!empty($this->featureContext->emailRecipients)) {
				foreach ($this->featureContext->emailRecipients as $emailRecipent) {
					EmailHelper::deleteAllEmailsForAMailbox(
						EmailHelper::getLocalEmailUrl(),
						$this->featureContext->getStepLineRef(),
						$emailRecipent
					);
				}
			}
		} catch (Exception $e) {
			echo __METHOD__ .
				" could not delete inbucket messages, is inbucket set up?\n" .
				$e->getMessage();
		}
	}
}
