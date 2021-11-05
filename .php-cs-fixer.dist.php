<?php

$dirToParse = 'tests/acceptance/';
$dirIterator = new DirectoryIterator(__DIR__ . '/' . $dirToParse);

$excludeDirs = [
    'node_modules'
];

$finder = PhpCsFixer\Finder::create()
    ->exclude($excludeDirs)
    ->in(__DIR__);

$config = new OC\CodingStandard\Config();
$config->setFinder($finder);
return $config;
