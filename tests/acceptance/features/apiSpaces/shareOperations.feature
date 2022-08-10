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


   Scenario: Uploading file to a group read-only share folder does not work
    Given using spaces DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has created a share with settings
      | path        | FOLDER |
      | shareType   | group  |
      | permissions | read   |
      | shareWith   | grp1   |
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"
    When user "Brian" uploads a file inside space "Shares Jail" with content "new description" to "/FOLDER/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "/FOLDER/textfile.txt" should not exist
   
  
  Scenario: Uploading file to a user upload-only share folder works
    Given using spaces DAV path
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has created a share with settings
      | path        | FOLDER |
      | shareType   | user   |
      | permissions | create |
      | shareWith   | Brian  |
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"
    When user "Brian" uploads a file inside space "Shares Jail" with content "new description" to "/FOLDER/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the following headers should match these regular expressions for user "Brian"
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    And the content of file "/FOLDER/textfile.txt" for user "Alice" should be:
    """
    new description
    """


  Scenario: Uploading file to a group upload-only share folder works
    Given using spaces DAV path
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has created a share with settings
      | path        | FOLDER |
      | shareType   | group  |
      | permissions | create |
      | shareWith   | grp1   |
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"
    When user "Brian" uploads a file inside space "Shares Jail" with content "new description" to "/FOLDER/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the following headers should match these regular expressions for user "Brian"
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    And the content of file "/FOLDER/textfile.txt" for user "Alice" should be:
    """
    new description
    """  


  Scenario: Uploading file to a user read/write share folder works
    Given using spaces DAV path
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has created a share with settings
      | path        | FOLDER |
      | shareType   | user   |
      | permissions | change |
      | shareWith   | Brian  |
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"
    When user "Brian" uploads a file inside space "Shares Jail" with content "new description" to "/FOLDER/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/FOLDER/textfile.txt" for user "Alice" should be:
    """
    new description
    """  


  Scenario: Uploading file to a group read/write share folder works
    Given using spaces DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has created a share with settings
      | path        | FOLDER |
      | shareType   | group  |
      | permissions | change |
      | shareWith   | grp1   |
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"
    When user "Brian" uploads a file inside space "Shares Jail" with content "new description" to "/FOLDER/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/FOLDER/textfile.txt" for user "Alice" should be:
    """
    new description
    """  
    

  Scenario: Check quota of owners parent directory of a shared file
    Given using spaces DAV path
    And the quota of user "Brian" has been set to "0"
    And user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/myfile.txt"
    And user "Alice" has shared file "myfile.txt" with user "Brian"
    And user "Brian" has accepted share "/myfile.txt" offered by user "Alice"
     When user "Brian" uploads a file inside space "Shares Jail" with content "new description" to "/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And the following headers should match these regular expressions for user "Brian"
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    And the content of file "/myfile.txt" for user "Alice" should be:
    """
    new description
    """


 Scenario: Uploading to a user shared folder with read/write permission when the sharer has unsufficient quota does not work
    Given using spaces DAV path
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has created a share with settings
      | path        | FOLDER |
      | shareType   | user   |
      | permissions | change |
      | shareWith   | Brian  |
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"
    And the quota of user "Alice" has been set to "0"
    When user "Brian" uploads a file inside space "Shares Jail" with content "new description" to "/FOLDER/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/textfile.txt" should not exist
  
  
  Scenario: Uploading to a user shared folder with upload-only permission when the sharer has insufficient quota does not work
    Given using spaces DAV path
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has created a share with settings
      | path        | FOLDER |
      | shareType   | user   |
      | permissions | create |
      | shareWith   | Brian  |
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"
    And the quota of user "Alice" has been set to "0"
    When user "Brian" uploads a file inside space "Shares Jail" with content "new description" to "/FOLDER/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/textfile.txt" should not exist

   
   Scenario: Uploading to a group shared folder with upload-only permission when the sharer has insufficient quota does not work
    Given using spaces DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has created a share with settings
      | path        | FOLDER |
      | shareType   | group  |
      | permissions | create |
      | shareWith   | grp1   |
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"
    And the quota of user "Alice" has been set to "0"
    When user "Brian" uploads a file inside space "Shares Jail" with content "new description" to "/FOLDER/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/textfile.txt" should not exist  

  
  Scenario Outline: Sharer can download file uploaded with different permission by sharee to a shared folder
    Given using spaces DAV path
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has created a share with settings
      | path        | FOLDER        |
      | shareType   | user          |
      | permissions | <permissions> |
      | shareWith   | Brian         |
    And user "Brian" has accepted share "/FOLDER" offered by user "Alice"
    When user "Brian" uploads a file inside space "Shares Jail" with content "some content" to "/FOLDER/textfile.txt" using the WebDAV API
    And user "Alice" downloads file "/FOLDER/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded content should be "some content" 
     Examples:
      | permissions |
      | change      |
      | create      |