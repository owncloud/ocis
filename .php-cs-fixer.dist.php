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
                'single_space_around_construct' => true,
                'no_unused_imports' => true,
                'array_indentation' => true,
                'method_chaining_indentation' => true,
                'trailing_comma_in_multiline' => [
                    'elements' => ['arrays', 'arguments', 'parameters'],
                ],
                'no_useless_else' => true,
                'single_line_comment_spacing' => true,
                'no_trailing_whitespace_in_comment' => true,
                'no_empty_comment' => true,
                'no_singleline_whitespace_before_semicolons' => true,
                'type_declaration_spaces' => true,
                'binary_operator_spaces' => true,
                'phpdoc_to_return_type' => true,
                'void_return' => true,
                'no_useless_concat_operator' => true,
                'concat_space' => [
                    "spacing" => "one",
                ],
            ]
        )
    );
$config->setFinder($finder);
return $config;
