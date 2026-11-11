Bugfix: Record the acting user in the audit log's User field

The audit log's `User` field is documented as the UID of the user performing
the action, but the file and folder converters filled it from the event's
`Owner`, which reva does not set, so the field was empty on every upload,
download, move, delete, restore and version restore. Where an event did carry
an owner, `User` named the owner rather than the person acting, for example a
share recipient uploading into a shared folder. Space, user and group events
hard-coded the field to an empty string although every one of those events
carries the actor.

`User` is now taken from the event's `Executant` for all of these events. The
`Owner` field of file events is unchanged. The failed public-link access event
still has no actor and keeps an empty `User`.

While adding tests for the user and group events, the group member added and
removed messages turned out to pass the user and group ids in the wrong order,
naming the group as the user and the user as the group; that is fixed too.

https://github.com/owncloud/ocis/issues/4789
