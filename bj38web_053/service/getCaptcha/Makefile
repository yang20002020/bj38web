
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc -I=. --micro_out=. --go_out=. proto/getCaptcha/getCaptcha.proto

.PHONY: build
build: proto

	go build -o getCaptcha-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t getCaptcha-srv:latest
