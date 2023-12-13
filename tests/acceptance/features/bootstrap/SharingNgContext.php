<?php declare(strict_types=1);
/**
 * @author Viktor Scharf <scharf.vi@gmail.com>
 *
 * @copyright Copyright (c) 2023, ownCloud GmbH
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
use TestHelpers\GraphHelper;

require_once 'bootstrap.php';

/**
 * Acceptance test steps related to testing sharing ng features
 */
class SharingNgContext implements Context {
	private FeatureContext $featureContext;
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
	 * @When /^user "([^"]*)" gets permissions list for (folder|file) "([^"]*)" of the space "([^"]*)" using the Graph API$/
	 *
	 * @param string $user
	 * @param string $fileOrFolder   (file|folder)
	 * @param string $resource
	 * @param string $space
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theUserPermissionsListOfResource(string $user, string $fileOrFolder, string $resource, string $space):void {
		$spaceId = ($this->spacesContext->getSpaceByName($user, $space))["id"];
	
		if ($fileOrFolder === 'folder') {
			$itemId = $this->spacesContext->getResourceId($user, $space, $resource);
		} else {
			$itemId = $this->spacesContext->getFileId($user, $space, $resource);
		}
		$this->featureContext->setResponse(
			GraphHelper::getPermissionsList(
				$this->featureContext->getBaseUrl(),
				$this->featureContext->getStepLineRef(),
				$user,
				$this->featureContext->getPasswordForUser($user),
				$spaceId,
				$itemId
			)
		);
	}
}
