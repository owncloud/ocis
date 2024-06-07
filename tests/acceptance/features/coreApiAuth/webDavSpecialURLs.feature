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
      | endpoint                                            |
      | //remote.php/webdav/textfile0.txt                   |
      | //remote.php//dav/files/%username%/textfile1.txt    |
      | /remote.php//dav/files/%username%/PARENT/parent.txt |
      | /remote.php//webdav/PARENT                          |
      | //remote.php/dav//files/%username%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,200,204,204,200" on oCIS or "204,204,204,204,204" on reva

  @skipOnRevaMaster
  Scenario: send DELETE requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "DELETE" using password "%regular%" about user "Alice"
      | endpoint                                             |
      | //remote.php/dav/spaces/%spaceid%/textfile0.txt      |
      | //remote.php//dav/spaces/%spaceid%/PARENT/parent.txt |
      | /remote.php//dav/spaces/%spaceid%/PARENT             |
      | //remote.php/dav//spaces/%spaceid%//FOLDER           |
    Then the HTTP status code of responses on each endpoint should be "200,200,204,200" on oCIS or "204,204,204,204" on reva


  Scenario: send GET requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "GET" using password "%regular%" about user "Alice"
      | endpoint                                            |
      | //remote.php/webdav/textfile0.txt                   |
      | //remote.php//dav/files/%username%/textfile1.txt    |
      | /remote.php//dav/files/%username%/PARENT/parent.txt |
      | /remote.php//webdav/PARENT                          |
      | //remote.php/dav//files/%username%//FOLDER          |
    Then the HTTP status code of responses on all endpoints should be "200"

  @skipOnRevaMaster
  Scenario: send GET requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "GET" using password "%regular%" about user "Alice"
      | endpoint                                             |
      | //remote.php/dav/spaces/%spaceid%/textfile0.txt      |
      | //remote.php//dav/spaces/%spaceid%/PARENT/parent.txt |
      | /remote.php//dav/spaces/%spaceid%/PARENT             |
      | //remote.php/dav//spaces/%spaceid%//FOLDER           |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: send LOCK requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "LOCK" to get property "d:shared" about user "Alice"
      | endpoint                                            |
      | //remote.php/webdav/textfile0.txt                   |
      | //remote.php//dav/files/%username%/textfile1.txt    |
      | /remote.php//dav/files/%username%/PARENT/parent.txt |
      | /remote.php//webdav/PARENT                          |
      | //remote.php/dav//files/%username%//FOLDER          |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: send LOCK requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "LOCK" to get property "d:shared" about user "Alice"
      | endpoint                                             |
      | //remote.php/dav/spaces/%spaceid%/textfile0.txt      |
      | //remote.php//dav/spaces/%spaceid%/PARENT/parent.txt |
      | /remote.php//dav/spaces/%spaceid%/PARENT             |
      | //remote.php/dav//spaces/%spaceid%//FOLDER           |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: send MKCOL requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "MKCOL" using password "%regular%" about user "Alice"
      | endpoint                                   |
      | //remote.php/webdav/PARENT1                |
      | /remote.php//webdav/PARENT2                |
      | //remote.php//webdav/PARENT3               |
      | //remote.php/dav//files/%username%/PARENT4 |
      | /remote.php/dav/files/%username%//PARENT5  |
      | /remote.php/dav//files/%username%/PARENT6  |
    Then the HTTP status code of responses on each endpoint should be "200,201,200,200,201,201" on oCIS or "201,201,201,201,201,201" on reva

  @skipOnRevaMaster
  Scenario: send MKCOL requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "MKCOL" using password "%regular%" about user "Alice"
      | endpoint                                   |
      | //remote.php/dav/spaces/%spaceid%/PARENT1  |
      | /remote.php//dav/spaces/%spaceid%/PARENT2  |
      | //remote.php//dav/spaces/%spaceid%/PARENT3 |
      | //remote.php/dav//spaces/%spaceid%/PARENT4 |
      | /remote.php/dav/spaces/%spaceid%//PARENT5  |
      | /remote.php/dav//spaces/%spaceid%/PARENT6  |
    Then the HTTP status code of responses on each endpoint should be "200,201,200,200,201,201" on oCIS or "201,201,201,201,201,201" on reva


  Scenario: send MOVE requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "MOVE" using password "%regular%" about user "Alice"
      | endpoint                                             | destination                                          |
      | //remote.php/webdav/textfile0.txt                    | /remote.php/webdav/textfileZero.txt                  |
      | /remote.php//dav/files/%username%/textfile1.txt      | /remote.php/dav/files/%username%/textfileOne.txt     |
      | /remote.php/webdav//PARENT                           | /remote.php/webdav/PARENT1                           |
      | //remote.php/dav/files/%username%//PARENT1           | /remote.php/dav/files/%username%/PARENT2             |
      | /remote.php/dav//files/%username%/PARENT2/parent.txt | /remote.php/dav/files/%username%/PARENT2/parent1.txt |
    Then the HTTP status code of responses on each endpoint should be "200,201,201,200,404" on oCIS or "201,201,201,201,201" on reva

  @skipOnRevaMaster
  Scenario: send MOVE requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "MOVE" using password "%regular%" about user "Alice"
      | endpoint                                             | destination                                          |
      | /remote.php//dav/spaces/%spaceid%/textfile1.txt      | /remote.php/dav/spaces/%spaceid%/textfileOne.txt     |
      | /remote.php/dav//spaces/%spaceid%/PARENT             | /remote.php/dav/spaces/%spaceid%/PARENT1             |
      | //remote.php/dav/spaces/%spaceid%//PARENT1           | /remote.php/dav/spaces/%spaceid%/PARENT2             |
      | //remote.php/dav/spaces/%spaceid%/PARENT2/parent.txt | /remote.php/dav/spaces/%spaceid%/PARENT2/parent1.txt |
    Then the HTTP status code of responses on each endpoint should be "201,201,200,200" on oCIS or "201,201,201,201" on reva


  Scenario: send POST requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "POST" including body "doesnotmatter" using password "%regular%" about user "Alice"
      | endpoint                                            |
      | //remote.php/webdav/textfile0.txt                   |
      | //remote.php//dav/files/%username%/textfile1.txt    |
      | /remote.php//dav/files/%username%/PARENT/parent.txt |
      | /remote.php//webdav/PARENT                          |
      | //remote.php/dav//files/%username%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,200,412,412,200" respectively

  @skipOnRevaMaster
  Scenario: send POST requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "POST" including body "doesnotmatter" using password "%regular%" about user "Alice"
      | endpoint                                            |
      | //remote.php//dav/spaces/%spaceid%/textfile1.txt    |
      | /remote.php//dav/spaces/%spaceid%/PARENT/parent.txt |
      | /remote.php//dav/spaces/%spaceid%/PARENT            |
      | //remote.php/dav//spaces/%spaceid%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,412,412,200" respectively


  Scenario: send PROPFIND requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "PROPFIND" to get property "d:href" about user "Alice"
      | endpoint                                            |
      | //remote.php/webdav/textfile0.txt                   |
      | //remote.php//dav/files/%username%/textfile1.txt    |
      | /remote.php//dav/files/%username%/PARENT/parent.txt |
      | /remote.php//webdav/PARENT                          |
      | //remote.php/dav//files/%username%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,200,207,207,200" on oCIS or "207,207,207,207,207" on reva

  @skipOnRevaMaster
  Scenario: send PROPFIND requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "PROPFIND" to get property "d:href" about user "Alice"
      | endpoint                                            |
      | //remote.php//dav/spaces/%spaceid%/textfile1.txt    |
      | /remote.php//dav/spaces/%spaceid%/PARENT/parent.txt |
      | /remote.php//dav/spaces/%spaceid%/PARENT            |
      | //remote.php/dav//spaces/%spaceid%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,207,207,200" on oCIS or "207,207,207,207" on reva


  Scenario: send PROPPATCH requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "PROPPATCH" to set property "d:getlastmodified" about user "Alice"
      | endpoint                                            |
      | //remote.php/webdav/textfile0.txt                   |
      | //remote.php//dav/files/%username%/textfile1.txt    |
      | /remote.php//dav/files/%username%/PARENT/parent.txt |
      | /remote.php//webdav/PARENT                          |
      | //remote.php/dav//files/%username%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,200,400,400,200" respectively

  @skipOnRevaMaster
  Scenario: send PROPPATCH requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "PROPPATCH" to set property "d:getlastmodified" about user "Alice"
      | endpoint                                            |
      | //remote.php//dav/spaces/%spaceid%/textfile1.txt    |
      | /remote.php//dav/spaces/%spaceid%/PARENT/parent.txt |
      | /remote.php//dav/spaces/%spaceid%/PARENT            |
      | //remote.php/dav//spaces/%spaceid%//FOLDER          |
    Then the HTTP status code of responses on each endpoint should be "200,400,400,200" respectively


  Scenario: send PUT requests to webDav endpoints with 2 slashes
    When user "Alice" requests these endpoints with "PUT" including body "doesnotmatter" using password "%regular%" about user "Alice"
      | endpoint                                             |
      | //remote.php/webdav/textfile0.txt                    |
      | /remote.php//webdav/textfile1.txt                    |
      | //remote.php//dav/files/%username%/textfile1.txt     |
      | /remote.php/dav/files/%username%/textfile7.txt       |
      | //remote.php/dav/files/%username%/PARENT//parent.txt |
    Then the HTTP status code of responses on each endpoint should be "200,204,200,201,200" on oCIS or "204,204,204,201,204" on reva

  @skipOnRevaMaster
  Scenario: send PUT requests to webDav endpoints with 2 slashes using the spaces WebDAV API
    When user "Alice" requests these endpoints with "PUT" including body "doesnotmatter" using password "%regular%" about user "Alice"
      | endpoint                                             |
      | //remote.php/dav/spaces/%spaceid%/textfile0.txt      |
      | /remote.php//dav/spaces/%spaceid%/textfile1.txt      |
      | //remote.php//dav/spaces/%spaceid%/textfile1.txt     |
      | /remote.php/dav/spaces/%spaceid%/textfile7.txt       |
      | //remote.php/dav/spaces/%spaceid%/PARENT//parent.txt |
    Then the HTTP status code of responses on each endpoint should be "200,204,200,201,200" on oCIS or "204,204,204,201,204" on reva
