FROM alpine

RUN apk --no-cache --update add ca-certificates

ENV TYPE API

ENV RABBITMQ_USER changeme
ENV RABBITMQ_PASS changeme
ENV RABBITMQ_HOSTNAME changeme
ENV REDIS_HOST changeme

# API only
ENV JWT_SECRET changeme
ENV PORT 80

# Worker only
ENV POSTGRES_HOST changeme
ENV POSTGRES_USER changeme
ENV POSTGRES_DBNAME changeme
ENV POSTGRES_PASSWORD changeme
ENV FACEBOOK_APP_ID changeme
ENV FACEBOOK_APP_SECRET changeme

ADD priv-neg/dist/priv-neg /priv-neg

ENTRYPOINT ["/priv-neg"]
