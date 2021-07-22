<?php
/**
 * @author Benedikt Kulmann <bkulmann@owncloud.com>
 *
 * @copyright Copyright (c) 2020, ownCloud GmbH
 * @license AGPL-3.0
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3,
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License, version 3,
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 *
 */

namespace OCA\Web\Controller;

use GuzzleHttp\Mimetypes;
use OC\AppFramework\Http;
use OCP\App\IAppManager;
use OCP\AppFramework\Controller;
use OCP\AppFramework\Http\ContentSecurityPolicy;
use OCP\AppFramework\Http\DataDisplayResponse;
use OCP\AppFramework\Http\DataResponse;
use OCP\AppFramework\Http\Response;
use OCP\IConfig;
use OCP\IRequest;

/**
 * Class FilesController
 *
 * @package OCA\Web\Controller
 */
class FilesController extends Controller {

	/**
	 * @var IConfig
	 */
	private $config;
    /**
     * @var IAppManager
     */
	private $appManager;

    /**
     * FilesController constructor.
     *
     * @param string $appName
     * @param IRequest $request
     * @param IConfig $config
     * @param IAppManager $appManager
     */
	public function __construct(string $appName, IRequest $request, IConfig $config, IAppManager  $appManager) {
		parent::__construct($appName, $request);
		$this->config = $config;
		$this->appManager = $appManager;
	}

	/**
	 * Tries to load a file by the given $path.
	 *
	 * @PublicPage
	 * @NoCSRFRequired
	 *
	 * @param $path string
	 * @return DataResponse
	 */
	public function getFile(string $path): Response {
		// don't allow directory traversal to parents
		if (\strpos($path, "..") !== false) {
			return new DataResponse(['error' => 'resource not found'], Http::STATUS_NOT_FOUND);
		}

		// check if path permitted
		$permittedPaths = ["css", "img", "js", "themes", "index.html", "manifest.json", "oidc-callback.html", "oidc-silent-redirect.html"];
		$found = false;
		foreach ($permittedPaths as $p) {
			if (\strpos($path, $p) === 0) {
				$found = true;
				break;
			}
		}
		if (!$found) {
			return new DataResponse(['error' => 'resource not found'], Http::STATUS_NOT_FOUND);
		}

		// check if path resolves to an actual file
		if (\is_dir($path)) {
			return new DataResponse(['error' => 'resource not found'], Http::STATUS_NOT_FOUND);
		}
		$basePath = \dirname(__DIR__,2);
		$absolutePath = \realpath( $basePath . '/' . $path);
		if ($absolutePath === false) {
			return new DataResponse(['error' => 'resource not found'], Http::STATUS_NOT_FOUND);
		}
		if (\strpos($absolutePath, $basePath) !== 0) {
			return new DataResponse(['error' => 'resource not found'], Http::STATUS_NOT_FOUND);
		}

		$response = new DataDisplayResponse(\file_get_contents($absolutePath), Http::STATUS_OK, [
			'Content-Type' => $this->getMimeType($absolutePath),
			'Content-Length' => \filesize($absolutePath),
			'Cache-Control' => 'max-age=0, no-cache, no-store, must-revalidate',
			'Pragma' => 'no-cache',
			'Expires' => 'Wed, 11 Jan 1984 05:00:00 GMT',
			'X-Frame-Options' => 'DENY'
		]);
		if (\strpos($path, "oidc-callback.html") === 0 || \strpos($path, "oidc-silent-redirect.html") === 0) {
			$csp = new ContentSecurityPolicy();
			$csp->allowInlineScript(true);
            $csp = $this->applyCSPOpenIDConnect($csp);
			$response->setContentSecurityPolicy($csp);
		}
		if (\strpos($path, "index.html") === 0) {
            $csp = new ContentSecurityPolicy();
            $csp->allowInlineScript(true);
            $csp = $this->applyCSPOpenIDConnect($csp);

            // for now we set CSP rules manually, until we have sufficient requirements for a generic solution.
            $csp = $this->applyCSPOnlyOffice($csp);
            $csp = $this->applyCSPRichDocuments($csp);

            $response->setContentSecurityPolicy($csp);
        }

		return $response;
	}

	private function getMimeType(string $filename): string {
		$mimeTypes = Mimetypes::getInstance();
		return $mimeTypes->fromFilename($filename);
	}

