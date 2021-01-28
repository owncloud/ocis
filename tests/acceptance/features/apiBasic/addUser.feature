@api @provisioning_api-app-required @skipOnLDAP	
Feature: add user	
  As an admin	
  I want to be able to add users	
  So that I can give people controlled individual access to resources on the ownCloud server	

  Scenario Outline: admin creates a user	
    Given using OCS API version "<ocs_api_version>"	
    And user "brand-new-user" has been deleted	
    When the administrator sends a user creation request for user "brand-new-user" password "%alt1%" using the provisioning API	
    Then the OCS status code should be "<ocs_status_code>"	
    And the HTTP status code should be "200"	
    And user "brand-new-user" should exist	
    And user "brand-new-user" should be able to access a skeleton file	
    Examples:	
      | ocs_api_version | ocs_status_code |	
      | 1               | 100             |	
      | 2               | 200             |