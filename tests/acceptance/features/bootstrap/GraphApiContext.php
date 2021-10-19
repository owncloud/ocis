<?php

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\HttpRequestHelper;
use TestHelpers\SetupHelper;
use PHPUnit\Framework\Assert;

require_once 'bootstrap.php';

/**
 * Context for GraphApi specific steps
 */
class GraphApiContext implements Context {

    /**
     * @var FeatureContext
     */
    private FeatureContext $featureContext;

    /**
     * @var string
     */
    private string $personalDriveWebDavUrl;

    /**
     * @return string
     */
    public function getPersonalDriveWebDavUrl(): string
    {
        return $this->personalDriveWebDavUrl;
    }

    /**
     * @param string $personalDriveWebDavUrl
     */
    public function setPersonalDriveWebDavUrl(string $personalDriveWebDavUrl): void
    {
        $this->personalDriveWebDavUrl = $personalDriveWebDavUrl;
    }
    /**
     * @BeforeScenario
     *
     * @param BeforeScenarioScope $scope
     *
     * @return void
     * @throws Exception
     */
    public function setUpScenario(BeforeScenarioScope $scope): void
    {
        // Get the environment
        $environment = $scope->getEnvironment();
        // Get all the contexts you need in this context
        $this->featureContext = $environment->getContext('FeatureContext');
        SetupHelper::init(
            $this->featureContext->getAdminUsername(),
            $this->featureContext->getAdminPassword(),
            $this->featureContext->getBaseUrl(),
            $this->featureContext->getOcPath()
        );
    }

    /**
     * Send Graph List Drives Request
     *
     * @param $baseUrl
     * @param $user
     * @param $password
     * @param $arguments
     * @param string $xRequestId
     * @param array $body
     * @param array $headers
     * @return ResponseInterface
     */
    public function listSpacesRequest(
        $baseUrl,
        $user,
        $password,
        $arguments,
        string $xRequestId = '',
        array $body = [],
        array $headers = []
    ) {
        $fullUrl = $baseUrl;
        if (!str_ends_with($fullUrl, '/')) {
            $fullUrl .= '/';
        }
        $fullUrl .= "graph/v1.0/me/drives/" . $arguments;

        return HttpRequestHelper::get($fullUrl, $xRequestId, $user, $password, $headers, $body);
    }

    /**
     * Send Graph List Drives Request
     *
     * @param $baseUrl
     * @param $user
     * @param $password
     * @param string $body
     * @param string $xRequestId
     * @param array $headers
     * @return ResponseInterface
     */
    public function sendCreateSpaceRequest(
        $baseUrl,
        $user,
        $password,
        string $body,
        string $xRequestId = '',
        array $headers = []
    ): ResponseInterface
    {
        $fullUrl = $baseUrl;
        if (!str_ends_with($fullUrl, '/')) {
            $fullUrl .= '/';
        }
        $fullUrl .= "graph/v1.0/drives/";

        return HttpRequestHelper::post($fullUrl, $xRequestId, $user, $password, $headers, $body);
    }

    /**
     * Send Propfind Request to Url
     *
     * @param $fullUrl
     * @param $user
     * @param $password
     * @param string $xRequestId
     * @param array $headers
     * @return ResponseInterface
     */
    public function sendPropfindRequestToUrl(
        $fullUrl,
        $user,
        $password,
        string $xRequestId = '',
        array $headers = []
    ): ResponseInterface
    {
        return HttpRequestHelper::sendRequest($fullUrl, $xRequestId, 'PROPFIND', $user, $password, $headers);
    }

    /**
     * @When /^user "([^"]*)" lists all available spaces via the GraphApi$/
     *
     * @param $user
     * @return void
     */
    public function theUserListsAllHisAvailableSpacesUsingTheGraphApi($user): void
    {
        $this->featureContext->setResponse(
            $this->listSpacesRequest(
                $this->featureContext->getBaseUrl(),
                $user,
                $this->featureContext->getPasswordForUser($user),
                "",
                ""
            )
        );
    }

    /**
     * @When /^user "([^"]*)" creates a space "([^"]*)" of type "([^"]*)" with the default quota using the GraphApi$/
     *
     * @param $user string
     * @param $spaceName string
     * @param $spaceType string
     *
     * @return void
     */
    public function theUserCreatesASpaceUsingTheGraphApi($user, $spaceName, $spaceType): void
    {
        $space = ["Name" => $spaceName, "driveType" => $spaceType];
        $body = json_encode($space);
        $this->featureContext->setResponse(
            $this->sendCreateSpaceRequest(
                $this->featureContext->getBaseUrl(),
                $user,
                $this->featureContext->getPasswordForUser($user),
                $body,
                ""
            )
        );
    }

