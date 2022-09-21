CGO_CPPFLAGS ?= ${CPPFLAGS}
export CGO_CPPFLAGS
CGO_CFLAGS ?= ${CFLAGS}
export CGO_CFLAGS
CGO_LDFLAGS ?= $(filter -g -L% -l% -O%,${LDFLAGS})
export CGO_LDFLAGS

EXE =
ifeq ($(GOOS),windows)
EXE = .exe
endif

## The following tasks delegate to `script/build.go` so they can be run cross-platform.

.PHONY: bin/bendsql$(EXE)
bin/bendsql$(EXE): script/build
	@script/build $@

script/build: script/build.go
	GOOS= GOARCH= GOARM= GOFLAGS= CGO_ENABLED= go build -o $@ $<

.PHONY: clean
clean: script/build
	@script/build $@

.PHONY: build

default: build

build: fmt vet
	go build -o bin/bendsql cmd/bendsql/main.go

test:
	GO111MODULE=on go test -p 1 -v -race ./...
	go vet ./...

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

DESTDIR :=
prefix  := /usr/local
bindir  := ${prefix}/bin


.PHONY: install
install: bin/bendsql
	install -d ${DESTDIR}${bindir}
	install -m755 bin/bendsql ${DESTDIR}${bindir}/

.PHONY: uninstall
uninstall:
	rm -f ${DESTDIR}${bindir}/bendsql