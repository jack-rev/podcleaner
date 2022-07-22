# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build

RUN mv podcleaner /usr/local/bin
COPY run.sh /usr/local/bin/

RUN chmod +x /usr/local/bin/run.sh

CMD ["run.sh"]
