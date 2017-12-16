BINARY=kontrol-main

VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
ENVLINUX=env GOOS=linux GOARCH=amd64
ENVPI=env GOOS=linux GOARCH=arm GOARM=6

build:
	go build ${LDFLAGS} -o ${BINARY}

linux:
	${ENVLINUX} go build ${LDFLAGS} -o ${BINARY}

pi:
	${ENVPI} go build ${LDFLAGS} -o ${BINARY}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

test:
	go test ./...
	
.PHONY: clean build