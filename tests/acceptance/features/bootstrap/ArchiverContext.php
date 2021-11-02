<?php
/**
 * ownCloud
 *
 * @author Artur Neumann <artur@jankaritech.com>
 * @copyright Copyright (c) 2021 Artur Neumann artur@jankaritech.com
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
use Behat\Gherkin\Node\TableNode;
use TestHelpers\HttpRequestHelper;
use TestHelpers\SetupHelper;
use wapmorgan\UnifiedArchive\UnifiedArchive;
use PHPUnit\Framework\Assert;

require_once 'bootstrap.php';

/**
 * Context for Archiver specific steps
 */
class ArchiverContext implements Context {

	/**
	 * @var FeatureContext
	 */
	private $featureContext;

	/**
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 *
	 * @throws Exception
	 */
	public function setUpScenario(BeforeScenarioScope $scope): void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context
		$this->featureContext = $environment->getContext('FeatureContext');
		SetupHelper::init(
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getOcPath()
		);
	}

	/**
	 * @When user :user downloads the archive of :resourceId using the resource id
	 *
	 * @param string $user
	 * @param string $resource
	 *
	 * @return void
	 *
	 * @throws \GuzzleHttp\Exception\GuzzleException
	 */
	public function userDownloadsTheArchiveOfUsingTheResourceId(string $user, string $resource): void {
		$resourceId = $this->featureContext->getFileIdForPath($user, $resource);
		$user = $this->featureContext->getActualUsername($user);
		$this->featureContext->setResponse(
			HttpRequestHelper::get(
				$this->featureContext->getBaseUrl() . '/archiver?id=' . $resourceId,
				'',
				$user,
				$this->featureContext->getPasswordForUser($user)
			)
		);
	}

	/**
	 * @Then the downloaded archive should contain these files:
	 *
	 * @param TableNode $expectedFiles
	 *
	 * @return void
	 *
	 * @throws Exception
	 */
	public function theDownloadedArchiveShouldContainTheseFiles(TableNode $expectedFiles) {
		$this->featureContext->verifyTableNodeColumns($expectedFiles, ['name', 'content']);
		$tempFile = \tempnam(\sys_get_temp_dir(), 'OcAcceptanceTests_');
		\unlink($tempFile); // we only need the name
		$tempFile = $tempFile . '.tar'; // it needs the extension
		\file_put_contents($tempFile, $this->featureContext->getResponse()->getBody()->getContents());
		$archive = UnifiedArchive::open($tempFile);
		foreach ($expectedFiles->getHash() as $expectedFile) {
			Assert::assertTrue(
				$archive->hasFile($expectedFile['name']),
				__METHOD__ .
				" archive does not contain '" . $expectedFile['name'] . "'"
			);
			Assert::assertEquals(
				$expectedFile['content'],
				$archive->getFileContent($expectedFile['name']),
				__METHOD__ .
				" content of '" . $expectedFile['name'] . "' not as expected"
			);
		}
		\unlink($tempFile);
	}
}
