.PHONY: all

BINARY="main.exe"

all: gotool build run

build:
	go build main.go

run:
	@./${BINARY}

gotool:
	go fmt ./
	go vet ./

clean:
	rm ${BINARY}

help:
	@echo "make - format go code, build project to binary file"
	@echo "make build - build project to binary file"
	@echo "make run - go run project"
	@echo "make clean - remove binary file and vim swap files"
	@echo "make gotool - run gotool, include: 'fmt' and 'vet'"
