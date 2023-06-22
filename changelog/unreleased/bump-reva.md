Enhancement: Update reva

Update reva to latest edge

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

https://github.com/owncloud/ocis/pull/6529
https://github.com/owncloud/ocis/pull/6544
https://github.com/owncloud/ocis/pull/6507
https://github.com/owncloud/ocis/pull/6572
https://github.com/owncloud/ocis/pull/6590
