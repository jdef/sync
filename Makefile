PKG_PREFIX := github.com/jdef/sync/pkg
LIBS :=	\
	future \
	promise
PACKAGES := ${LIBS:%=$(PKG_PREFIX)/%}

.PHONY: format all clean build

all: build

clean:
	go clean -i -r -v  ${PACKAGES}

build:
	go build -v ${PACKAGES}

format:
	go fmt ${PACKAGES}

test:
	go test ${PACKAGES}
