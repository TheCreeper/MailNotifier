#!/usr/bin/make -f

SHELL=/bin/sh
bin=bin
name=mnd

all: build

build:
	go build -v -o bin/$(name)

windows:
	GOOS=windows go build -v -o bin/$(name).exe

clean:
	go clean -x

remove:
	go clean -i