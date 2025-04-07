BINARY=stoglr
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.versionValue=${VERSION} -X main.buildTime=${BUILD}"

build: generate-templ build-pkg

check: generate-templ build-no-flags test vet

build-no-flags:
	go build -o ${BINARY}-no-flags
build-pkg:
	go build ${LDFLAGS} -o ${BINARY}
generate-templ:
	templ generate
run:
	./${BINARY}
test:
	go test ./... -coverprofile=coverage.out
vet:
	go vet ./...
fmt:
	go fmt ./...
clean:
	rm -rf ${BINARY} ${BINARY}-no-flags coverage.out