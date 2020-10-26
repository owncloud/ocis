@api
Feature: PROPFIND

  @issue-ocis-751
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: PROPFIND to "/remote.php/dav/files"
    Given user "Alice" has been created with default attributes and without skeleton files
    When user "Alice" requests "/remote.php/dav/files" with "PROPFIND" using basic auth
    Then the HTTP status code should be "500"
