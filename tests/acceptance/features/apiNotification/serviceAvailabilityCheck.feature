Feature: service health check


  Scenario: health check
    When a user requests these endpoints with "GET"
      | endpoint                      | service      |
      | http://127.0.0.1:9174/healthz | notification |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: ready check
    When a user requests these endpoints with "GET"
      | endpoint                     | service      |
      | http://127.0.0.1:9174/readyz | notification |
    Then the HTTP status code of responses on all endpoints should be "200"