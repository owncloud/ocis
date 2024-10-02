<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Amrita Shrestha <amrita@jankaritech.com>
 * @copyright Copyright (c) 2024 Amrita Shrestha amrita@jankaritech.com
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

use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\ResponseInterface;

/**
 * A helper class for managing wopi requests
 */
class CollaborationHelper {
	/**
	 * @param string $fileId
	 * @param string $app
	 * @param string $username
	 * @param string $password
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string|null $viewMode
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function sendPOSTRequestToAppOpen(
		string $fileId,
		string $app,
		string $username,
		string $password,
		string $baseUrl,
		string $xRequestId,
		?string $viewMode = null,
	): ResponseInterface {
		$url = $baseUrl . "/app/open?app_name=$app&file_id=$fileId";
		if ($viewMode) {
			$url .= "&view_mode=$viewMode";
		}

		return HttpRequestHelper::post(
			$url,
			$xRequestId,
			$username,
			$password,
			['Content-Type' => 'application/json']
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $user
	 * @param string $password
	 * @param string $parentContainerId
	 * @param string $file
	 * @param array|null $headers
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function createFile(
		string $baseUrl,
		string $xRequestId,
		string $user,
		string $password,
		string $parentContainerId,
		string $file,
		?array $headers = null
	): ResponseInterface {
		$url = $baseUrl . "/app/new?parent_container_id=$parentContainerId&filename=$file";
		return HttpRequestHelper::post(
			$url,
			$xRequestId,
			$user,
			$password,
			$headers
		);
	}
}
