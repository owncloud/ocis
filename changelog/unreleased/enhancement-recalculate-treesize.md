Enhancement: Add a command to recalculate the treesize of a space

The reported quota usage of a space can drift from the data actually stored on
disk. The treesize of a directory is maintained as a running counter that is
updated on every change, so a failed or incomplete update leaves the counter
wrong. Removing a node without propagating the size change, for example, leaves
its size counted against the quota forever, with no way to reclaim it.

The new command walks a space bottom up, recalculates the treesize of every
directory from the actual size of its children and corrects the stored values:

```
ocis storage-users spaces recalculate-treesize --space-id <id>
ocis storage-users spaces recalculate-treesize --space-id <id> --dry-run=false
```

Omit `--space-id` to process all spaces. As with the other maintenance commands,
`--dry-run` defaults to true, so the command reports what it would change until
it is explicitly told to write.

https://github.com/owncloud/ocis/pull/12746
