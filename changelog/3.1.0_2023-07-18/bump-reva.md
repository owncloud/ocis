Enhancement: Update reva to v2.15.0

*   Bugfix [cs3org/reva#4004](https://github.com/cs3org/reva/pull/4004): Add path to public link POST
*   Bugfix [cs3org/reva#3993](https://github.com/cs3org/reva/pull/3993): Add token to LinkAccessedEvent
*   Bugfix [cs3org/reva#4007](https://github.com/cs3org/reva/pull/4007): Close archive writer properly
*   Bugfix [cs3org/reva#3982](https://github.com/cs3org/reva/pull/3982): Fixed couple of smaller space lookup issues
*   Bugfix [cs3org/reva#4003](https://github.com/cs3org/reva/pull/4003): Don't connect ldap on startup
*   Bugfix [cs3org/reva#4032](https://github.com/cs3org/reva/pull/4032): Temporarily exclude ceph-iscsi when building revad-ceph image
*   Bugfix [cs3org/reva#4042](https://github.com/cs3org/reva/pull/4042): Fix writing 0 byte msgpack metadata
*   Bugfix [cs3org/reva#3970](https://github.com/cs3org/reva/pull/3970): Fix enforce-password issue
*   Bugfix [cs3org/reva#4057](https://github.com/cs3org/reva/pull/4057): Properly handle not-found errors when getting a public share
*   Bugfix [cs3org/reva#4048](https://github.com/cs3org/reva/pull/4048): Fix messagepack propagation
*   Bugfix [cs3org/reva#4056](https://github.com/cs3org/reva/pull/4056): Fix destroys data destination when moving issue
*   Bugfix [cs3org/reva#4012](https://github.com/cs3org/reva/pull/4012): Fix mtime if 0 size file uploaded
*   Bugfix [cs3org/reva#4010](https://github.com/cs3org/reva/pull/4010): Omit spaceroot when archiving
*   Bugfix [cs3org/reva#4047](https://github.com/cs3org/reva/pull/4047): Publish events synchrously
*   Bugfix [cs3org/reva#4039](https://github.com/cs3org/reva/pull/4039): Restart Postprocessing
*   Bugfix [cs3org/reva#3963](https://github.com/cs3org/reva/pull/3963): Treesize interger overflows
*   Bugfix [cs3org/reva#3943](https://github.com/cs3org/reva/pull/3943): When removing metadata always use correct database and table
*   Bugfix [cs3org/reva#3978](https://github.com/cs3org/reva/pull/3978): Decomposedfs no longer os.Stats when reading node metadata
*   Bugfix [cs3org/reva#3959](https://github.com/cs3org/reva/pull/3959): Drop unnecessary stat
*   Bugfix [cs3org/reva#3948](https://github.com/cs3org/reva/pull/3948): Handle the bad request status
*   Bugfix [cs3org/reva#3955](https://github.com/cs3org/reva/pull/3955): Fix panic
*   Bugfix [cs3org/reva#3977](https://github.com/cs3org/reva/pull/3977): Prevent direct access to trash items
*   Bugfix [cs3org/reva#3933](https://github.com/cs3org/reva/pull/3933): Concurrently invalidate mtime cache in jsoncs3 share manager
*   Bugfix [cs3org/reva#3985](https://github.com/cs3org/reva/pull/3985): Reduce jsoncs3 lock congestion
*   Bugfix [cs3org/reva#3960](https://github.com/cs3org/reva/pull/3960): Add trace span details
*   Bugfix [cs3org/reva#3951](https://github.com/cs3org/reva/pull/3951): Link context in metadata client
*   Bugfix [cs3org/reva#3950](https://github.com/cs3org/reva/pull/3950): Use plain otel tracing in metadata client
*   Bugfix [cs3org/reva#3975](https://github.com/cs3org/reva/pull/3975): Decomposedfs now resolves the parent without an os.Stat
*   Change [cs3org/reva#3947](https://github.com/cs3org/reva/pull/3947): Bump golangci-lint to 1.51.2
*   Change [cs3org/reva#3945](https://github.com/cs3org/reva/pull/3945): Revert golangci-lint back to 1.50.1
*   Enhancement [cs3org/reva#3966](https://github.com/cs3org/reva/pull/3966): Add space metadata to ocs shares list
*   Enhancement [cs3org/reva#3953](https://github.com/cs3org/reva/pull/3953): Client selector pool
*   Enhancement [cs3org/reva#3941](https://github.com/cs3org/reva/pull/3941): Adding tracing for jsoncs3
*   Enhancement [cs3org/reva#3965](https://github.com/cs3org/reva/pull/3965): ResumePostprocessing Event
*   Enhancement [cs3org/reva#3981](https://github.com/cs3org/reva/pull/3981): We have updated the UserFeatureChangedEvent to reflect value changes
*   Enhancement [cs3org/reva#3986](https://github.com/cs3org/reva/pull/3986): Allow disabling wopi chat
*   Enhancement [cs3org/reva#4060](https://github.com/cs3org/reva/pull/4060): We added a go-micro based app-provider registry
*   Enhancement [cs3org/reva#4013](https://github.com/cs3org/reva/pull/4013): Add new WebDAV permissions
*   Enhancement [cs3org/reva#3987](https://github.com/cs3org/reva/pull/3987): Cache space indexes
*   Enhancement [cs3org/reva#3973](https://github.com/cs3org/reva/pull/3973): More logging for metadata propagation
*   Enhancement [cs3org/reva#4059](https://github.com/cs3org/reva/pull/4059): Improve space index performance
*   Enhancement [cs3org/reva#3994](https://github.com/cs3org/reva/pull/3994): Load matching spaces concurrently
*   Enhancement [cs3org/reva#4049](https://github.com/cs3org/reva/pull/4049): Do not invalidate filemetadata cache early
*   Enhancement [cs3org/reva#4040](https://github.com/cs3org/reva/pull/4040): Allow to use external trace provider in micro service
*   Enhancement [cs3org/reva#4019](https://github.com/cs3org/reva/pull/4019): Allow to use external trace provider
*   Enhancement [cs3org/reva#4045](https://github.com/cs3org/reva/pull/4045): Log error message in grpc interceptor
*   Enhancement [cs3org/reva#3989](https://github.com/cs3org/reva/pull/3989): Parallelization of jsoncs3 operations
*   Enhancement [cs3org/reva#3809](https://github.com/cs3org/reva/pull/3809): Trace decomposedfs syscalls
*   Enhancement [cs3org/reva#4067](https://github.com/cs3org/reva/pull/4067): Trace upload progress
*   Enhancement [cs3org/reva#3887](https://github.com/cs3org/reva/pull/3887): Trace requests through datagateway
*   Enhancement [cs3org/reva#4052](https://github.com/cs3org/reva/pull/4052): Update go-ldap to v3.4.5
*   Enhancement [cs3org/reva#4065](https://github.com/cs3org/reva/pull/4065): Upload directly to dataprovider
*   Enhancement [cs3org/reva#4046](https://github.com/cs3org/reva/pull/4046): Use correct tracer name
*   Enhancement [cs3org/reva#3986](https://github.com/cs3org/reva/pull/3986): Allow disabling wopi chat writer properly

https://github.com/owncloud/ocis/pull/6829
https://github.com/owncloud/ocis/pull/6529
https://github.com/owncloud/ocis/pull/6544
https://github.com/owncloud/ocis/pull/6507
https://github.com/owncloud/ocis/pull/6572
https://github.com/owncloud/ocis/pull/6590
https://github.com/owncloud/ocis/pull/6812
