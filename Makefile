all: build test

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

test:
	go test -v ./...

image:
	docker build -t docker.io/goforbroke1006/postgres-exporter:latest ./
	docker push docker.io/goforbroke1006/postgres-exporter:latest