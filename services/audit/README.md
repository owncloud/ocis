# Audit service

The audit service logs all events of the system into an audit log. To be able to prove compliance with corporate guidelines as well as to enable reporting and auditing of operations, the Auditing extension takes note of actions conducted by users and administrators. Per default it will be logged to standard out but can also be configured to a file output. Supported log formats are json or a simple key value ("key1=value1 key2=value2").

The service is not started automatically when running `ocis server` (single binary setup), it has to be started explicitly.

Specifically, the application logs

-   file system operations (create/delete/move; including actions on the trash bin and versioning)
-   user management operations (creation/deletion of users)
-   sharing operations (user/group sharing, sharing via link, changing permissions, calls to sharing API from clients)
