package centers

import (
	"api/db"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type ResponseErr struct {
	Error string `json:"error"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	centroStr := r.URL.Query().Get("centro")
	if centroStr == "" {
		writeErr(w, "<none>", "missing required param 'centro'", http.StatusBadRequest)
		return
	}

	centro, err := strconv.Atoi(centroStr)
	if err != nil {
		writeErr(w, centroStr, "invalid 'centro' parameter", http.StatusBadRequest)
		return
	}

	dbc, err := db.New()
	if err != nil {
		writeErr(w, centroStr, "database connection error", http.StatusInternalServerError)
		return
	}
	defer dbc.Close()

	centerInfo, err := GetCenterInfo(dbc, centro)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeErr(w, centroStr, "no data found for the given centro", http.StatusNotFound)
		} else {
			writeErr(w, centroStr, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonResponse, err := json.Marshal(centerInfo)
	if err != nil {
		writeErr(w, centroStr, "error encoding response", http.StatusInternalServerError)
		return
	}

	slog.Info(fmt.Sprintf("successful request for centro %s", centroStr))
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func writeErr(w http.ResponseWriter, centro string, msg string, code int) {
	slog.Error(fmt.Sprintf("failed request for centro %s: %s", centro, msg))
	resp, _ := json.Marshal(ResponseErr{Error: msg})
	w.WriteHeader(code)
	w.Write(resp)
}