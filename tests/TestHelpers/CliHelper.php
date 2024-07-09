<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Sajan Gurung <sajan@jankaritech.com>
 * @copyright Copyright (c) 2024 Sajan Gurung sajan@jankaritech.com
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

use GuzzleHttp\Exception\ConnectException;
use GuzzleHttp\Exception\GuzzleException;
use GuzzleHttp\Psr7\Request;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\OcisConfigHelper;

/**
 * A helper class for running oCIS CLI commands
 */
class CliHelper {
	/**
	 * @param array $body
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function runCommand(array $body): ResponseInterface {
		$url = OcisConfigHelper::getWrapperUrl() . "/command";
		return OcisConfigHelper::sendRequest($url, "POST", \json_encode($body));
	}
}
