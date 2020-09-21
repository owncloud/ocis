config = {
  'modules': {
    'accounts': 'frontend',
    'glauth':'',
    'graph': '',
    'ocis-reva':'',
    'ocis-pkg':'',
    'ocis-phoenix':'',
    'store':'',
    'ocs':'',
    'webdav':'',
    'thumbnails':'',
    'proxy':'',
    'settings':'frontend',
    'konnectd':'',
  },
  'apiTests': {
    'coreBranch': 'master',
    'coreCommit': '2c5bb68fc689d7e9dd912125680c0fad99528fa9',
    'numberOfParts': 4
  },
  'uiTests': {
    'phoenixBranch': 'master',
    'phoenixCommit': '2ccec3966abebc7cff54cfdbd0aa9942253584a7',
    'suites': {
      'phoenixWebUI1': [
        'webUICreateFilesFolders',
        'webUIDeleteFilesFolders',
        'webUIFavorites',
        'webUIFiles',
        'webUILogin',
        'webUINotifications',
      ],
      'phoenixWebUI2': [
        'webUIPrivateLinks',
        'webUIRenameFiles',
        'webUIRenameFolders',
        'webUITrashbin',
        'webUIUpload',
        'webUIAccount',
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

def main(ctx):
  before = testPipelines(ctx)

  stages = [
    docker(ctx, 'amd64'),
    #docker(ctx, 'arm64'),
    #docker(ctx, 'arm'),
    binary(ctx, 'linux'),
    binary(ctx, 'darwin'),
    binary(ctx, 'windows'),
  ]

  after = [
  ]

  return before + stages + after

def testPipelines(ctx):
  pipelines = []

  for module in config['modules']:
    pipelines.append(testing(ctx, module))

  pipelines += [
    localApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], 'owncloud'),
    localApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], 'ocis')
  ]

  for runPart in range(1, config['apiTests']['numberOfParts'] + 1):
    pipelines.append(coreApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], runPart, config['apiTests']['numberOfParts'], 'owncloud'))
    pipelines.append(coreApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], runPart, config['apiTests']['numberOfParts'], 'ocis'))

  pipelines += uiTests(ctx, config['uiTests']['phoenixBranch'], config['uiTests']['phoenixCommit'])
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
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app',
          },
        ],
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
        'refs/tags/**',
        'refs/pull/**',
      ],
    },
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
          'OCIS_REVA_DATA_ROOT': '%s' % ('/srv/app/tmp/reva/' if storage == 'owncloud' else ''),
          'DELETE_USER_DATA_CMD': '%s' % ('rm -rf /srv/app/tmp/reva/data/*' if storage == 'owncloud' else 'rm -rf /srv/app/tmp/ocis/root/nodes/root/*'),
          'SKELETON_DIR': '/srv/app/tmp/testing/data/apiSkeleton',
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
        'refs/heads/master',
        'refs/tags/**',
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
          'OCIS_REVA_DATA_ROOT': '%s' % ('/srv/app/tmp/reva/' if storage == 'owncloud' else ''),
          'DELETE_USER_DATA_CMD': '%s' % ('rm -rf /srv/app/tmp/reva/data/*' if storage == 'owncloud' else 'rm -rf /srv/app/tmp/ocis/root/nodes/root/*'),
          'SKELETON_DIR': '/srv/app/tmp/testing/data/apiSkeleton',
          'TEST_OCIS':'true',
          'BEHAT_FILTER_TAGS': '~@notToImplementOnOCIS&&~@toImplementOnOCIS&&~comments-app-required&&~@federation-app-required&&~@notifications-app-required&&~systemtags-app-required&&~@local_storage',
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
        'refs/heads/master',
        'refs/tags/**',
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
          'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/reva',
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
        'refs/heads/master',
        'refs/tags/**',
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
            'from_secret': '*docker_username*',
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
              'refs/tags/*/v*',
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
    'depends_on': 
      getTestSuiteNames() + [
      'localApiTests-owncloud-storage',
      'localApiTests-ocis-storage',
    ] + getCoreApiTestPipelineNames() + getUITestSuiteNames(),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
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
        'from_secret': '*aws_secret_access_key*',
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
        'from_secret': '*aws_secret_access_key*',
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
          'title': ctx.build.ref.replace("refs/tags/", ""),
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
    'depends_on': 
      getTestSuiteNames() + [
      'localApiTests-owncloud-storage',
      'localApiTests-ocis-storage',
    ] + getCoreApiTestPipelineNames() + getUITestSuiteNames(),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**',
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
      #'docker-arm64',
      #'docker-arm',
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
      #'docker-arm64',
      #'docker-arm',
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
      #'docker-arm64',
      #'docker-arm',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
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
        'OCIS_LOG_LEVEL': 'debug',
        'REVA_STORAGE_HOME_DRIVER': '%s' % (storage),
        'REVA_STORAGE_HOME_DATA_DRIVER': '%s' % (storage),
        'REVA_STORAGE_OC_DRIVER': '%s' % (storage),
        'REVA_STORAGE_OC_DATA_DRIVER': '%s' % (storage),
        'REVA_STORAGE_HOME_DATA_TEMP_FOLDER': '/srv/app/tmp/',
        'REVA_STORAGE_OCIS_ROOT': '/srv/app/tmp/ocis/root',
        'REVA_STORAGE_LOCAL_ROOT': '/srv/app/tmp/reva/root',
        'REVA_STORAGE_OWNCLOUD_DATADIR': '/srv/app/tmp/reva/data',
        'REVA_STORAGE_OC_DATA_TEMP_FOLDER': '/srv/app/tmp/',
        'REVA_STORAGE_OWNCLOUD_REDIS_ADDR': 'redis:6379',
        'REVA_LDAP_IDP': 'https://ocis-server:9200',
        'REVA_OIDC_ISSUER': 'https://ocis-server:9200',
        'PROXY_OIDC_ISSUER': 'https://ocis-server:9200',
        'REVA_STORAGE_OC_DATA_SERVER_URL': 'http://ocis-server:9164/data',
        'REVA_DATAGATEWAY_URL': 'https://ocis-server:9200/data',
        'REVA_FRONTEND_URL': 'https://ocis-server:9200',
        'PHOENIX_WEB_CONFIG': '/drone/src/ocis/tests/config/drone/ocis-config.json',
        'KONNECTD_IDENTIFIER_REGISTRATION_CONF': '/drone/src/ocis/tests/config/drone/identifier-registration.yml',
        'KONNECTD_ISS': 'https://ocis-server:9200',
        'KONNECTD_TLS': 'true',
      },
      'commands': [
        'apk add mailcap', # install /etc/mime.types
        'mkdir -p /srv/app/tmp/reva',
        'mkdir -p /srv/app/tmp/ocis/root/nodes',
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
