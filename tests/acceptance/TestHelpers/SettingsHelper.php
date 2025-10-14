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

use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\ResponseInterface;
use PHPUnit\Framework\Assert;

/**
 * A helper class for ocis settings
 */
class SettingsHelper {
	private static string $settingsEndpoint = '/api/v0/settings/';

	private const BUNDLE_ID = "2a506de7-99bd-4f0d-994e-c38e72c28fd9";
	public const PROFILE_SETTINGS = [
		'language' => 'aa8cfbe5-95d4-4f7e-a032-c3c01f5f062f',
		'auto-accept-shares' => 'ec3ed4a3-3946-4efc-8f9f-76d38b12d3a9',
	];

	private const DEFAULT_SETTING_EVENTS = [
		'Disable Email Notifications' => '33ffb5d6-cd07-4dc0-afb0-84f7559ae438',
		'Auto Accept Shares' => 'ec3ed4a3-3946-4efc-8f9f-76d38b12d3a9',
		'File Rejected' => 'fe0a3011-d886-49c8-b797-33d02fa426ef',
		'Email Sending Interval' => '08dec2fe-3f97-42a9-9d1b-500855e92f25',
		'Share Received' => '872d8ef6-6f2a-42ab-af7d-f53cc81d7046',
		'Share Removed' => 'd7484394-8321-4c84-9677-741ba71e1f80',
		'Share Expired' => 'e1aa0b7c-1b0f-4072-9325-c643c89fee4e',
		'Added As Space Member' => '694d5ee1-a41c-448c-8d14-396b95d2a918',
		'Removed As Space Member' => '26c20e0e-98df-4483-8a77-759b3a766af0',
		'Space Deleted' => '094ceca9-5a00-40ba-bb1a-bbc7bccd39ee',
		'Space Disabled' => 'eb5c716e-03be-42c6-9ed1-1105d24e109f',
		'Space Membership Expired' => '7275921e-b737-4074-ba91-3c2983be3edd',
	];

	/**
	 * @return string
	 */
	public static function getBundleId(): string {
		return self::BUNDLE_ID;
	}

	/**
	 * @param string $event
	 *
	 * @return string
	 */
	public static function getSettingIdUsingEventName(string $event): string {
		return self::DEFAULT_SETTING_EVENTS[$event];
	}

