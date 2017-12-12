#!/bin/bash
# 
set -e
# Any subsequent(*) commands which fail will cause the shell script to exit immediately

export SOURCE=./
export TARGETPROGRAM=kontrol-main   
export TARGETUSER=kommitment
export TARGETSERVER=94.130.79.196
export SSHPORT=22
export SSHSERVER="ssh -p"$SSHPORT" $TARGETUSER@$TARGETSERVER"
export SPREADSHEET_KEY="1xkTQDGJkq9UKvZfFJTEK_W1EdM2AAy7xIFikxTCGhnk"
                        

echo "setting of variables done for deployment to >$TARGETSERVER"

# cross-compile the $TARGETPROGRAM server for rpi
# env GOOS=linux GOARCH=arm GOARM=6 go build -v $TARGETPROGRAM
echo "cross compilation with: env GOOS=linux GOARCH=amd64 go build -o $TARGETPROGRAM -v $SOURCE"
env GOOS=linux GOARCH=amd64 go build -o $TARGETPROGRAM -v $SOURCE
echo "done with cross compilation..."
echo
echo


# clear
echo "Deploy stuff to "$TARGETSERVER
echo " ... "$DEPLOYMENTTARGET
$SSHSERVER mv $TARGETPROGRAM $TARGETPROGRAM.old
scp $TARGETPROGRAM $TARGETUSER@$TARGETSERVER:$TARGETPROGRAM
$SSHSERVER ls -al
echo "done with file transfer..."
echo
echo

echo "get data file from google on "$TARGETSERVER
echo "curl https://docs.google.com/spreadsheets/d/$SPREADSHEET_KEY/export?exportFormat=csv > '2017-Buchungen-KG - Buchungen 2017.csv'"
$SSHSERVER "curl https://docs.google.com/spreadsheets/d/$SPREADSHEET_KEY/export?exportFormat=csv > '2017-Buchungen-KG - Buchungen 2017.csv'"
echo "done with data file..."
echo
echo

#
echo "run kontrol on "$TARGETSERVER 
set +e
$SSHSERVER killall $TARGETPROGRAM
echo killed running $TARGETPROGRAM
$SSHSERVER ./$TARGETPROGRAM > /tmp/$TARGETPROGRAM.log 2>&1 &
set -e
echo
echo


sleep 2
curl -s http://94.130.79.196:8991/kontrol/accounts/AN | python -m json.tool

