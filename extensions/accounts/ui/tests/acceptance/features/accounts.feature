Feature: Accounts

	Scenario: admin checks accounts list
		Given user "Moss" has logged in using the webUI
		When the user browses to the accounts page
		Then user "einstein" should be displayed in the accounts list on the WebUI
		And user "idp" should be displayed in the accounts list on the WebUI
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

	@skip @issue-product-167
	Scenario: regular user should not be able to see accounts list
		Given user "Marie" has logged in using the webUI
		When the user browses to the accounts page
		Then the user should not be able to see the accounts list on the WebUI

	@skip @issue-product-167
	Scenario: guest user should not be able to see accounts list
		Given user "Moss" has logged in using the webUI
		When the user browses to the accounts page
		Then user "einstein" should be displayed in the accounts list on the WebUI
		When the user changes the role of user "einstein" to "Guest" using the WebUI
		And the user logs out of the webUI
		And user "Einstein" logs in using the webUI
		And the user browses to the accounts page
		Then the user should not be able to see the accounts list on the WebUI

	# We want to separate this into own scenarios but because we do not have clean env for each scenario yet
	# we are resetting it manually by combining them into one
	Scenario: disable/enable account
		Given user "Moss" has logged in using the webUI
		When the user browses to the accounts page
		Then user "einstein" should be displayed in the accounts list on the WebUI
		When the user disables user "einstein" using the WebUI
		Then the status indicator of user "einstein" should be "disabled" on the WebUI
		# And user "einstein" should not be able to log in
		When the user enables user "einstein" using the WebUI
		Then the status indicator of user "einstein" should be "enabled" on the WebUI
		# And user "einstein" should be able to log in

	Scenario: disable/enable multiple accounts
		Given user "Moss" has logged in using the webUI
		When the user browses to the accounts page
		Then user "einstein" should be displayed in the accounts list on the WebUI
		And user "marie" should be displayed in the accounts list on the WebUI
		When the user disables users "einstein,marie" using the WebUI
		Then the status indicator of users "einstein,marie" should be "disabled" on the WebUI
		# And user "einstein" should not be able to log in
		# And user "marie" should not be able to log in
		When the user enables users "einstein,marie" using the WebUI
		Then the status indicator of user "einstein,marie" should be "enabled" on the WebUI
		# And user "einstein" should be able to log in
		# And user "marie" should be able to log in

	Scenario: create a user
		Given user "Moss" has logged in using the webUI
		And the user browses to the accounts page
		When the user creates a new user with username "bob", email "bob@example.org" and password "bob" using the WebUI
		Then user "bob" should be displayed in the accounts list on the WebUI

	Scenario: delete a user
		Given user "Moss" has logged in using the webUI
		And the user browses to the accounts page
		When the user deletes user "bob" using the WebUI
		Then user "bob" should not be displayed in the accounts list on the WebUI
