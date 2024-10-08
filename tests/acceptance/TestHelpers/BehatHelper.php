<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Sajan Gurung <sajan@jankaritech.com>
 * @copyright Copyright (c) 2017 Artur Neumann artur@jankaritech.com
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
namespace TestHelpers;

use Behat\Behat\Context\Context;
use Behat\Behat\Context\Environment\InitializedContextEnvironment;
use Behat\Behat\Context\Exception\ContextNotFoundException;
use Behat\Behat\Hook\Scope\ScenarioScope;

/**
 * Helper for Behat environment configuration
 *
 */
class BehatHelper {
	/**
	 * @param ScenarioScope $scope
	 * @param InitializedContextEnvironment $environment
	 * @param string $class
	 *
	 * @return Context
	 */
	public static function getContext(ScenarioScope $scope, InitializedContextEnvironment $environment, string $class): Context {
		try {
			return $environment->getContext($class);
		} catch (ContextNotFoundException $e) {
			print_r("[INFO] '$class' context not found. Registering...\n");
			$context = new $class();
			$environment->registerContext($context);
			if (\method_exists($context, 'before')) {
				$context->before($scope);
			}
			return $environment->getContext($class);
		}
	}
}
