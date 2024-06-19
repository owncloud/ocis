<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Viktor Scharf <scharf.vi@gmail.com>
 * @copyright Copyright (c) 2024 Viktor Scharf <scharf.vi@gmail.com>
 */

namespace TestHelpers;

use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\RequestInterface;
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
		$body['params'] = [
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
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function acceptInvitation(
		string $baseUrl,
		string $xRequestId,
		string $user,
		string $password,
		string $token
	): ResponseInterface {
		$body = [
		  "token" => $token,
		  "providerDomain" => 'ocis-server'
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
}
