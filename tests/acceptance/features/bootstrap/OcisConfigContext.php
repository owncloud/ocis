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

/**
 * steps needed to re-configure oCIS server
 */
class OcisConfigContext implements Context {
	/**
	 * @Given async upload has been enabled with post processing delayed to :delayTime seconds
	 *
	 * @param string $delayTime
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function asyncUploadHasbeenEnabledWithDelayedPostProcessing(string $delayTime): void {
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
	}

	/**
	 * @Given cors allowed origins has been set to :allowedOrigins
	 *
	 * @param string $allowedOrigins
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function corsAllowedOriginsHasBeenSet(string $allowedOrigins): void {
		$envs = [
			"OCIS_CORS_ALLOW_ORIGINS" => $allowedOrigins,
		];

		$response =  OcisConfigHelper::reConfigureOcis($envs);
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set OCIS_CORS_ALLOW_ORIGINS=" . $allowedOrigins
		);
	}

	/**
	 * @AfterScenario @env-config
	 *
	 * @return void
	 */
	public function rollbackOcis(): void {
		$response = OcisConfigHelper::rollbackOcis();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to rollback ocis server. Check if oCIS is started with ociswrapper."
		);
	}
}
