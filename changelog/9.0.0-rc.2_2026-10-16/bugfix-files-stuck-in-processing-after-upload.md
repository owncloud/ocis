Bugfix: Files no longer stay in processing after upload or copy

After uploading or copying files, some of them sometimes stayed greyed out as
"processing" until the page was reloaded. This happened when processing
finished while the file list was still being loaded: the update from the
server was ignored and the list then showed the file as still processing.

Files whose processing has finished are now refreshed automatically when they
show up in the file list as processing.

https://github.com/owncloud/ocis/pull/13134
