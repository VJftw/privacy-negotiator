FROM golang

ENV JWT_SECRET changeme
ENV RABBITMQ_USER changeme
ENV RABBITMQ_PASS changeme
ENV RABBITMQ_HOSTNAME changeme
ENV PORT 80

RUN curl https://glide.sh/get | sh

RUN go get github.com/pilu/fresh
RUN go get github.com/golang/lint/golint
