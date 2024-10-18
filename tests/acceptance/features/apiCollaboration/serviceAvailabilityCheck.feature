Feature: service health check


  Scenario: health check defauts service
    When a user requests these endpoints with "GET"
      | endpoint                                         | service       | comment |
      | %base_url_wÏithout_scheme_and_port%:9980/healthz | collaboration | extra   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: service ready check
    When a user requests these endpoints with "GET"
      | %base_url_wÏithout_scheme_and_port%:9980/readyz | collaboration |
    Then the HTTP status code of responses on all endpoints should be "200"