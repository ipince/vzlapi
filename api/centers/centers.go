package centers

import (
	"api/db"
)

type ActaInfo struct {
	Mesa             int    `json:"mesa"`
	VotosValidos     int    `json:"votos_validos"`
	VotosNulos       int    `json:"votos_nulos"`
	EdmundoGonzalez  int    `json:"edmundo_gonzalez"`
	NicolasMaduro    int    `json:"nicolas_maduro"`
	LuisMartinez     int    `json:"luis_martinez"`
	JavierBertucci   int    `json:"javier_bertucci"`
	JoseBrito        int    `json:"jose_brito"`
	AntonioEcarri    int    `json:"antonio_ecarri"`
	ClaudioFermin    int    `json:"claudio_fermin"`
	DanielCeballos   int    `json:"daniel_ceballos"`
	EnriqueMarquez   int    `json:"enrique_marquez"`
	BenjamminRausseo int    `json:"benjammin_rausseo"`
	URL              string `json:"url"`
}

type CenterInfo struct {
	Centro                int        `json:"centro"`
	Estado                string     `json:"estado"`
	Municipio             string     `json:"municipio"`
	Parroquia             string     `json:"parroquia"`
	TotalValidos          int        `json:"total_validos"`
	TotalNulos            int        `json:"total_nulos"`
	TotalEdmundoGonzalez  int        `json:"total_edmundo_gonzalez"`
	TotalNicolasMaduro    int        `json:"total_nicolas_maduro"`
	TotalLuisMartinez     int        `json:"total_luis_martinez"`
	TotalJavierBertucci   int        `json:"total_javier_bertucci"`
	TotalJoseBrito        int        `json:"total_jose_brito"`
	TotalAntonioEcarri    int        `json:"total_antonio_ecarri"`
	TotalClaudioFermin    int        `json:"total_claudio_fermin"`
	TotalDanielCeballos   int        `json:"total_daniel_ceballos"`
	TotalEnriqueMarquez   int        `json:"total_enrique_marquez"`
	TotalBenjamminRausseo int        `json:"total_benjammin_rausseo"`
	Actas                 []ActaInfo `json:"actas"`
}

func GetCenterInfo(dbc *db.Client, centro int) (*CenterInfo, error) {
	rows, err := db.GetCenterData(dbc, centro)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var centerInfo CenterInfo
	var actas []ActaInfo

	for rows.Next() {
		var acta ActaInfo
		var codigoEstado, codigoMunicipio, codParroquia int
		err := rows.Scan(
			&codigoEstado, &centerInfo.Estado, &codigoMunicipio, &centerInfo.Municipio,
			&codParroquia, &centerInfo.Parroquia, &acta.Mesa,
			&acta.VotosValidos, &acta.VotosNulos, &acta.EdmundoGonzalez,
			&acta.NicolasMaduro, &acta.LuisMartinez, &acta.JavierBertucci,
			&acta.JoseBrito, &acta.AntonioEcarri, &acta.ClaudioFermin,
			&acta.DanielCeballos, &acta.EnriqueMarquez, &acta.BenjamminRausseo, &acta.URL)
		if err != nil {
			return nil, err
		}

		// Aggregate totals in CenterInfo
		centerInfo.TotalValidos += acta.VotosValidos
		centerInfo.TotalNulos += acta.VotosNulos
		centerInfo.TotalEdmundoGonzalez += acta.EdmundoGonzalez
		centerInfo.TotalNicolasMaduro += acta.NicolasMaduro
		centerInfo.TotalLuisMartinez += acta.LuisMartinez
		centerInfo.TotalJavierBertucci += acta.JavierBertucci
		centerInfo.TotalJoseBrito += acta.JoseBrito
		centerInfo.TotalAntonioEcarri += acta.AntonioEcarri
		centerInfo.TotalClaudioFermin += acta.ClaudioFermin
		centerInfo.TotalDanielCeballos += acta.DanielCeballos
		centerInfo.TotalEnriqueMarquez += acta.EnriqueMarquez
		centerInfo.TotalBenjamminRausseo += acta.BenjamminRausseo

		actas = append(actas, acta)
	}

	centerInfo.Centro = centro
	centerInfo.Actas = actas

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &centerInfo, nil
}
