@api
Feature: users cannot move (rename) a folder to a blacklisted name
  As an administrator
  I want to be able to prevent users from moving (renaming) folders to specified names
  So that I can prevent unwanted folder names existing in the cloud storage

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: Rename a folder to a name that is banned by default
    Given using <dav_version> DAV path
    And user "Alice" has created folder "/testshare"
    When user "Alice" moves folder "/testshare" to "/.htaccess" using the WebDAV API
    Then the HTTP status code should be "403"
    And user "Alice" should see the following elements
      | /testshare/ |
    Examples:
      | dav_version |
      | old         |
      | new         |

    @skipOnOcV10 @personalSpace
    Examples:
      | dav_version |
      | spaces      |
