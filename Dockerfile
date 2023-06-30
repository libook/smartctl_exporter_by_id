##
## Build
##
FROM golang:1.18-alpine AS build

WORKDIR /app

ENV VERSION $VERSION
ENV BUILD_DATE $BUILD_DATE

COPY . ./
RUN go mod download
RUN go mod verify

RUN go build -o /go/bin/smartctl_exporter_by_id -ldflags "-w -s -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}"

##
## Deploy
##
FROM alpine:latest

RUN apk add --no-cache smartmontools

WORKDIR /opt/app/

COPY --from=build /go/bin/smartctl_exporter_by_id smartctl_exporter_by_id

EXPOSE 9111

ENTRYPOINT ["/opt/app/smartctl_exporter_by_id"]
