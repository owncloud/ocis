<?php declare(strict_types=1);

/**
 * ownCloud
 *
 * @author Sergio Bertolin <sbertolin@owncloud.com>
 * @author Phillip Davis <phil@jankaritech.com>
 * @copyright Copyright (c) 2018, ownCloud GmbH
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

use Behat\Behat\Hook\Scope\BeforeStepScope;
use GuzzleHttp\Exception\GuzzleException;
use Helmich\JsonAssert\JsonAssertions;
use rdx\behatvars\BehatVariablesContext;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Behat\Behat\Hook\Scope\AfterScenarioScope;
use Behat\Gherkin\Node\PyStringNode;
use Behat\Gherkin\Node\TableNode;
use Behat\Testwork\Hook\Scope\BeforeSuiteScope;
use Behat\Testwork\Hook\Scope\AfterSuiteScope;
use GuzzleHttp\Cookie\CookieJar;
use Psr\Http\Message\ResponseInterface;
use PHPUnit\Framework\Assert;
use TestHelpers\AppConfigHelper;
use TestHelpers\OcsApiHelper;
use TestHelpers\SetupHelper;
use TestHelpers\HttpRequestHelper;
use TestHelpers\HttpLogger;
use TestHelpers\UploadHelper;
use TestHelpers\OcisHelper;
use Laminas\Ldap\Ldap;
use TestHelpers\GraphHelper;
use TestHelpers\WebDavHelper;

require_once 'bootstrap.php';

/**
 * Features context.
 */
class FeatureContext extends BehatVariablesContext {
	use Provisioning;
	use Sharing;
	use WebDav;
	use JsonAssertions;

	/**
	 * Unix timestamp seconds
	 */
	private int $scenarioStartTime;
	private string $adminUsername = '';
	private string $adminPassword = '';
	private string $adminDisplayName = '';
	private string $adminEmailAddress = '';
	private string $originalAdminPassword = '';

	/**
	 * An array of values of replacement values of user attributes.
	 * These are only referenced when creating a user. After that, the
	 * run-time values are maintained and referenced in the $createdUsers array.
	 *
	 * Key is the username, value is an array of user attributes
	 */
	private ?array $userReplacements = null;
	private string $regularUserPassword = '';
	private string $alt1UserPassword = '';
	private string $alt2UserPassword = '';
	private string $alt3UserPassword = '';
	private string $alt4UserPassword = '';

	/**
	 * The password to use in tests that create a sub-admin user
	 */
	private string $subAdminPassword = '';

	/**
	 * The password to use in tests that create another admin user
	 */
	private string $alternateAdminPassword = '';

	/**
	 * The password to use in tests that create public link shares
	 */
	private string $publicLinkSharePassword = '';
	private string $ocPath = '';

	/**
	 * Location of the root folder of ownCloud on the local server under test
	 */
	private ?string $localServerRoot = null;
	private string $currentUser = '';
	private string $currentServer = '';

	/**
	 * The base URL of the current server under test,
	 * without any terminating slash
	 * e.g. http://localhost:8080
	 */
	private string $baseUrl = '';

	/**
	 * The base URL of the local server under test,
	 * without any terminating slash
	 * e.g. http://localhost:8080
	 */
	private string $localBaseUrl = '';

	/**
	 * The base URL of the remote (federated) server under test,
	 * without any terminating slash
	 * e.g. http://localhost:8180
	 */
	private string $remoteBaseUrl = '';

	/**
	 * The suite name, feature name and scenario line number.
	 * Example: apiComments/createComments.feature:24
	 */
	private string $scenarioString = '';

	/**
	 * A full unique reference to the step that is currently executing.
	 * Example: apiComments/createComments.feature:24-28
	 * That is line 28, in the scenario at line 24, in the createComments feature
	 * in the apiComments suite.
	 */
	private string $stepLineRef = '';
	private bool $sendStepLineRef = false;
	private bool $sendStepLineRefHasBeenChecked = false;

	/**
	 * @var boolean true if TEST_SERVER_FED_URL is defined
	 */
	private bool $federatedServerExists = false;
	private int $ocsApiVersion = 1;
	private ?ResponseInterface $response = null;
	private string $responseUser = '';
	private ?string $responseBodyContent = null;
	private array $userResponseBodyContents = [];
	public array $emailRecipients = [];
	private CookieJar $cookieJar;
	private string $requestToken;
	private array $storageIds = [];
	private array $createdFiles = [];

	/**
	 * The local source IP address from which to initiate API actions.
	 * Defaults to system-selected address matching IP address family and scope.
	 */
	private ?string $sourceIpAddress = null;
	private array $guzzleClientHeaders = [];
	public OCSContext $ocsContext;
	public AuthContext $authContext;
	public GraphContext $graphContext;
	public SpacesContext $spacesContext;
	private array $initialTrustedServer;

	/**
	 * The codes are stored as strings, even though they are numbers
	 */
	private array $lastHttpStatusCodesArray = [];
	private array $lastOCSStatusCodesArray = [];

	/**
	 * this is set true for db conversion tests
	 */
	private bool $dbConversion = false;

	/**
	 * @param bool $value
	 *
	 * @return void
	 */
	public function setDbConversionState(bool $value): void {
		$this->dbConversion = $value;
	}

	/**
	 * @return bool
	 */
	public function isRunningForDbConversion(): bool {
		return $this->dbConversion;
	}

	private string $oCSelector;

	/**
	 * @param string $selector
	 *
	 * @return void
	 */
	public function setOCSelector(string $selector): void {
		$this->oCSelector = $selector;
	}

	/**
	 * @return string
	 */
	public function getOCSelector(): string {
		return $this->oCSelector;
	}

	/**
	 * @param string|null $httpStatusCode
	 *
	 * @return void
	 */
	public function pushToLastHttpStatusCodesArray(?string $httpStatusCode = null): void {
		if ($httpStatusCode !== null) {
			$this->lastHttpStatusCodesArray[] = $httpStatusCode;
		} elseif ($this->getResponse()->getStatusCode() !== null) {
			$this->lastHttpStatusCodesArray[] = (string)$this->getResponse()->getStatusCode();
		}
	}

	/**
	 * @return void
	 */
	public function emptyLastHTTPStatusCodesArray(): void {
		$this->lastHttpStatusCodesArray = [];
	}

	/**
	 * @return void
	 */
	public function emptyLastOCSStatusCodesArray(): void {
		$this->lastOCSStatusCodesArray = [];
	}

	/**
	 * @return void
	 */
	public function clearStatusCodeArrays(): void {
		$this->emptyLastHTTPStatusCodesArray();
		$this->emptyLastOCSStatusCodesArray();
	}

	/**
	 * @param string $ocsStatusCode
	 *
	 * @return void
	 */
	public function pushToLastOcsCodesArray(string $ocsStatusCode): void {
		$this->lastOCSStatusCodesArray[] = $ocsStatusCode;
	}

	/**
	 * Add HTTP and OCS status code of the last response to the respective status code array
	 *
	 * @return void
	 */
	public function pushToLastStatusCodesArrays(): void {
		$this->pushToLastHttpStatusCodesArray(
			(string)$this->getResponse()->getStatusCode()
		);
		try {
			$this->pushToLastOcsCodesArray(
				$this->ocsContext->getOCSResponseStatusCode(
					$this->getResponse()
				)
			);
		} catch (Exception $exception) {
			// if response couldn't be converted into xml then push "notset" to last ocs status codes array
			$this->pushToLastOcsCodesArray("notset");
		}
	}

	/**
	 * @param string $emailAddress
	 *
	 * @return void
	 */
	public function pushEmailRecipientAsMailBox(string $emailAddress): void {
		$mailBox = explode("@", $emailAddress)[0];
		if (!\in_array($mailBox, $this->emailRecipients)) {
			$this->emailRecipients[] = $mailBox;
		}
	}

	private Ldap $ldap;
	private string $ldapBaseDN;
	private string $ldapHost;
	private int $ldapPort;
	private string $ldapAdminUser;
	private string $ldapAdminPassword = "";
	private string $ldapUsersOU;
	private string $ldapGroupsOU;
	private string $ldapGroupSchema;
	private bool $skipImportLdif;
	private array $toDeleteDNs = [];
	private array $ldapCreatedUsers = [];
	private array $ldapCreatedGroups = [];
	private array $toDeleteLdapConfigs = [];
	private array $oldLdapConfig = [];

	/**
	 * @return Ldap
	 */
	public function getLdap(): Ldap {
		return $this->ldap;
	}

	/**
	 * @param string $configId
	 *
	 * @return void
	 */
	public function setToDeleteLdapConfigs(string $configId): void {
		$this->toDeleteLdapConfigs[] = $configId;
	}

	/**
	 * @return array
	 */
	public function getToDeleteLdapConfigs(): array {
		return $this->toDeleteLdapConfigs;
	}

	/**
	 * @param string $setValue
	 *
	 * @return void
	 */
	public function setToDeleteDNs(string $setValue): void {
		$this->toDeleteDNs[] = $setValue;
	}

	/**
	 * @return string
	 */
	public function getLdapBaseDN(): string {
		return $this->ldapBaseDN;
	}

	/**
	 * @return string
	 */
	public function getLdapUsersOU(): string {
		return $this->ldapUsersOU;
	}

	/**
	 * @return string
	 */
	public function getLdapGroupsOU(): string {
		return $this->ldapGroupsOU;
	}

	/**
	 * @return array
	 */
	public function getOldLdapConfig(): array {
		return $this->oldLdapConfig;
	}

	/**
	 * @param string $configId
	 * @param string $configKey
	 * @param string $value
	 *
	 * @return void
	 */
	public function setOldLdapConfig(string $configId, string $configKey, string $value): void {
		$this->oldLdapConfig[$configId][$configKey] = $value;
	}

	/**
	 * @return string
	 */
	public function getLdapHost(): string {
		return $this->ldapHost;
	}

	/**
	 * @return string
	 */
	public function getLdapHostWithoutScheme(): string {
		return $this->removeSchemeFromUrl($this->ldapHost);
	}

	/**
	 * @return integer
	 */
	public function getLdapPort(): int {
		return $this->ldapPort;
	}

	/**
	 * @return bool
	 */
	public function isTestingWithLdap(): bool {
		return (\getenv("TEST_WITH_LDAP") === "true");
	}

	/**
	 * @return bool
	 */
	public function sendScenarioLineReferencesInXRequestId(): ?bool {
		if ($this->sendStepLineRefHasBeenChecked === false) {
			$this->sendStepLineRef = (\getenv("SEND_SCENARIO_LINE_REFERENCES") === "true");
			$this->sendStepLineRefHasBeenChecked = true;
		}
		return $this->sendStepLineRef;
	}

	/**
	 * @return bool
	 */
	public function isTestingReplacingUsernames(): bool {
		return (\getenv('REPLACE_USERNAMES') === "true");
	}

	/**
	 * @return array|null
	 */
	public function usersToBeReplaced(): ?array {
		if (($this->userReplacements === null) && $this->isTestingReplacingUsernames()) {
			$this->userReplacements = \json_decode(
				\file_get_contents("./tests/acceptance/usernames.json"),
				true
			);
			// Loop through the user replacements, and make entries for the lower
			// and upper case forms. This allows for steps that specifically
			// want to test that usernames like "alice", "Alice" and "ALICE" all work.
			// Such steps will make useful replacements for each form.
			foreach ($this->userReplacements as $key => $value) {
				$lowerKey = \strtolower($key);
				if ($lowerKey !== $key) {
					$this->userReplacements[$lowerKey] = $value;
					$this->userReplacements[$lowerKey]['username'] = \strtolower(
						$this->userReplacements[$lowerKey]['username']
					);
				}
				$upperKey = \strtoupper($key);
				if ($upperKey !== $key) {
					$this->userReplacements[$upperKey] = $value;
					$this->userReplacements[$upperKey]['username'] = \strtoupper(
						$this->userReplacements[$upperKey]['username']
					);
				}
			}
		}
		return $this->userReplacements;
	}

