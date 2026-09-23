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

require_once __DIR__ . '/../../../vendor-php/autoload.php';

use Exception;
use GuzzleHttp\Exception\GuzzleException;
use InvalidArgumentException;
use JsonException;
use OTPHP\TOTP;
use ParagonIE\ConstantTime\Base32;
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
	// Longest known step sequence in getCode()'s MFA loop, with one spare step of headroom:
	// password POST -> TOTP setup form -> submit -> intermediate redirect -> GET-follow ->
	// OTP challenge form -> submit -> final redirect (7 steps), so 8 never truncates a real
	// flow while still bounding it against an infinite loop.
	private const MAX_MFA_FLOW_STEPS = 8;
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
		$keycloakUrl = \getenv('KC_URL');
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
	 * Reset the cached admin access token
	 *
	 * @return void
	 */
	public static function resetAdminAccessToken(): void {
		self::$adminAccessToken = null;
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
	 * @return array
	 * @throws GuzzleException
	 */
	private static function getAuthorizationHeader(): array {
		return [ 'Authorization' => 'Bearer ' . self::getAdminAccessToken() ];
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
		$url = self::getKeycloakUrl() . '/admin/realms/oCIS/users';

		return HttpRequestHelper::post(
			$url,
			null,
			null,
			array_merge(
				self::getAuthorizationHeader(),
				[ 'Content-Type' => 'application/json' ],
			),
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
			array_merge(
				self::getAuthorizationHeader(),
				[ 'Content-Type' => 'application/json' ],
			),
			json_encode($body, JSON_THROW_ON_ERROR),
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
	public static function unassignRole(
		string $uuid,
		string $role,
	): ResponseInterface {
		$url = self::getKeycloakUrl() . "/admin/realms/oCIS/users/" . $uuid . "/role-mappings/realm";
		$ocisRole = self::OCIS_KEYCLOAK_USER_ROLES[$role]
			?? throw new InvalidArgumentException("Unknown oCIS role: $role");
		$body = [
			self::getRealmRole($ocisRole),
		];
		return HttpRequestHelper::delete(
			$url,
			null,
			null,
			array_merge(
				self::getAuthorizationHeader(),
				[ 'Content-Type' => 'application/json' ],
			),
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
	 * @param string|null $acrValues
	 *
	 * @return array
	 * @throws GuzzleException
	 */
	public static function getAuthorizationEndPoint(?string $acrValues = null): array {
		$loginParams = [
			'client_id' => 'web',
			'redirect_uri' => OcisHelper::getServerUrl() . '/oidc-callback.html',
			'response_mode' => 'query',
			'response_type' => 'code',
			'scope' => 'openid profile email acr',
		];
		if ($acrValues !== null) {
			$loginParams['acr_values'] = $acrValues;
		}
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
	 * @throws Exception
	 */
	public static function getCode(
		array $user,
		string $authorizationUrl,
		string $cookie,
	): string {
		$response = HttpRequestHelper::post(
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

		// Whatever step-up Keycloak asks for next (first-time "configure TOTP", since this
		// user has no credential yet, or a login challenge against one that already exists) is
		// handled here purely over HTTP, mirroring what a browser driving the same flow would do.
		// Keycloak's required-action flow follows Post/Redirect/Get: after accepting a form
		// submission it redirects back to its own login-actions endpoint (to be GET-followed)
		// rather than straight to our redirect_uri, so a Location header alone doesn't mean the
		// flow is finished - only a redirect that actually targets redirect_uri does.
		$redirectUriPrefix = OcisHelper::getServerUrl() . '/oidc-callback.html';
		$stepsTaken = [];
		for ($i = 0; $i < self::MAX_MFA_FLOW_STEPS; $i++) {
			$locationHeader = $response->getHeader('Location');
			if (!empty($locationHeader)) {
				$location = $locationHeader[0];
				if (\str_starts_with($location, $redirectUriPrefix)) {
					return self::extractCodeFromLocationHeader($location);
				}
				$stepsTaken[] = 'redirect: ' . $location;
				$cookie = $response->getHeader('Set-Cookie')[0] ?? $cookie;
				$response = HttpRequestHelper::get($location, null, null, [ 'Cookie' => $cookie ]);
				continue;
			}

			$html = $response->getBody()->getContents();
			// Keycloak may rotate the session cookie at every step of the flow
			$cookie = $response->getHeader('Set-Cookie')[0] ?? $cookie;

			if (\preg_match('/id="kc-totp-settings-form"/i', $html)) {
				$stepsTaken[] = 'kc-totp-settings-form';
				[$formUrl, $fields] = self::parseForm($html, 'kc-totp-settings-form');
				$rawSecret = $fields['totpSecret'] ?? null;
				if ($rawSecret === null) {
					throw new Exception('Could not find the "totpSecret" hidden field on the TOTP setup form.');
				}
				$knownTotpSecret = Base32::encodeUnpadded($rawSecret);
				$fields['totp'] = TOTP::createFromSecret($knownTotpSecret)->now();
				$fields['userLabel'] = 'test';
			} else {
				throw new Exception(
					'Unexpected response after username/password submission - no redirect, and '
					. 'neither a TOTP setup nor a TOTP login form was found. Status: '
					. $response->getStatusCode() . ', body: ' . $html,
				);
			}

			$response = HttpRequestHelper::post($formUrl, null, null, [ 'Cookie' => $cookie ], $fields);
		}
		throw new Exception(
			'Too many MFA steps without reaching a redirect to redirect_uri. Steps taken: '
			. \implode(' -> ', $stepsTaken) . '. Last response status: ' . $response->getStatusCode()
			. ', body: ' . $response->getBody()->getContents(),
		);
	}

	/**
	 * Parses the named form out of the given HTML, returning its action url and every field it
	 * carries (name/value pairs from every <input>/<button>, including hidden ones such as
	 * "credentialId" or "totpSecret") - a real browser submission sends all of them, and
	 * Keycloak's authenticators can silently re-render the same form (200, no redirect) if a
	 * required hidden field is missing from the submission.
	 *
	 * @param string $html
	 * @param string $formId
	 *
	 * @return array
	 * @throws Exception
	 */
	private static function parseForm(string $html, string $formId): array {
		// Capture the opening tag's attributes as one group so "id" and "action" can be pulled
		// out independent of which order Keycloak's template happens to render them in.
		if (!preg_match(
			'/<form\b([^>]*\bid="' . preg_quote($formId, '/') . '"[^>]*)>(.*?)<\/form>/is',
			$html,
			$formMatch,
		)
		) {
			throw new Exception("Could not find the \"$formId\" form in the HTML response body.");
		}
		$openingTagAttrs = $formMatch[1];
		$formBody = $formMatch[2];

		if (!preg_match('/\baction="([^"]+)"/i', $openingTagAttrs, $actionMatch)) {
			throw new Exception("No action url found on the \"$formId\" form.");
		}
		$formUrl = \html_entity_decode($actionMatch[1]);

		$fields = [];
		if (preg_match_all('/<(?:input|button)\b[^>]*>/i', $formBody, $fieldMatches)) {
			foreach ($fieldMatches[0] as $fieldTag) {
				if (!preg_match('/\bname="([^"]+)"/i', $fieldTag, $nameMatch)) {
					continue;
				}
				$value = '';
				if (preg_match('/\bvalue="([^"]*)"/i', $fieldTag, $valueMatch)) {
					$value = \html_entity_decode($valueMatch[1]);
				}
				$fields[$nameMatch[1]] = $value;
			}
		}
		return [$formUrl, $fields];
	}

	/**
	 * @param string $location
	 *
	 * @return string
	 * @throws Exception
	 */
	private static function extractCodeFromLocationHeader(string $location): string {
		$queryString = parse_url($location, PHP_URL_QUERY);
		parse_str((string)$queryString, $urlParams);
		if (!isset($urlParams['code'])) {
			throw new Exception("No 'code' parameter found in redirect location: $location");
		}
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
	 * @param string|null $acrValues
	 *
	 * @return array
	 * @throws GuzzleException
	 * @throws JsonException
	 * @throws Exception
	 */
	public static function setAccessTokenForKeycloakOcisUser(
		array $user,
		?string $acrValues = null,
	): array {
		[$authorizationUrl, $cookie] = self::getAuthorizationEndPoint($acrValues);
		$authorizationCode = self::getCode($user, $authorizationUrl, $cookie);
		$tokenResponse = self::getToken($authorizationCode);
		return json_decode($tokenResponse->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
	}

	/**
	 * @return array
	 * @throws GuzzleException
	 * @throws JsonException
	 * @throws Exception
	 */
	public static function getRealm(): array {
		$url = self::getKeycloakUrl() . '/admin/realms/oCIS';
		$response = HttpRequestHelper::get(
			$url,
			null,
			null,
			self::getAuthorizationHeader(),
		);
		if ($response->getStatusCode() !== 200) {
			throw new Exception("Failed to get realm roles.");
		}
		return json_decode($response->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
	}

	/**
	 * @param string $key
	 * @param string $value
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public static function updateRealmAttribute(string $key, string $value): ResponseInterface {
		$realm = self::getRealm();
		$attributes = $realm['attributes'] ?? [];
		$attributes[$key] = $value;
		$url = self::getKeycloakUrl() . '/admin/realms/oCIS';
		return HttpRequestHelper::put(
			$url,
			null,
			null,
			array_merge(
				self::getAuthorizationHeader(),
				[ 'Content-Type' => 'application/json' ],
			),
			json_encode(['attributes' => $attributes], JSON_THROW_ON_ERROR),
		);
	}

	/**
	 * @param string $key
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public static function deleteRealmAttribute(string $key): ResponseInterface {
		$realm = self::getRealm();
		$attributes = $realm['attributes'] ?? [];
		unset($attributes[$key]);
		$url = self::getKeycloakUrl() . '/admin/realms/oCIS';
		return HttpRequestHelper::put(
			$url,
			null,
			null,
			array_merge(
				self::getAuthorizationHeader(),
				[ 'Content-Type' => 'application/json' ],
			),
			json_encode(['attributes' => $attributes], JSON_THROW_ON_ERROR),
		);
	}

	/**
	 * @param string $username
	 *
	 * @return string
	 * @throws GuzzleException
	 * @throws JsonException
	 * @throws Exception
	 */
	public static function getUserIdByUsername(string $username): string {
		$url = self::getKeycloakUrl() . '/admin/realms/oCIS/users?username=' . \urlencode($username) . '&exact=true';
		$response = HttpRequestHelper::get(
			$url,
			null,
			null,
			self::getAuthorizationHeader(),
		);
		if ($response->getStatusCode() !== 200) {
			throw new Exception("Failed to look up Keycloak user '$username', status: " . $response->getStatusCode());
		}
		$users = json_decode($response->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
		if (empty($users)) {
			throw new Exception("Keycloak user '$username' not found.");
		}
		return $users[0]['id'];
	}

	/**
	 * Deletes all OTP/TOTP credentials for a Keycloak user, so MFA can be set up fresh.
	 *
	 * @param string $username
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 * @throws Exception
	 */
	public static function deleteUserTotpCredentials(string $username): void {
		$uuid = self::getUserIdByUsername($username);
		$url = self::getKeycloakUrl() . '/admin/realms/oCIS/users/' . $uuid . '/credentials';
		$response = HttpRequestHelper::get(
			$url,
			null,
			null,
			self::getAuthorizationHeader(),
		);
		if ($response->getStatusCode() !== 200) {
			throw new Exception("Failed to list credentials for Keycloak user '$username'.");
		}
		$credentials = json_decode($response->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
		foreach ($credentials as $credential) {
			if (isset($credential['type']) && \in_array($credential['type'], ['otp', 'totp'], true)) {
				$deleteUrl = self::getKeycloakUrl()
					. '/admin/realms/oCIS/users/' . $uuid
					. '/credentials/' . $credential['id'];
				HttpRequestHelper::delete(
					$deleteUrl,
					null,
					null,
					self::getAuthorizationHeader(),
				);
			}
		}
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
			self::getAuthorizationHeader(),
		);
	}
}
