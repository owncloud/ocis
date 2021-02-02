config = {
  'modules': {
    'accounts': 'frontend',
    'glauth':'',
    'idp':'',
    'ocis': '',
    'web':'',
    'ocis-pkg':'',
    'ocs':'',
    'proxy':'',
    'settings':'frontend',
    'storage':'',
    'store':'',
    'thumbnails':'',
    'webdav':'',
    'onlyoffice':'frontend'
  },
  'apiTests': {
    'numberOfParts': 10
  },
  'uiTests': {
      'suites': {
        'webUIBasic': [
          'webUILogin',
          'webUINotifications',
          'webUIPrivateLinks',
          'webUIPreview',
          'webUIAccount',
          # The following suites may have all scenarios currently skipped.
          # The suites are listed here so that scenarios will run when
          # they are enabled.
          'webUIAdminSettings',
          'webUIComments',
          'webUITags',
          'webUIWebdavLockProtection',
          'webUIWebdavLocks',
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
        'webUISharingFolderAdvPermissionsGrp': 'webUISharingFolderAdvancedPermissionsGroups',
        'webUISharingPermissionToRoot': 'webUISharingPermissionToRoot',
        'webUIResharing': 'webUIResharing',
        'webUIResharingToRoot': 'webUIResharingToRoot',
        'webUISharingPublic': 'webUISharingPublic',
        'webUISharingPublicExpire': 'webUISharingPublicExpire',
        'webUISharingPublicDifferentRoles': 'webUISharingPublicDifferentRoles',
        'webUITrashbinDelete': 'webUITrashbinDelete',
        'webUITrashbinFilesFolders': 'webUITrashbinFilesFolders',
        'webUITrashbinRestore': 'webUITrashbinRestore',
        'webUIUpload': 'webUIUpload',
        'webUISharingFilePermissionMultipleUsers': 'webUISharingFilePermissionMultipleUsers',
        'webUISharingFolderPermissionMultipleUsers': 'webUISharingFolderPermissionMultipleUsers',
        'webUISharingFolderAdvancedPermissionMU': 'webUISharingFolderAdvancedPermissionMultipleUsers',
        'webUIMoveFilesFolders': 'webUIMoveFilesFolders',
      },
  },
  'rocketchat': {
    'channel': 'ocis-internal',
    'from_secret': 'private_rocketchat',
  },
  'binaryReleases': {
    'os': ['linux', 'darwin', 'windows'],
  },
  'dockerReleases': {
    'architectures': ['arm', 'arm64', 'amd64'],
  },
}


# volume for steps to cache Go dependencies between steps of a pipeline
# GOPATH must be set to /srv/app inside the image, which is the case for webhippie/golang
stepVolumeGoWebhippie = \
  {
  'name': 'gopath',
  'path': '/srv/app',
  }

# volume for pipeline to cache Go dependencies between steps of a pipeline
# to be used in combination with stepVolumeGoWebhippie
pipelineVolumeGoWebhippie = \
  {
  'name': 'gopath',
  'temp': {},
  }

stepVolumeOC10Tests = \
  {
  'name': 'oC10Tests',
  'path': '/srv/app',
  }

pipelineVolumeOC10Tests = \
  {
  'name': 'oC10Tests',
  'temp': {},
  }

def pipelineDependsOn(pipeline, dependant_pipelines):
  pipeline['depends_on'] = getPipelineNames(dependant_pipelines)
  return pipeline

def pipelinesDependsOn(pipelines, dependant_pipelines):
  pipes = []
  for pipeline in pipelines:
    pipeline['depends_on'] = getPipelineNames(dependant_pipelines)
    pipes.append(pipeline)

  return pipes


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

  test_pipelines = \
    [ buildOcisBinaryForTesting(ctx) ] + \
    testOcisModules(ctx) + \
    testPipelines(ctx)

  build_release_pipelines = \
    dockerReleases(ctx) + \
    [dockerEos(ctx)] + \
    binaryReleases(ctx) + \
    [releaseSubmodule(ctx)]

  build_release_helpers = [
    changelog(ctx),
    docs(ctx),
    refreshDockerBadges(ctx)
  ]

  if ctx.build.event == "cron":
    pipelines = test_pipelines + [
      pipelineDependsOn(
        purgeBuildArtifactCache(ctx, 'ocis-binary-amd64'),
        testPipelines(ctx)
      )
    ] + example_deploys(ctx)

  elif \
  (ctx.build.event == "pull_request" and '[docs-only]' in ctx.build.title) \
  or \
  (ctx.build.event != "pull_request" and '[docs-only]' in (ctx.build.title + ctx.build.message)):
  # [docs-only] is not taken from PR messages, but from commit messages
    pipelines = [docs(ctx)]

  else:
    test_pipelines.append(
      pipelineDependsOn(
        purgeBuildArtifactCache(ctx, 'ocis-binary-amd64'),
        testPipelines(ctx)
      )
    )

    pipelines = test_pipelines + build_release_pipelines + build_release_helpers


    pipelines = \
      pipelines + \
      pipelinesDependsOn(
        example_deploys(ctx),
        pipelines
      )

  # always append notification step
  pipelines.append(
    pipelineDependsOn(
      notify(ctx),
      pipelines
    )
  )

  pipelineSanityChecks(ctx, pipelines)
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
    localApiTests(ctx, 'owncloud', 'apiBugDemonstration'),
    localApiTests(ctx, 'ocis', 'apiBugDemonstration'),
    localApiTests(ctx, 'owncloud', 'apiAccountsHashDifficulty', 'default'),
    localApiTests(ctx, 'ocis', 'apiAccountsHashDifficulty', 'default')
  ]

  for runPart in range(1, config['apiTests']['numberOfParts'] + 1):
    pipelines.append(coreApiTests(ctx, runPart, config['apiTests']['numberOfParts'], 'owncloud'))
    pipelines.append(coreApiTests(ctx, runPart, config['apiTests']['numberOfParts'], 'ocis'))

  pipelines += uiTests(ctx)
  pipelines.append(accountsUITests(ctx))
  return pipelines