	/**
	 * BasicStructure constructor.
	 *
	 * @param string $baseUrl
	 * @param string $adminUsername
	 * @param string $adminPassword
	 * @param string $regularUserPassword
	 * @param string $ocPath
	 *
	 */
	public function __construct(
		string $baseUrl,
		string $adminUsername,
		string $adminPassword,
		string $regularUserPassword,
		string $ocPath
	) {
		// Initialize your context here
		$this->baseUrl = \rtrim($baseUrl, '/');
		$this->adminUsername = $adminUsername;
		$this->adminPassword = $adminPassword;
		$this->regularUserPassword = $regularUserPassword;
		$this->localBaseUrl = $this->baseUrl;
		$this->currentServer = 'LOCAL';
		$this->cookieJar = new CookieJar();
		$this->ocPath = $ocPath;

		// PARALLEL DEPLOYMENT: ownCloud selector
		$this->oCSelector = "oc10";

		// These passwords are referenced in tests and can be overridden by
		// setting environment variables.
		$this->alt1UserPassword = "1234";
		$this->alt2UserPassword = "AaBb2Cc3Dd4";
		$this->alt3UserPassword = "aVeryLongPassword42TheMeaningOfLife";
		$this->alt4UserPassword = "ThisIsThe4thAlternatePwd";
		$this->subAdminPassword = "IamAJuniorAdmin42";
		$this->alternateAdminPassword = "IHave99LotsOfPriv";
		$this->publicLinkSharePassword = "publicPwd:1";

		// in case of CI deployment we take the server url from the environment
		$testServerUrl = \getenv('TEST_SERVER_URL');
		if ($testServerUrl !== false) {
			$this->baseUrl = \rtrim($testServerUrl, '/');
			$this->localBaseUrl = $this->baseUrl;
		}

		// federated server url from the environment
		$testRemoteServerUrl = \getenv('TEST_SERVER_FED_URL');
		if ($testRemoteServerUrl !== false) {
			$this->remoteBaseUrl = \rtrim($testRemoteServerUrl, '/');
			$this->federatedServerExists = true;
		} else {
			$this->remoteBaseUrl = $this->localBaseUrl;
			$this->federatedServerExists = false;
		}

		// get the admin username from the environment (if defined)
		$adminUsernameFromEnvironment = $this->getAdminUsernameFromEnvironment();
		if ($adminUsernameFromEnvironment !== false) {
			$this->adminUsername = $adminUsernameFromEnvironment;
		}

		// get the admin password from the environment (if defined)
		$adminPasswordFromEnvironment = $this->getAdminPasswordFromEnvironment();
		if ($adminPasswordFromEnvironment !== false) {
			$this->adminPassword = $adminPasswordFromEnvironment;
		}

		// get the regular user password from the environment (if defined)
		$regularUserPasswordFromEnvironment = $this->getRegularUserPasswordFromEnvironment();
		if ($regularUserPasswordFromEnvironment !== false) {
			$this->regularUserPassword = $regularUserPasswordFromEnvironment;
		}

		// get the alternate(1) user password from the environment (if defined)
		$alt1UserPasswordFromEnvironment = $this->getAlt1UserPasswordFromEnvironment();
		if ($alt1UserPasswordFromEnvironment !== false) {
			$this->alt1UserPassword = $alt1UserPasswordFromEnvironment;
		}

		// get the alternate(2) user password from the environment (if defined)
		$alt2UserPasswordFromEnvironment = $this->getAlt2UserPasswordFromEnvironment();
		if ($alt2UserPasswordFromEnvironment !== false) {
			$this->alt2UserPassword = $alt2UserPasswordFromEnvironment;
		}

		// get the alternate(3) user password from the environment (if defined)
		$alt3UserPasswordFromEnvironment = $this->getAlt3UserPasswordFromEnvironment();
		if ($alt3UserPasswordFromEnvironment !== false) {
			$this->alt3UserPassword = $alt3UserPasswordFromEnvironment;
		}

		// get the alternate(4) user password from the environment (if defined)
		$alt4UserPasswordFromEnvironment = $this->getAlt4UserPasswordFromEnvironment();
		if ($alt4UserPasswordFromEnvironment !== false) {
			$this->alt4UserPassword = $alt4UserPasswordFromEnvironment;
		}

		// get the sub-admin password from the environment (if defined)
		$subAdminPasswordFromEnvironment = $this->getSubAdminPasswordFromEnvironment();
		if ($subAdminPasswordFromEnvironment !== false) {
			$this->subAdminPassword = $subAdminPasswordFromEnvironment;
		}

		// get the alternate admin password from the environment (if defined)
		$alternateAdminPasswordFromEnvironment = $this->getAlternateAdminPasswordFromEnvironment();
		if ($alternateAdminPasswordFromEnvironment !== false) {
			$this->alternateAdminPassword = $alternateAdminPasswordFromEnvironment;
		}

		// get the public link share password from the environment (if defined)
		$publicLinkSharePasswordFromEnvironment = $this->getPublicLinkSharePasswordFromEnvironment();
		if ($publicLinkSharePasswordFromEnvironment !== false) {
			$this->publicLinkSharePassword = $publicLinkSharePasswordFromEnvironment;
		}
		$this->originalAdminPassword = $this->adminPassword;
	}

	/**
	 * Create log directory if it doesn't exist
	 *
	 * @BeforeSuite
	 *
	 * @param BeforeSuiteScope $scope
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function setupLogDir(BeforeSuiteScope $scope): void {
		if (!\file_exists(HttpLogger::getLogDir())) {
			\mkdir(HttpLogger::getLogDir(), 0777, true);
		}
	}

	/**
	 *
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function logScenario(BeforeScenarioScope $scope): void {
		$scenarioLine = self::getScenarioLine($scope);

		if ($scope->getScenario()->getNodeType() === "Example") {
			$scenario = "Scenario Outline: " . $scope->getScenario()->getOutlineTitle();
		} else {
			$scenario = $scope->getScenario()->getNodeType() . ": " . $scope->getScenario()->getTitle();
		}

		$logMessage = "## $scenario ($scenarioLine)\n";

		// Delete previous scenario's log file
		if (\file_exists(HttpLogger::getScenarioLogPath())) {
			\unlink(HttpLogger::getScenarioLogPath());
		}

		// Write the scenario log
		HttpLogger::writeLog(HttpLogger::getScenarioLogPath(), $logMessage);
	}

	/**
	 *
	 * @BeforeStep
	 *
	 * @param BeforeStepScope $scope
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function logStep(BeforeStepScope $scope): void {
		$step = $scope->getStep()->getKeyword() . " " . $scope->getStep()->getText();
		$logMessage = "\t### $step\n";
		HttpLogger::writeLog(HttpLogger::getScenarioLogPath(), $logMessage);
	}

	/**
	 * @param string $appTestCodeFullPath
	 *
	 * @return string the relative path from the core tests/acceptance dir
	 *                to the equivalent dir in the app
	 */
	public function getPathFromCoreToAppAcceptanceTests(
		string $appTestCodeFullPath
	): string {
		// $appTestCodeFullPath is something like:
		// '/somedir/anotherdir/core/apps/guests/tests/acceptance/features/bootstrap'
		// and we want to know the 'apps/guests/tests/acceptance' part

		$path = \dirname($appTestCodeFullPath, 2);
		$acceptanceDir = \basename($path);
		$path = \dirname($path);
		$testsDir = \basename($path);
		$path = \dirname($path);
		$appNameDir = \basename($path);
		$path = \dirname($path);
		// We specially are not sure about the name of the directory 'apps'
		// Sometimes the app could be installed in some alternate apps directory
		// like, for example, `apps-external`. So this really does need to be
		// resolved here at run-time.
		$appsDir = \basename($path);
		// To get from core tests/acceptance we go up 2 levels then down through
		// the above app dirs.
		return "../../$appsDir/$appNameDir/$testsDir/$acceptanceDir";
	}

	/**
	 * Get the externally-defined admin username, if any
	 *
	 * @return string|false
	 */
	private static function getAdminUsernameFromEnvironment() {
		return \getenv('ADMIN_USERNAME');
	}

	/**
	 * Get the externally-defined admin password, if any
	 *
	 * @return string|false
	 */
	private static function getAdminPasswordFromEnvironment() {
		return \getenv('ADMIN_PASSWORD');
	}

	/**
	 * Get the externally-defined regular user password, if any
	 *
	 * @return string|false
	 */
	private static function getRegularUserPasswordFromEnvironment() {
		return \getenv('REGULAR_USER_PASSWORD');
	}

	/**
	 * Get the externally-defined alternate(1) user password, if any
	 *
	 * @return string|false
	 */
	private static function getAlt1UserPasswordFromEnvironment() {
		return \getenv('ALT1_USER_PASSWORD');
	}

	/**
	 * Get the externally-defined alternate(2) user password, if any
	 *
	 * @return string|false
	 */
	private static function getAlt2UserPasswordFromEnvironment() {
		return \getenv('ALT2_USER_PASSWORD');
	}

	/**
	 * Get the externally-defined alternate(3) user password, if any
	 *
	 * @return string|false
	 */
	private static function getAlt3UserPasswordFromEnvironment() {
		return \getenv('ALT3_USER_PASSWORD');
	}

	/**
	 * Get the externally-defined alternate(4) user password, if any
	 *
	 * @return string|false
	 */
	private static function getAlt4UserPasswordFromEnvironment() {
		return \getenv('ALT4_USER_PASSWORD');
	}

	/**
	 * Get the externally-defined sub-admin password, if any
	 *
	 * @return string|false
	 */
	private static function getSubAdminPasswordFromEnvironment() {
		return \getenv('SUB_ADMIN_PASSWORD');
	}

	/**
	 * Get the externally-defined alternate admin password, if any
	 *
	 * @return string|false
	 */
	private static function getAlternateAdminPasswordFromEnvironment() {
		return \getenv('ALTERNATE_ADMIN_PASSWORD');
	}

	/**
	 * Get the externally-defined public link share password, if any
	 *
	 * @return string|false
	 */
	private static function getPublicLinkSharePasswordFromEnvironment() {
		return \getenv('PUBLIC_LINK_SHARE_PASSWORD');
	}

	/**
	 * removes the scheme "http(s)://" (if any) from the front of a URL
	 * note: only needs to handle http or https
	 *
	 * @param string $url
	 *
	 * @return string
	 */
	public function removeSchemeFromUrl(string $url): string {
		return \preg_replace(
			"(^https?://)",
			"",
			$url
		);
	}

	/**
	 * @return string
	 */
	public function getOcPath(): string {
		return $this->ocPath;
	}

	/**
	 * @return CookieJar
	 */
	public function getCookieJar(): CookieJar {
		return $this->cookieJar;
	}

	/**
	 * @return string
	 */
	public function getRequestToken(): string {
		return $this->requestToken;
	}

	/**
	 * returns the base URL (which is without a slash at the end)
	 *
	 * @return string
	 */
	public function getBaseUrl(): string {
		return $this->baseUrl;
	}

	/**
	 * returns the path of the base URL
	 * e.g. owncloud-core/10 if the baseUrl is http://localhost/owncloud-core/10
	 * the path is without a slash at the end and without a slash at the beginning
	 *
	 * @return string
	 */
	public function getBasePath(): string {
		$parsedUrl = \parse_url($this->getBaseUrl(), PHP_URL_PATH);
		// If the server-under-test is at the "top" of the domain then parse_url returns null.
		// For example, testing a server at http://localhost:8080 or http://example.com
		if ($parsedUrl === null) {
			$parsedUrl = '';
		}
		return \ltrim($parsedUrl, "/");
	}

	/**
	 * returns the OCS path
	 * the path is without a slash at the end and without a slash at the beginning
	 *
	 * @param string $ocsApiVersion
	 *
	 * @return string
	 */
	public function getOCSPath(string $ocsApiVersion): string {
		return \ltrim($this->getBasePath() . "/ocs/v$ocsApiVersion.php", "/");
	}

	/**
	 * returns the complete DAV path including the base path e.g. owncloud-core/remote.php/dav
	 *
	 * @return string
	 */
	public function getDAVPathIncludingBasePath(): string {
		return \ltrim($this->getBasePath() . "/" . $this->getDavPath(), "/");
	}

	/**
	 * returns the base URL but without "http(s)://" in front of it
	 *
	 * @return string
	 */
	public function getBaseUrlWithoutScheme(): string {
		return $this->removeSchemeFromUrl($this->getBaseUrl());
	}

	/**
	 * returns the local base URL (which is without a slash at the end)
	 *
	 * @return string
	 */
	public function getLocalBaseUrl(): string {
		return $this->localBaseUrl;
	}

	/**
	 * returns the local base URL but without "http(s)://" in front of it
	 *
	 * @return string
	 */
	public function getLocalBaseUrlWithoutScheme(): string {
		return $this->removeSchemeFromUrl($this->getLocalBaseUrl());
	}

	/**
	 * returns the remote base URL (which is without a slash at the end)
	 *
	 * @return string
	 */
	public function getRemoteBaseUrl(): string {
		return $this->remoteBaseUrl;
	}

