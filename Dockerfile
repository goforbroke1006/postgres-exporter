FROM golang:1.16 AS builder

WORKDIR /code/

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


FROM debian:stretch

COPY --from=builder /code/postgres-exporter /usr/local/bin/postgres-exporter

ENTRYPOINT [ "postgres-exporter" ]

EXPOSE 54380
