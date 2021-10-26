<?php

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\HttpRequestHelper;
use TestHelpers\SetupHelper;
use PHPUnit\Framework\Assert;
use \Behat\Gherkin\Node\TableNode;

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
    private $responseSpaceId;

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
     * Send Graph Create Space Request
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
        $this->rememberTheAvailableSpaces();
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
     * @When /^user "([^"]*)" creates a space "([^"]*)" of type "([^"]*)" with quota "([^"]*)" using the GraphApi$/
     *
     * @param $user string
     * @param $spaceName string
     * @param $spaceType string
     * @param $quota int
     *
     * @return void
     * @throws JsonException
     */
    public function theUserCreatesASpaceWithQuotaUsingTheGraphApi($user, $spaceName, $spaceType, $quota): void
    {
        $space = ["Name" => $spaceName, "driveType" => $spaceType, "quota" => ["total" => (int) $quota]];
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
    public function theAdministratorGivesUserTheRole(string $user, string $role): void
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
        if ($this->featureContext->getResponse()) {
            $rawBody =  $this->featureContext->getResponse()->getBody()->getContents();
        }
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
        if ($this->featureContext->getResponse()) {
            $rawBody = $this->featureContext->getResponse()->getBody()->getContents();
        }
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
        if ($this->featureContext->getResponse()) {
            $rawBody = $this->featureContext->getResponse()->getBody()->getContents();
        }

        $assignment = [];
        if (isset(\json_decode($rawBody,  true)["assignment"])) {
            $assignment = \json_decode($rawBody, true)["assignment"];
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
     * @param $user
     * @param $name
     * @return void
     */
    public function theUserListsTheContentOfAPersonalSpaceRootUsingTheWebDAvApi($user, $name): void
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
                [],
                []
            )
        );
        $this->setResponseSpaceId($spaceId);
        $this->setResponseXml(HttpRequestHelper::parseResponseAsXml($this->featureContext->getResponse())
        );
    }

    /**
     * @Then /^the (?:propfind|search) result of the space should (not|)\s?contain these (?:files|entries):$/
     *
     * @param string $user
     * @param string $shouldOrNot (not|)
     * @param TableNode $expectedFiles
     *
     * @return void
     * @throws Exception
     */
    public function thePropfindResultShouldContainEntries(
        $shouldOrNot,
        TableNode $expectedFiles
    ) {
        $this->propfindResultShouldContainEntries(
            $shouldOrNot,
            $expectedFiles,
        );
    }

    /**
     * @Then /the json responded should contain these key and value pairs/
     *
     * @param TableNode $table
     *
     * @return void
     */
    public function jsonRespondedShouldContain(TableNode $table) {
        $this->featureContext->verifyTableNodeColumns($table, ['key', 'value']);
        $responseJson = json_decode($this->featureContext->getResponse()->getBody(), true);

        foreach ($table->getHash() as $row) {
            if (empty($this->searchKeyValueInArray($responseJson, $row["key"], $row["value"]))){
                Assert::assertFalse($row["value"], ($row["value"] . ' not found'));
            }
        }
    }
    
    /**
     * Method search for a match $key->$value
     * 
     * @param array $array
     * @param string $key
     * @param string $value
     * @return array $results
     */
    public function searchKeyValueInArray($array, $key, $value)
    {
        $results = array();

        if (is_array($array)) {
            if (isset($array[$key]) && $array[$key] == $value) {
                $results[] = $array;
            }

            foreach ($array as $subarray) {
                $results = array_merge($results, $this->searchKeyValueInArray($subarray, $key, $value));
            }
        }
        return $results;
    }

    /**
     * @param string $shouldOrNot (not|)
     * @param TableNode $expectedFiles
     * @param string|null $user
     *
     * @return void
     * @throws Exception
     */
    public function propfindResultShouldContainEntries(
        $shouldOrNot,
        TableNode $expectedFiles
    ) {
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
    public function verifyTableNodeColumnsCount($table, $count) {
        if (!($table instanceof TableNode)) {
            throw new Exception("TableNode expected but got " . \gettype($table));
        }
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
     * @throws JsonException
     */
    public function theUserCreatesAFolderUsingTheGraphApi($user, $folder, $spaceName): void
    {
        $this->featureContext->setResponse(
            $this->sendCreateFolderRequest(
                $this->featureContext->getBaseUrl(),
                "",
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
     * @param $baseUrl
     * @param $user
     * @param $password
     * @param string $method
     * @param string $xRequestId
     * @param array $headers
     * @param string $folder
     * @param string $spaceName
     * @return ResponseInterface
     */
    public function sendCreateFolderRequest(
        $baseUrl,
        string $xRequestId = '',
        string $method,
        $user,
        $password,
        $folder,
        $spaceName,
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