	/**
	 * returns the remote base URL but without "http(s)://" in front of it
	 *
	 * @return string
	 */
	public function getRemoteBaseUrlWithoutScheme(): string {
		return $this->removeSchemeFromUrl($this->getRemoteBaseUrl());
	}

	/**
	 * returns the reference to the current line being executed.
	 *
	 * @return string
	 */
	public function getStepLineRef(): string {
		if (!$this->sendStepLineRef) {
			return '';
		}

		// If we are in BeforeScenario and possibly before any particular step
		// is being executed, then stepLineRef might be empty. In that case
		// return just the string for the scenario.
		if ($this->stepLineRef === '') {
			return $this->scenarioString;
		}
		return $this->stepLineRef;
	}

	/**
	 * returns the base URL without any sub-path e.g. http://localhost:8080
	 * of the base URL http://localhost:8080/owncloud
	 *
	 * @return string
	 */
	public function getBaseUrlWithoutPath(): string {
		$parts = \parse_url($this->getBaseUrl());
		$url = $parts ["scheme"] . "://" . $parts["host"];
		if (isset($parts["port"])) {
			$url = "$url:" . $parts["port"];
		}
		return $url;
	}

	/**
	 * @return int
	 */
	public function getOcsApiVersion(): int {
		return $this->ocsApiVersion;
	}

	/**
	 * @return string|null
	 */
	public function getSourceIpAddress(): ?string {
		return $this->sourceIpAddress;
	}

	/**
	 * @return array|null
	 */
	public function getStorageIds(): ?array {
		return $this->storageIds;
	}

	/**
	 * @param string $storageName
	 *
	 * @return integer
	 * @throws Exception
	 */
	public function getStorageId(string $storageName): int {
		$storageIds = $this->getStorageIds();
		$storageId = \array_search($storageName, $storageIds);
		Assert::assertNotFalse(
			$storageId,
			"Could not find storageId with storage name $storageName"
		);
		return $storageId;
	}

	/**
	 * @param integer $storageId
	 *
	 * @return void
	 */
	public function popStorageId(int $storageId): void {
		unset($this->storageIds[$storageId]);
	}

	/**
	 * @param string $sourceIpAddress
	 *
	 * @return void
	 */
	public function setSourceIpAddress(string $sourceIpAddress): void {
		$this->sourceIpAddress = $sourceIpAddress;
	}

	/**
	 * @return array
	 */
	public function getGuzzleClientHeaders(): array {
		return $this->guzzleClientHeaders;
	}

	/**
	 * @param array $guzzleClientHeaders ['X-Foo' => 'Bar']
	 *
	 * @return void
	 */
	public function setGuzzleClientHeaders(array $guzzleClientHeaders): void {
		$this->guzzleClientHeaders = $guzzleClientHeaders;
	}

	/**
	 * @param array $guzzleClientHeaders ['X-Foo' => 'Bar']
	 *
	 * @return void
	 */
	public function addGuzzleClientHeaders(array $guzzleClientHeaders): void {
		$this->guzzleClientHeaders = \array_merge(
			$this->guzzleClientHeaders,
			$guzzleClientHeaders
		);
	}

	/**
	 * @Given /^using OCS API version "([^"]*)"$/
	 *
	 * @param string $version
	 *
	 * @return void
	 */
	public function usingOcsApiVersion(string $version): void {
		$this->ocsApiVersion = (int)$version;
	}

	/**
	 * @Given /^as user "([^"]*)"$/
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function asUser(string $user): void {
		$this->currentUser = $this->getActualUsername($user);
	}

	/**
	 * @Given as the administrator
	 *
	 * @return void
	 */
	public function asTheAdministrator(): void {
		$this->currentUser = $this->getAdminUsername();
	}

	/**
	 * @return string
	 */
	public function getCurrentUser(): string {
		return $this->currentUser;
	}

	/**
	 * @param string $user
	 *
	 * @return void
	 */
	public function setCurrentUser(string $user): void {
		$this->currentUser = $user;
	}

	/**
	 * returns $this->response
	 * some steps use that private var to store the response for other steps
	 *
	 * @return ResponseInterface
	 */
	public function getResponse(): ?ResponseInterface {
		return $this->response;
	}

	/**
	 * let this class remember a response that was received elsewhere
	 * so that steps in this class can be used to examine the response
	 *
	 * @param ResponseInterface|null $response
	 * @param string $username of the user that received the response
	 *
	 * @return void
	 */
	public function setResponse(
		?ResponseInterface $response,
		string             $username = ""
	): void {
		$this->response = $response;
		//after a new response reset the response xml
		$this->responseXml = [];
		//after a new response reset the response xml object
		$this->responseXmlObject = null;
		// remember the user that received the response
		$this->responseUser = $username;
	}

	/**
	 * @return string
	 */
	public function getCurrentServer(): string {
		return $this->currentServer;
	}

	/**
	 * @Given /^using server "(LOCAL|REMOTE)"$/
	 *
	 * @param string|null $server
	 *
	 * @return string Previous used server
	 */
	public function usingServer(?string $server): string {
		$previousServer = $this->currentServer;
		if ($server === 'LOCAL') {
			$this->baseUrl = $this->localBaseUrl;
			$this->currentServer = 'LOCAL';
		} else {
			$this->baseUrl = $this->remoteBaseUrl;
			$this->currentServer = 'REMOTE';
		}
		return $previousServer;
	}

	/**
	 *
	 * @return boolean
	 */
	public function federatedServerExists(): bool {
		return $this->federatedServerExists;
	}

	/**
	 * Parses the response as XML
	 *
	 * @param ResponseInterface|null $response
	 * @param string|null $exceptionText text to put at the front of exception messages
	 *
	 * @return SimpleXMLElement
	 * @throws Exception
	 */
	public function getResponseXml(?ResponseInterface $response = null, ?string $exceptionText = ''): SimpleXMLElement {
		if ($response === null) {
			$response = $this->response;
		}

		if ($exceptionText === '') {
			$exceptionText = __METHOD__;
		}
		return HttpRequestHelper::getResponseXml($response, $exceptionText);
	}

	/**
	 * Parses the xml answer to get the requested key and sub-key
	 *
	 * @param ResponseInterface $response
	 * @param string $key1
	 * @param string $key2
	 *
	 * @return string
	 * @throws Exception
	 */
	public function getXMLKey1Key2Value(ResponseInterface $response, string $key1, string $key2): string {
		return (string)$this->getResponseXml($response, __METHOD__)->$key1->$key2;
	}

	/**
	 * Parses the xml answer to get the requested key sequence
	 *
	 * @param ResponseInterface $response
	 * @param string $key1
	 * @param string $key2
	 * @param string $key3
	 *
	 * @return string
	 * @throws Exception
	 */
	public function getXMLKey1Key2Key3Value(
		ResponseInterface $response,
		string            $key1,
		string            $key2,
		string            $key3
	): string {
		return (string)$this->getResponseXml($response, __METHOD__)->$key1->$key2->$key3;
	}

	/**
	 * Parses the xml answer to get the requested attribute value
	 *
	 * @param ResponseInterface $response
	 * @param string $key1
	 * @param string $key2
	 * @param string $key3
	 * @param string $attribute
	 *
	 * @return string
	 * @throws Exception
	 */
	public function getXMLKey1Key2Key3AttributeValue(
		ResponseInterface $response,
		string            $key1,
		string            $key2,
		string            $key3,
		string            $attribute
	): string {
		return (string)$this->getResponseXml($response, __METHOD__)->$key1->$key2->$key3->attributes()->$attribute;
	}

	/**
	 * This function is needed to use a vertical fashion in the gherkin tables.
	 *
	 * @param array $arrayOfArrays
	 *
	 * @return array
	 */
	public function simplifyArray(array $arrayOfArrays): array {
		$a = \array_map(
			function ($subArray) {
				return $subArray[0];
			},
			$arrayOfArrays
		);
		return $a;
	}

	/**
	 * @When /^user "([^"]*)" sends HTTP method "([^"]*)" to URL "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $verb
	 * @param string $url
	 *
	 * @return void
	 */
	public function userSendsHTTPMethodToUrl(string $user, string $verb, string $url): void {
		$user = $this->getActualUsername($user);
		$this->setResponse($this->sendingToWithDirectUrl($user, $verb, $url, null));
	}

	/**
	 * @Given /^user "([^"]*)" has sent HTTP method "([^"]*)" to URL "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $verb
	 * @param string $url
	 *
	 * @return void
	 */
	public function userHasSentHTTPMethodToUrl(string $user, string $verb, string $url): void {
		$this->userSendsHTTPMethodToUrl($user, $verb, $url);
		$this->theHTTPStatusCodeShouldBeSuccess();
	}

	/**
	 * @When user :user sends HTTP method :method to URL :davPath with content :content
	 *
	 * @param string $user
	 * @param string $method
	 * @param string $davPath
	 * @param string $content
	 *
	 * @return void
	 */
	public function userSendsHttpMethodToUrlWithContent(string $user, string $method, string $davPath, string $content): void {
		$this->setResponse($this->sendingToWithDirectUrl($user, $method, $davPath, $content));
	}

	/**
	 * @When /^user "([^"]*)" sends HTTP method "([^"]*)" to URL "([^"]*)" with password "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $verb
	 * @param string $url
	 * @param string $password
	 *
	 * @return void
	 */
	public function userSendsHTTPMethodToUrlWithPassword(string $user, string $verb, string $url, string $password): void {
		$this->setResponse($this->sendingToWithDirectUrl($user, $verb, $url, null, $password));
	}

	/**
	 * @Given /^user "([^"]*)" has sent HTTP method "([^"]*)" to URL "([^"]*)" with password "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $verb
	 * @param string $url
	 * @param string $password
	 *
	 * @return void
	 */
	public function userHasSentHTTPMethodToUrlWithPassword(string $user, string $verb, string $url, string $password): void {
		$this->userSendsHTTPMethodToUrlWithPassword($user, $verb, $url, $password);
		$this->theHTTPStatusCodeShouldBeSuccess();
	}

	/**
	 * @param string|null $user
	 * @param string|null $verb
	 * @param string|null $url
	 * @param string|null $body
	 * @param string|null $password
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public function sendingToWithDirectUrl(?string $user, ?string $verb, ?string $url, string $body = null, ?string $password = null): ResponseInterface {
		$fullUrl = $this->getBaseUrl() . $url;

		if ($password === null) {
			$password = $this->getPasswordForUser($user);
		}

		$headers = $this->guzzleClientHeaders;

		$config = null;
		if ($this->sourceIpAddress !== null) {
			$config = [
				'curl' => [
					CURLOPT_INTERFACE => $this->sourceIpAddress
				]
			];
		}

		$cookies = null;
		if (!empty($this->cookieJar->toArray())) {
			$cookies = $this->cookieJar;
		}

		if (isset($this->requestToken)) {
			$headers['requesttoken'] = $this->requestToken;
		}

		return HttpRequestHelper::sendRequest(
			$fullUrl,
			$this->getStepLineRef(),
			$verb,
			$user,
			$password,
			$headers,
			$body,
			$config,
			$cookies
		);
	}

	/**
	 * @param string $url
	 *
	 * @return bool
	 */
	public function isAPublicLinkUrl(string $url): bool {
		if (OcisHelper::isTestingOnReva()) {
			$urlEnding = \ltrim($url, '/');
		} else {
			if (\substr($url, 0, 4) !== "http") {
				return false;
			}
			$urlEnding = \substr($url, \strlen($this->getBaseUrl() . '/'));
		}

		$matchResult = \preg_match("%^(#/)?s/([a-zA-Z0-9]{15})$%", $urlEnding);

		// preg_match returns (int) 1 for a match, we want to return a boolean.
		if ($matchResult === 1) {
			$isPublicLinkUrl = true;
		} else {
			$isPublicLinkUrl = false;
		}
		return $isPublicLinkUrl;
	}

	/**
	 * Check that the status code in the saved response is the expected status
	 * code, or one of the expected status codes.
	 *
	 * @param int|int[]|string|string[] $expectedStatusCode
	 * @param string|null $message
	 * @param ResponseInterface|null $response
	 *
	 * @return void
	 */
	public function theHTTPStatusCodeShouldBe($expectedStatusCode, ?string $message = "", ?ResponseInterface $response = null): void {
		$response = $response ?? $this->response;
		$actualStatusCode = $response->getStatusCode();
		if (\is_array($expectedStatusCode)) {
			if ($message === "") {
				$message = "HTTP status code $actualStatusCode is not one of the expected values " . \implode(" or ", $expectedStatusCode);
			}

			Assert::assertContainsEquals(
				$actualStatusCode,
				$expectedStatusCode,
				$message
			);
		} else {
			if ($message === "") {
				$message = "HTTP status code $actualStatusCode is not the expected value $expectedStatusCode";
			}

			Assert::assertEquals(
				$expectedStatusCode,
				$actualStatusCode,
				$message
			);
		}
	}