	/**
	 * @param string $baseUrl
	 * @param string $path
	 *
	 * @return string
	 */
	public static function buildFullUrl(string $baseUrl, string $path): string {
		return $baseUrl . self::$settingsEndpoint . $path;
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 * @param array $headers
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function getBundlesList(
		string $baseUrl,
		string $user,
		string $password,
		array $headers = [],
	): ResponseInterface {
		$fullUrl = self::buildFullUrl($baseUrl, "bundles-list");
		return HttpRequestHelper::post(
			$fullUrl,
			$user,
			$password,
			$headers,
			"{}",
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 * @param string $bundleName
	 *
	 * @return array
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function getBundleByName(
		string $baseUrl,
		string $user,
		string $password,
		string $bundleName,
	): array {
		$response = self::getBundlesList($baseUrl, $user, $password);
		Assert::assertEquals(201, $response->getStatusCode(), "Failed to get bundles list");

		$bundlesList = HttpRequestHelper::getJsonDecodedResponseBodyContent($response);
		foreach ($bundlesList->bundles as $value) {
			if ($value->displayName === $bundleName) {
				return json_decode(json_encode($value), true);
			}
		}
		return [];
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 * @param array $headers
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function getRolesList(
		string $baseUrl,
		string $user,
		string $password,
		array $headers = [],
	): ResponseInterface {
		$fullUrl = self::buildFullUrl($baseUrl, "roles-list");
		return HttpRequestHelper::post(
			$fullUrl,
			$user,
			$password,
			$headers,
			"{}",
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 * @param string $assigneeId
	 * @param string $roleId
	 * @param array $headers
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function assignRoleToUser(
		string $baseUrl,
		string $user,
		string $password,
		string $assigneeId,
		string $roleId,
		array $headers = [],
	): ResponseInterface {
		$fullUrl = self::buildFullUrl($baseUrl, "assignments-add");
		$body = json_encode(["account_uuid" => $assigneeId, "role_id" => $roleId], JSON_THROW_ON_ERROR);
		return HttpRequestHelper::post(
			$fullUrl,
			$user,
			$password,
			$headers,
			$body,
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 * @param string $userId
	 * @param array $headers
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function getAssignmentsList(
		string $baseUrl,
		string $user,
		string $password,
		string $userId,
		array $headers = [],
	): ResponseInterface {
		$fullUrl = self::buildFullUrl($baseUrl, "assignments-list");
		$body = json_encode(["account_uuid" => $userId], JSON_THROW_ON_ERROR);
		return HttpRequestHelper::post(
			$fullUrl,
			$user,
			$password,
			$headers,
			$body,
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 * @param array $headers
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function getValuesList(
		string $baseUrl,
		string $user,
		string $password,
		array $headers = [],
	): ResponseInterface {
		$fullUrl = self::buildFullUrl($baseUrl, "values-list");
		$body = json_encode(["account_uuid" => "me"], JSON_THROW_ON_ERROR);
		return HttpRequestHelper::post(
			$fullUrl,
			$user,
			$password,
			$headers,
			$body,
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 *
	 * @return bool
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function getAutoAcceptSharesDefaultValue(
		string $baseUrl,
		string $user,
		string $password,
	): bool {
		$response = self::getValuesList($baseUrl, $user, $password);
		Assert::assertEquals(201, $response->getStatusCode(), "Failed to get values list");

		$valuesList = HttpRequestHelper::getJsonDecodedResponseBodyContent($response);

		if (empty($valuesList)) {
			// default is true
			return true;
		}

		foreach ($valuesList->values as $value) {
			if ($value->identifier->setting === "auto-accept-shares") {
				return $value->value->boolValue;
			}
		}

		return true;
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 *
	 * @return string
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function getLanguageSettingValue(
		string $baseUrl,
		string $user,
		string $password,
	): string {
		$response = self::getValuesBySettingID(
			$baseUrl,
			self::PROFILE_SETTINGS['language'],
			$user,
			$password,
		);
		Assert::assertEquals(201, $response->getStatusCode(), "Failed to get values list");

		$valuesList = HttpRequestHelper::getJsonDecodedResponseBodyContent($response);

		// if no language is set, return English as the default language
		if (empty($valuesList) || empty($valuesList->value)) {
			return "en";
		}
		if ($valuesList->value->identifier->setting === "language") {
			return $valuesList->value->value->listValue->values[0]->stringValue;
		}
		// if a language setting was still not found, return English
		return "en";
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 * @param string $body
	 * @param array $headers
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function updateSettings(
		string $baseUrl,
		string $user,
		string $password,
		string $body,
		array $headers = [],
	): ResponseInterface {
		$fullUrl = self::buildFullUrl($baseUrl, "values-save");
		return HttpRequestHelper::post(
			$fullUrl,
			$user,
			$password,
			$headers,
			$body,
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $settingId
	 * @param string $user
	 * @param string $password
	 *
	 * @return ResponseInterface
	 *
	 * @throws GuzzleException
	 * @throws \JsonException
	 */
	public static function getValuesBySettingID(
		string $baseUrl,
		string $settingId,
		string $user,
		string $password,
	): ResponseInterface {
		$fullUrl = self::buildFullUrl($baseUrl, "values-get-by-unique-identifiers");
		$body = json_encode(
			[
				"accountUuid" => "me",
				"settingId" => $settingId,
			],
			JSON_THROW_ON_ERROR,
		);
		return HttpRequestHelper::post(
			$fullUrl,
			$user,
			$password,
			[],
			$body,
		);
	}
}
