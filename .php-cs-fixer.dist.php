<?php

$dirToParse = 'tests/acceptance/';
$dirIterator = new DirectoryIterator(__DIR__ . '/' . $dirToParse);

$excludeDirs = [
    'node_modules',
    'vendor-php'
];

$finder = PhpCsFixer\Finder::create()
    ->exclude($excludeDirs)
    ->in(__DIR__);

$ocRule = (new OC\CodingStandard\Config())->getRules();
$config = new PhpCsFixer\Config();
$config->setFinder($finder)
    ->setIndent("\t")
    ->setRules(
        array_merge(
            $ocRule,
            [
                "return_type_declaration" => [
                    "space_before" => "none",
                ],
                'single_space_around_construct' => true
            ]
        )
    );
$config->setFinder($finder);
return $config;
