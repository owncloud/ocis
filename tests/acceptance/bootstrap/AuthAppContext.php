<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Niraj Acharya <niraj@jankaritech.com>
 * @copyright Copyright (c) 2024 Niraj Acharya niraj@jankaritech.com
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
use TestHelpers\BehatHelper;
use PHPUnit\Framework\Assert;
use TestHelpers\AuthAppHelper;

require_once 'bootstrap.php';

/**
 * AuthApp context
 */
class AuthAppContext implements Context {
	private FeatureContext $featureContext;
	private array $lastCreatedToken = [];
	private bool $usingAuthAppToken = false;

	/**
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
	}

	/**
	 * @return array
	 */
	public function getLastCreatedToken(): array {
		return $this->lastCreatedToken;
	}

	/**
	 * @return bool
	 */
	public function isUsingAuthAppToken(): bool {
		return $this->usingAuthAppToken;
	}

	/**
	 * @When /^user "([^"]*)" (?:creates|tries to create) auth-app token with expiration time "([^"]*)" using the auth-app API$/
	 *
	 * @param string $user
	 * @param string $expiration
	 *
	 * @return void
	 */
	public function userCreatesAppTokenWithExpirationTimeUsingTheAuthAppApi(string $user, string $expiration): void {
		$this->featureContext->setResponse(
			AuthAppHelper::createAppAuthToken(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getActualUsername($user),
				$this->featureContext->getPasswordForUser($user),
				["expiry" => $expiration],
			),
		);
	}

	/**
	 * @When /^user "([^"]*)" (?:creates|tries to create) auth-app token with expiration time "([^"]*)" and label "(.*)" using the auth-app API$/
	 *
	 * @param string $user
	 * @param string $expiration
	 * @param string $label
	 *
	 * @return void
	 */
	public function userCreatesAuthAppTokenWithExpirationTimeAndLabelUsingTheAuthAppApi(
		string $user,
		string $expiration,
		string $label,
	): void {
		$this->featureContext->setResponse(
			AuthAppHelper::createAppAuthToken(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getActualUsername($user),
				$this->featureContext->getPasswordForUser($user),
				["expiry" => $expiration, "label" => $label],
			),
		);
	}

	/**
	 * @Given user :user has created auth-app token with expiration time :expiration using the auth-app API
	 *
	 * @param string $user
	 * @param string $expiration
	 *
	 * @return void
	 */
	public function userHasCreatedAppTokenWithExpirationTime(string $user, string $expiration): void {
		$response = AuthAppHelper::createAppAuthToken(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			["expiry" => $expiration],
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "", $response);
		$this->lastCreatedToken = [
			"user" => strtolower($user),
			"token" => json_decode($response->getBody()->getContents())->token,
		];
	}

	/**
	 * @Given user :user has waited :expiration second(s) for auth-app token to expire
	 *
	 * @param string $user
	 * @param string $expiration
	 *
	 * @return void
	 */
	public function userWaitSecondForAuthAppTokenToExpire($user, $expiration): void {
		sleep((int)$expiration);
	}

