@api @skipOnOcV10
Feature: sharing

  Background:
    Given using OCS API version "1"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path  


  Scenario: Correct webdav share-permissions for received file with edit and reshare permissions
    Given user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has shared file "/tmp.txt" with user "Brian"
    And user "Brian" has accepted share "/tmp.txt" offered by user "Alice"
    When user "Brian" gets the following properties of file "/tmp.txt" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "19"


  Scenario: Correct webdav share-permissions for received group shared file with edit and reshare permissions
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has created a share with settings
      | path        | /tmp.txt          |
      | shareType   | group             |
      | permissions | share,update,read |
      | shareWith   | grp1              |
    And user "Brian" has accepted share "/tmp.txt" offered by user "Alice"
    When user "Brian" gets the following properties of file "/tmp.txt" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "19"


  Scenario: Correct webdav share-permissions for received file with edit permissions but no reshare permissions
    Given user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has shared file "tmp.txt" with user "Brian"
    And user "Brian" has accepted share "/tmp.txt" offered by user "Alice"
    When user "Alice" updates the last share using the sharing API with
      | permissions | update,read |
    Then the HTTP status code should be "200"
    And as user "Brian" file "/tmp.txt" inside space "Shares Jail" should contain a property "ocs:share-permissions" with value "3"


  Scenario: Correct webdav share-permissions for received group shared file with edit permissions but no reshare permissions
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has created a share with settings
      | path        | /tmp.txt    |
      | shareType   | group       |
      | permissions | update,read |
      | shareWith   | grp1        |
    And user "Brian" has accepted share "/tmp.txt" offered by user "Alice"
    When user "Brian" gets the following properties of file "/tmp.txt" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "3"


  Scenario: Correct webdav share-permissions for received file with reshare permissions but no edit permissions
    Given user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has shared file "tmp.txt" with user "Brian"
    And user "Brian" has accepted share "/tmp.txt" offered by user "Alice"
    When user "Alice" updates the last share using the sharing API with
      | permissions | share,read |
    Then the HTTP status code should be "200"
    And as user "Brian" file "/tmp.txt" inside space "Shares Jail" should contain a property "ocs:share-permissions" with value "17"
  

  Scenario: Correct webdav share-permissions for received group shared file with reshare permissions but no edit permissions
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has created a share with settings
      | path        | /tmp.txt   |
      | shareType   | group      |
      | permissions | share,read |
      | shareWith   | grp1       |
    And user "Brian" has accepted share "/tmp.txt" offered by user "Alice"
    When user "Brian" gets the following properties of file "/tmp.txt" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "17"


  Scenario: Correct webdav share-permissions for received folder with all permissions
    Given user "Alice" has created folder "/tmp"
    And user "Alice" has shared file "/tmp" with user "Brian"
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
    When user "Brian" gets the following properties of folder "/tmp" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "31"
    

  Scenario: Correct webdav share-permissions for received group shared folder with all permissions
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path      | tmp   |
      | shareType | group |
      | shareWith | grp1  |
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
    When user "Brian" gets the following properties of folder "/tmp" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "31"


  Scenario: Correct webdav share-permissions for received folder with all permissions but edit
    Given user "Alice" has created folder "/tmp"
    And user "Alice" has shared file "/tmp" with user "Brian"
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
    When user "Alice" updates the last share using the sharing API with
      | permissions | share,delete,create,read |
    Then the HTTP status code should be "200"
    And as user "Brian" folder "/tmp" inside space "Shares Jail" should contain a property "ocs:share-permissions" with value "29"


  Scenario: Correct webdav share-permissions for received group shared folder with all permissions but edit
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path        | tmp                      |
      | shareType   | group                    |
      | shareWith   | grp1                     |
      | permissions | share,delete,create,read |
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
     When user "Brian" gets the following properties of folder "/tmp" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "29"
   

  Scenario: Correct webdav share-permissions for received folder with all permissions but create
    Given user "Alice" has created folder "/tmp"
    And user "Alice" has shared file "/tmp" with user "Brian"
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
    When user "Alice" updates the last share using the sharing API with
      | permissions | share,delete,update,read |
    Then the HTTP status code should be "200"
     And as user "Brian" folder "/tmp" inside space "Shares Jail" should contain a property "ocs:share-permissions" with value "27"
   

  Scenario: Correct webdav share-permissions for received group shared folder with all permissions but create
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path        | tmp                      |
      | shareType   | group                    |
      | shareWith   | grp1                     |
      | permissions | share,delete,update,read |
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
    When user "Brian" gets the following properties of folder "/tmp" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "27"
   

  Scenario: Correct webdav share-permissions for received folder with all permissions but delete
    Given user "Alice" has created folder "/tmp"
    And user "Alice" has shared file "/tmp" with user "Brian"
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
    When user "Alice" updates the last share using the sharing API with
      | permissions | share,create,update,read |
    Then the HTTP status code should be "200"
    And as user "Brian" folder "/tmp" inside space "Shares Jail" should contain a property "ocs:share-permissions" with value "23"
   
   
  Scenario: Correct webdav share-permissions for received group shared folder with all permissions but delete
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path        | tmp                      |
      | shareType   | group                    |
      | shareWith   | grp1                     |
      | permissions | share,create,update,read |
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
    When user "Brian" gets the following properties of folder "/tmp" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "23"
    

  Scenario: Correct webdav share-permissions for received folder with all permissions but share
    Given user "Alice" has created folder "/tmp"
    And user "Alice" has shared file "/tmp" with user "Brian"
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
    When user "Alice" updates the last share using the sharing API with
      | permissions | change |
    Then the HTTP status code should be "200"
    And as user "Brian" folder "/tmp" inside space "Shares Jail" should contain a property "ocs:share-permissions" with value "15"
   

  Scenario: Correct webdav share-permissions for received group shared folder with all permissions but share
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path        | tmp    |
      | shareType   | group  |
      | shareWith   | grp1   |
      | permissions | change |
    And user "Brian" has accepted share "/tmp" offered by user "Alice"
    When user "Brian" gets the following properties of folder "/tmp" inside space "Shares Jail" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "15"