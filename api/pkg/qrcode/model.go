package qrcode

import "strconv"

type Result struct {
	ActaFilename string

	ActaCode   string
	CenterCode string
	Table      string

	ValidVotes   int
	NullVotes    int
	InvalidVotes int

	CandidateVotes map[string]int

	Votes map[Option]int
}

type Option struct {
	Candidate string
	Party     string
}

func (r *Result) CandidateTotals() map[string]int {
	tallies := map[string]int{}
	for opt, v := range r.Votes {
		tallies[opt.Candidate] += v
	}
	return tallies
}

func (r *Result) header() []string {
	return []string{"acta", "codigo", "centro", "mesa", "maduro", "edmundo", "martinez", "bertucci", "brito", "ecarri", "fermin", "ceballos", "marquez", "conde_pajuo", "total_validos", "total_nulo", "total_invalido"}
}

func (r *Result) asRow() []string {
	// For now write just a summary, but we should write all the data so we can dump into a DB.
	totals := r.CandidateTotals()
	nmm := totals[CandidateMaduro]
	egu := totals[CandidateGonzalez]
	lm := totals[CandidateMartinez]
	jber := totals[CandidateBertucci]
	jb := totals[CandidateBrito]
	ae := totals[CandidateEcarri]
	cf := totals[CandidateFermin]
	dc := totals[CandidateCeballos]
	em := totals[CandidateMarquez]
	ecp := totals[CandidateRausseo]
	return []string{
		r.ActaFilename,
		r.ActaCode,
		r.CenterCode,
		r.Table,
		strconv.Itoa(nmm),
		strconv.Itoa(egu),
		strconv.Itoa(lm),
		strconv.Itoa(jber),
		strconv.Itoa(jb),
		strconv.Itoa(ae),
		strconv.Itoa(cf),
		strconv.Itoa(dc),
		strconv.Itoa(em),
		strconv.Itoa(ecp),
		strconv.Itoa(r.ValidVotes),
		strconv.Itoa(r.NullVotes),
		strconv.Itoa(r.InvalidVotes),
	}
}

const (
	CandidateBertucci = "Javier Bertucci"
	CandidateBrito    = "Jose Brito"
	CandidateCeballos = "Daniel Ceballos"
	CandidateEcarri   = "Antonio Ecarri"
	CandidateFermin   = "Claudio Fermin"
	CandidateGonzalez = "Edmundo Gonzalez"
	CandidateMaduro   = "Nicolas Maduro"
	CandidateMartinez = "Luis Martinez"
	CandidateMarquez  = "Enrique Marquez"
	CandidateRausseo  = "Benjamin Rausseo"
)

var ballotOrder = []Option{
	{Candidate: CandidateMaduro, Party: "PSUV"},
	{Candidate: CandidateMaduro, Party: "PCV"},
	{Candidate: CandidateMaduro, Party: "TUPAMARO"},
	{Candidate: CandidateMaduro, Party: "PPT"},
	{Candidate: CandidateMaduro, Party: "MSV"},
	{Candidate: CandidateMaduro, Party: "PODEMOS"},
	{Candidate: CandidateMaduro, Party: "MEP"},
	{Candidate: CandidateMaduro, Party: "APC"},
	{Candidate: CandidateMaduro, Party: "ORA"},
	{Candidate: CandidateMaduro, Party: "UPV"},
	{Candidate: CandidateMaduro, Party: "EV"},
	{Candidate: CandidateMaduro, Party: "PVV"},
	{Candidate: CandidateMaduro, Party: "PFV"},
	{Candidate: CandidateMartinez, Party: "AD"},
	{Candidate: CandidateMartinez, Party: "COPEI"},
	{Candidate: CandidateMartinez, Party: "MR"},
	{Candidate: CandidateMartinez, Party: "BR"},
	{Candidate: CandidateMartinez, Party: "DDP"},
	{Candidate: CandidateMartinez, Party: "UNE"},
	{Candidate: CandidateBertucci, Party: "EL CAMBIO"},
	{Candidate: CandidateBrito, Party: "PV"},
	{Candidate: CandidateBrito, Party: "VU"},
	{Candidate: CandidateBrito, Party: "UVV"},
	{Candidate: CandidateBrito, Party: "MPJ"},
	{Candidate: CandidateEcarri, Party: "AP"},
	{Candidate: CandidateEcarri, Party: "MOVEV"},
	{Candidate: CandidateEcarri, Party: "CMC"},
	{Candidate: CandidateEcarri, Party: "FV"},
	{Candidate: CandidateEcarri, Party: "ALIANZA DEL LAPIZ"},
	{Candidate: CandidateEcarri, Party: "MIN UNIDAD"},
	{Candidate: CandidateFermin, Party: "SPV"},
	{Candidate: CandidateCeballos, Party: "VPA"},
	{Candidate: CandidateCeballos, Party: "AREPA"},
	{Candidate: CandidateGonzalez, Party: "UNTC"},
	{Candidate: CandidateGonzalez, Party: "MPV"},
	{Candidate: CandidateGonzalez, Party: "MUD"},
	{Candidate: CandidateMarquez, Party: "CENTRADOS"},
	{Candidate: CandidateRausseo, Party: "CONDE"},
}
