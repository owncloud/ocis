<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Kiran Parajuli <kiran@jankaritech.com>
 * @copyright Copyright (c) 2022 Kiran Parajuli kiran@jankaritech.com
 */

namespace TestHelpers;

use Exception;
use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\RequestInterface;
use Psr\Http\Message\ResponseInterface;

/**
 * A helper class for managing users and groups using the Graph API
 */
class GraphHelper {
	/**
	 * @param string $baseUrl
	 * @param string $path
	 *
	 * @return string
	 */
	private static function getFullUrl(string $baseUrl, string $path):string {
		$fullUrl = $baseUrl;
		if (\substr($fullUrl, -1) !== '/') {
			$fullUrl .= '/';
		}
		$fullUrl .= 'graph/v1.0/' . $path;
		return $fullUrl;
	}

    /**
     * @param string $baseUrl
     * @param string $xRequestId
     * @param string $method
     * @param string $path
     * @param string|null $body
     * @param array|null $headers
     *
     * @return RequestInterface
     */
	public static function createRequest(
		string $baseUrl,
		string $xRequestId,
		string $method,
		string $path,
		?string $body = null,
		?array $headers = []
	): RequestInterface {
		$fullUrl = self::getFullUrl($baseUrl, $path);
		return HttpRequestHelper::createRequest(
			$fullUrl,
			$xRequestId,
			$method,
			$headers,
			$body
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $adminUser
	 * @param string $adminPassword
	 * @param string $userName
	 * @param string $password
	 * @param string|null $email
	 * @param string|null $displayName
	 *
	 * @return ResponseInterface
	 */
	public static function createUser(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $userName,
		string $password,
		?string $email = null,
		?string $displayName = null
	):ResponseInterface {
		$payload = self::prepareCreateUserPayload(
			$userName,
			$password,
			$email,
			$displayName
		);

		$headers = ['Content-Type' => 'application/json'];
		$url = self::getFullUrl($baseUrl, 'users');
		return HttpRequestHelper::post(
			$url,
			$xRequestId,
			$adminUser,
			$adminPassword,
			$headers,
			$payload
		);
	}

    /**
     * @param string $baseUrl
     * @param string $xRequestId
     * @param string $adminUser
     * @param string $adminPassword
     * @param string $userId
     * @param string|null $userName
     * @param string|null $password
     * @param string|null $email
     * @param string|null $displayName
     *
     * @return ResponseInterface
     */
    public static function editUser(
        string $baseUrl,
        string $xRequestId,
        string $adminUser,
        string $adminPassword,
        string $userId,
        ?string $userName = null,
        ?string $password = null,
        ?string $email = null,
        ?string $displayName = null

    ): ResponseInterface {
        $payload = self::preparePatchUserPayload(
            $userName,
            $password,
            $email,
            $displayName
        );
        $headers = ['Content-Type' => 'application/json'];
        $url = self::getFullUrl($baseUrl, 'users/' . $userId);
        return HttpRequestHelper::sendRequest(
            $url,
            $xRequestId,
            "PATCH",
            $adminUser,
            $adminPassword,
            $headers,
            $payload
        );
    }

    /**
     * @param string $baseUrl
     * @param string $xRequestId
     * @param string $adminUser
     * @param string $adminPassword
     * @param string $userName
     *
     * @return ResponseInterface
     * @throws GuzzleException
     */
	public static function getUser(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $userName
	):ResponseInterface {
		$url = self::getFullUrl($baseUrl, 'users/' . $userName);
		return HttpRequestHelper::get(
			$url,
			$xRequestId,
			$adminUser,
			$adminPassword,
			["Content-Type" => "application/json"]
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $adminUser
	 * @param string $adminPassword
	 * @param string $userName
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function deleteUser(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $userName
	):ResponseInterface {
		$url = self::getFullUrl($baseUrl, 'users/' . $userName);
		return HttpRequestHelper::delete(
			$url,
			$xRequestId,
			$adminUser,
			$adminPassword,
		);
	}

	/**
	 * can send a request to the graph api to:
	 * - create a group
	 * - update a group
	 *
	 * displayName is the only field that can be assigned/updated
	 *
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $adminUser
	 * @param string $adminPassword
	 * @param string $groupName - the displayName of the group
	 * @param bool|null $update
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	private static function postPatchGroup(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $groupName,
		?bool $update = false
	): ResponseInterface {
		$url = ($update)
			? self::getFullUrl($baseUrl, 'groups/' . $groupName)
			: self::getFullUrl($baseUrl, 'groups');
		$method = ($update) ? 'PATCH' : 'POST';
		$headers = ['Content-Type' => 'application/json'];
		$payload['displayName'] = $groupName;
		return HttpRequestHelper::sendRequest(
			$url,
			$xRequestId,
			$method,
			$adminUser,
			$adminPassword,
			$headers,
			\json_encode($payload)
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $adminUser
	 * @param string $adminPassword
	 * @param string $groupName
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function createGroup(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $groupName
	):ResponseInterface {
		return self::postPatchGroup(
			$baseUrl,
			$xRequestId,
			$adminUser,
			$adminPassword,
			$groupName
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $adminUser
	 * @param string $adminPassword
	 * @param string $groupId
	 * @param string $displayName
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function updateGroup(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $groupId,
		string $displayName
	):ResponseInterface {
		return self::postPatchGroup(
			$baseUrl,
			$xRequestId,
			$adminUser,
			$adminPassword,
			$displayName,
			true
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $adminUser
	 * @param string $adminPassword
	 *
	 * @return array
	 * @throws Exception
	 */
	public static function getGroups(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword
	):array {
		$url = self::getFullUrl($baseUrl, 'groups');
		$response = HttpRequestHelper::get(
			$url,
			$xRequestId,
			$adminUser,
			$adminPassword
		);
		$groupsListEncoded = \json_decode($response->getBody()->getContents(), true);
		if (!isset($groupsListEncoded['value'])) {
			throw new Exception('No groups found');
		} else {
			return $groupsListEncoded['value'];
		}
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $adminUser
	 * @param string $adminPassword
	 * @param string $groupId
	 *
	 * @return ResponseInterface
	 * @throws GuzzleException
	 */
	public static function deleteGroup(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $groupId
	):ResponseInterface {
		$url = self::getFullUrl($baseUrl, 'groups/' . $groupId);
		return HttpRequestHelper::delete(
			$url,
			$xRequestId,
			$adminUser,
			$adminPassword,
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $xRequestId
	 * @param string $adminUser
	 * @param string $adminPassword
	 * @param string $groupId
	 * @param array $users expects users array with user ids [ [ 'id' => 'some_id' ], ]
	 *
	 * @return ResponseInterface
	 */
	public static function addUsersToGroup(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $groupId,
		array $users
	):ResponseInterface {
		$url = self::getFullUrl($baseUrl, 'groups/' . $groupId . '/users');
		$payload = [
			"members@odata.bind" => []
		];
		foreach ($users as $user) {
			$payload[0][] = self::getFullUrl($baseUrl, 'users/' . $user["id"]);
		}
		return HttpRequestHelper::post(
			$url,
			$xRequestId,
			$adminUser,
			$adminPassword,
			['Content-Type' => 'application/json'],
			\json_encode($payload)
		);
	}

	public static function addUserToGroup(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $userId,
		string $groupId
	):ResponseInterface {
		$url = self::getFullUrl($baseUrl, 'groups/' . $groupId . '/members/$ref');
		$body = [
			"@odata.id" => self::getFullUrl($baseUrl, 'users/' . $userId)
		];
		return HttpRequestHelper::post(
			$url,
			$xRequestId,
			$adminUser,
			$adminPassword,
			["application/json"],
			\json_encode($body)
		);
	}

	public static function removeUserFromGroup(
		string $baseUrl,
		string $xRequestId,
		string $adminUser,
		string $adminPassword,
		string $userId,
		string $groupId
	): ResponseInterface {
		$url = self::getFullUrl($baseUrl, 'groups/' . $groupId . '/members/' . $userId . '/$ref');
		return HttpRequestHelper::delete(
			$url,
			$xRequestId,
			$adminUser,
			$adminPassword,
		);
	}

    public static function getMembersList(
        string $baseUrl,
        string $xRequestId,
        string $adminUser,
        string $adminPassword,
        string $userId,
        string $groupId
    ): bool {
        $url = self::getFullUrl($baseUrl, 'groups/' . $groupId . '/members/' . $userId . '/$ref');
        return HttpRequestHelper::get(
            $url,
            $xRequestId,
            $adminUser,
            $adminPassword
        );
    }

    /**
     * @param string $baseUrl
     * @param string $xRequestId
     * @param string $adminUser
     * @param string $adminPassword
     * @param string $userId
     *
     * @return void
     */
    public static function getGroupListOfAUser(
        string $baseUrl,
        string $xRequestId,
        string $adminUser,
        string $adminPassword,
        string $userId
    ) {
        // TODO: endpoint not available https://github.com/owncloud/ocis/issues/3363
        // Not implemented yet
    }

	/**
	 * @param string|null $userName
	 * @param string|null $password
	 * @param string|null $email
	 * @param string|null $displayName
	 *
	 * @return string
	 */
	public static function prepareCreateUserPayload(
        string $userName,
        string $password,
        ?string $email,
        ?string $displayName
    ): string {
		$payload['onPremisesSamAccountName'] = $userName;
		$payload['passwordProfile'] = ['password' => $password];
		$payload['displayName'] = $displayName ?? $userName;
		$payload['mail'] = $email ?? $userName . '@example.com';
		return \json_encode($payload);
	}
    public static function preparePatchUserPayload(
        ?string $userName,
        ?string $password,
        ?string $email,
        ?string $displayName
    ): string {
        $payload = [];
        if ($userName) $payload['onPremisesSamAccountName'] = $userName;
        if ($password) $payload['passwordProfile'] = ['password' => $password];
        if ($displayName) $payload['displayName'] = $displayName;
        if ($email) $payload['mail'] = $email;
        return \json_encode($payload);
    }
}
