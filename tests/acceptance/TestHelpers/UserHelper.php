<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Artur Neumann <artur@jankaritech.com>
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

use Exception;
use GuzzleHttp\Exception\ClientException;
use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\ResponseInterface;

/**
 * Helper to administrate users (and groups) through the provisioning API
 *
 * @author Artur Neumann <artur@jankaritech.com>
 *
 */
class UserHelper {
	/**
	 *
	 * @param string|null $baseUrl
	 * @param string $user
	 * @param string $key
	 * @param string $value
	 * @param string $adminUser
	 * @param string $adminPassword
	 * @param int|null $ocsApiVersion
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function editUser(
		?string $baseUrl,
		string $user,
		string $key,
		string $value,
		string $adminUser,
		string $adminPassword,
		?int $ocsApiVersion = 2,
	): ResponseInterface {
		return OcsApiHelper::sendRequest(
			$baseUrl,
			$adminUser,
			$adminPassword,
			"PUT",
			"/cloud/users/" . $user,
			["key" => $key, "value" => $value],
			$ocsApiVersion,
		);
	}

	/**
	 * Send batch requests to edit the user.
	 * This will send multiple requests in parallel to the server which will be faster in comparison to waiting for each request to complete.
	 *
	 * @param string|null $baseUrl
	 * @param array|null $editData ['user' => '', 'key' => '', 'value' => '']
	 * @param string|null $adminUser
	 * @param string|null $adminPassword
	 * @param int|null $ocsApiVersion
	 *
	 * @return array
	 * @throws Exception
	 */
	public static function editUserBatch(
		?string $baseUrl,
		?array $editData,
		?string $adminUser,
		?string $adminPassword,
		?int $ocsApiVersion = 2,
	): array {
		$requests = [];
		$client = HttpRequestHelper::createClient(
			$adminUser,
			$adminPassword,
		);

		foreach ($editData as $data) {
			$path = "/cloud/users/" . $data['user'];
			$body = ["key" => $data['key'], 'value' => $data["value"]];
			// Create the OCS API requests and push them to an array.
			$requests[] = OcsApiHelper::createOcsRequest(
				$baseUrl,
				'PUT',
				$path,
				$body,
			);
		}
		// Send the array of requests at once in parallel.
		$results = HttpRequestHelper::sendBatchRequest($requests, $client);

		foreach ($results as $e) {
			if ($e instanceof ClientException) {
				$httpStatusCode = $e->getResponse()->getStatusCode();
				$reasonPhrase = $e->getResponse()->getReasonPhrase();
				throw new Exception(
					"Unexpected failure when editing a user: HTTP status $httpStatusCode HTTP reason $reasonPhrase",
				);
			}
		}
		return $results;
	}

	/**
	 *
	 * @param string|null $baseUrl
	 * @param string|null $userName
	 * @param string|null $adminUser
	 * @param string|null $adminPassword
	 * @param int|null $ocsApiVersion
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function getUser(
		?string $baseUrl,
		?string $userName,
		?string $adminUser,
		?string $adminPassword,
		?int $ocsApiVersion = 2,
	): ResponseInterface {
		return OcsApiHelper::sendRequest(
			$baseUrl,
			$adminUser,
			$adminPassword,
			"GET",
			"/cloud/users/" . $userName,
			[],
			$ocsApiVersion,
		);
	}

	/**
	 *
	 * @param string|null $baseUrl
	 * @param string|null $userName
	 * @param string|null $adminUser
	 * @param string|null $adminPassword
	 * @param int|null $ocsApiVersion
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function deleteUser(
		?string $baseUrl,
		?string $userName,
		?string $adminUser,
		?string $adminPassword,
		?int $ocsApiVersion = 2,
	): ResponseInterface {
		return OcsApiHelper::sendRequest(
			$baseUrl,
			$adminUser,
			$adminPassword,
			"DELETE",
			"/cloud/users/" . $userName,
			[],
			$ocsApiVersion,
		);
	}

	/**
	 *
	 * @param string|null $baseUrl
	 * @param string|null $user
	 * @param string|null $group
	 * @param string|null $adminUser
	 * @param string|null $adminPassword
	 * @param int|null $ocsApiVersion (1|2)
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function addUserToGroup(
		?string $baseUrl,
		?string $user,
		?string $group,
		?string $adminUser,
		?string $adminPassword,
		?int $ocsApiVersion = 2,
	): ResponseInterface {
		return OcsApiHelper::sendRequest(
			$baseUrl,
			$adminUser,
			$adminPassword,
			"POST",
			"/cloud/users/" . $user . "/groups",
			['groupid' => $group],
			$ocsApiVersion,
		);
	}

	/**
	 *
	 * @param string|null $baseUrl
	 * @param string|null $user
	 * @param string|null $group
	 * @param string|null $adminUser
	 * @param string|null $adminPassword
	 * @param int|null $ocsApiVersion (1|2)
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function removeUserFromGroup(
		?string $baseUrl,
		?string $user,
		?string $group,
		?string $adminUser,
		?string $adminPassword,
		?int $ocsApiVersion = 2,
	): ResponseInterface {
		return OcsApiHelper::sendRequest(
			$baseUrl,
			$adminUser,
			$adminPassword,
			"DELETE",
			"/cloud/users/" . $user . "/groups",
			['groupid' => $group],
			$ocsApiVersion,
		);
	}

	/**
	 *
	 * @param string|null $baseUrl
	 * @param string|null $adminUser
	 * @param string|null $adminPassword
	 * @param string|null $search
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function getGroups(
		?string $baseUrl,
		?string $adminUser,
		?string $adminPassword,
		?string $search = "",
	): ResponseInterface {
		return OcsApiHelper::sendRequest(
			$baseUrl,
			$adminUser,
			$adminPassword,
			"GET",
			"/cloud/groups?search=" . $search,
		);
	}

	/**
	 *
	 * @param string|null $baseUrl
	 * @param string|null $adminUser
	 * @param string|null $adminPassword
	 * @param string|null $search
	 *
	 * @return string[]
	 * @throws Exception|GuzzleException
	 */
	public static function getGroupsAsArray(
		?string $baseUrl,
		?string $adminUser,
		?string $adminPassword,
		?string $search = "",
	): array {
		$result = self::getGroups(
			$baseUrl,
			$adminUser,
			$adminPassword,
			$search,
		);
		$groups = HttpRequestHelper::getResponseXml($result, __METHOD__)->xpath(".//groups")[0];
		$return = [];
		foreach ($groups as $group) {
			$return[] = $group->__toString();
		}
		return $return;
	}
}
