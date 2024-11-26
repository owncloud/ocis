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
use GuzzleHttp\Exception\GuzzleException;
use TestHelpers\AuthAppHelper;

require_once 'bootstrap.php';

/**
 * AuthApp context
 */
class AuthAppContext implements Context {
	private FeatureContext $featureContext;
	private array $allCreatedTokens = [];

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
	 * @When the administrator creates app token with expiration time :expiration using the API
	 *
	 * @param string $expiration
	 *
	 * @return void
	 */
	public function theAdministratorCreatesAppTokenForUserWithExpirationTimeUsingTheApi(string $expiration): void {
		$this->featureContext->setResponse(
			AuthAppHelper::createAppAuthToken(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getAdminUsername(),
				$this->featureContext->getAdminPassword(),
				$expiration,
			)
		);
	}

	/**
	 * @Given the administrator has created app token with expiration time :expiration using the API
	 *
	 * @param string $expiration
	 *
	 * @return void
	 */
	public function theAdministratorHasCreatedAppTokenWithExpirationTimeUsingTheApi(string $expiration): void {
		$response = AuthAppHelper::createAppAuthToken(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$expiration,
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "", $response);
	}

	/**
	 * @When admin lists all created tokens
	 *
	 * @return void
	 */
	public function adminListsAllCreatedTokens(): void {
		$this->featureContext->setResponse(
			AuthAppHelper::listAllAppAuthToken(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getAdminUsername(),
				$this->featureContext->getAdminPassword(),
			)
		);
	}

	/**
	 * @return void
	 */
	public function deleteAllToken() : void {
		$response = AuthAppHelper::listAllAppAuthToken(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
		);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "", $response);
		$rawBody =  $response->getBody()->getContents();
		$tokens = json_decode($rawBody);
		foreach ($tokens as $token) {
			$this->featureContext->theHTTPStatusCodeShouldBe(
				200,
				"",
				AuthAppHelper::deleteAppAuthToken(
					$this->featureContext->getBaseUrl(),
					$this->featureContext->getAdminUsername(),
					$this->featureContext->getAdminPassword(),
					$token->token
				)
			);
		}
	}

	/**
	 * @AfterScenario
	 *
	 * @return void
	 *
	 * @throws Exception|GuzzleException
	 */
	public function cleanDataAfterTests(): void {
		$this->deleteAllToken();
	}
}