	/**
	 * @param PyStringNode|string $schemaString
	 *
	 * @return mixed
	 */
	public function getJSONSchema($schemaString) {
		if (\gettype($schemaString) !== 'string') {
			$schemaString = $schemaString->getRaw();
		}
		$schemaString = $this->substituteInLineCodes($schemaString);
		$schema = \json_decode($schemaString);
		Assert::assertNotNull($schema, 'schema is not valid JSON');
		return $schema;
	}

	/**
	 * returns json decoded body content of a json response as an object
	 *
	 * @param ResponseInterface|null $response
	 *
	 * @return object
	 */
	public function getJsonDecodedResponseBodyContent(ResponseInterface $response = null):?object {
		$response = $response ?? $this->response;
		if ($response !== null) {
			$response->getBody()->rewind();
			return json_decode($response->getBody()->getContents());
		}
		return null;
	}

	/**
	 * @Then the ocs JSON data of the response should match
	 *
	 * @param PyStringNode $schemaString
	 *
	 * @return void
	 *
	 * @throws Exception
	 */
	public function theOcsDataOfTheResponseShouldMatch(
		PyStringNode $schemaString
	): void {
		$jsonResponse = $this->getJsonDecodedResponseBodyContent();
		$this->assertJsonDocumentMatchesSchema(
			$jsonResponse->ocs->data,
			$this->getJSONSchema($schemaString)
		);
	}

	/**
	 * @Then the JSON data of the response should match
	 *
	 * @param PyStringNode $schemaString
	 * @param ResponseInterface|null $response
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theDataOfTheResponseShouldMatch(PyStringNode $schemaString, ResponseInterface $response=null): void {
		$responseBody = $this->getJsonDecodedResponseBodyContent($response);
		$this->assertJsonDocumentMatchesSchema(
			$responseBody,
			$this->getJSONSchema($schemaString)
		);
	}

	/**
	 * @Then /^the HTTP status code should be "([^"]*)"$/
	 *
	 * @param int|string $statusCode
	 *
	 * @return void
	 */
	public function thenTheHTTPStatusCodeShouldBe($statusCode): void {
		$this->theHTTPStatusCodeShouldBe($statusCode);
	}

	/**
	 * @Then /^the HTTP status code should be "([^"]*)" or "([^"]*)"$/
	 *
	 * @param int|string $statusCode1
	 * @param int|string $statusCode2
	 *
	 * @return void
	 */
	public function theHTTPStatusCodeShouldBeOr($statusCode1, $statusCode2): void {
		$this->theHTTPStatusCodeShouldBe(
			[$statusCode1, $statusCode2]
		);
	}

	/**
	 * @Then /^the HTTP status code should be between "(\d+)" and "(\d+)"$/
	 *
	 * @param int|string $minStatusCode
	 * @param int|string $maxStatusCode
	 * @param ResponseInterface|null $response
	 *
	 * @return void
	 */
	public function theHTTPStatusCodeShouldBeBetween(
		$minStatusCode,
		$maxStatusCode,
		?ResponseInterface $response= null
	): void {
		$response = $response ?? $this->response;
		$statusCode = $response->getStatusCode();
		$message = "The HTTP status code $statusCode is not between $minStatusCode and $maxStatusCode";
		Assert::assertGreaterThanOrEqual(
			$minStatusCode,
			$statusCode,
			$message
		);
		Assert::assertLessThanOrEqual(
			$maxStatusCode,
			$statusCode,
			$message
		);
	}

	/**
	 * @Then the HTTP status code should be success
	 *
	 * @return void
	 */
	public function theHTTPStatusCodeShouldBeSuccess(): void {
		$this->theHTTPStatusCodeShouldBeBetween(200, 299);
	}

	/**
	 * @Then the HTTP status code should be failure
	 *
	 * @return void
	 */
	public function theHTTPStatusCodeShouldBeFailure(): void {
		$statusCode = $this->response->getStatusCode();
		$message = "The HTTP status code $statusCode is not greater than or equals to 400";
		Assert::assertGreaterThanOrEqual(
			400,
			$statusCode,
			$message
		);
	}

	/**
	 *
	 * @return bool
	 */
	public function theHTTPStatusCodeWasSuccess(): bool {
		$statusCode = $this->response->getStatusCode();
		return (($statusCode >= 200) && ($statusCode <= 299));
	}

	/**
	 * Check the text in an HTTP responseXml message
	 *
	 * @Then /^the HTTP response message should be "([^"]*)"$/
	 *
	 * @param string $expectedMessage
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theHttpResponseMessageShouldBe(string $expectedMessage): void {
		$actualMessage = $this->responseXml['value'][1]['value'];
		Assert::assertEquals(
			$expectedMessage,
			$actualMessage,
			"Expected $expectedMessage HTTP response message but got $actualMessage"
		);
	}

	/**
	 * Check the text in an HTTP reason phrase
	 *
	 * @Then /^the HTTP reason phrase should be "([^"]*)"$/
	 *
	 * @param string $reasonPhrase
	 *
	 * @return void
	 */
	public function theHTTPReasonPhraseShouldBe(string $reasonPhrase): void {
		Assert::assertEquals(
			$reasonPhrase,
			$this->getResponse()->getReasonPhrase(),
			'Unexpected HTTP reason phrase in response'
		);
	}

	/**
	 * Check the text in an HTTP reason phrase
	 * Use this step form if the expected text contains double quotes,
	 * single quotes and other content that theHTTPReasonPhraseShouldBe()
	 * cannot handle.
	 *
	 * After the step, write the expected text in PyString form like:
	 *
	 * """
	 * File "abc.txt" can't be shared due to reason "xyz"
	 * """
	 *
	 * @Then /^the HTTP reason phrase should be:$/
	 *
	 * @param PyStringNode $reasonPhrase
	 *
	 * @return void
	 */
	public function theHTTPReasonPhraseShouldBePyString(
		PyStringNode $reasonPhrase
	): void {
		Assert::assertEquals(
			$reasonPhrase->getRaw(),
			$this->getResponse()->getReasonPhrase(),
			'Unexpected HTTP reason phrase in response'
		);
	}

	/**
	 * @Then /^the XML "([^"]*)" "([^"]*)" value should be "([^"]*)"$/
	 *
	 * @param string $key1
	 * @param string $key2
	 * @param string $idText
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theXMLKey1Key2ValueShouldBe(string $key1, string $key2, string $idText): void {
		$actualValue = $this->getXMLKey1Key2Value($this->response, $key1, $key2);
		Assert::assertEquals(
			$idText,
			$actualValue,
			"Expected $idText but got "
			. $actualValue
		);
	}

	/**
	 * @Then /^the XML "([^"]*)" "([^"]*)" "([^"]*)" value should be "([^"]*)"$/
	 *
	 * @param string $key1
	 * @param string $key2
	 * @param string $key3
	 * @param string $idText
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theXMLKey1Key2Key3ValueShouldBe(
		string $key1,
		string $key2,
		string $key3,
		string $idText
	) {
		$actualValue = $this->getXMLKey1Key2Key3Value($this->response, $key1, $key2, $key3);
		Assert::assertEquals(
			$idText,
			$actualValue,
			"Expected $idText but got "
			. $actualValue
		);
	}

	/**
	 * @Then /^the XML "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)" attribute value should be a valid version string$/
	 *
	 * @param string $key1
	 * @param string $key2
	 * @param string $key3
	 * @param string $attribute
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theXMLKey1Key2AttributeValueShouldBe(
		string $key1,
		string $key2,
		string $key3,
		string $attribute
	): void {
		$value = $this->getXMLKey1Key2Key3AttributeValue(
			$this->response,
			$key1,
			$key2,
			$key3,
			$attribute
		);
		Assert::assertTrue(
			\version_compare($value, '0.0.1') >= 0,
			"attribute $attribute value $value is not a valid version string"
		);
	}

	/**
	 * @param ResponseInterface $response
	 *
	 * @return void
	 */
	public function extractRequestTokenFromResponse(ResponseInterface $response): void {
		$this->requestToken = \substr(
			\preg_replace(
				'/(.*)data-requesttoken="(.*)">(.*)/sm',
				'\2',
				$response->getBody()->getContents()
			),
			0,
			89
		);
	}

	/**
	 * @Given /^user "([^"]*)" has logged in to a web-style session$/
	 *
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function userHasLoggedInToAWebStyleSessionUsingTheAPI(string $user): void {
		$user = $this->getActualUsername($user);
		$loginUrl = $this->getBaseUrl() . '/login';
		// Request a new session and extract CSRF token

		$config = null;
		if ($this->sourceIpAddress !== null) {
			$config = [
				'curl' => [
					CURLOPT_INTERFACE => $this->sourceIpAddress
				]
			];
		}

		$response = HttpRequestHelper::get(
			$loginUrl,
			$this->getStepLineRef(),
			null,
			null,
			$this->guzzleClientHeaders,
			null,
			$config,
			$this->cookieJar
		);
		$this->theHTTPStatusCodeShouldBeBetween(200, 299, $response);
		$this->extractRequestTokenFromResponse($response);

		// Login and extract new token
		$password = $this->getPasswordForUser($user);
		$body = [
			'user' => $user,
			'password' => $password,
			'requesttoken' => $this->requestToken
		];
		$response = HttpRequestHelper::post(
			$loginUrl,
			$this->getStepLineRef(),
			null,
			null,
			$this->guzzleClientHeaders,
			$body,
			$config,
			$this->cookieJar
		);
		$this->theHTTPStatusCodeShouldBeBetween(200, 299, $response);
		$this->extractRequestTokenFromResponse($response);
	}

	/**
	 * @When the client sends a :method to :url of user :user with requesttoken
	 *
	 * @param string $method
	 * @param string $url
	 * @param string $user
	 *
	 * @return void
	 * @throws GuzzleException
	 * @throws JsonException
	 */
	public function sendingAToWithRequesttoken(
		string $method,
		string $url,
		string $user
	): void {
		$headers = $this->guzzleClientHeaders;

		$config = null;
		if ($this->sourceIpAddress !== null) {
			$config = [
				'curl' => [
					CURLOPT_INTERFACE => $this->sourceIpAddress
				]
			];
		}

		$headers['requesttoken'] = $this->requestToken;

		$user = \strtolower($this->getActualUsername($user));
		$url = $this->getBaseUrl() . $url;
		$url = $this->substituteInLineCodes($url, $user);
		$this->response = HttpRequestHelper::sendRequest(
			$url,
			$this->getStepLineRef(),
			$method,
			null,
			null,
			$headers,
			null,
			$config,
			$this->cookieJar
		);
	}

	/**
	 * @Given the client has sent a :method to :url of user :user with requesttoken
	 *
	 * @param string $method
	 * @param string $url
	 * @param string $user
	 *
	 * @return void
	 */
	public function theClientHasSentAToWithRequesttoken(
		string $method,
		string $url,
		string $user
	): void {
		$this->sendingAToWithRequesttoken($method, $url, $user);
		$this->theHTTPStatusCodeShouldBeSuccess();
	}

	/**
	 * @When the client sends a :method to :url of user :user without requesttoken
	 *
	 * @param string $method
	 * @param string $url
	 * @param string|null $user
	 *
	 * @return void
	 * @throws JsonException
	 */
	public function sendingAToWithoutRequesttoken(string $method, string $url, ?string $user = null): void {
		$config = null;
		if ($this->sourceIpAddress !== null) {
			$config = [
				'curl' => [
					CURLOPT_INTERFACE => $this->sourceIpAddress
				]
			];
		}

		$user = \strtolower($this->getActualUsername($user));
		$url = $this->getBaseUrl() . $url;
		$url = $this->substituteInLineCodes($url, $user);
		$this->response = HttpRequestHelper::sendRequest(
			$url,
			$this->getStepLineRef(),
			$method,
			null,
			null,
			$this->guzzleClientHeaders,
			null,
			$config,
			$this->cookieJar
		);
	}

	/**
	 * @Given the client has sent a :method to :url without requesttoken
	 *
	 * @param string $method
	 * @param string $url
	 *
	 * @return void
	 */
	public function theClientHasSentAToWithoutRequesttoken(string $method, string $url): void {
		$this->sendingAToWithoutRequesttoken($method, $url);
		$this->theHTTPStatusCodeShouldBeSuccess();
	}

