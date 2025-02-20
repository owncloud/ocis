<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Amrita Shrestha <amrita@jankaritech.com>
 * @copyright Copyright (c) 2025 Amrita Shrestha amrita@jankaritech.com
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
	public static function getEmailFullUrl(string $path): string {
		return self::getEmailBaseUrl() . '/api/v1/' . $path;
	}

	/**
	 * Returns the host and port where Email messages can be read and deleted
	 * by the test runner.
	 *
	 * @return string
	 */
	public static function getEmailBaseUrl(): string {
		$localEmailHost = self::getLocalEmailHost();
		$emailPort = \getenv('EMAIL_PORT') ?: "9000";
		return "http://$localEmailHost:$emailPort";
	}

	/**
	 * Returns the host name or address of the Email server as seen from the
	 * point of view of the test runner.
	 *
	 * @return string
	 */
	public static function getLocalEmailHost(): string {
		return \getenv('LOCAL_EMAIL_HOST') ?: "127.0.0.1";
	}

	/**
	 * list all email
	 *
	 * @param string|null $xRequestId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function listAllEmails(
		?string $xRequestId,
	): ResponseInterface {
		return HttpRequestHelper::get(
			self::getEmailFullUrl("messages"),
			$xRequestId,
			null,
			null,
			['Content-Type' => 'application/json']
		);
	}

	/**
	 * @param string $id
	 * @param string|null $xRequestId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function getBodyOfAnEmailById(
		string $id,
		?string $xRequestId,
	): ResponseInterface {
		return HttpRequestHelper::getJsonDecodedResponseBodyContent(
			HttpRequestHelper::get(
				self::getEmailFullUrl("message/$id"),
				$xRequestId,
				null,
				null,
				['Content-Type' => 'application/json']
			)
		);
	}

	//    move to context
	/**
	 * Returns the body of the last received email for the provided receiver according to the provided email address and the serial number
	 * For email number, 1 means the latest one
	 *
	 * @param string $emailAddress
	 * @param string|null $xRequestId
	 * @param int|null $emailNumber For email number, 1 means the latest one
	 * @param int|null $waitTimeSec Time to wait for the email if the email has been delivered
	 *
	 * @return string
	 * @throws GuzzleException
	 * @throws Exception
	 */
	public static function getBodyOfLastEmail(
		string $emailAddress,
		string $xRequestId,
		?int $emailNumber = 1,
		?int $waitTimeSec = EMAIL_WAIT_TIMEOUT_SEC
	): string {
		$currentTime = \time();
		$endTime = $currentTime + $waitTimeSec;
		//        $mailBox = self::getMailBoxFromEmail($emailAddress);
		while ($currentTime <= $endTime) {
			$query = 'to:' . $emailAddress;
			$mailResponse = self::searchEmails($query, $xRequestId);
			if (!empty($mailResponse) && \sizeof($mailResponse) >= $emailNumber) {
				$response = self::getBodyOfAnEmailById("latest", $xRequestId);
				$body = \str_replace(
					"\r\n",
					"\n",
					\quoted_printable_decode($response->Text . "\n" . $response->HTML)
				);
				return $body;
			}
			\usleep(STANDARD_SLEEP_TIME_MICROSEC * 50);
			$currentTime = \time();
		}
		throw new Exception("Could not find the email to the address: " . $emailAddress);
	}

	/**
	 * search email
	 *
	 * @param string $query
	 * @param string|null $xRequestId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function searchEmails(
		string $query,
		?string $xRequestId,
	): ResponseInterface {
		$url = self::getEmailFullUrl("search") . "?query=$query";
		return HttpRequestHelper::get(
			$url,
			$xRequestId,
			null,
			null,
			['Content-Type' => 'application/json']
		);
	}

	/**
	 * Deletes all email
	 *
	 * @param string|null $xRequestId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function deleteAllEmails(
		?string $xRequestId,
	): ResponseInterface {
		return HttpRequestHelper::delete(
			self::getEmailFullUrl("messages"),
			$xRequestId,
			null,
			null,
			['Content-Type' => 'application/json']
		);
	}
}
