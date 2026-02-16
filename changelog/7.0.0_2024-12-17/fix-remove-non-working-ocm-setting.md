Bugfix: Removed 'OCM_OCM_PROVIDER_AUTHORIZER_VERIFY_REQUEST_HOSTNAME' setting

The config option 'OCM_OCM_PROVIDER_AUTHORIZER_VERIFY_REQUEST_HOSTNAME' was
removed from the OCM service. The additional security provided by this setting
is somewhat questionable and only provided in very specific setups.

We are not going through the normal deprecation process for this setting, as it
was never really working anyway. If you have this setting in your configuration,
it will be ignored. You can safely remove it.

https://github.com/owncloud/ocis/pull/10425
https://github.com/owncloud/ocis/issues/10355
