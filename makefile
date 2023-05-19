GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test
BINARY_NAME = md-doc
IMAGE_NAME = md-doc

install:
	$(GOMOD) tidy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME) .

image:
	docker build -t $(IMAGE_NAME) .

clean:
	rm -rf bin/$(BINARY_NAME)