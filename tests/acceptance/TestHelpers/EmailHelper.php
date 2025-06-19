<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Prajwol Amatya <prajwol@jankaritech.com>
 * @author Amrita Shrestha <amrita@jankaritech.com>
 * @copyright Copyright (c) 2023 JankariTech
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
 * A helper class for managing emails
 */
class EmailHelper {
	/**
	 * @param string $path
	 *
	 * @return string
	 */
	public static function getEmailAPIUrl(string $path): string {
		return self::getEmailBaseUrl() . '/api/v1/' . $path;
	}

	/**
	 * Returns the host and port where Email messages can be read and deleted
	 * by the test runner.
	 *
	 * @return string
	 */
	public static function getEmailBaseUrl(): string {
		$emailHost = self::getEmailHost();
		$emailPort = \getenv('EMAIL_PORT') ?: "8025";
		return "http://$emailHost:$emailPort";
	}

	/**
	 * Returns the host name or address of the Email server as seen from the
	 * point of view of the test runner.
	 *
	 * @return string
	 */
	public static function getEmailHost(): string {
		return \getenv('LOCAL_EMAIL_HOST') ?? \getenv('EMAIL_HOST') ?? "127.0.0.1";
	}

	/**
	 * list all email
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function listAllEmails(): ResponseInterface {
		return HttpRequestHelper::get(
			self::getEmailAPIUrl("messages"),
			null,
			null,
			['Content-Type' => 'application/json'],
		);
	}

	/**
	 * @param string $id when $id set to 'latest' returns the latest message.
	 * @param string $query
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function getEmailById(
		string $id,
		string $query,
	): ResponseInterface {
		return HttpRequestHelper::get(
			self::getEmailAPIUrl("message/$id") . "?query=$query",
			null,
			null,
			['Content-Type' => 'application/json'],
		);
	}

	/**
	 * search email
	 *
	 * @param string $query
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function searchEmails(
		string $query,
	): ResponseInterface {
		$url = self::getEmailAPIUrl("search") . "?query=$query";
		return HttpRequestHelper::get(
			$url,
			null,
			null,
			['Content-Type' => 'application/json'],
		);
	}

	/**
	 * Deletes all email
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function deleteAllEmails(): ResponseInterface {
		return HttpRequestHelper::delete(
			self::getEmailAPIUrl("messages"),
			null,
			null,
			['Content-Type' => 'application/json'],
		);
	}
}
