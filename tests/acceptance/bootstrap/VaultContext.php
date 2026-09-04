<?php

declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Prajwol Amatya <prajwol@jankaritech.com>
 * @copyright Copyright (c) 2026 Prajwol Amatya prajwol@jankaritech.com
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
use GuzzleHttp\Exception\GuzzleException;
use PHPUnit\Framework\Assert;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\BehatHelper;
use TestHelpers\GraphHelper;
use TestHelpers\HttpRequestHelper;
use TestHelpers\KeycloakHelper;
use TestHelpers\SettingsHelper;

require_once 'bootstrap.php';

/**
 * Context for ocis vault specific steps
 */
class VaultContext implements Context {
	private FeatureContext $featureContext;
	private array $originalRealmAttributes = [];

	/**
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 *
	 * @throws Exception
	 */
	public function before(BeforeScenarioScope $scope): void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context
		$this->featureContext = BehatHelper::getContext($scope, $environment, 'FeatureContext');
	}

	/**
	 * @Given the administrator has set the Keycloak realm attribute :key to :value
	 *
	 * @param string $key
	 * @param string $value
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function theAdministratorHasSetTheKeycloakRealmAttributeTo(string $key, string $value): void {
		$realm = KeycloakHelper::getRealm();
		$this->originalRealmAttributes[$key] = $realm['attributes'][$key] ?? null;
		$response = KeycloakHelper::updateRealmAttribute($key, $value);
		Assert::assertEquals(
			204,
			$response->getStatusCode(),
			"Failed to update Keycloak realm attribute $key. Response: " . $response->getBody()->getContents(),
		);
	}

	/**
	 * @AfterScenario @keycloak-config
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function restoreKeycloakRealmAttributes(): void {
		foreach ($this->originalRealmAttributes as $key => $value) {
			if ($value === null) {
				KeycloakHelper::deleteRealmAttribute($key);
				continue;
			}
			KeycloakHelper::updateRealmAttribute($key, $value);
		}
		$this->originalRealmAttributes = [];
	}

	/**
	 * @Then user :user should have a JWT token with an ACR value :acr
	 *
	 * @param string $user
	 * @param string $acr
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldHaveAJwtTokenWithAnAcrValue(string $user, string $acr): void {
		$accessToken = $this->featureContext->getOcisUserToken($user)['token']['accessToken'];

		// Decode JWT token
		$parts = explode('.', $accessToken);
		if (\count($parts) !== 3) {
			throw new Exception("Invalid JWT token format.");
		}
		$payload = $parts[1];
		$decodedPayload = base64_decode(strtr($payload, '-_', '+/'), true);
		$payloadArray = json_decode($decodedPayload, true);
		$actualAcr = $payloadArray['acr'] ?? null;
		Assert::assertEquals(
			$acr,
			$actualAcr,
			"Expected acr value to be $acr but got $actualAcr",
		);
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	private function getPermissionsList(string $user): ResponseInterface {
		$password = $this->featureContext->getPasswordForUser($user);
		$headers = [];
		$authUser = $user;
		if (KeycloakHelper::isTestingWithKeycloak()) {
			$this->authenticateKeycloakUserIfNeeded($user);
			$accessToken = $this->featureContext->getOcisUserToken($user)['token']['accessToken'];
			$headers['Authorization'] = 'Bearer ' . $accessToken;
			$authUser = null;
			$password = null;
		}
		$userId = $this->featureContext->getAttributeOfCreatedUser($user, 'id');
		return SettingsHelper::getPermissionsList(
			$this->featureContext->getBaseUrl(),
			$authUser,
			$password,
			$userId,
			$headers,
		);
	}

	/**
	 * Authenticates a Keycloak-backed user directly via the OIDC token endpoint (no browser,
	 * no MFA/vault-mode UI setup) so that the user's oCIS account (and id) is provisioned even
	 * for roles that don't have vault UI elements to interact with, e.g. User Light.
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	private function authenticateKeycloakUserIfNeeded(string $user): void {
		if ($this->featureContext->getAttributeOfCreatedUser($user, 'id')) {
			return;
		}
		$userAttribute = $this->featureContext->getCreatedKeycloakUsers()[strtolower($user)];
		$tokenData = KeycloakHelper::setAccessTokenForKeycloakOcisUser($userAttribute);
		$this->featureContext->setOcisUserToken($userAttribute, $tokenData);

		$response = HttpRequestHelper::get(
			GraphHelper::getFullUrl($this->featureContext->getBaseUrl(), 'me'),
			null,
			null,
			['Authorization' => 'Bearer ' . $tokenData['access_token']],
		);
		$userAttribute['id'] = $this->featureContext->getJsonDecodedResponse($response)['id'];
		$this->featureContext->addUserToCreatedUsersList(
			$user,
			$userAttribute['password'],
			$userAttribute['displayName'],
			$userAttribute['email'],
			$userAttribute['id'],
		);
	}

	/**
	 * @When /^user "([^"]*)" gets the permissions list using the settings API$/
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userGetsPermissionsList(string $user): void {
		$this->featureContext->setResponse($this->getPermissionsList($user));
	}
}
