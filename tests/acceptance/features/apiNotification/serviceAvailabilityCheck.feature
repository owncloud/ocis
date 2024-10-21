Feature: service health check


  Scenario: health check
    When a user requests these endpoints
      | endpoint                      | service      |
      | http://localhost:9174/healthz | notification |


  Scenario: ready check
    When a user requests these endpoints
      | endpoint                     | service      |
      | http://localhost:9174/readyz | notification |
