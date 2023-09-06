<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Prajwol Amatya <prajwol@jankaritech.com>
 * @copyright Copyright (c) 2023 Prajwol Amatya prajwol@jankaritech.com
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
	 * @param string $emailAddress
	 *
	 * @return string
	 */
	public static function getMailBoxFromEmail(string $emailAddress):string {
		return explode("@", $emailAddress)[0];
	}

	/**
	 * Returns the host and port where Email messages can be read and deleted
	 * by the test runner.
	 *
	 * @return string
	 */
	public static function getLocalEmailUrl():string {
		$localEmailHost = self::getLocalEmailHost();
		$emailPort = \getenv('EMAIL_PORT');
		if ($emailPort === false) {
			$emailPort = "9000";
		}
		return "http://$localEmailHost:$emailPort";
	}

	/**
	 * Returns the host name or address of the Email server as seen from the
	 * point of view of the system-under-test.
	 *
	 * @return string
	 */
	public static function getEmailHost():string {
		$emailHost = \getenv('EMAIL_HOST');
		if ($emailHost === false) {
			$emailHost = "127.0.0.1";
		}
		return $emailHost;
	}

	/**
	 * Returns the host name or address of the Email server as seen from the
	 * point of view of the test runner.
	 *
	 * @return string
	 */
	public static function getLocalEmailHost():string {
		$localEmailHost = \getenv('LOCAL_EMAIL_HOST');
		if ($localEmailHost === false) {
			$localEmailHost = self::getEmailHost();
		}
		return $localEmailHost;
	}

	/**
	 * Returns general response information about the provided mailbox
	 * A mailbox is created automatically in InBucket for every unique email sender|receiver
	 *
	 * @param string $mailBox
	 * @param string|null $xRequestId
	 *
	 * @return array
	 * @throws GuzzleException
	 */
	public static function getMailBoxInformation(string $mailBox, ?string $xRequestId = null):array {
		$response = HttpRequestHelper::get(
			self::getLocalEmailUrl() . "/api/v1/mailbox/" . $mailBox,
			$xRequestId,
			null,
			null,
			['Content-Type' => 'application/json']
		);
		return \json_decode($response->getBody()->getContents());
	}

	/**
	 * returns body content of a specific email (mailBox) with email ID (mailbox Id)
	 *
	 * @param string $mailBox
	 * @param string $mailboxId
	 * @param string|null $xRequestId
	 *
	 * @return object
	 * @throws GuzzleException
	 */
	public static function getBodyOfAnEmailById(string $mailBox, string $mailboxId, ?string $xRequestId = null):object {
		$response = HttpRequestHelper::get(
			self::getLocalEmailUrl() . "/api/v1/mailbox/" . $mailBox . "/" . $mailboxId,
			$xRequestId,
			null,
			null,
			['Content-Type' => 'application/json']
		);
		return \json_decode($response->getBody()->getContents());
	}

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
	):string {
		$currentTime = \time();
		$endTime = $currentTime + $waitTimeSec;
		$mailBox = self::getMailBoxFromEmail($emailAddress);
		while ($currentTime <= $endTime) {
			$mailboxResponse = self::getMailboxInformation($mailBox, $xRequestId);
			if (!empty($mailboxResponse) && \sizeof($mailboxResponse) >= $emailNumber) {
				$mailboxId = $mailboxResponse[\sizeof($mailboxResponse) - $emailNumber]->id;
				$response = self::getBodyOfAnEmailById($mailBox, $mailboxId, $xRequestId);
				$body = \str_replace(
					"\r\n",
					"\n",
					\quoted_printable_decode($response->body->text . "\n" . $response->body->html)
				);
				return $body;
			}
			\usleep(STANDARD_SLEEP_TIME_MICROSEC * 50);
			$currentTime = \time();
		}
		throw new Exception("Could not find the email to the address: " . $emailAddress);
	}

	/**
	 * Deletes all the emails for the provided mailbox
	 *
	 * @param string $localInbucketUrl
	 * @param string|null $xRequestId
	 * @param string $mailBox
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function deleteAllEmailsForAMailbox(
		string $localInbucketUrl,
		?string $xRequestId,
		string $mailBox
	):ResponseInterface {
		return HttpRequestHelper::delete(
			$localInbucketUrl . "/api/v1/mailbox/" . $mailBox,
			$xRequestId
		);
	}
}
