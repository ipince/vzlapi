# api.vzlapi.com

# Objetivos

1. Tener la **[data](https://docs.google.com/spreadsheets/d/1l6ThiQQZXog_8fBw3z5RwqThG7QAy0AqF4wPYpvGUWA/edit?gid=1712398662#gid=1712398662
   )** mas completa y confiable sobre las elecciones presidenciales 2024 en Venezuela.
2. Hacer la data publica a traves de un **[API](https://api.vzlapi.com/actas?cedula?4000000)**, para que quien quiera la pueda accesar..

A futuro, me gustaria expandir el "scope" y recolectar y difundir mas data relacionada Venezuela, no tan solo enfocada en las elecciones presidenciales.

Preguntas/problemas, o quieres ayudar? Agrega un [issue](https://github.com/ipince/vzlapi/issues).

# 1. Data

## Contenido

Por ahora, contamos con la siguiente data, mas facilmente accesible en Google Sheets:
https://docs.google.com/spreadsheets/d/1l6ThiQQZXog_8fBw3z5RwqThG7QAy0AqF4wPYpvGUWA/edit?gid=1712398662#gid=1712398662

Tambien contamos con un archivo de ~12GB con todas las imagenes de ~24mil actas. (Tengo que ponerlo en un lugar publico, pero es el mismo que otras personas han publicado).

### Contiene la siguiente data
- Centros CNE: Los 15.962 centros electorales en el registro para el 2024, uncluido el numero de electores por centro. Data proveniente del CNE.
- Centros convzla: Todos (o casi todos) los centros electorales para los cuales existen actas digitalizadas en https://resultadosconvzla.com. Incluye un mapeo entre los IDs de resultadosconvzla y los IDs del CNE para cada centro. Data proviene de scraping.
- Actas (data): La data extraida del codigo QR de casi todas (~22mil de ~24mil) de las actas digitalizadas. Data proviene del procesamiento de imagenes.

### Algunas cosas por destacar
- La data de las actas (los votos) NO es tomada de resultadosconvzla.com. TODA es extraida por nosotros de las imagenes de las actas utilizando lector de codigo QR.
- Por ende, la data es mas completa que la data de resultadosconvzla.com. Aqui tambien tenemos los votos por cada uno de los 10 candidatos, y los votos nulos (no todos 0) e invalidos (todos 0).

## Herramientas

Hemos procesado actas de distintas fuentas (ademas de la nuestra propia), e intentado varias herramientas para extraer la data, cada con distinta eficacia.
- Las imagenes de las actas provienen de varios archivos que distintas personas han publicado. TODO: poner un listado
- La mayoria de las actas fueron procesadas utilizando: https://github.com/xaiki/resultadosvzla-tools. Es la que tuvo el mayor porcentaje de exito.


# 2. API -- https://api.vzlapi.com/

Por ahora solo hay un endpoint, `/actas`, que acepta una cedula y regresa informacion del votante, su centro electoral, y la informacion que tenemos de ese centro electoral.

Ejemplo:
```
// https://api.vzlapi.com/actas?cedula=4000000

{
  "Cedula": "V4000000",
  
  // El centro de votacion asignado al elector
  "StateID": "18",
  "StateName": "EDO. TACHIRA",
  "CountyID": "247",
  "CountyName": "MP. SAN CRISTOBAL",
  "ParishID": "815",
  "ParishName": "PQ. LA CONCORDIA",
  "CenterOldID": "7433",
  "CenterID": "180801028",
  "CenterName": "ESCUELA BOLIVARIANA RITA ELISA MEDINA DE USECHE",
  "CenterAddress": "BARRIO LA VICTORIA DERECHA CALLE DETRAS DEL HOGAR DON BOSCO. IZQUIERDA CALLE 1 CENTRO MEDICO ROTARY. FRENTE CALLE PRINCIPAL ENTRADA AL ALBERGUE JUVENIL CIUDAD DE LOS MUCHACHOS CASA",
  
  // Incluida el numero de mesa (no siempre disponible)
  "TableID": "7362",
  "TableNumber": "1",
  
  // Elace al acta digitalizada (si es que existe)
  "ActaFilename": "788144_568516_0998Acta0563.jpg",
  "ActaBucketURL": "https://elecciones2024ve.s3.amazonaws.com/788144_568516_0998Acta0563.jpg",
  "ActaStaticURL": "https://static.resultadosconvzla.com/788144_568516_0998Acta0563.jpg",
  
  // Data extraida del acta (si existe), con todos los votos por candidato
  "ActaValidVotes": 311,
  "ActaNullVotes": 0,
  "ActaInvalidVotes": 0,
  "ActaCandidateVotes": {
    "Antonio Ecarri": 1,
    "Benjamin Rausseo": 1,
    "Claudio Fermin": 2,
    "Daniel Ceballos": 0,
    "Edmundo Gonzalez": 238,
    "Enrique Marquez": 0,
    "Javier Bertucci": 2,
    "Jose Brito": 0,
    "Luis Martinez": 1,
    "Nicolas Maduro": 66
  },
  
  // Enlaces a las paginas correspondientes en resultadosconvzla.com
  "ResultsStateURL": "https://resultadosconvzla.com/estado/18",
  "ResultsCountyURL": "https://resultadosconvzla.com/municipio/247",
  "ResultsParishURL": "https://resultadosconvzla.com/parroquia/815",
  "ResultsCenterURL": "https://resultadosconvzla.com/centro/7433",
  "ResultsTableURL": "https://resultadosconvzla.com/mesa/7433/7362"
}
```

## Corer localmente

Lo mas facil es utilizar devbox, pero tambien puedes instalar todas las dependencias directamente.

Usando devbox
```
# Instalar devbox si no lo tienes: https://www.jetify.com/devbox/docs/installing_devbox/

devbox run start
```

o
```
devbox shell
cd api
go run ./main.go
```

## Como lanzar/deploy a produccion

Cuando se hace merge a `main`, el API hace deploy automaticamente, hosteado por https://cloud.jetify.com