def testOcisModule(ctx, module):
  steps = makeGenerate(module) + [
    {
      'name': 'vet',
      'image': 'webhippie/golang:1.15',
      'pull': 'always',
      'commands': [
        'make -C %s vet' % (module),
      ],
      'volumes': [stepVolumeGoWebhippie,],
    },
    {
      'name': 'staticcheck',
      'image': 'webhippie/golang:1.15',
      'pull': 'always',
      'commands': [
        'make -C %s staticcheck' % (module),
      ],
      'volumes': [stepVolumeGoWebhippie,],
    },
    {
      'name': 'lint',
      'image': 'webhippie/golang:1.15',
      'pull': 'always',
      'commands': [
        'make -C %s lint' % (module),
      ],
      'volumes': [stepVolumeGoWebhippie,],
    },
    {
        'name': 'test',
        'image': 'webhippie/golang:1.15',
        'pull': 'always',
        'commands': [
          'make -C %s test' % (module),
          'mv %s/coverage.out %s_coverage.out' % (module, module),
        ],
        'volumes': [stepVolumeGoWebhippie,],
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
    'volumes': [pipelineVolumeGoWebhippie],
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
    'volumes': [pipelineVolumeGoWebhippie],
  }

def uploadCoverage(ctx):
  sonar_env = {
      'SONAR_TOKEN': {
        'from_secret': 'sonar_token',
      },
    }
  if ctx.build.event == "pull_request":
    sonar_env.update({
      'SONAR_PULL_REQUEST_BASE': '%s' % (ctx.build.target),
      'SONAR_PULL_REQUEST_BRANCH': '%s' % (ctx.build.source),
      'SONAR_PULL_REQUEST_KEY': '%s' % (ctx.build.ref.replace("refs/pull/", "").split("/")[0]),
    })

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
        'image': 'minio/mc:RELEASE.2020-12-10T01-26-17Z',
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
        'environment': sonar_env,
      },
      {
        'name': 'purge-cache',
        'image': 'minio/mc:RELEASE.2020-12-10T01-26-17Z',
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

def localApiTests(ctx, storage = 'owncloud', suite = 'apiBugDemonstration', accounts_hash_difficulty = 4):
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
      ocisServer(storage, accounts_hash_difficulty, [stepVolumeOC10Tests]) +
      cloneCoreRepos() + [
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
        'volumes': [stepVolumeOC10Tests],
      },
    ],
    'services':
      redis(),
    'depends_on': getPipelineNames([buildOcisBinaryForTesting(ctx)]),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
    'volumes': [pipelineVolumeOC10Tests],
  }

