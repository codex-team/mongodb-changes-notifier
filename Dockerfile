##
## Build
##

FROM golang:1.17.0-buster AS build

# disable CGO
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

COPY Makefile ./

RUN make build

##
## Deploy
##

FROM alpine:3.14.1

WORKDIR /

COPY --from=build /app/bin/mongodb-changes-notifier /mongodb-changes-notifier

RUN adduser -S appuser

USER appuser

ENTRYPOINT ["/mongodb-changes-notifier"]