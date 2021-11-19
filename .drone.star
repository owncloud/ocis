# images
OCIS_IMG = "owncloud/ocis:latest"
OC10_IMG = "owncloud/server:latest"
OC10_DB_IMG = "mariadb:10.6"
OPENLDAP_IMG = "osixia/openldap:latest"
KEYCLOAK_IMG = "quay.io/keycloak/keycloak:latest"
KEYCLOAK_DB_IMG = "postgres:alpine"
REDIS_IMG = "redis:6"
OC_CI_PHP = "owncloudci/php:7.4"
OC_CI_UBUNTU = "owncloud/ubuntu:18.04"
OC_CI_WAITFOR = "owncloudci/wait-for:latest"
OC_CI_GOLANG = "owncloudci/golang:1.17"
OC_CI_NODEJS = "owncloudci/nodejs:14"
OC_CI_ALPINE = "owncloudci/alpine:latest"

# Settings
OCIS_URL = "https://ocis:9200"
OCIS_DOMAIN = "ocis:9200"
OC10_URL = "http://oc10:8080"

DRONE_CONFIG_PATH = "/drone/src/tests/parallelDeployAcceptance/drone"

def main(ctx):
    pipelines = []
    pipelines = acceptancePipeline()
    return pipelines

def acceptancePipeline():
    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "Parallel-Deploy-API-Tests",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps":
            composerInstall() +
            copyConfigs() +
            # makeNodeGenerate() +
            # makeGoGenerate() +
            # buildOCIS() + 
            waitForServices() + 
            oC10Service() + 
            waitForOC10() +
            owncloudLog() +
            fixPermissions() +
            ocisServer() + 
            waitForOCIS() + 
            apiTests(),
        "services": keycloakDbService() +
                    oc10DbService() +
                    ldapService() +
                    redisService() +
                    keycloakService(),
        "volumes": [
            {
                "name": "config-templates",
                "temp": {},
            },
            {
                "name": "preserver-config",
                "temp": {},
            },
            {
                "name": "crons",
                "temp": {},
            },
            {
                "name": "data",
                "temp": {},
            },
            {
                "name": "core-apps",
                "temp": {},
            },
            {
                "name": "proxy-config",
                "temp": {},
            },
            {
                "name": "gopath",
                "temp": {},
            },
        ],
        "trigger": {
            "ref": [
                "refs/pull/**",
            ]
        }
    }

def apiTests():
    return [{
        "name": "API Tests",
        "image": OC_CI_PHP,
        "environment": {
            "TEST_SERVER_URL": OCIS_URL,
            "TEST_OC10_URL": OC10_URL,
            "PARALLEL_DEPLOY": "true",
            "TEST_OCIS": "true",
            "TEST_WITH_LDAP": "true",
            "REVA_LDAP_PORT" : 636,
            "REVA_LDAP_BASE_DN": "dc=owncloud,dc=com",
            "REVA_LDAP_HOSTNAME": "openldap",
            "REVA_LDAP_BIND_DN": "cn=admin,dc=owncloud,dc=com",
            "SKELETON_DIR": "/var/www/owncloud/apps/testing/data/apiSkeleton",
        },
        "commands": [
            "make -C ./tests/parallelDeployAcceptance test-paralleldeployment-api",
        ],
        "depends_on": ["composer-install", "wait-for-oc10", "wait-for-ocis"],
        "volumes": [
            {
                "name": "core-apps",
                "path": "/var/www/owncloud/apps",
            },
        ]
    }]

def makeNodeGenerate():
    return [{
        "name": "generate-nodejs",
        "image": OC_CI_NODEJS,
        "commands": [
            "make ci-node-generate",
        ],
        "volumes": [
            {
                "name": "gopath",
                "path": "/go",
            },
        ],
    }]

def makeGoGenerate():
    return [{
        "name": "generate-go",
        "image": OC_CI_GOLANG,
        "commands": [
            "whoami",
            "make ci-go-generate",
        ],
        "volumes": [
            {
                "name": "gopath",
                "path": "/go",
            },
        ],
    }]