def coreApiTests(ctx, part_number = 1, number_of_parts = 1, storage = 'owncloud', accounts_hash_difficulty = 4):
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
      ocisServer(storage, accounts_hash_difficulty, [stepVolumeOC10Tests]) +
      cloneCoreRepos() + [
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
          'EXPECTED_FAILURES_FILE': '/drone/src/tests/acceptance/expected-failures-API-on-%s-storage.md' % (storage.upper()),
        },
        'commands': [
          'make -C /srv/app/testrunner test-acceptance-api',
        ],
        'volumes': [stepVolumeOC10Tests],
      },
    ],
    'services':
      redis(),
    'depends_on': getPipelineNames([buildOcisBinaryForTesting(ctx)]),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
    'volumes': [pipelineVolumeOC10Tests],
  }

def uiTests(ctx):
  suiteNames = config['uiTests']['suites'].keys()
  return [uiTestPipeline(ctx, suiteName) for suiteName in suiteNames]

def uiTestPipeline(ctx, suiteName, storage = 'owncloud', accounts_hash_difficulty = 4):
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
      ocisServer(storage, accounts_hash_difficulty, [stepVolumeOC10Tests]) + [
      {
        'name': 'webUITests',
        'image': 'webhippie/nodejs:latest',
        'pull': 'always',
        'environment': {
          'SERVER_HOST': 'https://ocis-server:9200',
          'BACKEND_HOST': 'https://ocis-server:9200',
          'RUN_ON_OCIS': 'true',
          'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/ocis/owncloud/data',
          'OCIS_SKELETON_DIR': '/srv/app/testing/data/webUISkeleton',
          'WEB_UI_CONFIG': '/drone/src/tests/config/drone/ocis-config.json',
          'TEST_TAGS': 'not @skipOnOCIS and not @skip and not @notToImplementOnOCIS',
          'LOCAL_UPLOAD_DIR': '/uploads',
          'NODE_TLS_REJECT_UNAUTHORIZED': 0,
          'TEST_PATHS': paths,
          'EXPECTED_FAILURES_FILE': '/drone/src/tests/acceptance/expected-failures-webUI-on-%s-storage.md' % (storage.upper()),
        },
        'commands': [
          'source /drone/src/.drone.env',
          'git clone -b master --depth=1 https://github.com/owncloud/testing.git /srv/app/testing',
          'git clone -b $WEB_BRANCH --single-branch --no-tags https://github.com/owncloud/web.git /srv/app/web',
          'cd /srv/app/web',
          'git checkout $WEB_COMMITID',
          'cp -r tests/acceptance/filesForUpload/* /uploads',
          'yarn install-all',
          './tests/acceptance/run.sh'
        ],
        'volumes':
          [stepVolumeOC10Tests] +
          [{
            'name': 'uploads',
            'path': '/uploads'
          }]
      },
    ],
    'services':
      redis() +
      selenium(),
    'volumes':
      [pipelineVolumeOC10Tests] +
      [{
        'name': 'uploads',
        'temp': {}
      }],
    'depends_on': getPipelineNames([buildOcisBinaryForTesting(ctx)]),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def accountsUITests(ctx, storage = 'owncloud', accounts_hash_difficulty = 4):
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
      ocisServer(storage, accounts_hash_difficulty, [stepVolumeOC10Tests]) + [
      {
        'name': 'WebUIAcceptanceTests',
        'image': 'webhippie/nodejs:latest',
        'pull': 'always',
        'environment': {
          'SERVER_HOST': 'https://ocis-server:9200',
          'BACKEND_HOST': 'https://ocis-server:9200',
          'RUN_ON_OCIS': 'true',
          'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/ocis/owncloud/data',
          'OCIS_SKELETON_DIR': '/srv/app/testing/data/webUISkeleton',
          'WEB_UI_CONFIG': '/drone/src/tests/config/drone/ocis-config.json',
          'TEST_TAGS': 'not @skipOnOCIS and not @skip',
          'LOCAL_UPLOAD_DIR': '/uploads',
          'NODE_TLS_REJECT_UNAUTHORIZED': 0,
          'WEB_PATH': '/srv/app/web',
          'FEATURE_PATH': '/drone/src/accounts/ui/tests/acceptance/features',
        },
        'commands': [
          'source /drone/src/.drone.env',
          'git clone -b master --depth=1 https://github.com/owncloud/testing.git /srv/app/testing',
          'git clone -b $WEB_BRANCH --single-branch --no-tags https://github.com/owncloud/web.git /srv/app/web',
          'cd /srv/app/web',
          'git checkout $WEB_COMMITID',
          'cp -r tests/acceptance/filesForUpload/* /uploads',
          'yarn install-all',
          'cd /drone/src/accounts',
          'yarn install --all',
          'make test-acceptance-webui'
        ],
        'volumes':
          [stepVolumeOC10Tests] +
          [{
            'name': 'uploads',
            'path': '/uploads'
          }]
      },
    ],
    'services':
      redis() +
      selenium(),
    'volumes':
      [stepVolumeOC10Tests] +
      [{
        'name': 'uploads',
        'temp': {}
      }],
    'depends_on': getPipelineNames([buildOcisBinaryForTesting(ctx)]),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
  }

