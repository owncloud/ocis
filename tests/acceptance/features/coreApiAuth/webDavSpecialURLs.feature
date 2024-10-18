Feature: make webdav request with special urls
  As a user
  I want to make webdav request with special urls
  So that I can make sure that they work

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "some data" to "/textfile1.txt"
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has uploaded file with content "some data" to "/PARENT/parent.txt"


  Scenario: send DELETE requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "DELETE" using password "%regular%" about user "Alice"
      | endpoint                                 |
      | //webdav/textfile0.txt                   |
      | //dav//files/%username%/textfile1.txt    |
      | /dav//files/%username%/PARENT/parent.txt |
      | /webdav//PARENT                          |
      | //dav/files/%username%//FOLDER           |
    Then the HTTP status code of responses on each endpoint should be "200,200,204,204,200" on oCIS or "204,204,204,204,204" on reva


  Scenario: send DELETE requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "DELETE" using password "%regular%" about user "Alice"
      | endpoint                                  |
      | //dav/spaces/%spaceid%/textfile0.txt      |
      | //dav//spaces/%spaceid%/PARENT/parent.txt |
      | /dav//spaces/%spaceid%/PARENT             |
      | //dav/spaces/%spaceid%//FOLDER            |
    Then the HTTP status code of responses on each endpoint should be "200,200,204,200" on oCIS or "204,204,204,204" on reva


  Scenario: send GET requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "GET" using password "%regular%" about user "Alice"
      | endpoint                                 |
      | //webdav/textfile0.txt                   |
      | //dav//files/%username%/textfile1.txt    |
      | /dav//files/%username%/PARENT/parent.txt |
      | //webdav/PARENT                          |
      | //dav/files/%username%//FOLDER           |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: send GET requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "GET" using password "%regular%" about user "Alice"
      | endpoint                                  |
      | //dav/spaces/%spaceid%/textfile0.txt      |
      | //dav//spaces/%spaceid%/PARENT/parent.txt |
      | /dav//spaces/%spaceid%/PARENT             |
      | //dav/spaces/%spaceid%//FOLDER            |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: send LOCK requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "LOCK" to get property "d:shared" about user "Alice"
      | endpoint                                 |
      | //webdav/textfile0.txt                   |
      | //dav//files/%username%/textfile1.txt    |
      | /dav//files/%username%/PARENT/parent.txt |
      | //webdav/PARENT                          |
      | //dav/files/%username%//FOLDER           |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: send LOCK requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "LOCK" to get property "d:shared" about user "Alice"
      | endpoint                                  |
      | //dav/spaces/%spaceid%/textfile0.txt      |
      | //dav//spaces/%spaceid%/PARENT/parent.txt |
      | /dav//spaces/%spaceid%/PARENT             |
      | //dav/spaces/%spaceid%//FOLDER            |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: send MKCOL requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "MKCOL" using password "%regular%" about user "Alice"
      | endpoint                        |
      | //webdav/PARENT1                |
      | /webdav//PARENT2                |
      | //webdav//PARENT3               |
      | //dav//files/%username%/PARENT4 |
      | /dav/files/%username%//PARENT5  |
      | /dav//files/%username%/PARENT6  |
    Then the HTTP status code of responses on each endpoint should be "200,201,200,200,201,201" on oCIS or "201,201,201,201,201,201" on reva


  Scenario: send MKCOL requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "MKCOL" using password "%regular%" about user "Alice"
      | endpoint                        |
      | //dav/spaces/%spaceid%/PARENT1  |
      | /dav//spaces/%spaceid%/PARENT2  |
      | //dav//spaces/%spaceid%/PARENT3 |
      | //dav/spaces//%spaceid%/PARENT4 |
      | /dav/spaces/%spaceid%//PARENT5  |
      | /dav//spaces/%spaceid%/PARENT6  |
    Then the HTTP status code of responses on each endpoint should be "200,201,200,200,201,201" on oCIS or "201,201,201,201,201,201" on reva


  Scenario: send MOVE requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "MOVE" using password "%regular%" about user "Alice"
      | endpoint                                  | destination                               |
      | //webdav/textfile0.txt                    | /webdav/textfileZero.txt                  |
      | /dav//files/%username%/textfile1.txt      | /dav/files/%username%/textfileOne.txt     |
      | /webdav//PARENT                           | /webdav/PARENT1                           |
      | //dav/files//%username%//PARENT1          | /dav/files/%username%/PARENT2             |
      | /dav//files/%username%/PARENT2/parent.txt | /dav/files/%username%/PARENT2/parent1.txt |
    Then the HTTP status code of responses on each endpoint should be "200,201,201,200,404" on oCIS or "201,201,201,201,201" on reva


  Scenario: send MOVE requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "MOVE" using password "%regular%" about user "Alice"
      | endpoint                                  | destination                               |
      | /dav//spaces/%spaceid%/textfile1.txt      | /dav/spaces/%spaceid%/textfileOne.txt     |
      | /dav/spaces/%spaceid%//PARENT             | /dav/spaces/%spaceid%/PARENT1             |
      | //dav/spaces/%spaceid%//PARENT1           | /dav/spaces/%spaceid%/PARENT2             |
      | //dav/spaces/%spaceid%/PARENT2/parent.txt | /dav/spaces/%spaceid%/PARENT2/parent1.txt |
    Then the HTTP status code of responses on each endpoint should be "201,201,200,200" on oCIS or "201,201,201,201" on reva


  Scenario: send POST requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "POST" including body "doesnotmatter" using password "%regular%" about user "Alice"
      | endpoint                                 |
      | //webdav/textfile0.txt                   |
      | //dav//files/%username%/textfile1.txt    |
      | /dav//files/%username%/PARENT/parent.txt |
      | /webdav//PARENT                          |
      | //dav/files//%username%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,200,412,412,200" respectively


  Scenario: send POST requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "POST" including body "doesnotmatter" using password "%regular%" about user "Alice"
      | endpoint                                 |
      | //dav//spaces/%spaceid%/textfile1.txt    |
      | /dav//spaces/%spaceid%/PARENT/parent.txt |
      | /dav//spaces/%spaceid%/PARENT            |
      | //dav//spaces/%spaceid%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,412,412,200" respectively


  Scenario: send PROPFIND requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "PROPFIND" to get property "d:href" about user "Alice"
      | endpoint                                 |
      | //webdav/textfile0.txt                   |
      | //dav//files/%username%/textfile1.txt    |
      | /dav//files/%username%/PARENT/parent.txt |
      | /webdav//PARENT                          |
      | //dav/files//%username%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,200,207,207,200" on oCIS or "207,207,207,207,207" on reva


  Scenario: send PROPFIND requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "PROPFIND" to get property "d:href" about user "Alice"
      | endpoint                                 |
      | //dav//spaces/%spaceid%/textfile1.txt    |
      | /dav//spaces/%spaceid%/PARENT/parent.txt |
      | /dav//spaces/%spaceid%/PARENT            |
      | //dav/spaces//%spaceid%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,207,207,200" on oCIS or "207,207,207,207" on reva


  Scenario: send PROPPATCH requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "PROPPATCH" to set property "d:getlastmodified" about user "Alice"
      | endpoint                                 |
      | //webdav/textfile0.txt                   |
      | //dav//files/%username%/textfile1.txt    |
      | /dav//files/%username%/PARENT/parent.txt |
      | /webdav//PARENT                          |
      | //dav//files/%username%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,200,400,400,200" respectively


  Scenario: send PROPPATCH requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "PROPPATCH" to set property "d:getlastmodified" about user "Alice"
      | endpoint                                 |
      | //dav//spaces/%spaceid%/textfile1.txt    |
      | /dav//spaces/%spaceid%/PARENT/parent.txt |
      | /dav//spaces/%spaceid%/PARENT            |
      | //dav/spaces//%spaceid%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,400,400,200" respectively


  Scenario: send PUT requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "PUT" including body "doesnotmatter" using password "%regular%" about user "Alice"
      | endpoint                                   |
      | //webdav/textfile0.txt                     |
      | /webdav//textfile1.txt                     |
      | //dav//files/%username%/textfile1.txt      |
      | /dav/files//%username%/textfile7.txt       |
      | //dav//files/%username%/PARENT//parent.txt |
    Then the HTTP status code of responses on each endpoint should be "200,204,200,201,200" on oCIS or "204,204,204,201,204" on reva


  Scenario: send PUT requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "PUT" including body "doesnotmatter" using password "%regular%" about user "Alice"
      | endpoint                                   |
      | //dav/spaces/%spaceid%/textfile0.txt       |
      | /dav//spaces/%spaceid%/textfile1.txt       |
      | //dav//spaces/%spaceid%/textfile1.txt      |
      | /dav/spaces//%spaceid%/textfile7.txt       |
      | //dav/spaces//%spaceid%/PARENT//parent.txt |
    Then the HTTP status code of responses on each endpoint should be "200,204,200,201,200" on oCIS or "204,204,204,201,204" on reva
