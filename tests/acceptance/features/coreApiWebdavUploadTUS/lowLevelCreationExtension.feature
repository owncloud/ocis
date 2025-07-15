Feature: low level tests of the creation extension see https://tus.io/protocols/resumable-upload.html#creation
  As a user
  I want to able to create resources
  So that I can manage my resources

  Background:
    Given user "Alice" has been created with default attributes


  Scenario Outline: creating a new upload resource
    Given using <dav-path-version> DAV path
    When user "Alice" creates a new TUS resource on the WebDAV API with these headers:
      | Upload-Length   | 100                                           |
      #    d29ybGRfZG9taW5hdGlvbl9wbGFuLnBkZg== is the base64 encode of world_domination_plan.pdf
      | Upload-Metadata | filename d29ybGRfZG9taW5hdGlvbl9wbGFuLnBkZg== |
      | Tus-Resumable   | 1.0.0                                         |
    Then the HTTP status code should be "201"
    And the following headers should match these regular expressions
      | Tus-Resumable | /1\.0\.0/                       |
      | Location      | /http[s]?:\/\/.*:\d+\/data\/.*/ |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: creating a new upload resource without upload length
    Given using <dav-path-version> DAV path
    When user "Alice" creates a new TUS resource on the WebDAV API with these headers:
      | Tus-Resumable   | 1.0.0                                         |
      #    d29ybGRfZG9taW5hdGlvbl9wbGFuLnBkZg== is the base64 encode of world_domination_plan.pdf
      | Upload-Metadata | filename d29ybGRfZG9taW5hdGlvbl9wbGFuLnBkZg== |
    Then the HTTP status code should be "412"
    And the following headers should not be set
      | header   |
      | Location |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |
