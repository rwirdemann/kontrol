#!/bin/bash
#
export SPREADSHEET_KEY="1xkTQDGJkq9UKvZfFJTEK_W1EdM2AAy7xIFikxTCGhnk"
export SPREADSHEET_KEY_Sven="1-p7QnCtwv0TrOINDgG2v_KqfB6wtSPbNNaXQ6PX1u3g"
export FILENAME="Buchungen-KG.csv"
export URL="https://docs.google.com/spreadsheets/d/$SPREADSHEET_KEY/export?exportFormat=csv"

export CONTENTLENGTH=$(curl -sL -I $URL | grep -i Content-Length | awk '{print $2}')
export CURRENTFILESIZE=$(ls -al "$FILENAME" | awk '{print $5}')
touch "$FILENAME"

echo $URL
echo $CONTENTLENGTH
echo $CURRENTFILESIZE

case "$CONTENTLENGTH" in
  "$CURRENTFILESIZE"*)
  echo "no action taken" ;;
  *)
  echo "The remote file differs, I will get the spreadsheet from the google drive now"
  curl -sL $URL > "$FILENAME" 2>&1 /dev/null
esac
