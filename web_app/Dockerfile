FROM nginx:alpine

ADD https://github.com/just-containers/s6-overlay/releases/download/v1.19.1.1/s6-overlay-amd64.tar.gz /tmp/
RUN tar xzf /tmp/s6-overlay-amd64.tar.gz -C /
ENTRYPOINT ["/init"]

COPY docker/nginx.sh /etc/services.d/nginx/run
RUN chmod 554 /etc/services.d/nginx/run

COPY docker/init.sh /etc/cont-init.d/set_api
RUN chmod 554 /etc/cont-init.d/set_api

ENV API_ENDPOINT ""
ENV FB_APP_ID ""

COPY priv-neg/dist /usr/share/nginx/html
