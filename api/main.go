package main

import (
	"api/actas"
	"api/centers"
	"net/http"
)

func main() {

	//actas.Resolve("4000000")
	runServer()
}

func runServer() {
	http.HandleFunc("/actas", actas.Handler)
	// example: http://localhost:8080/actas/centros?centro=10113042
	http.HandleFunc("/actas/centros", centers.Handler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("invalid endpoint; please see https://github.com/ipince/vzlapi for docs and to submit issues"))
}
