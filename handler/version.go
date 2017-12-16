package handler

import "net/http"

func MakeVersionHandler(githash string, buidstamp string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
