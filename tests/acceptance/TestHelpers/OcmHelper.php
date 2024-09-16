<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Viktor Scharf <scharf.vi@gmail.com>
 * @copyright Copyright (c) 2024 Viktor Scharf <scharf.vi@gmail.com>
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
 * A helper class for managing federation server requests
 */
class OcmHelper {
	/**
	 * @return string[]
	 */
	private static function getRequestHeaders(): array {
		return [
			'Content-Type' => 'application/json',
		];
	}

	/**
	 * @param string $baseUrl
	 * @param string $path
	 *
	 * @return string
	 */
	public static function getFullUrl(string $baseUrl, string $path): string {
		$fullUrl = $baseUrl;
		if (\substr($fullUrl, -1) !== '/') {
			$fullUrl .= '/';
		}
		$fullUrl .= 'sciencemesh/' . $path;
		return $fullUrl;
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $user
	 * @param string $password
	 * @param string|null $email
	 * @param string|null $description
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function createInvitation(
		string $baseUrl,
		string $xRequestId,
		string $user,
		string $password,
		?string $email = null,
		?string $description = null
	): ResponseInterface {
		$body = [
		  "description" => $description,
		  "recipient" => $email
		];
		$url = self::getFullUrl($baseUrl, 'generate-invite');
		return HttpRequestHelper::post(
			$url,
			$xRequestId,
			$user,
			$password,
			self::getRequestHeaders(),
			\json_encode($body)
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $user
	 * @param string $password
	 * @param string $token
	 * @param string $providerDomain
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function acceptInvitation(
		string $baseUrl,
		string $xRequestId,
		string $user,
		string $password,
		string $token,
		string $providerDomain
	): ResponseInterface {
		$body = [
		  "token" => $token,
		  "providerDomain" => $providerDomain
		];
		$url = self::getFullUrl($baseUrl, 'accept-invite');
		return HttpRequestHelper::post(
			$url,
			$xRequestId,
			$user,
			$password,
			self::getRequestHeaders(),
			\json_encode($body)
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $user
	 * @param string $password
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function findAcceptedUsers(
		string $baseUrl,
		string $xRequestId,
		string $user,
		string $password
	): ResponseInterface {
		$url = self::getFullUrl($baseUrl, 'find-accepted-users');
		return HttpRequestHelper::get(
			$url,
			$xRequestId,
			$user,
			$password,
			self::getRequestHeaders()
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $user
	 * @param string $password
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function listInvite(
		string $baseUrl,
		string $xRequestId,
		string $user,
		string $password
	): ResponseInterface {
		$url = self::getFullUrl($baseUrl, 'list-invite');
		return HttpRequestHelper::get(
			$url,
			$xRequestId,
			$user,
			$password,
			self::getRequestHeaders()
		);
	}
}
