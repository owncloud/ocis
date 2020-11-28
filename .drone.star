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
    'coreCommit': '31105a0d3e6e6ef0a0da2c16c66d4261fd91c069',
    'numberOfParts': 10
  },
  'uiTests': {
    'phoenixBranch': 'master',
    'phoenixCommit': 'fd281418cc30af9c9795e627692caf45c0d3bf30',
      'suites': {
        'webUIBasic': [
          'webUILogin',
          'webUINotifications',
          'webUIPrivateLinks',
          'webUIPreview',
          'webUIAccount',
        ],
        'webUICreateFilesFolders': 'webUICreateFilesFolders',
        'webUIDeleteFilesFolders': 'webUIDeleteFilesFolders',
        'webUIRename': [
          'webUIRenameFiles',
          'webUIRenameFolders',
        ],
        'webUISharingBasic': [
          'webUISharingAcceptShares',
          'webUISharingAcceptSharesToRoot',
        ],
        'webUIRestrictSharing': 'webUIRestrictSharing',
        'webUISharingNotifications': [
          'webUISharingNotifications',
          'webUISharingNotificationsToRoot',
        ],
        'webUIFavorites': 'webUIFavorites',
        'webUIFiles': 'webUIFiles',
        'webUISharingAutocompletion': 'webUISharingAutocompletion',
        'webUISharingInternalGroups': [
          'webUISharingInternalGroups',
          'webUISharingInternalGroupsEdgeCases',
          'webUISharingInternalGroupsSharingIndicator',
          'webUISharingInternalGroupsToRoot',
          'webUISharingInternalGroupsToRootEdgeCases',
          'webUISharingInternalGroupsToRootSharingIndicator',
        ],
        'webUISharingInternalUsers': [
          'webUISharingInternalUsers',
          'webUISharingInternalUsersBlacklisted',
          'webUISharingInternalUsersSharingIndicator',
          'webUISharingInternalUsersToRoot',
          'webUISharingInternalUsersToRootBlacklisted',
          'webUISharingInternalUsersToRootSharingIndicator',
        ],
        'webUISharingInternalUsersExpire': 'webUISharingInternalUsersExpire',
        'webUISharingInternalUsersExpireToRoot': 'webUISharingInternalUsersExpireToRoot',
        'webUISharingPermissionsUsers': 'webUISharingPermissionsUsers',
        'webUISharingFilePermissionsGroups': 'webUISharingFilePermissionsGroups',
        'webUISharingFolderPermissionsGroups': 'webUISharingFolderPermissionsGroups',
        'webUISharingFolderAdvancedPermissionsGroups': 'webUISharingFolderAdvPermissionsGrp',
        'webUISharingPermissionToRoot': 'webUISharingPermissionToRoot',
        'webUIResharing': 'webUIResharing',
        'webUIResharingToRoot': 'webUIResharingToRoot',
        'webUISharingPublic': 'webUISharingPublic',
        'webUISharingPublicExpire': 'webUISharingPublicExpire',
        'webUISharingPublicDifferentRoles': 'webUISharingPublicDifferentRoles',
        'webUITrashbin': 'webUITrashbin',
        'webUITrashbinFilesFolders': 'webUITrashbinFilesFolders',
        'webUITrashbinRestore': 'webUITrashbinRestore',
        'webUIUpload': 'webUIUpload',
        'webUISharingFilePermissionMultipleUsers': 'webUISharingFilePermissionMultipleUsers',
        'webUISharingFolderPermissionMultipleUsers': 'webUISharingFolderPermissionMultipleUsers',
        'webUISharingFolderAdvancedPermissionMultipleUsers': [
            'webUISharingFolderAdvancedPermissionMU',
        ],
        'webUIMoveFilesFolders': 'webUIMoveFilesFolders',
      },
  },
  'rocketchat': {
    'channel': 'ocis-internal',
    'from_secret': 'private_rocketchat',
  },
}

