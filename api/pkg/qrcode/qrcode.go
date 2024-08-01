package qrcode

import (
	"encoding/csv"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	_ "image/jpeg"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pkg/errors"
)

func ExtractData(dir string, resultsFile string) {
	filenames, err := processed(resultsFile)
	if err != nil {
		log.Printf("failed to read results from %s", resultsFile)
		panic(err)
	}

	results, err := processDir(filenames, dir, resultsFile)
	if err != nil {
		panic(err)
	}

	log.Printf("processed %d new actas", len(results))
	log.Printf("DONE :)")
}

func processed(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	filenames := []string{}
	for i, row := range records {
		if i == 0 { // skip header
			continue
		}
		filenames = append(filenames, row[0])
	}
	log.Printf("read %d existing processed actas", len(filenames))

	return filenames, nil
}

func appendResults(results []*Result, path string) error {
	// We append so we avoid clobbering existing data, and punt on parsing existing csv data
	csvFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	for _, r := range results {
		err = writer.Write(r.asRow())
		if err != nil {
			return err
		}
	}
	writer.Flush()
	log.Printf("flushed %d rows", len(results))

	return nil
}

func processDir(processed []string, dir string, resultsFile string) ([]*Result, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	results := []*Result{}
	unflushed := []*Result{}
	for _, file := range files {
		if len(unflushed) > 0 && len(unflushed)%100 == 0 { // flush every 100 results
			err = appendResults(unflushed, resultsFile)
			if err != nil {
				return nil, err
			}
			unflushed = []*Result{}
		}

		if !file.IsDir() {
			filename := file.Name()
			if slices.Contains(processed, filename) {
				log.Printf("skipping %s, already extracted", filename)
				continue
			}
			log.Printf("processing %s...", filename)
			result, err := process(filename, filepath.Join(dir, filename))
			if err != nil {
				log.Printf("failed to process %s: %s", filename, err)
				continue
			}
			byCandidate := result.CandidateTotals()
			log.Printf("successfully processed %s... maduro %d, edmundo %d", filename, byCandidate[CandidateMaduro], byCandidate[CandidateGonzalez])

			results = append(results, result)
			unflushed = append(unflushed, result)
		}
	}
	err = appendResults(unflushed, resultsFile)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func process(filename, path string) (*Result, error) {
	data, err := readQR(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read qr code from image")
	}

	return parse(filename, data)
}

func readQR(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func parse(filename, data string) (*Result, error) {
	// Example data:
	// 110601011.04.1.0001!122,1,0,0,4,2,0,0,2,1,0,1,2,1,0,0,0,5,0,2,0,0,0,0,0,0,0,0,1,0,0,0,0,8,22,406,0,1!0!0
	parts := strings.Split(data, "!")
	if len(parts) != 4 {
		return nil, errors.New(fmt.Sprintf("did not find 4 parts in data: %s", data))
	}

	actaCode := parts[0] // 110601011.04.1.0001 (first part is the voting center code)
	actaCodeParts := strings.Split(actaCode, ".")
	centerCode := actaCodeParts[0]
	table := actaCodeParts[1]

	validVotes := parts[1]
	nullVotes, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse null votes")
	}
	invalidVotes, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse invalid votes")
	}

	votes := strings.Split(validVotes, ",")
	if len(votes) != 38 {
		return nil, errors.New(fmt.Sprintf("found unexpected number of votes in data: %s", data))
	}

	result := &Result{
		Code:         actaCode,
		CenterCode:   centerCode,
		Table:        table,
		ActaFilename: filename,
		NullVotes:    nullVotes,
		InvalidVotes: invalidVotes,
		Votes:        map[Option]int{},
	}
	sum := 0
	for i, v := range votes {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse vote")
		}
		opt := ballotOrder[i]
		result.Votes[opt] = vInt
		sum += vInt
	}
	result.ValidVotes = sum

	return result, nil
}
