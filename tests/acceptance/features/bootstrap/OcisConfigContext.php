<?php declare(strict_types=1);
/**
 * ownCloud
 * 
 * @author Sajan Gurung <sajan@jankaritech.com>
 * @copyright Copyright (c) 2023 Sajan Gurung sajan@jankaritech.com
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
use TestHelpers\OcisConfigHelper;
use PHPUnit\Framework\Assert;

 class OcisConfigContext implements Context {
	/**
	 * This is used to determine if the ocis config has been changed via tests
	 * 
	 * @var bool
	 */
	private static $touched = false;

	/**
	 * @return bool
	 */
	public static function hasTouched(): bool {
		return self::$touched;
	}

	/**
	 * set the touched flag
	 */
	public static function setTouched(): void {
		self::$touched = true;
	}

	/**
	 * reset the touched flag
	 */
	public static function reset(): void {
		self::$touched = false;
	}

	/**
	 * @Given async upload has been enabled with post processing delayed to :delayTime seconds
	 *
	 * @param string $delayTime
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function asyncUploadHasbeenEnabledWithDeplayedPostProcessing(string $delayTime): void {
		$envs = [
			"STORAGE_USERS_OCIS_ASYNC_UPLOADS" => true,
			"OCIS_EVENTS_ENABLE_TLS" => false,
			"POSTPROCESSING_DELAY" => $delayTime . "s",
		];

		$response =  OcisConfigHelper::reConfigureOcis($envs);
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set async upload with delayed post processing"
		);
		// ocis config has been changed
		// set the touched flag
		self::setTouched();
	}

	/**
	 * @Given cors allowed origins has been set to :allowedOrigins
	 *
	 * @param string $allowedOrigins
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function corsAllowedOriginsHasbeenSet(string $allowedOrigins): void {
		$envs = [
			"CORS_ALLOWED_ORIGINS" => $allowedOrigins,
		];

		$response =  OcisConfigHelper::reConfigureOcis($envs);
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set CORS_ALLOWED_ORIGINS"
		);
		// ocis config has been changed
		// set the touched flag
		self::setTouched();
	}

	/**
	 * @AfterScenario
	 *
	 * @return void
	 */
	public function rollbackOcis(): void {
		$response = OcisConfigHelper::rollbackOcis();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to rollback ocis server"
		);

		self::reset();
	}
 }