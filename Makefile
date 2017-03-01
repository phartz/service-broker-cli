BINARY=sb

VERSION=0.0.1
BUILD=`git log -n 1 | head -n 1 | sed -e 's/^commit //' | head -c 10`
BUILD_TIME=`date +%Y-%m-%d:%H:%M:%S`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

# Build project
build:
	go build ${LDFLAGS} -o ${BINARY}

# Install project
install:
	go install ${LDFLAGS}

# Clean our project
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install