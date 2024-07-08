<?php declare(strict_types=1);
/**
 * @author Viktor Scharf <scharf.vi@gmail.com>
 *
 * @copyright Copyright (c) 2024, ownCloud GmbH
 * @license AGPL-3.0
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3,
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License, version 3,
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 *
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\OcisHelper;
use TestHelpers\OcmHelper;

/**
 * Acceptance test steps related to testing federation share(ocm) features
 */
class OcmContext implements Context {
	private FeatureContext $featureContext;
	private string $invitationToken;

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
		if (OcisHelper::isTestingOnReva()) {
			return;
		}
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context from here
		$this->featureContext = $environment->getContext('FeatureContext');
	}

	/**
	 * @return string
	 */
	public function getOcisDomain(): string {
		return $this->extractDomain(\getenv('TEST_SERVER_URL'));
	}

	/**
	 * @return string
	 */
	public function getFedOcisDomain(): string {
		return $this->extractDomain(\getenv('TEST_SERVER_FED_URL'));
	}

	/**
	 * @param string $url
	 *
	 * @return string
	 */
	public function extractDomain($url): string {
		if (!$url) {
			return "localhost";
		}
		return parse_url($url)["host"];
	}

	/**
	 * @param string $user
	 * @param string $email
	 * @param string $description
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function createInvitation(string $user, $email = null, $description = null): ResponseInterface {
		$response = OcmHelper::createInvitation(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$email,
			$description
		);
		$responseData = \json_decode($response->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
		if (isset($responseData["token"])) {
			$this->invitationToken = $responseData["token"];
		} else {
			throw new Exception(__METHOD__ . " response doesn't contain token");
		}
		return $response;
	}

	/**
	 * @When :user creates the federation share invitation
	 * @When :user creates the federation share invitation with email :email and description :description
	 *
	 * @param string $user
	 * @param string $email
	 * @param string $description
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userCreatesTheFederationShareInvitation(string $user, $email = null, $description = null): void {
		$this->featureContext->setResponse($this->createInvitation($user, $email, $description));
	}

	/**
	 * @Given :user has created the federation share invitation
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userHasCreatedTheFederationShareInvitation(string $user): void {
		$response = $this->createInvitation($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function acceptInvitation(string $user): ResponseInterface {
		$providerDomain = ($this->featureContext->getCurrentServer() === "LOCAL") ? $this->getFedOcisDomain() : $this->getOcisDomain();
		return OcmHelper::acceptInvitation(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$this->invitationToken,
			$providerDomain
		);
	}

	/**
	 * @When :user accepts the last federation share invitation
	 * @When :user tries to accept the last federation share invitation
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userAcceptsTheLastFederationShareInvitation(string $user): void {
		$this->featureContext->setResponse($this->acceptInvitation($user));
	}

	/**
	 * @Given :user has accepted invitation
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userHasAcceptedTheLastFederationShareInvitation(string $user): void {
		$response = $this->acceptInvitation($user);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @param string $user
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function findAcceptedUsers(string $user): ResponseInterface {
		return OcmHelper::findAcceptedUsers(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user)
		);
	}

	/**
	 * @When :user searches for accepted users
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userFindsAcceptedUsers(string $user): void {
		$this->featureContext->setResponse($this->findAcceptedUsers($user));
	}
}
