GOOS ?= linux
GOARCH ?= amd64
GIT_COMMIT=$(shell git rev-list -1 HEAD)
GIT_URL=$(shell git config --get remote.origin.url)

getdeps:
	go get -u -v ./...

build:	getdeps
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build -ldflags " \
	-X main.gitCommit=${GIT_COMMIT} \
	-X main.gitUrl=${GIT_URL}"

zip: build
	 zip s3-file-backup.zip s3-file-backup

clean:
	go clean