@api
Feature: check etag propagation after different file alterations

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And user "Alice" has created folder "/upload"

  Scenario: copying a file inside a folder as a share receiver changes its etag for all collaborators
    Given user "Alice" has uploaded file with content "uploaded content" to "/upload/file.txt"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Alice" has stored etag of element "/upload/file.txt" inside space "Personal"
    And user "Alice" has stored etag of element "/upload/file.txt" on path "/upload/renamed.txt" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload/file.txt" inside space "Shares Jail"
    And user "Brian" has stored etag of element "upload/file.txt" on path "/upload/renamed.txt" inside space "Shares Jail"
    When user "Brian" copies file "/upload/file.txt" to "/upload/renamed.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path                | space       |
      | Alice | /                   | Personal    |
      | Alice | /upload             | Personal    |
      | Alice | /upload/renamed.txt | Personal    |
      | Brian | /                   | Shares Jail |
      | Brian | /upload             | Shares Jail |
      | Brian | /upload/renamed.txt | Shares Jail |
    And these etags should not have changed
      | user  | path             | space       |
      | Alice | /upload/file.txt | Personal    |
      | Brian | /upload/file.txt | Shares Jail |

  Scenario: copying a file inside a folder as a sharer changes its etag for all collaborators
    Given user "Alice" has uploaded file with content "uploaded content" to "/upload/file.txt"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Alice" has stored etag of element "/upload/file.txt" inside space "Personal"
    And user "Alice" has stored etag of element "/upload/file.txt" on path "/upload/renamed.txt" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload/file.txt" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload/file.txt" on path "/upload/renamed.txt" inside space "Shares Jail"
    When user "Alice" copies file "/upload/file.txt" to "/upload/renamed.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path                | space       |
      | Alice | /                   | Personal    |
      | Alice | /upload             | Personal    |
      | Alice | /upload/renamed.txt | Personal    |
      | Brian | /                   | Shares Jail |
      | Brian | /upload             | Shares Jail |
      | Brian | /upload/renamed.txt | Shares Jail |
    And these etags should not have changed
      | user  | path             | space       |
      | Alice | /upload/file.txt | Personal    |
      | Brian | /upload/file.txt | Shares Jail |


  Scenario: as share receiver renaming a file inside a folder changes its etag for all collaborators
    Given user "Alice" has uploaded file with content "uploaded content" to "/upload/file.txt"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Alice" has stored etag of element "/upload/file.txt" on path "/upload/renamed.txt" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload/file.txt" on path "/upload/renamed.txt" inside space "Shares Jail"
    When user "Brian" moves file "/upload/file.txt" to "/upload/renamed.txt" in space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path                | space       |
      | Alice | /                   | Personal    |
      | Alice | /upload             | Personal    |
      | Brian | /                   | Shares Jail |
      | Brian | /upload             | Shares Jail | 
    And these etags should not have changed
      | user  | path                | space       |
      | Alice | /upload/renamed.txt | Personal    |
      | Brian | /upload/renamed.txt | Shares Jail |


  Scenario: as sharer renaming a file inside a folder changes its etag for all collaborators
    Given user "Alice" has uploaded file with content "uploaded content" to "/upload/file.txt"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Alice" has stored etag of element "/upload/file.txt" on path "/upload/renamed.txt" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload/file.txt" on path "/upload/renamed.txt" inside space "Shares Jail"
    When user "Alice" moves file "/upload/file.txt" to "/upload/renamed.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path                | space       |
      | Alice | /                   | Personal    |
      | Alice | /upload             | Personal    |
      | Brian | /                   | Shares Jail |
      | Brian | /upload             | Shares Jail | 
    And these etags should not have changed
      | user  | path                | space       |
      | Alice | /upload/renamed.txt | Personal    |
      | Brian | /upload/renamed.txt | Shares Jail |


  Scenario: as sharer moving a file from one folder to an other changes the etags of both folders for all collaborators
    Given user "Alice" has created folder "/dst"
    And user "Alice" has uploaded file with content "uploaded content" to "/upload/file.txt"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has shared folder "/dst" with user "Brian"
    And user "Brian" has accepted share "/dst" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Alice" has stored etag of element "/dst" inside space "Personal"
    And user "Alice" has stored etag of element "/upload/file.txt" on path "/dst/file.txt" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/dst" inside space "Shares Jail"
    When user "Alice" moves file "/upload/file.txt" to "/dst/file.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path        | space       |
      | Alice | /           | Personal    |
      | Alice | /upload     | Personal    |
      | Alice | /dst        | Personal    |
      | Brian | /           | Shares Jail |
      | Brian | /upload     | Shares Jail |
      | Brian | /dst        | Shares Jail |


  Scenario: as share receiver moving a file from one folder to an other changes the etags of both folders for all collaborators
    Given user "Alice" has created folder "/dst"
    And user "Alice" has uploaded file with content "uploaded content" to "/upload/file.txt"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has shared folder "/dst" with user "Brian"
    And user "Brian" has accepted share "/dst" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Alice" has stored etag of element "/dst" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/dst" inside space "Shares Jail"
    When user "Brian" moves file "/upload/file.txt" to "/dst/file.txt" in space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path        | space       |
      | Alice | /           | Personal    |
      | Alice | /upload     | Personal    |
      | Alice | /dst        | Personal    |
      | Brian | /           | Shares Jail |
      | Brian | /upload     | Shares Jail |
      | Brian | /dst        | Shares Jail |


  Scenario: as sharer moving a folder from one folder to an other changes the etags of both folders for all collaborators
    Given user "Alice" has created folder "/dst"
    And user "Alice" has created folder "/upload/toMove"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has shared folder "/dst" with user "Brian"
    And user "Brian" has accepted share "/dst" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Alice" has stored etag of element "/dst" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/dst" inside space "Shares Jail"
    When user "Alice" moves file "/upload/toMove" to "/dst/toMove" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path        | space       |
      | Alice | /           | Personal    |
      | Alice | /upload     | Personal    |
      | Alice | /dst        | Personal    |
      | Brian | /           | Shares Jail |
      | Brian | /upload     | Shares Jail |
      | Brian | /dst        | Shares Jail |
  

  Scenario: as share reciever moving a folder from one folder to an other changes the etags of both folders for all collaborators
    Given user "Alice" has created folder "/dst"
    And user "Alice" has created folder "/upload/toMove"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has shared folder "/dst" with user "Brian"
    And user "Brian" has accepted share "/dst" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Alice" has stored etag of element "/dst" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/dst" inside space "Shares Jail"
    When user "Brian" moves file "/upload/toMove" to "/dst/toMove" in space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path        | space       |
      | Alice | /           | Personal    |
      | Alice | /upload     | Personal    |
      | Alice | /dst        | Personal    |
      | Brian | /           | Shares Jail |
      | Brian | /upload     | Shares Jail |
      | Brian | /dst        | Shares Jail |


  Scenario: as share receiver creating a folder inside a folder received as a share changes its etag for all collaborators
    Given user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    When user "Brian" creates a subfolder "/upload/new" in space "Shares Jail" using the WebDav Api
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path           | space       |
      | Alice | /              | Personal    |
      | Alice | /upload        | Personal    |
      | Brian | /              | Shares Jail |
      | Brian | /upload        | Shares Jail |

  
  Scenario: as sharer creating a folder inside a shared folder changes etag for all collaborators
    Given user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    When user "Alice" creates folder "/upload/new" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path           | space       |
      | Alice | /              | Personal    |
      | Alice | /upload        | Personal    |
      | Brian | /              | Shares Jail |
      | Brian | /upload        | Shares Jail |


  Scenario: as share receiver uploading a file inside a folder received as a share changes its etag for all collaborators
    Given user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    When user "Brian" uploads a file inside space "Shares Jail" with content "uploaded content" to "/upload/file.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path           | space       |
      | Alice | /              | Personal    |
      | Alice | /upload        | Personal    |
      | Brian | /              | Shares Jail |
      | Brian | /upload        | Shares Jail |    


  Scenario: as sharer uploading a file inside a shared folder should update etags for all collaborators    
    Given user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    When user "Alice" uploads file with content "uploaded content" to "/upload/file.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And these etags should have changed
      | user  | path           | space       |
      | Alice | /              | Personal    |
      | Alice | /upload        | Personal    |
      | Brian | /              | Shares Jail |
      | Brian | /upload        | Shares Jail |


  Scenario: as share receiver overwriting a file inside a received shared folder should update etags for all collaborators
    Given user "Alice" has uploaded file with content "uploaded content" to "/upload/file.txt"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    When user "Brian" uploads a file inside space "Shares Jail" with content "new content" to "/upload/file.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And these etags should have changed
      | user  | path           | space       |
      | Alice | /              | Personal    |
      | Alice | /upload        | Personal    |
      | Brian | /              | Shares Jail |
      | Brian | /upload        | Shares Jail |

  
  Scenario: as sharer overwriting a file inside a shared folder should update etags for all collaborators    
    Given user "Alice" has uploaded file with content "uploaded content" to "/upload/file.txt"
    And user "Alice" has shared folder "/upload" with user "Brian"
    And user "Brian" has accepted share "/upload" offered by user "Alice"
    And user "Alice" has stored etag of element "/" inside space "Personal"
    And user "Alice" has stored etag of element "/upload" inside space "Personal"
    And user "Brian" has stored etag of element "/" inside space "Shares Jail"
    And user "Brian" has stored etag of element "/upload" inside space "Shares Jail"
    When user "Alice" uploads file with content "new content" to "/upload/file.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And these etags should have changed
      | user  | path           | space       |
      | Alice | /              | Personal    |
      | Alice | /upload        | Personal    |
      | Brian | /              | Shares Jail |
      | Brian | /upload        | Shares Jail |
