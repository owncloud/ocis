# Deployment scenario ocis with external konnectd service on separate node and traefik as reverse proxy

## Setup on ocis server

- Clone ocis repository

  `git clone https://github.com/owncloud/ocis.git`

- Copy example sub folder for ocisnode to /opt

  `cp deployment/examples/ocis_external_konnectd/ocisnode /opt/`

- Overwrite OCIS_DOMAIN and IDP_DOMAIN in .env with your-ocis.domain.com and your-idp.domain.com

  `sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/ocisnode/.env`

  `sed -i 's/idp.domain.com/your-idp.domain.com/g' /opt/ocisnode/.env`

- Change into deployment folder

  `cd /opt/ocisnode`

- Start application stack

  `docker-compose up -d`

## Setup on idp server

- Clone ocis repository

  `git clone https://github.com/owncloud/ocis.git`

- Copy example sub folder for idpnode to /opt

  `cp deployment/examples/ocis_external_konnectd/idpnode /opt/`

- Overwrite OCIS_DOMAIN and IDP_DOMAIN in .env with your-ocis.domain.com and your-idp.domain.com

  `sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/idpnode/.env`

  `sed -i 's/idp.domain.com/your-idp.domain.com/g' /opt/idpnode/.env`

- Overwrite redirect uri with your-ocis.domain.com in identifier-registration.yml

  `sed -i 's/ocis.domain.com/your-ocis.domain.com/g' /opt/idpnode/config/identifier-registration.yml
  `

- Change into deployment folder

  `cd /opt/idpnode`

- Start application stack

  `docker-compose up -d`
