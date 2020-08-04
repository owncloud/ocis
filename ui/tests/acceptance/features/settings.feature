Feature: Set user specific settings
	As a user
	I want to set user specific settings
	So that I can customize my OCIS experience to my liking

  Background:
    Given these users have been created with default attributes:
      | username |
      | user1    |
      | user2    |

	Scenario: Check the default settings
		Given user "user1" has logged in using the webUI
		And the user browses to the settings page
		Then the setting "Language" should have value "Please select"
		When the user browses to the files page
		Then the files menu should be listed in language "English"

	Scenario: changing the language
		Given user "user1" has logged in using the webUI
		And the user browses to the settings page
		When the user changes the language to "Deutsch"
		Then the setting "Language" should have value "Deutsch"
		When the user browses to the files page
		And the user reloads the current page of the webUI
		Then the files menu should be listed in language "Deutsch"

	Scenario: changing the language only affects one user
		Given user "user2" has logged in using the webUI
		And the user browses to the settings page
		When the user changes the language to "Español"
		Then the setting "Language" should have value "Español"
		When the user browses to the files page
		And the user reloads the current page of the webUI
		Then the files menu should be listed in language "Español"
		When the user re-logs in as "user1" using the webUI
		And the user reloads the current page of the webUI
		Then the files menu should be listed in language "English"

	Scenario: Check the accounts menu when the language is changed
		Given user "user2" has logged in using the webUI
		And the user browses to the settings page
		When the user changes the language to "Deutsch"
		And the user reloads the current page of the webUI
		Then the setting "Language" should have value "Deutsch"
		And the account menu should be listed in language "Deutsch"
		When the user changes the language to "Français"
		Then the account menu should be listed in language "Français"

	Scenario: Check the files table header menu when the language is changed
		Given user "user2" has logged in using the webUI
		And the user browses to the settings page
		When the user changes the language to "Deutsch"
		Then the setting "Language" should have value "Deutsch"
		When the user browses to the files page
		And the user reloads the current page of the webUI
		Then the files header should be displayed in language "Deutsch"
