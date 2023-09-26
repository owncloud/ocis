@env-config
Feature: enforce password on public link
  As a user
  I want to enforce passwords on public links shared with upload, edit, or contribute permission
  So that the password is required to access the contents of the link

  Scenario Outline: create a public link with edit permission without a password when enforce-password is enabled
    Given the config "OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD" has been set to "true"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And using OCS API version "<ocs-api-version>"
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | /testfile.txt |
      | permissions | 3             |
    Then the HTTP status code should be "<http-code>"
    And the OCS status code should be "400"
    And the OCS status message should be "missing required password"
    Examples:
      | ocs-api-version | http-code |
      | 1               | 200       |
      | 2               | 400       |


  Scenario Outline: update a public link to edit permission without a password
    Given the config "OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD" has been set to "true"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And using OCS API version "<ocs-api-version>"
    And user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3 |
    Then the HTTP status code should be "<http-code>"
    And the OCS status code should be "400"
    And the OCS status message should be "missing required password"
    Examples:
      | ocs-api-version | http-code |
      | 1               | 200       |
      | 2               | 400       |


  Scenario Outline: updates a public link to edit permission with a password
    Given the config "OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD" has been set to "true"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And using OCS API version "<ocs-api-version>"
    And user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3            |
      | password    | testpassword |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs-code>"
    And the OCS status message should be "OK"
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API without a password
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "wrong pass"
    But the public should be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "testpassword"
    Examples:
      | ocs-api-version | ocs-code |
      | 1               | 100      |
      | 2               | 200      |


  Scenario Outline: create a public link with a password in accordance with the password policy
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
      | FRONTEND_PASSWORD_POLICY_MIN_CHARACTERS                | 13    |
      | FRONTEND_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS      | 3     |
      | FRONTEND_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS      | 2     |
      | FRONTEND_PASSWORD_POLICY_MIN_DIGITS                    | 2     |
      | FRONTEND_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS        | 2     |
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And using OCS API version "<ocs-api-version>"
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | /testfile.txt |
      | permissions | 3             |
      | password    | 3s:5WW9uE5h=A |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs-code>"
    And the OCS status message should be "OK"
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API without a password
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "wrong pass"
    But the public should be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "3s:5WW9uE5h=A"
    Examples:
      | ocs-api-version | ocs-code |
      | 1               | 100      |
      | 2               | 200      |


  Scenario Outline: try to create a public link with a password that does not comply with the password policy
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
      | FRONTEND_PASSWORD_POLICY_MIN_CHARACTERS                | 13    |
      | FRONTEND_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS      | 3     |
      | FRONTEND_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS      | 2     |
      | FRONTEND_PASSWORD_POLICY_MIN_DIGITS                    | 2     |
      | FRONTEND_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS        | 2     |
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And using OCS API version "<ocs-api-version>"
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | /testfile.txt |
      | permissions | 3             |
      | password    | Pas1          |
    Then the HTTP status code should be "<http-code>"
    And the OCS status code should be "400"
    And the OCS status message should be:
      """
      At least 13 characters are required
      At least 3 lowercase letters are required
      At least 2 uppercase letters are required
      At least 2 numbers are required
      At least 2 special characters are required  !"#$%&'()*+,-./:;<=>?@[\]^_`{|}~
      """
    Examples:
      | ocs-api-version | http-code |
      | 1               | 200       |
      | 2               | 400       |


  Scenario Outline: update a public link with a password in accordance with the password policy
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
      | FRONTEND_PASSWORD_POLICY_MIN_CHARACTERS                | 13    |
      | FRONTEND_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS      | 3     |
      | FRONTEND_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS      | 2     |
      | FRONTEND_PASSWORD_POLICY_MIN_DIGITS                    | 1     |
      | FRONTEND_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS        | 2     |
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And using OCS API version "<ocs-api-version>"
    And user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3             |
      | password    | 6a0Q;A3 +i^m[ |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs-code>"
    And the OCS status message should be "OK"
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API without a password
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "wrong pass"
    But the public should be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "6a0Q;A3 +i^m["
    Examples:
      | ocs-api-version | ocs-code |
      | 1               | 100      |
      | 2               | 200      |


  Scenario Outline: try to update a public link with a password that does not comply with the password policy
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
      | FRONTEND_PASSWORD_POLICY_MIN_CHARACTERS                | 13    |
      | FRONTEND_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS      | 3     |
      | FRONTEND_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS      | 2     |
      | FRONTEND_PASSWORD_POLICY_MIN_DIGITS                    | 1     |
      | FRONTEND_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS        | 2     |
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And using OCS API version "<ocs-api-version>"
    And user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3    |
      | password    | Pws^ |
    Then the HTTP status code should be "<http-code>"
    And the OCS status code should be "400"
    And the OCS status message should be:
      """
      At least 13 characters are required
      At least 3 lowercase letters are required
      At least 2 uppercase letters are required
      At least 1 numbers are required
      At least 2 special characters are required  !"#$%&'()*+,-./:;<=>?@[\]^_`{|}~
      """
    Examples:
      | ocs-api-version | http-code |
      | 1               | 200       |
      | 2               | 400       |


  Scenario Outline: update a public link with a password in accordance with the password policy (valid cases)
    Given the following configs have been set:
      | config                                                 | value          |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true           |
      | <config>                                               | <config-value> |
    And using OCS API version "2"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3          |
      | password    | <password> |
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API without a password
    And the public should not be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "wrong pass"
    But the public should be able to download file "/textfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "<password>"
    Examples:
      | config                                            | config-value | password                          |
      | FRONTEND_PASSWORD_POLICY_MIN_CHARACTERS           | 5            | 12345                             |
      | FRONTEND_PASSWORD_POLICY_MIN_CHARACTERS           | 10           | with space                        |
      | FRONTEND_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS | 3            | Test                              |
      | FRONTEND_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS | 3            | TeST                              |
      | FRONTEND_PASSWORD_POLICY_MIN_DIGITS               | 1            | test1                             |
      | FRONTEND_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 1            | test pass                         |
      | FRONTEND_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 33           | ! #$%&'()*+,-./:;<=>?@[\]^_`{  }~ |
      | FRONTEND_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 5            | sameCharacterShouldWork!!!!!      |


  Scenario Outline: try to update a public link with a password that does not comply with the password policy (invalid cases)
    Given the following configs have been set:
      | config                                                 | value          |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true           |
      | <config>                                               | <config-value> |
    And using OCS API version "2"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3          |
      | password    | <password> |
    Then the HTTP status code should be "400"
    And the OCS status code should be "400"
    And the OCS status message should be "<message>"
    Examples:
      | config                                            | config-value | password | message                                   |
      | FRONTEND_PASSWORD_POLICY_MIN_CHARACTERS           | 5            | 1234     | at least 5 characters are required        |
      | FRONTEND_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS | 3            | TesT     | at least 3 lowercase letters are required |
      | FRONTEND_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS | 3            | TesT     | at least 3 uppercase letters are required |
      | FRONTEND_PASSWORD_POLICY_MIN_DIGITS               | 2            | test1    | at least 2 numbers are required           |


  Scenario Outline: update a public link with a password that is listed in the Banned-Password-List
    Given the config "FRONTEND_PASSWORD_POLICY_BANNED_PASSWORDS_LIST" has been set to path "bannedPassword/banned-password-list.txt"
    And using OCS API version "2"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And user "Alice" has created a public link share with settings
      | path        | /testfile.txt |
      | permissions | 1             |
    When user "Alice" updates the last public link share using the sharing API with
      | permissions | 3          |
      | password    | <password> |
    Then the HTTP status code should be "<http-code>"
    And the OCS status code should be "<ocs-code>"
    And the OCS status message should be "<message>"
    Examples:
      | password | http-code | ocs-code | message                                                                                               |
      | 123      | 400       | 400      | unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety |
      | password | 400       | 400      | unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety |
      | ownCloud | 400       | 400      | unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety |
      | 1234     | 200       | 200      | OK                                                                                                    |
      | OwnCloud | 200       | 200      | OK                                                                                                    |
      | p@ssword | 200       | 200      | OK                                                                                                    |


  Scenario Outline: create  a public link with a password that is listed in the Banned-Password-List
    Given the config "FRONTEND_PASSWORD_POLICY_BANNED_PASSWORDS_LIST" has been set to path "bannedPassword/banned-password-list.txt"
    And using OCS API version "2"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | /testfile.txt |
      | permissions | 3             |
      | password    | <password>    |
    Then the HTTP status code should be "<http-code>"
    And the OCS status code should be "<ocs-code>"
    And the OCS status message should be "<message>"
    Examples:
      | password | http-code | ocs-code | message                                                                                               |
      | 123      | 400       | 400      | unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety |
      | password | 400       | 400      | unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety |
      | ownCloud | 400       | 400      | unfortunately, your password is commonly used. please pick a harder-to-guess password for your safety |
      | 1234     | 200       | 200      | OK                                                                                                    |
      | OwnCloud | 200       | 200      | OK                                                                                                    |
      | p@ssword | 200       | 200      | OK                                                                                                    |
