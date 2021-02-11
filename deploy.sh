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
scp ./valueMagnets/kommitment.json ${TARGETUSER}@${TARGETSERVER}:./kommitment.json
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
${SSHSERVER} "./${TARGETPROGRAM} -httpPort=20211 -httpsPort=20212 -year=2021> /tmp/${TARGETPROGRAM}-2021.log 2>&1 &"
${SSHSERVER} "./${TARGETPROGRAM} -httpPort=20201 -httpsPort=20202 -year=2020> /tmp/${TARGETPROGRAM}-2020.log 2>&1 &"
${SSHSERVER} "./${TARGETPROGRAM} -httpPort=20191 -httpsPort=20192 -year=2019> /tmp/${TARGETPROGRAM}-2019.log 2>&1 &"
${SSHSERVER} "./${TARGETPROGRAM} -httpPort=20181 -httpsPort=20182 -year=2018> /tmp/${TARGETPROGRAM}-2018.log 2>&1 &"
${SSHSERVER} "./${TARGETPROGRAM} -httpPort=20171 -httpsPort=20172 -year=2017 > /tmp/${TARGETPROGRAM}-2017.log 2>&1 &"
set +x
set -e
echo
echo

#   make sure to have a dayly reboot...
#   0 4   *   *   *    /sbin/shutdown -r +5
echo "filling the crontab @reboot..."
$SSHSERVER "rm -f crontab.del"
$SSHSERVER "if crontab -l  | grep -v '{TARGETPROGRAM}' | grep -v '@reboot' > crontab.del; then echo "crontab exists"; fi"
$SSHSERVER "echo '@reboot cd /home/$TARGETUSER; ./${TARGETPROGRAM} -httpPort=20211 -httpsPort=20212 -year=2021> /tmp/${TARGETPROGRAM}-2021.log 2>&1 &' >> crontab.del"
$SSHSERVER "echo '@reboot cd /home/$TARGETUSER; ./${TARGETPROGRAM} -httpPort=20201 -httpsPort=20202 -year=2020> /tmp/${TARGETPROGRAM}-2020.log 2>&1 &' >> crontab.del"
$SSHSERVER "echo '@reboot cd /home/$TARGETUSER; ./${TARGETPROGRAM} -httpPort=20191 -httpsPort=20192 -year=2019> /tmp/${TARGETPROGRAM}-2019.log 2>&1 &' >> crontab.del"
$SSHSERVER "echo '@reboot cd /home/$TARGETUSER; ./${TARGETPROGRAM} -httpPort=20181 -httpsPort=20182 -year=2018> /tmp/${TARGETPROGRAM}-2018.log 2>&1 &' >> crontab.del"
$SSHSERVER "echo '@reboot cd /home/$TARGETUSER; ./${TARGETPROGRAM} -httpPort=20171 -httpsPort=20172 -year=2017 > /tmp/${TARGETPROGRAM}-2017.log 2>&1 &' >> crontab.del"
$SSHSERVER "echo '@reboot cd /home/kommitment; ./simpleServer -httpPort=8042 -httpsPort=8043 > /dev/null 2>&1 &' >> crontab.del"
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
