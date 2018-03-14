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
	go install bitbucket.org/rwirdemann/kontrol
	go install bitbucket.org/rwirdemann/kontrol/cli

SPREADSHEET_KEY="1xkTQDGJkq9UKvZfFJTEK_W1EdM2AAy7xIFikxTCGhnk"
getdatafile:
	curl "https://docs.google.com/spreadsheets/d/$(SPREADSHEET_KEY)/export?exportFormat=csv" > "2017-Buchungen-KG - Buchungen 2017.csv"
	
.PHONY: clean build linux pi clean test run install getdatafile