    private function applyCSPOpenIDConnect(ContentSecurityPolicy $csp): ContentSecurityPolicy {
        $oidcUrl = $this->getOpenIDConnectServerUrl();
        $oidcServerUrl = $this->extractDomain($oidcUrl);
        if (!empty($oidcServerUrl)) {
            $csp->addAllowedConnectDomain($oidcServerUrl);
        }
        return $csp;
    }

    /**
     * Extracts the openidconnect provider-url from the app
     *
     * @return string
     * @throws \OCP\AppFramework\QueryException
     */
    private function getOpenIDConnectServerUrl(): string {
        if (!$this->isAppEnabled("openidconnect")) {
            return "";
        }
        if (!class_exists("\OCA\OpenIdConnect\Application")) {
            return "";
        }
        $oidcClient = \OC::$server->query(\OCA\OpenIdConnect\Client::class);
        return $oidcClient->getProviderURL();
    }

	private function applyCSPOnlyOffice(ContentSecurityPolicy $csp): ContentSecurityPolicy {
        $ooUrl = $this->getOnlyOfficeDocumentServerUrl();
        $documentServerUrl = $this->extractDomain($ooUrl);
        if (!empty($documentServerUrl)) {
            $csp->addAllowedScriptDomain($documentServerUrl);
            $csp->addAllowedFrameDomain($documentServerUrl);
        } else if (!empty($ooUrl)) {
            $csp->addAllowedFrameDomain("'self'");
	}
        return $csp;
    }

    /**
     * Extracts the onlyoffice document server URL from the app
     *
     * @return string
     * @throws \OCP\AppFramework\QueryException
     */
    private function getOnlyOfficeDocumentServerUrl(): string {
        if (!$this->isAppEnabled("onlyoffice")) {
            return "";
        }
        if (!class_exists("\OCA\Onlyoffice\AppConfig")) {
            return "";
        }
        $onlyofficeConfig = \OC::$server->query(\OCA\Onlyoffice\AppConfig::class);
        return $onlyofficeConfig->GetDocumentServerUrl();
    }

    private function applyCSPRichDocuments(ContentSecurityPolicy $csp): ContentSecurityPolicy {
        $documentServerUrl = $this->extractDomain($this->getRichDocumentsServerUrl());
        if (!empty($documentServerUrl)) {
            $csp->addAllowedFrameDomain($documentServerUrl);
        }
        return $csp;
    }

    /**
     * Extracts the richdocuments document server URL from the app-config, in the same manner like
     * the richdocuments app:
     * - https://github.com/owncloud/richdocuments/blob/9a23f426048c540793fc16119f71a44c26077f16/lib/Controller/DocumentController.php#L122
     * - https://github.com/owncloud/richdocuments/blob/9a23f426048c540793fc16119f71a44c26077f16/lib/Controller/DocumentController.php#L393
     *
     * @return string
     * @throws \OCP\AppFramework\QueryException
     */
    private function getRichDocumentsServerUrl(): string {
        if (!$this->isAppEnabled("richdocuments")) {
            return "";
        }
        if (!class_exists("\OCA\Richdocuments\AppConfig")) {
            return "";
        }
        $richdocumentsConfig = \OC::$server->query(\OCA\Richdocuments\AppConfig::class);
        if (empty($richdocumentsConfig)) {
            return "";
        }
        return $richdocumentsConfig->getAppValue('wopi_url');
    }

    /**
     * Extracts the domain part from a url.
     *
     * @param string $url
     * @return string
     */
    private function extractDomain(string $url): string {
        $parsedUrl = \parse_url($url);
        if (empty($parsedUrl['host'])) {
            return "";
        }
        return (isset($parsedUrl['scheme']) ? $parsedUrl['scheme'] . '://' : '')
            . $parsedUrl['host']
            . (isset($parsedUrl['port']) ? ':' . $parsedUrl['port'] : '');
    }

    /**
     * Checks whether the given app is installed and enabled.
     *
     * @param string $appName
     * @return bool
     */
    private function isAppEnabled(string $appName): bool {
        if (!$this->appManager->isInstalled($appName)) {
            return false;
        }
        return $this->appManager->isEnabledForUser($appName);
    }
}
