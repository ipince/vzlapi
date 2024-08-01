# vzlapi

Serving at https://api.vzlapi.com, an open API for venezuelan-related data.

## Example

The `/actas` returns information about actas (voting machine receipts) for the 7/28/2024 Presidential elections.

```
// https://api.vzlapi.com/actas?cedula=4000002

{
  // The ID of the voter
  "Cedula": "V4000002",
  
  // Information on where the voter is registered to vote
  "StateID": "21",
  "StateName": "EDO. ZULIA",
  "CountyID": "307",
  "CountyName": "MP. MARACAIBO",
  "ParishID": "996",
  "ParishName": "PQ. COQUIVACOA",
  "CenterOldID": "10532",
  "CenterID": "210502010",
  "CenterName": "ESCUELA BASICA NACIONAL MONSEÑOR FRANCISCO ANTONIO GRANADILLO",
  "CenterAddress": "SECTOR 18 DE OCTUBRE FRENTE AVENIDA 2. DERECHA CALLE E. IZQUIERDA CALLE F AVENIDA 2 ENTRE CALLES E Y F SECTOR 18 DE OCTUBRE CASA",
  "TableID": "10439",
  "TableNumber": "1",
  
  // If a scanned acta exists, returns its locations
  "ActaFilename": "676839_898940_0234Acta0234.jpg",
  "ActaBucketURL": "https://elecciones2024ve.s3.amazonaws.com/676839_898940_0234Acta0234.jpg",
  "ActaStaticURL": "https://static.resultadosconvzla.com/676839_898940_0234Acta0234.jpg",
  
  // Deeplinks to the opposition-run results website
  "ResultsStateURL": "https://resultadosconvzla.com/estado/21",
  "ResultsCountyURL": "https://resultadosconvzla.com/municipio/307",
  "ResultsParishURL": "https://resultadosconvzla.com/parroquia/996",
  "ResultsCenterURL": "https://resultadosconvzla.com/centro/10532",
  "ResultsTableURL": "https://resultadosconvzla.com/mesa/10532/10439"
}
```

## Run locally

```
devbox run start
```

or
```
devbox shell
cd api
go run ./main.go
```

## Deploy

All merges to `main` are automatically deployed to production using https://cloud.jetify.com
