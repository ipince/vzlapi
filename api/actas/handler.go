package actas

import (
	"api/db"
	"api/pkg/qrcode"
	"database/sql"
	"encoding/json"
	"errors"
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
		writeErr(w, "<none>", "missing required param 'cedula'", http.StatusBadRequest)
		return
	}

	info, err := Resolve(cedula)
	if err != nil {
		writeErr(w, cedula, err.Error(), http.StatusInternalServerError)
		return
	}

	dbc, err := db.New()
	if err != nil {
		writeErr(w, cedula, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbc.Close()

	actaCode := fmt.Sprintf("%s.0%s.1.0001", info.CenterID, info.TableNumber)
	acta, err := dbc.GetActa(actaCode)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		writeErr(w, cedula, err.Error(), http.StatusInternalServerError)
		return
	}
	if acta != nil {
		fillVotes(info, acta)
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

func fillVotes(info *CedulaInfo, acta *qrcode.Result) {
	info.ActaValidVotes = acta.ValidVotes
	info.ActaNullVotes = acta.NullVotes
	info.ActaInvalidVotes = acta.InvalidVotes

	info.ActaCandidateVotes = acta.CandidateVotes
}
