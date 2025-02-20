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
	 * @AfterScenario @email
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function clearAllEmails(): void {
		try {
			EmailHelper::deleteAllEmails($this->featureContext->getStepLineRef());
		} catch (Exception $e) {
			echo __METHOD__ .
				" could not delete email messages?\n" .
				$e->getMessage();
		}
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
		$this->spacesContext = BehatHelper::getContext($scope, $environment, 'SpacesContext');
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
		PyStringNode $content
	): void {
		$rawExpectedEmailBodyContent = \str_replace("\r\n", "\n", $content->getRaw());
		$this->featureContext->setResponse(
			GraphHelper::getMySpaces(
				$this->featureContext->getBaseUrl(),
				$user,
				$this->featureContext->getPasswordForUser($user),
				'',
				$this->featureContext->getStepLineRef()
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
			]
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
		PyStringNode $content
	): void {
		$rawExpectedEmailBodyContent = \str_replace("\r\n", "\n", $content->getRaw());
		$expectedEmailBodyContent = $this->featureContext->substituteInLineCodes(
			$rawExpectedEmailBodyContent,
			$sender
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
		$mailResponse = EmailHelper::searchEmails($query, $this->featureContext->getStepLineRef());
		$totalMail = $this->featureContext->getJsonDecodedResponse($mailResponse)["messages_count"];
		Assert::assertSame(
			$count,
			$totalMail,
			"Expected '$address' received mail total '$count' mail but got '$totalMail' mail"
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
		PyStringNode $content
	): void {
		$rawExpectedEmailBodyContent = \str_replace("\r\n", "\n", $content->getRaw());
		$expectedEmailBodyContent = $this->featureContext->substituteInLineCodes(
			$rawExpectedEmailBodyContent,
			$sender
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
		$ignoreWhiteSpace = false
	): void {
		$address = $this->featureContext->getEmailAddressForUser($user);
		$this->featureContext->pushEmailRecipientAsMailBox($address);
		$actualEmailBodyContent = EmailHelper::getBodyOfLastEmail($address, $this->featureContext->getStepLineRef());
		if ($ignoreWhiteSpace) {
			$expectedEmailBodyContent = preg_replace('/\s+/', '', $expectedEmailBodyContent);
			$actualEmailBodyContent = preg_replace('/\s+/', '', $actualEmailBodyContent);
		}
		Assert::assertStringContainsString(
			$expectedEmailBodyContent,
			$actualEmailBodyContent,
			"The email address '$address' should have received an"
			. "email with the body containing $expectedEmailBodyContent
			but the received email is $actualEmailBodyContent"
		);
	}

}
