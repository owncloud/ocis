config = {
  'apiTests': {
    'coreBranch': 'master',
    'coreCommit': 'e3aba53ada536915cd0e1bdbf33527d9970a13d3',
    'numberOfParts': 2
  },
  'uiTests': {
    'phoenixBranch': 'master',
    'phoenixCommit': 'ccdca163c7e9e6ecc57e08a298b08b1d1175f1d5',
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

def getUITestSuiteNames():
  return config['uiTests']['suites'].keys()

def getUITestSuites():
  return config['uiTests']['suites']

def getCoreApiTestPipelineNames():
  names = []
  for runPart in range(1, config['apiTests']['numberOfParts'] + 1):
    names.append('coreApiTests-%s' % runPart)
  return names

def main(ctx):
  before = testPipelines(ctx)

  stages = [
    docker(ctx, 'amd64'),
    docker(ctx, 'arm64'),
    docker(ctx, 'arm'),
    binary(ctx, 'linux'),
    binary(ctx, 'darwin'),
    binary(ctx, 'windows'),
  ]

  after = [
    manifest(ctx),
    changelog(ctx),
    readme(ctx),
    badges(ctx),
    website(ctx),
  ]

  return before + stages + after

def testPipelines(ctx):

  pipelines = [
    linting(ctx),
    unitTests(ctx),
    localApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'])
  ]

  for runPart in range(1, config['apiTests']['numberOfParts'] + 1):
    pipelines.append(coreApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], runPart, config['apiTests']['numberOfParts']))

  pipelines += uiTests(ctx, config['uiTests']['phoenixBranch'], config['uiTests']['phoenixCommit'])
  return pipelines

