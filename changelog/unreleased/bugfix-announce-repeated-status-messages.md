Bugfix: Announce repeated status messages

Repeating an action that shows the same status message, such as copying a
permanent link several times in a row, was only read out by the Orca screen
reader the first time. Orca skips a message whose text matches the previous
one, even when the announcement area was cleared in between.

Consecutive status messages now always differ by an inaudible trailing space,
so every repeated message is announced.

https://github.com/owncloud/ocis/pull/13261
