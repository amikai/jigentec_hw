GOCMD ?= go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BINARY_NAME = downloader
BINARY_UNIX = $(BINARY_NAME)_unix

TAG=1.0


default: all

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v -ldflags="-X main.VERSION=$(TAG)" jigentec.homework/cmd

test:
	$(GOTEST) -v jigentec.homework/...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run: build
	./$(BINARY_NAME)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v jigentec.homework/cmd
