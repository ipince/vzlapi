package actas

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ResponseErr struct {
	Error string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cedula := r.URL.Query().Get("cedula")
	if cedula == "" {
		writeErr(w, "<none", "missing required param 'cedula'", http.StatusBadRequest)
		return
	}

	info, err := Handle(cedula)
	if err != nil {
		writeErr(w, cedula, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(info)
	if err != nil {
		writeErr(w, cedula, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info(fmt.Sprintf("successful request for %s", cedula))
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func writeErr(w http.ResponseWriter, cedula string, msg string, code int) {
	slog.Error(fmt.Sprintf("failed request for cedula %s: %s", cedula, msg))
	resp, _ := json.Marshal(ResponseErr{Error: msg})
	w.WriteHeader(code)
	w.Write(resp)
}

func Handle(cedula string) (*CedulaInfo, error) {
	info, err := Resolve(cedula)
	if err != nil {
		return nil, err
	}

	FillUrls(info)

	return info, nil
}
