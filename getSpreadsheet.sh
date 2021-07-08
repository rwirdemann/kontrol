#!/bin/bash
#

#https://docs.google.com/spreadsheets/d/1xkTQDGJkq9UKvZfFJTEK_W1EdM2AAy7xIFikxTCGhnk/export?exportFormat=csv

export SPREADSHEET_KEY="1xkTQDGJkq9UKvZfFJTEK_W1EdM2AAy7xIFikxTCGhnk"
export FILENAME="Buchungen-KG.csv"
export URL="https://docs.google.com/spreadsheets/d/$SPREADSHEET_KEY/export?exportFormat=csv"

touch "$FILENAME"
export CONTENTLENGTH=$(curl -s -I $URL | grep -i Content-Length | awk '{print $2}')
export CURRENTFILESIZE=$(ls -al "$FILENAME" | awk '{print $5}')

echo $CONTENTLENGTH
echo $CURRENTFILESIZE
echo "getting file from $URL"

case "$CONTENTLENGTH" in
  "$CURRENTFILESIZE"*)
  echo "no action taken" ;;
  *)
  echo "The remote file differs, I will get the spreadsheet from the google drive now"
  echo "curl -s -L $URL"
  curl -s -L $URL > "$FILENAME" 2>&1 /dev/null
esac


if pgrep -x "kontrol-main" > /dev/null
then
    echo "$(date): kontrol-main running"
else
    echo "$(date): kontrol-main stopped, restarting"
fi
