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
use Behat\Gherkin\Node\PyStringNode;
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
	 * @Given the administrator has cleaned the search service data
	 *
	 * @return void
	 */
	public function theAdministratorHasCleanedTheSearchServiceData(): void {
		$dataPath = $this->featureContext->getOcisDataPath() . "/search";
		\exec("rm -rf $dataPath");
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
		string $password,
	): void {
		$command = "idm resetpassword -u $user";
		$body = [
			"command" => $command,
			"inputs" => [$password, $password],
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
			"command" => $command,
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
			"command" => $command,
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator lists all the unified roles using the CLI
	 *
	 * @return void
	 */
	public function theAdministratorListsAllTheUnifiedRolesUsingTheCli(): void {
		$command = "graph list-unified-roles";
		$body = [
			"command" => $command,
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
		string $expirationTime,
	): void {
		$user = $this->featureContext->getActualUserName($user);
		$command = "auth-app create --user-name=$user --expiration=$expirationTime";
		$body = [
			"command" => $command,
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
		$expirationTime,
	): void {
		$user = $this->featureContext->getActualUserName($user);
		$command = "auth-app create --user-name=$user --expiration=$expirationTime";
		$body = [
			"command" => $command,
		];
		$response = CliHelper::runCommand($body);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
		$jsonResponse = $this->featureContext->getJsonDecodedResponse($response);
		Assert::assertSame("OK", $jsonResponse["status"]);
		Assert::assertSame(
			0,
			$jsonResponse["exitCode"],
			"Expected exit code to be 0, but got " . $jsonResponse["exitCode"],
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
		string $expirationTime,
	): void {
		$user = $this->featureContext->getActualUserName($user);
		$command = "auth-app create --user-name=$user --expiration=$expirationTime";
		$body = [
			"command" => $command,
		];

		$response = CliHelper::runCommand($body);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
		$jsonResponse = $this->featureContext->getJsonDecodedResponse($response);
		Assert::assertSame("OK", $jsonResponse["status"]);
		Assert::assertSame(
			0,
			$jsonResponse["exitCode"],
			"Expected exit code to be 0, but got " . $jsonResponse["exitCode"],
		);
	}

	/**
	 * @When the administrator removes all the file versions using the CLI
	 *
	 * @return void
	 */
	public function theAdministratorRemovesAllVersionsOfResources(): void {
		$path = $this->featureContext->getStorageUsersRoot();
		$command = "revisions purge -p $path --dry-run=false";
		$body = [
			"command" => $command,
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
	public function theAdministratorRemovesTheVersionsOfFileUsingFileId($file, $user, $space): void {
		$path = $this->featureContext->getStorageUsersRoot();
		$fileId = $this->spacesContext->getFileId($user, $space, $file);
		$command = "revisions purge -p $path -r $fileId --dry-run=false";
		$body = [
			"command" => $command,
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
			"command" => $command,
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
			"command" => $command,
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
			"command" => $command,
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
		$this->featureContext->theHTTPStatusCodeShouldBe(200);
		$jsonResponse = $this->featureContext->getJsonDecodedResponse();

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
			"Expected exit code to be 0, but got " . $jsonResponse["exitCode"],
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
	 * @Then the command output should include the following roles:
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theCommandOutputShouldIncludeTheFollowingRoles(TableNode $table): void {
		$this->featureContext->verifyTableNodeColumns(
			$table,
			['LABEL', 'ENABLED', 'DESCRIPTION'],
		);

		$expectedRoles = $table->getColumnsHash();
		$response = $this->featureContext->getResponse();
		$decodedResponse = $this->featureContext->getJsonDecodedResponse($response);

		// Regex pattern to extract LABEL, ENABLED, and DESCRIPTION
		$pattern = '/│\s*\d+\s*│\s*([^\|]+?)\s*│\s*([a-f0-9\-]{36})\s*│\s*(enabled|disabled)\s*│\s*(.*?)\s*│/i';
		preg_match_all($pattern, $decodedResponse['message'], $matches, PREG_SET_ORDER);

		$actualRoles = [];
		foreach ($matches as $match) {
			$actualRoles[] = [
				'LABEL' => trim($match[1]),
				'ENABLED' => trim($match[3]),
				'DESCRIPTION' => trim($match[4]),
			];
		}

		// Compare expected roles with actual roles by LABEL and assert equality
		foreach ($expectedRoles as $expected) {
			$label = $expected['LABEL'];
			$actual = null;
			foreach ($actualRoles as $role) {
				if ($role['LABEL'] === $label) {
					$actual = $role;
					break;
				}
			}

			Assert::assertNotNull(
				$actual,
				"Role with LABEL '$label' not found in command output.",
			);

			Assert::assertEquals(
				$expected,
				$actual,
				"Mismatch for LABEL '$label':\nExpected: " . json_encode($expected) .
				"\nActual: " . json_encode($actual),
			);
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
			"command" => $command,
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
			"command" => $command,
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
			"command" => $command,
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When /^the administrator (restarts|resumes) the upload session of file "([^"]*)" using the CLI$/
	 *
	 * @param string $flag
	 * @param string $file
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function theAdministratorResumesOrRestartsUploadSessionOfFileUsingTheCLI(string $flag, string $file): void {
		$flag = rtrim($flag, "s");
		$response = CliHelper::runCommand(["command" => "storage-users uploads sessions --json"]);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $response);
		$responseArray = $this->getJSONDecodedCliMessage($response);

		foreach ($responseArray as $item) {
			if ($item->filename === $file) {
				$uploadId = $item->id;
				break;
			}
		}

		$command = "storage-users uploads sessions --id=$uploadId --$flag --json";
		$body = [
			"command" => $command,
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When /^the administrator (resumes|restarts) the upload session of file "([^"]*)" using postprocessing command$/
	 *
	 * @param string $flag
	 * @param string $file
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function theAdministratorResumesOrRestartsUploadSessionOfFileUsingPostprocessingCommand(
		string $flag,
		string $file,
	): void {
		$uploadSessions = CliHelper::runCommand(["command" => "storage-users uploads sessions --json"]);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, '', $uploadSessions);
		$responseArray = $this->getJSONDecodedCliMessage($uploadSessions);

		foreach ($responseArray as $item) {
			if ($item->filename === $file) {
				$uploadId = $item->id;
				break;
			}
		}

		$command = "postprocessing resume" . ($flag === "restarts" ? " -r" : "") . " -u $uploadId";
		$body = [
			"command" => $command,
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
					"The resource '$expectedFile' was found in the response.",
				);
			}
		} else {
			foreach ($expectedFiles as $expectedFile) {
				Assert::assertTrue(
					\in_array($expectedFile, $resourceNames),
					"The resource '$expectedFile' was not found in the response.",
				);
			}
		}
	}

	/**
	 * @Given the administrator waits until file :filename is no longer in processing upload sessions
	 *
	 * @param string $filename
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function theAdministratorWaitsUntilFileIsNoLongerInProcessingUploadSessions(string $filename): void {
		$timeout = 30; // seconds
		$interval = 1; // seconds
		$startTime = \time();

		// Poll until file is no longer in --processing list (allows time for state transitions).
		// Replaces fixed wait that assumed deterministic timing.
		while (true) {
			$this->theAdministratorListsAllTheUploadSessions("processing");
			$this->featureContext->theHTTPStatusCodeShouldBe(200);
			$responseArray = $this->getJSONDecodedCliMessage($this->featureContext->getResponse());

			$found = false;
			foreach ($responseArray as $item) {
				if (isset($item->filename) && $item->filename === $filename) {
					$found = true;
					break;
				}
			}

			if (!$found) {
				// File is no longer in processing list - abort completed and state updated
				return;
			}

			$elapsed = \time() - $startTime;
			if ($elapsed >= $timeout) {
				Assert::fail(
					"Timeout after {$timeout}s: file '{$filename}' is still in processing upload sessions list. " .
					"Virus scan + abort may not have completed within timeout period.",
				);
			}

			\sleep($interval);
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
			"command" => $command,
		];
		$response = CliHelper::runCommand($body);
		Assert::assertEquals("200", $response->getStatusCode(), "Failed to clean upload sessions");
	}

	/**
	 * @AfterScenario @cli-stale-uploads
	 *
	 * @return void
	 */
	public function cleanUpStaleUploads(): void {
		$response = $this->deleteStaleUploads();
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "Failed to cleanup stale upload", $response);
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
			"command" => $command,
		];

		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @Given the administrator has created stale upload
	 *
	 * @return void
	 */
	public function theAdministratorHasCreatedStaleUpload(): void {
		$folderPath = $this->featureContext->getStorageUsersRoot() . "/uploads";
		$infoFiles = glob($folderPath . '/*.info');
		foreach ($infoFiles as $file) {
			if (!unlink($file)) {
				Assert::fail("Fail to delete info file");
			}
		}
	}

	/**
	 * @param string|null $spaceId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	protected function listStaleUploads(?string $spaceId = null): ResponseInterface {
		$command = "storage-users uploads delete-stale-nodes --dry-run=true";

		if ($spaceId !== null) {
			$command .= " --spaceid=$spaceId";
		}

		$body = [
			"command" => $command,
		];
		return CliHelper::runCommand($body);
	}

	/**
	 * @param string|null $spaceId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	protected function deleteStaleUploads(?string $spaceId = null): ResponseInterface {
		$command = "storage-users uploads delete-stale-nodes --dry-run=false";
		if ($spaceId !== null) {
			$command .= " --spaceid=$spaceId";
		}

		$body = [
			"command" => $command,
		];
		return CliHelper::runCommand($body);
	}

	/**
	 * @When the administrator lists all the stale uploads
	 *
	 * @return void
	 */
	public function theAdministratorListsAllTheStaleUploads(): void {
		$this->featureContext->setResponse($this->listStaleUploads());
	}

	/**
	 * @When the administrator lists all the stale uploads of space :space owned by user :user
	 *
	 * @param string $spaceName
	 * @param string $user
	 *
	 * @return void
	 */
	public function theAdministratorListsTheStaleUploadsOfSpace(
		string $spaceName,
		string $user,
	): void {
		$spaceOwnerId = $this->spacesContext->getSpaceOwnerUserIdByName(
			$user,
			$spaceName,
		);
		$this->featureContext->setResponse($this->listStaleUploads($spaceOwnerId));
	}

	/**
	 * @Then the CLI response should contain the following message:
	 *
	 * @param PyStringNode $content
	 *
	 * @return void
	 */
	public function theCLIResponseShouldContainTheseMessage(PyStringNode $content): void {
		$response = $this->featureContext->getJsonDecodedResponseBodyContent();
		$expectedMessage = str_replace("\r\n", "\n", trim($content->getRaw()));
		$actualMessage = str_replace("\r\n", "\n", trim($response->message ?? ''));

		Assert::assertSame(
			$expectedMessage,
			$actualMessage,
			"Expected cli output to be $expectedMessage but found $actualMessage",
		);
	}

	/**
	 * @When the administrator deletes all the stale uploads
	 *
	 * @return void
	 */
	public function theAdministratorDeletesAllTheStaleUploads(): void {
		$this->featureContext->setResponse($this->deleteStaleUploads());
	}

	/**
	 * @When the administrator deletes all the stale uploads of space :spaceName owned by user :user
	 *
	 * @param string $spaceName
	 * @param string $user
	 *
	 * @return void
	 */
	public function theAdministratorDeletesTheStaleUploadsOfSpaceOwnedByUser(
		string $spaceName,
		string $user,
	): void {
		$spaceOwnerId = $this->spacesContext->getSpaceOwnerUserIdByName(
			$user,
			$spaceName,
		);
		$this->featureContext->setResponse($this->deleteStaleUploads($spaceOwnerId));
	}

	/**
	 * @Then there should be :number stale uploads
	 * @Then there should be :number stale uploads of space :spaceName owned by user :user
	 *
	 * @param int $number
	 * @param string $spaceName
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function thereShouldBeStaleUploadsOfSpaceOwnedByUser(
		int $number,
		string $spaceName = '',
		string $user = '',
	): void {
		$spaceOwnerId = null;
		if ($spaceName !== '' && $user !== '') {
			$spaceOwnerId = $this->spacesContext->getSpaceOwnerUserIdByName(
				$user,
				$spaceName,
			);
		}

		$response = $this->listStaleUploads($spaceOwnerId);
		$jsonDecodedResponse = $this->featureContext->getJsonDecodedResponseBodyContent($response);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "", $response);

		$expectedMessage = "Total stale nodes: $number";

		Assert::assertStringContainsString(
			$expectedMessage,
			$jsonDecodedResponse->message ?? '',
			"Expected message to contain '$expectedMessage', but got: " . ($jsonDecodedResponse->message ?? 'null'),
		);

	}

	/**
	 * @When the administrator lists all the trashed resources of space :space owned by user :user
	 *
	 * @param string $spaceName
	 * @param string $user
	 *
	 * @return void
	 */
	public function theAdministratorListsAllTrashedResourceOfSpaceOwnedByUser(
		string $spaceName,
		string $user,
	): void {
		$spaceOwnerId = $this->spacesContext->getSpaceOwnerUserIdByName(
			$user,
			$spaceName,
		);
		$this->featureContext->setResponse($this->listTrashedResource($spaceOwnerId));
	}

	/**
	 * @param string $spaceId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	protected function listTrashedResource(string $spaceId): ResponseInterface {
		$body = [
			"command" => "storage-users trash-bin list $spaceId",
		];
		return CliHelper::runCommand($body);
	}

	/**
	 * @param ResponseInterface|null $response
	 *
	 * @return array
	 */
	protected function getTrashedResourceFromCliCommandResponse(
		?ResponseInterface $response = null,
	): array {
		$responseArray = $this->featureContext->getJsonDecodedResponseBodyContent($response);
		$lines = explode("\n", $responseArray->message);
		$items = [];
		$totalCount = 0;
		foreach ($lines as $line) {
			if (preg_match(
				'/^\s*\│\s*([a-f0-9\-]{36})\s*\│\s*(.*?)\s*\│\s*(file|folder)\s*\│\s*([\d\-T:Z]+)\s*\│/',
				$line,
				$matches,
			)
			) {
				$items[] = [
					'itemID' => $matches[1],
					'path' => $matches[2],
					'type' => $matches[3],
					'delete at' => $matches[4],
				];
			}

			if (preg_match('/total count:\s*(\d+)/', $line, $countMatch)) {
				$totalCount = (int)$countMatch[1];
			}
		}
		return [$items, $totalCount];
	}

	/**
	 * @Then /^the command output should contain "([^"]*)" trashed resources with the following information:$/
	 *
	 * @param int $count
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function theCommandOutputShouldContainTheFollowingTrashResource(
		int $count,
		TableNode $table,
	): void {
		[$items, $totalCount] = $this->getTrashedResourceFromCliCommandResponse();

		Assert::assertSame($totalCount, $count, "Expected total trashed resource");

		foreach ($table->getHash() as $expectedRow) {
			$matchFound = false;

			foreach ($items as $item) {
				if ($item['path'] === $expectedRow['resource']
					&& $item['type'] === $expectedRow['type']
				) {
					$matchFound = true;
					break;
				}
			}

			Assert::assertTrue(
				$matchFound,
				"Could not find expected resource '{$expectedRow['resource']}' of type '{$expectedRow['type']}'",
			);
		}
	}

	/**
	 * @When the administrator restores all the trashed resources of space :space owned by user :user
	 *
	 * @param string $spaceName
	 * @param string $user
	 *
	 * @return void
	 */
	public function theAdministratorRestoresAllTheTrashedResourcesOfSpaceOwnedByUser(
		string $spaceName,
		string $user,
	): void {
		$spaceOwnerId = $this->spacesContext->getSpaceOwnerUserIdByName(
			$user,
			$spaceName,
		);
		$body = [
			"command" => "storage-users trash-bin restore-all -y $spaceOwnerId",
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @Then there should be no trashed resources of space :spaceName owned by user :user
	 *
	 * @param string $spaceName
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function thereShouldBeNoTrashedResourcesOfSpaceOwnedByUser(
		string $spaceName,
		string $user,
	): void {
		$spaceOwnerId = $this->spacesContext->getSpaceOwnerUserIdByName(
			$user,
			$spaceName,
		);
		$response = $this->listTrashedResource($spaceOwnerId);
		[$items, $totalCount] = $this->getTrashedResourceFromCliCommandResponse($response);

		Assert::assertSame($totalCount, 0, "Expected total trashed resource");
	}

	/**
	 * @Then there should be :number trashed resources of space :spaceName owned by user :user:
	 *
	 * @param int $number
	 * @param string $spaceName
	 * @param string $user
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function thereShouldBeTrashedResourcesOfSpaceOwnedByUser(
		int $number,
		string $spaceName,
		string $user,
		TableNode $table,
	): void {
		$spaceOwnerId = $this->spacesContext->getSpaceOwnerUserIdByName(
			$user,
			$spaceName,
		);
		$response = $this->listTrashedResource($spaceOwnerId);
		[$items, $totalCount] = $this->getTrashedResourceFromCliCommandResponse($response);

		Assert::assertSame($totalCount, $number, "Expected total trashed resource");

		foreach ($table->getHash() as $expectedRow) {
			$matchFound = false;

			foreach ($items as $item) {
				if ($item['path'] === $expectedRow['resource']
					&& $item['type'] === $expectedRow['type']
				) {
					$matchFound = true;
					break;
				}
			}

			Assert::assertTrue(
				$matchFound,
				"Could not find expected resource '{$expectedRow['resource']}' of type '{$expectedRow['type']}'",
			);
		}
	}

	/**
	 * @When the administrator restores the trashed resources :resource of space :space owned by user :user
	 *
	 * @param string $resource
	 * @param string $spaceName
	 * @param string $user
	 *
	 * @return void
	 */
	public function theAdministratorRestoresTheTrashedResourcesOfSpaceOwnedByUser(
		string $resource,
		string $spaceName,
		string $user,
	): void {
		$spaceOwnerId = $this->spacesContext->getSpaceOwnerUserIdByName(
			$user,
			$spaceName,
		);
		$response = $this->listTrashedResource($spaceOwnerId);
		[$items, $totalCount] = $this->getTrashedResourceFromCliCommandResponse($response);
		$trashItemId = null;
		foreach ($items as $item) {
			if ($item['path'] === $resource) {
				$trashItemId = $item['itemID'];
				break;
			}
		}
		$body = [
			"command" => "storage-users trash-bin restore $spaceOwnerId $trashItemId",
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator purges the expired trash resources
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theAdministratorPurgesTheExpiredTrashResources(): void {
		$body = [
			"command" => "storage-users trash-bin purge-expired",
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator resumes all the upload sessions using the CLI
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theAdministratorResumesAllTheUploadSessionsUsingTheCLI(): void {
		$command = "storage-users uploads sessions --resume --json";
		$body = [
			"command" => $command,
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When the administrator restarts the expired upload sessions using the CLI
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theAdministratorRestartsTheExpiredUploadSessionsUsingTheCLI(): void {
		$command = "storage-users uploads sessions --expired --restart --json";
		$body = [
			"command" => $command,
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When administrator deletes :space space using the CLI
	 *
	 * @param string $space
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theAdministratorDeletesSpaceUsingTheCLI(string $space): void {
		$command = "storage-users spaces purge -t=" . strtolower($space) . " -r=0s -v -d=false";
		$body = [
			"command" => $command,
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When administrator deletes :space space of user :user with space-id using the CLI
	 *
	 * @param string $spaceName
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function administratorDeletesSpaceWithSpaceIdUsingTheCli(string $spaceName, string $user): void {
		$space = $this->featureContext->spacesContext->getSpaceByName($user, $spaceName);

		if ($spaceName === "Personal") {
			$spaceType = "personal";
		} else {
			$spaceType = 'project';
		}

		$command = "storage-users spaces purge -t=$spaceType -s=" . $space['id'] . " -r=0s -v -d=false";
		$body = [
			"command" => $command,
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}

	/**
	 * @When administrator deletes :space space with :retentionPeriod retention period using the CLI
	 *
	 * @param string $space
	 * @param string $retentionPeriod
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function administratorDeletesSpaceWithRetentionPeriodOfUsingTheCli(
		string $space,
		string $retentionPeriod,
	): void {
		$command = "storage-users spaces purge -t=" . strtolower($space) . " -r=$retentionPeriod -v -d=false";
		$body = [
			"command" => $command,
		];
		$this->featureContext->setResponse(CliHelper::runCommand($body));
	}
}
