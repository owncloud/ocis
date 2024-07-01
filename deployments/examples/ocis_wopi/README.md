Please refer to our [admin documentation](https://doc.owncloud.com/ocis/latest/depl-examples/ubuntu-compose/ubuntu-compose-prod.html) for instructions on how to deploy this scenario.

Note: This deployment setup is highly configurable. At minimum, it starts `traefik`, `ocis`, `tika`, the `wopiserver` and `collabora`. Additional services can be started by removing the respective comment in the `.env` file. Depending on the service added, related variables need to be configured.

