config = {
  'modules': {
    'accounts': 'frontend',
    'glauth':'',
    'konnectd':'',
    'ocis-phoenix':'',
    'ocis-pkg':'',
    'storage':'',
    'ocs':'',
    'proxy':'',
    'settings':'frontend',
    'store':'',
    'thumbnails':'',
    'webdav':'',
  },
  'apiTests': {
    'coreBranch': 'master',
    'coreCommit': '79fea80eece798d9cec7dadcc0cf02cf741cfbe6',
    'numberOfParts': 6
  },
  'uiTests': {
    'phoenixBranch': 'master',
    'phoenixCommit': '693fa6d275243fc860b4b82e5146d25f62d927fb',
    'suites': {
      'phoenixWebUI1': [
        'webUICreateFilesFolders',
        'webUIDeleteFilesFolders',
        'webUIFavorites',
      ],
      'phoenixWebUI2' : [
        'webUIFiles',
        'webUILogin',
        'webUINotifications',
        'webUIPrivateLinks',
      ],
      'phoenixWebUI3': [
        'webUIAccount',
        'webUIRenameFiles',
        'webUIRenameFolders',
      ],
       'phoenixWebUI4': [
        'webUITrashbin',
        'webUIUpload',
        # All tests in the following suites are skipped currently
        # so they won't run now but when they are enabled they will run
        'webUIRestrictSharing',
        'webUISharingAutocompletion',
        'webUISharingInternalGroups',
        'webUISharingInternalUsers',
        'webUISharingPermissionsUsers',
        'webUISharingFilePermissionsGroups',
        'webUISharingFolderPermissionsGroups',
        'webUISharingFolderAdvancedPermissionsGroups',
        'webUIResharing',
        'webUISharingPublic',
        'webUISharingPublicDifferentRoles',
        'webUISharingAcceptShares',
        'webUISharingFilePermissionMultipleUsers',
        'webUISharingFolderPermissionMultipleUsers',
        'webUISharingFolderAdvancedPermissionMultipleUsers',
        'webUISharingNotifications',
      ],
    }
  }
}
def getTestSuiteNames():
  keys = config['modules'].keys()
  names = []
  for key in keys:
    names.append('linting&unitTests-%s' % (key))
  return names

def getUITestSuiteNames():
  return config['uiTests']['suites'].keys()

def getUITestSuites():
  return config['uiTests']['suites']

def getCoreApiTestPipelineNames():
  names = []
  for runPart in range(1, config['apiTests']['numberOfParts'] + 1):
    names.append('Core-API-Tests-owncloud-storage-%s' % runPart)
    names.append('Core-API-Tests-ocis-storage-%s' % runPart)
  return names

def getDependsOnAllTestPipelines(ctx):
  dependencies = getTestSuiteNames() + [ 'upload-coverage' ]
  if ctx.build.ref != "refs/heads/master":
    dependencies = getTestSuiteNames() + [
      'upload-coverage',
      'localApiTests-owncloud-storage',
      'localApiTests-ocis-storage',
    ] + getCoreApiTestPipelineNames() + getUITestSuiteNames() + ['accountsUITests']

  return dependencies

def main(ctx):
  pipelines = []
  before = testPipelines(ctx)

  stages = [
    docker(ctx, 'amd64'),
    docker(ctx, 'arm64'),
    docker(ctx, 'arm'),
    dockerEos(ctx),
    binary(ctx, 'linux'),
    binary(ctx, 'darwin'),
    binary(ctx, 'windows'),
    releaseSubmodule(ctx),
  ]

  after = [
    manifest(ctx),
    changelog(ctx),
    readme(ctx),
    badges(ctx),
    docs(ctx),
    updateDeployment(ctx)
  ]

  if '[docs-only]' in (ctx.build.title + ctx.build.message):
    pipelines = docs(ctx)
    pipelines['depends_on'] = []
  else:
    pipelines = before + stages + after

  return pipelines

def testPipelines(ctx):
  pipelines = []

  for module in config['modules']:
    pipelines.append(testing(ctx, module))

  pipelines += [
    uploadCoverage(ctx),
    localApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], 'owncloud'),
    localApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], 'ocis')
  ]

  for runPart in range(1, config['apiTests']['numberOfParts'] + 1):
    pipelines.append(coreApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], runPart, config['apiTests']['numberOfParts'], 'owncloud'))
    pipelines.append(coreApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], runPart, config['apiTests']['numberOfParts'], 'ocis'))

  pipelines += uiTests(ctx, config['uiTests']['phoenixBranch'], config['uiTests']['phoenixCommit'])
  pipelines.append(accountsUITests(ctx, config['uiTests']['phoenixBranch'], config['uiTests']['phoenixCommit'], 'owncloud'))
  return pipelines

