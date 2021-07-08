MAIN=main.go
MAIN_CLI=cli/main.go
BINARY=kontrol-main

VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-X main.githash=${VERSION} -X main.buildstamp=${BUILD}"
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

run:
	go run ${MAIN} > /tmp/${MAIN}.log 2>&1 &

install:
	go install github.com/ahojsenn/kontrol
	go install github.com/ahojsenn/kontrol

SPREADSHEET_KEY="1xkTQDGJkq9UKvZfFJTEK_W1EdM2AAy7xIFikxTCGhnk"
getdatafile:
	curl "https://docs.google.com/spreadsheets/d/$(SPREADSHEET_KEY)/export?exportFormat=csv" > "Buchungen-KG.csv"

.PHONY: clean build linux pi clean test run install getdatafile
