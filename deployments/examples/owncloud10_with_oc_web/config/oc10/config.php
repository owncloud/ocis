<?php

function getConfigFromEnv() {
  if (isset($_SERVER['HTTP_X_FORWARDED_HOST'])) {
    $domain = trim(
      explode(
        ",",
        $_SERVER['HTTP_X_FORWARDED_HOST']
      )[0]
    );
  } else if (isset($_SERVER['SERVER_NAME'])) {
    $domain = $_SERVER['SERVER_NAME'];
  } else {
    $domain = 'localhost';
  }

  $config = [
    'apps_paths' => [
      0 => [
        "path" => OC::$SERVERROOT . "/apps",
        "url" => "/apps",
        "writable" => false
      ],
      1 => [
        "path" => OC::$SERVERROOT . "/custom",
        "url" => "/custom",
        "writable" => true
      ]
    ],

    'trusted_domains' => [
      0 => $domain
    ],
    'openid-connect' => [
        'provider-url' => getenv('OCIS_URL'),
        'client-id' => 'oc10',
        'client-secret' => 'super',
        'loginButtonName' => 'OpenId Connect',
        'search-attribute' => 'preferred_username',
        'mode' => 'userid',
        'autoRedirectOnLoginPage' => true,
        'insecure' => true,
        'post_logout_redirect_uri' => getenv('OWNCLOUD_DOMAIN') . '/',
      ],
    'datadirectory' => getenv('OWNCLOUD_VOLUME_FILES'),
    'dbtype' => getenv('OWNCLOUD_DB_TYPE'),
    'dbhost' => getenv('OWNCLOUD_DB_HOST'),
    'dbname' => getenv('OWNCLOUD_DB_NAME'),
    'dbuser' => getenv('OWNCLOUD_DB_USERNAME'),
    'dbpassword' => getenv('OWNCLOUD_DB_PASSWORD'),
    'dbtableprefix' => getenv('OWNCLOUD_DB_PREFIX'),

    'web.baseUrl' => getenv('OCIS_URL'),
    'cors.allowed-domains' => [getenv('OCIS_URL')],

    'log_type' => 'owncloud',

    'supportedDatabases' => [
      'sqlite',
      'mysql',
      'pgsql',
    ],

    'upgrade.disable-web' => true,
  ];

  if (getenv('OWNCLOUD_CORS_ALLOWED_DOMAINS') != '') {
    $config['cors.allowed-domains'] = explode(',', getenv('OWNCLOUD_CORS_ALLOWED_DOMAINS'));
  }

  if (getenv('OWNCLOUD_VERSION_HIDE') != '') {
    $config['version.hide'] = getenv('OWNCLOUD_VERSION_HIDE') == 'true';
  }

  if (getenv('OWNCLOUD_SHOW_SERVER_HOSTNAME') != '') {
    $config['show_server_hostname'] = getenv('OWNCLOUD_SHOW_SERVER_HOSTNAME') == 'true';
  }

  if (getenv('OWNCLOUD_DEFAULT_LANGUAGE') != '') {
    $config['default_language'] = getenv('OWNCLOUD_DEFAULT_LANGUAGE');
  }

  if (getenv('OWNCLOUD_DEFAULT_APP') != '') {
    $config['defaultapp'] = getenv('OWNCLOUD_DEFAULT_APP');
  }

  if (getenv('OWNCLOUD_KNOWLEDGEBASE_ENABLED') != '') {
    $config['knowledgebaseenabled'] = getenv('OWNCLOUD_KNOWLEDGEBASE_ENABLED') == 'true';
  }

  if (getenv('OWNCLOUD_ENABLE_AVATARS') != '') {
    $config['enable_avatars'] = getenv('OWNCLOUD_ENABLE_AVATARS') == 'true';
  }

  if (getenv('OWNCLOUD_ALLOW_USER_TO_CHANGE_DISPLAY_NAME') != '') {
    $config['allow_user_to_change_display_name'] = getenv('OWNCLOUD_ALLOW_USER_TO_CHANGE_DISPLAY_NAME') == 'true';
  }

  if (getenv('OWNCLOUD_REMEMBER_LOGIN_COOKIE_LIFETIME') != '') {
    $config['remember_login_cookie_lifetime'] = (int) getenv('OWNCLOUD_REMEMBER_LOGIN_COOKIE_LIFETIME');
  }

  if (getenv('OWNCLOUD_SESSION_LIFETIME') != '') {
    $config['session_lifetime'] = (int) getenv('OWNCLOUD_SESSION_LIFETIME');
  }

  if (getenv('OWNCLOUD_SESSION_KEEPALIVE') != '') {
    $config['session_keepalive'] = getenv('OWNCLOUD_SESSION_KEEPALIVE') == 'true';
  }

  if (getenv('OWNCLOUD_TOKEN_AUTH_ENFORCED') != '') {
    $config['token_auth_enforced'] = getenv('OWNCLOUD_TOKEN_AUTH_ENFORCED') == 'true';
  }

  if (getenv('OWNCLOUD_CSRF_DISABLED') != '') {
    $config['csrf.disabled'] = getenv('OWNCLOUD_CSRF_DISABLED') == 'true';
  }

  if (getenv('OWNCLOUD_SKELETON_DIRECTORY') != '') {
    $config['skeletondirectory'] = getenv('OWNCLOUD_SKELETON_DIRECTORY');
  }

  if (getenv('OWNCLOUD_LOST_PASSWORD_LINK') != '') {
    $config['lost_password_link'] = getenv('OWNCLOUD_LOST_PASSWORD_LINK');
  }

  if (getenv('OWNCLOUD_ACCOUNTS_ENABLE_MEDIAL_SEARCH') != '') {
    $config['accounts.enable_medial_search'] = getenv('OWNCLOUD_ACCOUNTS_ENABLE_MEDIAL_SEARCH') == 'true';
  }

  if (getenv('OWNCLOUD_USER_SEARCH_MIN_LENGTH') != '') {
    $config['user.search_min_length'] = (int) getenv('OWNCLOUD_USER_SEARCH_MIN_LENGTH');
  }

  if (getenv('OWNCLOUD_MAIL_DOMAIN') != '') {
    $config['mail_domain'] = getenv('OWNCLOUD_MAIL_DOMAIN');
  }

  if (getenv('OWNCLOUD_MAIL_FROM_ADDRESS') != '') {
    $config['mail_from_address'] = getenv('OWNCLOUD_MAIL_FROM_ADDRESS');
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_DEBUG') != '') {
    $config['mail_smtpdebug'] = getenv('OWNCLOUD_MAIL_SMTP_DEBUG') == 'true';
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_MODE') != '') {
    $config['mail_smtpmode'] = getenv('OWNCLOUD_MAIL_SMTP_MODE');
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_HOST') != '') {
    $config['mail_smtphost'] = getenv('OWNCLOUD_MAIL_SMTP_HOST');
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_PORT') != '') {
    $config['mail_smtpport'] = (int) getenv('OWNCLOUD_MAIL_SMTP_PORT');
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_TIMEOUT') != '') {
    $config['mail_smtptimeout'] = (int) getenv('OWNCLOUD_MAIL_SMTP_TIMEOUT');
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_SECURE') != '') {
    $config['mail_smtpsecure'] = getenv('OWNCLOUD_MAIL_SMTP_SECURE');
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_AUTH') != '') {
    $config['mail_smtpauth'] = getenv('OWNCLOUD_MAIL_SMTP_AUTH') == 'true';
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_AUTH_TYPE') != '') {
    $config['mail_smtpauthtype'] = getenv('OWNCLOUD_MAIL_SMTP_AUTH_TYPE');
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_NAME') != '') {
    $config['mail_smtpname'] = getenv('OWNCLOUD_MAIL_SMTP_NAME');
  }

  if (getenv('OWNCLOUD_MAIL_SMTP_PASSWORD') != '') {
    $config['mail_smtppassword'] = getenv('OWNCLOUD_MAIL_SMTP_PASSWORD');
  }

  if (getenv('OWNCLOUD_OVERWRITE_HOST') != '') {
    $config['overwritehost'] = getenv('OWNCLOUD_OVERWRITE_HOST');
  }

  if (getenv('OWNCLOUD_OVERWRITE_PROTOCOL') != '') {
    $config['overwriteprotocol'] = getenv('OWNCLOUD_OVERWRITE_PROTOCOL');
  }

  if (getenv('OWNCLOUD_OVERWRITE_WEBROOT') != '') {
    $config['overwritewebroot'] = getenv('OWNCLOUD_OVERWRITE_WEBROOT');
  }

  if (getenv('OWNCLOUD_OVERWRITE_COND_ADDR') != '') {
    $config['overwritecondaddr'] = getenv('OWNCLOUD_OVERWRITE_COND_ADDR');
  }

  if (getenv('OWNCLOUD_OVERWRITE_CLI_URL') != '') {
    $config['overwrite.cli.url'] = getenv('OWNCLOUD_OVERWRITE_CLI_URL');
  }

  if (getenv('OWNCLOUD_HTACCESS_REWRITE_BASE') != '') {
    $config['htaccess.RewriteBase'] = getenv('OWNCLOUD_HTACCESS_REWRITE_BASE');
  }

  if (getenv('OWNCLOUD_PROXY') != '') {
    $config['proxy'] = getenv('OWNCLOUD_PROXY');
  }

  if (getenv('OWNCLOUD_PROXY_USERPWD') != '') {
    $config['proxyuserpwd'] = getenv('OWNCLOUD_PROXY_USERPWD');
  }

  if (getenv('OWNCLOUD_TRASHBIN_RETENTION_OBLIGATION') != '') {
    $config['trashbin_retention_obligation'] = getenv('OWNCLOUD_TRASHBIN_RETENTION_OBLIGATION');
  }

  if (getenv('OWNCLOUD_TRASHBIN_PURGE_LIMIT') != '') {
    $config['trashbin_purge_limit'] = (int) getenv('OWNCLOUD_TRASHBIN_PURGE_LIMIT');
  }

  if (getenv('OWNCLOUD_VERSIONS_RETENTION_OBLIGATION') != '') {
    $config['versions_retention_obligation'] = getenv('OWNCLOUD_VERSIONS_RETENTION_OBLIGATION');
  }

  if (getenv('OWNCLOUD_UPDATE_CHECKER') != '') {
    $config['updatechecker'] = getenv('OWNCLOUD_UPDATE_CHECKER') == 'true';
  }

  if (getenv('OWNCLOUD_UPDATER_SERVER_URL') != '') {
    $config['updater.server.url'] = getenv('OWNCLOUD_UPDATER_SERVER_URL');
  }

  if (getenv('OWNCLOUD_HAS_INTERNET_CONNECTION') != '') {
    $config['has_internet_connection'] = getenv('OWNCLOUD_HAS_INTERNET_CONNECTION') == 'true';
  }

  if (getenv('OWNCLOUD_CHECK_FOR_WORKING_WELLKNOWN_SETUP') != '') {
    $config['check_for_working_wellknown_setup'] = getenv('OWNCLOUD_CHECK_FOR_WORKING_WELLKNOWN_SETUP') == 'true';
  }

  if (getenv('OWNCLOUD_OPERATION_MODE') != '') {
    $config['operation.mode'] = getenv('OWNCLOUD_OPERATION_MODE');
  }

  if (getenv('OWNCLOUD_LOG_FILE') != '') {
    $config['logfile'] = getenv('OWNCLOUD_LOG_FILE');
  }

  if (getenv('OWNCLOUD_LOG_LEVEL') != '') {
    $config['loglevel'] = (int) getenv('OWNCLOUD_LOG_LEVEL');
  }

  if (getenv('OWNCLOUD_LOG_DATE_FORMAT') != '') {
    $config['logdateformat'] = getenv('OWNCLOUD_LOG_DATE_FORMAT');
  }

  if (getenv('OWNCLOUD_LOG_TIMEZONE') != '') {
    $config['logtimezone'] = getenv('OWNCLOUD_LOG_TIMEZONE');
  }

  if (getenv('OWNCLOUD_CRON_LOG') != '') {
    $config['cron_log'] = getenv('OWNCLOUD_CRON_LOG') == 'true';
  }

  if (getenv('OWNCLOUD_LOG_ROTATE_SIZE') != '') {
    $config['log_rotate_size'] = (int) getenv('OWNCLOUD_LOG_ROTATE_SIZE');
  }

  if (getenv('OWNCLOUD_ENABLE_PREVIEWS') != '') {
    $config['enable_previews'] = getenv('OWNCLOUD_ENABLE_PREVIEWS') == 'true';
  }

  if (getenv('OWNCLOUD_PREVIEW_MAX_X') != '') {
    $config['preview_max_x'] = (int) getenv('OWNCLOUD_PREVIEW_MAX_X');
  }

  if (getenv('OWNCLOUD_PREVIEW_MAX_Y') != '') {
    $config['preview_max_y'] = (int) getenv('OWNCLOUD_PREVIEW_MAX_Y');
  }

  if (getenv('OWNCLOUD_PREVIEW_MAX_SCALE_FACTOR') != '') {
    $config['preview_max_scale_factor'] = (int) getenv('OWNCLOUD_PREVIEW_MAX_SCALE_FACTOR');
  }

  if (getenv('OWNCLOUD_PREVIEW_MAX_FILESIZE_IMAGE') != '') {
    $config['preview_max_filesize_image'] = getenv('OWNCLOUD_PREVIEW_MAX_FILESIZE_IMAGE');
  }

  if (getenv('OWNCLOUD_PREVIEW_LIBREOFFICE_PATH') != '') {
    $config['preview_libreoffice_path'] = getenv('OWNCLOUD_PREVIEW_LIBREOFFICE_PATH');
  }

  if (getenv('OWNCLOUD_PREVIEW_OFFICE_CL_PARAMETERS') != '') {
    $config['preview_office_cl_parameters'] = getenv('OWNCLOUD_PREVIEW_OFFICE_CL_PARAMETERS');
  }

  if (getenv('OWNCLOUD_ENABLED_PREVIEW_PROVIDERS') != '') {
    $config['enabledPreviewProviders'] = explode(',', getenv('OWNCLOUD_ENABLED_PREVIEW_PROVIDERS'));
  }

  if (getenv('OWNCLOUD_COMMENTS_MANAGER_FACTORY') != '') {
    $config['comments.managerFactory'] = getenv('OWNCLOUD_COMMENTS_MANAGER_FACTORY');
  }

  if (getenv('OWNCLOUD_SYSTEMTAGS_MANAGER_FACTORY') != '') {
    $config['systemtags.managerFactory'] = getenv('OWNCLOUD_SYSTEMTAGS_MANAGER_FACTORY');
  }

  if (getenv('OWNCLOUD_MAINTENANCE') != '') {
    $config['maintenance'] = getenv('OWNCLOUD_MAINTENANCE') == 'true';
  }

  if (getenv('OWNCLOUD_SINGLEUSER') != '') {
    $config['singleuser'] = getenv('OWNCLOUD_SINGLEUSER');
  }

  if (getenv('OWNCLOUD_ENABLE_CERTIFICATE_MANAGEMENT') != '') {
    $config['enable_certificate_management'] = getenv('OWNCLOUD_ENABLE_CERTIFICATE_MANAGEMENT');
  }

  if (getenv('OWNCLOUD_MEMCACHE_LOCAL') != '') {
    $config['memcache.local'] = getenv('OWNCLOUD_MEMCACHE_LOCAL');
  }

  if (getenv('OWNCLOUD_CACHE_PATH') != '') {
    $config['cache_path'] = getenv('OWNCLOUD_CACHE_PATH');
  }

  if (getenv('OWNCLOUD_CACHE_CHUNK_GC_TTL') != '') {
    $config['cache_chunk_gc_ttl'] = (int) getenv('OWNCLOUD_CACHE_CHUNK_GC_TTL');
  }

  if (getenv('OWNCLOUD_DAV_CHUNK_BASE_DIR') != '') {
    $config['dav.chunk_base_dir'] = getenv('OWNCLOUD_DAV_CHUNK_BASE_DIR');
  }

  if (getenv('OWNCLOUD_SHARING_MANAGER_FACTORY') != '') {
    $config['sharing.managerFactory'] = getenv('OWNCLOUD_SHARING_MANAGER_FACTORY');
  }

  if (getenv('OWNCLOUD_SHARING_FEDERATION_ALLOW_HTTP_FALLBACK') != '') {
    $config['sharing.federation.allowHttpFallback'] = getenv('OWNCLOUD_SHARING_FEDERATION_ALLOW_HTTP_FALLBACK') == 'true';
  }

  if (getenv('OWNCLOUD_SQLITE_JOURNAL_MODE') != '') {
    $config['sqlite.journal_mode'] = getenv('OWNCLOUD_SQLITE_JOURNAL_MODE');
  }

  if (getenv('OWNCLOUD_MYSQL_UTF8MB4') != '') {
    $config['mysql.utf8mb4'] = getenv('OWNCLOUD_MYSQL_UTF8MB4') == 'true';
  }

  if (getenv('OWNCLOUD_TEMP_DIRECTORY') != '') {
    $config['tempdirectory'] = getenv('OWNCLOUD_TEMP_DIRECTORY');
  }

  if (getenv('OWNCLOUD_HASHING_COST') != '') {
    $config['hashingCost'] = (int) getenv('OWNCLOUD_HASHING_COST');
  }

  if (getenv('OWNCLOUD_BLACKLISTED_FILES') != '') {
    $config['blacklisted_files'] = explode(',', getenv('OWNCLOUD_BLACKLISTED_FILES'));
  }

  if (getenv('OWNCLOUD_EXCLUDED_DIRECTORIES') != '') {
    $config['excluded_directories'] = explode(',', getenv('OWNCLOUD_EXCLUDED_DIRECTORIES'));
  }

  if (getenv('OWNCLOUD_INTEGRITY_EXCLUDED_FILES') != '') {
    $config['integrity.excluded.files'] = explode(',', getenv('OWNCLOUD_INTEGRITY_EXCLUDED_FILES'));
  }

  if (getenv('OWNCLOUD_INTEGRITY_IGNORE_MISSING_APP_SIGNATURE') != '') {
    $config['integrity.ignore.missing.app.signature'] = explode(',', getenv('OWNCLOUD_INTEGRITY_IGNORE_MISSING_APP_SIGNATURE'));
  }

  if (getenv('OWNCLOUD_SHARE_FOLDER') != '') {
    $config['share_folder'] = getenv('OWNCLOUD_SHARE_FOLDER');
  }

  if (getenv('OWNCLOUD_CIPHER') != '') {
    $config['cipher'] = getenv('OWNCLOUD_CIPHER');
  }

  if (getenv('OWNCLOUD_MINIMUM_SUPPORTED_DESKTOP_VERSION') != '') {
    $config['minimum.supported.desktop.version'] = getenv('OWNCLOUD_MINIMUM_SUPPORTED_DESKTOP_VERSION');
  }

  if (getenv('OWNCLOUD_QUOTA_INCLUDE_EXTERNAL_STORAGE') != '') {
    $config['quota_include_external_storage'] = getenv('OWNCLOUD_QUOTA_INCLUDE_EXTERNAL_STORAGE') == 'true';
  }

  if (getenv('OWNCLOUD_FILESYSTEM_CHECK_CHANGES') != '') {
    $config['filesystem_check_changes'] = (int) getenv('OWNCLOUD_FILESYSTEM_CHECK_CHANGES');
  }

  if (getenv('OWNCLOUD_PART_FILE_IN_STORAGE') != '') {
    $config['part_file_in_storage'] = getenv('OWNCLOUD_PART_FILE_IN_STORAGE') == 'true';
  }

  if (getenv('OWNCLOUD_MOUNT_FILE') != '') {
    $config['mount_file'] = getenv('OWNCLOUD_MOUNT_FILE');
  }

  if (getenv('OWNCLOUD_FILESYSTEM_CACHE_READONLY') != '') {
    $config['filesystem_cache_readonly'] = getenv('OWNCLOUD_FILESYSTEM_CACHE_READONLY') == 'true';
  }

  if (getenv('OWNCLOUD_SECRET') != '') {
    $config['secret'] = getenv('OWNCLOUD_SECRET');
  }

  if (getenv('OWNCLOUD_TRUSTED_PROXIES') != '') {
    $config['trusted_proxies'] = explode(',', getenv('OWNCLOUD_TRUSTED_PROXIES'));
  }

  if (getenv('OWNCLOUD_FORWARDED_FOR_HEADERS') != '') {
    $config['forwarded_for_headers'] = explode(',', getenv('OWNCLOUD_FORWARDED_FOR_HEADERS'));
  }

  if (getenv('OWNCLOUD_MAX_FILESIZE_ANIMATED_GIFS_PUBLIC_SHARING') != '') {
    $config['max_filesize_animated_gifs_public_sharing'] = (int) getenv('OWNCLOUD_MAX_FILESIZE_ANIMATED_GIFS_PUBLIC_SHARING');
  }

  if (getenv('OWNCLOUD_FILELOCKING_ENABLED') != '') {
    $config['filelocking.enabled'] = getenv('OWNCLOUD_FILELOCKING_ENABLED') == 'true';
  }

  if (getenv('OWNCLOUD_FILELOCKING_TTL') != '') {
    $config['filelocking.ttl'] = getenv('OWNCLOUD_FILELOCKING_TTL');
  }

  if (getenv('OWNCLOUD_MEMCACHE_LOCKING') != '') {
    $config['memcache.locking'] = getenv('OWNCLOUD_MEMCACHE_LOCKING');
  }

  if (getenv('OWNCLOUD_UPGRADE_AUTOMATIC_APP_UPDATES') != '') {
    $config['upgrade.automatic-app-update'] = getenv('OWNCLOUD_UPGRADE_AUTOMATIC_APP_UPDATES') == 'true';
  }

  if (getenv('OWNCLOUD_DEBUG') != '') {
    $config['debug'] = getenv('OWNCLOUD_DEBUG') == 'true';
  }

  if (getenv('OWNCLOUD_FILES_EXTERNAL_ALLOW_NEW_LOCAL') != '') {
    $config['files_external_allow_create_new_local'] = getenv('OWNCLOUD_FILES_EXTERNAL_ALLOW_NEW_LOCAL') == 'true';
  }

  if (getenv('OWNCLOUD_SMB_LOGGING_ENABLE') != '') {
    $config['smb.logging.enable'] = getenv('OWNCLOUD_SMB_LOGGING_ENABLE');
  }

  if (getenv('OWNCLOUD_DAV_ENABLE_ASYNC') != '') {
    $config['dav.enable.async'] = getenv('OWNCLOUD_DAV_ENABLE_ASYNC');
  }

  if (getenv('OWNCLOUD_LICENSE_KEY') != '') {
    $config['license-key'] = getenv('OWNCLOUD_LICENSE_KEY');
  }

  if (getenv('OWNCLOUD_MARKETPLACE_KEY') != '') {
    $config['marketplace.key'] = getenv('OWNCLOUD_MARKETPLACE_KEY');
  }

  if (getenv('OWNCLOUD_MARKETPLACE_CA') != '') {
    $config['marketplace.ca'] = getenv('OWNCLOUD_MARKETPLACE_CA');
  }

  if (getenv('OWNCLOUD_APPSTORE_URL') != '') {
    $config['appstoreurl'] = getenv('OWNCLOUD_APPSTORE_URL');
  }

  if (getenv('OWNCLOUD_LOGIN_ALTERNATIVES') != '') {
    $rows = explode(',', getenv('OWNCLOUD_LOGIN_ALTERNATIVES'));

    foreach ($rows as $key => $value) {
      parse_str($value, $opts);
      $config['login.alternatives'][$key] = $opts;
    }
  }

  switch (true) {
    case getenv('OWNCLOUD_REDIS_ENABLED') && getenv('OWNCLOUD_REDIS_ENABLED') == 'true':
      $config = array_merge_recursive($config, [
        'memcache.distributed' => '\OC\Memcache\Redis',
        'memcache.locking' => '\OC\Memcache\Redis',
      ]);
      switch (true) {
        case getenv('OWNCLOUD_REDIS_SEEDS') != '':
          $config['redis.cluster']['seeds'] = explode(',', getenv('OWNCLOUD_REDIS_SEEDS'));

          if (getenv('OWNCLOUD_REDIS_TIMEOUT') != '') {
            $config['redis.cluster']['timeout'] = (float) getenv('OWNCLOUD_REDIS_TIMEOUT');
          }

          if (getenv('OWNCLOUD_REDIS_READ_TIMEOUT') != '') {
            $config['redis.cluster']['read_timeout'] = (float) getenv('OWNCLOUD_REDIS_READ_TIMEOUT');
          }

          if (getenv('OWNCLOUD_REDIS_FAILOVER_MODE') != '') {
            switch (getenv('OWNCLOUD_REDIS_FAILOVER_MODE')) {
              case 'FAILOVER_NONE':
                $config['redis.cluster']['failover_mode'] = \RedisCluster::FAILOVER_NONE;
              case 'FAILOVER_ERROR':
                $config['redis.cluster']['failover_mode'] = \RedisCluster::FAILOVER_ERROR;
              case 'FAILOVER_DISTRIBUTE':
                $config['redis.cluster']['failover_mode'] = \RedisCluster::FAILOVER_DISTRIBUTE;
            }
          }

        case getenv('OWNCLOUD_REDIS_HOST') != '':
          $config['redis']['host'] = getenv('OWNCLOUD_REDIS_HOST');
          $config['redis']['port'] = getenv('OWNCLOUD_REDIS_PORT');

          if (getenv('OWNCLOUD_REDIS_DB') != '') {
            $config['redis']['dbindex'] = getenv('OWNCLOUD_REDIS_DB');
          }

          if (getenv('OWNCLOUD_REDIS_PASSWORD') != '') {
            $config['redis']['password'] = getenv('OWNCLOUD_REDIS_PASSWORD');
          }

          if (getenv('OWNCLOUD_REDIS_TIMEOUT') != '') {
            $config['redis']['timeout'] = (float) getenv('OWNCLOUD_REDIS_TIMEOUT');
          }
      }

      break;
    case getenv('OWNCLOUD_MEMCACHED_ENABLED') && getenv('OWNCLOUD_MEMCACHED_ENABLED') == 'true':
      $config = array_merge_recursive($config, [
        'memcache.distributed' => '\OC\Memcache\Memcached',
        'memcache.locking' => '\OC\Memcache\Memcached',

        'memcached_servers' => [
          [
            getenv('OWNCLOUD_MEMCACHED_HOST'),
            getenv('OWNCLOUD_MEMCACHED_PORT'),
          ],
        ],
      ]);

      if (getenv('OWNCLOUD_MEMCACHED_OPTIONS') != '') {
        parse_str(getenv('OWNCLOUD_MEMCACHED_OPTIONS'), $opts);

        foreach($opts as $key => $value) {
          $config['memcached_options'][constant($key)] = $value;
        }
      }

      break;
  }

  return $config;
}

$CONFIG = getConfigFromEnv();
