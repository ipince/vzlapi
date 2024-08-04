package db

import (
	"api/pkg/qrcode"
	"database/sql"
	"encoding/csv"
	"errors"
	"os"
	"strconv"

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

func (c *Client) Close() {
	c.db.Close()
}

func (c *Client) InitTables() error {

	createActasTableSQL := `CREATE TABLE IF NOT EXISTS actas (
        code TEXT PRIMARY KEY,
        
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
        code, valid_votes, null_votes, invalid_votes,
        votes_bertucci, votes_brito, votes_ceballos, votes_ecarri, votes_fermin, votes_gonzalez, votes_maduro, votes_martinez, votes_marquez, votes_rausseo
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := c.db.Exec(insertSQL,
		acta.Code,

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
	res, err := c.GetActa(acta.Code)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if res == nil {
		return c.InsertActa(acta)
	}
	return nil // no-op, acta already in db
}

func (c *Client) GetActa(code string) (*qrcode.Result, error) {
	querySQL := `SELECT * FROM actas WHERE code = ?`

	row := c.db.QueryRow(querySQL, code)

	var qr qrcode.Result

	var c1, c2, c3, c4, c5, c6, c7, c8, c9, c10 int
	err := row.Scan(
		&qr.Code,

		&qr.ValidVotes,
		&qr.NullVotes,
		&qr.InvalidVotes,

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
func (c *Client) InitTablesMCM() error {
	createActasMariaCorinaTableSQL := `CREATE TABLE IF NOT EXISTS actas_mcm (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        CODIGO_ESTADO INTEGER,
        ESTADO TEXT,
        CODIGO_MUNICIPIO INTEGER,
        MUNICIPIO TEXT,
        COD_PARROQUIA INTEGER,
        PARROQUIA TEXT,
        CENTRO INTEGER,
        MESA INTEGER,
        RE INTEGER,
        VOTOS_VALIDOS INTEGER,
        VOTOS_NULOS INTEGER,
        EDMUNDO_GONZALEZ INTEGER,
        NICOLAS_MADURO INTEGER,
        LUIS_MARTINEZ INTEGER,
        JAVIER_BERTUCCI INTEGER,
        JOSE_BRITO INTEGER,
        ANTONIO_ECARRI INTEGER,
        CLAUDIO_FERMIN INTEGER,
        DANIEL_CEBALLOS INTEGER,
        ENRIQUE_MARQUEZ INTEGER,
        BENJAMMIN_RAUSSEO INTEGER,
        URL TEXT
    );`

	_, err := c.db.Exec(createActasMariaCorinaTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) InsertActaMariaCorina(record []string) error {
	insertSQL := `INSERT INTO actas_mcm (
        CODIGO_ESTADO, ESTADO, CODIGO_MUNICIPIO, MUNICIPIO, COD_PARROQUIA, PARROQUIA, CENTRO, MESA, RE,
        VOTOS_VALIDOS, VOTOS_NULOS, EDMUNDO_GONZALEZ, NICOLAS_MADURO, LUIS_MARTINEZ, JAVIER_BERTUCCI,
        JOSE_BRITO, ANTONIO_ECARRI, CLAUDIO_FERMIN, DANIEL_CEBALLOS, ENRIQUE_MARQUEZ, BENJAMMIN_RAUSSEO, URL
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	codigoEstado, _ := strconv.Atoi(record[0])
	codigoMunicipio, _ := strconv.Atoi(record[2])
	codParroquia, _ := strconv.Atoi(record[4])
	centro, _ := strconv.Atoi(record[6])
	mesa, _ := strconv.Atoi(record[7])
	re, _ := strconv.Atoi(record[8])
	votosValidos, _ := strconv.Atoi(record[9])
	votosNulos, _ := strconv.Atoi(record[10])
	edmundoGonzalez, _ := strconv.Atoi(record[11])
	nicolasMaduro, _ := strconv.Atoi(record[12])
	luisMartinez, _ := strconv.Atoi(record[13])
	javierBertucci, _ := strconv.Atoi(record[14])
	joseBrito, _ := strconv.Atoi(record[15])
	antonioEcarri, _ := strconv.Atoi(record[16])
	claudioFermin, _ := strconv.Atoi(record[17])
	danielCeballos, _ := strconv.Atoi(record[18])
	enriqueMarquez, _ := strconv.Atoi(record[19])
	benjamminRausseo, _ := strconv.Atoi(record[20])
	url := record[21]

	_, err := c.db.Exec(insertSQL, codigoEstado, record[1], codigoMunicipio, record[3], codParroquia,
		record[5], centro, mesa, re, votosValidos, votosNulos, edmundoGonzalez, nicolasMaduro,
		luisMartinez, javierBertucci, joseBrito, antonioEcarri, claudioFermin, danielCeballos,
		enriqueMarquez, benjamminRausseo, url)
	return err
}

func (c *Client) UpsertActaMariaCorina(record []string) error {

	return c.InsertActaMariaCorina(record)
}

func (c *Client) ImportCSVForMariaCorina() error {
	file, err := os.Open("./actas_mcm.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		err := c.UpsertActaMariaCorina(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCenterData(dbc *Client, centro int) (*sql.Rows, error) {
	query := `SELECT CODIGO_ESTADO, ESTADO, CODIGO_MUNICIPIO, MUNICIPIO, 
	COD_PARROQUIA, PARROQUIA, MESA, VOTOS_VALIDOS, VOTOS_NULOS, 
	EDMUNDO_GONZALEZ, NICOLAS_MADURO, LUIS_MARTINEZ, JAVIER_BERTUCCI, 
	JOSE_BRITO, ANTONIO_ECARRI, CLAUDIO_FERMIN, DANIEL_CEBALLOS, 
	ENRIQUE_MARQUEZ, BENJAMMIN_RAUSSEO, URL 
	FROM actas_mcm WHERE CENTRO = ?`
	dbc.db.Query(query, centro)
	return dbc.db.Query(query, centro)
}
