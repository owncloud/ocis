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
use TestHelpers\OcisHelper;
use TestHelpers\OcmHelper;

/**
 * Acceptance test steps related to testing sharing ng features
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
	 * @When :user generates invitation
	 * @When :user generates invitation with email :email and description :description
	 *
	 * @param string $user
	 * @param string $email
	 * @param string $description
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userGeneratesInvitation(string $user, $email = null, $description = null): void {
		$response = OcmHelper::createInvitation(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$email,
			$description
		);
		$this->featureContext->setResponse($response);
		$this->invitationToken = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse())['token'];
	}

	/**
	 * @When :user accepts invitation
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userAcceptsInvitation(string $user): void {
		$response = OcmHelper::acceptInvitation(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$this->invitationToken
		);
		$this->featureContext->setResponse($response);
	}
}
