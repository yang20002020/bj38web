
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc -I=. --micro_out=. --go_out=. proto/house/house.proto

.PHONY: build
build: proto

	go build -o house-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t house-srv:latest
