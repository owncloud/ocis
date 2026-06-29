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

require __DIR__ . '/../../../vendor-php/autoload.php';
use Exception;
use Playwright\Playwright;
use Zxing\QrReader;
use OTPHP\TOTP;

/**
 * A helper class for setting up user and store user access token using web UI
 */
class WebUIHelper {
	/**
	 * @param string $ocisUrl
	 * @param string $username
	 * @param string $password
	 *
	 * @return array
	 * @throws Exception
	 */
	public static function setUpUser(string $ocisUrl, string $username, string $password): array {
		$context = Playwright::chromium(
			[
				'headless' => true,
				'args' => ['--ignore-certificate-errors', '--no-sandbox'],
			],
		);
		$screenshotPath = '/tmp/qr_' . uniqid() . '.png';
		try {
			$page = $context->newPage();
			$page->goto($ocisUrl, ['waitUntil' => 'networkidle']);
			$page->waitForSelector('#kc-header', ['timeout' => 3000]);
			$page->locator('#username')->fill($username);
			$page->locator('#password')->fill($password);
			$page->locator('#kc-login')->click();
			$page->waitForSelector('#files-view', ['timeout' => 30000]);

			// change to vault mode
			$page->locator('#oc-topbar-mode-switch-btn')->click();
			$page->locator(
				"//button[contains(@class, 'oc-topbar-mode-switch-option')]" .
				"[.//span[text()='Vault']]",
			)->click();
			$page->waitForSelector('#kc-totp-secret-qr-code', ['timeout' => 3000]);
			$qrLocator = $page->locator('#kc-totp-secret-qr-code');

			// setup mfa
			$qrLocator->screenshot($screenshotPath);
			if (!file_exists($screenshotPath)) {
				throw new Exception("Failed to save QR code screenshot to: " . $screenshotPath);
			}
			$otp = self::extractOtpFromQr($screenshotPath);
			$page->locator('#totp')->fill((string)$otp);
			$page->locator('#userLabel')->fill('test');
			$page->locator('#saveTOTPBtn')->click();
			$page->waitForSelector('#files-view', ['timeout' => 3000]);
			return $context->storageState();
		} catch (\Exception $e) {
			throw new Exception("Login failed for user '$username': " . $e->getMessage(), 0, $e);
		} finally {
			if (file_exists($screenshotPath)) {
				unlink($screenshotPath);
			}
			$context->close();
		}
	}

	/**
	 * @param string $imagePath
	 *
	 * @return string
	 * @throws Exception
	 */
	public static function extractOtpFromQr(string $imagePath): string {
		$qrReader = new QrReader($imagePath);
		$qrData = $qrReader->text();

		if (empty($qrData)) {
			throw new Exception("Could not decode QR code from image: " . $imagePath);
		}

		$parsedUrl = parse_url($qrData);
		if (!isset($parsedUrl['query'])) {
			throw new Exception("QR code data does not contain a valid otpauth URL: " . $qrData);
		}

		parse_str($parsedUrl['query'], $queryParams);
		$secret = $queryParams['secret'] ?? null;

		if ($secret === null) {
			throw new Exception("No 'secret' parameter found in QR code URL: " . $qrData);
		}

		$totp = TOTP::createFromSecret($secret);

		return $totp->now();
	}
}
