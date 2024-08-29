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
use TestHelpers\CliHelper;
use TestHelpers\OcisConfigHelper;

/**
 * CLI context
 */
class CliContext implements Context {
	private FeatureContext $featureContext;
	private SpacesContext $spacesContext;

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
		$this->spacesContext = $environment->getContext('SpacesContext');
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
	 * @Given /^the administrator (?:starts|has started) the server$/
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
	 * @When the administrator deletes the empty trashbin folders using the CLI
	 *
	 * @return void
	 */
	public function theAdministratorDeletesEmptyTrashbinFoldersUsingTheCli():void {
		$path = $this->featureContext->getStorageUsersRoot();
		$command = "trash purge-empty-dirs -p $path --dry-run=false";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator checks the backup consistency using the CLI
	 *
	 * @return void
	 */
	public function theAdministratorChecksTheBackupConsistencyUsingTheCli():void {
		$path = $this->featureContext->getStorageUsersRoot();
		$command = "backup consistency -p $path";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator removes all the file versions using the CLI
	 *
	 * @return void
	 */
	public function theAdministratorRemovesAllVersionsOfResources() {
		$path = $this->featureContext->getStorageUsersRoot();
		$command = "revisions purge -p $path --dry-run=false";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator removes the versions of file :file of user :user from space :space using the CLI
	 *
	 * @param string $file
	 * @param string $user
	 * @param string $space
	 *
	 * @return void
	 */
	public function theAdministratorRemovesTheVersionsOfFileUsingFileId($file, $user, $space) {
		$path = $this->featureContext->getStorageUsersRoot();
		$fileId = $this->spacesContext->getFileId($user, $space, $file);
		$command = "revisions purge -p $path -r $fileId --dry-run=false";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When /^the administrator reindexes all spaces using the CLI$/
	 *
	 * @return void
	 */
	public function theAdministratorReindexesAllSpacesUsingTheCli(): void {
		$command = "search index --all-spaces";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When /^the administrator reindexes a space "([^"]*)" using the CLI$/
	 *
	 * @param string $spaceName
	 *
	 * @return void
	 */
	public function theAdministratorReindexesASpaceUsingTheCli(string $spaceName): void {
		$spaceId = $this->spacesContext->getSpaceIdByName($this->featureContext->getAdminUsername(), $spaceName);
		$command = "search index --space $spaceId";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator removes the file versions of space :space using the CLI
	 *
	 * @param string $space
	 *
	 * @return void
	 */
	public function theAdministratorRemovesTheVersionsOfFilesInSpaceUsingSpaceId(string $space):void {
		$path = $this->featureContext->getStorageUsersRoot();
		$adminUsername = $this->featureContext->getAdminUsername();
		$spaceId = $this->spacesContext->getSpaceIdByName($adminUsername, $space);
		$command = "revisions purge -p $path -r $spaceId --dry-run=false";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
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
		$output = $this->featureContext->substituteInLineCodes($output);

		if ($shouldOrNot === "should") {
			Assert::assertStringContainsString($output, $jsonResponse["message"]);
		} else {
			Assert::assertStringNotContainsString($output, $jsonResponse["message"]);
		}
	}
}
