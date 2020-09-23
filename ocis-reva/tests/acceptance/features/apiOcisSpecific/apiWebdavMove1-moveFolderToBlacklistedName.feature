@api @issue-ocis-reva-14
Feature: users cannot move (rename) a folder to a blacklisted name
  As an administrator
  I want to be able to prevent users from moving (renaming) folders to specified names
  So that I can prevent unwanted folder names existing in the cloud storage

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @issue-ocis-reva-211 @skipOnOcis-EOS-Storage @issue-ocis-reva-269 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Renaming a folder to a name that is banned by default is allowed
    Given using <dav_version> DAV path
    And user "Alice" has created folder "/testshare"
    When user "Alice" moves folder "/testshare" to "/.htaccess" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" folder "/.htaccess" should exist
    Examples:
      | dav_version |
      | old         |
      | new         |
