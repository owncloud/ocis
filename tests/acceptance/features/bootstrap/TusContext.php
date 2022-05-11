<?php

declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Viktor Scharf <v.scharf@owncloud.com>
 * @copyright Copyright (c) 2022 Viktor Scharf v.scharf@owncloud.com
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use GuzzleHttp\Exception\GuzzleException;
use TestHelpers\HttpRequestHelper;

require_once 'bootstrap.php';

/**
 * Context for the provisioning specific steps using the Graph API
 */
class TusContext implements Context
{

    /**
     * @var FeatureContext
     */
    private FeatureContext $featureContext;

    /**
     * @var SpacesContext
     */
    private SpacesContext $spacesContext;

    /**
     * @var string
     */
    private string $baseUrl;

    /**
     * @return string
     */
    public function acceptanceTestsDirLocation(): string
    {
        return \dirname(__FILE__) . "/../../filesForUpload/";
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
     */
    public function before(BeforeScenarioScope $scope): void
    {
        // Get the environment
        $environment = $scope->getEnvironment();
        // Get all the contexts you need in this context from here
        $this->featureContext = $environment->getContext('FeatureContext');
        $this->spacesContext = $environment->getContext('SpacesContext');
        $this->baseUrl = \trim($this->featureContext->getBaseUrl(), "/");
    }

    /**
     * @Given /^user "([^"]*)" has uploaded a file "([^"]*)" via TUS inside of the space "([^"]*)" using  WebDAV API$/
     *
     * @param string $user
     * @param string $resource
     * @param string $spaceName
     *
     * @return void
     *
     * @throws Exception
     * @throws GuzzleException
     */
    public function uploadFileViaTus(string $user, string $resource, string $spaceName): void
    {
        $resourceLocation = $this->getResourceLocation($user, $resource, $spaceName);
        $file = \fopen($this->acceptanceTestsDirLocation() . $resource, 'r');
        
        $this->featureContext->setResponse(
            HttpRequestHelper::sendRequest(
                $resourceLocation,
                "",
                'HEAD',
                $user,
                $this->featureContext->getPasswordForUser($user),
                [],
                ""
            )
        );
        $this->featureContext->theHTTPStatusCodeShouldBe(200, "Expected response status code should be 200");

        
        $this->featureContext->setResponse(
            HttpRequestHelper::sendRequest(
                $resourceLocation,
                "",
                'PATCH',
                $user,
                $this->featureContext->getPasswordForUser($user),
                ["Tus-Resumable" => "1.0.0", "Upload-Offset" => 0, 'Content-Type' => 'application/offset+octet-stream'],
                $file
            )
        );
        $this->featureContext->theHTTPStatusCodeShouldBe(204, "Expected response status code should be 204");
    }

    /**
     * send POST and return the url of the resource location in the response header
     * 
     * @param string $user
     * @param string $resource
     * @param string $spaceName
     *
     * @return string
     */
    public function getResourceLocation(string $user, string $resource, string $spaceName): string
    {
        $space = $this->spacesContext->getSpaceByName($user, $spaceName);
        $fullUrl = $this->baseUrl . "/remote.php/dav/spaces/" . $space["id"];

        $tusEndpoint = "tusEndpoint " . base64_encode(str_replace("$", "%", $fullUrl));
        $fileName = "filename " . base64_encode($resource);

        $headers = [
            "Tus-Resumable" => "1.0.0",
            "Upload-Metadata" => $tusEndpoint . ',' . $fileName,
            "Upload-Length" => filesize($this->acceptanceTestsDirLocation() . $resource)
        ];

        $this->featureContext->setResponse(
            HttpRequestHelper::post(
                $fullUrl,
                "",
                $this->featureContext->getActualUsername($user),
                $this->featureContext->getUserPassword($user),
                $headers,
                ''
            )
        );
        $this->featureContext->theHTTPStatusCodeShouldBe(201, "Expected response status code should be 201");

        $locationHeader = $this->featureContext->getResponse()->getHeader('Location');
        if (\sizeof($locationHeader) > 0) {
            return $locationHeader[0];
        }
    }
}