def dockerReleases(ctx):
  pipelines = []
  for arch in config['dockerReleases']['architectures']:
    pipelines.append(dockerRelease(ctx, arch))

  manifest = releaseDockerManifest(ctx)
  manifest['depends_on'] = getPipelineNames(pipelines)
  pipelines.append(manifest)

  readme = releaseDockerReadme(ctx)
  readme['depends_on'] = getPipelineNames(pipelines)
  pipelines.append(readme)

  return pipelines

def dockerRelease(ctx, arch):
  build_args = [
    'REVISION=%s' % (ctx.build.commit),
    'VERSION=%s' % (ctx.build.ref.replace("refs/tags/", "") if ctx.build.event == "tag" else "latest")
  ]

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
          'build_args': build_args,
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
          'repo': ctx.build.commit,
          'build_args': build_args,
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
    'depends_on': getPipelineNames(testOcisModules(ctx) + testPipelines(ctx)),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
    'volumes': [pipelineVolumeGoWebhippie],
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
    'depends_on': getPipelineNames(testOcisModules(ctx) + testPipelines(ctx)),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
    'volumes': [pipelineVolumeGoWebhippie],
  }

def binaryReleases(ctx):
  pipelines = []
  for os in config['binaryReleases']['os']:
    pipelines.append(binaryRelease(ctx, os))

  return pipelines

def binaryRelease(ctx, name):
  # uploads binary to https://download.owncloud.com/ocis/ocis/testing/
  target = '/ocis/%s/testing' % (ctx.repo.name.replace("ocis-", ""))
  if ctx.build.event == "tag":
    # uploads binary to eg. https://download.owncloud.com/ocis/ocis/1.0.0-beta9/
    target = '/ocis/%s/%s' % (ctx.repo.name.replace("ocis-", ""), ctx.build.ref.replace("refs/tags/v", ""))

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
    'target': target,
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
        'image': 'webhippie/golang:1.15',
        'pull': 'always',
        'commands': [
          'make -C ocis release-%s' % (name),
        ],
      },
      {
        'name': 'finish',
        'image': 'webhippie/golang:1.15',
        'pull': 'always',
        'commands': [
          'make -C ocis release-finish',
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
    'depends_on': getPipelineNames(testOcisModules(ctx) + testPipelines(ctx)),
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
        'refs/pull/**',
      ],
    },
    'volumes': [pipelineVolumeGoWebhippie],
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
    'depends_on': depends,
    'trigger': {
      'ref': [
        'refs/tags/*/v*',
      ],
    },
  }