def buildOCIS():
    return [{
        "name": "build-ocis",
        "image": OC_CI_GOLANG,
        "commands": [
            "make -C ocis build",
        ],
        "volumes": [
            {
                "name": "gopath",
                "path": "/go",
            },
        ],
        "depends_on": ["generate-nodejs", "generate-go"]
    }]

def ocisServer():
    environment = {
        "PROXY_ENABLE_BASIC_AUTH": "true",
        # Keycloak IDP specific configuration
        "PROXY_OIDC_ISSUER": "https://keycloak/auth/realmsowncloud",
        "WEB_OIDC_AUTHORITY": "https://keycloak/auth/realms/owncloud",
        "WEB_OIDC_CLIENT_ID": "ocis-web",
        "WEB_OIDC_METADATA_URL": "https://keycloak/auth/realms/owncloud/.well-known/openid-configuration",
        "STORAGE_OIDC_ISSUER": "https://keycloak",
        "STORAGE_LDAP_IDP": "https://keycloak/auth/realms/owncloud",
        "WEB_OIDC_SCOPE": "openid profile email owncloud",
        # LDAP bind
        "STORAGE_LDAP_HOSTNAME": "openldap",
        "STORAGE_LDAP_PORT": 636,
        "STORAGE_LDAP_INSECURE": "true",
        "STORAGE_LDAP_BIND_DN": "cn=admin,dc=owncloud,dc=com",
        "STORAGE_LDAP_BIND_PASSWORD": "admin",
        # LDAP user settings
        "PROXY_AUTOPROVISION_ACCOUNTS": "true", # automatically create users when they login
        "PROXY_ACCOUNT_BACKEND_TYPE": "cs3", # proxy should get users from CS3APIS (which gets it from LDAP)
        "PROXY_USER_OIDC_CLAIM": "ocis.user.uuid", # claim was added in Keycloak
        "PROXY_USER_CS3_CLAIM": "userid", # equals STORAGE_LDAP_USER_SCHEMA_UID
        "STORAGE_LDAP_BASE_DN": "dc=owncloud,dc=com",
        "STORAGE_LDAP_GROUP_SCHEMA_DISPLAYNAME": "cn",
        "STORAGE_LDAP_GROUP_SCHEMA_GID_NUMBER": "gidnumber",
        "STORAGE_LDAP_GROUP_SCHEMA_GID": "cn",
        "STORAGE_LDAP_GROUP_SCHEMA_MAIL": "mail",
        "STORAGE_LDAP_GROUPATTRIBUTEFILTER": "(&(objectclass=posixGroup)(objectclass=owncloud)({{attr}}={{value}}))",
        "STORAGE_LDAP_GROUPFILTER": "(&(objectclass=groupOfUniqueNames)(objectclass=owncloud)(ownclouduuid={{.OpaqueId}}*))",
        "STORAGE_LDAP_GROUPMEMBERFILTER": "(&(objectclass=posixAccount)(objectclass=owncloud)(ownclouduuid={{.OpaqueId}}*))",
        "STORAGE_LDAP_USERGROUPFILTER": "(&(objectclass=posixGroup)(objectclass=owncloud)(ownclouduuid={{.OpaqueId}}*))",
        "STORAGE_LDAP_USER_SCHEMA_CN": "cn",
        "STORAGE_LDAP_USER_SCHEMA_DISPLAYNAME": "displayname",
        "STORAGE_LDAP_USER_SCHEMA_GID_NUMBER": "gidnumber",
        "STORAGE_LDAP_USER_SCHEMA_MAIL": "mail",
        "STORAGE_LDAP_USER_SCHEMA_UID_NUMBER": "uidnumber",
        "STORAGE_LDAP_USER_SCHEMA_UID": "ownclouduuid",
        "STORAGE_LDAP_LOGINFILTER": "(&(objectclass=posixAccount)(objectclass=owncloud)(|(uid={{login}})(mail={{login}})))",
        "STORAGE_LDAP_USERATTRIBUTEFILTER": "(&(objectclass=posixAccount)(objectclass=owncloud)({{attr}}={{value}}))",
        "STORAGE_LDAP_USERFILTER": "(&(objectclass=posixAccount)(objectclass=owncloud)(|(ownclouduuid={{.OpaqueId}})(uid={{.OpaqueId}})))",
        "STORAGE_LDAP_USERFINDFILTER": "(&(objectclass=posixAccount)(objectclass=owncloud)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))",
        # ownCloud storage driver
        "STORAGE_HOME_DRIVER": "owncloudsql",
        "STORAGE_USERS_DRIVER": "owncloudsql",
        "STORAGE_METADATA_DRIVER": "ocis",
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_DATADIR": "/mnt/data/files",
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_UPLOADINFO_DIR": "/tmp",
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_SHARE_FOLDER": "/Shares",
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_LAYOUT": "{{.Username}}",
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBUSERNAME": "owncloud",
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBPASSWORD": "owncloud",
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBHOST": "oc10-db",
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBPORT": 3306,
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_DBNAME": "owncloud",
        # TODO: redis is not yet supported
        "STORAGE_USERS_DRIVER_OWNCLOUDSQL_REDIS_ADDR": "redis:6379",
        # ownCloud storage readonly
        # TODO: conflict with OWNCLOUDSQL -> https://github.com/owncloud/ocis/issues/2303
        "OCIS_STORAGE_READ_ONLY": "false",
        # General oCIS config
        "OCIS_LOG_LEVEL": "error",
        "OCIS_URL": OCIS_URL,
        "PROXY_TLS": "true",
        # change default secrets
        "OCIS_JWT_SECRET": "Pive-Fumkiu4",
        "STORAGE_TRANSFER_SECRET": "replace-me-with-a-transfer-secret",
        "OCIS_MACHINE_AUTH_API_KEY": "change-me-please",
        "OCIS_INSECURE": "true",
    }

    return [{
        "name": "ocis",
        "image": OCIS_IMG,
        "environment": environment,
        "detach": True,
        "commands": [
            "whoami",
            "cd /mnt/data",
            "ls -al",
            "%s/ocis/server.sh" % (DRONE_CONFIG_PATH),
        ],
        "volumes": [
            {
                "name": "data",
                "path": "/mnt/data",
            },
            {
                "name": "proxy-config",
                "path": "/etc/ocis",
            },
            {
                "name": "gopath",
                "path": "/go",
            },
        ],
        "user": "33:33",
        "depends_on": ["fix-permissions"],
    }]

