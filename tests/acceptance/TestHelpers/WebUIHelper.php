<?php declare(strict_types=1);

namespace TestHelpers;

require __DIR__ . '/../../../vendor-php/autoload.php';
use Exception;
use Playwright\Playwright;
use Zxing\QrReader;
use OTPHP\TOTP;

class WebUIHelper {
    /**
     * @throws Exception
     */
    public static function setUpUser(string $ocisUrl, string $username, string $password): array
    {
        $context = Playwright::chromium([
            'headless' => false,
            'args' => ['--ignore-certificate-errors', '--no-sandbox']
        ]);
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
            $page->locator("//button[contains(@class, 'oc-topbar-mode-switch-option')][.//span[text()='Vault']]")->click();
            $page->waitForSelector('#kc-totp-secret-qr-code', ['timeout' => 3000]);
            $qrLocator = $page->locator('#kc-totp-secret-qr-code');

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
//            // get the access token
//            $state = $context->storageState();
//            $stateData = \json_decode($state['origins'][0]['localStorage'][2]['value']);
//            return $stateData->access_token;
        } catch (\Exception $e) {
            throw new Exception("Login failed for user '$username': " . $e->getMessage(), 0, $e);
        } finally {
            if (file_exists($screenshotPath)) {
                unlink($screenshotPath);
            }
            $context->close();
        }
    }

    public static function extractOtpFromQr(string $imagePath): string
    {
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
