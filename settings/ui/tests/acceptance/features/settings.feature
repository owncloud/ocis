Feature: Set user specific settings
	As a user
	I want to set user specific settings
	So that I can customize my OCIS experience to my liking

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | user1    |
      | user2    |
	And user "user1" has created folder "simple-folder"

	Scenario: Check the default settings
		Given user "user1" has logged in using the webUI
		And the user browses to the settings page
		Then the setting "Language" should not have any value
		When the user browses to the files page
		Then the files menu should be listed in language "English"

	Scenario: changing the language (reactive and with page reload)
		Given user "user1" has logged in using the webUI
		And the user browses to the settings page
		When the user changes the language to "Deutsch"
		Then the setting "Language" should have value "Deutsch"
		When the user browses to the files page
		Then the files menu should be listed in language "Deutsch"
		And the account menu should be listed in language "Deutsch"
		And the files header should be displayed in language "Deutsch"
		When the user reloads the current page of the webUI
		Then the files menu should be listed in language "Deutsch"
		And the account menu should be listed in language "Deutsch"
		And the files header should be displayed in language "Deutsch"
		When the user browses to the settings page
		And the user changes the language to "English"
		And the user browses to the files page
		When the user browses to the files page
		Then the files menu should be listed in language "English"

	Scenario: changing the language only affects one user
		Given user "user2" has logged in using the webUI
		And the user browses to the settings page
		When the user changes the language to "Español"
		Then the setting "Language" should have value "Español"
		When the user browses to the files page
		Then the files menu should be listed in language "Español"
		When the user re-logs in as "user1" using the webUI
		Then the files menu should be listed in language "English"
