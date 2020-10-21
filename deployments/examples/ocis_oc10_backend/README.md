# Deployment scenario ocis with oc10 backend and traefik as reverse proxy

## Setup on server

* Clone ocis repository

  `git clone https://github.com/owncloud/ocis.git`

* Copy example folder to /opt
  `cp deployment/examples/ocis_oc10_backend /opt/`

* Overwrite OCIS_DOMAIN and OC10_DOMAIN in .env with your-ocis.domain.com and your-oc10.domain.com

  `sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/ocis_oc10_backend/.env`

  `sed -i 's/oc10.domain.com/your-oc10.domain.com/g' /opt/ocis_oc10_backend/.env`

* Overwrite redirect uris with your-ocis.domain.com and your-oc10.domain.com in identifier-registration.yml

  `sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/ocis_oc10_backend/ocis/identifier-registration.yml`

  `sed -i 's/oc10.domain.com/your-oc10.domain.com/g' /opt/ocis_oc10_backend/ocis/identifier-registration.yml`

* Change into deployment folder

  `cd /opt/ocis_oc10_backend`

* Start application stack

  `docker-compose up -d`
