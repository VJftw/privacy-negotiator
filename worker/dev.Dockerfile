FROM golang

ENV RABBITMQ_USER changeme
ENV RABBITMQ_PASS changeme
ENV RABBITMQ_HOSTNAME changeme

RUN curl https://glide.sh/get | sh

RUN go get github.com/pilu/fresh
RUN go get github.com/golang/lint/golint
