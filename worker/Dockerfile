FROM alpine

RUN apk --no-cache --update add ca-certificates

ENV RABBITMQ_USER changeme
ENV RABBITMQ_PASS changeme
ENV RABBITMQ_HOSTNAME changeme

ENV POSTGRES_HOST changeme
ENV POSTGRES_USER changeme
ENV POSTGRES_DBNAME changeme
ENV POSTGRES_PASSWORD changeme

ADD priv-neg/dist/priv-neg /priv-neg

ENTRYPOINT ["/priv-neg"]
