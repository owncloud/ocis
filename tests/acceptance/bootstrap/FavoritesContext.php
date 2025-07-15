<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Artur Neumann <artur@jankaritech.com>
 * @copyright Copyright (c) 2018 Artur Neumann artur@jankaritech.com
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
use Psr\Http\Message\ResponseInterface;
use TestHelpers\WebDavHelper;
use TestHelpers\BehatHelper;

require_once 'bootstrap.php';

/**
 * context containing favorites related API steps
 */
class FavoritesContext implements Context {
	private FeatureContext $featureContext;
	private WebDavPropertiesContext $webDavPropertiesContext;

	/**
	 * @param string$user
	 * @param string $path
	 * @param string|null $spaceId
	 *
	 * @return ResponseInterface
	 */
	public function userFavoritesElement(string $user, string $path, ?string $spaceId = null): ResponseInterface {
		return $this->changeFavStateOfAnElement(
			$user,
			$path,
			1,
			$spaceId,
		);
	}

	/**
	 * @When user :user favorites element :path using the WebDAV API
	 *
	 * @param string $user
	 * @param string $path
	 *
	 * @return void
	 */
	public function userFavoritesElementUsingWebDavApi(string $user, string $path): void {
		$this->featureContext->setResponse($this->userFavoritesElement($user, $path));
	}

	/**
	 * @Given user :user has favorited element :path
	 *
	 * @param string $user
	 * @param string $path
	 *
	 * @return void
	 */
	public function userHasFavoritedElementUsingWebDavApi(string $user, string $path): void {
		$this->featureContext->theHTTPStatusCodeShouldBe(207, '', $this->userFavoritesElement($user, $path));
	}

	/**
	 * @param string $user
	 * @param string $path
	 *
	 * @return ResponseInterface
	 */
	public function userUnfavoritesElement(string $user, string $path): ResponseInterface {
		return $this->changeFavStateOfAnElement(
			$user,
			$path,
			0,
			null,
		);
	}

	/**
	 * @When user :user unfavorites element :path using the WebDAV API
	 *
	 * @param string $user
	 * @param string $path
	 *
	 * @return void
	 */
	public function userUnfavoritesElementUsingWebDavApi(string $user, string $path): void {
		$this->featureContext->setResponse($this->userUnfavoritesElement($user, $path));
	}

	/**
	 * @Then /^user "([^"]*)" should (not|)\s?have the following favorited items$/
	 *
	 * @param string $user
	 * @param string $shouldOrNot (not|)
	 * @param TableNode $expectedElements
	 *
	 * @return void
	 */
	public function checkFavoritedElements(
		string $user,
		string $shouldOrNot,
		TableNode $expectedElements,
	): void {
		$user = $this->featureContext->getActualUsername($user);
		$this->userListsFavorites($user);
		$this->featureContext->propfindResultShouldContainEntries(
			$shouldOrNot,
			$expectedElements,
			$user,
		);
	}

	/**
	 * @When /^user "([^"]*)" lists the favorites and limits the result to ([\d*]) elements using the WebDAV API$/
	 *
	 * @param string $user
	 * @param int|null $limit
	 *
	 * @return void
	 */
	public function userListsFavorites(string $user, ?int $limit = null): void {
		$renamedUser = $this->featureContext->getActualUsername($user);
		$baseUrl = $this->featureContext->getBaseUrl();
		$password = $this->featureContext->getPasswordForUser($user);
		$body
			= "<?xml version='1.0' encoding='utf-8' ?>\n" .
			"	<oc:filter-files xmlns:a='DAV:' xmlns:oc='http://owncloud.org/ns' >\n" .
			"		<a:prop><oc:favorite/></a:prop>\n" .
			"		<oc:filter-rules><oc:favorite>1</oc:favorite></oc:filter-rules>\n";

		if ($limit !== null) {
			$body .= "		<oc:search>\n" .
				"			<oc:limit>$limit</oc:limit>\n" .
				"		</oc:search>\n";
		}

		$body .= "	</oc:filter-files>";
		$response = WebDavHelper::makeDavRequest(
			$baseUrl,
			$renamedUser,
			$password,
			"REPORT",
			"/",
			null,
			null,
			$body,
			$this->featureContext->getDavPathVersion(),
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Then /^as user "([^"]*)" (?:file|folder|entry) "([^"]*)" should be favorited$/
	 *
	 * @param string $user
	 * @param string $path
	 * @param integer $expectedValue 0|1
	 * @param string|null $spaceId
	 *
	 * @return void
	 */
	public function asUserFileOrFolderShouldBeFavorited(
		string $user,
		string $path,
		int $expectedValue = 1,
		?string $spaceId = null,
	): void {
		$property = "oc:favorite";
		$this->webDavPropertiesContext->checkPropertyOfAFolder(
			$user,
			$path,
			$property,
			(string)$expectedValue,
			null,
			$spaceId,
		);
	}

	/**
	 * @Then /^as user "([^"]*)" (?:file|folder|entry) "([^"]*)" should not be favorited$/
	 *
	 * @param string $user
	 * @param string $path
	 *
	 * @return void
	 */
	public function asUserFileShouldNotBeFavorited(string $user, string $path): void {
		$this->asUserFileOrFolderShouldBeFavorited($user, $path, 0);
	}

	/**
	 * Set the elements of a proppatch
	 *
	 * @param string $user
	 * @param string $path
	 * @param int|null $favOrUnfav 1 = favorite, 0 = unfavorite
	 * @param string|null $spaceId
	 *
	 * @return ResponseInterface
	 */
	public function changeFavStateOfAnElement(
		string $user,
		string $path,
		?int $favOrUnfav,
		?string $spaceId,
	): ResponseInterface {
		$renamedUser = $this->featureContext->getActualUsername($user);
		return WebDavHelper::proppatch(
			$this->featureContext->getBaseUrl(),
			$renamedUser,
			$this->featureContext->getPasswordForUser($user),
			$path,
			'favorite',
			(string)$favOrUnfav,
			"oc='http://owncloud.org/ns'",
			$this->featureContext->getDavPathVersion(),
			'files',
			$spaceId,
		);
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
	public function before(BeforeScenarioScope $scope): void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context
		$this->featureContext = BehatHelper::getContext($scope, $environment, 'FeatureContext');
		$this->webDavPropertiesContext = BehatHelper::getContext(
			$scope,
			$environment,
			'WebDavPropertiesContext',
		);
	}
}
