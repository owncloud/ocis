<?php

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Behat\Gherkin\Node\TableNode;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\HttpRequestHelper;
use TestHelpers\SetupHelper;
use PHPUnit\Framework\Assert;

require_once 'bootstrap.php';

/**
 * Context for ocis spaces specific steps
 */
class SpacesContext implements Context {

    /**
     * @var FeatureContext
     */
    private FeatureContext $featureContext;

    /**
     * @var array
     */
    private array $availableSpaces;

    /**
     * @return array
     */
    public function getAvailableSpaces(): array
    {
        return $this->availableSpaces;
    }

    /**
     * @param array $availableSpaces
     */
    public function setAvailableSpaces(array $availableSpaces): void
    {
        $this->availableSpaces = $availableSpaces;
    }

    /**
     * response content parsed from XML to an array
     *
     * @var array
     */
    private array $responseXml = [];

    /**
     * @return array
     */
    public function getResponseXml(): array
    {
        return $this->responseXml;
    }

    /**
     * @param array $responseXml
     */
    public function setResponseXml(array $responseXml): void
    {
        $this->responseXml = $responseXml;
    }

    /**
     * space id from last propfind request
     *
     * @var string
     */
    private string $responseSpaceId;

    /**
     * @param string $responseSpaceId
     */
    public function setResponseSpaceId(string $responseSpaceId): void
    {
        $this->responseSpaceId = $responseSpaceId;
    }

    /**
     * @return string
     */
    public function getResponseSpaceId(): string
    {
        return $this->responseSpaceId;
    }

    /**
     * Get SpaceId by Name
     *
     * @param $name string
     * @return string
     * @throws Exception
     */
    public function getSpaceIdByName(string $name): string
    {
        $response = json_decode($this->featureContext->getResponse()->getBody(), true);
        if (isset($response['name']) && $response['name'] === $name) {
            return $response["id"];
        }
        foreach ($response["value"] as $spaceCandidate) {
            if ($spaceCandidate['name'] === $name) {
                return $spaceCandidate["id"];
            }
        }
        throw new Exception(__METHOD__ . " space with name $name not found");
    }

