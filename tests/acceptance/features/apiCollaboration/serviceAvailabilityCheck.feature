Feature: service health check


  Scenario: health check defauts service
    When a user requests these endpoints
      | endpoint                      | service       |
      | http://localhost:9980/healthz | collaboration |


  Scenario: service ready check
    When a user requests these endpoints
      | endpoint                     | service       |
      | http://localhost:9980/readyz | collaboration |
