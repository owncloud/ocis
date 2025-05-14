
---
title: "Migrate Data using rclone"
date: 2020-06-12T14:35:00+01:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/guides
geekdocFilePath: migrate-data-rclone.md
geekdocCollapseSection: true
---

People keep asking on how to migrate data from other cloud storage to Infinite Scale.

There are too many cloud variants and use cases out there to have a migration path for all at hand, but let's see what we can start with: There is the famous Swiss army knife for clouds called rclone available.

The awesome rclone tool makes it easy to migrate data from one installation to another on a user per user base. A very good first step.

This article explains by the example of Nextcloud how you would migrate your data from an running NC to Infinite Scale. A prerequisite is that you have Infinite Scale already set up and running on a different domain.

{{< hint warning >}}
Be prepared that migration can take some time. Also, check the size of your data. This example was around 1.5 GiB of data, that obviously went fast.

And of course: Have a backup! Even if this method only reads from the source, you never know.
{{< /hint >}}

## Install rclone

Check the [rclone website on how to install](https://rclone.org/install/) rclone.

## Create Users

First, decide on a user that you want to migrate. To create the user on Infinite Scale, log in as an admin user and create the desired user. Create and set groups accordingly.

For the next step, you need the user credentials on both the source- and the destination cloud.

## Configure rclone Remotes

To be able to address the clouds in rclone, you need to configure so called _remotes_. It is nothing else than a shortcut for the combination of
- which kind of cloud are you talking to
- the URL
- the username
- the password, if one is set

You need to add a configuration for both the source cloud (Nextcloud) and the target (Infinite Scale). As both talk WebDAV, the [rclone manual page](https://rclone.org/webdav/) is accurate to follow.

For both, use an URL in the form of `https://my.host.de/remote.php/webdav`.

Once that is finished, the command `rclone config show` should give output similar to this:

```bash
[:~/] ± rclone config show

[NCoC]
type = webdav
url = https://nc.this.de/remote.php/webdav
vendor = nextcloud
user = wilma
pass = zfdsaiewrafdskfjdasfxdasffdafdsafas

[ocis]
type = webdav
url = https://infinitescale.works/remote.php/webdav
vendor = owncloud
user = wilma
pass = cdsfasrefdsadaGkxTXjksfpqQFI5nQawqs

```

Now, for example the directories on the Nextcloud root can be checked with `rclone lsd NCoC:/`.

## Copy Data

To migrate the data, rclone provides the command `copy`. It transfers data from one remote to the other. Use the following command example to transfer the entire cloud data from Nextcloud to Infinite Scale:
```
rclone  copy NCoC:/ ocis:/ --no-check-certificate -P
```
The --no-check-certificate can and should be skipped if your clouds have proper certificates. The `-P` however, provides you with interesting statistics about the copy progress.
Once you are finished, this might be the result:
```
[:~/] $ rclone copy NCoC:/ ocis:/ --no-check-certificate -P
Transferred:        1.228 GiB / 1.228 GiB, 100%, 10.170 MiB/s, ETA 0s
Transferred:          411 / 411, 100%
Elapsed time:      2m19.3s
```

Note that while testing this, occasionally the Nextcloud was returning a `404 not found` for some files. While the reason for that was not completely clear, it does not matter, because the rclone command can be repeated. It is clever enough to only copy what has changed!

## Enjoy!

All done! Now you have your data on Infinite Scale.

Obviously this method has a few downsides, such as:
- This migration requires a little of "quiet time" for migrating data.
- It is a user by user method where provisioning of users has to be done manually.
- Only data is migrated, and there is probably a data size limit in real life using this way.
- Private- and public shares are not migrated
- The trashbin, versions, comments and favorites are not migrated

These are shortcomings but this is a good first step to start investigating. The other parts will be sorted as we move along.

---
To improve this guide, you are welcome to file an issue or even send a pull request. See the [getting started guide](https://owncloud.dev/ocis/development/build-docs/) how easy it is to build the documentation.