	/**
	 * @When user :user lists all created tokens using the auth-app API
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function userListsAllCreatedTokensUsingTheAuthAppApi(string $user): void {
		$this->featureContext->setResponse(
			AuthAppHelper::listAllAppAuthTokensForUser(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getActualUsername($user),
				$this->featureContext->getPasswordForUser($user),
			),
		);
	}

	/**
	 * @Given user :user has created auth-app token for user :impersonatedUser with expiration time :expiration using the auth-app API
	 *
	 * @param string $user
	 * @param string $impersonatedUser
	 * @param string $expiration
	 *
	 * @return void
	 */
	public function userHasCreatedAppTokenWithExpirationTimeImpersonatingUserUsingTheAuthAppApi(
		string $user,
		string $impersonatedUser,
		string $expiration,
	): void {
		$response = AuthAppHelper::createAppAuthToken(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
			[
				"expiry" => $expiration,
				"userName" => $this->featureContext->getActualUsername($impersonatedUser),
			],
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(
			200,
			"Failed creating auth-app token\n"
			. "HTTP status code 200 is not the expected value " . $response->getStatusCode(),
			$response,
		);
		$this->lastCreatedToken = [
			"user" => strtolower($impersonatedUser),
			"token" => json_decode($response->getBody()->getContents())->token,
		];
	}

	/**
	 * @When /^user "([^"]*)" (?:creates|tries to create) auth-app token for user "([^"]*)" with expiration time "([^"]*)" using the auth-app API$/
	 *
	 * @param string $user
	 * @param string $impersonatedUser
	 * @param string $expiration
	 *
	 * @return void
	 */
	public function theUserCreatesAppTokenForUserWithExpirationTimeViaAuthAppApi(
		string $user,
		string $impersonatedUser,
		string $expiration,
	): void {
		$this->featureContext->setResponse(
			AuthAppHelper::createAppAuthToken(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getActualUsername($user),
				$this->featureContext->getPasswordForUser($user),
				[
					"expiry" => $expiration,
					"userName" => $this->featureContext->getActualUsername($impersonatedUser),
				],
			),
		);
	}

	/**
	 * @When user :user deletes all the created auth-app tokens using the auth-app API
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function userDeletesAllCreatedAuthAppTokenUsingAuthAppAPI(string $user): void {
		$baseUrl = $this->featureContext->getBaseUrl();
		$user = $this->featureContext->getActualUsername($user);
		$password = $this->featureContext->getPasswordForUser($user);
		$response = AuthAppHelper::listAllAppAuthTokensForUser(
			$baseUrl,
			$user,
			$password,
		);
		$authAppTokens = json_decode($response->getBody()->getContents());
		foreach ($authAppTokens as $tokenObj) {
			$deleteResponse = AuthAppHelper::deleteAppAuthToken(
				$baseUrl,
				$user,
				$password,
				$tokenObj->token,
			);
			$this->featureContext->setResponse($deleteResponse);
			$this->featureContext->pushToLastHttpStatusCodesArray((string)$deleteResponse->getStatusCode());
		}
	}

	/**
	 * @Then user :user should have :count auth-app tokens
	 *
	 * @param string $user
	 * @param integer $count
	 *
	 * @return void
	 */
	public function userShouldHaveAuthAppTokens(string $user, int $count): void {
		$response = AuthAppHelper::listAllAppAuthTokensForUser(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getActualUsername($user),
			$this->featureContext->getPasswordForUser($user),
		);
		$authAppTokens = json_decode($response->getBody()->getContents());
		Assert::assertCount(
			$count,
			$authAppTokens,
			"Expected the count to be $count but got " . \count($authAppTokens),
		);
	}

	/**
	 * @When user :user creates app token with user-id for user :impersonatedUser with expiration time :expiration using the auth-app API
	 * @When user :user tries to create app token with user-id for user :impersonatedUser with expiration time :expiration using the auth-app API
	 *
	 * @param string $user
	 * @param string $impersonatedUser
	 * @param string $expiration
	 *
	 * @return void
	 */
	public function userCreatesAppTokenWithUserIdForUserWithExpirationTimeUsingTheAuthAppApi(
		string $user,
		string $impersonatedUser,
		string $expiration,
	): void {
		$this->featureContext->setResponse(
			AuthAppHelper::createAppAuthToken(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getActualUsername($user),
				$this->featureContext->getPasswordForUser($user),
				[
					"expiry" => $expiration,
					"userID" => $this->featureContext->getUserIdByUserName($impersonatedUser),
				],
			),
		);
	}
	/**
	 * @When user :user tries to delete the last created auth-app token using the auth-app API
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function userTriesToDeleteLastCreatedAuthAppTokensUsingTheAuthAppApi(string $user): void {
		$baseUrl = $this->featureContext->getBaseUrl();
		$user = $this->featureContext->getActualUsername($user);
		$password = $this->featureContext->getPasswordForUser($user);
		$deleteResponse = AuthAppHelper::deleteAppAuthToken(
			$baseUrl,
			$user,
			$password,
			$this->getLastCreatedToken()['token'],
		);
		$this->featureContext->setResponse($deleteResponse);
		$this->featureContext->pushToLastHttpStatusCodesArray((string)$deleteResponse->getStatusCode());
	}

	/**
	 * @Given using auth-app token
	 *
	 * @return void
	 */
	public function usingAuthAppToken(): void {
		$this->usingAuthAppToken = true;
	}
}