	/**
	 * @param string $path
	 * @param string $filename
	 *
	 * @return void
	 */
	public static function removeFile(string $path, string $filename): void {
		if (\file_exists("$path$filename")) {
			\unlink("$path$filename");
		}
	}

	/**
	 * Creates a file locally in the file system of the test runner
	 * The file will be available to upload to the server
	 *
	 * @param string $name
	 * @param string $size
	 * @param string $endData
	 *
	 * @return void
	 */
	public function createLocalFileOfSpecificSize(string $name, string $size, string $endData = 'a'): void {
		$folder = $this->workStorageDirLocation();
		if (!\is_dir($folder)) {
			\mkDir($folder);
		}
		$file = \fopen($folder . $name, 'w');
		\fseek($file, $size - \strlen($endData), SEEK_CUR);
		\fwrite($file, $endData); // write the end data to force the file size
		\fclose($file);
	}

	/**
	 * Make a directory under the server root on the ownCloud server
	 *
	 * @param string $dirPathFromServerRoot e.g. 'apps2/myapp/appinfo'
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function mkDirOnServer(string $dirPathFromServerRoot): void {
		SetupHelper::mkDirOnServer(
			$dirPathFromServerRoot,
			$this->getStepLineRef(),
			$this->getBaseUrl(),
			$this->getAdminUsername(),
			$this->getAdminPassword()
		);
	}

	/**
	 * @param string $filePathFromServerRoot
	 * @param string $content
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function createFileOnServerWithContent(
		string $filePathFromServerRoot,
		string $content
	): void {
		SetupHelper::createFileOnServer(
			$filePathFromServerRoot,
			$content,
			$this->getStepLineRef(),
			$this->getBaseUrl(),
			$this->getAdminUsername(),
			$this->getAdminPassword()
		);
	}

	/**
	 * @param string $user
	 *
	 * @return boolean
	 */
	public function isAdminUsername(string $user): bool {
		return ($user === $this->getAdminUsername());
	}

	/**
	 * @return string
	 */
	public function getAdminUsername(): string {
		return $this->adminUsername;
	}

	/**
	 * @return string
	 */
	public function getAdminPassword(): string {
		return $this->adminPassword;
	}

	/**
	 * @param string $password
	 *
	 * @return void
	 */
	public function rememberNewAdminPassword(string $password): void {
		$this->adminPassword = $password;
	}

	/**
	 * @param string|null $userName
	 *
	 * @return string
	 */
	public function getPasswordForUser(?string $userName): string {
		$userNameNormalized = $this->normalizeUsername($userName);
		$username = $this->getActualUsername($userNameNormalized);
		if ($username === $this->getAdminUsername()) {
			return $this->getAdminPassword();
		} elseif (\array_key_exists($username, $this->createdUsers)) {
			return (string)$this->createdUsers[$username]['password'];
		} elseif (\array_key_exists($username, $this->createdRemoteUsers)) {
			return (string)$this->createdRemoteUsers[$username]['password'];
		}

		// The user has not been created yet, see if there is a replacement
		// defined for the user.
		$usernameReplacements = $this->usersToBeReplaced();
		if (isset($usernameReplacements)) {
			if (isset($usernameReplacements[$userNameNormalized])) {
				return $usernameReplacements[$userNameNormalized]['password'];
			}
		}

		// Fall back to the default password used for the well-known users.
		if ($username === 'regularuser') {
			return $this->regularUserPassword;
		} elseif ($username === 'alice') {
			return $this->regularUserPassword;
		} elseif ($username === 'brian') {
			return $this->alt1UserPassword;
		} elseif ($username === 'carol') {
			return $this->alt2UserPassword;
		} elseif ($username === 'david') {
			return $this->alt3UserPassword;
		} elseif ($username === 'emily') {
			return $this->alt4UserPassword;
		} elseif ($username === 'usergrp') {
			return $this->regularUserPassword;
		} elseif ($username === 'sharee1') {
			return $this->regularUserPassword;
		}

		// The user has not been created yet and is not one of the pre-known
		// users. So let the caller have the default password.
		return (string)$this->getActualPassword($this->regularUserPassword);
	}

	/**
	 * Get the display name of the user.
	 *
	 * For users that have already been created, return their display name.
	 * For special known usernames, return the display name that is also used by LDAP tests.
	 * For other users, return null. They will not be assigned any particular
	 * display name by this function.
	 *
	 * @param string $userName
	 *
	 * @return string|null
	 */
	public function getDisplayNameForUser(string $userName): ?string {
		$userNameNormalized = $this->normalizeUsername($userName);
		$username = $this->getActualUsername($userNameNormalized);
		if (\array_key_exists($username, $this->createdUsers)) {
			if (isset($this->createdUsers[$username]['displayname'])) {
				return (string)$this->createdUsers[$username]['displayname'];
			}
			return $userName;
		}
		if (\array_key_exists($username, $this->createdRemoteUsers)) {
			if (isset($this->createdRemoteUsers[$username]['displayname'])) {
				return (string)$this->createdRemoteUsers[$username]['displayname'];
			}
			return $userName;
		}

		// The user has not been created yet, see if there is a replacement
		// defined for the user.
		$usernameReplacements = $this->usersToBeReplaced();
		if (isset($usernameReplacements)) {
			if (isset($usernameReplacements[$userNameNormalized])) {
				return $usernameReplacements[$userNameNormalized]['displayname'];
			} elseif (isset($usernameReplacements[$userName])) {
				return $usernameReplacements[$userName]['displayname'];
			}
		}

		// Fall back to the default display name used for the well-known users.
		if ($username === 'regularuser') {
			return 'Regular User';
		} elseif ($username === 'alice') {
			return 'Alice Hansen';
		} elseif ($username === 'brian') {
			return 'Brian Murphy';
		} elseif ($username === 'carol') {
			return 'Carol King';
		} elseif ($username === 'david') {
			return 'David Lopez';
		} elseif ($username === 'emily') {
			return 'Emily Wagner';
		} elseif ($username === 'usergrp') {
			return 'User Grp';
		} elseif ($username === 'sharee1') {
			return 'Sharee One';
		} elseif ($username === 'sharee2') {
			return 'Sharee Two';
		} elseif (\in_array($username, ["grp1", "***redacted***"])) {
			return $username;
		}
		return null;
	}

	/**
	 * Get the email address of the user.
	 *
	 * For users that have already been created, return their email address.
	 * For special known usernames, return the email address that is also used by LDAP tests.
	 * For other users, return null. They will not be assigned any particular
	 * email address by this function.
	 *
	 * @param string $userName
	 *
	 * @return string|null
	 */
	public function getEmailAddressForUser(string $userName): ?string {
		$userNameNormalized = $this->normalizeUsername($userName);
		$username = $this->getActualUsername($userNameNormalized);
		if (\array_key_exists($username, $this->createdUsers)) {
			return (string)$this->createdUsers[$username]['email'];
		}
		if (\array_key_exists($username, $this->createdRemoteUsers)) {
			return (string)$this->createdRemoteUsers[$username]['email'];
		}

		// The user has not been created yet, see if there is a replacement
		// defined for the user.
		$usernameReplacements = $this->usersToBeReplaced();
		if (isset($usernameReplacements)) {
			if (isset($usernameReplacements[$userNameNormalized])) {
				return $usernameReplacements[$userNameNormalized]['email'];
			} elseif (isset($usernameReplacements[$userName])) {
				return $usernameReplacements[$userName]['email'];
			}
		}

		// Fall back to the default display name used for the well-known users.
		if ($username === 'regularuser') {
			return 'regularuser@example.org';
		} elseif ($username === 'alice') {
			return 'alice@example.org';
		} elseif ($username === 'brian') {
			return 'brian@example.org';
		} elseif ($username === 'carol') {
			return 'carol@example.org';
		} elseif ($username === 'david') {
			return 'david@example.org';
		} elseif ($username === 'emily') {
			return 'emily@example.org';
		} elseif ($username === 'usergrp') {
			return 'usergrp@example.org';
		} elseif ($username === 'sharee1') {
			return 'sharee1@example.org';
		} else {
			return null;
		}
	}

	// TODO do similar for other usernames for e.g. %regularuser% or %test-user-1%

	/**
	 * @param string|null $functionalUsername
	 *
	 * @return string|null
	 * @throws JsonException
	 */
	public function getActualUsername(?string $functionalUsername): ?string {
		if ($functionalUsername === null) {
			return null;
		}
		$usernames = $this->usersToBeReplaced();
		if (isset($usernames)) {
			if (isset($usernames[$functionalUsername])) {
				return $usernames[$functionalUsername]['username'];
			}
			$normalizedUsername = $this->normalizeUsername($functionalUsername);
			if (isset($usernames[$normalizedUsername])) {
				return $usernames[$normalizedUsername]['username'];
			}
		}
		if ($functionalUsername === "%admin%") {
			return $this->getAdminUsername();
		}
		return $functionalUsername;
	}

	/**
	 * @param string|null $functionalPassword
	 *
	 * @return string|null
	 */
	public function getActualPassword(?string $functionalPassword): ?string {
		if ($functionalPassword === "%regular%") {
			return $this->regularUserPassword;
		} elseif ($functionalPassword === "%alt1%") {
			return $this->alt1UserPassword;
		} elseif ($functionalPassword === "%alt2%") {
			return $this->alt2UserPassword;
		} elseif ($functionalPassword === "%alt3%") {
			return $this->alt3UserPassword;
		} elseif ($functionalPassword === "%alt4%") {
			return $this->alt4UserPassword;
		} elseif ($functionalPassword === "%subadmin%") {
			return $this->subAdminPassword;
		} elseif ($functionalPassword === "%admin%") {
			return $this->getAdminPassword();
		} elseif ($functionalPassword === "%altadmin%") {
			return $this->alternateAdminPassword;
		} elseif ($functionalPassword === "%public%") {
			return $this->publicLinkSharePassword;
		} elseif ($functionalPassword === "%remove%") {
			return "";
		} else {
			return $functionalPassword;
		}
	}

	/**
	 * @param string $userName
	 *
	 * @return array
	 */
	public function getAuthOptionForUser(string $userName): array {
		return [$userName, $this->getPasswordForUser($userName)];
	}

	/**
	 * @return array
	 */
	public function getAuthOptionForAdmin(): array {
		return $this->getAuthOptionForUser($this->getAdminUsername());
	}

	/**
	 * @When the administrator requests status.php
	 *
	 * @return void
	 */
	public function theAdministratorRequestsStatusPhp(): void {
		$this->response = $this->getStatusPhp();
	}

