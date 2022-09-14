@api @skipOnOcV10
Feature: upload resources using TUS protocol
  As a user
  I want to be able to upload files
  So that I can store and share files between multiple client systems

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And using spaces DAV path


  Scenario: upload a file within the set quota to a project space
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "10000"
    When user "Alice" uploads a file with content "uploaded content" to "/upload.txt" via TUS inside of the space "Project Jupiter" using the WebDAV API
    Then the HTTP status code should be "200"
    And for user "Alice" the space "Project Jupiter" should contain these entries:
      | upload.txt |


  Scenario: upload a file bigger than the set quota to a project space
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "10"
    When user "Alice" creates a new TUS resource for the space "Project Jupiter" using the WebDAV API with these headers:
      | Upload-Length   | 100                       |
      # dXBsb2FkLnR4dA== is the base64 encoded value of filename upload.txt
      | Upload-Metadata | filename dXBsb2FkLnR4dA== |
      | Tus-Resumable   | 1.0.0                     |
    Then the HTTP status code should be "507"
    And for user "Alice" the space "Project Jupiter" should not contain these entries:
      | upload.txt |

