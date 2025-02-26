Bugfix: Mask user email in output

We have fixed a bug where the user email was not masked in the output and the user emails could be enumerated through
the sharee search. This is the ocis side which adds an suiting config option to mask user emails in the output.

https://github.com/owncloud/ocis/issues/8726
https://github.com/cs3org/reva/pull/4603
https://github.com/owncloud/ocis/pull/8764
