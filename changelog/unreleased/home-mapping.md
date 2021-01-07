Enhancement: Functionality to map home directory to different storage providers

We added a parameter in reva that allows us to redirect /home requests to
different storage providers based on a mapping derived from the user attributes,
which was previously not possible since we hardcode the /home path for all
users. This PR adds the config for that parameter.

https://github.com/owncloud/ocis/pull/1186
https://github.com/cs3org/reva/pull/1142
