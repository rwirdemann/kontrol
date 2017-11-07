package util

import (
	"encoding/json"
	"fmt"
	"log"
)

func Json(entities interface{}) string {
	var b []byte
	var err error
	if b, err = json.Marshal(entities); err == nil {
		json := fmt.Sprintf("%s", string(b[:]))
		return json
	}
	log.Fatal(err)
	return ""
}