	/**
	 * Copy a file from the test-runner to the temporary storage directory on
	 * the system-under-test. This uses the testing app to push the file into
	 * the backend of the server, where it can be seen by occ commands done in
	 * the server-under-test.
	 *
	 * @Given the administrator has copied file :localPath to :destination in temporary storage on the system under test
	 *
	 * @param string $localPath relative to the core "root" folder
	 * @param string $destination
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theAdministratorHasCopiedFileToTemporaryStorageOnTheSystemUnderTest(
		string $localPath,
		string $destination
	): void {
		// FeatureContext is in tests/acceptance/features/bootstrap so go up 4
		// levels to the test-runner root
		$testRunnerRoot = \dirname(__DIR__, 4);
		// The local path is specified down from the root - e.g. tests/data/file.txt
		$content = \file_get_contents("$testRunnerRoot/$localPath");
		Assert::assertNotFalse(
			$content,
			"Local file $localPath cannot be read"
		);
		$this->copyContentToFileInTemporaryStorageOnSystemUnderTest($destination, $content);
		$this->theFileWithContentShouldExistInTheServerRoot(TEMPORARY_STORAGE_DIR_ON_REMOTE_SERVER . "/$destination", $content);
	}

	/**
	 * @param string $destination
	 * @param string $content
	 *
	 * @return ResponseInterface
	 * @throws Exception
	 */
	public function copyContentToFileInTemporaryStorageOnSystemUnderTest(
		string $destination,
		string $content
	): ResponseInterface {
		$this->mkDirOnServer(TEMPORARY_STORAGE_DIR_ON_REMOTE_SERVER);

		return OcsApiHelper::sendRequest(
			$this->getBaseUrl(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			'POST',
			"/apps/testing/api/v1/file",
			$this->getStepLineRef(),
			[
				'file' => TEMPORARY_STORAGE_DIR_ON_REMOTE_SERVER . "/$destination",
				'content' => $content
			],
			$this->getOcsApiVersion()
		);
	}

	/**
	 * @Given a file with the size of :size bytes and the name :name has been created locally
	 *
	 * @param int $size if not int given it will be cast to int
	 * @param string $name
	 *
	 * @return void
	 * @throws InvalidArgumentException
	 */
	public function aFileWithSizeAndNameHasBeenCreatedLocally(int $size, string $name): void {
		$fullPath = UploadHelper::getUploadFilesDir($name);
		if (\file_exists($fullPath)) {
			throw new InvalidArgumentException(
				__METHOD__ . " could not create '$fullPath' file exists"
			);
		}
		UploadHelper::createFileSpecificSize($fullPath, $size);
		$this->createdFiles[] = $fullPath;
	}

	/**
	 *
	 * @return ResponseInterface
	 */
	public function getStatusPhp(): ResponseInterface {
		$fullUrl = $this->getBaseUrl() . "/status.php";

		$config = null;
		if ($this->sourceIpAddress !== null) {
			$config = [
				'curl' => [
					CURLOPT_INTERFACE => $this->sourceIpAddress
				]
			];
		}

		return HttpRequestHelper::get(
			$fullUrl,
			$this->getStepLineRef(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			$this->guzzleClientHeaders,
			null,
			$config
		);
	}

	/**
	 * @Then the json responded should match with
	 *
	 * @param PyStringNode $jsonExpected
	 *
	 * @return void
	 */
	public function jsonRespondedShouldMatch(PyStringNode $jsonExpected): void {
		$jsonExpectedEncoded = \json_encode($jsonExpected->getRaw());
		$jsonRespondedEncoded = \json_encode((string)$this->response->getBody());
		Assert::assertEquals(
			$jsonExpectedEncoded,
			$jsonRespondedEncoded,
			"The json responded: $jsonRespondedEncoded does not match with json expected: $jsonExpectedEncoded"
		);
	}

	/**
	 * send request to read a server file for core
	 *
	 * @param string $path
	 *
	 * @return void
	 */
	public function readFileInServerRootForCore(string $path): ResponseInterface {
		return OcsApiHelper::sendRequest(
			$this->getBaseUrl(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			'GET',
			"/apps/testing/api/v1/file?file=$path",
			$this->getStepLineRef()
		);
	}

	/**
	 * read a server file for ocis
	 *
	 * @param string $path
	 *
	 * @return string
	 * @throws Exception
	 */
	public function readFileInServerRootForOCIS(string $path): string {
		$pathToOcis = \getenv("PATH_TO_OCIS");
		$targetFile = \rtrim($pathToOcis, "/") . "/" . "services/web/assets" . "/" . ltrim($path, '/');
		if (!\file_exists($targetFile)) {
			throw new Exception('Target File ' . $targetFile . ' could not be found');
		}
		return \file_get_contents($targetFile);
	}

	/**
	 * send request to list a server file
	 *
	 * @param string $path
	 *
	 * @return void
	 */
	public function listTrashbinFileInServerRoot(string $path): void {
		$response = OcsApiHelper::sendRequest(
			$this->getBaseUrl(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			'GET',
			"/apps/testing/api/v1/dir?dir=$path",
			$this->getStepLineRef()
		);
		$this->setResponse($response);
	}

	/**
	 * move file in server root
	 *
	 * @param string $path
	 * @param string $target
	 *
	 * @return void
	 */
	public function moveFileInServerRoot(string $path, string $target): void {
		$response = OcsApiHelper::sendRequest(
			$this->getBaseUrl(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			"MOVE",
			"/apps/testing/api/v1/file",
			$this->getStepLineRef(),
			[
				'source' => $path,
				'target' => $target
			]
		);

		$this->setResponse($response);
	}

	/**
	 * @Then the file :path with content :content should exist in the server root
	 *
	 * @param string $path
	 * @param string $content
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theFileWithContentShouldExistInTheServerRoot(string $path, string $content): void {
		$fileContent = $this->readFileInServerRootForOCIS($path);
		Assert::assertSame(
			$content,
			$fileContent,
			"The content of the file does not match with '$content'"
		);
	}

	/**
	 * @Then /^the content in the response should match with the content of file "([^"]*)" in the server root$/
	 *
	 * @param string $path
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theContentInTheRespShouldMatchWithFileInTheServerRoot(string $path): void {
		$content = $this->getResponse()->getBody()->getContents();
		$this->theFileWithContentShouldExistInTheServerRoot($path, $content);
	}

	/**
	 * @Then the file :path should not exist in the server root
	 *
	 * @param string $path
	 *
	 * @return void
	 */
	public function theFileShouldNotExistInTheServerRoot(string $path): void {
		$this->readFileInServerRootForCore($path);
		Assert::assertSame(
			404,
			$this->getResponse()->getStatusCode(),
			"The file '$path' exists in the server root but was not expected to exist"
		);
	}

	/**
	 * @Then the body of the response should be empty
	 *
	 * @return void
	 */
	public function theResponseBodyShouldBeEmpty(): void {
		Assert::assertEmpty(
			$this->getResponse()->getBody()->getContents(),
			"The response body was expected to be empty but got "
			. $this->getResponse()->getBody()->getContents()
		);
	}

	/**
	 * @param ResponseInterface|null $response
	 *
	 * @return array
	 */
	public function getJsonDecodedResponse(?ResponseInterface $response = null): array {
		if ($response === null) {
			$response = $this->getResponse();
		}
		return \json_decode(
			(string)$response->getBody(),
			true
		);
	}

	/**
	 *
	 * @return array
	 */
	public function getJsonDecodedStatusPhp(): array {
		return $this->getJsonDecodedResponse(
			$this->getStatusPhp()
		);
	}

	/**
	 * @return string
	 */
	public function getEditionFromStatus(): string {
		$decodedResponse = $this->getJsonDecodedStatusPhp();
		if (isset($decodedResponse['edition'])) {
			return $decodedResponse['edition'];
		}
		return '';
	}

	/**
	 * @return string|null
	 */
	public function getProductNameFromStatus(): ?string {
		$decodedResponse = $this->getJsonDecodedStatusPhp();
		if (isset($decodedResponse['productname'])) {
			return $decodedResponse['productname'];
		}
		return '';
	}

	/**
	 * @return string|null
	 */
	public function getVersionFromStatus(): ?string {
		$decodedResponse = $this->getJsonDecodedStatusPhp();
		if (isset($decodedResponse['version'])) {
			return $decodedResponse['version'];
		}
		return '';
	}

	/**
	 * @return string|null
	 */
	public function getVersionStringFromStatus(): ?string {
		$decodedResponse = $this->getJsonDecodedStatusPhp();
		if (isset($decodedResponse['versionstring'])) {
			return $decodedResponse['versionstring'];
		}
		return '';
	}

	/**
	 * returns a string that can be used to check a URL of comments with
	 * regular expression (without delimiter)
	 *
	 * @return string
	 */
	public function getCommentUrlRegExp(): string {
		$basePath = \ltrim($this->getBasePath() . "/", "/");
		return "/{$basePath}remote.php/dav/comments/files/([0-9]+)";
	}

	/**
	 * substitutes codes like %base_url% with the value
	 * if the given value does not have anything to be substituted
	 * then it is returned unmodified
	 *
	 * @param string|null $value
	 * @param string|null $user
	 * @param array|null $functions associative array of functions and parameters to be
	 *                              called on every replacement string before the
	 *                              replacement
	 *                              function name has to be the key and the parameters an
	 *                              own array
	 *                              the replacement itself will be used as first parameter
	 *                              e.g. substituteInLineCodes($value, ['preg_quote' => ['/']])
	 * @param array|null $additionalSubstitutions
	 *                         array of additional substitution configurations
	 *                           [
	 *                             [
	 *                               "code" => "%my_code%",
	 *                               "function" => [
	 *                                                $myClass,
	 *                                                "myFunction"
	 *                               ],
	 *                               "parameter" => []
	 *                             ],
	 *                           ]
	 * @param string|null $group
	 * @param string|null $userName
	 *
	 * @return string
	 */
	public function substituteInLineCodes(
		?string $value,
		?string $user = null,
		?array  $functions = [],
		?array  $additionalSubstitutions = [],
		?string $group = null,
		?string $userName = null
	): ?string {
		$substitutions = [
			[
				"code" => "%base_url%",
				"function" => [
					$this,
					"getBaseUrl"
				],
				"parameter" => []
			],
			[
				"code" => "%base_url_without_scheme%",
				"function" => [
					$this,
					"getBaseUrlWithoutScheme"
				],
				"parameter" => []
			],
			[
				"code" => "%remote_server%",
				"function" => [
					$this,
					"getRemoteBaseUrl"
				],
				"parameter" => []
			],
			[
				"code" => "%remote_server_without_scheme%",
				"function" => [
					$this,
					"getRemoteBaseUrlWithoutScheme"
				],
				"parameter" => []
			],
			[
				"code" => "%local_server%",
				"function" => [
					$this,
					"getLocalBaseUrl"
				],
				"parameter" => []
			],
			[
				"code" => "%local_server_without_scheme%",
				"function" => [
					$this,
					"getLocalBaseUrlWithoutScheme"
				],
				"parameter" => []
			],
			[
				"code" => "%base_path%",
				"function" => [
					$this,
					"getBasePath"
				],
				"parameter" => []
			],
			[
				"code" => "%dav_path%",
				"function" => [
					$this,
					"getDAVPathIncludingBasePath"
				],
				"parameter" => []
			],
			[
				"code" => "%ocs_path_v1%",
				"function" => [
					$this,
					"getOCSPath"
				],
				"parameter" => ["1"]
			],
			[
				"code" => "%ocs_path_v2%",
				"function" => [
					$this,
					"getOCSPath"
				],
				"parameter" => ["2"]
			],
			[
				"code" => "%productname%",
				"function" => [
					$this,
					"getProductNameFromStatus"
				],
				"parameter" => []
			],
			[
				"code" => "%edition%",
				"function" => [
					$this,
					"getEditionFromStatus"
				],
				"parameter" => []
			],
			[
				"code" => "%version%",
				"function" => [
					$this,
					"getVersionFromStatus"
				],
				"parameter" => []
			],
			[
				"code" => "%versionstring%",
				"function" => [
					$this,
					"getVersionStringFromStatus"
				],
				"parameter" => []
			],
			[
				"code" => "%a_comment_url%",
				"function" => [
					$this,
					"getCommentUrlRegExp"
				],
				"parameter" => []
			],
			[
				"code" => "%last_share_id%",
				"function" => [
					$this,
					"getLastCreatedUserGroupShareId"
				],
				"parameter" => []
			],
			[
				"code" => "%last_public_share_token%",
				"function" => [
					$this,
					"getLastCreatedPublicShareToken"
				],
				"parameter" => []
			],
			[
				"code" => "%group_id_pattern%",
				"function" => [
					__NAMESPACE__ . '\TestHelpers\GraphHelper',
					"getUUIDv4Regex"
				],
				"parameter" => []
			],
			[
				"code" => "%space_id_pattern%",
				"function" => [
					__NAMESPACE__ . '\TestHelpers\GraphHelper',
					"getSpaceIdRegex"
				],
				"parameter" => []
			],
			[
				"code" => "%user_id_pattern%",
				"function" => [
					__NAMESPACE__ . '\TestHelpers\GraphHelper',
					"getUUIDv4Regex"
				],
				"parameter" => []
			],
			[
				"code" => "%user_id%",
				"function" => [
					$this, "getUserIdByUserName"
				],
					"parameter" => [$userName]
				],
			[
				"code" => "%group_id%",
				"function" => [
					$this, "getGroupIdByGroupName"
				],
				"parameter" => [$group]
			]
		];
		if ($user !== null) {
			array_push(
				$substitutions,
				[
					"code" => "%username%",
					"function" => [
						$this,
						"getActualUsername"
					],
					"parameter" => [$user]
				],
				[
					"code" => "%displayname%",
					"function" => [
						$this,
						"getDisplayNameForUser"
					],
					"parameter" => [$user]
				],
				[
					"code" => "%password%",
					"function" => [
						$this,
						"getPasswordForUser"
					],
					"parameter" => [$user]
				],
				[
					"code" => "%emailaddress%",
					"function" => [
						$this,
						"getEmailAddressForUser"
					],
					"parameter" => [$user]
				],
				[
					"code" => "%spaceid%",
					"function" => [
						$this,
						"getPersonalSpaceIdForUser",
					],
					"parameter" => [$user, true]
				],
				[
				"code" => "%user_id%",
				"function" =>
				[$this, "getUserIdByUserName"],
				"parameter" => [$userName]
				],
				[
				"code" => "%group_id%",
				"function" =>
				[$this, "getGroupIdByGroupName"],
				"parameter" => [$group]
				]
			);
		}

		if (!empty($additionalSubstitutions)) {
			$substitutions = \array_merge($substitutions, $additionalSubstitutions);
		}

		foreach ($substitutions as $substitution) {
			if (strpos($value, $substitution['code']) === false) {
				continue;
			}

			$replacement = \call_user_func_array(
				$substitution["function"],
				$substitution["parameter"]
			);
			foreach ($functions as $function => $parameters) {
				$replacement = \call_user_func_array(
					$function,
					\array_merge([$replacement], $parameters)
				);
			}
			$value = \str_replace(
				$substitution["code"],
				$replacement,
				$value
			);
		}
		return $value;
	}

	/**
	 * returns personal space id for user if the test is using the spaces dav path
	 * or if alwaysDoIt is set to true,
	 * otherwise it returns null.
	 *
	 * @param string $user
	 * @param bool $alwaysDoIt default false. Set to true
	 *
	 * @return string|null
	 * @throws GuzzleException
	 */
	public function getPersonalSpaceIdForUser(string $user, bool $alwaysDoIt = false): ?string {
		if ($alwaysDoIt || ($this->getDavPathVersion() === WebDavHelper::DAV_VERSION_SPACES)) {
			return WebDavHelper::getPersonalSpaceIdForUserOrFakeIfNotFound(
				$this->getBaseUrl(),
				$user,
				$this->getPasswordForUser($user),
				$this->getStepLineRef()
			);
		}
		return null;
	}

	/**
	 * @return string
	 */
	public function temporaryStorageSubfolderName(): string {
		return "work_tmp";
	}

	/**
	 * @return string
	 */
	public function acceptanceTestsDirLocation(): string {
		return \dirname(__FILE__) . "/../../";
	}

	/**
	 * @return string
	 */
	public function workStorageDirLocation(): string {
		return $this->acceptanceTestsDirLocation() . $this->temporaryStorageSubfolderName() . "/";
	}

	/**
	 * Get the path of the ownCloud server root directory
	 *
	 * @return string
	 * @throws Exception
	 */
	public function getServerRoot(): string {
		if ($this->localServerRoot === null) {
			$this->localServerRoot = SetupHelper::getServerRoot(
				$this->getBaseUrl(),
				$this->getAdminUsername(),
				$this->getAdminPassword(),
				$this->getStepLineRef()
			);
		}
		return $this->localServerRoot;
	}

	/**
	 * @Then the config key :key of app :appID should have value :value
	 *
	 * @param string $key
	 * @param string $appID
	 * @param string $value
	 *
	 * @return void
	 * @throws Exception
	 */
	public function theConfigKeyOfAppShouldHaveValue(string $key, string $appID, string $value): void {
		$response = OcsApiHelper::sendRequest(
			$this->getBaseUrl(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			'GET',
			"/apps/testing/api/v1/app/$appID/$key",
			$this->getStepLineRef(),
			[],
			$this->getOcsApiVersion()
		);
		$configkeyValue = (string)$this->getResponseXml($response, __METHOD__)->data[0]->element->value;
		Assert::assertEquals(
			$value,
			$configkeyValue,
			"The config key $key of app $appID was expected to have value $value but got $configkeyValue"
		);
	}

	/**
	 * Parse list of config keys from the given XML response
	 *
	 * @param SimpleXMLElement $responseXml
	 *
	 * @return array
	 */
	public function parseConfigListFromResponseXml(SimpleXMLElement $responseXml): array {
		$configkeyData = \json_decode(\json_encode($responseXml->data), true);
		if (isset($configkeyData['element'])) {
			$configkeyData = $configkeyData['element'];
		} else {
			// There are no keys for the app
			return [];
		}
		if (isset($configkeyData[0])) {
			$configkeyValues = $configkeyData;
		} else {
			// There is just 1 key for the app
			$configkeyValues[0] = $configkeyData;
		}
		return $configkeyValues;
	}

	/**
	 * Returns a list of config keys for the given app
	 *
	 * @param string $appID
	 * @param string $exceptionText text to put at the front of exception messages
	 *
	 * @return array
	 * @throws Exception
	 */
	public function getConfigKeyList(string $appID, string $exceptionText = ''): array {
		if ($exceptionText === '') {
			$exceptionText = __METHOD__;
		}
		$response = OcsApiHelper::sendRequest(
			$this->getBaseUrl(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			'GET',
			"/apps/testing/api/v1/app/$appID",
			$this->getStepLineRef(),
			[],
			$this->getOcsApiVersion()
		);
		return $this->parseConfigListFromResponseXml(
			$this->getResponseXml($response, $exceptionText)
		);
	}

	/**
	 * Check if given config key is present for given app
	 *
	 * @param string $key
	 * @param string $appID
	 *
	 * @return bool
	 * @throws Exception
	 */
	public function checkConfigKeyInApp(string $key, string $appID): bool {
		$configkeyList = $this->getConfigKeyList($appID);
		foreach ($configkeyList as $config) {
			if ($config['configkey'] === $key) {
				return true;
			}
		}
		return false;
	}

	/**
	 * @Then /^app ((?:'[^']*')|(?:"[^"]*")) should (not|)\s?have config key ((?:'[^']*')|(?:"[^"]*"))$/
	 *
	 * @param string $appID
	 * @param string $shouldOrNot
	 * @param string $key
	 *
	 * @return void
	 * @throws Exception
	 */
	public function appShouldHaveConfigKey(string $appID, string $shouldOrNot, string $key): void {
		$appID = \trim($appID, $appID[0]);
		$key = \trim($key, $key[0]);

		$should = ($shouldOrNot !== "not");

		if ($should) {
			Assert::assertTrue(
				$this->checkConfigKeyInApp($key, $appID),
				"App $appID does not have config key $key"
			);
		} else {
			Assert::assertFalse(
				$this->checkConfigKeyInApp($key, $appID),
				"App $appID has config key $key but was not expected to"
			);
		}
	}

	/**
	 * @Then /^following config keys should (not|)\s?exist$/
	 *
	 * @param string $shouldOrNot
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws Exception
	 */
	public function followingConfigKeysShouldExist(string $shouldOrNot, TableNode $table): void {
		$should = ($shouldOrNot !== "not");
		if ($should) {
			foreach ($table as $item) {
				Assert::assertTrue(
					$this->checkConfigKeyInApp($item['configkey'], $item['appid']),
					"{$item['appid']} was expected to have config key {$item['configkey']} but does not"
				);
			}
		} else {
			foreach ($table as $item) {
				Assert::assertFalse(
					$this->checkConfigKeyInApp($item['configkey'], $item['appid']),
					"Expected : {$item['appid']} should not have config key {$item['configkey']}"
				);
			}
		}
	}

	/**
	 * @param string $user
	 * @param string|null $asUser
	 * @param string|null $password
	 *
	 * @return void
	 */
	public function sendUserSyncRequest(string $user, ?string $asUser = null, ?string $password = null): void {
		$user = $this->getActualUsername($user);
		$asUser = $asUser ?? $this->getAdminUsername();
		$password = $password ?? $this->getPasswordForUser($asUser);
		$response = OcsApiHelper::sendRequest(
			$this->getBaseUrl(),
			$asUser,
			$password,
			'POST',
			"/cloud/user-sync/$user",
			$this->getStepLineRef(),
			[],
			$this->getOcsApiVersion()
		);
		$this->setResponse($response);
	}

	/**
	 * @When the administrator tries to sync user :user using the OCS API
	 *
	 * @param string $user
	 *
	 * @return void
	 */
	public function theAdministratorTriesToSyncUserUsingTheOcsApi(string $user): void {
		$this->sendUserSyncRequest($user);
	}

	/**
	 * @When user :asUser tries to sync user :user using the OCS API
	 *
	 * @param string $asUser
	 * @param string $user
	 *
	 * @return void
	 */
	public function userTriesToSyncUserUsingTheOcsApi(string $asUser, string $user): void {
		$asUser = $this->getActualUsername($asUser);
		$user = $this->getActualUsername($user);
		$this->sendUserSyncRequest($user, $asUser);
	}

	/**
	 * @When the administrator tries to sync user :user using password :password and the OCS API
	 *
	 * @param string|null $user
	 * @param string|null $password
	 *
	 * @return void
	 */
	public function theAdministratorTriesToSyncUserUsingPasswordAndTheOcsApi(?string $user, ?string $password): void {
		$this->sendUserSyncRequest($user, null, $password);
	}

	/**
	 * This will run before EVERY scenario.
	 * It will set the properties for this object.
	 *
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 * @throws Exception
	 */
	public function before(BeforeScenarioScope $scope): void {
		$this->scenarioStartTime = \time();
		// Get the environment
		$environment = $scope->getEnvironment();
		// registers context in every suite, as every suite has FeatureContext
		// that calls BasicStructure.php
		$this->ocsContext = new OCSContext();
		$this->authContext = new AuthContext();
		$this->ocsContext->before($scope);
		$this->authContext->setUpScenario($scope);
		$environment->registerContext($this->ocsContext);
		$environment->registerContext($this->authContext);
		$scenarioLine = $scope->getScenario()->getLine();
		$featureFile = $scope->getFeature()->getFile();
		$suiteName = $scope->getSuite()->getName();
		$featureFileName = \basename($featureFile);

		if (!OcisHelper::isTestingOnReva()) {
			$this->spacesContext = new SpacesContext();
			$this->spacesContext->setUpScenario($scope);
			$environment->registerContext($this->spacesContext);
		}

		if ($this->sendScenarioLineReferencesInXRequestId()) {
			$this->scenarioString = $suiteName . '/' . $featureFileName . ':' . $scenarioLine;
		} else {
			$this->scenarioString = '';
		}

		// Initialize SetupHelper
		SetupHelper::init(
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			$this->getBaseUrl(),
			$this->getOcPath()
		);

		if ($this->isTestingWithLdap()) {
			$suiteParameters = SetupHelper::getSuiteParameters($scope);
			$this->connectToLdap($suiteParameters);
		}

		if (OcisHelper::isTestingWithGraphApi()) {
			$this->graphContext = new GraphContext();
			$this->graphContext->before($scope);
			$environment->registerContext($this->graphContext);
		}
	}

	/**
	 * This will run before EVERY step.
	 *
	 * @BeforeStep
	 *
	 * @param BeforeStepScope $scope
	 *
	 * @return void
	 */
	public function beforeEachStep(BeforeStepScope $scope): void {
		if ($this->sendScenarioLineReferencesInXRequestId()) {
			$this->stepLineRef = $this->scenarioString . '-' . $scope->getStep()->getLine();
		} else {
			$this->stepLineRef = '';
		}
	}

	/**
	 * @AfterScenario
	 *
	 * @return void
	 */
	public function restoreAdminPassword(): void {
		if ($this->adminPassword !== $this->originalAdminPassword) {
			$this->resetUserPasswordAsAdminUsingTheProvisioningApi(
				$this->getAdminUsername(),
				$this->originalAdminPassword
			);
			$this->adminPassword = $this->originalAdminPassword;
		}
	}

	/**
	 * @AfterScenario
	 *
	 * @return void
	 */
	public function deleteAllResourceCreatedByAdmin(): void {
		foreach ($this->adminResources as $resource) {
			$this->deleteFile("admin", $resource);
		}
	}

	/**
	 * @BeforeScenario @temporary_storage_on_server
	 *
	 * @return void
	 * @throws Exception
	 */
	public function makeTemporaryStorageOnServerBefore(): void {
		$this->mkDirOnServer(
			TEMPORARY_STORAGE_DIR_ON_REMOTE_SERVER
		);
	}

	/**
	 * @AfterScenario @temporary_storage_on_server
	 *
	 * @return void
	 * @throws Exception
	 */
	public function removeTemporaryStorageOnServerAfter(): void {
		SetupHelper::rmDirOnServer(
			TEMPORARY_STORAGE_DIR_ON_REMOTE_SERVER,
			$this->getStepLineRef()
		);
	}

	/**
	 * @AfterScenario
	 *
	 * @return void
	 */
	public function removeCreatedFilesAfter(): void {
		foreach ($this->createdFiles as $file) {
			\unlink($file);
		}
	}

	/**
	 *
	 * @param string $serverUrl
	 *
	 * @return void
	 */
	public function clearFileLocksForServer(string $serverUrl): void {
		$response = OcsApiHelper::sendRequest(
			$serverUrl,
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			'delete',
			"/apps/testing/api/v1/lockprovisioning",
			$this->getStepLineRef(),
			["global" => "true"]
		);
		Assert::assertEquals("200", $response->getStatusCode());
	}

	/**
	 * @AfterScenario
	 *
	 * clear space id reference
	 *
	 * @return void
	 * @throws Exception
	 */
	public function clearSpaceId(): void {
		if (\count(WebDavHelper::$spacesIdRef) > 0) {
			WebDavHelper::$spacesIdRef = [];
		}
		WebDavHelper::$SPACE_ID_FROM_OCIS = '';
	}

	/**
	 * @BeforeSuite
	 *
	 * @param BeforeSuiteScope $scope
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function useBigFileIDs(BeforeSuiteScope $scope): void {
		return;
	}

	/**
	 * Verify that the tableNode contains expected headers
	 *
	 * @param TableNode|null $table
	 * @param array|null $requiredHeader
	 * @param array|null $allowedHeader
	 *
	 * @return void
	 * @throws Exception
	 */
	public function verifyTableNodeColumns(?TableNode $table, ?array $requiredHeader = [], ?array $allowedHeader = []): void {
		if ($table === null || \count($table->getHash()) < 1) {
			throw new Exception("Table should have at least one row.");
		}
		$tableHeaders = $table->getRows()[0];
		$allowedHeader = \array_unique(\array_merge($requiredHeader, $allowedHeader));
		if ($requiredHeader != []) {
			foreach ($requiredHeader as $element) {
				if (!\in_array($element, $tableHeaders)) {
					throw new Exception("Row with header '$element' expected to be in table but not found");
				}
			}
		}

		if ($allowedHeader != []) {
			foreach ($tableHeaders as $element) {
				if (!\in_array($element, $allowedHeader)) {
					throw new Exception("Row with header '$element' is not allowed in table but found");
				}
			}
		}
	}

	/**
	 * Verify that the tableNode contains expected rows
	 *
	 * @param TableNode $table
	 * @param array $requiredRows
	 * @param array $allowedRows
	 *
	 * @return void
	 * @throws Exception
	 */
	public function verifyTableNodeRows(TableNode $table, array $requiredRows = [], array $allowedRows = []): void {
		if (\count($table->getRows()) < 1) {
			throw new Exception("Table should have at least one row.");
		}
		$tableHeaders = $table->getColumn(0);
		$allowedRows = \array_unique(\array_merge($requiredRows, $allowedRows));
		if ($requiredRows != []) {
			foreach ($requiredRows as $element) {
				if (!\in_array($element, $tableHeaders)) {
					throw new Exception("Row with name '$element' expected to be in table but not found");
				}
			}
		}

		if ($allowedRows != []) {
			foreach ($tableHeaders as $element) {
				if (!\in_array($element, $allowedRows)) {
					throw new Exception("Row with name '$element' is not allowed in table but found");
				}
			}
		}
	}

	/**
	 * Verify that the tableNode contains expected number of columns
	 *
	 * @param TableNode $table
	 * @param int $count
	 *
	 * @return void
	 * @throws Exception
	 */
	public function verifyTableNodeColumnsCount(TableNode $table, int $count): void {
		if (\count($table->getRows()) < 1) {
			throw new Exception("Table should have at least one row.");
		}
		$rowCount = \count($table->getRows()[0]);
		if ($count !== $rowCount) {
			throw new Exception("Table expected to have $count rows but found $rowCount");
		}
	}

	/**
	 * @return void
	 */
	public function resetAppConfigs(): void {
		// Set the required starting values for testing
		$this->setCapabilities($this->getCommonSharingConfigs());
	}

	/**
	 * @Given the administrator has set the last login date for user :user to :days days ago
	 * @When the administrator sets the last login date for user :user to :days days ago using the testing API
	 *
	 * @param string $user
	 * @param string $days
	 *
	 * @return void
	 */
	public function theAdministratorSetsTheLastLoginDateForUserToDaysAgoUsingTheTestingApi(string $user, string $days): void {
		$user = $this->getActualUsername($user);
		$adminUser = $this->getAdminUsername();
		$baseUrl = "/apps/testing/api/v1/lastlogindate/$user";
		$response = OcsApiHelper::sendRequest(
			$this->getBaseUrl(),
			$adminUser,
			$this->getAdminPassword(),
			'POST',
			$baseUrl,
			$this->getStepLineRef(),
			['days' => $days],
			$this->getOcsApiVersion()
		);
		$this->setResponse($response);
	}

	/**
	 * @param array $capabilitiesArray with each array entry containing keys for:
	 *                                 ['capabilitiesApp'] the "app" name in the capabilities response
	 *                                 ['capabilitiesParameter'] the parameter name in the capabilities response
	 *                                 ['testingApp'] the "app" name as understood by "testing"
	 *                                 ['testingParameter'] the parameter name as understood by "testing"
	 *                                 ['testingState'] boolean state the parameter must be set to for the test
	 *
	 * @return void
	 */
	public function setCapabilities(array $capabilitiesArray): void {
		AppConfigHelper::setCapabilities(
			$this->getBaseUrl(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			$capabilitiesArray,
			$this->getStepLineRef()
		);
	}

	/**
	 * @param string $sourceUser
	 * @param string $targetUser
	 *
	 * @return string|null
	 * @throws Exception
	 */
	public function findLastTransferFolderForUser(string $sourceUser, string $targetUser): ?string {
		$foundPaths = [];
		$responseXmlObject = $this->listFolderAndReturnResponseXml(
			$targetUser,
			'',
			'1'
		);
		$transferredElements = $responseXmlObject->xpath(
			"//d:response/d:href[contains(., '/transferred%20from%20$sourceUser%20on%')]"
		);
		foreach ($transferredElements as $transferredElement) {
			// $transferredElement is an XML object. We want to work with the string in the XML element.
			$path = \rawurldecode((string)$transferredElement);
			$parts = \explode(' ', $path);
			// store timestamp as key
			$foundPaths[] = [
				'date' => \strtotime(\trim($parts[4], '/')),
				'path' => $path,
			];
		}
		if (empty($foundPaths)) {
			return null;
		}

		\usort(
			$foundPaths,
			function ($a, $b) {
				return $a['date'] - $b['date'];
			}
		);

		$davPath = \rtrim($this->getFullDavFilesPath($targetUser), '/');

		$foundPath = \end($foundPaths)['path'];
		// strip DAV path
		return \substr($foundPath, \strlen($davPath) + 1);
	}

	/**
	 * Get the array of trusted servers in format ["url" => "id"]
	 *
	 * @param string $server 'LOCAL'/'REMOTE'
	 *
	 * @return array
	 * @throws Exception
	 */
	public function getTrustedServers(string $server = 'LOCAL'): array {
		if ($server === 'LOCAL') {
			$url = $this->getLocalBaseUrl();
		} elseif ($server === 'REMOTE') {
			$url = $this->getRemoteBaseUrl();
		} else {
			throw new Exception(__METHOD__ . " Invalid value for server : $server");
		}
		$adminUser = $this->getAdminUsername();
		$response = OcsApiHelper::sendRequest(
			$url,
			$adminUser,
			$this->getAdminPassword(),
			'GET',
			"/apps/testing/api/v1/trustedservers",
			$this->getStepLineRef()
		);
		if ($response->getStatusCode() !== 200) {
			throw new Exception("Could not get the list of trusted servers" . $response->getBody()->getContents());
		}
		$responseXml = HttpRequestHelper::getResponseXml(
			$response,
			__METHOD__
		);
		$serverData = \json_decode(
			\json_encode(
				$responseXml->data
			),
			true
		);
		if (!\array_key_exists('element', $serverData)) {
			return [];
		} else {
			return isset($serverData['element'][0]) ?
				\array_column($serverData['element'], 'id', 'url') :
				\array_column($serverData, 'id', 'url');
		}
	}

	/**
	 * @param string $method http request method
	 * @param string $property property in form d:getetag
	 *                         if property is `doesnotmatter` body is also set `doesnotmatter`
	 *
	 * @return string
	 */
	public function getBodyForOCSRequest(string $method, string $property): ?string {
		$body = null;
		if ($method === 'PROPFIND') {
			$body = '<?xml version="1.0"?><d:propfind  xmlns:d="DAV:" xmlns:oc="http://owncloud.org/ns"><d:prop><' . $property . '/></d:prop></d:propfind>';
		} elseif ($method === 'LOCK') {
			$body = "<?xml version='1.0' encoding='UTF-8'?><d:lockinfo xmlns:d='DAV:'> <d:lockscope><" . $property . " /></d:lockscope></d:lockinfo>";
		} elseif ($method === 'PROPPATCH') {
			if ($property === 'favorite') {
				$property = '<oc:favorite xmlns:oc="http://owncloud.org/ns">1</oc:favorite>';
			}
			$body = '<?xml version="1.0"?><d:propertyupdate xmlns:d="DAV:" xmlns:oc="http://owncloud.org/ns"><d:set><d:prop>' . $property . '</d:prop></d:set></d:propertyupdate>';
		}
		if ($property === '') {
			$body = '';
		}
		return $body;
	}

	/**
	 * The method returns userId
	 *
	 * @param string $userName
	 *
	 * @return string
	 * @throws Exception|GuzzleException
	 */
	public function getUserIdByUserName(string $userName): string {
		$response = GraphHelper::getUser(
			$this->getBaseUrl(),
			$this->getStepLineRef(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			$userName
		);
		$data = \json_decode($response->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
		if (isset($data["id"])) {
			return $data["id"];
		} else {
			throw new Exception(__METHOD__ . " accounts-list is empty");
		}
	}

	/**
	 * The method returns groupId
	 *
	 * @param string $groupName
	 *
	 * @return string
	 * @throws Exception|GuzzleException
	 */
	public function getGroupIdByGroupName(string $groupName):string {
		$response = GraphHelper::getGroup(
			$this->getBaseUrl(),
			$this->getStepLineRef(),
			$this->getAdminUsername(),
			$this->getAdminPassword(),
			$groupName
		);
		$data = $this->getJsonDecodedResponse($response);
		if (isset($data["id"])) {
			return $data["id"];
		} else {
			throw new Exception(__METHOD__ . " accounts-list is empty");
		}
	}

	/**
	 *
	 * @AfterSuite
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function clearScenarioLog(): void {
		if (\file_exists(HttpLogger::getScenarioLogPath())) {
			\unlink(HttpLogger::getScenarioLogPath());
		}
	}

	/**
	 * Log request and response logs if scenario fails
	 *
	 * @AfterScenario
	 *
	 * @param AfterScenarioScope $scope
	 *
	 * @return void
	 * @throws Exception
	 */
	public static function checkScenario(AfterScenarioScope $scope): void {
		if (($scope->getTestResult()->getResultCode() !== 0)
			&& (!self::isExpectedToFail(self::getScenarioLine($scope)))
		) {
			$logs = \file_get_contents(HttpLogger::getScenarioLogPath());
			// add new lines
			$logs = \rtrim($logs, "\n") . "\n\n\n";
			HttpLogger::writeLog(HttpLogger::getFailedLogPath(), $logs);
		}
	}

	/**
	 * @param BeforeScenarioScope|AfterScenarioScope $scope
	 *
	 * @return string
	 */
	public static function getScenarioLine($scope): string {
		$feature = $scope->getFeature()->getFile();
		$feature = \explode('/', $feature);
		$feature = \array_slice($feature, -2);
		$feature = \implode('/', $feature);
		$scenarioLine = $scope->getScenario()->getLine();
		// Example: apiGraph/createUser.feature:24
		return $feature . ':' . $scenarioLine;
	}

	/**
	 * @param string $scenarioLine
	 *
	 * @return bool
	 */
	public static function isExpectedToFail(string $scenarioLine): bool {
		$expectedFailFile = \getenv('EXPECTED_FAILURES_FILE');
		if (!$expectedFailFile) {
			$expectedFailFile = __DIR__ . '/../../expected-failures-localAPI-on-OCIS-storage.md';
			if (\strpos($scenarioLine, "coreApi") === 0) {
				$expectedFailFile = __DIR__ . '/../../expected-failures-API-on-OCIS-storage.md';
			}
		}

		$reader = \fopen($expectedFailFile, 'r');
		if ($reader) {
			while (($line = \fgets($reader)) !== false) {
				if (\strpos($line, $scenarioLine) !== false) {
					\fclose($reader);
					return true;
				}
			}
			\fclose($reader);
		}
		return false;
	}
}