def releaseDockerManifest(ctx):
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
        'image': 'webhippie/golang:1.15',
        'pull': 'always',
        'commands': [
          'make -C ocis changelog',
        ],
      },
      {
        'name': 'diff',
        'image': 'owncloudci/alpine:latest',
        'pull': 'always',
        'commands': [
          'git diff',
        ],
      },
      {
        'name': 'output',
        'image': 'owncloudci/alpine:latest',
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
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/pull/**',
      ],
    },
  }

def releaseDockerReadme(ctx):
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
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
      ],
    },
  }

def refreshDockerBadges(ctx):
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
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
      ],
    },
    'depends_on': getPipelineNames(dockerReleases(ctx)),
  }

def docs(ctx):
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
        'name': 'generate-config-docs',
        'image': 'webhippie/golang:1.15',
        'commands': ['make -C %s config-docs-generate' % (module) for module in config['modules']],
      },
      {
        'name': 'prepare',
        'image': 'owncloudci/alpine:latest',
        'commands': [
          'make -C docs docs-copy'
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
          'pages_directory': 'docs/hugo/content',
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
        'name': 'list and remove temporary files',
        'image': 'owncloudci/alpine:latest',
        'commands': [
          'tree hugo/public',
          'rm -rf docs/hugo',
        ],
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
      'image': 'webhippie/golang:1.15',
      'pull': 'always',
      'commands': [
        'make -C %s generate' % (module),
      ],
      'volumes': [stepVolumeGoWebhippie,],
    }
  ]

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

def ocisServer(storage, accounts_hash_difficulty = 4, volumes=[]):
  environment = {
    #'OCIS_LOG_LEVEL': 'debug',
    'OCIS_URL': 'https://ocis-server:9200',
    'STORAGE_HOME_DRIVER': '%s' % (storage),
    'STORAGE_USERS_DRIVER': '%s' % (storage),
    'STORAGE_DRIVER_OCIS_ROOT': '/srv/app/tmp/ocis/storage/users',
    'STORAGE_DRIVER_LOCAL_ROOT': '/srv/app/tmp/ocis/local/root',
    'STORAGE_METADATA_ROOT': '/srv/app/tmp/ocis/metadata',
    'STORAGE_DRIVER_OWNCLOUD_DATADIR': '/srv/app/tmp/ocis/owncloud/data',
    'STORAGE_DRIVER_OWNCLOUD_REDIS_ADDR': 'redis:6379',
    'STORAGE_HOME_DATA_SERVER_URL': 'http://ocis-server:9155/data',
    'STORAGE_USERS_DATA_SERVER_URL': 'http://ocis-server:9158/data',
    'STORAGE_SHARING_USER_JSON_FILE': '/srv/app/tmp/ocis/shares.json',
    'PROXY_ENABLE_BASIC_AUTH': True,
    'WEB_UI_CONFIG': '/drone/src/tests/config/drone/ocis-config.json',
    'IDP_IDENTIFIER_REGISTRATION_CONF': '/drone/src/tests/config/drone/identifier-registration.yml',
    'IDP_TLS': 'true',
    'OCIS_LOG_LEVEL': 'warn',
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
      'image': 'owncloudci/alpine:latest',
      'pull': 'always',
      'detach': True,
      'environment' : environment,
      'commands': [
        'apk add mailcap', # install /etc/mime.types
        'ocis/bin/ocis server'
      ],
      'volumes': volumes,
    },
    {
      'name': 'wait-for-ocis-server',
      'image': 'thegeeklab/wait-for:latest',
      'pull': 'always',
      'commands': [
        'wait-for ocis-server:9200 -t 300',
      ],
    },
  ]

def cloneCoreRepos():
  return [
    {
      'name': 'clone-core-repos',
      'image': 'owncloudci/alpine:latest',
      'pull': 'always',
      'commands': [
        'source /drone/src/.drone.env',
        'git clone -b master --depth=1 https://github.com/owncloud/testing.git /srv/app/tmp/testing',
        'git clone -b $CORE_BRANCH --single-branch --no-tags https://github.com/owncloud/core.git /srv/app/testrunner',
        'cd /srv/app/testrunner',
        'git checkout $CORE_COMMITID'
      ],
      'volumes': [stepVolumeOC10Tests],
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
      'image': 'webhippie/golang:1.15',
      'pull': 'always',
      'commands': [
        'make -C ocis build',
      ],
      'volumes': [stepVolumeGoWebhippie,],
    },
  ]

def example_deploys(ctx):
  latest_configs = [
    'cs3_users_ocis/latest.yml',
    'ocis_keycloak/latest.yml',
    'ocis_traefik/latest.yml',
  ]
  released_configs = [
    'cs3_users_ocis/released.yml',
    'ocis_keycloak/released.yml',
    'ocis_traefik/released.yml',
  ]

  # if on master branch:
  configs = latest_configs
  rebuild = "false"

  if ctx.build.event == "tag":
    configs = released_configs
    rebuild = 'false'

  if ctx.build.event == "cron":
    configs = latest_configs + released_configs
    rebuild = 'true'

  deploys = []
  for config in configs:
    deploys.append(deploy(ctx, config, rebuild))

  return deploys

def deploy(ctx, config, rebuild):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'deploy_%s' % (config),
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'clone continuous deployment playbook',
        'image': 'alpine/git',
        'commands': [
          'cd deployments/continuous-deployment-config',
          'git clone https://github.com/owncloud-devops/continuous-deployment.git',
        ]
      },
      {
        'name': 'deploy',
        'image': 'owncloudci/drone-ansible',
        'failure': 'ignore',
        'environment': {
          'CONTINUOUS_DEPLOY_SERVERS_CONFIG': '../%s' % (config),
          "REBUILD": '%s' % (rebuild),
          'HCLOUD_API_TOKEN': {
            'from_secret': 'hcloud_api_token'
          },
          'CLOUDFLARE_API_TOKEN': {
            'from_secret': 'cloudflare_api_token'
          }
        },
        'settings': {
          'playbook': 'deployments/continuous-deployment-config/continuous-deployment/playbook-all.yml',
          'galaxy': 'deployments/continuous-deployment-config/continuous-deployment/requirements.yml',
          'requirements': 'deployments/continuous-deployment-config/continuous-deployment/py-requirements.txt',
          'inventory': 'localhost',
          'private_key': {
            'from_secret': 'ssh_private_key'
          }
        }
      },
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/v*',
      ],
    },
  }

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
        'region': 'us-east-1', # not used at all, but fails if not given!
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
        'image': 'minio/mc:RELEASE.2020-12-10T01-26-17Z',
        'failure': 'ignore',
        'environment': {
          'MC_HOST_cache': {
            'from_secret': 'cache_s3_connection_url'
          }
        },
        'commands': [
          'mc rm --recursive --force cache/cache/%s/%s' % (ctx.repo.name, cache_key),
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

def pipelineSanityChecks(ctx, pipelines):
  """pipelineSanityChecks helps the CI developers to find errors before running it

  These sanity checks are only executed on when converting starlark to yaml.
  Error outputs are only visible when the conversion is done with the drone cli.

  Args:
    ctx: drone passes a context with information which the pipeline can be adapted to
    pipelines: pipelines to be checked, normally you should run this on the return value of main()

  Returns:
    none
  """

  # check if name length of pipeline and steps are exceeded.
  max_name_length = 50
  for pipeline in pipelines:
    pipeline_name = pipeline["name"]
    if len(pipeline_name) > max_name_length:
      print("Error: pipeline name %s is longer than 50 characters" %(pipeline_name))

    for step in pipeline["steps"]:
      step_name = step["name"]
      if len(step_name) > max_name_length:
        print("Error: step name %s in pipeline %s is longer than 50 characters" %(step_name, pipeline_name))

  # check for non existing depends_on
  possible_depends = []
  for pipeline in pipelines:
    possible_depends.append(pipeline["name"])

  for pipeline in pipelines:
    if "depends_on" in pipeline.keys():
      for depends in pipeline["depends_on"]:
        if not depends in possible_depends:
          print("Error: depends_on %s for pipeline %s is not defined" %(depends, pipeline["name"]))

  # check for non declared volumes
  for pipeline in pipelines:
    pipeline_volumes = []
    if "volumes" in pipeline.keys():
      for volume in pipeline['volumes']:
        pipeline_volumes.append(volume['name'])

    for step in pipeline["steps"]:
      if "volumes" in step.keys():
        for volume in step["volumes"]:
          if not volume['name'] in pipeline_volumes:
            print("Warning: volume %s for step %s is not defined in pipeline %s" %(volume['name'], step['name'], pipeline['name']))

  # list used docker images
  print("")
  print("List of used docker images:")

  images = {}

  for pipeline in pipelines:
    for step in pipeline['steps']:
      image = step["image"]
      if image in images.keys():
        images[image] = images[image] + 1
      else:
        images[image] = 1

  for image in images.keys():
    print(" %sx\t%s" %(images[image], image))
