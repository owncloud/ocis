# Audit service

The audit service logs all events of the system as an audit log. Per default, it will be logged to standard out, but can also be configured to a file output. Supported log formats are json or a simple key-value pair ("key1=value1 key2=value2").

With audit logs you are able to prove compliance with corporate guidelines as well as to enable reporting and auditing of operations. The audit service takes note of actions conducted by users and administrators.

The service is not started automatically when running as single binary started via `ocis server` or when running as docker container and must be started and stopped manually on demand.

Specifically, the audit service logs:

-   File system operations (create/delete/move; including actions on the trash bin and versioning)
-   User management operations (creation/deletion of users)
-   Sharing operations (user/group sharing, sharing via link, changing permissions, calls to sharing API from clients)
