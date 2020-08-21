def main(ctx):
  before = [
    testing(ctx),
    UITests(ctx, 'basic-roles-support', '', 'master', '0075dbd7b14360de203ec2a327c6b3f5c5844364')
  ]

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

def testing(ctx):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'testing',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'frontend',
        'image': 'webhippie/nodejs:latest',
        'pull': 'always',
        'commands': [
          'yarn install --frozen-lockfile',
          'yarn lint',
          'yarn test',
          'yarn build',
        ],
      },
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
      },
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
      {
        'name': 'sonarcloud',
        'image': 'sonarsource/sonar-scanner-cli',
        'pull': 'always',
        'environment': {
          'SONAR_TOKEN': {
            'from_secret': 'sonar_token',
          },
          'SONAR_PULL_REQUEST_BASE': 'master' if ctx.build.event == 'pull_request' else None,
          'SONAR_PULL_REQUEST_BRANCH': ctx.build.source if ctx.build.event == 'pull_request' else None,
          'SONAR_PULL_REQUEST_KEY': ctx.build.ref.replace("refs/pull/", "").split("/")[0] if ctx.build.event == 'pull_request' else None,
        },
      },
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
    'steps': [
      {
        'name': 'frontend',
        'image': 'webhippie/nodejs:latest',
        'pull': 'always',
        'commands': [
          'yarn install --frozen-lockfile',
          'yarn build',
        ],
      },
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
      },
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
      'testing',
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**',
      ],
    },
  }

def UITests(ctx, ocisBranch, ocisCommitId, phoenixBranch, phoenixCommitId):
  return {
   'kind': 'pipeline',
   'type': 'docker',
   'name': 'UiTests',
   'platform': {
     'os': 'linux',
     'arch': 'amd64',
    },
   'steps': [
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
     {
       'name': 'build-ocis',
       'image': 'webhippie/golang:1.13',
       'pull': 'always',
       'commands': [
         'git clone -b %s --single-branch --no-tags https://github.com/owncloud/ocis /srv/app/ocis' % (ocisBranch),
         'cd /srv/app/ocis',
         'git checkout %s' % (ocisCommitId),
         'make build',
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
         'REVA_STORAGE_HOME_DATA_TEMP_FOLDER': '/srv/app/tmp/',
         'REVA_STORAGE_LOCAL_ROOT': '/srv/app/tmp/reva/root',
         'REVA_STORAGE_OWNCLOUD_DATADIR': '/srv/app/tmp/reva/data',
         'REVA_STORAGE_OC_DATA_TEMP_FOLDER': '/srv/app/tmp/',
         'REVA_STORAGE_OWNCLOUD_REDIS_ADDR': 'redis:6379',
         'REVA_OIDC_ISSUER': 'https://ocis-server:9200',
         'PROXY_OIDC_ISSUER': 'https://ocis-server:9200',
         'REVA_STORAGE_OC_DATA_SERVER_URL': 'http://ocis-server:9164/data',
         'REVA_DATAGATEWAY_URL': 'https://ocis-server:9200/data',
         'REVA_FRONTEND_URL': 'https://ocis-server:9200',
         'REVA_LDAP_IDP': 'https://ocis-server:9200',
         'PHOENIX_WEB_CONFIG': '/drone/src/ui/tests/config/drone/ocis-config.json',
         'KONNECTD_IDENTIFIER_REGISTRATION_CONF': '/drone/src/ui/tests/config/drone/identifier-registration.yml',
         'KONNECTD_ISS': 'https://ocis-server:9200',
         'KONNECTD_TLS': 'true',
       },
       'commands': [
         'mkdir -p /srv/app/tmp/reva',
         # First run settings service because accounts need it to register the settings bundles
         '/srv/app/ocis/bin/ocis settings &',
         # Now start the accounts service
         'bin/ocis-accounts server &',
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
         'OCIS_REVA_DATA_ROOT': '/srv/app/tmp/reva',
         'OCIS_SKELETON_DIR': '/srv/app/testing/data/webUISkeleton',
         'PHOENIX_CONFIG': '/drone/src/ui/tests/config/drone/ocis-config.json',
         'TEST_TAGS': 'not @skipOnOCIS and not @skip',
         'LOCAL_UPLOAD_DIR': '/uploads',
         'PHOENIX_PATH': '/srv/app/phoenix',
         'FEATURE_PATH': 'ui/tests/acceptance/features',
         'NODE_TLS_REJECT_UNAUTHORIZED': '0'
       },
       'commands': [
         'git clone --depth=1 https://github.com/owncloud/testing.git /srv/app/testing',
         'git clone -b %s --single-branch https://github.com/owncloud/phoenix /srv/app/phoenix' % (phoenixBranch),
         'cd /srv/app/phoenix',
         'git checkout %s' % (phoenixCommitId),
         'cp -r /srv/app/phoenix/tests/acceptance/filesForUpload/* /uploads',
         'yarn install-all',
         'cd /drone/src',
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
       }]
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
       'volumes': [{
           'name': 'uploads',
           'path': '/uploads'
       }],
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
    'steps': [
      {
        'name': 'frontend',
        'image': 'webhippie/nodejs:latest',
        'pull': 'always',
        'commands': [
          'yarn install --frozen-lockfile',
          'yarn build',
        ],
      },
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
      },
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
      'testing',
    ],
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
