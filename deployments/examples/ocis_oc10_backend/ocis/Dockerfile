FROM owncloud/ocis:latest

ARG OCIS_DOMAIN
ENV OCIS_DOMAIN $OCIS_DOMAIN
ARG OC10_DOMAIN
ENV OC10_DOMAIN $OC10_DOMAIN

ADD config /config
RUN sed -i "s/ocis.example.org/${OCIS_DOMAIN}/g" /config/web/config.json \
    && sed -i "s/ocis.example.org/${OCIS_DOMAIN}/g" /config/identifier-registration.yaml \
    && sed -i "s/oc10.example.org/${OC10_DOMAIN}/g" /config/identifier-registration.yaml \
