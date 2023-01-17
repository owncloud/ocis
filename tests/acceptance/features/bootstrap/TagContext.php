<?php declare(strict_types=1);
/**
 * @author Viktor Scharf <scharf.vi@gmail.com>
 *
 * @copyright Copyright (c) 2022, ownCloud GmbH
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
use PHPUnit\Framework\Assert;
use TestHelpers\GraphHelper;

require_once 'bootstrap.php';

/**
 * Acceptance test steps related to testing tags features
 */
class TagContext implements Context {
	/**
	 *
	 * @var FeatureContext
	 */
	private $featureContext;

	/**
	 * @var SpacesContext
	 */
	private SpacesContext $spacesContext;

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
		// Get all the contexts you need in this context from here
		$this->featureContext = $environment->getContext('FeatureContext');
		$this->spacesContext = $environment->getContext('SpacesContext');
	}

	/**
	 * @var array
	 */
	private $createdTags = [];

	/**
	 * @When /^user "([^"]*)" creates the following tags for (folder|file)\s?"([^"]*)" of space "([^"]*)":$/
	 *
	 * @param string $user
	 * @param string $fileOrFolder   (file|folder)
	 * @param string $resource
	 * @param string $space
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theUserCreatesFollowingTags(string $user, string $fileOrFolder, string $resource, string $space, TableNode $table):void {
		$tagNameArray = [];
		foreach ($table->getRows() as $value) {
			array_push($tagNameArray, $value[0]);
		}

		if ($fileOrFolder === 'folder') {
			$resourceId = $this->spacesContext->getFolderId($user, $space, $resource);
		} else {
			$resourceId = $this->spacesContext->getFileId($user, $space, $resource);
		}

		$response = GraphHelper::createTags(
			$this->featureContext->getBaseUrl(),
			$this->featureContext->getStepLineRef(),
			$user,
			$this->featureContext->getPasswordForUser($user),
			$resourceId,
			$tagNameArray
		);
		$this->featureContext->setResponse($response);
	}

	/**
	 * @Given /^user "([^"]*)" has created the following tags a (folder|file)\s?"([^"]*)" of the space "([^"]*)":$/
	 *
	 * @param string $user
	 * @param string $fileOrFolder   (file|folder)
	 * @param string $resource
	 * @param string $space
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theUserHasCreatedFollowingTags(string $user, string $fileOrFolder, string $resource, string $space, TableNode $table):void {
		$this->theUserCreatesFollowingTags($user, $fileOrFolder, $resource, $space, $table);
		$this->featureContext->theHttpStatusCodeShouldBe(200);
	}

	/**
	 * @When user :user lists all available tag(s) via the GraphApi
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theUserGetsAllAvailableTags(string $user):void {
		$this->featureContext->setResponse(
			GraphHelper::getTags(
				$this->featureContext->getBaseUrl(),
				$user,
				$this->featureContext->getPasswordForUser($user)
			)
		);
	}

	/**
	 * @Then the response should contain following tag(s):
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theFollowingTagsShouldExistForUser(TableNode $table):void {
		$rows = $table->getRows();
		foreach ($rows as $row) {
			$responseArray = $this->featureContext->getJsonDecodedResponse($this->featureContext->getResponse())['value'];
			Assert::assertTrue(\in_array($row[0], $responseArray), "the response does not contain the tag $row[0]");
		}
	}
}
