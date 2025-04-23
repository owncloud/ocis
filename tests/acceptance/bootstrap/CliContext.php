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
use Behat\Gherkin\Node\TableNode;
use GuzzleHttp\Exception\GuzzleException;
use PHPUnit\Framework\Assert;
use TestHelpers\CliHelper;
use TestHelpers\OcisConfigHelper;
use TestHelpers\BehatHelper;
use Psr\Http\Message\ResponseInterface;

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
	public function before(BeforeScenarioScope $scope): void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context
		$this->featureContext = BehatHelper::getContext($scope, $environment, 'FeatureContext');
		$this->spacesContext = BehatHelper::getContext($scope, $environment, 'SpacesContext');
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
	public function theAdministratorResetsThePasswordOfUserUsingTheCLI(
		string $status,
		string $user,
		string $password
	): void {
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
	public function theAdministratorDeletesEmptyTrashbinFoldersUsingTheCli(): void {
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
	public function theAdministratorChecksTheBackupConsistencyUsingTheCli(): void {
		$path = $this->featureContext->getStorageUsersRoot();
		$command = "backup consistency -p $path";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator creates auth-app token for user :user with expiration time :expirationTime using the auth-app CLI
	 * @When the administrator tries to create auth-app token for user :user with expiration time :expirationTime using the auth-app CLI
	 *
	 * @param string $user
	 * @param string $expirationTime
	 *
	 * @return void
	 */
	public function theAdministratorCreatesAppTokenForUserWithExpirationTimeUsingTheAuthAppCLI(
		string $user,
		string $expirationTime
	): void {
		$user = $this->featureContext->getActualUserName($user);
		$command = "auth-app create --user-name=$user --expiration=$expirationTime";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @Given the administrator has created app token for user :user with expiration time :expirationTime using the auth-app CLI
	 *
	 * @param string $user
	 * @param string $expirationTime
	 *
	 * @return void
	 */
	public function theAdministratorHasCreatedAppTokenForUserWithExpirationTimeUsingTheAuthAppCli(
		$user,
		$expirationTime
	): void {
		$user = $this->featureContext->getActualUserName($user);
		$command = "auth-app create --user-name=$user --expiration=$expirationTime";
		$body = [
			"command" => $command
		];
		$response = CliHelper::runCommand($body);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
		$jsonResponse = $this->featureContext->getJsonDecodedResponse($response);
		Assert::assertSame("OK", $jsonResponse["status"]);
		Assert::assertSame(
			0,
			$jsonResponse["exitCode"],
			"Expected exit code to be 0, but got " . $jsonResponse["exitCode"]
		);
		$output = $this->featureContext->substituteInLineCodes("App token created for $user");
		Assert::assertStringContainsString($output, $jsonResponse["message"]);
	}

	/**
	 * @Given user :user has created app token with expiration time :expirationTime using the auth-app CLI
	 *
	 * @param string $user
	 * @param string $expirationTime
	 *
	 * @return void
	 */
	public function userHasCreatedAppTokenWithExpirationTimeUsingTheAuthAppCLI(
		string $user,
		string $expirationTime
	): void {
		$user = $this->featureContext->getActualUserName($user);
		$command = "auth-app create --user-name=$user --expiration=$expirationTime";
		$body = [
			"command" => $command
		];

		$response = CliHelper::runCommand($body);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
		$jsonResponse = $this->featureContext->getJsonDecodedResponse($response);
		Assert::assertSame("OK", $jsonResponse["status"]);
		Assert::assertSame(
			0,
			$jsonResponse["exitCode"],
			"Expected exit code to be 0, but got " . $jsonResponse["exitCode"]
		);
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
	public function theAdministratorRemovesTheVersionsOfFilesInSpaceUsingSpaceId(string $space): void {
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
	 * @Then /^the command should be (successful|unsuccessful)$/
	 *
	 * @param string $successfulOrNot
	 *
	 * @return void
	 */
	public function theCommandShouldBeSuccessful(string $successfulOrNot): void {
		$response = $this->featureContext->getResponse();
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);

		$jsonResponse = $this->featureContext->getJsonDecodedResponse($response);

		$expectedStatus = 'OK';
		$expectedExitCode = 0;
		if ($successfulOrNot === "unsuccessful") {
			$expectedStatus = "ERROR";
			$expectedExitCode = 1;
		}

		Assert::assertSame($expectedStatus, $jsonResponse["status"]);
		Assert::assertSame(
			$expectedExitCode,
			$jsonResponse["exitCode"],
			"Expected exit code to be 0, but got " . $jsonResponse["exitCode"]
		);
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

	/**
	 * @When the administrator lists all the upload sessions
	 * @When the administrator lists all the upload sessions with flag :flag
	 *
	 * @param string|null $flag
	 *
	 * @return void
	 */
	public function theAdministratorListsAllTheUploadSessions(?string $flag = null): void {
		if ($flag) {
			$flag = "--$flag";
		}
		$command = "storage-users uploads sessions --json $flag";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator cleans upload sessions with the following flags:
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function theAdministratorCleansUploadSessionsWithTheFollowingFlags(TableNode $table): void {
		$flag = "";
		foreach ($table->getRows() as $row) {
			$flag .= "--$row[0] ";
		}
		$flagString = trim($flag);
		$command = "storage-users uploads sessions $flagString --clean --json";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator restarts the upload sessions that are in postprocessing
	 *
	 * @return void
	 */
	public function theAdministratorRestartsTheUploadSessionsThatAreInPostprocessing(): void {
		$command = "storage-users uploads sessions --processing --restart --json";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator restarts the upload sessions of file :file
	 *
	 * @param string $file
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function theAdministratorRestartsUploadSessionsOfFile(string $file): void {
		$response = CliHelper::runCommand(["command" => "storage-users uploads sessions --json"]);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
		$responseArray = $this->getJSONDecodedCliMessage($response);

		foreach ($responseArray as $item) {
			if ($item->filename === $file) {
				$uploadId = $item->id;
			}
		}

		$command = "storage-users uploads sessions --id=$uploadId --restart --json";
		$body = [
			"command" => $command
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @Then /^the CLI response (should|should not) contain these entries:$/
	 *
	 * @param string $shouldOrNot
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function theCLIResponseShouldContainTheseEntries(string $shouldOrNot, TableNode $table): void {
		$expectedFiles = $table->getColumn(0);
		$responseArray = $this->getJSONDecodedCliMessage($this->featureContext->getResponse());

		$resourceNames = [];
		foreach ($responseArray as $item) {
			if (isset($item->filename)) {
				$resourceNames[] = $item->filename;
			}
		}

		if ($shouldOrNot === "should not") {
			foreach ($expectedFiles as $expectedFile) {
				Assert::assertNotTrue(
					\in_array($expectedFile, $resourceNames),
					"The resource '$expectedFile' was found in the response."
				);
			}
		} else {
			foreach ($expectedFiles as $expectedFile) {
				Assert::assertTrue(
					\in_array($expectedFile, $resourceNames),
					"The resource '$expectedFile' was not found in the response."
				);
			}
		}
	}

	/**
	 * @param ResponseInterface $response
	 *
	 * @return array
	 * @throws JsonException
	 */
	public function getJSONDecodedCliMessage(ResponseInterface $response): array {
		$responseBody = $this->featureContext->getJsonDecodedResponse($response);

		// $responseBody["message"] contains a message info with the array of output json of the upload sessions command
		// Example Output: "INFO memory is not limited, skipping package=github.com/KimMachineGun/automemlimit/memlimit [{<output-json>}]"
		// So, only extracting the array of output json from the message
		\preg_match('/(\[.*\])/', $responseBody["message"], $matches);
		return \json_decode($matches[1], null, 512, JSON_THROW_ON_ERROR);
	}

	/**
	 * @AfterScenario @cli-uploads-sessions
	 *
	 * @return void
	 */
	public function cleanUploadsSessions(): void {
		$command = "storage-users uploads sessions --clean";
		$body = [
			"command" => $command
		];
		$response = CliHelper::runCommand($body);
		Assert::assertEquals("200", $response->getStatusCode(), "Failed to clean upload sessions");
	}

	/**
	 * @When /^the administrator triggers "([^"]*)" email notifications using the CLI$/
	 *
	 * @param string $interval
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theAdministratorTriggersEmailNotificationsUsingTheCLI(string $interval): void {
		$command = "notifications send-email --$interval";
		$body = [
			"command" => $command
		];

		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}
}
