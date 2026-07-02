const withHttp = (url) => (/^https?:\/\//i.test(url) ? url : `https://${url}`)

export const config = {
  // environment
  assets: './tests/e2e/filesForUpload',
  tempAssetsPath: './tests/e2e/filesForUpload/temp',
  baseUrlOcis: process.env.BASE_URL_OCIS ?? 'host.docker.internal:9200',
  basicAuth: process.env.BASIC_AUTH === 'true',
  testType: process.env.TEST_TYPE ?? 'playwright',
  // admin user
  adminUsername: process.env.ADMIN_USERNAME ?? 'admin',
  adminPassword: process.env.ADMIN_PASSWORD ?? 'admin',
  // use predefined users
  // if set to true, tests will not create the users
  // all users are expected to exist beforehand
  predefinedUsers: process.env.PREDEFINED_USERS === 'true',
  // json file with predefined users
  // The json file MUST contain the list of users matching the key defined in tests/e2e/support/store/user.ts
  // useful where useres are predefined but different from the default ones
  predefinedUsersFile: process.env.PREDEFINED_USERS_FILE,
  // keycloak config
  keycloak: process.env.KEYCLOAK === 'true',
  keycloakHost: process.env.KEYCLOAK_HOST ?? 'keycloak.owncloud.test',
  keycloakRealm: process.env.KEYCLOAK_REALM ?? 'oCIS',
  keycloakAdminUser: process.env.KEYCLOAK_ADMIN_USER ?? 'admin',
  keycloakAdminPassword: process.env.KEYCLOAK_ADMIN_PASSWORD ?? 'admin',
  get keycloakUrl() {
    return withHttp(this.keycloakHost)
  },
  get keycloakLoginUrl() {
    return withHttp(this.keycloakHost + '/admin/master/console')
  },
  // ocm config
  federatedBaseUrlOcis: process.env.FEDERATED_BASE_URL_OCIS ?? 'host.docker.internal:10200',
  federatedServer: false,
  get baseUrl() {
    return withHttp(this.federatedServer ? this.federatedBaseUrlOcis : this.baseUrlOcis)
  },
  debug: process.env.DEBUG === 'true',
  logLevel: process.env.LOG_LEVEL || 'silent',
  retry: parseInt(process.env.RETRY) || 0,
  // playwright
  slowMo: parseInt(process.env.SLOW_MO) || 0,
  timeout: parseInt(process.env.TIMEOUT) || 180,
  minTimeout: parseInt(process.env.MIN_TIMEOUT) || 5,
  tokenTimeout: parseInt(process.env.TOKEN_TIMEOUT) || 40,
  headless: process.env.HEADLESS === 'true',
  acceptDownloads: process.env.DOWNLOADS !== 'false',
  browser: process.env.BROWSER ?? 'chrome',
  reportDir: process.env.REPORT_DIR || 'reports/e2e',
  get tracingReportDir() {
    return this.reportDir + '/playwright/tracing'
  },
  reportVideo: process.env.REPORT_VIDEO === 'true',
  reportHar: process.env.REPORT_HAR === 'true',
  reportTracing: process.env.REPORT_TRACING === 'true',
  failOnUncaughtConsoleError: process.env.FAIL_ON_UNCAUGHT_CONSOLE_ERR === 'true',
  skipA11y: process.env.SKIP_A11Y_TESTS === 'true',
  mfa: process.env.MFA === 'true'
}