def oC10Service():
    return [{
        "name": "oc10",
        "image": OC10_IMG,
        "pull": "always",
        "detach": True,
        "environment": {
            # can be switched to "web"
            "OWNCLOUD_DEFAULT_APP": "files",
            "OWNCLOUD_WEB_REWRITE_LINKS": "false",
            # script / config variables
            "IDP_OIDC_ISSUER": "https://keycloak/auth/realms/owncloud",
            "IDP_OIDC_CLIENT_SECRET": "oc10-oidc-secret",
            "CLOUD_DOMAIN": OCIS_DOMAIN,
            # LDAP bind configuration
            "LDAP_HOST": "openldap",
            "LDAP_PORT": 389,
            "STORAGE_LDAP_BIND_DN": "cn=admin,dc=owncloud,dc=com",
            "STORAGE_LDAP_BIND_PASSWORD": "admin",
            # LDAP user configuration 
            "LDAP_BASE_DN": "dc=owncloud,dc=com",
            "LDAP_USER_SCHEMA_DISPLAYNAME": "displayname",
            "LDAP_LOGINFILTER": "(&(objectclass=owncloud)(|(uid=%uid)(mail=%uid)))",
            "LDAP_GROUP_SCHEMA_DISPLAYNAME": "cn",
            "LDAP_USER_SCHEMA_NAME_ATTR": "uid",
            "LDAP_GROUPFILTER": "(&(objectclass=groupOfUniqueNames)(objectclass=owncloud))",
            "LDAP_USER_SCHEMA_UID": "ownclouduuid",
            "LDAP_USERATTRIBUTEFILTERS": "uid", # ownCloudUUID;cn;uid;mail
            "LDAP_USER_SCHEMA_MAIL": "mail",
            "LDAP_USERFILTER": "(&(objectclass=owncloud))",
            "LDAP_GROUP_MEMBER_ASSOC_ATTR": "uniqueMember",
            # database
            "OWNCLOUD_DB_TYPE": "mysql",
            "OWNCLOUD_DB_NAME": "owncloud",
            "OWNCLOUD_DB_USERNAME": "owncloud",
            "OWNCLOUD_DB_PASSWORD": "owncloud",
            "OWNCLOUD_DB_HOST": "oc10-db",
            "OWNCLOUD_ADMIN_USERNAME": "admin",
            "OWNCLOUD_ADMIN_PASSWORD": "admin",
            "OWNCLOUD_MYSQL_UTF8MB4": "true",
            # redis
            "OWNCLOUD_REDIS_ENABLED": "true",
            "OWNCLOUD_REDIS_HOST": "redis",
            # ownCloud config
            "OWNCLOUD_TRUSTED_PROXIES": OCIS_DOMAIN,
            "OWNCLOUD_OVERWRITE_PROTOCOL": "https",
            "OWNCLOUD_OVERWRITE_HOST": OCIS_DOMAIN,
            "OWNCLOUD_APPS_ENABLE": "openidconnect,oauth2,user_ldap,graphapi",
            "OWNCLOUD_LOG_LEVEL": 2,
            "OWNCLOUD_LOG_FILE": "/mnt/data/owncloud.log",
            
        },
        "volumes": [
            {
                "name": "data",
                "path": "/mnt/data",
            },
            {
                "name": "core-apps",
                "path": "/var/www/owncloud/apps",
            },
            {
                "name": "config-templates",
                "path": "/etc/templates",
            },
            {
                "name": "preserver-config",
                "path": "/etc/pre_server.d",
            },
            {
                "name": "crons",
                "path": "/tmp",
            },
        ],
        "depends_on": ["wait-for-services", "copy-configs"],
    }]

