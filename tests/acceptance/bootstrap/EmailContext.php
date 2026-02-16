<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Amrita Shrestha <amrita@jankaritech.com>
 * @copyright Copyright (c) 2025 Amrita Shrestha amrita@jankaritech.com
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Behat\Gherkin\Node\PyStringNode;
use PHPUnit\Framework\Assert;
use GuzzleHttp\Exception\GuzzleException;
use TestHelpers\EmailHelper;
use TestHelpers\GraphHelper;
use TestHelpers\BehatHelper;

require_once 'bootstrap.php';

/**
 * Defines application features from the specific context.
 */
class EmailContext implements Context {
	private FeatureContext $featureContext;
	private SpacesContext $spacesContext;

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
		$this->spacesContext = BehatHelper::getContext($scope, $environment, 'SpacesContext');
	}

	/**
	 * @BeforeScenario @email
	 * @AfterScenario @email
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function clearAllEmails(): void {
		EmailHelper::deleteAllEmails();
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
	public function userShouldHaveReceivedTheFollowingEmailFromUserAboutTheShareOfProjectSpace(
		string $user,
		string $sender,
		string $spaceName,
		PyStringNode $content,
	): void {
		$rawExpectedEmailBodyContent = \str_replace("\r\n", "\n", $content->getRaw());
		$this->featureContext->setResponse(
			GraphHelper::getMySpaces(
				$this->featureContext->getBaseUrl(),
				$user,
				$this->featureContext->getPasswordForUser($user),
				'',
			),
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
					"parameter" => [$sender, $spaceName],
				],
			],
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
	public function userShouldHaveReceivedTheFollowingEmailFromUser(
		string $user,
		string $sender,
		PyStringNode $content,
	): void {
		$rawExpectedEmailBodyContent = \str_replace("\r\n", "\n", $content->getRaw());
		$expectedEmailBodyContent = $this->featureContext->substituteInLineCodes(
			$rawExpectedEmailBodyContent,
			$sender,
		);
		$this->assertEmailContains($user, $expectedEmailBodyContent);
	}

	/**
	 * @Then user :user should have :count emails
	 *
	 * @param string $user
	 * @param int $count
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userShouldHaveEmail(string $user, int $count): void {
		$address = $this->featureContext->getEmailAddressForUser($user);
		$query = 'to:' . $address;
		$response = EmailHelper::searchEmails($query);
		$emails = $this->featureContext->getJsonDecodedResponse($response);
		if ($emails["messages_count"] <= $count) {
			echo "[INFO] Mailbox is empty...\n";
			// Wait for 1 second and try again
			// the mailbox might not be created yet
			sleep(1);
			$response = EmailHelper::searchEmails($query);
			$emails = $this->featureContext->getJsonDecodedResponse($response);
		}

		Assert::assertSame(
			$count,
			$emails["messages_count"],
			"Expected '$address' received mail total '$count' email but got " . $emails["messages_count"] . " email",
		);
	}

	/**
	 * @Then user :user should have received the following email from user :sender ignoring whitespaces
	 *
	 * @param string $user
	 * @param string $sender
	 * @param PyStringNode $content
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldHaveReceivedTheFollowingEmailFromUserIgnoringWhitespaces(
		string $user,
		string $sender,
		PyStringNode $content,
	): void {
		$rawExpectedEmailBodyContent = \str_replace("\r\n", "\n", $content->getRaw());
		$expectedEmailBodyContent = $this->featureContext->substituteInLineCodes(
			$rawExpectedEmailBodyContent,
			$sender,
		);
		$this->assertEmailContains($user, $expectedEmailBodyContent, true);
	}

	/***
	 * @param string $user
	 * @param string $expectedEmailBodyContent
	 * @param bool $ignoreWhiteSpace
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function assertEmailContains(
		string $user,
		string $expectedEmailBodyContent,
		$ignoreWhiteSpace = false,
	): void {
		$address = $this->featureContext->getEmailAddressForUser($user);
		$actualEmailBodyContent = $this->getBodyOfLastEmail($address);
		if ($ignoreWhiteSpace) {
			$expectedEmailBodyContent = preg_replace('/\s+/', '', $expectedEmailBodyContent);
			$actualEmailBodyContent = preg_replace('/\s+/', '', $actualEmailBodyContent);
		}
		Assert::assertStringContainsString(
			$expectedEmailBodyContent,
			$actualEmailBodyContent,
			"The email address '$address' should have received an"
			. "email with the body containing $expectedEmailBodyContent
			but the received email is $actualEmailBodyContent",
		);
	}

	/**
	 * Returns the body of the last received email for the provided receiver according to the provided email address and the serial number
	 * For email number, 1 means the latest one
	 *
	 * @param string $emailAddress
	 * @param int|null $waitTimeSec Time to wait for the email if the email has been delivered
	 *
	 * @return string
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function getBodyOfLastEmail(
		string $emailAddress,
		?int $waitTimeSec = EMAIL_WAIT_TIMEOUT_SEC,
	): string {
		$currentTime = \time();
		$endTime = $currentTime + $waitTimeSec;
		while ($currentTime <= $endTime) {
			$query = 'to:' . $emailAddress;
			$mailResponse = $this->featureContext->getJsonDecodedResponse(
				EmailHelper::searchEmails($query),
			);
			if ($mailResponse["messages_count"] > 0) {
				$lastEmail = $this->featureContext->getJsonDecodedResponse(
					EmailHelper::getEmailById("latest", $query),
				);
				$body = \str_replace(
					"\r\n",
					"\n",
					\quoted_printable_decode($lastEmail["Text"] . "\n" . $lastEmail["HTML"]),
				);
				return $body;
			}
			\usleep(STANDARD_SLEEP_TIME_MICROSEC * 50);
			$currentTime = \time();
		}
		throw new Exception("Could not find the email to the address: " . $emailAddress);
	}

	/**
	 * @Then user :user should have received the following grouped email
	 *
	 * @param string $user
	 * @param PyStringNode $content
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function userShouldHaveReceivedTheFollowingGroupedEmail(
		string $user,
		PyStringNode $content,
	): void {
		$rawExpectedEmailBodyContent = \str_replace("\r\n", "\n", $content->getRaw());
		$expectedEmailBodyContent = $this->featureContext->substituteInLineCodes(
			$rawExpectedEmailBodyContent,
		);
		$this->assertEmailContains($user, $expectedEmailBodyContent);
	}
}
