<?php declare(strict_types=1);

/**
 * @author Artur Neumann <artur@jankaritech.com>
 *
 * @copyright Copyright (c) 2020, ownCloud GmbH
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
use Behat\Gherkin\Node\TableNode;
use GuzzleHttp\Exception\GuzzleException;
use TestHelpers\HttpRequestHelper;
use TestHelpers\WebDavHelper;
use TusPhp\Exception\ConnectionException;
use TusPhp\Exception\TusException;
use TusPhp\Tus\Client;
use PHPUnit\Framework\Assert;
use Psr\Http\Message\ResponseInterface;

require_once 'bootstrap.php';

/**
 * TUS related test steps
 */
class TUSContext implements Context {
	private FeatureContext $featureContext;

	private ?string $resourceLocation = null;

	/**
	 * @return string
	 */
	public function getTusResourceLocation(): string {
		return $this->resourceLocation ?: "";
	}

	/**
	 * @param string $user
	 * @param TableNode $headers
	 * @param string $content
	 *
	 * @return ResponseInterface
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function createNewTUSResourceWithHeaders(string $user, TableNode $headers, string $content = ''): ResponseInterface {
		$this->featureContext->verifyTableNodeColumnsCount($headers, 2);
		$user = $this->featureContext->getActualUsername($user);
		$password = $this->featureContext->getUserPassword($user);
		$this->resourceLocation = null;

		$response = $this->featureContext->makeDavRequest(
			$user,
			"POST",
			null,
			$headers->getRowsHash(),
			$content,
			"files",
			null,
			false,
			$password
		);
		$locationHeader = $response->getHeader('Location');
		if (\sizeof($locationHeader) > 0) {
			$this->resourceLocation = $locationHeader[0];
		}
		return $response;
	}

	/**
	 * @When user :user creates a new TUS resource on the WebDAV API with these headers:
	 *
	 * @param string $user
	 * @param TableNode $headers
	 * @param string $content
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userCreateNewTUSResourceWithHeaders(string $user, TableNode $headers, string $content = ''): void {
		$response = $this->createNewTUSResourceWithHeaders($user, $headers, $content);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given user :user has created a new TUS resource on the WebDAV API with these headers:
	 *
	 * @param string $user
	 * @param TableNode $headers Tus-Resumable: 1.0.0 header is added automatically
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userHasCreatedNewTUSResourceWithHeaders(string $user, TableNode $headers): void {
		$response = $this->createNewTUSResource($user, $headers);
		$this->featureContext->theHTTPStatusCodeShouldBe(201, "", $response);
	}

	/**
	 * @param string $user
	 * @param TableNode $headers
	 *
	 * @return ResponseInterface
	 */
	public function createNewTUSResource(string $user, TableNode $headers):ResponseInterface {
		$rows = $headers->getRows();
		$rows[] = ['Tus-Resumable', '1.0.0'];
		return $this->createNewTUSResourceWithHeaders($user, new TableNode($rows));
	}

	/**
	 * @param string $user
	 * @param string $offset
	 * @param string $data
	 * @param string $checksum
	 * @param array|null $extraHeaders
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function sendsAChunkToTUSLocationWithOffsetAndData(string $user, string $offset, string $data, string $checksum = '', ?array $extraHeaders = null): ResponseInterface {
		$user = $this->featureContext->getActualUsername($user);
		$password = $this->featureContext->getUserPassword($user);
		$headers = [
		'Content-Type' => 'application/offset+octet-stream',
		'Tus-Resumable' => '1.0.0',
		'Upload-Checksum' => $checksum,
		'Upload-Offset' => $offset
		];
		$headers = empty($extraHeaders) ? $headers : array_merge($headers, $extraHeaders);

		return HttpRequestHelper::sendRequest(
			$this->resourceLocation,
			$this->featureContext->getStepLineRef(),
			'PATCH',
			$user,
			$password,
			$headers,
			$data
		);
	}

	/**
	 * @When user :user sends a chunk to the last created TUS Location with offset :offset and data :data using the WebDAV API
	 *
	 * @param string $user
	 * @param string $offset
	 * @param string $data
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userSendsAChunkToTUSLocationWithOffsetAndData(string $user, string $offset, string $data): void {
		$response = $this->sendsAChunkToTUSLocationWithOffsetAndData($user, $offset, $data);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user uploads file :source to :destination using the TUS protocol on the WebDAV API
	 *
	 * @param string|null $user
	 * @param string $source
	 * @param string $destination
	 * @param array $uploadMetadata array of metadata to be placed in the
	 *                              `Upload-Metadata` header.
	 *                              see https://tus.io/protocols/resumable-upload.html#upload-metadata
	 *                              Don't Base64 encode the value.
	 * @param int $noOfChunks
	 * @param int|null $bytes
	 * @param string $checksum
	 *
	 * @return void
	 * @throws ConnectionException
	 * @throws GuzzleException
	 * @throws JsonException
	 * @throws ReflectionException
	 * @throws TusException
	 */
	public function userUploadsUsingTusAFileTo(
		?string $user,
		string  $source,
		string  $destination,
		array   $uploadMetadata = [],
		int     $noOfChunks = 1,
		int     $bytes = null,
		string  $checksum = ''
	): void {
		$this->uploadFileUsingTus($user, $source, $destination, $uploadMetadata, $noOfChunks, $bytes, $checksum);
		$this->featureContext->setLastUploadDeleteTime(\time());
	}