    /**
     * Get Space Array by name
     *
     * @param string $name
     * @return array
     */
    public function getSpaceByName(string $name): array
    {
        $response = json_decode($this->featureContext->getResponse()->getBody(), true);
        $spaceAsArray = $response;
        if (isset($response['name']) && $response['name'] === $name) {
            return $response;
        }
        foreach ($spaceAsArray["value"] as $spaceCandidate) {
            if ($spaceCandidate['name'] === $name) {
                return $spaceCandidate;
            }
        }
        return [];
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
     * Send Graph List Spaces Request
     *
     * @param string $baseUrl
     * @param string $user
     * @param string $password
     * @param string $urlArguments
     * @param string $xRequestId
     * @param array $body
     * @param array $headers
     * @return ResponseInterface
     */
    public function listSpacesRequest(
        string $baseUrl,
        string $user,
        string $password,
        string $urlArguments,
        string $xRequestId = '',
        array  $body = [],
        array  $headers = []
    ): ResponseInterface {
        $fullUrl = $baseUrl;
        if (!str_ends_with($fullUrl, '/')) {
            $fullUrl .= '/';
        }
        $fullUrl .= "graph/v1.0/me/drives/" . $urlArguments;

        return HttpRequestHelper::get($fullUrl, $xRequestId, $user, $password, $headers, $body);
    }

    /**
     * Send Graph Create Space Request
     *
     * @param string $baseUrl
     * @param string $user
     * @param string $password
     * @param string $body
     * @param string $xRequestId
     * @param array $headers
     * @return ResponseInterface
     */
    public function sendCreateSpaceRequest(
        string $baseUrl,
        string $user,
        string $password,
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
     * @param string $fullUrl
     * @param string $user
     * @param string $password
     * @param string $xRequestId
     * @param array $headers
     * @return ResponseInterface
     */
    public function sendPropfindRequestToUrl(
        string $fullUrl,
        string $user,
        string $password,
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
    public function theUserListsAllHisAvailableSpacesUsingTheGraphApi(string $user): void
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
        $this->rememberTheAvailableSpaces();
    }

    /**
     * @When /^user "([^"]*)" creates a space "([^"]*)" of type "([^"]*)" with the default quota using the GraphApi$/
     *
     * @param string $user
     * @param string $spaceName
     * @param string $spaceType
     *
     * @return void
     */
    public function theUserCreatesASpaceUsingTheGraphApi(
        string $user,
        string $spaceName,
        string $spaceType): void
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
     * @When /^user "([^"]*)" creates a space "([^"]*)" of type "([^"]*)" with quota "([^"]*)" using the GraphApi$/
     *
     * @param string $user
     * @param string $spaceName
     * @param string $spaceType
     * @param int $quota
     *
     * @return void
     */
    public function theUserCreatesASpaceWithQuotaUsingTheGraphApi(
        string $user,
        string $spaceName,
        string $spaceType,
        int $quota): void
    {
        $space = ["Name" => $spaceName, "driveType" => $spaceType, "quota" => ["total" => $quota]];
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
     * @param string $user
     * @param string $role
     *
     * @return void
     */
    public function theAdministratorGivesUserTheRole(string $user, string $role): void
    {
        $admin = $this->featureContext->getAdminUsername();
        $password = $this->featureContext->getAdminPassword();
        $headers = [];
        $bundles = [];
        $accounts = [];
        $assignment = [];

        $baseUrl = $this->featureContext->getBaseUrl();
        if (!str_ends_with($baseUrl, '/')) {
            $baseUrl .= '/';
        }
        // get the roles list first
        $fullUrl = $baseUrl . "api/v0/settings/roles-list";
        $this->featureContext->setResponse(HttpRequestHelper::post($fullUrl, "", $admin, $password, $headers, "{}"));
        if ($this->featureContext->getResponse()) {
            $rawBody =  $this->featureContext->getResponse()->getBody()->getContents();
            if (isset(\json_decode($rawBody,  true)["bundles"])) {
                $bundles = \json_decode($rawBody, true)["bundles"];
            }
        }
        $roleToAssign = "";
        foreach($bundles as $value) {
            // find the selected role
            if ($value["displayName"] === $role) {
                $roleToAssign = $value;
            }
        }
        Assert::assertNotEmpty($roleToAssign, "The selected role $role could not be found");

        // get the accounts list first
        $fullUrl = $baseUrl . "api/v0/accounts/accounts-list";
        $this->featureContext->setResponse(HttpRequestHelper::post($fullUrl, "", $admin, $password, $headers, "{}"));
        if ($this->featureContext->getResponse()) {
            $rawBody = $this->featureContext->getResponse()->getBody()->getContents();
            if (isset(\json_decode($rawBody,  true)["accounts"])) {
                $accounts = \json_decode($rawBody, true)["accounts"];
            }
        }
        $accountToChange = "";
        foreach($accounts as $account) {
            // find the selected user
            if ($account["preferredName"] === $user) {
                $accountToChange = $account;
            }
        }
        Assert::assertNotEmpty($accountToChange, "The selected account $user does not exist");

        // set the new role
        $fullUrl = $baseUrl . "api/v0/settings/assignments-add";
        $body = json_encode(["account_uuid" => $accountToChange["id"], "role_id" => $roleToAssign["id"]]);

        $this->featureContext->setResponse(HttpRequestHelper::post($fullUrl, "", $admin, $password, $headers, $body));
        if ($this->featureContext->getResponse()) {
            $rawBody = $this->featureContext->getResponse()->getBody()->getContents();
            if (isset(\json_decode($rawBody,  true)["assignment"])) {
                $assignment = \json_decode($rawBody, true)["assignment"];
            }
        }

        Assert::assertEquals($accountToChange["id"], $assignment["accountUuid"]);
        Assert::assertEquals($roleToAssign["id"], $assignment["roleId"]);
    }


    /**
     * Remember the available Spaces
     *
     * @return void
     */
    public function rememberTheAvailableSpaces(): void
    {
        $rawBody =  $this->featureContext->getResponse()->getBody()->getContents();
        $drives = json_decode($rawBody,  true);
        if (isset($drives["value"])) {
            $drives = $drives["value"];
        }

        Assert::assertArrayHasKey(0, $drives, "No drives were found on that endpoint");
        $spaces = [];
        foreach($drives as $drive) {
            $spaces[$drive["name"]] = $drive;
        }
        $this->setAvailableSpaces($spaces);
        Assert::assertNotEmpty($spaces, "No spaces have been found");
    }

    /**
     * @When /^user "([^"]*)" lists the content of the space with the name "([^"]*)" using the WebDav Api$/
     *
     * @param string $user
     * @param string $name
     * @return void
     */
    public function theUserListsTheContentOfAPersonalSpaceRootUsingTheWebDAvApi(
        string $user,
        string $name
    ): void
    {
        $spaceId = $this->getAvailableSpaces()[$name]["id"];
        $spaceWebDavUrl = $this->getAvailableSpaces()[$name]["root"]["webDavUrl"];
        $this->featureContext->setResponse(
            $this->sendPropfindRequestToUrl(
                $spaceWebDavUrl,
                $user,
                $this->featureContext->getPasswordForUser($user),
                "",
                [],
            )
        );
        $this->setResponseSpaceId($spaceId);
        $this->setResponseXml(HttpRequestHelper::parseResponseAsXml($this->featureContext->getResponse())
        );
    }

    /**
     * @Then /^the (?:propfind|search) result of the space should (not|)\s?contain these (?:files|entries):$/
     *
     * @param string $shouldOrNot (not|)
     * @param TableNode $expectedFiles
     *
     * @return void
     * @throws Exception
     */
    public function thePropfindResultShouldContainEntries(
        string $shouldOrNot,
        TableNode $expectedFiles
    ):void {
        $this->propfindResultShouldContainEntries(
            $shouldOrNot,
            $expectedFiles,
        );
    }

    /**
     * @Then /^the json responded should contain a space "([^"]*)" with these key and value pairs:$/
     *
     * @param string $spaceName
     * @param TableNode $table
     *
     * @return void
     */
    public function jsonRespondedShouldContain(
        string $spaceName,
        TableNode $table
    ): void {
        $this->featureContext->verifyTableNodeColumns($table, ['key', 'value']);
        Assert::assertIsArray($spaceAsArray = $this->getSpaceByName($spaceName), "No space with name $spaceName found");
        foreach ($table->getHash() as $row) {
            // remember the original Space Array
            $original = $spaceAsArray;
            $row['value'] = $this->featureContext->substituteInLineCodes(
                $row['value'],
                $this->featureContext->getCurrentUser(),
                [],
                [
                    [
                        "code" => "%space_id%",
                        "function" =>
                            [$this, "getSpaceIdByName"],
                        "parameter" => ["$spaceName"]
                    ]
                ]
            );
            $segments = explode("@@@", $row["key"]);
            // traverse down in the array
            foreach ($segments as $segment) {
                $arrayKeyExists = array_key_exists($segment, $spaceAsArray);
                $key = $row["key"];
                Assert::assertTrue($arrayKeyExists, "The key $key does not exist on the response");
                if ($arrayKeyExists) {
                    $spaceAsArray = $spaceAsArray[$segment];
                }
            }
            Assert::assertEquals($row["value"], $spaceAsArray);
            // set the spaceArray to the point before traversing
            $spaceAsArray = $original;
        }
    }

    /**
     * @param string $shouldOrNot (not|)
     * @param TableNode $expectedFiles
     *
     * @return void
     * @throws Exception
     */
    public function propfindResultShouldContainEntries(
        string $shouldOrNot,
        TableNode $expectedFiles
    ): void {
        $this->verifyTableNodeColumnsCount($expectedFiles, 1);
        $elementRows = $expectedFiles->getRows();
        $should = ($shouldOrNot !== "not");

        foreach ($elementRows as $expectedFile) {
            $fileFound = $this->findEntryFromPropfindResponse(
                $expectedFile[0]
            );
            if ($should) {
                Assert::assertNotEmpty(
                    $fileFound,
                    "response does not contain the entry '$expectedFile[0]'"
                );
            } else {
                Assert::assertFalse(
                    $fileFound,
                    "response does contain the entry '$expectedFile[0]' but should not"
                );
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
    public function verifyTableNodeColumnsCount(
        TableNode $table,
        int $count
    ): void {
        if (\count($table->getRows()) < 1) {
            throw new Exception("Table should have at least one row.");
        }
        $rowCount = \count($table->getRows()[0]);
        if ($count !== $rowCount) {
            throw new Exception("Table expected to have $count rows but found $rowCount");
        }
    }

    /**
     * parses a PROPFIND response from $this->response into xml
     * and returns found search results if found else returns false
     *
     * @param string|null $entryNameToSearch
     * @return string|array|boolean
     * string if $entryNameToSearch is given and is found
     * array if $entryNameToSearch is not given
     * boolean false if $entryNameToSearch is given and is not found
     */
    public function findEntryFromPropfindResponse(
        string $entryNameToSearch = null
    ) {
        $spaceId = $this->getResponseSpaceId();
        //if we are using that step the second time in a scenario e.g. 'But ... should not'
        //then don't parse the result again, because the result in a ResponseInterface
        if (empty($this->getResponseXml())) {
            $this->setResponseXml(
                HttpRequestHelper::parseResponseAsXml($this->featureContext->getResponse())
            );
        }
        Assert::assertNotEmpty($this->getResponseXml(), __METHOD__ . ' Response is empty');
        Assert::assertNotEmpty($spaceId, __METHOD__ . ' SpaceId is empty');

        // trim any leading "/" passed by the caller, we can just match the "raw" name
        $trimmedEntryNameToSearch = \trim($entryNameToSearch, "/");

        // topWebDavPath should be something like /remote.php/webdav/ or
        // /remote.php/dav/files/alice/
        $topWebDavPath = "/" . "dav/spaces/" . $spaceId . "/";

        Assert::assertIsArray(
            $this->responseXml,
            __METHOD__ . " responseXml for space $spaceId is not an array"
        );
        Assert::assertArrayHasKey(
            "value",
            $this->responseXml,
            __METHOD__ . " responseXml for space $spaceId does not have key 'value'"
        );
        $multistatusResults = $this->responseXml["value"];
        $results = [];
        if ($multistatusResults !== null) {
            foreach ($multistatusResults as $multistatusResult) {
                $entryPath = $multistatusResult['value'][0]['value'];
                $entryName = \str_replace($topWebDavPath, "", $entryPath);
                $entryName = \rawurldecode($entryName);
                $entryName = \trim($entryName, "/");
                if ($trimmedEntryNameToSearch === $entryName) {
                    return $multistatusResult;
                }
                $results[] = $entryName;
            }
        }
        if ($entryNameToSearch === null) {
            return $results;
        }
        return false;
    }

    /**
     * @When /^user "([^"]*)" creates a folder "([^"]*)" in space "([^"]*)" using the WebDav Api$/
     *
     * @param string $user
     * @param string $folder
     * @param string $spaceName
     *
     * @return void
     */
    public function theUserCreatesAFolderUsingTheGraphApi(
        string $user,
        string $folder,
        string $spaceName
    ): void
    {
        $this->featureContext->setResponse(
            $this->sendCreateFolderRequest(
                $this->featureContext->getBaseUrl(),
                "MKCOL",
                $user,
                $this->featureContext->getPasswordForUser($user),
                $folder,
                $spaceName
            )
        );
    }

    /**
     * Send Graph Create Space Request
     *
     * @param string $baseUrl
     * @param string $method
     * @param string $user
     * @param string $password
     * @param string $folder
     * @param string $spaceName
     * @param string $xRequestId
     * @param array $headers
     * @return ResponseInterface
     */
    public function sendCreateFolderRequest(
        string $baseUrl,
        string $method,
        string $user,
        string $password,
        string $folder,
        string $spaceName,
        string $xRequestId = '',
        array $headers = []
    ): ResponseInterface
    {
        $spaceId = $this->getAvailableSpaces()[$spaceName]["id"];
        $fullUrl = $baseUrl;
        if (!str_ends_with($fullUrl, '/')) {
            $fullUrl .= '/';
        }
        $fullUrl .= "dav/spaces/" .  $spaceId . '/' . $folder;

        return HttpRequestHelper::sendRequest($fullUrl, $xRequestId, $method, $user, $password, $headers);
    }
}