def getPipelineNames(pipelines=[]):
  """getPipelineNames returns names of pipelines as a string array

  Args:
    pipelines: array of drone pipelines

  Returns:
    names of the given pipelines as string array
  """
  names = []
  for pipeline in pipelines:
    names.append(pipeline['name'])
  return names

def main(ctx):
  """main is the entrypoint for drone

  Args:
    ctx: drone passes a context with information which the pipeline can be adapted to

  Returns:
    none
  """

  pipelines = []

  before = \
    [ buildOcisBinaryForTesting(ctx) ] + \
    testOcisModules(ctx) + \
    testPipelines(ctx)

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

  purge = purgeBuildArtifactCache(ctx, 'ocis-binary-amd64')
  purge['depends_on'] = getPipelineNames(testPipelines(ctx))

  after = [
    manifest(ctx),
    changelog(ctx),
    readme(ctx),
    badges(ctx),
    docs(ctx),
    updateDeployment(ctx),
    purge,
  ]

  if ctx.build.event == "cron":
    notify_pipeline = notify(ctx)
    notify_pipeline['depends_on'] = \
      getPipelineNames(before)

    pipelines = before + [ notify_pipeline ]

  elif \
  (ctx.build.event == "pull" and '[docs-only]' in ctx.build.title) \
  or \
  (ctx.build.event != "pull" and '[docs-only]' in (ctx.build.title + ctx.build.message)):
  # [docs-only] is not taken from PR messages, but from commit messages

    docs_pipeline = docs(ctx)
    docs_pipeline['depends_on'] = []
    docs_pipelines = [ docs_pipeline ]

    notify_pipeline = notify(ctx)
    notify_pipeline['depends_on'] = \
      getPipelineNames(docs_pipelines)

    pipelines = docs_pipelines + [ notify_pipeline ]

  else:
    pipelines = before + stages + after

    notify_pipeline = notify(ctx)
    notify_pipeline['depends_on'] = \
      getPipelineNames(pipelines)

    pipelines = pipelines + [ notify_pipeline ]

  return pipelines

def testOcisModules(ctx):
  pipelines = []
  for module in config['modules']:
    pipelines.append(testOcisModule(ctx, module))

  coverage_upload = uploadCoverage(ctx)
  coverage_upload['depends_on'] = getPipelineNames(pipelines)

  return pipelines + [coverage_upload]

def testPipelines(ctx):
  pipelines = [
    localApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], 'owncloud', 'apiOcisSpecific'),
    localApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], 'ocis', 'apiOcisSpecific'),
    localApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], 'owncloud', 'apiBasic', 'default'),
    localApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], 'ocis', 'apiBasic', 'default')
  ]

  for runPart in range(1, config['apiTests']['numberOfParts'] + 1):
    pipelines.append(coreApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], runPart, config['apiTests']['numberOfParts'], 'owncloud'))
    pipelines.append(coreApiTests(ctx, config['apiTests']['coreBranch'], config['apiTests']['coreCommit'], runPart, config['apiTests']['numberOfParts'], 'ocis'))

  pipelines += uiTests(ctx, config['uiTests']['phoenixBranch'], config['uiTests']['phoenixCommit'])
  pipelines.append(accountsUITests(ctx, config['uiTests']['phoenixBranch'], config['uiTests']['phoenixCommit']))
  return pipelines

