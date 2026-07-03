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
use PHPUnit\Framework\Assert;
use TestHelpers\BehatHelper;

require_once 'bootstrap.php';

/**
 * Context for ocis vault specific steps
 */
class VaultContext implements Context {
	private FeatureContext $featureContext;

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
	 * @Then user :user should have acr value :acr
	 *
	 * @param string $user
	 * @param string $acr
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userShouldHaveAcrVaule(string $user, string $acr): void {
		$access_token = $this->featureContext->getOcisUserToken($user)['token']['accessToken'];

		// Decode JWT token
		$parts = explode('.', $access_token);
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
			"Exected acr value to be $acr but got $actualAcr",
		);
	}
}
