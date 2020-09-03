Feature: Accounts

	Scenario: admin checks accounts list
		Given user "Moss" has logged in using the webUI
		When the user browses to the accounts page
		Then user "einstein" should be displayed in the accounts list on the WebUI
		And user "konnectd" should be displayed in the accounts list on the WebUI
		And user "marie" should be displayed in the accounts list on the WebUI
		And user "reva" should be displayed in the accounts list on the WebUI
		And user "richard" should be displayed in the accounts list on the WebUI

	Scenario: admin changes non-admin user's role to admin
		Given user "Moss" has logged in using the webUI
		When the user browses to the accounts page
		Then user "einstein" should be displayed in the accounts list on the WebUI
		When the user changes the role of user "einstein" to "Admin" using the WebUI
		Then the displayed role of user "einstein" should be "Admin" on the WebUI
		When the user reloads the current page of the webUI
		Then the displayed role of user "einstein" should be "Admin" on the WebUI

	Scenario: regular user should not be able to see accounts list
		Given user "Marie" has logged in using the webUI
		When the user browses to the accounts page
		Then the user should not be able to see the accounts list on the WebUI

	Scenario: guest user should not be able to see accounts list
		Given user "Moss" has logged in using the webUI
		When the user browses to the accounts page
		Then user "einstein" should be displayed in the accounts list on the WebUI
		When the user changes the role of user "einstein" to "Guest" using the WebUI
		And the user logs out of the webUI
		And user "Einstein" logs in using the webUI
		And the user browses to the accounts page
		Then the user should not be able to see the accounts list on the WebUI
