<?php declare(strict_types=1);
/**
 * ownCloud
 *
 * @author Artur Neumann <artur@jankaritech.com>
 * @copyright Copyright (c) 2018 Artur Neumann artur@jankaritech.com
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
namespace TestHelpers\Asserts;

use PHPUnit\Framework\Assert;
use SimpleXMLElement;

/**
 * WebDAV related asserts
 */
class WebDav extends Assert {
	/**
	 *
	 * @param string|null $element exception|message|reason
	 * @param string|null $expectedValue
	 * @param array|null $responseXml
	 * @param string|null $extraErrorText
	 *
	 * @return void
	 */
	public static function assertDavResponseElementIs(
		?string $element,
		?string $expectedValue,
		?array $responseXml,
		?string $extraErrorText = ''
	):void {
		if ($extraErrorText !== '') {
			$extraErrorText = $extraErrorText . " ";
		}
		self::assertArrayHasKey(
			'value',
			$responseXml,
			$extraErrorText . "responseXml does not have key 'value'"
		);
		if ($element === "exception") {
			$result = $responseXml['value'][0]['value'];
		} elseif ($element === "message") {
			$result = $responseXml['value'][1]['value'];
		} elseif ($element === "reason") {
			$result = $responseXml['value'][3]['value'];
		} else {
			self::fail(__METHOD__ . " element must be one of exception, response or reason. But '$element' was passed in.");
		}

		self::assertEquals(
			$expectedValue,
			$result,
			__METHOD__ . " " . $extraErrorText . "Expected '$expectedValue' in element $element got '$result'"
		);
	}

	/**
	 *
	 * @param SimpleXMLElement $responseXmlObject
	 * @param array|null $expectedShareTypes
	 *
	 * @return void
	 */
	public static function assertResponseContainsShareTypes(
		SimpleXMLElement $responseXmlObject,
		?array $expectedShareTypes
	):void {
		foreach ($expectedShareTypes as $row) {
			$xmlPart = $responseXmlObject->xpath(
				"//d:prop/oc:share-types/oc:share-type[.=" . $row[0] . "]"
			);
			self::assertNotEmpty(
				$xmlPart,
				"cannot find share-type '" . $row[0] . "'"
			);
		}
	}
}
