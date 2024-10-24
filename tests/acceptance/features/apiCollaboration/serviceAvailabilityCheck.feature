Feature: service health check

  @skip
  Scenario: check service health
    When a user requests these endpoints:
      | endpoint                     | service       | comment                                             |
      | wopi-onlyoffice:9304/healthz | collaboration | wopi-onlyoffice port 9304: Connection refused in CI |
      | wopi-collabora:9304/healthz  | collaboration | wopi-collabora port 9304: Connection refused in CI  |
      | wopi-fakeoffice:9304/healthz | collaboration | wopi-fakeoffice port 9304: Connection refused in CI |
    Then the HTTP status code of responses on all endpoints should be "200"

  @skip
  Scenario: check service readiness
    When a user requests these endpoints:
      | endpoint                    | service       | comment                                             |
      | wopi-onlyoffice:9304/readyz | collaboration | wopi-onlyoffice port 9304: Connection refused in CI |
      | wopi-collabora:9304/readyz  | collaboration | wopi-collabora port 9304: Connection refused in CI  |
      | wopi-fakeoffice:9304/readyz | collaboration | wopi-fakeoffice port 9304: Connection refused in CI |
    Then the HTTP status code of responses on all endpoints should be "200"