Feature: an user changes its own password
  As a user
  I want to change my password
  So that I can use new combination as password


  Scenario Outline: change own password
    Given user "Alice" has been created with default attributes and without skeleton files
    When the user "Alice" changes its own password "<current-password>" to "<new-password>" using the Graph API
    Then the HTTP status code should be "<http-status-code>"
    Examples:
      | current-password | new-password | http-status-code |
      | 123456           | validPass    | 204              |
      | 123456           | кириллица    | 204              |
      | 123456           | 密码           | 204              |
      | 123456           | ?&^%0        | 204              |
      | 123456           |              | 400              |
      | 123456           | 123456       | 400              |
      | wrongPass        | 123456       | 400              |
      |                  | validPass    | 400              |