    /**
     * @When /^the administrator gives "([^"]*)" the role "([^"]*)" using the settings api$/
     *
     * @param $user string
     * @param $role string
     *
     * @return void
     */
    public function theAdministratorGivesUserTheRole($user, $role): void
    {
        $admin = $this->featureContext->getAdminUsername();
        $password = $this->featureContext->getAdminPassword();
        $headers = [];

        $baseUrl = $this->featureContext->getBaseUrl();
        if (!str_ends_with($baseUrl, '/')) {
            $baseUrl .= '/';
        }
        // get the roles list first
        $fullUrl = $baseUrl . "api/v0/settings/roles-list";
        $this->featureContext->setResponse(HttpRequestHelper::post($fullUrl, "", $admin, $password, $headers, "{}"));
        $rawBody =  $this->featureContext->getResponse()->getBody()->getContents();
        $bundles = [];
        if (isset(\json_decode($rawBody,  true)["bundles"])) {
            $bundles = \json_decode($rawBody, true)["bundles"];
        }
        $roleToAssign = "";
        foreach($bundles as $bundle => $value) {
            // find the selected role
            if ($value["displayName"] === $role) {
                $roleToAssign = $value;
            }
        }
        Assert::assertNotEmpty($roleToAssign, "The selected role $role could not be found");

        // get the accounts list first
        $fullUrl = $baseUrl . "api/v0/accounts/accounts-list";
        $this->featureContext->setResponse(HttpRequestHelper::post($fullUrl, "", $admin, $password, $headers, "{}"));
        $rawBody =  $this->featureContext->getResponse()->getBody()->getContents();
        $accounts = [];
        if (isset(\json_decode($rawBody,  true)["accounts"])) {
            $accounts = \json_decode($rawBody, true)["accounts"];
        }
        $accountToChange = "";
        foreach($accounts as $account) {
            // find the selected user
            if ($account["preferredName"] === $user) {
                $accountToChange = $account;
            }
        }
        Assert::assertNotEmpty($accountToChange, "The seleted account $user does not exist");

        // set the new role
        $fullUrl = $baseUrl . "api/v0/settings/assignments-add";
        $body = json_encode(["account_uuid" => $accountToChange["id"], "role_id" => $roleToAssign["id"]]);

        $this->featureContext->setResponse(HttpRequestHelper::post($fullUrl, "", $admin, $password, $headers, $body));
        $rawBody =  $this->featureContext->getResponse()->getBody()->getContents();

        $assignment = [];
        if (isset(\json_decode($rawBody,  true)["assignment"])) {
            $assignment = \json_decode($rawBody, true)["assignment"];
        }

        Assert::assertEquals($accountToChange["id"], $assignment["accountUuid"]);
        Assert::assertEquals($roleToAssign["id"], $assignment["roleId"]);
    }


    /**
     * Get the webDavUrl of the personal space
     *
     * @return void
     */
    public function theWebDavUrlOfThePersonalSpaceHasBeenFound(): void
    {
        $rawBody =  $this->featureContext->getResponse()->getBody()->getContents();
        $drives = [];
        if (isset(\json_decode($rawBody,  true)["value"])) {
            $drives = \json_decode($rawBody, true)["value"];
        }

        Assert::assertArrayHasKey(0, $drives, "No drives were found on that endpoint");

        foreach($drives as $drive) {
            if (isset($drive["driveType"]) && $drive["driveType"] === "personal") {
                $this->setPersonalDriveWebDavUrl($drive["root"]["webDavUrl"]);

                Assert::assertNotEmpty(
                    $drive["root"]["webDavUrl"],
                    "The personal space attributes contain no webDavUrl"
                );
            }
        }
    }

    /**
     * @When /^user "([^"]*)" lists the content of the personal space root using the WebDav Api$/
     *
     * @param $user
     *
     * @return void
     */
    public function theUserListsTheContentOfAPersonalSpaceRootUsingTheWebDAvApi($user): void
    {
        $this->theUserListsAllHisAvailableSpacesUsingTheGraphApi($user);
        $this->theWebDavUrlOfThePersonalSpaceHasBeenFound();
        $this->featureContext->setResponse(
            $this->sendPropfindRequestToUrl(
                $this->getPersonalDriveWebDavUrl(),
                $user,
                $this->featureContext->getPasswordForUser($user),
                "",
                [],
                [],
                []
            )
        );
    }
}