	/**
	 * @param string $user
	 * @param string $source
	 * @param string $destination
	 * @param array $uploadMetadata
	 * @param integer $noOfChunks
	 * @param integer $bytes
	 * @param string $checksum
	 *
	 * @return void
	 */
	public function uploadFileUsingTus(
		?string $user,
		string  $source,
		string  $destination,
		array   $uploadMetadata = [],
		int     $noOfChunks = 1,
		int     $bytes = null,
		string  $checksum = ''
	) {
		$user = $this->featureContext->getActualUsername($user);
		$password = $this->featureContext->getUserPassword($user);
		$headers = [
			'Authorization' => 'Basic ' . \base64_encode($user . ':' . $password)
		];
		if ($bytes !== null) {
			$creationWithUploadHeader = [
				'Content-Type' => 'application/offset+octet-stream',
				'Tus-Resumable' => '1.0.0'
			];
			$headers = \array_merge($headers, $creationWithUploadHeader);
		}
		if ($checksum != '') {
			$checksumHeader = [
				'Upload-Checksum' => $checksum
			];
			$headers = \array_merge($headers, $checksumHeader);
		}

		$client = new Client(
			$this->featureContext->getBaseUrl(),
			[
				'verify' => false,
				'headers' => $headers
			]
		);
		$client->setChecksumAlgorithm('sha1');
		$client->setApiPath(
			WebDavHelper::getDavPath(
				$user,
				$this->featureContext->getDavPathVersion(),
				"files",
				WebDavHelper::$SPACE_ID_FROM_OCIS
					?: $this->featureContext->getPersonalSpaceIdForUser($user)
			)
		);
		WebDavHelper::$SPACE_ID_FROM_OCIS = '';
		$client->setMetadata($uploadMetadata);
		$sourceFile = $this->featureContext->acceptanceTestsDirLocation() . $source;
		$client->setKey((string)rand())->file($sourceFile, $destination);
		$this->featureContext->pauseUploadDelete();

		if ($bytes !== null) {
			$client->file($sourceFile, $destination)->createWithUpload($client->getKey(), $bytes);
		} elseif (\filesize($sourceFile) === 0) {
			$client->file($sourceFile, $destination)->createWithUpload($client->getKey(), 0);
		} elseif ($noOfChunks === 1) {
			$client->file($sourceFile, $destination)->upload();
		} else {
			$bytesPerChunk = (int)\ceil(\filesize($sourceFile) / $noOfChunks);
			for ($i = 0; $i < $noOfChunks; $i++) {
				$client->upload($bytesPerChunk);
			}
		}
	}

