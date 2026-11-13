Bugfix: Copy folders file by file

Copying or duplicating a folder sent one single request, in which the server
copied the whole folder. For big folders, that request took longer than the
request timeout in front of the server and got cut off, so only part of the
files were copied. The web UI now lists the folder, creates the folder and its
subfolders itself and copies every file with its own request, so no request
runs long. When some of the content can't be copied, the rest is still copied,
and duplicating reports the folder as incomplete. Replacing an existing folder,
copying a folder into one of its own subfolders (which the server rejects) and
moving a folder still use one request.

https://github.com/owncloud/ocis/pull/13445
