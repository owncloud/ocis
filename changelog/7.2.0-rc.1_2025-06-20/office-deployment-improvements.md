Enhancement: Harden office deployment process

Office deployment will use proof keys by default to ensure requests to the
collaboration service come from a trusted source. In addition, OnlyOffice will
use ip filters to ensure requests come from the collaboration service (with the
exception of the editor). Lastly, the collaboration service won't be exposed
to the outside and will remain in the docker network.

https://github.com/owncloud/ocis/pull/11339
