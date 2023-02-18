VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-X main.version=${VERSION} \
	-X main.date=${BUILD} \
	-X main.commit=${COMMIT} \
	-X main.builtBy=Makefile

LDFLAGS_BUILD=-ldflags "${LDFLAGS}"
LDFLAGS_RELEASE=-ldflags "-s -w ${LDFLAGS}"

OUTPUT_OBJ=-o build/urlshortener-ui

MAIN_GO=./main.go

.PHONY: build
build: build_dir tidy analyze
	go build ${LDFLAGS_BUILD} ${OUTPUT_OBJ} ${MAIN_GO}

.PHONY: release
release: clean build_dir
	go build ${LDFLAGS_RELEASE} ${OUTPUT_OBJ} ${MAIN_GO}

.PHONY: clean
clean:
	rm -rf ./build
	go clean -cache

.PHONY: analyze
analyze: vet lint

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy
	go mod verify

.PHONY: lint
lint:
	golint ./...

.PHONY: build_dir
build_dir:
	mkdir -p ./build/