#!/bin/bash
#
export SPREADSHEET_KEY="1xkTQDGJkq9UKvZfFJTEK_W1EdM2AAy7xIFikxTCGhnk"
export FILENAME="2017-Buchungen-KG - Buchungen 2017.csv"
export URL="https://docs.google.com/spreadsheets/d/$SPREADSHEET_KEY/export?exportFormat=csv"

export CONTENTLENGTH=$(curl -s -I $URL | grep -i Content-Length | awk '{print $2}')
export CURRENTFILESIZE=$(ls -al "$FILENAME" | awk '{print $5}')
touch "$FILENAME"

echo $CONTENTLENGTH
echo $CURRENTFILESIZE

case "$CONTENTLENGTH" in
  "$CURRENTFILESIZE"*)
  echo "no action taken" ;;
  *)
  echo "The remote file differs, I will get the spreadsheet from the google drive now"
  curl -s https://docs.google.com/spreadsheets/d/$SPREADSHEET_KEY/export?exportFormat=csv > "$FILENAME" 2>&1 /dev/null
esac


if pgrep -x "kontrol-main" > /dev/null
then
    echo "$(date): kontrol-main running"
else
    echo "$(date): kontrol-main stopped, restarting"
    /home/kommitment/kontrol-main > /tmp/kontrol-main.log &
fi
