<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Sajan Gurung <sajan@jankaritech.com>
 * @copyright Copyright (c) 2024 Sajan Gurung sajan@jankaritech.com
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

use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Behat\Behat\Context\Context;
use PHPUnit\Framework\Assert;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\CliHelper;
use TestHelpers\OcisConfigHelper;

/**
 * CLI context
 */
class CliContext implements Context {
	private FeatureContext $featureContext;

	/**
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 */
	public function setUpScenario(BeforeScenarioScope $scope): void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context
		$this->featureContext = $environment->getContext('FeatureContext');
	}

	/**
	 * @Given the administrator has stopped the server
	 *
	 * @return void
	 */
	public function theAdministratorHasStoppedTheServer(): void {
		$response = OcisConfigHelper::stopOcis();
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @Given the administrator has started the server
	 *
	 * @return void
	 */
	public function theAdministratorHasStartedTheServer(): void {
		$response = OcisConfigHelper::startOcis();
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
	}

	/**
	 * @When /^the administrator resets the password of (non-existing|existing) user "([^"]*)" to "([^"]*)" using the CLI$/
	 *
	 * @param string $status
	 * @param string $user
	 * @param string $password
	 *
	 * @return void
	 */
	public function theAdministratorResetsThePasswordOfUserUsingTheCLI(string $status, string $user, string $password): void {
		$command = "idm resetpassword -u $user";
		$body = [
			"command" => $command,
			"inputs" => [$password, $password]
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
		if ($status === "non-existing") {
			return;
		}
		$this->featureContext->updateUserPassword($user, $password);
	}

	/**
	 * @Then the command should be successful
	 *
	 * @return void
	 */
	public function theCommandShouldBeSuccessful(): void {
		$response = $this->featureContext->getResponse();
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);

		$jsonResponse = $this->featureContext->getJsonDecodedResponse($response);

		Assert::assertSame("OK", $jsonResponse["status"]);
		Assert::assertSame(0, $jsonResponse["exitCode"], "Expected exit code to be 0, but got " . $jsonResponse["exitCode"]);
	}

	/**
	 * @Then /^the command output (should|should not) contain "([^"]*)"$/
	 *
	 * @param string $shouldOrNot
	 * @param string $output
	 *
	 * @return void
	 */
	public function theCommandOutputShouldContain(string $shouldOrNot, string $output): void {
		$response = $this->featureContext->getResponse();
		$jsonResponse = $this->featureContext->getJsonDecodedResponse($response);

		if ($shouldOrNot === "should") {
			Assert::assertStringContainsString($output, $jsonResponse["message"]);
		} else {
			Assert::assertStringNotContainsString($output, $jsonResponse["message"]);
		}
	}
}
