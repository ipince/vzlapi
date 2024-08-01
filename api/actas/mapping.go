package actas

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
)

//go:embed estructura_completa.json
var structure []byte
var urlmap = map[string]string{}

var stateMapping = map[string]int{
	"DTTO. CAPITAL":   1,
	"EDO. AMAZONAS":   22,
	"EDO. ANZOATEGUI": 2,
	"EDO. APURE":      3,
	"EDO. ARAGUA":     4,
	"EDO. BARINAS":    5,
	"EDO. BOLIVAR":    6,
	"EDO. CARABOBO":   7,
	"EDO. COJEDES":    8,
	"EDO. DELTA AMAC": 23,
	"EDO. FALCON":     9,
	"EDO. GUARICO":    10,
	"EDO. LA GUAIRA":  24,
	"EDO. LARA":       11,
	"EDO. MERIDA":     12,
	"EDO. MIRANDA":    13,
	"EDO. MONAGAS":    14,
	"EDO. PORTUGUESA": 16,
	"EDO. SUCRE":      17,
	"EDO. TACHIRA":    18,
	"EDO. TRUJILLO":   19,
	"EDO. YARACUY":    20,
	"EDO. ZULIA":      21,
	"EDO.NVA.ESPARTA": 15,
}

func init() {
	sitemap := map[string]map[string]map[string]interface{}{}
	err := json.Unmarshal(structure, &sitemap)
	if err != nil {
		panic(err)
	}
	for _, muns := range sitemap {
		for mun, data := range muns {
			urlmap[mun] = data["url"].(string)
			parishes := data["parroquias"].(map[string]interface{})
			for parish, data2 := range parishes {
				urlmap[parish] = data2.(map[string]interface{})["url"].(string)
				centers := data2.(map[string]interface{})["centros"].(map[string]interface{})
				for center, url := range centers {
					urlmap[center] = url.(string)
				}
			}
		}
	}
}

func FillUrls(info *CedulaInfo) {
	if url, ok := stateMapping[info.StateName]; ok {
		info.ResultsStateURL = fmt.Sprintf("https://resultadosconvzla.com/estado/%d", url)
	}
	if url, ok := urlmap[info.CountyName]; ok {
		info.ResultsCountyURL = url
	}
	if url, ok := urlmap[info.ParishName]; ok {
		info.ResultsParishURL = url
	}
	if url, ok := urlmap[info.CenterName]; ok {
		info.ResultsCenterURL = url
	} else {
		slog.Warn(fmt.Sprintf("failed to find results url for center %s", info.CenterName))
	}
}
