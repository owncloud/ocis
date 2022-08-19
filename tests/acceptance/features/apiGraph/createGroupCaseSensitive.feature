@api
Feature: create groups, group names are case insensitive

  Scenario Outline: group names are case insensitive, creating groups with different upper and lower case names
    Given using OCS API version "<ocs_api_version>"
    And group "<group_id1>" has been created
    When the administrator creates a group "<group_id2>" using the Graph API
    And the administrator creates a group "<group_id3>" using the Graph API
    Then the HTTP status code of responses on all endpoints should be "400"
    And these groups should not exist:
    | groupname   |
    | <group_id2> |
    | <group_id3> |
    Examples:
      | ocs_api_version | group_id1            | group_id2            | group_id3            |
      | 1               | case-sensitive-group | Case-Sensitive-Group | CASE-SENSITIVE-GROUP |
      | 1               | Case-Sensitive-Group | CASE-SENSITIVE-GROUP | case-sensitive-group |
      | 1               | CASE-SENSITIVE-GROUP | case-sensitive-group | Case-Sensitive-Group |
      | 2               | case-sensitive-group | Case-Sensitive-Group | CASE-SENSITIVE-GROUP |
      | 2               | Case-Sensitive-Group | CASE-SENSITIVE-GROUP | case-sensitive-group |
      | 2               | CASE-SENSITIVE-GROUP | case-sensitive-group | Case-Sensitive-Group |