def testOcisModule(ctx, module):
  steps = makeGenerate(module) + [
    {
      'name': 'vet',
      'image': 'webhippie/golang:1.14',
      'pull': 'always',
      'commands': [
        'make -C %s vet' % (module),
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
      'image': 'webhippie/golang:1.14',
      'pull': 'always',
      'commands': [
        'make -C %s staticcheck' % (module),
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
      'image': 'webhippie/golang:1.14',
      'pull': 'always',
      'commands': [
        'make -C %s lint' % (module),
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
        'image': 'webhippie/golang:1.14',
        'pull': 'always',
        'commands': [
          'make -C %s test' % (module),
          'mv %s/coverage.out %s_coverage.out' % (module, module),
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
          'source': '%s_coverage.out' % (module),
          'target': '%s/%s/coverage' % (ctx.repo.slug, ctx.build.commit + '-${DRONE_BUILD_NUMBER}'),
          'path_style': True,
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
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
  }

def buildOcisBinaryForTesting(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'build_ocis_binary_for_testing',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      makeGenerate('ocis') +
      build() +
      rebuildBuildArtifactCache(ctx, 'ocis-binary-amd64', 'ocis/bin/ocis'),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
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
  }

def localApiTests(ctx, coreBranch = 'master', coreCommit = '', storage = 'owncloud', suite = 'apiOcisSpecific', accounts_hash_difficulty = 4):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'localApiTests-%s-%s' % (suite, storage),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      restoreBuildArtifactCache(ctx, 'ocis-binary-amd64', 'ocis/bin/ocis') +
      ocisServer(storage, accounts_hash_difficulty) +
      cloneCoreRepos(coreBranch, coreCommit) + [
      {
        'name': 'localApiTests-%s-%s' % (suite, storage),
        'image': 'owncloudci/php:7.4',
        'pull': 'always',
        'environment' : {
          'TEST_SERVER_URL': 'https://ocis-server:9200',
          'OCIS_REVA_DATA_ROOT': '%s' % ('/srv/app/tmp/ocis/owncloud/data/' if storage == 'owncloud' else ''),
          'DELETE_USER_DATA_CMD': '%s' % ('' if storage == 'owncloud' else 'rm -rf /srv/app/tmp/ocis/storage/users/nodes/root/* /srv/app/tmp/ocis/storage/users/nodes/*-*-*-*'),
          'SKELETON_DIR': '/srv/app/tmp/testing/data/apiSkeleton',
          'OCIS_SKELETON_STRATEGY': '%s' % ('copy' if storage == 'owncloud' else 'upload'),
          'TEST_OCIS':'true',
          'BEHAT_SUITE': suite,
          'BEHAT_FILTER_TAGS': '~@skipOnOcis-%s-Storage' % ('OC' if storage == 'owncloud' else 'OCIS'),
          'PATH_TO_CORE': '/srv/app/testrunner',
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
      redis(),
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'depends_on': getPipelineNames([buildOcisBinaryForTesting(ctx)]),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def coreApiTests(ctx, coreBranch = 'master', coreCommit = '', part_number = 1, number_of_parts = 1, storage = 'owncloud', accounts_hash_difficulty = 4):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'Core-API-Tests-%s-storage-%s' % (storage, part_number),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      restoreBuildArtifactCache(ctx, 'ocis-binary-amd64', 'ocis/bin/ocis') +
      ocisServer(storage, accounts_hash_difficulty) +
      cloneCoreRepos(coreBranch, coreCommit) + [
      {
        'name': 'oC10ApiTests-%s-storage-%s' % (storage, part_number),
        'image': 'owncloudci/php:7.4',
        'pull': 'always',
        'environment' : {
          'TEST_SERVER_URL': 'https://ocis-server:9200',
          'OCIS_REVA_DATA_ROOT': '%s' % ('/srv/app/tmp/ocis/owncloud/data/' if storage == 'owncloud' else ''),
          'DELETE_USER_DATA_CMD': '%s' % ('' if storage == 'owncloud' else 'rm -rf /srv/app/tmp/ocis/storage/users/nodes/root/* /srv/app/tmp/ocis/storage/users/nodes/*-*-*-*'),
          'SKELETON_DIR': '/srv/app/tmp/testing/data/apiSkeleton',
          'OCIS_SKELETON_STRATEGY': '%s' % ('copy' if storage == 'owncloud' else 'upload'),
          'TEST_OCIS':'true',
          'BEHAT_FILTER_TAGS': '~@notToImplementOnOCIS&&~@toImplementOnOCIS&&~comments-app-required&&~@federation-app-required&&~@notifications-app-required&&~systemtags-app-required&&~@local_storage&&~@skipOnOcis-%s-Storage' % ('OC' if storage == 'owncloud' else 'OCIS'),
          'DIVIDE_INTO_NUM_PARTS': number_of_parts,
          'RUN_PART': part_number,
          'EXPECTED_FAILURES_FILE': '/drone/src/tests/acceptance/expected-failures-on-%s-storage.txt' % (storage.upper()),
        },
        'commands': [
          'make -C /srv/app/testrunner test-acceptance-api',
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
    'depends_on': getPipelineNames([buildOcisBinaryForTesting(ctx)]),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def uiTests(ctx, phoenixBranch, phoenixCommit):
  suiteNames = config['uiTests']['suites'].keys()
  return [uiTestPipeline(ctx, suiteName, phoenixBranch, phoenixCommit) for suiteName in suiteNames]

def uiTestPipeline(ctx, suiteName, phoenixBranch = 'master', phoenixCommit = '', storage = 'owncloud', accounts_hash_difficulty = 4):
  suites = config['uiTests']['suites']
  paths = ""
  suite = suites[suiteName]
  if type(suite) == "list":
    for path in suite:
      paths = paths + "tests/acceptance/features/" + path + " "
  else:
    paths = paths + "tests/acceptance/features/" + suite + " "

  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': suiteName,
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      restoreBuildArtifactCache(ctx, 'ocis-binary-amd64', 'ocis/bin/ocis') +
      ocisServer(storage, accounts_hash_difficulty) + [
      {
        'name': 'webUITests',
        'image': 'owncloudci/nodejs:11',
        'pull': 'always',
        'environment': {
          'SERVER_HOST': 'https://ocis-server:9200',
          'BACKEND_HOST': 'https://ocis-server:9200',
          'RUN_ON_OCIS': 'true',
          'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/ocis/owncloud/data',
          'OCIS_SKELETON_DIR': '/srv/app/testing/data/webUISkeleton',
          'PHOENIX_CONFIG': '/drone/src/tests/config/drone/ocis-config.json',
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
    'depends_on': getPipelineNames([buildOcisBinaryForTesting(ctx)]),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def accountsUITests(ctx, phoenixBranch, phoenixCommit, storage = 'owncloud', accounts_hash_difficulty = 4):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'accountsUITests',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps':
      restoreBuildArtifactCache(ctx, 'ocis-binary-amd64', 'ocis/bin/ocis') +
      ocisServer(storage, accounts_hash_difficulty) + [
      {
        'name': 'WebUIAcceptanceTests',
        'image': 'owncloudci/nodejs:11',
        'pull': 'always',
        'environment': {
          'SERVER_HOST': 'https://ocis-server:9200',
          'BACKEND_HOST': 'https://ocis-server:9200',
          'RUN_ON_OCIS': 'true',
          'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/ocis/owncloud/data',
          'OCIS_SKELETON_DIR': '/srv/app/testing/data/webUISkeleton',
          'PHOENIX_CONFIG': '/drone/src/tests/config/drone/ocis-config.json',
          'TEST_TAGS': 'not @skipOnOCIS and not @skip',
          'LOCAL_UPLOAD_DIR': '/uploads',
          'NODE_TLS_REJECT_UNAUTHORIZED': 0,
          'PHOENIX_PATH': '/srv/app/phoenix',
          'FEATURE_PATH': '/drone/src/accounts/ui/tests/acceptance/features',
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
          'cd /drone/src/accounts',
          'yarn install --all',
          'make test-acceptance-webui'
        ],
        'volumes': [{
          'name': 'gopath',
          'path': '/srv/app',
        },
        {
          'name': 'uploads',
          'path': '/uploads'
        }],
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
    'depends_on': getPipelineNames([buildOcisBinaryForTesting(ctx)]),
    'trigger': {
      'ref': [
        'refs/heads/master',
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
      makeGenerate('ocis') +
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
    'depends_on': getPipelineNames(testOcisModules(ctx) + testPipelines(ctx)),
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
      makeGenerate('ocis') +
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
    'depends_on': getPipelineNames(testOcisModules(ctx) + testPipelines(ctx)),
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
      makeGenerate('ocis') + [
      {
        'name': 'build',
        'image': 'webhippie/golang:1.14',
        'pull': 'always',
        'commands': [
          'make -C ocis release-%s' % (name),
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
        'image': 'webhippie/golang:1.14',
        'pull': 'always',
        'commands': [
          'make -C ocis release-finish',
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
          'calens --version %s -o ocis/dist/CHANGELOG.md' % ctx.build.ref.replace("refs/tags/v", "").split("-")[0],
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
    'depends_on': getPipelineNames(testOcisModules(ctx) + testPipelines(ctx)),
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
        'image': 'webhippie/golang:1.14',
        'pull': 'always',
        'commands': [
          'make -C ocis changelog',
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
    generateConfigDocs.append('make -C %s config-docs-generate' % (module))

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
          'make -C docs docs-copy'
        ],
      },
      {
        'name': 'generate-config-docs',
        'image': 'webhippie/golang:1.14',
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
          'cd docs/hugo',
          'hugo',
        ],
      },
      {
        'name': 'list and remove temporary files',
        'image': 'owncloudci/alpine:latest',
        'commands': [
          'tree hugo/public',
          'rm -rf docs/hugo',
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

def makeGenerate(module):
  return [
    {
      'name': 'generate',
      'image': 'webhippie/golang:1.14',
      'pull': 'always',
      'commands': [
        'make -C %s generate' % (module),
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

def notify(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'chat-notifications',
    'clone': {
      'disable': True
    },
    'steps': [
      {
        'name': 'notify-rocketchat',
        'image': 'plugins/slack:1',
        'pull': 'always',
        'settings': {
          'webhook': {
            'from_secret': config['rocketchat']['from_secret']
          },
          'channel': config['rocketchat']['channel']
        },
        'when': {
          'status': [
            'failure',
          ],
        },
      },
    ],
    'depends_on': [],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/heads/release*',
        'refs/tags/**',
      ],
      'status': [
				'failure'
			]
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

def ocisServer(storage, accounts_hash_difficulty = 4):
  environment = {
    #'OCIS_LOG_LEVEL': 'debug',
    'STORAGE_HOME_DRIVER': '%s' % (storage),
    'STORAGE_USERS_DRIVER': '%s' % (storage),
    'STORAGE_DRIVER_OCIS_ROOT': '/srv/app/tmp/ocis/storage/users',
    'STORAGE_DRIVER_LOCAL_ROOT': '/srv/app/tmp/ocis/local/root',
    'STORAGE_METADATA_ROOT': '/srv/app/tmp/ocis/metadata',
    'STORAGE_DRIVER_OWNCLOUD_DATADIR': '/srv/app/tmp/ocis/owncloud/data',
    'STORAGE_DRIVER_OWNCLOUD_REDIS_ADDR': 'redis:6379',
    'STORAGE_LDAP_IDP': 'https://ocis-server:9200',
    'STORAGE_OIDC_ISSUER': 'https://ocis-server:9200',
    'PROXY_OIDC_ISSUER': 'https://ocis-server:9200',
    'STORAGE_HOME_DATA_SERVER_URL': 'http://ocis-server:9155/data',
    'STORAGE_DATAGATEWAY_PUBLIC_URL': 'https://ocis-server:9200/data',
    'STORAGE_USERS_DATA_SERVER_URL': 'http://ocis-server:9158/data',
    'STORAGE_FRONTEND_PUBLIC_URL': 'https://ocis-server:9200',
    'STORAGE_SHARING_USER_JSON_FILE': '/srv/app/tmp/ocis/shares.json',
    'PROXY_ENABLE_BASIC_AUTH': True,
    'PHOENIX_WEB_CONFIG': '/drone/src/tests/config/drone/ocis-config.json',
    'KONNECTD_IDENTIFIER_REGISTRATION_CONF': '/drone/src/tests/config/drone/identifier-registration.yml',
    'KONNECTD_ISS': 'https://ocis-server:9200',
    'KONNECTD_TLS': 'true',
  }

  # Pass in "default" accounts_hash_difficulty to not set this environment variable.
  # That will allow OCIS to use whatever its built-in default is.
  # Otherwise pass in a value from 4 to about 11 or 12 (default 4, for making regular tests fast)
  # The high values cause lots of CPU to be used when hashing passwords, and really slow down the tests.
  if (accounts_hash_difficulty != 'default'):
    environment['ACCOUNTS_HASH_DIFFICULTY'] = accounts_hash_difficulty

  return [
    {
      'name': 'ocis-server',
      'image': 'webhippie/golang:1.14',
      'pull': 'always',
      'detach': True,
      'environment' : environment,
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
      'image': 'owncloudci/php:7.4',
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
      'image': 'webhippie/golang:1.14',
      'pull': 'always',
      'commands': [
        'make -C ocis build',
      ],
      'volumes': [
        {
          'name': 'gopath',
          'path': '/srv/app',
        },
      ],
    },
  ]

def genericCache(name, action, mounts, cache_key):
  rebuild = 'false'
  restore = 'false'
  if action == 'rebuild':
    rebuild = 'true'
    action = 'rebuild'
  else:
    restore = 'true'
    action = 'restore'

  step = {
      'name': '%s_%s' %(action, name),
      'image': 'meltwater/drone-cache:v1',
      'pull': 'always',
      'environment': {
        'AWS_ACCESS_KEY_ID': {
          'from_secret': 'cache_s3_access_key',
        },
        'AWS_SECRET_ACCESS_KEY': {
          'from_secret': 'cache_s3_secret_key',
        },
      },
      'settings': {
        'endpoint': {
            'from_secret': 'cache_s3_endpoint'
          },
        'bucket': 'cache',
        'region': 'us-east-1', # not used at all, but failes if not given!
        'path_style': 'true',
        'cache_key': cache_key,
        'rebuild': rebuild,
        'restore': restore,
        'mount': mounts,
      },
    }
  return step

def genericCachePurge(ctx, name, cache_key):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'purge_%s' %(name),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'purge-cache',
        'image': 'minio/mc',
        'failure': 'ignore',
        'environment': {
          'MC_HOST_cache': {
            'from_secret': 'cache_s3_connection_url'
          }
        },
        'commands': [
          'mc rm --recursive --force %s/%s' % (ctx.repo.name, cache_key),
        ]
      },
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
      'status': [
        'success',
        'failure',
      ]
    },
  }

def genericBuildArtifactCache(ctx, name, action, path):
  name = '%s_build_artifact_cache' %(name)
  cache_key = '%s/%s/%s' % (ctx.repo.slug, ctx.build.commit + '-${DRONE_BUILD_NUMBER}', name)
  if action == "rebuild" or action == "restore":
    return genericCache(name, action, [path], cache_key)
  if action == "purge":
    return genericCachePurge(ctx, name, cache_key)
  return []

def restoreBuildArtifactCache(ctx, name, path):
  return [genericBuildArtifactCache(ctx, name, 'restore', path)]

def rebuildBuildArtifactCache(ctx, name, path):
  return [genericBuildArtifactCache(ctx, name, 'rebuild', path)]

def purgeBuildArtifactCache(ctx, name):
  return genericBuildArtifactCache(ctx, name, 'purge', [])
