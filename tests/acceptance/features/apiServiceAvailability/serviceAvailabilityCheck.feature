@env-config
Feature: service health check


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/healthz | audit   |
      | http://%base_url_hostname%:9229/readyz  | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check services health and readiness while running separately
    Given the following configs have been set:
      | config                    | value |
      | OCIS_EXCLUDE_RUN_SERVICES | audit |
      | OCIS_LOG_LEVEL            | info  |
    And the administrator has started service "audit" separately with the following configs:
      | config           | value        |
      | AUDIT_LOG_LEVEL  | info         |
      | AUDIT_DEBUG_ADDR | 0.0.0.0:9229 |
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service |
      | http://%base_url_hostname%:9229/readyz  | audit   |
      | http://%base_url_hostname%:9229/healthz | audit   |
    Then the HTTP status code of responses on all endpoints should be "200"
