BINARY=sb

VERSION=0.0.2
BUILD=`git log -n 1 | head -n 1 | sed -e 's/^commit //' | head -c 10`
BUILD_TIME=`date +%Y-%m-%d:%H:%M:%S`
TAG=`git describe --tags --abbrev=0`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

# Build project
build:
	go build ${LDFLAGS} -o ${BINARY}

pipeline_build:
	mkdir -p binaries
	go build ${LDFLAGS} -o binaries/${BINARY}_${TAG}_${BUILD}

# Install project
install:
	go install ${LDFLAGS}
	mv ${GOPATH}/bin/service-broker-cli ${GOPATH}/bin/sb

# Clean our project
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
