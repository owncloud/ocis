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
		if (\getenv('K8S') === "true") {
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

		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set async upload with delayed post processing",
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

		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set config $configVariable=$configValue",
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

		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to enable role $role",
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

		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to enable roles: " . implode(', ', $roles),
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
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to disable role $role",
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
			// In K8s, use the path where the ConfigMap is mounted
			// For example: The banned-password-list.txt is mounted at /etc/ocis/config/drone/
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
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set config $configVariable=$path",
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

		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to set config",
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
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to start service $service.",
		);
	}

	/**
	 * @AfterScenario @env-config
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function rollback(): void {
		if (\getenv('K8S') === "true") {
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
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to rollback ocis server. Check if oCIS is started with ociswrapper.",
		);
	}

	/**
	 * @return void
	 * @throws GuzzleException
	 */
	public function rollbackServices(): void {
		$response = OcisConfigHelper::rollbackServices();
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to rollback services.",
		);
	}

	/**
	 * @return void
	 * @throws GuzzleException
	 */
	public function rollbackK8sServices(): void {
		$url = OcisConfigHelper::getWrapperUrl() . "/k8s/rollback";
		$response = OcisConfigHelper::sendRequest($url, "DELETE");
		Assert::assertEquals(
			200,
			$response->getStatusCode(),
			"Failed to rollback services.",
		);
	}
}