def keycloakService():
    return [{
        "name": "keycloak",
        "image": KEYCLOAK_IMG,
        "pull": "always",
        "environment": {
            "CLOUD_DOMAIN": OCIS_DOMAIN,
            "OC10_OIDC_CLIENT_SECRET": "oc10-oidc-secret",
            "LDAP_ADMIN_PASSWORD": "admin",
            "DB_VENDOR": "POSTGRES",
            "DB_ADDR": "keycloak-db",
            "DB_DATABASE": "keycloak",
            "DB_USER": "keycloak",
            "DB_SCHEMA": "public",
            "DB_PASSWORD": "keycloak",
            "KEYCLOAK_USER": "admin",
            "KEYCLOAK_PASSWORD": "admin",
            "PROXY_ADDRESS_FORWARDING": "true",
            "KEYCLOAK_IMPORT": "%s/keycloak/owncloud-realm.json" % (DRONE_CONFIG_PATH),
        },
    }]

def ldapService():
    return [{
        "name": "openldap",
        "image": OPENLDAP_IMG,
        "pull": "always",
        "environment": {
            "LDAP_TLS_VERIFY_CLIENT": "never",
            "LDAP_DOMAIN": "owncloud.com",
            "LDAP_ORGANISATION": "owncloud",
            "LDAP_ADMIN_PASSWORD": "admin",
            "LDAP_RFC2307BIS_SCHEMA": "true",
            "LDAP_REMOVE_CONFIG_AFTER_SETUP": "false",
            "LDAP_SEED_INTERNAL_LDIF_PATH": "%s/ldap/ldif" % (DRONE_CONFIG_PATH),
        },
        "command": [
            "--copy-service",
            "--loglevel",
            "debug",
        ],
    }]

def keycloakDbService():
    return [{
        "name": "keycloak-db",
        "image": KEYCLOAK_DB_IMG,
        "pull": "always",
        "environment": {
            "POSTGRES_DB": "keycloak",
            "POSTGRES_USER": "keycloak",
            "POSTGRES_PASSWORD": "keycloak",
        }
    }]

