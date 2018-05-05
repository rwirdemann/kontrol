#!/bin/bash
#
set -e
# Any subsequent(*) commands which fail will cause the shell script to exit immediately

export SOURCE=./
export TARGETPROGRAM=kontrol-main
export TARGETUSER=kommitment
export TARGETSERVER=94.130.79.196
export TARGETSERVER=kommitment.dyn.amicdns.de
export SSHPORT=22
export SSHSERVER="ssh -p"${SSHPORT}" $TARGETUSER@$TARGETSERVER"
export SPREADSHEET_KEY="1xkTQDGJkq9UKvZfFJTEK_W1EdM2AAy7xIFikxTCGhnk"

go test ./...

echo "setting of variables done for deployment to >$TARGETSERVER"

echo "make linux..."
make linux
echo "done."
echo
echo


# clear
echo "Deploy stuff to "${TARGETSERVER}
echo " ... "${DEPLOYMENTTARGET}
${SSHSERVER} mv ${TARGETPROGRAM} ${TARGETPROGRAM}.old
scp ${TARGETPROGRAM} ${TARGETUSER}@${TARGETSERVER}:${TARGETPROGRAM}
${SSHSERVER} ls -al
echo "done with file transfer..."
echo
echo

echo "get data file from google on "${TARGETSERVER}
echo "curl https://docs.google.com/spreadsheets/d/$SPREADSHEET_KEY/export?exportFormat=csv > '2017-Buchungen-KG - Buchungen 2017.csv'"
${SSHSERVER} "curl https://docs.google.com/spreadsheets/d/$SPREADSHEET_KEY/export?exportFormat=csv > '2017-Buchungen-KG - Buchungen 2017.csv'"
echo "done with data file..."
echo
echo

#
echo "run kontrol on "${TARGETSERVER}
set +e
${SSHSERVER} killall ${TARGETPROGRAM}
echo killed running ${TARGETPROGRAM}
${SSHSERVER} ./${TARGETPROGRAM} > /tmp/${TARGETPROGRAM}.log 2>&1 &
set -e
echo
echo


sleep 2
curl -s http://${TARGETSERVER}:8991/kontrol/version
echo