def testing(ctx, module):
  steps = generate(module) + [
    {
      'name': 'vet',
      'image': 'webhippie/golang:1.13',
      'pull': 'always',
      'commands': [
        'cd %s' % (module),
        'make vet',
      ],
      'volumes': [
        {
          'name': 'gopath',
          'path': '/srv/app',
        },
      ],
    },
    {
      'name': 'staticcheck',
      'image': 'webhippie/golang:1.13',
      'pull': 'always',
      'commands': [
        'cd %s' % (module),
        'make staticcheck',
      ],
      'volumes': [
        {
          'name': 'gopath',
          'path': '/srv/app',
        },
      ],
    },
    {
      'name': 'lint',
      'image': 'webhippie/golang:1.13',
      'pull': 'always',
      'commands': [
        'cd %s' % (module),
        'make lint',
      ],
      'volumes': [
        {
          'name': 'gopath',
          'path': '/srv/app',
        },
      ],
    },
    {
        'name': 'test',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'commands': [
          'cd %s' % (module),
          'make test',
          'mv coverage.out %s_coverage.out' % (module),
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app',
          },
        ],
      },
      {
        'name': 'coverage-cache',
        'image': 'plugins/s3',
        'settings': {
          'endpoint': {
            'from_secret': 'cache_s3_endpoint'
          },
          'bucket': 'cache',
          'source': '%s/%s_coverage.out' % (module, module),
          'target': '%s/%s/coverage' % (ctx.repo.slug, ctx.build.commit + '-${DRONE_BUILD_NUMBER}'),
          'path_style': True,
          'strip_prefix': module,
          'access_key': {
            'from_secret': 'cache_s3_access_key'
          },
          'secret_key': {
            'from_secret': 'cache_s3_secret_key'
          }
        }
      }
  ]

  if config['modules'][module] == 'frontend':
    steps = frontend(module) + steps

  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'linting&unitTests-%s' % (module),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': steps,
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/tags/%s/v*' % (module),
        'refs/pull/**',
      ],
    },
  }

def uploadCoverage(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'upload-coverage',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'sync-from-cache',
        'image': 'minio/mc',
        'environment': {
          'MC_HOST_cache': {
            'from_secret': 'cache_s3_connection_url'
          }
        },
        'commands': [
          'mkdir -p coverage',
          'mc mirror cache/cache/%s/%s/coverage coverage/' % (ctx.repo.slug, ctx.build.commit + '-${DRONE_BUILD_NUMBER}'),
        ]
      },
      {
        'name': 'codacy',
        'image': 'plugins/codacy:1',
        'pull': 'always',
        'settings': {
          'token': {
            'from_secret': 'codacy_token',
          },
        },
      },
      {
        'name': 'sonarcloud',
        'image': 'sonarsource/sonar-scanner-cli',
        'pull': 'always',
        'environment': {
          'SONAR_TOKEN': {
            'from_secret': 'sonar_token',
          },
          'SONAR_PULL_REQUEST_BASE': '%s' % ('master' if ctx.build.event == 'pull_request' else None),
          'SONAR_PULL_REQUEST_BRANCH': '%s' % (ctx.build.source if ctx.build.event == 'pull_request' else None),
          'SONAR_PULL_REQUEST_KEY': '%s' % (ctx.build.ref.replace("refs/pull/", "").split("/")[0] if ctx.build.event == 'pull_request' else None),
        },
      },
      {
        'name': 'purge-cache',
        'image': 'minio/mc',
        'environment': {
          'MC_HOST_cache': {
            'from_secret': 'cache_s3_connection_url'
          }
        },
        'commands': [
          'mc rm --recursive --force cache/cache/%s/%s/coverage' % (ctx.repo.slug, ctx.build.commit + '-${DRONE_BUILD_NUMBER}'),
        ]
      },
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
    'depends_on': getTestSuiteNames(),
  }

