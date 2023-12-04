Uberspace and ownCloud Infinite Scale in 50 seconds
---------------------------------------------------

If you want to set up ownCloud Infinite Scale for a quick test, here's the video that shows the fastest possible way:  https://cloud.owncloud.com/s/tsieyFn70U3ySm6

Basically, it's all done in three steps – assuming you already have an account for Ubernauten (Join us here: https://dashboard.uberspace.de/register?lang=en)

. Download the Infinite Scale binary and make it executable

. Set some environment variables related to Uberspace

. Start the `ocis` binary, first with the `init` parameter (which also gives you your unique login password for the user `admin`, then again with `ocis start`:

. Visit the url of your uberspace server and login:
+
image:https://cloud.owncloud.com/index.php/apps/files_sharing/ajax/publicpreview.php?x=1920&y=645&a=true&file=uberspace-login.png&t=bGkiHY25YGFtBQQ&scalingup=0[alt="Login to ownCloud Infinite Scale",width=400,height=400]


These are the commands needed: 

[source,bash]
---------
curl https://download.owncloud.com/ocis/ocis/stable/4.0.2/ocis-4.0.2-linux-amd64 --output ocis
chmod +x ocis
uberspace web backend set / --http --port 9200
export OCIS_URL=https://ocis.uber.space
export PROXY_TLS=false
export PROXY_HTTP_ADDR=0.0.0.0:9200
export PROXY_LOG_LEVEL=debug
export THUMBNAILS_WEBDAVSOURCEBASE_URL=http://localhost:9200/remote.php/webdav/
./ocis init
./ocis server
---------

If you omit (or forget) the `ocis init` command, you will get the following error message:

[source,bash]
---------
[mfeilner@apus ~]$ ./ocis server
The jwt_secret has not been set properly in your config for ocis. Make sure your /home/mfeilner/.ocis/config config contains the proper values (e.g. by running ocis init or setting it manually in the config/corresponding environment variable).
The jwt_secret has not been set properly in your config for ocis. Make sure your /home/mfeilner/.ocis/config config contains the proper values (e.g. by running ocis init or setting it manually in the config/corresponding environment variable).
---------

For copy and paste, here's these commands in a script I called `ocis.start`: 

[source,bash]
-----------
#!/bin/bash
# This file is named ocis.install
# It downloads ocis, configures the environment varialbes and starts ownCloud Infinite Scale
curl https://download.owncloud.com/ocis/ocis/stable/4.0.2/ocis-4.0.2-linux-amd64 --output ocis
chmod +x ocis
uberspace web backend set / --http --port 9200
export OCIS_URL=https://ocis.uber.space
export PROXY_TLS=false
export PROXY_HTTP_ADDR=0.0.0.0:9200
export PROXY_LOG_LEVEL=debug
export THUMBNAILS_WEBDAVSOURCEBASE_URL=http://localhost:9200/remote.php/webdav/
./ocis init
./ocis server
-----------
Service Management with Supervisord
-----------------------------------

If you want ocis to run continuously, you need to configure `supervisord` (http://supervisord.org) which is the tool Uberspace is using for service management. 

You can start and stop services with `supervisorctl`, it will (re)read configuration files it finds in your home directory, under `etc/services.d/`, in `.ini` files. The content of these files is very simple, you only have to enter three lines, here is the example for Infinite Scale in `/home/ocis/etc/services.d/ocis.ini`.  

[source,bash]
--------
[program:ocis]
command="/home/ocis/ocis.start"
startsecs=60
--------

`ocis.start` is a script that combines all of the commands above except for the download of the ocis binary. It looks like this:

[source,bash]
------------
#!/bin/bash
# This file is named ocis.start
/usr/bin/uberspace web backend set / --http --port 9200 &
export OCIS_URL=https://ocis.uber.space
export PROXY_TLS=false
export PROXY_HTTP_ADDR=0.0.0.0:9200
export PROXY_LOG_LEVEL=debug
export THUMBNAILS_WEBDAVSOURCEBASE_URL=http://localhost:9200/remote.php/webdav/
/home/ocis/ocis server
------------

There are three supervisorctl commands that you will find useful (many more can be found in its documentation). You can use `supervisorctl status` to check which services managed by supervisorctl are running at the moment, a `reread` will be necessary after you changed the `ini` files, and an `update` is applying the changes:
 
[source,bash]
-----------
[ocis@rigel ~]$ supervisorctl status
ocis                             RUNNING   pid 9813, uptime 0:01:40
[ocis@rigel ~]$ supervisorctl reread
No config updates to processes
[ocis@rigel ~]$ supervisorctl update
-----------

Updating ownCloud Infinite Scale
--------------------------------

Updating the ocis binary is simple: When a new version comes to life, just download the new `ocis` binary from the download server, replacing the old `ocis` executable on your uberspace server. 

Make a backup of your data and make sure you have read and understood the release notes of your new version , especially the "breaking changes" section before starting the binary. Don't worry, you can always go back to the older version you had installed, there's a long list of older versions available. 

Wiping and Clean Restart from Scratch
-------------------------------------

This little script is removing your ocis installation (and *all of your data*!), replacing it with a new, clean ocis installation. Be careful and only use it for testing purposes. Specify your desired ocis version in the curl command. 

---------
[source,bash]
#!/bin/bash
# This file is named ocis.reinstall 
rm -rf .ocis
curl https://download.owncloud.com/ocis/ocis/stable/4.0.2/ocis-4.0.2-linux-amd64 --output ocis
chmod +x ocis
uberspace web backend set / --http --port 9200
export OCIS_URL=https://ocis.uber.space
export PROXY_TLS=false
export PROXY_HTTP_ADDR=0.0.0.0:9200
export PROXY_LOG_LEVEL=debug
export THUMBNAILS_WEBDAVSOURCEBASE_URL=http://localhost:9200/remote.php/webdav/
./ocis init
./ocis server
---------

