@skipOnReva
Feature: sharing works when a username and group name are the same
  As a user
  I want to share resources with group and users having same name
  So that I can make sure that the sharing works

  Background:
    Given user "Alice" has been created with default attributes


  Scenario: creating a new share with user and a group having same name
    Given these users have been created with default attributes:
      | username |
      | Brian    |
      | Carol    |
    And group "Brian" has been created
    And user "Carol" has been added to group "Brian"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | randomfile.txt |
      | space           | Personal       |
      | sharee          | Brian          |
      | shareType       | group          |
      | permissionsRole | File Editor    |
    And user "Carol" has a share "randomfile.txt" synced
    When user "Alice" shares file "randomfile.txt" with user "Brian" using the sharing API
    Then the HTTP status code should be "200"
    And the OCS status code should be "100"
    And user "Brian" should see the following elements
      | /Shares/randomfile.txt |
    And user "Carol" should see the following elements
      | /Shares/randomfile.txt |
    And the content of file "/Shares/randomfile.txt" for user "Brian" should be "Random data"
    And the content of file "/Shares/randomfile.txt" for user "Carol" should be "Random data"


  Scenario: creating a new share with group and a user having same name
    Given these users have been created with default attributes:
      | username |
      | Brian    |
      | Carol    |
    And group "Brian" has been created
    And user "Carol" has been added to group "Brian"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | randomfile.txt |
      | space           | Personal       |
      | sharee          | Brian          |
      | shareType       | user           |
      | permissionsRole | File Editor    |
    And user "Brian" has a share "randomfile.txt" synced
    When user "Alice" shares file "randomfile.txt" with group "Brian" using the sharing API
    Then the HTTP status code should be "200"
    And the OCS status code should be "100"
    And user "Brian" should see the following elements
      | /Shares/randomfile.txt |
    And user "Carol" should see the following elements
      | /Shares/randomfile.txt |
    And the content of file "/Shares/randomfile.txt" for user "Brian" should be "Random data"
    And the content of file "/Shares/randomfile.txt" for user "Carol" should be "Random data"


  Scenario: creating a new share with user and a group having same name but different case
    Given these users have been created with default attributes:
      | username |
      | Brian    |
      | Carol    |
    And group "brian" has been created
    And user "Carol" has been added to group "brian"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | randomfile.txt |
      | space           | Personal       |
      | sharee          | brian          |
      | shareType       | group          |
      | permissionsRole | File Editor    |
    And user "Carol" has a share "randomfile.txt" synced
    When user "Alice" shares file "randomfile.txt" with user "Brian" using the sharing API
    Then the HTTP status code should be "200"
    And the OCS status code should be "100"
    And user "Brian" should see the following elements
      | /Shares/randomfile.txt |
    And user "Carol" should see the following elements
      | /Shares/randomfile.txt |
    And the content of file "/Shares/randomfile.txt" for user "Brian" should be "Random data"
    And the content of file "/Shares/randomfile.txt" for user "Carol" should be "Random data"


  Scenario: creating a new share with group and a user having same name but different case
    Given these users have been created with default attributes:
      | username |
      | Brian    |
      | Carol    |
    And group "brian" has been created
    And user "Carol" has been added to group "brian"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | randomfile.txt |
      | space           | Personal       |
      | sharee          | Brian          |
      | shareType       | user           |
      | permissionsRole | File Editor    |
    And user "Brian" has a share "randomfile.txt" synced
    When user "Alice" shares file "randomfile.txt" with group "brian" using the sharing API
    Then the HTTP status code should be "200"
    And the OCS status code should be "100"
    And user "Carol" should see the following elements
      | /Shares/randomfile.txt |
    And user "Brian" should see the following elements
      | /Shares/randomfile.txt |
    And the content of file "/Shares/randomfile.txt" for user "Carol" should be "Random data"
    And the content of file "/Shares/randomfile.txt" for user "Brian" should be "Random data"