def localApiTests(ctx, coreBranch = 'master', coreCommit = '', storage = 'owncloud'):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'localApiTests-%s-storage' % (storage),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate('ocis') +
      build() +
      ocisServer(storage) +
      cloneCoreRepos(coreBranch, coreCommit) + [
      {
        'name': 'localApiTests-%s-storage' % (storage),
        'image': 'owncloudci/php:7.2',
        'pull': 'always',
        'environment' : {
          'TEST_SERVER_URL': 'https://ocis-server:9200',
          'OCIS_REVA_DATA_ROOT': '%s' % ('/srv/app/tmp/ocis/owncloud/' if storage == 'owncloud' else ''),
          'DELETE_USER_DATA_CMD': '%s' % ('' if storage == 'owncloud' else 'rm -rf /srv/app/tmp/ocis/storage/users/nodes/root/*'),
          'SKELETON_DIR': '/srv/app/tmp/testing/data/apiSkeleton',
          'OCIS_SKELETON_STRATEGY': '%s' % ('copy' if storage == 'owncloud' else 'upload'),
          'TEST_OCIS':'true',
          'BEHAT_FILTER_TAGS': '~@skipOnOcis-%s-Storage' % ('OC' if storage == 'owncloud' else 'OCIS'),
          'PATH_TO_CORE': '/srv/app/testrunner'
        },
        'commands': [
          'cd ocis',
          'make test-acceptance-api',
        ],
        'volumes': [{
          'name': 'gopath',
          'path': '/srv/app',
        }]
      },
    ],
    'services':
      redis(),
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'trigger': {
      'ref': [
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def coreApiTests(ctx, coreBranch = 'master', coreCommit = '', part_number = 1, number_of_parts = 1, storage = 'owncloud'):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'Core-API-Tests-%s-storage-%s' % (storage, part_number),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate('ocis') +
      build() +
      ocisServer(storage) +
      cloneCoreRepos(coreBranch, coreCommit) + [
      {
        'name': 'oC10ApiTests-%s-storage-%s' % (storage, part_number),
        'image': 'owncloudci/php:7.2',
        'pull': 'always',
        'environment' : {
          'TEST_SERVER_URL': 'https://ocis-server:9200',
          'OCIS_REVA_DATA_ROOT': '%s' % ('/srv/app/tmp/ocis/owncloud/' if storage == 'owncloud' else ''),
          'DELETE_USER_DATA_CMD': '%s' % ('' if storage == 'owncloud' else 'rm -rf /srv/app/tmp/ocis/storage/users/nodes/root/*'),
          'SKELETON_DIR': '/srv/app/tmp/testing/data/apiSkeleton',
          'OCIS_SKELETON_STRATEGY': '%s' % ('copy' if storage == 'owncloud' else 'upload'),
          'TEST_OCIS':'true',
          'BEHAT_FILTER_TAGS': '~@notToImplementOnOCIS&&~@toImplementOnOCIS&&~comments-app-required&&~@federation-app-required&&~@notifications-app-required&&~systemtags-app-required&&~@local_storage&&~@skipOnOcis-%s-Storage' % ('OC' if storage == 'owncloud' else 'OCIS'),
          'DIVIDE_INTO_NUM_PARTS': number_of_parts,
          'RUN_PART': part_number,
          'EXPECTED_FAILURES_FILE': '/drone/src/ocis/tests/acceptance/expected-failures-on-%s-storage.txt' % (storage.upper())
        },
        'commands': [
          'cd /srv/app/testrunner',
          'make test-acceptance-api',
        ],
        'volumes': [{
          'name': 'gopath',
          'path': '/srv/app',
        }]
      },
    ],
    'services':
      redis(),
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'trigger': {
      'ref': [
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def uiTests(ctx, phoenixBranch, phoenixCommit):
  suiteNames = getUITestSuiteNames()
  return [uiTestPipeline(suiteName, phoenixBranch, phoenixCommit) for suiteName in suiteNames]

def uiTestPipeline(suiteName, phoenixBranch = 'master', phoenixCommit = '', storage = 'owncloud'):
  suites = getUITestSuites()
  paths = ""
  for path in suites[suiteName]:
    paths = paths + "tests/acceptance/features/" + path + " "

  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': suiteName,
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate('ocis') +
      build() +
      ocisServer(storage) + [
      {
        'name': 'webUITests',
        'image': 'owncloudci/nodejs:11',
        'pull': 'always',
        'environment': {
          'SERVER_HOST': 'https://ocis-server:9200',
          'BACKEND_HOST': 'https://ocis-server:9200',
          'RUN_ON_OCIS': 'true',
          'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/ocis/owncloud',
          'OCIS_SKELETON_DIR': '/srv/app/testing/data/webUISkeleton',
          'PHOENIX_CONFIG': '/drone/src/ocis/tests/config/drone/ocis-config.json',
          'TEST_TAGS': 'not @skipOnOCIS and not @skip',
          'LOCAL_UPLOAD_DIR': '/uploads',
          'NODE_TLS_REJECT_UNAUTHORIZED': 0,
          'TEST_PATHS': paths,
        },
        'commands': [
          'git clone -b master --depth=1 https://github.com/owncloud/testing.git /srv/app/testing',
          'git clone -b %s --single-branch --no-tags https://github.com/owncloud/phoenix.git /srv/app/phoenix' % (phoenixBranch),
          'cp -r /srv/app/phoenix/tests/acceptance/filesForUpload/* /uploads',
          'cd /srv/app/phoenix',
        ] + ([
          'git checkout %s' % (phoenixCommit)
        ] if phoenixCommit != '' else []) + [
          'yarn install-all',
          'yarn run acceptance-tests-drone'
        ],
        'volumes': [{
          'name': 'gopath',
          'path': '/srv/app',
        },
        {
          'name': 'uploads',
          'path': '/uploads'
        }]
      },
    ],
    'services':
      redis() +
      selenium(),
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
      {
        'name': 'uploads',
        'temp': {}
      }
    ],
    'trigger': {
      'ref': [
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def accountsUITests(ctx, phoenixBranch, phoenixCommitId, storage):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'accountsUITests',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'build-ocis',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'commands': [
          'cd ocis',
          'make build',
          'mkdir -p /srv/app/ocis/bin',
          'cp bin/ocis /srv/app/ocis/bin',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app'
          },
        ]
      },
      {
        'name': 'ocis-server',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'detach': True,
        'environment' : {
          #'OCIS_LOG_LEVEL': 'debug',
          'STORAGE_STORAGE_HOME_DRIVER': '%s' % (storage),
          'STORAGE_STORAGE_HOME_DATA_DRIVER': '%s' % (storage),
          'STORAGE_STORAGE_OC_DRIVER': '%s' % (storage),
          'STORAGE_STORAGE_OC_DATA_DRIVER': '%s' % (storage),
          'STORAGE_STORAGE_HOME_DATA_TEMP_FOLDER': '/srv/app/tmp/',
          'STORAGE_STORAGE_OCIS_ROOT': '/srv/app/tmp/ocis/storage/users',
          'STORAGE_STORAGE_LOCAL_ROOT': '/srv/app/tmp/ocis/reva/root',
          'STORAGE_STORAGE_OWNCLOUD_DATADIR': '/srv/app/tmp/ocis/owncloud/data',
          'STORAGE_STORAGE_OC_DATA_TEMP_FOLDER': '/srv/app/tmp/',
          'STORAGE_STORAGE_OWNCLOUD_REDIS_ADDR': 'redis:6379',
          'STORAGE_OIDC_ISSUER': 'https://ocis-server:9200',
          'STORAGE_LDAP_IDP': 'https://ocis-server:9200',
          'PROXY_OIDC_ISSUER': 'https://ocis-server:9200',
          'STORAGE_STORAGE_OC_DATA_SERVER_URL': 'http://ocis-server:9164/data',
          'STORAGE_DATAGATEWAY_URL': 'https://ocis-server:9200/data',
          'STORAGE_FRONTEND_URL': 'https://ocis-server:9200',
          'PHOENIX_WEB_CONFIG': '/drone/src/accounts/ui/tests/config/drone/ocis-config.json',
          'KONNECTD_IDENTIFIER_REGISTRATION_CONF': '/drone/src/accounts/ui/tests/config/drone/identifier-registration.yml',
          'KONNECTD_ISS': 'https://ocis-server:9200',
          #'ACCOUNTS_STORAGE_DISK_PATH': '/srv/app/tmp/ocis-accounts', # Temporary workaround, don't use metadata storage
        },
        'commands': [
          'mkdir -p /srv/app/tmp/reva',
          # First run settings service because accounts need it to register the settings bundles
          '/srv/app/ocis/bin/ocis settings &',

          # Wait for the settings service to start
          "while [[ \"$(curl -s -o /dev/null -w ''%{http_code}'' localhost:9190)\" != \"404\" ]]; do sleep 2; done",

          # Now start the accounts service
          '/srv/app/ocis/bin/ocis accounts &',

          # Wait for the accounts service to start
          "while [[ \"$(curl -s -o /dev/null -w ''%{http_code}'' localhost:9181)\" != \"404\" ]]; do sleep 2; done",

          # Now run all the ocis services except the accounts and settings because they are already running
          '/srv/app/ocis/bin/ocis server',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app'
          },
        ]
      },
      {
        'name': 'WebUIAcceptanceTests',
        'image': 'owncloudci/nodejs:10',
        'pull': 'always',
        'environment': {
          'SERVER_HOST': 'https://ocis-server:9200',
          'BACKEND_HOST': 'https://ocis-server:9200',
          'RUN_ON_OCIS': 'true',
          'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/ocis/owncloud',
          'OCIS_SKELETON_DIR': '/srv/app/testing/data/webUISkeleton',
          'PHOENIX_CONFIG': '/drone/src/accounts/ui/tests/config/drone/ocis-config.json',
          'TEST_TAGS': 'not @skipOnOCIS and not @skip',
          'LOCAL_UPLOAD_DIR': '/uploads',
          'PHOENIX_PATH': '/srv/app/phoenix',
          'FEATURE_PATH': '/drone/src/accounts/ui/tests/acceptance/features',
          'NODE_TLS_REJECT_UNAUTHORIZED': '0'
        },
        'commands': [
          'git clone --depth=1 https://github.com/owncloud/testing.git /srv/app/testing',
          'git clone -b %s --single-branch https://github.com/owncloud/phoenix /srv/app/phoenix' % (phoenixBranch),
          'cd /srv/app/phoenix',
          'git checkout %s' % (phoenixCommitId),
          'cp -r /srv/app/phoenix/tests/acceptance/filesForUpload/* /uploads',
          'yarn install-all',
          'cd /drone/src/accounts',
          'yarn install --all',
          'make test-acceptance-webui'
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app',
          },
          {
            'name': 'uploads',
            'path': '/uploads'
          }
        ],
      },
    ],
    'services': [
      {
        'name': 'redis',
        'image': 'webhippie/redis',
        'pull': 'always',
        'environment': {
          'REDIS_DATABASES': 1
        },
      },
      {
        'name': 'selenium',
        'image': 'selenium/standalone-chrome-debug:3.141.59-20200326',
        'pull': 'always',
        'volumes': [
          {
            'name': 'uploads',
            'path': '/uploads'
          }
        ],
      },
    ],
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
      {
        'name': 'uploads',
        'temp': {}
      }
    ],
    'trigger': {
      'ref': [
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def docker(ctx, arch):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'docker-%s' % (arch),
    'platform': {
      'os': 'linux',
      'arch': arch,
    },
    'steps':
      generate('ocis') +
      build() + [
      {
        'name': 'dryrun',
        'image': 'plugins/docker:18.09',
        'pull': 'always',
        'settings': {
          'dry_run': True,
          'context': 'ocis',
          'tags': 'linux-%s' % (arch),
          'dockerfile': 'ocis/docker/Dockerfile.linux.%s' % (arch),
          'repo': ctx.repo.slug,
        },
        'when': {
          'ref': {
            'include': [
              'refs/pull/**',
            ],
          },
        },
      },
      {
        'name': 'docker',
        'image': 'plugins/docker:18.09',
        'pull': 'always',
        'settings': {
          'username': {
            'from_secret': 'docker_username',
          },
          'password': {
            'from_secret': 'docker_password',
          },
          'auto_tag': True,
          'context': 'ocis',
          'auto_tag_suffix': 'linux-%s' % (arch),
          'dockerfile': 'ocis/docker/Dockerfile.linux.%s' % (arch),
          'repo': ctx.repo.slug,
        },
        'when': {
          'ref': {
            'exclude': [
              'refs/pull/**',
            ],
          },
        },
      },
    ],
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'depends_on': getDependsOnAllTestPipelines(ctx),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def dockerEos(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'docker-eos-ocis',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate('ocis') +
      build() + [
        {
          'name': 'dryrun-eos-ocis',
          'image': 'plugins/docker:18.09',
          'pull': 'always',
          'settings': {
            'dry_run': True,
            'context': 'ocis/docker/eos-ocis',
            'tags': 'linux-eos-ocis',
            'dockerfile': 'ocis/docker/eos-ocis/Dockerfile',
            'repo': 'owncloud/eos-ocis',
          },
          'when': {
            'ref': {
              'include': [
                'refs/pull/**',
              ],
            },
          },
        },
        {
          'name': 'docker-eos-ocis',
          'image': 'plugins/docker:18.09',
          'pull': 'always',
          'settings': {
            'username': {
              'from_secret': 'docker_username',
            },
            'password': {
              'from_secret': 'docker_password',
            },
            'auto_tag': True,
            'context': 'ocis/docker/eos-ocis',
            'auto_tag_suffix': 'linux-amd64',
            'dockerfile': 'ocis/docker/eos-ocis/Dockerfile',
            'repo': 'owncloud/eos-ocis',
          },
          'when': {
            'ref': {
              'exclude': [
                'refs/pull/**',
              ],
            },
          },
        },
      ],
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'depends_on': getDependsOnAllTestPipelines(ctx),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def binary(ctx, name):
  if ctx.build.event == "tag":
    settings = {
      'endpoint': {
        'from_secret': 's3_endpoint',
      },
      'access_key': {
        'from_secret': 'aws_access_key_id',
      },
      'secret_key': {
        'from_secret': 'aws_secret_access_key',
      },
      'bucket': {
        'from_secret': 's3_bucket',
      },
      'path_style': True,
      'strip_prefix': 'ocis/dist/release/',
      'source': 'ocis/dist/release/*',
      'target': '/ocis/%s/%s' % (ctx.repo.name.replace("ocis-", ""), ctx.build.ref.replace("refs/tags/v", "")),
    }
  else:
    settings = {
      'endpoint': {
        'from_secret': 's3_endpoint',
      },
      'access_key': {
        'from_secret': 'aws_access_key_id',
      },
      'secret_key': {
        'from_secret': 'aws_secret_access_key',
      },
      'bucket': {
        'from_secret': 's3_bucket',
      },
      'path_style': True,
      'strip_prefix': 'dist/release/',
      'source': 'ocis/dist/release/*',
      'target': '/ocis/%s/testing' % (ctx.repo.name.replace("ocis-", "")),
    }

  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'binaries-%s' % (name),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate('ocis') + [
      {
        'name': 'build',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'commands': [
          'cd ocis',
          'make release-%s' % (name),
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app',
          },
        ],
      },
      {
        'name': 'finish',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'commands': [
          'cd ocis',
          'make release-finish',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app',
          },
        ],
        'when': {
          'ref': [
            'refs/heads/master',
            'refs/tags/v*',
          ],
        },
      },
      {
        'name': 'upload',
        'image': 'plugins/s3:1',
        'pull': 'always',
        'settings': settings,
        'when': {
          'ref': [
            'refs/heads/master',
            'refs/tags/v*',
          ],
        },
      },
      {
        'name': 'changelog',
        'image': 'toolhippie/calens:latest',
        'pull': 'always',
        'commands': [
          'cd ocis',
          'calens --version %s -o dist/CHANGELOG.md' % ctx.build.ref.replace("refs/tags/v", "").split("-")[0],
        ],
        'when': {
          'ref': [
            'refs/tags/v*',
          ],
        },
      },
      {
        'name': 'release',
        'image': 'plugins/github-release:1',
        'pull': 'always',
        'settings': {
          'api_key': {
            'from_secret': 'github_token',
          },
          'files': [
            'ocis/dist/release/*',
          ],
          'title': ctx.build.ref.replace("refs/tags/v", ""),
          'note': 'ocis/dist/CHANGELOG.md',
          'overwrite': True,
          'prerelease': len(ctx.build.ref.split("-")) > 1,
        },
        'when': {
          'ref': [
            'refs/tags/v*',
          ],
        },
      },
    ],
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'depends_on': getDependsOnAllTestPipelines(ctx),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def releaseSubmodule(ctx):
  depends = []
  if len(ctx.build.ref.replace("refs/tags/", "").split("/")) == 2:
    depends = ['linting&unitTests-%s' % (ctx.build.ref.replace("refs/tags/", "").split("/")[0])]

  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'release-%s' % (ctx.build.ref.replace("refs/tags/", "")),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps' : [
      {
        'name': 'release-submodule',
        'image': 'plugins/github-release:1',
        'pull': 'always',
        'settings': {
          'api_key': {
            'from_secret': 'github_token',
          },
          'files': [
          ],
          'title': ctx.build.ref.replace("refs/tags/", "").replace("/v", " "),
          'note': 'Release %s submodule' % (ctx.build.ref.replace("refs/tags/", "").replace("/v", " ")),
          'overwrite': True,
          'prerelease': len(ctx.build.ref.split("-")) > 1,
        },
        'when': {
          'ref': [
            'refs/tags/*/v*',
          ],
        },
      },
    ],
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'depends_on': depends,
    'trigger': {
      'ref': [
        'refs/tags/*/v*',
      ],
    },
  }


