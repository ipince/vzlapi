package db

import (
	"api/pkg/qrcode"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	db *sql.DB
}

func New() (*Client, error) {
	// Connect to SQLite database (creates it if it doesn't exist)
	db, err := sql.Open("sqlite3", "./db/db.sqlite")
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) InitTables() error {

	createActasTableSQL := `CREATE TABLE IF NOT EXISTS actas (
        filename TEXT PRIMARY KEY,
        
        center_code TEXT,
        table_number TEXT,
        
        valid_votes INTEGER,
        null_votes INTEGER,
        invalid_votes INTEGER,
        
        votes_ceballos INTEGER,
        votes_bertucci INTEGER,
        votes_brito INTEGER,
        votes_ecarri INTEGER,
        votes_fermin INTEGER,
        votes_gonzalez INTEGER,
        votes_maduro INTEGER,
        votes_martinez INTEGER,
        votes_marquez INTEGER,
        votes_rausseo INTEGER
    );`

	_, err := c.db.Exec(createActasTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) InsertActa(acta *qrcode.Result) error {
	insertSQL := `INSERT INTO actas (
        filename, center_code, table_number, valid_votes, null_votes, invalid_votes,
        votes_bertucci, votes_brito, votes_ceballos, votes_ecarri, votes_fermin, votes_gonzalez, votes_maduro, votes_martinez, votes_marquez, votes_rausseo
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := c.db.Exec(insertSQL,
		acta.ActaFilename,
		acta.CenterCode,
		acta.Table,

		acta.ValidVotes,
		acta.NullVotes,
		acta.InvalidVotes,

		acta.CandidateVotes[qrcode.CandidateBertucci],
		acta.CandidateVotes[qrcode.CandidateBrito],
		acta.CandidateVotes[qrcode.CandidateCeballos],
		acta.CandidateVotes[qrcode.CandidateEcarri],
		acta.CandidateVotes[qrcode.CandidateFermin],
		acta.CandidateVotes[qrcode.CandidateGonzalez],
		acta.CandidateVotes[qrcode.CandidateMaduro],
		acta.CandidateVotes[qrcode.CandidateMartinez],
		acta.CandidateVotes[qrcode.CandidateMarquez],
		acta.CandidateVotes[qrcode.CandidateRausseo],
	)
	return err

}

func (c *Client) UpsertActa(acta *qrcode.Result) error {
	res, err := c.GetActa(acta.ActaFilename)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if res == nil {
		return c.InsertActa(acta)
	}
	return nil // no-op, acta already in db
}

func (c *Client) GetActa(filename string) (*qrcode.Result, error) {
	querySQL := `SELECT * FROM actas WHERE filename = ?`

	row := c.db.QueryRow(querySQL, filename)

	var qr qrcode.Result

	var c1, c2, c3, c4, c5, c6, c7, c8, c9, c10 int
	err := row.Scan(
		&qr.ActaFilename,

		&qr.CenterCode,
		&qr.Table,

		&qr.ValidVotes,
		&qr.ValidVotes,
		&qr.ValidVotes,

		&c1, &c2, &c3, &c4, &c5, &c6, &c7, &c8, &c9, &c10,
	)
	if err != nil {
		return nil, err
	}

	qr.CandidateVotes = map[string]int{
		qrcode.CandidateBertucci: c1,
		qrcode.CandidateBrito:    c2,
		qrcode.CandidateCeballos: c3,
		qrcode.CandidateEcarri:   c4,
		qrcode.CandidateFermin:   c5,
		qrcode.CandidateGonzalez: c6,
		qrcode.CandidateMaduro:   c7,
		qrcode.CandidateMartinez: c8,
		qrcode.CandidateMarquez:  c9,
		qrcode.CandidateRausseo:  c10,
	}

	return &qr, nil
}
