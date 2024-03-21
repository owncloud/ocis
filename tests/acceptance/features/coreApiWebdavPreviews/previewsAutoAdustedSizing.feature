@skipOnReva
Feature: sizing of previews of files downloaded through the webdav API
  As a user
  I want the aspect-ratio of previews to be preserved even when I ask for an unusual preview size
  So that the previews always have a similar look-and-feel to the original file

  This is optional behavior of an implementation. OCIS happens like this,
  but oC10 does not do this auto-fix of the aspect ratio.

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: download different sizes of previews of file
    Given user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
    When user "Alice" downloads the preview of "/parent.txt" with width <request-width> and height <request-height> using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded image should be <return-width> pixels wide and <return-height> pixels high
    Examples:
      | request-width | request-height | return-width | return-height |
      | 1             | 1              | 16           | 16            |
      | 32            | 32             | 32           | 32            |
      | 1024          | 1024           | 640          | 480           |
      | 1             | 1024           | 16           | 16            |
      | 1024          | 1              | 640          | 480           |
