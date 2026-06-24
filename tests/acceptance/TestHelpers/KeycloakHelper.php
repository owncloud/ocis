<?php declare(strict_types=1);

namespace TestHelpers;

use Exception;
use GuzzleHttp\Exception\GuzzleException;
use GuzzleHttp\Psr7\Query;
use Psr\Http\Message\ResponseInterface;

/**
 * A helper class for Keycloak admin API requests.
 */
class KeycloakHelper {
    const OCIS_KEYCLOAK_USER_ROLES = [
        'Admin' => 'ocisAdmin',
        'Space Admin' => 'ocisSpaceAdmin',
        'User' => 'ocisUser',
        'User Light' => 'ocisGuest'
    ];
    private static ?string $adminAccessToken = null;
//    private static array $userTokens = [];

    /**
     * @return bool
     */
    public static function isTestingWithKeycloak(): bool {
        return (\getenv('KEYCLOAK') === "true");
    }

    public static function getKeycloakUrl(): string {
        $keycloakUrl = \getenv('KEYCLOAK_URL');
        if ($keycloakUrl !== false && $keycloakUrl !== '') {
            return $keycloakUrl;
        }
        return 'https://keycloak.owncloud.test';
    }

    public static function setAdminAccessToken(string $accessToken): void {
        self::$adminAccessToken = $accessToken;
    }

    /**
     * @throws GuzzleException
     */
    public static function getAdminAccessToken(): string {
        // Refresh when token is missing or about to expire.
        if (self::$adminAccessToken === null) {
            $tokenData = self::generateAdminAccessToken();
            self::setAdminAccessToken($tokenData['accessToken']);
        }

        return (string)self::$adminAccessToken;
    }

    /**
     * @throws Exception
     */
//    public static function setOcisUserToken(array $user, array $tokenData): void {
//        $userId = $user['userid'] ?? null;
//        if ($userId === null) {
//            throw new Exception('User ID is required to store token');
//        }
//
//        self::$userTokens[$userId] = [
//            'user' => $user,
//            'token' => [
//                'userid' => $userId,
//                'accessToken' => $tokenData['access_token'],
//                'refreshToken' => $tokenData['refresh_token']
//            ]
//        ];
//    }
//
//    public static function getOcisUserToken(string $userId): array {
//        return self::$userTokens[$userId];
//    }

    /**
     * @param string $username
     * @param string $password
     * @param string|null $email
     * @param string|null $displayName
     *
     * @return ResponseInterface
     * @throws Exception
     * @throws GuzzleException
     */
    public static function createUser(
        string $username,
        string $password,
        ?string $email = null,
        ?string $displayName = null,
    ): ResponseInterface {
        $accessToken = self::getAdminAccessToken();
        $url = self::getKeycloakUrl() . '/admin/realms/oCIS/users';

        $response = HttpRequestHelper::post(
            $url,
            null,
            null,
            [
                'Authorization' => 'Bearer ' . $accessToken,
                'Content-Type' => 'application/json',
            ],
            self::prepareCreateUserPayload($username, $password, $email, $displayName),
        );
        return $response;
    }

    public static function assignRole(
        string $uuid,
        string $role
    ): ResponseInterface {
        $url = self::getKeycloakUrl() . "/admin/realms/oCIS/users/" . $uuid . "/role-mappings/realm";
        $body = [
            [
                'id' => '8c79ff81-c256-48fd-b0b9-795c7941eedf',
                'name' => self::OCIS_KEYCLOAK_USER_ROLES[$role]
            ],
            [
                'id' => 'e2145b30-bf6f-49fb-af3f-1b40168bfcef',
                'name' => 'offline_access'
            ]
        ];
        return HttpRequestHelper::post(
            $url,
            null,
            null,
            [
                "Content-Type" => "application/json",
                "Authorization" => "Bearer " . self::getAdminAccessToken()
            ],
            json_encode($body, JSON_THROW_ON_ERROR)
        );
    }

    /**
     * @return string
     * @throws Exception
     * @throws GuzzleException
     */
    private static function generateAdminAccessToken(): array {
        $url = self::getKeycloakUrl()
            . '/realms/master/protocol/openid-connect/token';
        $response = HttpRequestHelper::post(
            $url,
            null,
            null,
            ['Content-Type' => 'application/x-www-form-urlencoded'],
            [
                'client_id' => 'admin-cli',
                'username' => 'admin',
                'password' => 'admin',
                'grant_type' => 'password',
            ],
        );

        if ($response->getStatusCode() >= 400) {
            throw new Exception(
                __METHOD__
                . ' failed to get Keycloak admin access token, status '
                . $response->getStatusCode()
                . ', response '
                . (string)$response->getBody(),
            );
        }

        $decodedResponse = HttpRequestHelper::getJsonDecodedResponseBodyContent($response);
        if (!isset($decodedResponse->access_token)) {
            throw new Exception(__METHOD__ . ' could not find access_token in Keycloak token response');
        }

        $expiresInSeconds = 60;
        if (isset($decodedResponse->expires_in)) {
            $expiresInSeconds = (int)$decodedResponse->expires_in;
        }

        return [
            'accessToken' => (string)$decodedResponse->access_token,
            'expiresInSeconds' => $expiresInSeconds,
        ];
    }

