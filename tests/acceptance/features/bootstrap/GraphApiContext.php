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
     * @param string $spaceName
     * @param string $xRequestId
     * @param array $headers
     * @return ResponseInterface
     */
    public function sendCreateSpaceRequest(
        $baseUrl,
        $user,
        $password,
        string $spaceName,
        string $xRequestId = '',
        array $headers = []
    ): ResponseInterface
    {
        $fullUrl = $baseUrl;
        if (!str_ends_with($fullUrl, '/')) {
            $fullUrl .= '/';
        }
        $fullUrl .= "drives/" . $spaceName;

        return HttpRequestHelper::sendRequest($fullUrl, $xRequestId, 'POST', $user, $password, $headers);
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
     * Get the webDavUrl of the personal space has been found
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
