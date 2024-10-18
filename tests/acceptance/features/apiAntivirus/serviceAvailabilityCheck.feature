Feature: service health check


  Scenario: health check
    When a user requests these endpoints:
      | endpoint                                        | service   |
      | %base_url_without_scheme_and_port%:9297/healthz | antivirus |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: ready check
    When a user requests these endpoints:
      | endpoint                                       | service   |
      | %base_url_without_scheme_and_port%:9297/readyz | antivirus |
    Then the HTTP status code of responses on all endpoints should be "200"
