# Deployment scenario ocis with traefik

## Setup on server

- Clone ocis repository

  `git clone https://github.com/owncloud/ocis.git`

- Copy example folder to /opt

  `cp deployment/examples/ocis_traefik /opt/`

- Overwrite OCIS_DOMAIN in .env with your.domain.com

  `sed -i 's/ocis.domain.com/your.domain.com/g' /opt/ocis_traefik/.env`

- Overwrite redirect uri with your.domain.com in identifier-registration.yml

  `sed -i 's/ocis.domain.com/your.domain.com/g' /opt/ocis_traefik/config/identifier-registration.yml`

- Change into deployment folder

  `cd /opt/ocis_traefik`

- Start application stack

  `docker-compose up -d`