def oc10DbService():
    return [
        {
            "name": "oc10-db",
            "image": OC10_DB_IMG,
            "pull": "always",
            "environment": {
                "MYSQL_ROOT_PASSWORD": "owncloud",
                "MYSQL_USER": "owncloud",
                "MYSQL_PASSWORD": "owncloud",
                "MYSQL_DATABASE": "owncloud",
            },
            "command": [
                "--max-allowed-packet=128M",
                "--innodb-log-file-size=64M",
                "--innodb-read-only-compressed=OFF",
            ]
        },
        
    ]

def redisService():
    return [{
        "name": "redis",
        "image": REDIS_IMG,
        "pull": "always",
        "command": [
            "--databases",
            "1"
        ],
    }]

def composerInstall():
    return [{
        "name": "composer-install",
        "image": OC_CI_PHP,
        "commands": [
            "cd ./vendor-bin/behat",
            "composer install",
        ],
    }]

def copyConfigs():
    return [{
        "name": "copy-configs",
        "image": OC10_IMG,
        "pull": "always",
        "commands": [
            # ocis proxy config
            "mkdir -p /etc/ocis",
            "cp %s/ocis/proxy.json /etc/ocis/proxy.json" % (DRONE_CONFIG_PATH),
            # oc10 configs
            "mkdir -p /etc/templates",
            "mkdir -p /etc/pre_server.d",
            "cp %s/oc10/oidc.config.php /etc/templates/oidc.config.php" % (DRONE_CONFIG_PATH),
            "cp %s/oc10/ldap-config.tmpl.json /etc/templates/ldap-config.tmpl.json" % (DRONE_CONFIG_PATH),
            "cp %s/oc10/web.config.php /etc/templates/web.config.php" % (DRONE_CONFIG_PATH),
            "cp %s/oc10/web-config.tmpl.json /etc/templates/web-config.tmpl.json" % (DRONE_CONFIG_PATH),
            "cp %s/oc10/ldap-sync-cron /tmp/ldap-sync-cron" % (DRONE_CONFIG_PATH),
            "cp %s/oc10/10-custom-config.sh /etc/pre_server.d/10-custom-config.sh" % (DRONE_CONFIG_PATH),
        ],
        "volumes": [
            {
                "name": "proxy-config",
                "path": "/etc/ocis",
            },
            {
                "name": "config-templates",
                "path": "/etc/templates",
            },
            {
                "name": "preserver-config",
                "path": "/etc/pre_server.d",
            },
            {
                "name": "crons",
                "path": "/tmp",
            },
        ]
    }]

def owncloudLog():
    return [{
        "name": "owncloud-log",
        "image": OC_CI_UBUNTU,
        "pull": "always",
        "detach": True,
        "commands": [
            "tail -f /mnt/data/owncloud.log",
        ],
        "volumes": [
            {
                "name": "data",
                "path": "/mnt/data",
            }
        ],
        "depends_on": ["wait-for-oc10"]
    }]

def fixPermissions():
    return [{
        "name": "fix-permissions",
        "image": OC_CI_PHP,
        "pull": "always",
        "commands": [
            "chown -R www-data:www-data /var/www/owncloud/apps",
            "chmod -R 777 /var/www/owncloud/apps",
            "chmod -R 777 /mnt/data/",
            "cd /mnt/data",
            "ls -al",
        ],
        "volumes": [
            {
                "name": "core-apps",
                "path": "/var/www/owncloud/apps",
            },
            {
                "name": "data",
                "path": "/mnt/data",
            }
        ],
        "depends_on": ["wait-for-oc10"],
    }]

def waitForServices():
    return [{
        "name": "wait-for-services",
        "image": OC_CI_WAITFOR,
        "commands": [
            "wait-for -it oc10-db:3306 -t 300",
            "wait-for -it openldap:636 -t 300",
        ],
    }]

def waitForOC10():
    return [{
        "name": "wait-for-oc10",
        "image": OC_CI_WAITFOR,
        "commands": [
            "wait-for -it oc10:8080 -t 300",
        ],
        "depends_on": ["wait-for-services"]
    }]

def waitForOCIS():
    return [{
        "name": "wait-for-ocis",
        "image": OC_CI_WAITFOR,
        "commands": [
            "wait-for -it ocis:9200 -t 300",
        ],
        "depends_on": ["wait-for-oc10"],
    }]