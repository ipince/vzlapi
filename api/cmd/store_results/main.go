package main

import (
	"api/db"
	"api/pkg/qrcode"
	"encoding/csv"
	"log/slog"
	"os"
	"strconv"
)

func main() {

	dbc, err := db.New()
	if err != nil {
		panic(err)
	}
	err = dbc.InitTables()
	if err != nil {
		panic(err)
	}

	f, err := os.Open("./cmd/store_results/resultados.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for i, record := range records {
		if i == 0 {
			continue
		} // skip header

		qr := &qrcode.Result{
			ActaFilename: record[0],
			CenterCode:   record[2],
			Table:        record[3],

			CandidateVotes: map[string]int{
				qrcode.CandidateMaduro:   mustInt(record[4]),
				qrcode.CandidateGonzalez: mustInt(record[5]),
				qrcode.CandidateMartinez: mustInt(record[6]),
				qrcode.CandidateBertucci: mustInt(record[7]),
				qrcode.CandidateBrito:    mustInt(record[8]),
				qrcode.CandidateEcarri:   mustInt(record[9]),
				qrcode.CandidateFermin:   mustInt(record[10]),
				qrcode.CandidateCeballos: mustInt(record[11]),
				qrcode.CandidateMarquez:  mustInt(record[12]),
				qrcode.CandidateRausseo:  mustInt(record[13]),
			},

			ValidVotes:   mustInt(record[14]),
			NullVotes:    mustInt(record[15]),
			InvalidVotes: mustInt(record[16]),
		}

		err := dbc.UpsertActa(qr)
		if err != nil {
			slog.Warn("failed to insert acta, skipping", "err", err)
			continue
		}
	}
}

func mustInt(s string) int {
	return must(func() (int, error) {
		return strconv.Atoi(s)
	})
}

func must[T any](f func() (T, error)) T {
	t, e := f()
	if e != nil {
		panic(e)
	}
	return t
}
