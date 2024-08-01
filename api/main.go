package main

import (
	"api/actas"
	"net/http"
)

func main() {
	runServer()
}

func runServer() {
	http.HandleFunc("/actas", actas.Handler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("invalid endpoint; please see https://github.com/ipince/vzlapi for docs and to submit issues"))
}
