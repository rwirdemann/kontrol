#!/bin/bash
#
set -e
# Any subsequent(*) commands which fail will cause the shell script to exit immediately

export SOURCE=./
export TARGETPROGRAM=kontrol-main
export TARGETUSER=kommitment
export TARGETSERVER=94.130.79.196
export TARGETSERVER=kontrol.kommitment.biz # kommitment hetzner server
# export TARGETSERVER=kommitment.dyn.amicdns.de
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
${SSHSERVER} "[ -f ${TARGETPROGRAM} ] &&  mv ${TARGETPROGRAM} ${TARGETPROGRAM}.old"
scp ${TARGETPROGRAM} ${TARGETUSER}@${TARGETSERVER}:${TARGETPROGRAM}
scp httpsconfig.env ${TARGETUSER}@${TARGETSERVER}:.
scp ./valueMagnets/kommitmenschen.json ${TARGETUSER}@${TARGETSERVER}:./kommitmenschen.json
scp getSpreadsheet.sh ${TARGETUSER}@${TARGETSERVER}:.
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
set -x
${SSHSERVER} "./${TARGETPROGRAM} -httpPort=20181 -httpsPort=20182 -year=2018> /tmp/${TARGETPROGRAM}.log 2>&1 &"
${SSHSERVER} "./${TARGETPROGRAM} -httpPort=20171 -httpsPort=20172 -year=2017 > /tmp/${TARGETPROGRAM}-2017.log 2>&1 &"
set +x
set -e
echo
echo

#
echo "filling the crontab @reboot..."
$SSHSERVER "rm -f crontab.del"
$SSHSERVER "if crontab -l  | grep -v '{TARGETPROGRAM}' > crontab.del; then echo "crontab exists"; fi"
$SSHSERVER "echo '@reboot cd /home/$TARGETUSER; ./${TARGETPROGRAM} -httpPort=20181 -httpsPort=20182 -year=2018> /tmp/${TARGETPROGRAM}.log 2>&1 &' >> crontab.del"
$SSHSERVER "echo '@reboot cd /home/$TARGETUSER; ./${TARGETPROGRAM} -httpPort=20171 -httpsPort=20172 -year=2017 > /tmp/${TARGETPROGRAM}-2017.log 2>&1 &' >> crontab.del"
$SSHSERVER "cat crontab.del | crontab -"
$SSHSERVER "rm -f crontab.del"
$SSHSERVER crontab -l
echo " "


#
echo "getting the latest data"
${SSHSERVER}  "chmod +x ./getSpreadsheet.sh"
${SSHSERVER}  "./getSpreadsheet.sh"

sleep 2
curl -s http://${TARGETSERVER}:8991/kontrol/version
echo
