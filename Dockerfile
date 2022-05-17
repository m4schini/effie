# Build Stage
FROM golang:1.18 as build

ENV APP_NAME effie

RUN mkdir /src
COPY . /src
WORKDIR /src

RUN go mod tidy
RUN go build -v -o /$APP_NAME /src/

# Run Stage
FROM golang:1.18

ENV APP_NAME effie

COPY --from=build /$APP_NAME .

RUN touch ./blocked_summoners.txt

ENTRYPOINT ./$APP_NAME