def linting(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'linting',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate() + [
      {
        'name': 'vet',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'commands': [
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
          'make lint',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app',
          },
        ],
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

def unitTests(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'unitTests',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate() +
      build() + [
      {
        'name': 'test',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'commands': [
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
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**',
      ],
    },
  }

def localApiTests(ctx, coreBranch = 'master', coreCommit = ''):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'localApiTests',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate() +
      build() +
      ocisServer() +
      cloneCoreRepos(coreBranch, coreCommit) + [
      {
        'name': 'LocalApiTests',
        'image': 'owncloudci/php:7.2',
        'pull': 'always',
        'environment' : {
          'TEST_SERVER_URL': 'http://ocis-server:9140',
          'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/reva/',
          'SKELETON_DIR': '/srv/app/tmp/testing/data/apiSkeleton',
          'TEST_EXTERNAL_USER_BACKENDS':'true',
          'REVA_LDAP_HOSTNAME':'ldap',
          'TEST_OCIS':'true',
          'PATH_TO_CORE': '/srv/app/testrunner'
        },
        'commands': [
          'make test-acceptance-api',
        ],
        'volumes': [{
          'name': 'gopath',
          'path': '/srv/app',
        }]
      },
    ],
    'services':
      ldap() +
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

def coreApiTests(ctx, coreBranch = 'master', coreCommit = '', part_number = 1, number_of_parts = 1):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'coreApiTests-%s' % (part_number),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate() +
      build() +
      ocisServer() +
      cloneCoreRepos(coreBranch, coreCommit) + [
      {
        'name': 'oC10ApiTests-%s' % (part_number),
        'image': 'owncloudci/php:7.2',
        'pull': 'always',
        'environment' : {
          'TEST_SERVER_URL': 'http://ocis-server:9140',
          'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/reva/',
          'SKELETON_DIR': '/srv/app/tmp/testing/data/apiSkeleton',
          'TEST_EXTERNAL_USER_BACKENDS':'true',
          'REVA_LDAP_HOSTNAME':'ldap',
          'TEST_OCIS':'true',
          'BEHAT_FILTER_TAGS': '~@notToImplementOnOCIS&&~@toImplementOnOCIS&&~@preview-extension-required',
          'DIVIDE_INTO_NUM_PARTS': number_of_parts,
          'RUN_PART': part_number,
          'EXPECTED_FAILURES_FILE': '/drone/src/tests/acceptance/expected-failures.txt'
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
      ldap() +
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

def uiTestPipeline(suiteName, phoenixBranch = 'master', phoenixCommit = ''):
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
      generate() +
      build() +
      ocisServer() + [
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
          'PHOENIX_CONFIG': '/drone/src/tests/config/drone/ocis-config.json',
          'LDAP_SERVER_URL': 'ldap://ldap',
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
      ldap() +
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
    'name': arch,
    'platform': {
      'os': 'linux',
      'arch': arch,
    },
    'steps':
      generate() +
      build() + [
      {
        'name': 'dryrun',
        'image': 'plugins/docker:18.09',
        'pull': 'always',
        'settings': {
          'dry_run': True,
          'tags': 'linux-%s' % (arch),
          'dockerfile': 'docker/Dockerfile.linux.%s' % (arch),
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
          'auto_tag_suffix': 'linux-%s' % (arch),
          'dockerfile': 'docker/Dockerfile.linux.%s' % (arch),
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
    'depends_on': [
      'linting',
      'unitTests',
      'localApiTests',
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
        'from_secret': 'aws_secret_access_key',
      },
      'bucket': {
        'from_secret': 's3_bucket',
      },
      'path_style': True,
      'strip_prefix': 'dist/release/',
      'source': 'dist/release/*',
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
      'source': 'dist/release/*',
      'target': '/ocis/%s/testing' % (ctx.repo.name.replace("ocis-", "")),
    }

  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': name,
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      generate() + [
      {
        'name': 'build',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'commands': [
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
          'make release-finish',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/srv/app',
          },
        ],
      },
      {
        'name': 'upload',
        'image': 'plugins/s3:1',
        'pull': 'always',
        'settings': settings,
        'when': {
          'ref': [
            'refs/heads/master',
            'refs/tags/**',
          ],
        },
      },
      {
        'name': 'changelog',
        'image': 'toolhippie/calens:latest',
        'pull': 'always',
        'commands': [
          'calens --version %s -o dist/CHANGELOG.md' % ctx.build.ref.replace("refs/tags/v", "").split("-")[0],
        ],
        'when': {
          'ref': [
            'refs/tags/**',
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
            'dist/release/*',
          ],
          'title': ctx.build.ref.replace("refs/tags/v", ""),
          'note': 'dist/CHANGELOG.md',
          'overwrite': True,
          'prerelease': len(ctx.build.ref.split("-")) > 1,
        },
        'when': {
          'ref': [
            'refs/tags/**',
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
    'depends_on': [
      'linting',
      'unitTests',
      'localApiTests',
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
          'spec': 'docker/manifest.tmpl',
          'auto_tag': True,
          'ignore_missing': True,
        },
      },
    ],
    'depends_on': [
      'amd64',
      'arm64',
      'arm',
      'linux',
      'darwin',
      'windows',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
      ],
    },
  }

def changelog(ctx):
  repo_slug = ctx.build.source_repo if ctx.build.source_repo else ctx.repo.slug
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'changelog',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'clone': {
      'disable': True,
    },
    'steps': [
      {
        'name': 'clone',
        'image': 'plugins/git-action:1',
        'pull': 'always',
        'settings': {
          'actions': [
            'clone',
          ],
          'remote': 'https://github.com/%s' % (repo_slug),
          'branch': ctx.build.source if ctx.build.event == 'pull_request' else 'master',
          'path': '/drone/src',
          'netrc_machine': 'github.com',
          'netrc_username': {
            'from_secret': 'github_username',
          },
          'netrc_password': {
            'from_secret': 'github_token',
          },
        },
      },
      {
        'name': 'generate',
        'image': 'webhippie/golang:1.13',
        'pull': 'always',
        'commands': [
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
      'changelog',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
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
      'readme',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
      ],
    },
  }

def website(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'website',
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
        'name': 'test',
        'image': 'webhippie/hugo:latest',
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
          'server': 'https://cloud.drone.io/',
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
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/pull/**',
      ],
    },
  }

def generate():
  return [
    {
      'name': 'generate',
      'image': 'webhippie/golang:1.13',
      'pull': 'always',
      'commands': [
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

def build():
  return [
    {
      'name': 'build',
      'image': 'webhippie/golang:1.13',
      'pull': 'always',
      'commands': [
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

def ocisServer():
  return [
    {
      'name': 'ocis-server',
      'image': 'webhippie/golang:1.13',
      'pull': 'always',
      'detach': True,
      'environment' : {
        'REVA_LDAP_HOSTNAME': 'ldap',
        'REVA_LDAP_PORT': 636,
        'REVA_LDAP_BIND_PASSWORD': 'admin',
        'REVA_LDAP_BIND_DN': 'cn=admin,dc=owncloud,dc=com',
        'REVA_LDAP_BASE_DN': 'dc=owncloud,dc=com',
        'REVA_LDAP_SCHEMA_UID': 'uid',
        'REVA_LDAP_SCHEMA_MAIL': 'mail',
        'REVA_LDAP_SCHEMA_DISPLAYNAME': 'displayName',
        'REVA_STORAGE_HOME_DATA_TEMP_FOLDER': '/srv/app/tmp/',
        'REVA_STORAGE_LOCAL_ROOT': '/srv/app/tmp/reva/root',
        'REVA_STORAGE_OWNCLOUD_DATADIR': '/srv/app/tmp/reva/data',
        'REVA_STORAGE_OC_DATA_TEMP_FOLDER': '/srv/app/tmp/',
        'REVA_STORAGE_OWNCLOUD_REDIS_ADDR': 'redis:6379',
        'REVA_OIDC_ISSUER': 'https://ocis-server:9200',
        'REVA_STORAGE_OC_DATA_SERVER_URL': 'http://ocis-server:9164/data',
        'REVA_DATAGATEWAY_URL': 'https://ocis-server:9200/data',
        'REVA_FRONTEND_URL': 'https://ocis-server:9200',
        'PHOENIX_WEB_CONFIG': '/drone/src/tests/config/drone/ocis-config.json',
        'KONNECTD_IDENTIFIER_REGISTRATION_CONF': '/drone/src/tests/config/drone/identifier-registration.yml',
        'KONNECTD_ISS': 'https://ocis-server:9200',
        'KONNECTD_TLS': 'true',
        'LDAP_URI': 'ldap://ldap',
        'LDAP_BINDDN': 'cn=admin,dc=owncloud,dc=com',
        'LDAP_BINDPW': 'admin',
        'LDAP_BASEDN': 'dc=owncloud,dc=com'
      },
      'commands': [
        'apk add mailcap', # install /etc/mime.types
        'mkdir -p /srv/app/tmp/reva',
        'bin/ocis server'
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

def ldap():
  return [
    {
      'name': 'ldap',
      'image': 'osixia/openldap',
      'pull': 'always',
      'environment': {
        'LDAP_DOMAIN': 'owncloud.com',
        'LDAP_ORGANISATION': 'ownCloud',
        'LDAP_ADMIN_PASSWORD': 'admin',
        'LDAP_TLS_VERIFY_CLIENT': 'never',
      },
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
