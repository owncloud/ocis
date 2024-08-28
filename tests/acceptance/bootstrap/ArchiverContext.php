<?php declare(strict_types=1);
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
use GuzzleHttp\Exception\GuzzleException;
use TestHelpers\HttpRequestHelper;
use TestHelpers\SetupHelper;
use PHPUnit\Framework\Assert;
use Psr\Http\Message\ResponseInterface;
use splitbrain\PHPArchive\Tar;
use splitbrain\PHPArchive\Zip;
use splitbrain\PHPArchive\Archive;

require_once 'bootstrap.php';

/**
 * Context for Archiver specific steps
 */
class ArchiverContext implements Context {
	/**
	 * @var FeatureContext
	 */
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
	 * @param string $type
	 *
	 * @return Archive
	 */
	public function getArchiveClass(string $type): Archive {
		if ($type === 'zip') {
			return new Zip();
		} elseif ($type === 'tar') {
			return new Tar();
		} else {
			throw new Exception('Unknown archive type: ' . $type);
		}
	}

	/**
	 * @param string $dir
	 *
	 * @return void
	 */
	public function removeDir(string $dir): void {
		$items = array_diff(scandir($dir), ['.', '..']);
		foreach ($items as $item) {
			$itemPath = $dir . DIRECTORY_SEPARATOR . $item;
			if (\is_dir($itemPath)) {
				$this->removeDir($itemPath);
			} else {
				\unlink($itemPath);
			}
		}
		\rmdir($dir);
	}

	/**
	 * @param string $user
	 * @param string $resource
	 * @param string $addressType id|ids|path|paths
	 *
	 * @return string
	 *
	 * @throws Exception
	 */
	private function getArchiverQueryString(
		string $user,
		string $resource,
		string $addressType
	): string {
		switch ($addressType) {
			case 'id':
			case 'ids':
				return 'id=' . $this->featureContext->getFileIdForPath($user, $resource);
			case 'path':
			case 'paths':
				return 'path=' . $resource;
			default:
				throw new Exception(
					'"' . $addressType .
					'" is not a legal value for $addressType, must be id|ids|path|paths'
				);
		}
	}

	/**
	 * @When /^user "([^"]*)" downloads the (zip|tar) archive of "([^"]*)" using the resource (id|ids|path|paths) and setting these headers:$/
	 *
	 * @param string $user
	 * @param string $archiveType
	 * @param string $resource
	 * @param string $addressType id|path
	 * @param TableNode $headersTable
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function userDownloadsTheZipOrTarArchiveOfResourceUsingResourceIdOrPathAndSettingTheseHeaders(
		string $user,
		string $archiveType,
		string $resource,
		string $addressType,
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
		$this->featureContext->setResponse($this->downloadArchive($user, $resource, $addressType, $archiveType, null, $headers));
	}

	/**
	 * @When user :downloader downloads the archive of :item of user :owner using the resource :addressType
	 *
	 * @param string $downloader Who sends the request
	 * @param string $resource
	 * @param string $owner Who is the real owner of the file
	 * @param string $addressType id|path
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function userDownloadsTheArchiveOfItemOfUser(
		string $downloader,
		string $resource,
		string $owner,
		string $addressType
	): void {
		$this->featureContext->setResponse($this->downloadArchive($downloader, $resource, $addressType, null, $owner));
	}

	/**
	 * @param string $downloader
	 * @param string $resource
	 * @param string $addressType
	 * @param string|null $archiveType
	 * @param string|null $owner
	 * @param array|null $headers
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function downloadArchive(
		string $downloader,
		string $resource,
		string $addressType,
		?string $archiveType = null,
		?string $owner = null,
		?array $headers = null
	): ResponseInterface {
		$owner = $owner ??  $downloader;
		$downloader = $this->featureContext->getActualUsername($downloader);
		$queryString = $this->getArchiverQueryString($owner, $resource, $addressType);
		if ($archiveType !== null) {
			$queryString .= '&output-format=' . $archiveType;
		}
		return HttpRequestHelper::get(
			$this->featureContext->getBaseUrl() . '/archiver?' . $queryString,
			$this->featureContext->getStepLineRef(),
			$downloader,
			$this->featureContext->getPasswordForUser($downloader),
			$headers
		);
	}

	/**
	 * @When user :user downloads the archive of these items using the resource :addressType
	 *
	 * @param string $user
	 * @param TableNode $items
	 * @param string $addressType ids|paths
	 *
	 * @return void
	 *
	 * @throws GuzzleException|Exception
	 */
	public function userDownloadsTheArchiveOfTheseItems(
		string $user,
		TableNode $items,
		string $addressType
	): void {
		$user = $this->featureContext->getActualUsername($user);
		$queryString = '';
		foreach ($items->getRows() as $item) {
			$queryString .= $this->getArchiverQueryString($user, $item[0], $addressType) . '&';
		}

		$queryString = \rtrim($queryString, '&');
		$this->featureContext->setResponse(
			HttpRequestHelper::get(
				$this->featureContext->getBaseUrl() . '/archiver?' . $queryString,
				$this->featureContext->getStepLineRef(),
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
	public function theDownloadedArchiveShouldContainTheseFiles(string $type, TableNode $expectedFiles):void {
		$this->featureContext->verifyTableNodeColumns($expectedFiles, ['name', 'content']);
		$contents = $this->featureContext->getResponse()->getBody()->getContents();
		$tempFile = \tempnam(\sys_get_temp_dir(), 'OcAcceptanceTests_');
		$tempExtractFolder = $tempFile;
		\unlink($tempFile); // we only need the name
		$tempFile = $tempFile . '.' . $type; // it needs the extension
		\file_put_contents($tempFile, $contents);

		// open the archive
		$tar = $this->getArchiveClass($type);
		$tar->open($tempFile);
		$archiveData = $tar->contents();
		
		// extract the archive
		$tar->open($tempFile);
		$tar->extract($tempExtractFolder);
		$tar->close();

		foreach ($expectedFiles->getHash() as $expectedItem) {
			$expectedPath = trim($expectedItem['name'], "/");
			$found = false;
			foreach ($archiveData as $info) {
				// get only the parent folder path for the given item
				$actualPath = $info->getPath();

				if ($expectedPath === $actualPath) {
					if (!$info->getIsdir()) {
						$fileContent = \file_get_contents("$tempExtractFolder/$actualPath");
						Assert::assertEquals(
							$expectedItem['content'],
							$fileContent,
							__METHOD__ .
							" content of '" . $expectedPath . "' not as expected"
						);
					}
					$found = true;
					break;
				}
			}
			if (!$found) {
				Assert::fail("Resource '" . $expectedPath . "' is not in the downloaded archive.");
			}
		}
		\unlink($tempFile);
		$this->removeDir($tempExtractFolder);
	}
}
