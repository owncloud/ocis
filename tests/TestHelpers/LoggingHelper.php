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
use InvalidArgumentException;

/**
 * Helper to read and analyze the owncloud log file
 *
 * @author Artur Neumann <artur@jankaritech.com>
 *
 */
class LoggingHelper {
	/**
	 * @var array
	 */
	public const LOG_LEVEL_ARRAY = [
		"debug",
		"info",
		"warning",
		"error",
		"fatal"
	];

	/**
	 * returns the log file path local to the system ownCloud is running on
	 *
	 * @param string|null $xRequestId
	 *
	 * @return string
	 * @throws Exception
	 */
	public static function getLogFilePath(
		?string $xRequestId = ''
	):string {
		// Currently we don't interact with the log file on reva or OCIS
		return "";
	}

	/**
	 *
	 * @param string|null $baseUrl
	 * @param string|null $adminUsername
	 * @param string|null $adminPassword
	 * @param string|null $xRequestId
	 * @param int|null $noOfLinesToRead
	 *
	 * @return array
	 * @throws Exception
	 */
	public static function getLogFileContent(
		?string $baseUrl,
		?string $adminUsername,
		?string $adminPassword,
		?string $xRequestId = '',
		?int $noOfLinesToRead = 0
	):array {
		$result = OcsApiHelper::sendRequest(
			$baseUrl,
			$adminUsername,
			$adminPassword,
			"GET",
			"/apps/testing/api/v1/logfile/$noOfLinesToRead",
			$xRequestId
		);
		if ($result->getStatusCode() !== 200) {
			throw new Exception(
				"could not get logfile content " . $result->getReasonPhrase()
			);
		}
		$response = HttpRequestHelper::getResponseXml($result, __METHOD__);

		$result = [];
		foreach ($response->data->element as $line) {
			array_push($result, (string)$line);
		}
		return $result;
	}

	/**
	 * returns the currently set log level [debug, info, warning, error, fatal]
	 *
	 * @param string|null $xRequestId
	 *
	 * @return string
	 * @throws Exception
	 */
	public static function getLogLevel(
		?string $xRequestId = ''
	):string {
		return "debug";
	}

	/**
	 *
	 * @param string|null $logLevel (debug|info|warning|error|fatal)
	 * @param string|null $xRequestId
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function setLogLevel(
		?string $logLevel,
		?string $xRequestId = ''
	):void {
		// Currently we can't manage log file settings on reva or OCIS
		return;
	}

	/**
	 * returns the currently set logging backend (owncloud|syslog|errorlog)
	 *
	 * @param string|null $xRequestId
	 *
	 * @return string
	 * @throws Exception
	 */
	public static function getLogBackend(
		?string $xRequestId = ''
	):string {
		return "errorlog";
	}

	/**
	 *
	 * @param string|null $backend (owncloud|syslog|errorlog)
	 * @param string|null $xRequestId
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function setLogBackend(
		?string  $backend,
		?string  $xRequestId = ''
	):void {
		if (!\in_array($backend, ["owncloud", "syslog", "errorlog"])) {
			throw new InvalidArgumentException("invalid log backend");
		}
		// Currently we can't manage log file settings on reva or OCIS
		return;
	}

	/**
	 * returns the currently set logging timezone
	 *
	 * @param string|null $xRequestId
	 *
	 * @return string
	 * @throws Exception
	 */
	public static function getLogTimezone(
		?string  $xRequestId = ''
	):string {
		return "UTC";
	}

	/**
	 *
	 * @param string|null $timezone
	 * @param string|null $xRequestId
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function setLogTimezone(
		?string $timezone,
		?string $xRequestId = ''
	):void {
		// Currently we can't manage log file settings on reva or OCIS
		return;
	}

	/**
	 *
	 * @param string|null $baseUrl
	 * @param string|null $adminUsername
	 * @param string|null $adminPassword
	 * @param string|null $xRequestId
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function clearLogFile(
		?string $baseUrl,
		?string $adminUsername,
		?string $adminPassword,
		?string $xRequestId = ''
	):void {
		// Currently we don't interact with the log file on reva or OCIS
		return;
	}

	/**
	 *
	 * @param string|null $logLevel
	 * @param string|null $backend
	 * @param string|null $timezone
	 * @param string|null $xRequestId
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function restoreLoggingStatus(
		?string $logLevel,
		?string $backend,
		?string $timezone,
		?string $xRequestId = ''
	):void {
		// Currently we don't interact with the log file on reva or OCIS
		return;
	}

	/**
	 * returns the currently set log level, backend and timezone
	 *
	 * @param string|null $xRequestId
	 *
	 * @return array|string[]
	 * @throws Exception
	 */
	public static function getLogInfo(
		?string $xRequestId = ''
	):array {
		return [
			"level" => "debug",
			"backend" => "errorlog",
			"timezone" => "UTC"
		];
	}
}
