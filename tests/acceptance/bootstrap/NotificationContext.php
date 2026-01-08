<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Viktor Scharf <vscharf@owncloud.com>
 * @copyright Copyright (c) 2023 Viktor Scharf vscharf@owncloud.com
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Behat\Gherkin\Node\TableNode;
use Behat\Gherkin\Node\PyStringNode;
use PHPUnit\Framework\Assert;
use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\OcsApiHelper;
use TestHelpers\SettingsHelper;
use TestHelpers\BehatHelper;

require_once 'bootstrap.php';

/**
 * Defines application features from the specific context.
 */
class NotificationContext implements Context {
	private FeatureContext $featureContext;
	private string $notificationEndpointPath = '/apps/notifications/api/v1/notifications';
	private string $globalNotificationEndpointPath = '/apps/notifications/api/v1/notifications/global';

	private array $notificationIds;

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
	 * @return array[]
	 */
	public function getNotificationIds(): array {
		return $this->notificationIds;
	}

	/**
	 * @return array[]
	 */
	public function getLastNotificationId(): array {
		return \end($this->notificationIds);
	}

	/**
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 * @throws Exception
	 */
	public function before(BeforeScenarioScope $scope): void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context
		$this->featureContext = BehatHelper::getContext($scope, $environment, 'FeatureContext');
	}

	/**
	 * delete all in-app notifications
	 *
	 * @AfterScenario @notification
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function deleteDeprovisioningNotification(): void {
		$payload["ids"] = ["deprovision"];

		OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			'DELETE',
			$this->globalNotificationEndpointPath,
			json_encode($payload),
		);
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 */
	public function listAllNotifications(string $user): ResponseInterface {
		$this->setUserRecipient($user);
		$language = SettingsHelper::getLanguageSettingValue(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
		);
		$headers = ["accept-language" => $language];
		return OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			'GET',
			$this->notificationEndpointPath . '?format=json',
			[],
			2,
			$headers,
		);
	}

	/**
	 * @When /^user "([^"]*)" lists all notifications$/
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function userListAllNotifications(string $user): void {
		$response = $this->listAllNotifications($user);
		$this->featureContext->setResponse($response);
		$this->featureContext->pushToLastHttpStatusCodesArray();
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function deleteAllInAppNotifications(string $user): ResponseInterface {
		$response = $this->listAllNotifications($user);
		if (isset($this->featureContext->getJsonDecodedResponseBodyContent($response)->ocs->data)) {
			$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent($response)->ocs->data;
			foreach ($responseBody as $value) {
				// set notificationId
				$this->notificationIds[] = $value->notification_id;
			}
		}
		return $this->userDeletesNotification($user);
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
	public function userDeletesAllNotifications(string $user): void {
		$response = $this->deleteAllInAppNotifications($user);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given user :user has deleted all notifications
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userHasDeletedAllNotifications(string $user): void {
		$response = $this->deleteAllInAppNotifications($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "", $response);
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
	public function userDeletesNotificationOfResourceAndSubject(string $user, string $resource, string $subject): void {
		$response = $this->listAllNotifications($user);
		$this->filterNotificationsBySubjectAndResource($subject, $resource, $response);
		$this->featureContext->setResponse($this->userDeletesNotification($user));
	}

	/**
	 * deletes notification
	 *
	 * @param string $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userDeletesNotification(string $user): ResponseInterface {
		$this->setUserRecipient($user);
		$payload["ids"] = $this->getNotificationIds();
		return OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			'DELETE',
			$this->notificationEndpointPath . '?format=json',
			\json_encode($payload),
			2,
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
				. " Failed to get user notification list" . $response,
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
		$response = $this->listAllNotifications($user);
		$notifications = $this->featureContext->getJsonDecodedResponseBodyContent($response)->ocs->data;
		Assert::assertNull($notifications, "response should not contain any notification");
	}

	/**
	 * @Then /^there should be "([^"]*)" notifications$/
	 *
	 * @param int $numberOfNotification
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldHaveNotifications(int $numberOfNotification): void {
		if (!isset($this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data)) {
			throw new Exception("Notification is empty");
		}
		$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent()->ocs->data;
		$actualNumber = \count($responseBody);
		Assert::assertEquals(
			$numberOfNotification,
			$actualNumber,
			"Expected number of notifications was '$numberOfNotification', but got '$actualNumber'",
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
		PyStringNode $schemaString,
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
			$this->featureContext->getJSONSchema($schemaString),
		);
	}

	/**
	 * filter notification according to subject
	 *
	 * @param string $subject
	 * @param ResponseInterface|null $response
	 *
	 * @return object
	 */
	public function filterResponseAccordingToNotificationSubject(
		string $subject,
		?ResponseInterface $response = null,
	): object {
		$response = $response ?? $this->featureContext->getResponse();
		if (isset($this->featureContext->getJsonDecodedResponseBodyContent($response)->ocs->data)) {
			$responseBody = $this->featureContext->getJsonDecodedResponseBodyContent($response)->ocs->data;
			foreach ($responseBody as $value) {
				if (isset($value->subject) && $value->subject === $subject) {
					// set notificationId
					$this->notificationIds[] = $value->notification_id;
					return $value;
				}
			}
		}
		return new StdClass();
	}

	/**
	 * filter notification according to subject and resource
	 *
	 * @param string $subject
	 * @param string $resource
	 * @param ResponseInterface|null $response
	 *
	 * @return array
	 */
	public function filterNotificationsBySubjectAndResource(
		string $subject,
		string $resource,
		?ResponseInterface $response = null,
	): array {
		$filteredNotifications = [];
		$response = $response ?? $this->featureContext->getResponse();
		$responseObject = $this->featureContext->getJsonDecodedResponseBodyContent($response);

		if (!isset($responseObject->ocs->data)) {
			Assert::fail("Response doesn't contain notification: " . print_r($responseObject, true));
		}

		$notifications = $responseObject->ocs->data;
		foreach ($notifications as $notification) {
			if (isset($notification->subject) && $notification->subject === $subject
				&& isset($notification->messageRichParameters->resource->name)
				&& $notification->messageRichParameters->resource->name === $resource
			) {
				$this->notificationIds[] = $notification->notification_id;
				$filteredNotifications[] = $notification;
			}
		}
		return $filteredNotifications;
	}

	/**
	 * filter notification according to subject and space
	 *
	 * @param string $subject
	 * @param string $space
	 * @param ResponseInterface|null $response
	 *
	 * @return array
	 */
	public function filterNotificationsBySubjectAndSpace(
		string $subject,
		string $space,
		?ResponseInterface $response = null,
	): array {
		$filteredNotifications = [];
		$response = $response ?? $this->featureContext->getResponse();
		$responseObject = $this->featureContext->getJsonDecodedResponseBodyContent($response);
		if (!isset($responseObject->ocs->data)) {
			Assert::fail("Response doesn't contain notification: " . print_r($responseObject, true));
		}

		$notifications = $responseObject->ocs->data;
		foreach ($notifications as $notification) {
			if (isset($notification->subject) && $notification->subject === $subject
				&& isset($notification->messageRichParameters->space->name)
				&& $notification->messageRichParameters->space->name === $space
			) {
				$this->notificationIds[] = $notification->notification_id;
				$filteredNotifications[] = $notification;
			}
		}
		return $filteredNotifications;
	}

	/**
	 * @Then /^user "([^"]*)" should (?:get|have) a notification with subject "([^"]*)" and message:$/
	 *
	 * @param string $user
	 * @param string $subject
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldGetANotificationWithMessage(string $user, string $subject, TableNode $table): void {
		$count = 0;
		// Sometimes the test might try to get the notifications before the server has created the notification.
		// To prevent the test from failing because of that, try to list the notifications again
		do {
			if ($count > 0) {
				\sleep(1);
			}
			$this->featureContext->setResponse(null);
			$response = $this->listAllNotifications($user);
			$this->featureContext->theHTTPStatusCodeShouldBe(200, "", $response);
			++$count;
		} while (!isset($this->filterResponseAccordingToNotificationSubject($subject, $response)->message)
			&& $count <= 5
		);
		if (isset($this->filterResponseAccordingToNotificationSubject($subject, $response)->message)) {
			$actualMessage = str_replace(
				["\r", "\n"],
				" ",
				$this->filterResponseAccordingToNotificationSubject($subject, $response)->message,
			);
		} else {
			throw new \Exception("Notification was not found even after retrying for 5 seconds.");
		}
		$expectedMessage = $table->getColumnsHash()[0]['message'];
		Assert::assertStringContainsString(
			$expectedMessage,
			$actualMessage,
			__METHOD__ . "expected message to be '$expectedMessage' but found'$actualMessage'",
		);
	}

	/**
	 * @Then user :user should get a notification for resource :resource with subject :subject and message:
	 *
	 * @param string $user
	 * @param string $resource
	 * @param string $subject
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldGetNotificationForResourceWithMessage(
		string $user,
		string $resource,
		string $subject,
		TableNode $table,
	): void {
		$response = $this->listAllNotifications($user);
		$notification = $this->filterNotificationsBySubjectAndResource($subject, $resource, $response);

		if (\count($notification) === 1) {
			$actualMessage = str_replace(["\r", "\r"], " ", $notification[0]->message);
			$expectedMessage = $table->getColumnsHash()[0]['message'];
			Assert::assertStringContainsString(
				$expectedMessage,
				$actualMessage,
				__METHOD__ . "expected message to be '$expectedMessage' but found'$actualMessage'",
			);
			$response = $this->userDeletesNotification($user);
			$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
		} elseif (\count($notification) === 0) {
			throw new \Exception(
				"Response doesn't contain any notification with resource '$resource' and subject '$subject'.\n"
				. print_r($notification, true),
			);
		} else {
			throw new \Exception(
				"Response contains more than one notification with resource '$resource' and subject '$subject'.\n"
				. print_r($notification, true),
			);
		}
	}

	/**
	 * @Then /^user "([^"]*)" should not have a notification related to (resource|space) "([^"]*)" with subject "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $resourceOrSpace
	 * @param string $resource
	 * @param string $subject
	 *
	 * @return void
	 */
	public function userShouldNotHaveANotificationRelatedToResourceOrSpaceWithSubject(
		string $user,
		string $resourceOrSpace,
		string $resource,
		string $subject,
	): void {
		$response = $this->listAllNotifications($user);
		if ($resourceOrSpace === "space") {
			$filteredResponse = $this->filterNotificationsBySubjectAndSpace($subject, $resource, $response);
		} else {
			$filteredResponse = $this->filterNotificationsBySubjectAndResource($subject, $resource, $response);
		}
		Assert::assertCount(
			0,
			$filteredResponse,
			"Response should not contain notification related to resource '$resource' with subject '$subject' but found"
			. print_r($filteredResponse, true),
		);
	}

	/**
	 *
	 * @param string|null $user
	 * @param string|null $deprovision_date
	 * @param string|null $deprovision_date_format
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 *
	 * @throws JsonException
	 */
	public function userCreatesDeprovisioningNotification(
		?string $user = null,
		?string $deprovision_date = "2043-07-04T11:23:12Z",
		?string $deprovision_date_format = "2006-01-02T15:04:05Z07:00",
	): ResponseInterface {
		$payload["type"] = "deprovision";
		$payload["data"] = [
			"deprovision_date" => $deprovision_date, "deprovision_date_format" => $deprovision_date_format];
		return OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$user ? $this->featureContext->getActualUsername($user) : $this->featureContext->getAdminUsername(),
			$user ? $this->featureContext->getPasswordForUser($user) : $this->featureContext->getAdminPassword(),
			'POST',
			$this->globalNotificationEndpointPath,
			json_encode($payload),
		);
	}

	/**
	 * @When the administrator creates a deprovisioning notification
	 * @When user :user tries to create a deprovisioning notification
	 *
	 * @param string|null $user
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function theAdministratorCreatesADeprovisioningNotification(?string $user = null): void {
		$response = $this->userCreatesDeprovisioningNotification($user);
		$this->featureContext->setResponse($response);
		$this->featureContext->pushToLastHttpStatusCodesArray();
	}

	/**
	 * @When the administrator creates a deprovisioning notification for date :deprovision_date of format :deprovision_date_format
	 *
	 * @param $deprovision_date
	 * @param $deprovision_date_format
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function theAdministratorCreatesADeprovisioningNotificationUsingDateFormat(
		$deprovision_date,
		$deprovision_date_format,
	): void {
		$response = $this->userCreatesDeprovisioningNotification(null, $deprovision_date, $deprovision_date_format);
		$this->featureContext->setResponse($response);
		$this->featureContext->pushToLastHttpStatusCodesArray();
	}

	/**
	 * @Given the administrator has created a deprovisioning notification
	 *
	 * @return void
	 */
	public function userHasCreatedDeprovisioningNotification(): void {
		$response = $this->userCreatesDeprovisioningNotification();
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "", $response);
	}

	/**
	 * @When the administrator deletes the deprovisioning notification
	 * @When user :user tries to delete the deprovisioning notification
	 *
	 * @param string|null $user
	 *
	 * @return void
	 */
	public function userDeletesDeprovisioningNotification(?string $user = null): void {
		$payload["ids"] = ["deprovision"];

		$response = OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$user ? $this->featureContext->getActualUsername($user) : $this->featureContext->getAdminUsername(),
			$user ? $this->featureContext->getPasswordForUser($user) : $this->featureContext->getAdminPassword(),
			'DELETE',
			$this->globalNotificationEndpointPath,
			json_encode($payload),
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * deletes notification using id
	 *
	 * @param string $user
	 * @param string $notificationId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function deleteNotificationUsingId(string $user, string $notificationId): ResponseInterface {
		$deleteNotificationEndpoint = $this->notificationEndpointPath . '/' . $notificationId;
		return OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			'DELETE',
			$deleteNotificationEndpoint,
		);
	}

	/**
	 * @When user :user deletes a notification related to resource :resource with subject :subject using id
	 *
	 * @param string $user
	 * @param string $resource
	 * @param string $subject
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userDeletesNotificationOfResourceAndSubjectById(
		string $user,
		string $resource,
		string $subject,
	): void {
		$allNotifications = $this->listAllNotifications($user);
		$filteredNotificationId = $this->filterNotificationsBySubjectAndResource(
			$subject,
			$resource,
			$allNotifications,
		)[0]->notification_id;
		$this->featureContext->setResponse($this->deleteNotificationUsingId($user, $filteredNotificationId));
	}
}
