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
	 * @When user :user downloads the archive of :resourceId using the resource id and setting these headers
	 *
	 * @param string $user
	 * @param string $resource
	 * @param TableNode $headersTable
	 *
	 * @return void
	 *
	 * @throws \GuzzleHttp\Exception\GuzzleException
	 */
	public function userDownloadsTheArchiveOfUsingTheResourceId(
		string $user,
		string $resource,
		TableNode $headersTable
	): void {
		$this->featureContext->verifyTableNodeColumns(
			$headersTable,
			['header', 'value']
		);
		$headers = [];
		foreach ($headersTable as $row) {
			$headers[$row['header']] = $row ['value'];
		}
		$resourceId = $this->featureContext->getFileIdForPath($user, $resource);
		$user = $this->featureContext->getActualUsername($user);
		$this->featureContext->setResponse(
			HttpRequestHelper::get(
				$this->featureContext->getBaseUrl() . '/archiver?id=' . $resourceId,
				'',
				$user,
				$this->featureContext->getPasswordForUser($user),
				$headers
			)
		);
	}

	/**
	 * @When user :downloader downloads the archive of :item of user :owner using the resource id
	 *
	 * @param string $downloader Who sends the request
	 * @param string $resource
	 * @param string $owner Who is the real owner of the file
	 *
	 * @return void
	 *
	 * @throws \GuzzleHttp\Exception\GuzzleException
	 */
	public function userDownloadsTheArchiveOfItemOfUserUsingTheResourceId(
		string $downloader,
		string $resource,
		string $owner
	): void {
		$resourceId = $this->featureContext->getFileIdForPath($owner, $resource);
		$downloader = $this->featureContext->getActualUsername($downloader);
		$this->featureContext->setResponse(
			HttpRequestHelper::get(
				$this->featureContext->getBaseUrl() . '/archiver?id=' . $resourceId,
				'',
				$downloader,
				$this->featureContext->getPasswordForUser($downloader),
			)
		);
	}

	/**
	 * @When user :arg1 downloads the archive of these items using the resource ids
	 *
	 * @param string $user
	 * @param TableNode $items
	 *
	 * @return void
	 *
	 * @throws \GuzzleHttp\Exception\GuzzleException
	 */
	public function userDownloadsTheArchiveOfTheseItemsUsingTheResourceIds(
		string $user,
		TableNode $items
	): void {
		$user = $this->featureContext->getActualUsername($user);
		$resourceIdsString = '';
		foreach ($items->getRows() as $item) {
			$fileId = $this->featureContext->getFileIdForPath($user, $item[0]);
			$resourceIdsString .= 'id=' . $fileId . '&';
		}
		$resourceIdsString = \rtrim($resourceIdsString, '&');
		$this->featureContext->setResponse(
			HttpRequestHelper::get(
				$this->featureContext->getBaseUrl() . '/archiver?' . $resourceIdsString,
				'',
				$user,
				$this->featureContext->getPasswordForUser($user),
			)
		);
	}

	/**
	 * @Then the downloaded :type archive should contain these files:
	 *
	 * @param string $type
	 * @param TableNode $expectedFiles
	 *
	 * @return void
	 *
	 * @throws Exception
	 */
	public function theDownloadedArchiveShouldContainTheseFiles(string $type, TableNode $expectedFiles) {
		$this->featureContext->verifyTableNodeColumns($expectedFiles, ['name', 'content']);
		$tempFile = \tempnam(\sys_get_temp_dir(), 'OcAcceptanceTests_');
		\unlink($tempFile); // we only need the name
		$tempFile = $tempFile . '.' . $type; // it needs the extension
		\file_put_contents($tempFile, $this->featureContext->getResponse()->getBody()->getContents());
		$archive = UnifiedArchive::open($tempFile);
		foreach ($expectedFiles->getHash() as $expectedFile) {
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
