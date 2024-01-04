Feature: moving/renaming file using file id
  As a user
  I want to be able to move or rename files using file id
  So that I can manage my file system

  Background:
    Given using spaces DAV path
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: move a file into a folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file into a sub-folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder/sub-folder" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder/sub-folder/" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from sub-folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/sub-folder/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder/sub-folder" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from personal to share space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has shared folder "/folder" with user "Brian" with permissions "<permissions>"
    And user "Brian" has uploaded file with content "some data" to "/test.txt"
    And we save it into "FILEID"
    When user "Brian" moves a file "test.txt" into "folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And the value of the item "/d:error/s:message" in the response about user "Brian" should be "cross storage moves are not permitted, use copy and delete"
    And for user "Brian" folder "/" of the space "Personal" should contain these files:
      | test.txt |
    But for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | test.txt |
    Examples:
      | permissions | dav-path                          |
      | all         | /remote.php/dav/spaces/<<FILEID>> |
      | all         | /dav/spaces/<<FILEID>>            |
      | change      | /remote.php/dav/spaces/<<FILEID>> |
      | change      | /dav/spaces/<<FILEID>>            |
      | read        | /remote.php/dav/spaces/<<FILEID>> |
      | read        | /dav/spaces/<<FILEID>>            |

  @issue-7618
  Scenario Outline: move a file from personal to project space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Brian" has uploaded a file inside space "Personal" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "project-space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" moves a file "textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "/" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | role    | http-status-code | dav-path                          |
      | manager | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | viewer  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | manager | 502              | /dav/spaces/<<FILEID>>            |
      | editor  | 502              | /dav/spaces/<<FILEID>>            |
      | viewer  | 403              | /dav/spaces/<<FILEID>>            |

  @issue-7618
  Scenario Outline: move a file to different name from personal space to project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Personal" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "/textfile.txt" into "/renamed.txt" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And the value of the item "/d:error/s:message" in the response about user "Alice" should be "move:error: not supported: cannot move across spaces"
    And for user "Alice" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Alice" folder "/" of the space "project-space" should not contain these files:
      | renamed.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file into a folder inside project space (manager/editor)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "/folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "project-space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" moves a file "/textfile.txt" into "/folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "project-space" should not contain these entries:
      | textfile.txt |
    Examples:
      | role    | dav-path                          |
      | manager | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | /remote.php/dav/spaces/<<FILEID>> |
      | manager | /dav/spaces/<<FILEID>>            |
      | editor  | /dav/spaces/<<FILEID>>            |


  Scenario Outline: try to move a file into a folder inside project space (viewer)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "/folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "project-space" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    When user "Brian" moves a file "/textfile.txt" into "/folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And for user "Alice" the space "project-space" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file into a sub-folder inside project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder/sub-folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder/sub-folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder/sub-folder/" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from folder to root inside project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "project-space" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from sub-folder to root inside project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder/sub-folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/sub-folder/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "project-space" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder/sub-folder" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |

  @issue-8116
  Scenario Outline: move a file between two project spaces
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "first-project-space" with the default quota using the Graph API
    And user "Alice" has created a space "second-project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "first-project-space" with content "first project space" to "file.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "first-project-space" with settings:
      | shareWith | Brian       |
      | role      | <from_role> |
    And user "Alice" has shared a space "second-project-space" with settings:
      | shareWith | Brian       |
      | role      | <to_role>   |
    When user "Brian" moves a file "file.txt" into "/" inside space "second-project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" the space "second-project-space" should not contain these entries:
      | file.txt |
    But for user "Brian" the space "first-project-space" should contain these entries:
      | file.txt |
    Examples:
      | from_role | to_role | http-status-code | dav-path                          |
      | manager   | manager | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | editor    | manager | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | manager   | editor  | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | editor    | editor  | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | manager   | viewer  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | editor    | viewer  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | viewer    | manager | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | viewer    | editor  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | viewer    | viewer  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | manager   | manager | 502              | /dav/spaces/<<FILEID>>            |
      | editor    | manager | 502              | /dav/spaces/<<FILEID>>            |
      | manager   | editor  | 502              | /dav/spaces/<<FILEID>>            |
      | editor    | editor  | 502              | /dav/spaces/<<FILEID>>            |
      | manager   | viewer  | 403              | /dav/spaces/<<FILEID>>            |
      | editor    | viewer  | 403              | /dav/spaces/<<FILEID>>            |
      | viewer    | manager | 403              | /dav/spaces/<<FILEID>>            |
      | viewer    | editor  | 403              | /dav/spaces/<<FILEID>>            |
      | viewer    | viewer  | 403              | /dav/spaces/<<FILEID>>            |

  @issue-8116
  Scenario Outline: move a file to different name between project spaces
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "first-project-space" with the default quota using the Graph API
    And user "Alice" has created a space "second-project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "first-project-space" with content "data from first project space" to "firstProjectSpacetextfile.txt"
    And user "Alice" has uploaded a file inside space "second-project-space" with content "data from second project space" to "secondProjectSpacetextfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "/secondProjectSpacetextfile.txt" into "/renamedSecondProjectSpacetextfile.txt" inside space "first-project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And the value of the item "/d:error/s:message" in the response about user "Alice" should be "move:error: not supported: cannot move across spaces"
    And for user "Alice" folder "/" of the space "first-project-space" should contain these files:
      | firstProjectSpacetextfile.txt |
    And for user "Alice" folder "/" of the space "second-project-space" should contain these files:
      | secondProjectSpacetextfile.txt |
    But for user "Alice" the space "first-project-space" should not contain these entries:
      | renamedSecondProjectSpacetextfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from project to shares space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "project-space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    And user "Alice" has created folder "testshare"
    And user "Alice" has shared folder "testshare" with user "Brian" with permissions "<permissions>"
    When user "Brian" moves a file "textfile.txt" into "testshare" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And for user "Brian" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "testshare" of the space "Shares" should not contain these files:
      | textfile.txt |
    Examples:
      | role    | permissions | dav-path                          |
      | manager | all         | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | all         | /remote.php/dav/spaces/<<FILEID>> |
      | viewer  | all         | /remote.php/dav/spaces/<<FILEID>> |
      | manager | change      | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | change      | /remote.php/dav/spaces/<<FILEID>> |
      | viewer  | change      | /remote.php/dav/spaces/<<FILEID>> |
      | manager | read        | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | read        | /remote.php/dav/spaces/<<FILEID>> |
      | viewer  | read        | /remote.php/dav/spaces/<<FILEID>> |
      | manager | all         | /dav/spaces/<<FILEID>>            |
      | editor  | all         | /dav/spaces/<<FILEID>>            |
      | viewer  | all         | /dav/spaces/<<FILEID>>            |
      | manager | change      | /dav/spaces/<<FILEID>>            |
      | editor  | change      | /dav/spaces/<<FILEID>>            |
      | viewer  | change      | /dav/spaces/<<FILEID>>            |
      | manager | read        | /dav/spaces/<<FILEID>>            |
      | editor  | read        | /dav/spaces/<<FILEID>>            |
      | viewer  | read        | /dav/spaces/<<FILEID>>            |

  @issue-7618
  Scenario Outline: move a file from project to personal space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "project-space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" moves a file "/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "/" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | role    | http-status-code | dav-path                          |
      | manager | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | viewer  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | manager | 502              | /dav/spaces/<<FILEID>>            |
      | editor  | 502              | /dav/spaces/<<FILEID>>            |
      | viewer  | 403              | /dav/spaces/<<FILEID>>            |

  @issue-7618
  Scenario Outline: move a file to different name from project space to personal space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "/textfile.txt" into "/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And the value of the item "/d:error/s:message" in the response about user "Alice" should be "move:error: not supported: cannot move across spaces"
    And for user "Alice" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Alice" folder "/" of the space "Personal" should not contain these files:
      | renamed.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |

  @issue-7617
  Scenario Outline: move a file into a folder within a shared folder (all/change permissions)
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "folder" with user "Brian" with permissions "<permissions>"
    When user "Brian" moves a file "Shares/folder/test.txt" into "folder/sub-folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Brian" folder "folder/sub-folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should contain these files:
      | test.txt |
    But for user "Brian" folder "folder" of the space "Shares" should not contain these files:
      | test.txt |
    And for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | test.txt |
    Examples:
      | permissions | dav-path                          |
      | all         | /remote.php/dav/spaces/<<FILEID>> |
      | all         | /dav/spaces/<<FILEID>>            |
      | change      | /remote.php/dav/spaces/<<FILEID>> |
      | change      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: try to move a file into a folder within a shared folder (read permissions)
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "folder" with user "Brian" with permissions "read"
    When user "Brian" moves a file "Shares/folder/test.txt" into "folder/sub-folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And for user "Brian" folder "folder/sub-folder" of the space "Shares" should not contain these files:
      | test.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should not contain these files:
      | test.txt |
    But for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | test.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from one shared folder to another shared folder
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "testshare1"
    And user "Alice" has created folder "testshare2"
    And user "Alice" has uploaded file with content "some data" to "testshare1/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "testshare1" with user "Brian" with permissions "<from_permissions>"
    And user "Alice" has shared folder "testshare2" with user "Brian" with permissions "<to_permissions>"
    When user "Brian" moves a file "Shares/testshare1/textfile.txt" into "testshare2" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And for user "Brian" folder "testshare1" of the space "Shares" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "testshare2" of the space "Shares" should not contain these files:
      | textfile.txt |
    Examples:
      | from_permissions | to_permissions | dav-path                          |
      | all              | all            | /remote.php/dav/spaces/<<FILEID>> |
      | all              | change         | /remote.php/dav/spaces/<<FILEID>> |
      | all              | read           | /remote.php/dav/spaces/<<FILEID>> |
      | change           | all            | /remote.php/dav/spaces/<<FILEID>> |
      | change           | change         | /remote.php/dav/spaces/<<FILEID>> |
      | change           | read           | /remote.php/dav/spaces/<<FILEID>> |
      | read             | all            | /remote.php/dav/spaces/<<FILEID>> |
      | read             | change         | /remote.php/dav/spaces/<<FILEID>> |
      | read             | read           | /remote.php/dav/spaces/<<FILEID>> |
      | all              | all            | /dav/spaces/<<FILEID>>            |
      | all              | change         | /dav/spaces/<<FILEID>>            |
      | all              | read           | /dav/spaces/<<FILEID>>            |
      | change           | all            | /dav/spaces/<<FILEID>>            |
      | change           | change         | /dav/spaces/<<FILEID>>            |
      | change           | read           | /dav/spaces/<<FILEID>>            |
      | read             | all            | /dav/spaces/<<FILEID>>            |
      | read             | change         | /dav/spaces/<<FILEID>>            |
      | read             | read           | /dav/spaces/<<FILEID>>            |

  @issue-8124
  Scenario Outline: move a file from share to personal space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "/folder" with user "Brian" with permissions "<permissions>"
    When user "Brian" moves a file "Shares/folder/test.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "/" of the space "Personal" should not contain these files:
      | test.txt |
    Examples:
      | permissions | dav-path                          | http-status-code |
      | all         | /remote.php/dav/spaces/<<FILEID>> | 502              |
      | all         | /dav/spaces/<<FILEID>>            | 502              |
      | change      | /remote.php/dav/spaces/<<FILEID>> | 502              |
      | change      | /dav/spaces/<<FILEID>>            | 502              |
      | read        | /remote.php/dav/spaces/<<FILEID>> | 403              |
      | read        | /dav/spaces/<<FILEID>>            | 403              |

  @issue-8125
  Scenario Outline: move a file from shares to project space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has shared a space "project-space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    And user "Alice" has created folder "testshare"
    And user "Alice" has uploaded file with content "some data" to "testshare/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "testshare" with user "Brian" with permissions "<permissions>"
    When user "Brian" moves a file "Shares/testshare/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" folder "testshare" of the space "Shares" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "/" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | role    | permissions | http-status-code | dav-path                          |
      | manager | all         | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | all         | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | viewer  | all         | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | manager | change      | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | change      | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | viewer  | change      | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | manager | read        | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | read        | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | viewer  | read        | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | manager | all         | 502              | /dav/spaces/<<FILEID>>            |
      | editor  | all         | 502              | /dav/spaces/<<FILEID>>            |
      | viewer  | all         | 403              | /dav/spaces/<<FILEID>>            |
      | manager | change      | 502              | /dav/spaces/<<FILEID>>            |
      | editor  | change      | 502              | /dav/spaces/<<FILEID>>            |
      | viewer  | change      | 403              | /dav/spaces/<<FILEID>>            |
      | manager | read        | 403              | /dav/spaces/<<FILEID>>            |
      | editor  | read        | 403              | /dav/spaces/<<FILEID>>            |
      | viewer  | read        | 403              | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a root file inside personal space
    Given user "Alice" has uploaded file with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "textfile.txt" into "renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | renamed.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move into a folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "textfile.txt" into "/folder/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | renamed.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move into a sub-folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "textfile.txt" into "/folder/sub-folder/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should contain these files:
      | renamed.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move from a folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "folder/textfile.txt" into "/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | renamed.txt |
    But for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | renamed.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move from sub-folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "folder/sub-folder/textfile.txt" into "/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these files:
      | renamed.txt |
    But for user "Alice" folder "folder/sub-folder" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |

  @issue-7617
  Scenario Outline: move a file to a different name into a sub-folder inside share space (all/change permissions)
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has created folder "/folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "/folder" with user "Brian" with permissions "<permissions>"
    When user "Brian" renames a file "Shares/folder/test.txt" into "folder/sub-folder/renamed.txt" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Brian" folder "folder/sub-folder" of the space "Shares" should contain these files:
      | renamed.txt |
    But for user "Brian" folder "folder" of the space "Shares" should not contain these files:
      | test.txt |
    Examples:
      | permissions | dav-path                          |
      | all         | /remote.php/dav/spaces/<<FILEID>> |
      | all         | /dav/spaces/<<FILEID>>            |
      | change      | /remote.php/dav/spaces/<<FILEID>> |
      | change      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file to a different name into a sub-folder inside share space (read permissions)
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has created folder "/folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "/folder" with user "Brian" with permissions "read"
    When user "Brian" renames a file "Shares/folder/test.txt" into "folder/sub-folder/renamed.txt" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    But for user "Brian" folder "folder/sub-folder" of the space "Shares" should not contain these files:
      | renamed.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |
