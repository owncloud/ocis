<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Prajwol Amatya <prajwol@jankaritech.com>
 * @copyright Copyright (c) 2022 Prajwol Amatya prajwol@jankaritech.com
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
use GuzzleHttp\Exception\GuzzleException;
use InvalidArgumentException;
use JsonException;
use Psr\Http\Message\ResponseInterface;

/**
 * A helper class for Keycloak admin API requests.
 */
class KeycloakHelper {
	private const OCIS_KEYCLOAK_USER_ROLES = [
		'Admin' => 'ocisAdmin',
		'Space Admin' => 'ocisSpaceAdmin',
		'User' => 'ocisUser',
		'User Light' => 'ocisGuest',
	];
	private const OCIS_KEYCLOAK_USER_ROLE_IDS = [
		'ocisAdmin' => '0bb40fa2-4490-4687-9159-b1d27ec7423a',
		'ocisSpaceAdmin' => 'bd5f5012-48bb-4ea4-bfe6-0623e3ca0552',
		'ocisUser' => '8c79ff81-c256-48fd-b0b9-795c7941eedf',
		'ocisGuest' => '7eedfa6d-a2d9-4296-b6db-e75e4e9c0963',
		'offline_access' => 'e2145b30-bf6f-49fb-af3f-1b40168bfcef',
	];
	private static ?string $adminAccessToken = null;

	/**
	 * @return bool
	 */
	public static function isTestingWithKeycloak(): bool {
		return (\getenv('KEYCLOAK') === "true");
	}

	/**
	 * @return string
	 */
	public static function getKeycloakUrl(): string {
		$keycloakUrl = \getenv('KEYCLOAK_URL');
		if ($keycloakUrl !== false && $keycloakUrl !== '') {
			return $keycloakUrl;
		}
		return 'https://keycloak.owncloud.test';
	}

	/**
	 * @param string $accessToken
	 *
	 * @return void
	 */
	public static function setAdminAccessToken(string $accessToken): void {
		self::$adminAccessToken = $accessToken;
	}

	/**
	 * @return string
	 * @throws GuzzleException
	 */
	public static function getAdminAccessToken(): string {
		if (self::$adminAccessToken === null) {
			self::setAdminAccessToken(self::generateAdminAccessToken());
		}

		return (string)self::$adminAccessToken;
	}

	/**
	 * @param string $roleName
	 *
	 * @return array
	 */
	private static function getRealmRole(string $roleName): array {
		$roleId = self::OCIS_KEYCLOAK_USER_ROLE_IDS[$roleName]
			?? throw new InvalidArgumentException("Unknown Keycloak role: {$roleName}");
		return [
			'id' => $roleId,
			'name' => $roleName,
		];
	}

	/**
	 * @param string $username
	 * @param string $password
	 * @param string|null $email
	 * @param string|null $displayName
	 *
	 * @return ResponseInterface
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public static function createUser(
		string $username,
		string $password,
		?string $email = null,
		?string $displayName = null,
	): ResponseInterface {
		$accessToken = self::getAdminAccessToken();
		$url = self::getKeycloakUrl() . '/admin/realms/oCIS/users';

		return HttpRequestHelper::post(
			$url,
			null,
			null,
			[
				'Authorization' => 'Bearer ' . $accessToken,
				'Content-Type' => 'application/json',
			],
			self::prepareCreateUserPayload($username, $password, $email, $displayName),
		);
	}

	/**
	 * @param string $uuid
	 * @param string $role
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public static function assignRole(
		string $uuid,
		string $role,
	): ResponseInterface {
		$url = self::getKeycloakUrl() . "/admin/realms/oCIS/users/" . $uuid . "/role-mappings/realm";
		$ocisRole = self::OCIS_KEYCLOAK_USER_ROLES[$role]
			?? throw new InvalidArgumentException("Unknown oCIS role: {$role}");
		$body = [
			self::getRealmRole($ocisRole),
			self::getRealmRole('offline_access'),
		];
		return HttpRequestHelper::post(
			$url,
			null,
			null,
			[
				"Content-Type" => "application/json",
				"Authorization" => "Bearer " . self::getAdminAccessToken(),
			],
			json_encode($body, JSON_THROW_ON_ERROR),
		);
	}

	/**
	 * @return string
	 * @throws Exception
	 * @throws GuzzleException
	 */
	private static function generateAdminAccessToken(): string {
		$url = self::getKeycloakUrl()
			. '/realms/master/protocol/openid-connect/token';
		$response = HttpRequestHelper::post(
			$url,
			null,
			null,
			['Content-Type' => 'application/x-www-form-urlencoded'],
			[
				'client_id' => 'admin-cli',
				'username' => 'admin',
				'password' => 'admin',
				'grant_type' => 'password',
			],
		);

		if ($response->getStatusCode() >= 400) {
			throw new Exception(
				__METHOD__
				. ' failed to get Keycloak admin access token, status '
				. $response->getStatusCode()
				. ', response '
				. (string)$response->getBody(),
			);
		}

		$decodedResponse = HttpRequestHelper::getJsonDecodedResponseBodyContent($response);
		if (!isset($decodedResponse->access_token)) {
			throw new Exception(__METHOD__ . ' could not find access_token in Keycloak token response');
		}
		return (string)$decodedResponse->access_token;
	}

