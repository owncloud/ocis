Enhancement: Functionality to map home directory to different storage providers

We added a parameter in reva that allows us to redirect /home requests to
different storage providers based on a mapping derived from the user attributes,
which was previously not possible since we hardcode the /home path for all
users. For example, having its value as `/home/{{substr 0 1 .Username}}` can be
used to redirect home requests for different users to different storage
providers.

https://github.com/owncloud/ocis/pull/1186
https://github.com/cs3org/reva/pull/1142
