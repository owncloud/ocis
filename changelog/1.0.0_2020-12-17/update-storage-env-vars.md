Change: Clarify storage driver env vars

After renaming ocsi-reva to storage and combining the storage and data providers some env vars were confusingly named `STORAGE_STORAGE_...`. We are changing the prefix for driver related env vars to `STORAGE_DRIVER_...`. This makes changing the storage driver using eg.: `STORAGE_HOME_DRIVER=eos` and setting driver options using `STORAGE_DRIVER_EOS_LAYOUT=...` less confusing.

https://github.com/owncloud/ocis/pull/729