FROM golang:1.9 as build

RUN cd / && go get -v github.com/Jeffail/gabs github.com/influxdata/influxdb/client/v2
COPY main.go /main.go
RUN cd /&& CGO_ENABLED=0 GOOS=linux go build -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -s" -o c2i main.go

FROM ubuntu as cert
RUN apt-get update && apt-get install ca-certificates -y \
  && mkdir -p /usr/local/share/ca-certificates/

FROM scratch
LABEL maintainer "Jan Garaj <jan.garaj@gmail.com>"

ENV \
  APP_PORT=80 \
  INFLUXDB_URL=http://localhost:8086 \
  INFLUXDB_USERNAME= \
  INFLUXDB_PASSWORD= \
  INFLUXDB_DB_DATA=catchpoint-data \
  INFLUXDB_DB_ALERT=catchpoint-alert

CMD ["/c2i"]
COPY --from=cert /etc/ssl/certs /etc/ssl/certs
COPY --from=build /c2i /

