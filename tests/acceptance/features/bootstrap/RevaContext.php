<?php

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use TestHelpers\AppConfigHelper;
use TestHelpers\SetupHelper;

require_once 'bootstrap.php';

/**
 * Context for Reva specific steps
 */
class RevaContext implements Context {

    /**
     * @var FeatureContext
     */
    private $featureContext;

    /**
     * @BeforeScenario
     *
     * @param BeforeScenarioScope $scope
     *
     * @return void
     * @throws Exception
     */
    public function setUpScenario(BeforeScenarioScope $scope): void {
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
}
