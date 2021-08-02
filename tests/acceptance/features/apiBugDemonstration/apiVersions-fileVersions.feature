@api @files_versions-app-required @skipOnOcis-EOS-Storage @issue-ocis-reva-275

Feature: dav-versions

  Background:
    Given using OCS API version "2"
    And using new DAV path
    And user "Alice" has been created with default attributes and without skeleton files


