<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Sajan Gurung <sajan@jankaritech.com>
 * @copyright Copyright (c) 2023 Sajan Gurung sajan@jankaritech.com
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

use Behat\Behat\Context\Context;
use Behat\Gherkin\Node\TableNode;
use GuzzleHttp\Exception\GuzzleException;
use Psr\Http\Message\ResponseInterface;
use TestHelpers\HttpRequestHelper;
use TestHelpers\KeycloakHelper;
use TestHelpers\OcisConfigHelper;
use TestHelpers\GraphHelper;
use PHPUnit\Framework\Assert;

/**
 * steps needed to re-configure oCIS server
 */
class OcisConfigContext implements Context {
	private array $enabledPermissionsRoles = [];

	/**
	 * @return array
	 */
	public function getEnabledPermissionsRoles(): array {
		return $this->enabledPermissionsRoles;
	}

	/**
	 * @param array $enabledPermissionsRoles
	 *
	 * @return void
	 */
	public function setEnabledPermissionsRoles(array $enabledPermissionsRoles): void {
		$this->enabledPermissionsRoles = $enabledPermissionsRoles;
	}

	/**
	 * @Given async upload has been enabled with post-processing delayed to :delayTime seconds
	 *
	 * @param string $delayTime
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function asyncUploadHasBeenEnabledWithDelayedPostProcessing(string $delayTime): void {
		if (OcisConfigHelper::isK8s()) {
			$envs = [
				"storageusers" => ["OCIS_ASYNC_UPLOADS" => true, "OCIS_EVENTS_ENABLE_TLS" => false],
				"postprocessing" => ["POSTPROCESSING_DELAY" => $delayTime . "s"],
			];
			$response = OcisConfigHelper::reConfigureOcisK8s($envs);
		} else {
			$envs = [
				"OCIS_ASYNC_UPLOADS" => true,
				"OCIS_EVENTS_ENABLE_TLS" => false,
				"POSTPROCESSING_DELAY" => $delayTime . "s",
			];
			$response = OcisConfigHelper::reConfigureOcis($envs);
		}

		$resBody = $response->getBody()->getContents();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set async upload with delayed post processing. Response: $resBody",
		);
	}

	/**
	 * @Given the config :configVariable has been set to :configValue
	 * @Given the config :configVariable has been set to :configValue for :serviceName service
	 *
	 * @param string $configVariable
	 * @param string $configValue
	 * @param string|null $serviceName
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theConfigHasBeenSetTo(
		string $configVariable,
		string $configValue,
		?string $serviceName = null,
	): void {
		if (getenv("K8S") === "true") {
			$envs = [
				$serviceName => [$configVariable => $configValue],
			];
			$response = OcisConfigHelper::reConfigureOcisK8s($envs);

		} else {
			$envs = [
				$configVariable => $configValue,
			];
			$response = OcisConfigHelper::reConfigureOcis($envs);
		}

		$this->assertOcisRestarted(
			$response,
			"Failed to set config $configVariable=$configValue.",
		);
	}

	/**
	 * @Given the administrator has enabled the permissions role :role
	 *
	 * @param string $role
	 *
	 * @return void
	 */
	public function theAdministratorHasEnabledTheRole(string $role): void {
		$roleId = GraphHelper::getPermissionsRoleIdByName($role);
		$defaultRoles = array_values(GraphHelper::DEFAULT_PERMISSIONS_ROLES);
		if (!\in_array($roleId, $defaultRoles)) {
			$defaultRoles[] = $roleId;
		}

		if (getenv("K8S") === "true") {
			$envs = [
				'graph' => [ "GRAPH_AVAILABLE_ROLES" => implode(',', $defaultRoles)],
			];
			$response = OcisConfigHelper::reConfigureOcisK8s($envs);

		} else {
			$envs = [
				"GRAPH_AVAILABLE_ROLES" => implode(',', $defaultRoles),
			];
			$response = OcisConfigHelper::reConfigureOcis($envs);
		}

		$resBody = $response->getBody()->getContents();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to enable role $role. Response: $resBody",
		);
		$this->setEnabledPermissionsRoles($defaultRoles);
	}

	/**
	 * @Given the administrator has enabled the following share permissions roles:
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 */
	public function theAdministratorHasEnabledTheFollowingSharePermissionsRoles(TableNode $table): void {
		$defaultRoles = array_values(GraphHelper::DEFAULT_PERMISSIONS_ROLES);
		$roles = [];
		foreach ($table->getHash() as $row) {
			$roles[] = $row['permissions-role'];
			$roleId = GraphHelper::getPermissionsRoleIdByName($row['permissions-role']);
			if (!\in_array($roleId, $defaultRoles)) {
				$defaultRoles[] = $roleId;
			}
		}

		if (getenv("K8S") === "true") {
			$envs = [
				'graph' => ["GRAPH_AVAILABLE_ROLES" => implode(',', $defaultRoles)],
			];
			$response = OcisConfigHelper::reConfigureOcisK8s($envs);

		} else {
			$envs = [
				"GRAPH_AVAILABLE_ROLES" => implode(',', $defaultRoles),
			];
			$response = OcisConfigHelper::reConfigureOcis($envs);
		}

		$resBody = $response->getBody()->getContents();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to enable roles: " . implode(', ', $roles) . ". Response: $resBody",
		);
		$this->setEnabledPermissionsRoles($defaultRoles);
	}

	/**
	 * @Given the administrator has disabled the permissions role :role
	 *
	 * @param string $role
	 *
	 * @return void
	 */
	public function theAdministratorHasDisabledThePermissionsRole(string $role): void {
		$roleId = GraphHelper::getPermissionsRoleIdByName($role);
		$availableRoles = $this->getEnabledPermissionsRoles();

		if ($key = array_search($roleId, $availableRoles)) {
			unset($availableRoles[$key]);
		}
		if (getenv("K8S") === "true") {
			$envs = [
				'graph' => [ "GRAPH_AVAILABLE_ROLES" => implode(',', $availableRoles)],
			];
			$response = OcisConfigHelper::reConfigureOcisK8s($envs);

		} else {
			$envs = [
				"GRAPH_AVAILABLE_ROLES" => implode(',', $availableRoles),
			];
			$response = OcisConfigHelper::reConfigureOcis($envs);
		}
		$resBody = $response->getBody()->getContents();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to disable role $role. Response: $resBody",
		);
		$this->setEnabledPermissionsRoles($availableRoles);
	}

	/**
	 * @Given the config :configVariable has been set to path :path for :serviceName service
	 *
	 * @param string $configVariable
	 * @param string $path
	 * @param string|null $serviceName
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theConfigHasBeenSetToPath(string $configVariable, string $path, ?string $serviceName = null): void {
		if (getenv("K8S") === "true") {
			if (\str_ends_with($configVariable, "PASSWORD_POLICY_BANNED_PASSWORDS_LIST")) {
				// The banned password list is already configured in K8s setup.
				return;
			}
			// In K8s, we MUST use the mounted paths
			// All files for test MUST be mounted in '/etc/ocis/' + $path
			$k8sPath = "/etc/ocis/" . $path;
			$envs = [
				$serviceName => [$configVariable => $k8sPath],
			];
			$response = OcisConfigHelper::reConfigureOcisK8s($envs);

		} else {
			$path = \dirname(__FILE__) . "/../../" . $path;
			$envs = [
				$configVariable => $path,
			];
			$response = OcisConfigHelper::reConfigureOcis($envs);
		}
		$resBody = $response->getBody()->getContents();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set config $configVariable=$path. Response: $resBody",
		);
	}

	/**
	 * @Given the following configs have been set:
	 *
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theConfigHasBeenSetToValue(TableNode $table): void {
		$envs = [];
		if (getenv("K8S") === "true") {
			foreach ($table->getHash() as $row) {
				$envs[$row['service']][$row['config']] = $row['value'];
			}
			$response = OcisConfigHelper::reConfigureOcisK8s($envs);
		} else {
			foreach ($table->getHash() as $row) {
				$envs[$row['config']] = $row['value'];
			}
			$response = OcisConfigHelper::reConfigureOcis($envs);
		}

		$resBody = $response->getBody()->getContents();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set config. Response: $resBody",
		);
	}

	/**
	 * @Given the administrator has started service :service separately with the following configs:
	 *
	 * @param string $service
	 * @param TableNode $table
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function theAdministratorHasStartedServiceSeparatelyWithTheFollowingConfig(
		string $service,
		TableNode $table,
	): void {
		$envs = [];
		foreach ($table->getHash() as $row) {
			$envs[$row['config']] = $row['value'];
		}

		$response = OcisConfigHelper::startService($service, $envs);
		$resBody = $response->getBody()->getContents();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to start service $service. Response: $resBody",
		);
	}

	/**
	 * @AfterScenario @env-config
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function rollback(): void {
		if (OcisConfigHelper::isK8s()) {
			$this->rollbackK8sServices();
			return;
		}
		$this->rollbackServices();
		$this->rollbackOcis();
	}

	/**
	 * @return void
	 * @throws GuzzleException
	 */
	public function rollbackOcis(): void {
		$response = OcisConfigHelper::rollbackOcis();
		$this->assertOcisRestarted(
			$response,
			"Failed to rollback ocis server. Check if oCIS is started with ociswrapper.",
		);
	}

	/**
	 * Asserts that an oCIS restart triggered by ociswrapper succeeded.
	 * In vault mode the wrapper's admin health-check uses basic-auth which is
	 * disabled (OIDC role assignment only), so the wrapper always returns HTTP 500
	 * even after a successful restart. When Keycloak/vault mode is detected we
	 * poll the unauthenticated proxy debug readyz endpoint instead of failing.
	 *
	 * @param ResponseInterface $response
	 * @param string $errorMessage
	 *
	 * @return void
	 */
	private function assertOcisRestarted(ResponseInterface $response, string $errorMessage): void {
		$statusCode = $response->getStatusCode();
		if ($statusCode === 200) {
			return;
		}
		if (KeycloakHelper::isTestingWithKeycloak()) {
			$this->waitForOcisProxyReady();
			return;
		}
		Assert::assertEquals(
			200,
			$statusCode,
			$errorMessage . " Response: " . $response->getBody()->getContents(),
		);
	}

	/**
	 * Poll the unauthenticated proxy debug readyz endpoint until oCIS is ready.
	 * Used in vault/Keycloak mode where basic-auth health checks are unavailable.
	 *
	 * @param int $timeoutSeconds
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	private function waitForOcisProxyReady(int $timeoutSeconds = 60): void {
		$readyzUrl = 'http://localhost:9205/readyz';
		$deadline = time() + $timeoutSeconds;
		while (time() < $deadline) {
			try {
				$response = HttpRequestHelper::get($readyzUrl);
				if ($response->getStatusCode() === 200) {
					return;
				}
				echo "oCIS not ready yet. Retrying in 1s...\n";
			} catch (\Exception $e) {
				throw new Exception("oCIS not ready. Error: $e");
			}
			sleep(1);
		}
		throw new \RuntimeException(
			"Timed out after {$timeoutSeconds}s waiting for oCIS proxy readyz at {$readyzUrl}",
		);
	}

	/**
	 * @return void
	 * @throws GuzzleException
	 */
	public function rollbackServices(): void {
		$response = OcisConfigHelper::rollbackServices();
		$resBody = $response->getBody()->getContents();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to rollback services. Response: $resBody",
		);
	}

	/**
	 * @return void
	 * @throws GuzzleException
	 */
	public function rollbackK8sServices(): void {
		$url = OcisConfigHelper::getWrapperUrl() . "/k8s/rollback";
		$response = OcisConfigHelper::sendRequest($url, "DELETE");
		$resBody = $response->getBody()->getContents();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to rollback services. Response: $resBody",
		);
	}
}