	/**
	 * @param string $username
	 * @param string $password
	 * @param string|null $email
	 * @param string|null $displayName
	 *
	 * @return string
	 * @throws JsonException
	 */
	private static function prepareCreateUserPayload(
		string $username,
		string $password,
		?string $email = null,
		?string $displayName = null,
	): string {
		$firstName = $username;
		$lastName = '';
		if ($displayName !== null && \trim($displayName) !== '') {
			$nameParts = \preg_split('/\s+/', \trim($displayName), 2);
			if ($nameParts !== false && isset($nameParts[0])) {
				$firstName = $nameParts[0];
			}
			if ($nameParts !== false && isset($nameParts[1])) {
				$lastName = $nameParts[1];
			}
		}

		$payload = [
			'username' => $username,
			'credentials' => [[
				'value' => $password,
				'type' => 'password',
			]],
			'firstName' => $firstName,
			'lastName' => $lastName,
			'emailVerified' => true,
			'enabled' => true,
		];

		if ($email !== null) {
			$payload['email'] = $email;
		}

		return \json_encode($payload, JSON_THROW_ON_ERROR);
	}

	/**
	 * @return array
	 * @throws GuzzleException
	 */
	public static function getAuthorizationEndPoint(): array {
		$loginParams = [
			'client_id' => 'web',
			'redirect_uri' => OcisHelper::getServerUrl() . '/oidc-callback.html',
			'response_mode' => 'query',
			'response_type' => 'code',
			'scope' => 'openid profile email acr',
		];
		$queryString = \http_build_query($loginParams);
		$authUrl = self::getKeycloakUrl() . "/realms/oCIS/protocol/openid-connect/auth?" . $queryString;
		$response = HttpRequestHelper::get(
			$authUrl,
		);
		$cookie = $response->getHeader("Set-Cookie")[0];
		$htmlData = $response->getBody()->getContents();
		if (!preg_match('/action="([^"]+)"/i', $htmlData, $match)) {
			throw new Exception('No authorization url found in the HTML response body.');
		}
		$authorizationUrl = $match[1];
		return [$authorizationUrl, $cookie];
	}

	/**
	 * @param array $user
	 * @param string $authorizationUrl
	 * @param string $cookie
	 *
	 * @return string
	 * @throws GuzzleException
	 */
	public static function getCode(array $user, string $authorizationUrl, string $cookie): string {
		$authCodeResponse = HttpRequestHelper::post(
			$authorizationUrl,
			null,
			null,
			[
				'Cookie' => $cookie,
			],
			[
				'username' => $user['actualUsername'],
				'password' => $user['password'],
			],
		);
		$locationHeader = $authCodeResponse->getHeader('Location');
		$queryString = parse_url($locationHeader[0], PHP_URL_QUERY);
		parse_str($queryString, $urlParams);
		return $urlParams['code'];
	}

	/**
	 * @param string $authorizationCode
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function getToken(string $authorizationCode): ResponseInterface {
		$url = self::getKeycloakUrl() . '/realms/oCIS/protocol/openid-connect/token';
		$tokenResponse = HttpRequestHelper::post(
			$url,
			null,
			null,
			null,
			[
				'client_id' => 'web',
				'code' => $authorizationCode,
				'redirect_uri' => OcisHelper::getServerUrl() . '/oidc-callback.html',
				'grant_type' => 'authorization_code',
			],
		);
		if ($tokenResponse->getStatusCode() !== 200) {
			throw new Exception(
				'Failed to retrieve token: Expected status code to be 200 but received ' .
				$tokenResponse->getStatusCode() .
				'.\nMessage: ' .
				$tokenResponse->getBody()->getContents(),
			);
		}
		return $tokenResponse;
	}

	/**
	 * @param array $user
	 *
	 * @return array
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public static function setAccessTokenForKeycloakOcisUser(array $user): array {
		[$authorizationUrl, $cookie] = self::getAuthorizationEndPoint();
		$authorizationCode = self::getCode($user, $authorizationUrl, $cookie);
		$tokenResponse = self::getToken($authorizationCode);
		return json_decode($tokenResponse->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
	}

	/**
	 * @param string $uuid
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function deleteKeycloakUser(string $uuid): ResponseInterface {
		$url = self::getKeycloakUrl() . '/admin/realms/oCIS/users/' . $uuid;
		return HttpRequestHelper::delete(
			$url,
			null,
			null,
			[
				"Authorization" => "Bearer  " . self::getAdminAccessToken(),
			],
		);
	}
}
