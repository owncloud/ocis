Bugfix: Make storage users mount ids unique by default

The mount ID of the storage users provider needs to be unique by default. We made this value configurable and added it to ocis init to be sure that we have a random uuid v4. This is important for federated instances.

> **Warning**
>BREAKING Change: In order to  make every ocis storage provider ID unique by default, we needed to use a random uuidv4 during ocis init. Existing installations need to set this value explicitly or ocis will terminate after the upgrade.
> To upgrade from 2.0.0-rc.1 to 2.0.0-rc.2, 2.0.0 or later you need to set `GATEWAY_STORAGE_USERS_MOUNT_ID` and `STORAGE_USERS_MOUNT_ID` to the same random uuidv4.
>
>You can also add
>```
>storage_users:
>  mount_id: some-random-uuid
>gateway:
>  storage_registry:
>    storage_users_mount_id: some-random-uuid
>```
>to the ocis.yaml file which was created during initialisation
>
>Changing the ID of the storage-users provider will change all
>- WebDAV Urls
>- FileIDs
>- SpaceIDs
>- Bookmarks
>- and will make all existing shares invalid.
>
>The Android, Web and iOS clients will continue to work without interruptions. The Desktop Client sync connections need to be deleted and recreated.
>Sorry for the inconvenience 😅
>
>WORKAROUND - Not Recommended: You can avoid this by setting
>`GATEWAY_STORAGE_USERS_MOUNT_ID=1284d238-aa92-42ce-bdc4-0b0000009157` and
>`STORAGE_USERS_MOUNT_ID=1284d238-aa92-42ce-bdc4-0b0000009157`
>But this will cause problems later when two ocis instances want to federate.

https://github.com/owncloud/ocis/pull/5091