    /**
     * @param string $username
     * @param string $password
     * @param string|null $email
     * @param string|null $displayName
     *
     * @return string
     */
    private static function prepareCreateUserPayload(
        string $username,
        string $password,
        ?string $email = null,
        ?string $displayName = null,
    ): string {
        $firstName = $username;
        $lastName = '';
        if ($displayName !== null && \trim($displayName) !== '') {
            $nameParts = \preg_split('/\s+/', \trim($displayName), 2);
            if ($nameParts !== false && isset($nameParts[0])) {
                $firstName = $nameParts[0];
            }
            if ($nameParts !== false && isset($nameParts[1])) {
                $lastName = $nameParts[1];
            }
        }

        $payload = [
            'username' => $username,
            'credentials' => [[
                'value' => $password,
                'type' => 'password',
            ]],
            'firstName' => $firstName,
            'lastName' => $lastName,
            'emailVerified' => true,
            'enabled' => true,
        ];

        if ($email !== null) {
            $payload['email'] = $email;
        }

        return \json_encode($payload, JSON_THROW_ON_ERROR);
    }

    /**
     * @throws GuzzleException
     */
    public static function getAuthorizationEndPoint(): array {
        $loginParams = [
            'client_id' => 'web',
            'redirect_uri' => 'https://ocis.owncloud.test/oidc-callback.html',
            'response_mode' => 'query',
            'response_type'=> 'code',
            'scope' => 'openid profile email acr'
        ];
        $queryString = \http_build_query($loginParams);
        $authUrl = self::getKeycloakUrl() . "/realms/oCIS/protocol/openid-connect/auth?" . $queryString;
        $response = HttpRequestHelper::get(
            $authUrl
        );
        $cookie = $response->getHeader("Set-Cookie")[0];
        $htmlData = $response->getBody()->getContents();
        if (!preg_match('/action="([^"]+)"/i', $htmlData, $match)) {
            throw new Exception('No authorization url found in the HTML response body.');
        }
        $authorizationUrl = $match[1];
        return [$authorizationUrl, $cookie];
    }

    public static function getCode(array $user, string $authorizationUrl, string $cookie): string {
        $authCodeResponse = HttpRequestHelper::post(
            $authorizationUrl,
            null,
            null,
            [
                'Cookie' => $cookie
            ],
            [
                'username' => $user['userid'],
                'password' => $user['password']
            ]
        );
        $locationHeader = $authCodeResponse->getHeader('Location');
        $queryString = parse_url($locationHeader[0], PHP_URL_QUERY);
        parse_str($queryString, $urlParams);
        return $urlParams['code'];
    }

    /**
     * @throws GuzzleException
     */
    public static function getToken(string $authorizationCode): ResponseInterface {
        $tokenResponse = HttpRequestHelper::post(
            'https://keycloak.owncloud.test/realms/oCIS/protocol/openid-connect/token',
            null,
            null,
            null,
            [
                'client_id' => 'web',
                'code' => $authorizationCode,
                'redirect_uri' => 'https://ocis.owncloud.test/oidc-callback.html',
                'grant_type' => 'authorization_code'
            ]
        );
        if ($tokenResponse->getStatusCode() !== 200) {
            throw new Exception('Failed to retrieve token: Expected status code to be 200 but received' . $tokenResponse->getStatusCode() . '.\nMessage: ' . $tokenResponse->getBody()->getContents());
        }
        return $tokenResponse;
    }

    /**
     * @throws GuzzleException
     * @throws Exception
     */
    public static function setAccessTokenForKeycloakOcisUser(array $user): array {
        [$authorizationUrl, $cookie] = self::getAuthorizationEndPoint();
        $authorizationCode = self::getCode($user, $authorizationUrl, $cookie);
        $tokenResponse = self::getToken($authorizationCode);
        $tokenData = json_decode($tokenResponse->getBody()->getContents(), true, 512, JSON_THROW_ON_ERROR);
        return $tokenData;
    }
}
