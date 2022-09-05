all: build
start: restapi
clean: rm restapi

build:
	go build