	/**
	 * @When user :user uploads file with content :content to :destination using the TUS protocol on the WebDAV API
	 *
	 * @param string $user
	 * @param string $content
	 * @param string $destination
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function userUploadsAFileWithContentToUsingTus(
		string $user,
		string $content,
		string $destination
	): void {
		$temporaryFileName = $this->writeDataToTempFile($content);
		try {
			$this->uploadFileUsingTus(
				$user,
				\basename($temporaryFileName),
				$destination
			);
			$this->featureContext->setLastUploadDeleteTime(\time());
		} catch (Exception $e) {
			Assert::assertStringContainsString('TusPhp\Exception\FileException: Unable to create resource', (string)$e);
		}
		\unlink($temporaryFileName);
	}

	/**
	 * @When user :user uploads file with content :content in :noOfChunks chunks to :destination using the TUS protocol on the WebDAV API
	 *
	 * @param string|null $user
	 * @param string $content
	 * @param int|null $noOfChunks
	 * @param string $destination
	 *
	 * @return void
	 * @throws ConnectionException
	 * @throws GuzzleException
	 * @throws JsonException
	 * @throws ReflectionException
	 * @throws TusException
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userUploadsAFileWithContentInChunksUsingTus(
		?string $user,
		string  $content,
		?int    $noOfChunks,
		string  $destination
	): void {
		$temporaryFileName = $this->writeDataToTempFile($content);
		$this->uploadFileUsingTus(
			$user,
			\basename($temporaryFileName),
			$destination,
			[],
			$noOfChunks
		);
		$this->featureContext->setLastUploadDeleteTime(\time());
		\unlink($temporaryFileName);
	}

	/**
	 * @When user :user uploads file :source to :destination with mtime :mtime using the TUS protocol on the WebDAV API
	 *
	 * @param string $user
	 * @param string $source
	 * @param string $destination
	 * @param string $mtime Time in human-readable format is taken as input which is converted into milliseconds that is used by API
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userUploadsFileWithContentToWithMtimeUsingTUS(
		string $user,
		string $source,
		string $destination,
		string $mtime
	): void {
		$mtime = new DateTime($mtime);
		$mtime = $mtime->format('U');
		$user = $this->featureContext->getActualUsername($user);
		$this->uploadFileUsingTus(
			$user,
			$source,
			$destination,
			['mtime' => $mtime]
		);
		$this->featureContext->setLastUploadDeleteTime(\time());
	}

	/**
	 * @param string $content
	 *
	 * @return string the file name
	 * @throws Exception
	 */
	public function writeDataToTempFile(string $content): string {
		$temporaryFileName = \tempnam(
			$this->featureContext->acceptanceTestsDirLocation(),
			"tus-upload-test-"
		);
		if ($temporaryFileName === false) {
			throw new \Exception("could not create a temporary filename");
		}
		$temporaryFile = \fopen($temporaryFileName, "w");
		if ($temporaryFile === false) {
			throw new \Exception("could not open " . $temporaryFileName . " for write");
		}
		\fwrite($temporaryFile, $content);
		\fclose($temporaryFile);
		return $temporaryFileName;
	}

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
	 * @When user :user creates a new TUS resource with content :content on the WebDAV API with these headers:
	 *
	 * @param string $user
	 * @param string $content
	 * @param TableNode $headers
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userCreatesWithUpload(
		string    $user,
		string    $content,
		TableNode $headers
	): void {
		$response = $this->createNewTUSResourceWithHeaders($user, $headers, $content);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user creates file :source and uploads content :content in the same request using the TUS protocol on the WebDAV API
	 *
	 * @param string $user
	 * @param string $source
	 * @param string $content
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userUploadsWithCreatesWithUpload(
		string $user,
		string $source,
		string $content
	): void {
		$temporaryFileName = $this->writeDataToTempFile($content);
		$this->uploadFileUsingTus(
			$user,
			\basename($temporaryFileName),
			$source,
			[],
			1,
			-1
		);
		$this->featureContext->setLastUploadDeleteTime(\time());
		\unlink($temporaryFileName);
	}

	/**
	 * @When user :user uploads file with checksum :checksum to the last created TUS Location with offset :offset and content :content using the TUS protocol on the WebDAV API
	 *
	 * @param string $user
	 * @param string $checksum
	 * @param string $offset
	 * @param string $content
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userUploadsFileWithChecksum(
		string $user,
		string $checksum,
		string $offset,
		string $content
	): void {
		$response = $this->sendsAChunkToTUSLocationWithOffsetAndData($user, $offset, $content, $checksum);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given user :user has uploaded file with checksum :checksum to the last created TUS Location with offset :offset and content :content using the TUS protocol on the WebDAV API
	 *
	 * @param string $user
	 * @param string $checksum
	 * @param string $offset
	 * @param string $content
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userHasUploadedFileWithChecksum(
		string $user,
		string $checksum,
		string $offset,
		string $content
	): void {
		$response = $this->sendsAChunkToTUSLocationWithOffsetAndData($user, $offset, $content, $checksum);
		$this->featureContext->theHTTPStatusCodeShouldBe(204, "", $response);
	}

	/**
	 * @When user :user sends a chunk to the last created TUS Location with offset :offset and data :data with checksum :checksum using the TUS protocol on the WebDAV API
	 *
	 * @param string $user
	 * @param string $offset
	 * @param string $data
	 * @param string $checksum
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userUploadsChunkFileWithChecksum(string $user, string $offset, string $data, string $checksum): void {
		$response = $this->sendsAChunkToTUSLocationWithOffsetAndData($user, $offset, $data, $checksum);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given user :user has uploaded a chunk to the last created TUS Location with offset :offset and data :data with checksum :checksum using the TUS protocol on the WebDAV API
	 *
	 * @param string $user
	 * @param string $offset
	 * @param string $data
	 * @param string $checksum
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userHasUploadedChunkFileWithChecksum(string $user, string $offset, string $data, string $checksum): void {
		$response = $this->sendsAChunkToTUSLocationWithOffsetAndData($user, $offset, $data, $checksum);
		$this->featureContext->theHTTPStatusCodeShouldBe(204, "", $response);
	}

	/**
	 * @When user :user overwrites recently shared file with offset :offset and data :data with checksum :checksum using the TUS protocol on the WebDAV API with these headers:
	 * @When user :user overwrites existing file with offset :offset and data :data with checksum :checksum using the TUS protocol on the WebDAV API with these headers:
	 *
	 * @param string $user
	 * @param string $offset
	 * @param string $data
	 * @param string $checksum
	 * @param TableNode $headers Tus-Resumable: 1.0.0 header is added automatically
	 *
	 * @return void
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public function userOverwritesFileWithChecksum(string $user, string $offset, string $data, string $checksum, TableNode $headers): void {
		$createResponse = $this->createNewTUSResource($user, $headers);
		$this->featureContext->theHTTPStatusCodeShouldBe(201, "", $createResponse);
		$response = $this->sendsAChunkToTUSLocationWithOffsetAndData($user, $offset, $data, $checksum);
		$this->featureContext->setResponse($response);
	}
}