def manifest(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'manifest',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'execute',
        'image': 'plugins/manifest:1',
        'pull': 'always',
        'settings': {
          'username': {
            'from_secret': 'docker_username',
          },
          'password': {
            'from_secret': 'docker_password',
          },
          'spec': 'ocis/docker/manifest.tmpl',
          'auto_tag': True,
          'ignore_missing': True,
        },
      },
    ],
    'depends_on': [
      'docker-amd64',
      'docker-arm64',
      'docker-arm',
      'docker-eos-ocis',
      'binaries-linux',
      'binaries-darwin',
      'binaries-windows',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
      ],
    },
  }

def changelog(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'changelog',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'generate',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'commands': [
          'cd ocis',
          'make changelog',
        ],
      },
      {
        'name': 'diff',
        'image': 'owncloud/alpine:latest',
        'pull': 'always',
        'commands': [
          'git diff',
        ],
      },
      {
        'name': 'output',
        'image': 'owncloud/alpine:latest',
        'pull': 'always',
        'commands': [
          'cat CHANGELOG.md',
        ],
      },
      {
        'name': 'publish',
        'image': 'plugins/git-action:1',
        'pull': 'always',
        'settings': {
          'actions': [
            'commit',
            'push',
          ],
          'message': 'Automated changelog update [skip ci]',
          'branch': 'master',
          'author_email': 'devops@owncloud.com',
          'author_name': 'ownClouders',
          'netrc_machine': 'github.com',
          'netrc_username': {
            'from_secret': 'github_username',
          },
          'netrc_password': {
            'from_secret': 'github_token',
          },
        },
        'when': {
          'ref': {
            'exclude': [
              'refs/pull/**',
            ],
          },
        },
      },
    ],
    'depends_on': [
      'manifest',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/pull/**',
      ],
    },
  }

def readme(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'readme',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'execute',
        'image': 'sheogorath/readme-to-dockerhub:latest',
        'pull': 'always',
        'environment': {
          'DOCKERHUB_USERNAME': {
            'from_secret': 'docker_username',
          },
          'DOCKERHUB_PASSWORD': {
            'from_secret': 'docker_password',
          },
          'DOCKERHUB_REPO_PREFIX': ctx.repo.namespace,
          'DOCKERHUB_REPO_NAME': ctx.repo.name,
          'SHORT_DESCRIPTION': 'Docker images for %s' % (ctx.repo.name),
          'README_PATH': 'README.md',
        },
      },
    ],
    'depends_on': [
      'docker-amd64',
      'docker-arm64',
      'docker-arm',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
      ],
    },
  }

def badges(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'badges',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'execute',
        'image': 'plugins/webhook:1',
        'pull': 'always',
        'settings': {
          'urls': {
            'from_secret': 'microbadger_url',
          },
        },
      },
    ],
    'depends_on': [
      'docker-amd64',
      'docker-arm64',
      'docker-arm',
      'docker-eos-ocis',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
      ],
    },
  }

def docs(ctx):
  generateConfigDocs = []

  for module in config['modules']:
    generateConfigDocs += ['cd /drone/src/%s' % (module), 'make config-docs-generate']

  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'docs',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'prepare',
        'image': 'owncloudci/alpine:latest',
        'commands': [
          'make docs-copy'
        ],
      },
      {
        'name': 'generate-config-docs',
        'image': 'webhippie/golang:1.13',
        'commands': generateConfigDocs,
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app',
          },
        ],
      },
      {
        'name': 'test',
        'image': 'owncloudci/hugo:0.71.0',
        'commands': [
          'cd hugo',
          'hugo',
        ],
      },
      {
        'name': 'list',
        'image': 'owncloudci/alpine:latest',
        'commands': [
          'tree hugo/public',
        ],
      },
      {
        'name': 'publish',
        'image': 'plugins/gh-pages:1',
        'pull': 'always',
        'settings': {
          'username': {
            'from_secret': 'github_username',
          },
          'password': {
            'from_secret': 'github_token',
          },
          'pages_directory': 'docs/',
          'target_branch': 'docs',
        },
        'when': {
          'ref': {
            'exclude': [
              'refs/pull/**',
            ],
          },
        },
      },
      {
        'name': 'downstream',
        'image': 'plugins/downstream',
        'settings': {
          'server': 'https://drone.owncloud.com/',
          'token': {
            'from_secret': 'drone_token',
          },
          'repositories': [
            'owncloud/owncloud.github.io@source',
          ],
        },
        'when': {
          'ref': {
            'exclude': [
              'refs/pull/**',
            ],
          },
        },
      },
    ],
    'depends_on': [
      'badges',
    ],
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/pull/**',
      ],
    },
  }

def generate(module):
  return [
    {
      'name': 'generate',
      'image': 'webhippie/golang:1.13',
      'pull': 'always',
      'commands': [
        'cd %s' % (module),
        'make generate',
      ],
      'volumes': [
        {
          'name': 'gopath',
          'path': '/srv/app',
        },
      ],
    }
  ]

def updateDeployment(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'updateDeployment',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'webhook',
        'image': 'plugins/webhook',
        'settings': {
          'username': {
            'from_secret': 'webhook_username',
          },
          'password': {
            'from_secret': 'webhook_password',
          },
          'method': 'GET',
          'urls': 'https://ocis.owncloud.works/hooks/update-ocis',
        }
      }
    ],
    'depends_on': [
      'docker-amd64',
      'docker-arm64',
      'docker-arm',
      'binaries-linux',
      'binaries-darwin',
      'binaries-windows',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
      ],
    }
  }

def frontend(module):
  return [
    {
      'name': 'frontend',
      'image': 'webhippie/nodejs:latest',
      'pull': 'always',
      'commands': [
        'cd %s' % (module),
        'yarn install --frozen-lockfile',
        'yarn lint',
        'yarn test',
        'yarn build',
      ],
    }
  ]

def ocisServer(storage):
  return [
    {
      'name': 'ocis-server',
      'image': 'webhippie/golang:1.13',
      'pull': 'always',
      'detach': True,
      'environment' : {
        #'OCIS_LOG_LEVEL': 'debug',
        'STORAGE_STORAGE_HOME_DRIVER': '%s' % (storage),
        'STORAGE_STORAGE_HOME_DATA_DRIVER': '%s' % (storage),
        'STORAGE_STORAGE_OC_DRIVER': '%s' % (storage),
        'STORAGE_STORAGE_OC_DATA_DRIVER': '%s' % (storage),
        'STORAGE_STORAGE_HOME_DATA_TEMP_FOLDER': '/srv/app/tmp/',
        'STORAGE_STORAGE_OCIS_ROOT': '/srv/app/tmp/ocis/storage/users',
        'STORAGE_STORAGE_LOCAL_ROOT': '/srv/app/tmp/ocis/local/root',
        'STORAGE_STORAGE_OWNCLOUD_DATADIR': '/srv/app/tmp/ocis/owncloud/data',
        'STORAGE_STORAGE_OC_DATA_TEMP_FOLDER': '/srv/app/tmp/',
        'STORAGE_STORAGE_OWNCLOUD_REDIS_ADDR': 'redis:6379',
        'STORAGE_LDAP_IDP': 'https://ocis-server:9200',
        'STORAGE_OIDC_ISSUER': 'https://ocis-server:9200',
        'PROXY_OIDC_ISSUER': 'https://ocis-server:9200',
        'STORAGE_STORAGE_OC_DATA_SERVER_URL': 'http://ocis-server:9164/data',
        'STORAGE_DATAGATEWAY_URL': 'https://ocis-server:9200/data',
        'STORAGE_FRONTEND_URL': 'https://ocis-server:9200',
        'PHOENIX_WEB_CONFIG': '/drone/src/ocis/tests/config/drone/ocis-config.json',
        'KONNECTD_IDENTIFIER_REGISTRATION_CONF': '/drone/src/ocis/tests/config/drone/identifier-registration.yml',
        'KONNECTD_ISS': 'https://ocis-server:9200',
        'KONNECTD_TLS': 'true',
        'ACCOUNTS_DATA_PATH': '/srv/app/tmp/ocis-accounts/',
      },
      'commands': [
        'apk add mailcap', # install /etc/mime.types
        'mkdir -p /srv/app/tmp/ocis/owncloud/data/',
        'mkdir -p /srv/app/tmp/ocis/storage/users/',
        'ocis/bin/ocis server'
      ],
      'volumes': [
        {
          'name': 'gopath',
          'path': '/srv/app'
        },
      ]
    },
  ]

def cloneCoreRepos(coreBranch, coreCommit):
  return [
    {
      'name': 'clone-core-repos',
      'image': 'owncloudci/php:7.2',
      'pull': 'always',
      'commands': [
        'git clone -b master --depth=1 https://github.com/owncloud/testing.git /srv/app/tmp/testing',
        'git clone -b %s --single-branch --no-tags https://github.com/owncloud/core.git /srv/app/testrunner' % (coreBranch),
        'cd /srv/app/testrunner',
      ] + ([
        'git checkout %s' % (coreCommit)
      ] if coreCommit != '' else []),
      'volumes': [{
        'name': 'gopath',
        'path': '/srv/app',
      }]
    }
  ]

def redis():
  return [
    {
      'name': 'redis',
      'image': 'webhippie/redis',
      'pull': 'always',
      'environment': {
        'REDIS_DATABASES': 1
      },
    }
  ]

def selenium():
  return [
    {
      'name': 'selenium',
      'image': 'selenium/standalone-chrome-debug:3.141.59-20200326',
      'pull': 'always',
      'volumes': [{
          'name': 'uploads',
          'path': '/uploads'
      }],
    }
  ]

def build():
  return [
    {
      'name': 'build',
      'image': 'webhippie/golang:1.13',
      'pull': 'always',
      'commands': [
        'cd ocis',
        'make build',
      ],
      'volumes': [
        {
          'name': 'gopath',
          'path': '/srv/app',
        },
      ],
    },
  ]
