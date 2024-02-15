<?php declare(strict_types=1);
/**
 * @author Christoph Wurst <christoph@owncloud.com>
 *
 * @copyright Copyright (c) 2018, ownCloud GmbH
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

use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use TestHelpers\HttpRequestHelper;
use Behat\Gherkin\Node\TableNode;
use Behat\Behat\Context\Context;
use TestHelpers\SetupHelper;
use \Psr\Http\Message\ResponseInterface;
use TestHelpers\WebDavHelper;

/**
 * Authentication functions
 */
class AuthContext implements Context {
	private string $clientToken;
	private string $appToken;
	private array $appTokens;
	private FeatureContext $featureContext;

	/**
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 */
	public function setUpScenario(BeforeScenarioScope $scope):void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context
		$this->featureContext = $environment->getContext('FeatureContext');

		// Reset ResponseXml
		$this->featureContext->setResponseXml([]);

		// Initialize SetupHelper class
		SetupHelper::init(
			$this->featureContext->getAdminUsername(),
			$this->featureContext->getAdminPassword(),
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getOcPath()
		);
	}

	/**
	 * @param string $user
	 * @param string $password
	 *
	 * @return array
	 */
	public function createBasicAuthHeader(string $user, string $password): array {
		$header = [];
		$authString = \base64_encode("$user:$password");
		$header["Authorization"] = "basic $authString";
		return $header;
	}

	/**
	 * @param string $url
	 * @param string $method
	 * @param string|null $body
	 * @param array|null $headers
	 *
	 * @return ResponseInterface
	 */
	public function sendRequest(
		string $url,
		string $method,
		?string $body = null,
		?array $headers = []
	): ResponseInterface {
		$fullUrl = $this->featureContext->getBaseUrl() . $url;

		return HttpRequestHelper::sendRequest(
			$fullUrl,
			$this->featureContext->getStepLineRef(),
			$method,
			null,
			null,
			$headers,
			$body,
		);
	}

	/**
	 * @param string $user
	 * @param string $url
	 * @param string $method
	 * @param string|null $body
	 * @param array|null $headers
	 * @param string|null $property
	 *
	 * @return ResponseInterface
	 */
	public function requestUrlWithBasicAuth(
		string $user,
		string $url,
		string $method,
		?string $body = null,
		?array $headers = null,
		?string $property = null
	): ResponseInterface {
		$user = $this->featureContext->getActualUsername($user);
		$url = $this->featureContext->substituteInLineCodes(
			$url,
			$user
		);
		$authHeader = $this->createBasicAuthHeader($user, $this->featureContext->getPasswordForUser($user));
		$headers = \array_merge($headers ?? [], $authHeader);

		if ($property !== null) {
			$body = $this->featureContext->getBodyForOCSRequest($method, $property);
		}

		return $this->sendRequest(
			$url,
			$method,
			$body,
			$headers
		);
	}

	/**
	 * @When a user requests :url with :method and no authentication
	 *
	 * @param string $url
	 * @param string $method
	 *
	 * @return void
	 */
	public function userRequestsURLWithNoAuth(string $url, string $method):void {
		$this->featureContext->setResponse($this->sendRequest($url, $method));
	}

	/**
	 * @When a user requests these endpoints with :method with body :body and no authentication about user :user
	 *
	 * @param string $method
	 * @param string $body
	 * @param string $ofUser
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function userRequestsEndpointsWithBodyAndNoAuthThenStatusCodeAboutUser(string $method, string $body, string $ofUser, TableNode $table): void {
		$ofUser = \strtolower($this->featureContext->getActualUsername($ofUser));
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);
		foreach ($table->getHash() as $row) {
			$row['endpoint'] = $this->featureContext->substituteInLineCodes(
				$row['endpoint'],
				$ofUser
			);
			$response = $this->sendRequest($row['endpoint'], $method, $body);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When a user requests these endpoints with :method with no authentication about user :user
	 *
	 * @param string $method
	 * @param string $ofUser
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsEndpointsWithoutBodyAndNoAuthAboutUser(string $method, string $ofUser, TableNode $table): void {
		$ofUser = \strtolower($this->featureContext->getActualUsername($ofUser));
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);
		foreach ($table->getHash() as $row) {
			$row['endpoint'] = $this->featureContext->substituteInLineCodes(
				$row['endpoint'],
				$ofUser
			);
			$response = $this->sendRequest($row['endpoint'], $method);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When a user requests these endpoints with :method and no authentication
	 *
	 * @param string $method
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsEndpointsWithNoAuthentication(string $method, TableNode $table):void {
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);
		foreach ($table->getHash() as $row) {
			$this->featureContext->setResponse($this->sendRequest($row['endpoint'], $method));
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When the user :user requests these endpoints with :method with basic auth
	 *
	 * @param string $user
	 * @param string $method
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsEndpointsWithBasicAuth(string $user, string $method, TableNode $table):void {
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);
		foreach ($table->getHash() as $row) {
			$response = $this->requestUrlWithBasicAuth($user, $row['endpoint'], $method);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When /^user "([^"]*)" requests these endpoints with "([^"]*)" to (?:get|set) property "([^"]*)" about user "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $method
	 * @param string $property
	 * @param string $ofUser
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theUserRequestsTheseEndpointsToGetOrSetPropertyAboutUser(
		string $user,
		string $method,
		string $property,
		string $ofUser,
		TableNode $table
	):void {
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);

		foreach ($table->getHash() as $row) {
			$row['endpoint'] = $this->featureContext->substituteInLineCodes(
				$row['endpoint'],
				$ofUser
			);
			$response = $this->requestUrlWithBasicAuth($user, $row['endpoint'], $method, null, null, $property);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When the administrator requests these endpoints with :method
	 *
	 * @param string $method
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theAdminRequestsTheseEndpointsWithMethod(string $method, TableNode $table):void {
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);
		foreach ($table->getHash() as $row) {
			$response = $this->requestUrlWithBasicAuth(
				$this->featureContext->getAdminUsername(),
				$row['endpoint'],
				$method
			);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When user :user requests :url with :method using basic auth
	 *
	 * @param string $user
	 * @param string $url
	 * @param string $method
	 *
	 * @return void
	 */
	public function userRequestsURLUsingBasicAuth(
		string $user,
		string $url,
		string $method
	):void {
		$response = $this->requestUrlWithBasicAuth($user, $url, $method);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When user :user requests :url with :method using basic auth and with headers
	 *
	 * @param string $user
	 * @param string $url
	 * @param string $method
	 * @param TableNode $headersTable
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsURLWithUsingBasicAuthAndDepthHeader(string $user, string $url, string $method, TableNode $headersTable):void {
		$user = $this->featureContext->getActualUsername($user);
		$url = $this->featureContext->substituteInLineCodes(
			$url,
			$user
		);
		$this->featureContext->verifyTableNodeColumns(
			$headersTable,
			['header', 'value']
		);
		$headers = [];
		foreach ($headersTable as $row) {
			$headers[$row['header']] = $row ['value'];
		}
		$this->featureContext->setResponse(
			$this->requestUrlWithBasicAuth(
				$user,
				$url,
				$method,
				null,
				$headers
			)
		);
	}

	/**
	 * @When user :user requests these endpoints with :method using password :password
	 *
	 * @param string $user
	 * @param string $method
	 * @param string $password
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsTheseEndpointsWithPassword(
		string $user,
		string $method,
		string $password,
		TableNode $table
	): void {
		$user = $this->featureContext->getActualUsername($user);
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);

		$authHeader = $this->createBasicAuthHeader($user, $this->featureContext->getActualPassword($password));

		foreach ($table->getHash() as $row) {
			$response = $this->sendRequest($row['endpoint'], $method, null, $authHeader);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When user :user requests these endpoints with :method using password :password about user :ofUser
	 *
	 * @param string $user
	 * @param string $method
	 * @param string $password
	 * @param string $ofUser
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsTheseEndpointsUsingPasswordAboutUser(
		string $user,
		string $method,
		string $password,
		string $ofUser,
		TableNode $table
	):void {
		$user = $this->featureContext->getActualUsername($user);
		$ofUser = $this->featureContext->getActualUsername($ofUser);
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint'], ['destination']);

		$headers = $this->createBasicAuthHeader($user, $this->featureContext->getActualPassword($password));

		foreach ($table->getHash() as $row) {
			$row['endpoint'] = $this->featureContext->substituteInLineCodes(
				$row['endpoint'],
				$ofUser
			);
			if (isset($row['destination'])) {
				$headers['Destination'] = $this->featureContext->substituteInLineCodes(
					$this->featureContext->getBaseUrl() . $row['destination'],
					$ofUser
				);
			}
			$response = $this->sendRequest(
				$row['endpoint'],
				$method,
				null,
				$headers
			);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When user :user requests these endpoints with :method including body :body using password :password about user :ofUser
	 *
	 * @param string $user
	 * @param string $method
	 * @param string $body
	 * @param string $password
	 * @param string $ofUser
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsTheseEndpointsWithBodyUsingPasswordAboutUser(
		string $user,
		string $method,
		string $body,
		string $password,
		string $ofUser,
		TableNode $table
	):void {
		$user = $this->featureContext->getActualUsername($user);
		$ofUser = $this->featureContext->getActualUsername($ofUser);
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint'], ['destination']);

		$headers = $this->createBasicAuthHeader($user, $this->featureContext->getActualPassword($password));

		foreach ($table->getHash() as $row) {
			$row['endpoint'] = $this->featureContext->substituteInLineCodes(
				$row['endpoint'],
				$ofUser
			);
			if (isset($row['destination'])) {
				$headers['Destination'] = $this->featureContext->substituteInLineCodes(
					$this->featureContext->getBaseUrl() . $row['destination'],
					$ofUser
				);
			}
			$response = $this->sendRequest(
				$row['endpoint'],
				$method,
				$body,
				$headers
			);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When user :user requests these endpoints with :method including body :body about user :ofUser
	 *
	 * @param string $user
	 * @param string $method
	 * @param string $body
	 * @param string $ofUser
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsTheseEndpointsIncludingBodyAboutUser(string $user, string $method, string $body, string $ofUser, TableNode $table):void {
		$user = $this->featureContext->getActualUsername($user);
		$ofUser = $this->featureContext->getActualUsername($ofUser);
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);

		$headers = [];
		if ($method === 'MOVE' || $method === 'COPY') {
			$headers['Destination'] = '/path/to/destination';
		}

		foreach ($table->getHash() as $row) {
			$row['endpoint'] = $this->featureContext->substituteInLineCodes(
				$row['endpoint'],
				$ofUser
			);
			$response = $this->requestUrlWithBasicAuth(
				$user,
				$row['endpoint'],
				$method,
				$body,
				$headers
			);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When user :asUser requests these endpoints with :method using the password of user :ofUser
	 *
	 * @param string $asUser
	 * @param string $method
	 * @param string $ofUser
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsTheseEndpointsWithoutBodyUsingThePasswordOfUser(string $asUser, string $method, string $ofUser, TableNode $table):void {
		$asUser = $this->featureContext->getActualUsername($asUser);
		$ofUser = $this->featureContext->getActualUsername($ofUser);
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);

		// do request as $asUser using password of $ofUser
		$authHeader = $this->createBasicAuthHeader($asUser, $this->featureContext->getPasswordForUser($ofUser));

		foreach ($table->getHash() as $row) {
			$row['endpoint'] = $this->featureContext->substituteInLineCodes(
				$row['endpoint'],
				$ofUser
			);
			$response = $this->sendRequest(
				$row['endpoint'],
				$method,
				null,
				$authHeader
			);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When user :asUser requests these endpoints with :method including body :body using the password of user :user
	 *
	 * @param string $asUser
	 * @param string $method
	 * @param string $body
	 * @param string $ofUser
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsTheseEndpointsIncludingBodyUsingPasswordOfUser(string $asUser, string $method, ?string $body, string $ofUser, TableNode $table):void {
		$asUser = $this->featureContext->getActualUsername($asUser);
		$ofUser = $this->featureContext->getActualUsername($ofUser);
		$this->featureContext->verifyTableNodeColumns($table, ['endpoint']);

		// do request as $asUser using password of $ofUser
		$authHeader = $this->createBasicAuthHeader($asUser, $this->featureContext->getPasswordForUser($ofUser));

		foreach ($table->getHash() as $row) {
			$row['endpoint'] = $this->featureContext->substituteInLineCodes(
				$row['endpoint'],
				$ofUser
			);
			$response = $this->sendRequest(
				$row['endpoint'],
				$method,
				$body,
				$authHeader
			);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When user :user requests these endpoints with :method about user :ofUser
	 *
	 * @param string $user
	 * @param string $method
	 * @param string $ofUser
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userRequestsTheseEndpointsAboutUser(string $user, string $method, string $ofUser, TableNode $table):void {
		$headers = [];
		if ($method === 'MOVE' || $method === 'COPY') {
			$baseUrl = $this->featureContext->getBaseUrl();
			$suffix = "";
			if ($this->featureContext->getDavPathVersion() === WebDavHelper::DAV_VERSION_SPACES){
				$suffix = $this->featureContext->spacesContext->getSpaceIdByName($user, "Personal") . "/";
			}
			$davPath = WebDavHelper::getDavPath($user, $this->featureContext->getDavPathVersion());
			$headers['Destination'] = "{$baseUrl}{$davPath}{$suffix}moved";
		}

		foreach ($table->getHash() as $row) {
			$row['endpoint'] = $this->featureContext->substituteInLineCodes(
				$row['endpoint'],
				$ofUser
			);
			$response = $this->requestUrlWithBasicAuth(
				$user,
				$row['endpoint'],
				$method,
				null,
				$headers
			);
			$this->featureContext->setResponse($response);
			$this->featureContext->pushToLastStatusCodesArrays();
		}
	}

	/**
	 * @When user :user requests :endpoint with :method without retrying
	 *
	 * @param string $user
	 * @param string $endpoint
	 * @param string $method
	 *
	 * @return void
	 */
	public function userRequestsURLWithoutRetry(
		string $user,
		string $endpoint,
		string $method
	):void {
		$username = $this->featureContext->getActualUsername($user);
		$endpoint = $this->featureContext->substituteInLineCodes(
			$endpoint,
			$username
		);
		$fullUrl = $this->featureContext->getBaseUrl() . $endpoint;
		$response = HttpRequestHelper::sendRequestOnce(
			$fullUrl,
			$this->featureContext->getStepLineRef(),
			$method,
			$username,
			$this->featureContext->getPasswordForUser($user)
		);
		$this->featureContext->setResponse($response);
	}
}
