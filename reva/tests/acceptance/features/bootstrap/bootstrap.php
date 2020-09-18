<?php
$pathToCore = \getenv('PATH_TO_CORE');
if ($pathToCore === false) {
    $pathToCore = "../core";
}

require_once $pathToCore . '/tests/acceptance/features/bootstrap/bootstrap.php';

$classLoader = new \Composer\Autoload\ClassLoader();
$classLoader->addPsr4(
    "", $pathToCore . "/tests/acceptance/features/bootstrap", true
);

$classLoader->register();
