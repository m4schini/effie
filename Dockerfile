# Build Stage
FROM golang:1.18-alpine as build

ENV APP_NAME effie

RUN mkdir /src
COPY . /src
WORKDIR /src

RUN go mod tidy
RUN go build -v -o /$APP_NAME /src/

# Run Stage
FROM alpine:latest

ENV APP_NAME effie

COPY --from=build /$APP_NAME .

RUN mkdir /config

ENTRYPOINT ./$APP_NAME