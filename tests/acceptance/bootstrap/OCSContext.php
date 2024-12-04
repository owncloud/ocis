<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Artur Neumann <artur@jankaritech.com>
 * @copyright Copyright (c) 2019, ownCloud GmbH
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
use Behat\Gherkin\Node\PyStringNode;
use Behat\Gherkin\Node\TableNode;
use Psr\Http\Message\ResponseInterface;
use PHPUnit\Framework\Assert;
use TestHelpers\HttpRequestHelper;
use TestHelpers\OcsApiHelper;
use TestHelpers\TranslationHelper;
use TestHelpers\BehatHelper;

require_once 'bootstrap.php';

/**
 * steps needed to send requests to the OCS API
 */
class OCSContext implements Context {
	private FeatureContext $featureContext;

	/**
	 * @When /^the user sends HTTP method "([^"]*)" to OCS API endpoint "([^"]*)"$/
	 *
	 * @param string $verb
	 * @param string $url
	 *
	 * @return void
	 */
	public function theUserSendsToOcsApiEndpoint(string $verb, string $url):void {
		$response = $this->theUserSendsToOcsApiEndpointWithBody($verb, $url);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^user "([^"]*)" sends HTTP method "([^"]*)" to OCS API endpoint "([^"]*)"$/
	 * @When /^user "([^"]*)" sends HTTP method "([^"]*)" to OCS API endpoint "([^"]*)" using password "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $verb
	 * @param string $url
	 * @param string|null $password
	 *
	 * @return void
	 */
	public function userSendsToOcsApiEndpoint(string $user, string $verb, string $url, ?string $password = null):void {
		$response = $this->sendRequestToOcsEndpoint(
			$user,
			$verb,
			$url,
			null,
			$password
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @param string $user
	 * @param string $verb
	 * @param string $url
	 * @param TableNode|null $body
	 * @param string|null $password
	 * @param array|null $headers
	 *
	 * @return ResponseInterface
	 */
	public function sendRequestToOcsEndpoint(
		string $user,
		string $verb,
		string $url,
		?TableNode $body = null,
		?string $password = null,
		?array $headers = null
	):ResponseInterface {
		/**
		 * array of the data to be sent in the body.
		 * contains $body data converted to an array
		 */
		$bodyArray = [];
		if ($body instanceof TableNode) {
			$bodyArray = $body->getRowsHash();
		}

		if ($user !== 'UNAUTHORIZED_USER') {
			if ($password === null) {
				$password = $this->featureContext->getPasswordForUser($user);
			}
			$user = $this->featureContext->getActualUsername($user);
		} else {
			$user = null;
			$password = null;
		}
		return OcsApiHelper::sendRequest(
			$this->featureContext->getBaseUrl(),
			$user,
			$password,
			$verb,
			$url,
			$this->featureContext->getStepLineRef(),
			$bodyArray,
			$this->featureContext->getOcsApiVersion(),
			$headers
		);
	}

	/**
	 * @param string $verb
	 * @param string $url
	 * @param TableNode|null $body
	 *
	 * @return ResponseInterface
	 */
	public function adminSendsHttpMethodToOcsApiEndpointWithBody(
		string $verb,
		string $url,
		?TableNode $body
	):ResponseInterface {
		$admin = $this->featureContext->getAdminUsername();
		return $this->sendRequestToOcsEndpoint(
			$admin,
			$verb,
			$url,
			$body
		);
	}

	/**
	 * @param string $verb
	 * @param string $url
	 * @param TableNode|null $body
	 *
	 * @return ResponseInterface
	 */
	public function theUserSendsToOcsApiEndpointWithBody(string $verb, string $url, ?TableNode $body = null):ResponseInterface {
		return $this->sendRequestToOcsEndpoint(
			$this->featureContext->getCurrentUser(),
			$verb,
			$url,
			$body
		);
	}

	/**
	 * @When /^user "([^"]*)" sends HTTP method "([^"]*)" to OCS API endpoint "([^"]*)" with body$/
	 *
	 * @param string $user
	 * @param string $verb
	 * @param string $url
	 * @param TableNode|null $body
	 * @param string|null $password
	 *
	 * @return void
	 */
	public function userSendHTTPMethodToOcsApiEndpointWithBody(
		string $user,
		string $verb,
		string $url,
		?TableNode $body = null,
		?string $password = null
	):void {
		$response = $this->sendRequestToOcsEndpoint(
			$user,
			$verb,
			$url,
			$body,
			$password
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When the administrator sends HTTP method :verb to OCS API endpoint :url
	 * @When the administrator sends HTTP method :verb to OCS API endpoint :url using password :password
	 *
	 * @param string $verb
	 * @param string $url
	 * @param string|null $password
	 *
	 * @return void
	 */
	public function theAdministratorSendsHttpMethodToOcsApiEndpoint(
		string $verb,
		string $url,
		?string $password = null
	):void {
		$this->featureContext->setResponse(
			$this->sendRequestToOcsEndpoint(
				$this->featureContext->getAdminUsername(),
				$verb,
				$url,
				null,
				$password
			)
		);
	}

	/**
	 * @When /^user "([^"]*)" sends HTTP method "([^"]*)" to OCS API endpoint "([^"]*)" with headers$/
	 *
	 * @param string $user
	 * @param string $verb
	 * @param string $url
	 * @param TableNode $headersTable
	 *
	 * @return void
	 * @throws Exception
	 */
	public function userSendsToOcsApiEndpointWithHeaders(
		string $user,
		string $verb,
		string $url,
		TableNode $headersTable
	):void {
		$user = $this->featureContext->getActualUsername($user);
		$password = $this->featureContext->getPasswordForUser($user);
		$this->featureContext->setResponse(
			$this->sendRequestToOcsEndpoint(
				$user,
				$verb,
				$url,
				null,
				$password,
				$headersTable->getRowsHash()
			)
		);
	}

	/**
	 * @When /^the administrator sends HTTP method "([^"]*)" to OCS API endpoint "([^"]*)" with headers$/
	 *
	 * @param string $verb
	 * @param string $url
	 * @param TableNode $headersTable
	 *
	 * @return void
	 * @throws Exception
	 */
	public function administratorSendsToOcsApiEndpointWithHeaders(
		string $verb,
		string $url,
		TableNode $headersTable
	):void {
		$user = $this->featureContext->getAdminUsername();
		$password = $this->featureContext->getPasswordForUser($user);
		$this->featureContext->setResponse(
			$this->sendRequestToOcsEndpoint(
				$user,
				$verb,
				$url,
				null,
				$password,
				$headersTable->getRowsHash()
			)
		);
	}

	/**
	 * @When /^the administrator sends HTTP method "([^"]*)" to OCS API endpoint "([^"]*)" with headers using password "([^"]*)"$/
	 *
	 * @param string $verb
	 * @param string $url
	 * @param string $password
	 * @param TableNode $headersTable
	 *
	 * @return void
	 * @throws Exception
	 */
	public function administratorSendsToOcsApiEndpointWithHeadersAndPassword(
		string $verb,
		string $url,
		string $password,
		TableNode $headersTable
	):void {
		$this->featureContext->setResponse(
			$this->sendRequestToOcsEndpoint(
				$this->featureContext->getAdminUsername(),
				$verb,
				$url,
				null,
				$password,
				$headersTable->getRowsHash()
			)
		);
	}

	/**
	 * @When the administrator sends HTTP method :verb to OCS API endpoint :url with body
	 *
	 * @param string $verb
	 * @param string $url
	 * @param TableNode|null $body
	 *
	 * @return void
	 */
	public function theAdministratorSendsHttpMethodToOcsApiEndpointWithBody(
		string $verb,
		string $url,
		?TableNode $body
	):void {
		$response = $this->adminSendsHttpMethodToOcsApiEndpointWithBody(
			$verb,
			$url,
			$body
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^the user sends HTTP method "([^"]*)" to OCS API endpoint "([^"]*)" with body$/
	 *
	 * @param string $verb
	 * @param string $url
	 * @param TableNode $body
	 *
	 * @return void
	 */
	public function theUserSendsHTTPMethodToOcsApiEndpointWithBody(string $verb, string $url, TableNode $body):void {
		$response = $this->theUserSendsToOcsApiEndpointWithBody(
			$verb,
			$url,
			$body
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When the administrator sends HTTP method :verb to OCS API endpoint :url with body using password :password
	 *
	 * @param string $verb
	 * @param string $url
	 * @param string $password
	 * @param TableNode $body
	 *
	 * @return void
	 */
	public function theAdministratorSendsHttpMethodToOcsApiWithBodyAndPassword(
		string $verb,
		string $url,
		string $password,
		TableNode $body
	):void {
		$admin = $this->featureContext->getAdminUsername();
		$response = $this->sendRequestToOcsEndpoint(
			$admin,
			$verb,
			$url,
			$body,
			$password
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @When /^user "([^"]*)" sends HTTP method "([^"]*)" to OCS API endpoint "([^"]*)" with body using password "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $verb
	 * @param string $url
	 * @param string $password
	 * @param TableNode $body
	 *
	 * @return void
	 */
	public function userSendsHTTPMethodToOcsApiEndpointWithBodyAndPassword(
		string $user,
		string $verb,
		string $url,
		string $password,
		TableNode $body
	):void {
		$response = $this->sendRequestToOcsEndpoint(
			$user,
			$verb,
			$url,
			$body,
			$password
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Then /^the OCS status code should be "([^"]*)"$/
	 *
	 * @param string $statusCode
	 * @param string $message
	 * @param ResponseInterface|null $response
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theOCSStatusCodeShouldBe(string $statusCode, string $message = "", ?ResponseInterface $response = null):void {
		$statusCodes = explode(",", $statusCode);
		$response = $response ?? $this->featureContext->getResponse();
		$responseStatusCode = $this->getOCSResponseStatusCode(
			$response
		);
		if (\is_array($statusCodes)) {
			if ($message === "") {
				$message = "OCS status code is not any of the expected values " . \implode(",", $statusCodes) . " got " . $responseStatusCode;
			}
			Assert::assertContainsEquals(
				$responseStatusCode,
				$statusCodes,
				$message
			);
			$this->featureContext->emptyLastOCSStatusCodesArray();
		} else {
			if ($message === "") {
				$message = "OCS status code is not the expected value " . $statusCodes . " got " . $responseStatusCode;
			}

			Assert::assertEquals(
				$statusCodes,
				$responseStatusCode,
				$message
			);
		}
	}

	/**
	 * @Then /^the OCS status code should be "([^"]*)" or "([^"]*)"$/
	 *
	 * @param string $statusCode1
	 * @param string $statusCode2
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theOcsStatusCodeShouldBeOr(string $statusCode1, string $statusCode2):void {
		$statusCodes = [$statusCode1,$statusCode1];
		$response = $this->featureContext->getResponse();
		$responseStatusCode = $this->getOCSResponseStatusCode(
			$response
		);
		Assert::assertContainsEquals(
			$responseStatusCode,
			$statusCodes,
			"OCS status code is not any of the expected values " . \implode(",", $statusCodes) . " got " . $responseStatusCode
		);
		$this->featureContext->emptyLastOCSStatusCodesArray();
	}

	/**
	 * Check the text in an OCS status message
	 *
	 * @Then /^the OCS status message should be "([^"]*)"$/
	 * @Then /^the OCS status message should be "([^"]*)" in language "([^"]*)"$/
	 *
	 * @param string $statusMessage
	 * @param string|null $language
	 *
	 * @return void
	 */
	public function theOCSStatusMessageShouldBe(string $statusMessage, ?string $language = null):void {
		$language = TranslationHelper::getLanguage($language);
		$statusMessage = $this->getActualStatusMessage($statusMessage, $language);

		Assert::assertEquals(
			$statusMessage,
			$this->getOCSResponseStatusMessage(
				$this->featureContext->getResponse()
			),
			'Unexpected OCS status message :"' . $this->getOCSResponseStatusMessage(
				$this->featureContext->getResponse()
			) . '" in response'
		);
	}

	/**
	 * Check the text in an OCS status message
	 *
	 * @Then /^the OCS status message about user "([^"]*)" should be "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $statusMessage
	 *
	 * @return void
	 */
	public function theOCSStatusMessageAboutUserShouldBe(string $user, string $statusMessage):void {
		$user = \strtolower($this->featureContext->getActualUsername($user));
		$statusMessage = $this->featureContext->substituteInLineCodes(
			$statusMessage,
			$user
		);
		Assert::assertEquals(
			$statusMessage,
			$this->getOCSResponseStatusMessage(
				$this->featureContext->getResponse()
			),
			'Unexpected OCS status message :"' . $this->getOCSResponseStatusMessage(
				$this->featureContext->getResponse()
			) . '" in response'
		);
	}

	/**
	 * Check the text in an OCS status message.
	 * Use this step form if the expected text contains double quotes,
	 * single quotes and other content that theOCSStatusMessageShouldBe()
	 * cannot handle.
	 *
	 * After the step, write the expected text in PyString form like:
	 *
	 * """
	 * File "abc.txt" can't be shared due to reason "xyz"
	 * """
	 *
	 * @Then /^the OCS status message should be:$/
	 *
	 * @param PyStringNode $statusMessage
	 *
	 * @return void
	 */
	public function theOCSStatusMessageShouldBePyString(
		PyStringNode $statusMessage
	):void {
		Assert::assertEquals(
			$statusMessage->getRaw(),
			$this->getOCSResponseStatusMessage(
				$this->featureContext->getResponse()
			),
			'Unexpected OCS status message: "' . $this->getOCSResponseStatusMessage(
				$this->featureContext->getResponse()
			) . '" in response'
		);
	}

	/**
	 * Parses the xml answer to get ocs response which doesn't match with
	 * http one in v1 of the api.
	 *
	 * @param ResponseInterface $response
	 *
	 * @return string
	 * @throws Exception
	 */
	public function getOCSResponseStatusCode(ResponseInterface $response):string {
		try {
			$jsonResponse = $this->featureContext->getJsonDecodedResponseBodyContent($response);
		} catch (JsonException $e) {
			$jsonResponse = null;
		}

		if (\is_object($jsonResponse) && $jsonResponse->ocs->meta->statuscode) {
			return (string) $jsonResponse->ocs->meta->statuscode;
		}
		// go to xml response when json response is null (it means not formatted and get status code)
		$responseXmlObject = HttpRequestHelper::getResponseXml($response, __METHOD__);
		if (isset($responseXmlObject->meta[0], $responseXmlObject->meta[0]->statuscode)) {
			return (string) $responseXmlObject->meta[0]->statuscode;
		}
		Assert::fail("No OCS status code found in response");
	}

	/**
	 * Parses the xml answer to return data items from ocs response
	 *
	 * @param ResponseInterface $response
	 *
	 * @return SimpleXMLElement
	 * @throws Exception
	 */
	public function getOCSResponseData(ResponseInterface $response): SimpleXMLElement {
		$responseXmlObject = HttpRequestHelper::getResponseXml($response, __METHOD__);
		if (isset($responseXmlObject->data)) {
			return $responseXmlObject->data;
		}
		Assert::fail("No OCS data items found in response");
	}

	/**
	 * Parses the xml answer to get ocs response message which doesn't match with
	 * http one in v1 of the api.
	 *
	 * @param ResponseInterface $response
	 *
	 * @return string
	 */
	public function getOCSResponseStatusMessage(ResponseInterface $response):string {
		return (string) HttpRequestHelper::getResponseXml($response, __METHOD__)->meta[0]->message;
	}

	/**
	 * convert status message in the desired language
	 *
	 * @param string $statusMessage
	 * @param string|null $language
	 *
	 * @return string
	 */
	public function getActualStatusMessage(string $statusMessage, ?string $language = null):string {
		if ($language !== null) {
			$multiLingualMessage = \json_decode(
				\file_get_contents(__DIR__ . "/../fixtures/multiLanguageErrors.json"),
				true
			);

			if (isset($multiLingualMessage[$statusMessage][$language])) {
				$statusMessage = $multiLingualMessage[$statusMessage][$language];
			}
		}
		return $statusMessage;
	}

	/**
	 * check if the HTTP status code and the OCS status code indicate that the request was successful
	 * this function is aware of the currently used OCS version
	 *
	 * @param string|null $message
	 * @param ResponseInterface|null $response
	 *
	 * @return void
	 * @throws Exception
	 */
	public function assertOCSResponseIndicatesSuccess(?string $message = "", ?ResponseInterface $response = null):void {
		$response = $response ?? $this->featureContext->getResponse();
		$this->featureContext->theHTTPStatusCodeShouldBe('200', $message, $response);
		if ($this->featureContext->getOcsApiVersion() === 1) {
			$this->theOCSStatusCodeShouldBe('100', $message, $response);
		} else {
			$this->theOCSStatusCodeShouldBe('200', $message, $response);
		}
	}

	/**
	 * This will run before EVERY scenario.
	 * It will set the properties for this object.
	 *
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 */
	public function before(BeforeScenarioScope $scope):void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context
		$this->featureContext = BehatHelper::getContext($scope, $environment, 'FeatureContext');
	}
}
