package main

import (
	"api/actas"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/pkg/errors"
)

var apiEndpoint = "https://gdp.theempire.tech/api/data"

var shutdown = false

const writeResults = true

const thousand = 1000
const shutdownThreshold = 750
const startThousand = 30_000
const endThousand = 31_000
const numWorkers = 100

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println("Received signal:", sig)
		shutdown = true

		sig = <-sigs
		log.Println("Received signal (again):", sig)
		os.Exit(0)
	}()

	tasks := make(chan int)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, &wg)
	}

	for thou := startThousand; thou < endThousand; thou++ {
		tasks <- thou
	}

	close(tasks) // Close the tasks channel to indicate no more tasks will be sent
	wg.Wait()

	log.Println("All tasks completed")
}

func worker(id int, thous <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for thou := range thous {
		if shutdown { // do not start new workers
			return
		}
		log.Printf("Worker %d processing thousand %d\n", id, thou)
		err := fetchThousand(thou)
		if err != nil {
			panic(err)
		}
		log.Printf("Worker %d DONE with thousand %d\n", id, thou)
	}
}

func fetchThousand(prefix int) error {
	results := map[string]interface{}{}
	for i := 0; i < thousand; i++ {
		if shutdown && i < shutdownThreshold {
			log.Printf("received shutdown signal and we're only %d cedulas into this thousand, stopping", i)
			return nil // skip writing current progress for now
		}

		cedula := strconv.Itoa(prefix*1000 + i)
		resp, err := Resolve(cedula)
		if err != nil {
			log.Print(err)
			continue
		}
		log.Printf("successfully fetched cedula %s", cedula)

		results[cedula] = resp
	}

	// flush to file
	filename := fmt.Sprintf("cedulas/%d000.json", prefix)
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	if writeResults {
		return os.WriteFile(filename, data, 0664)
	} else {
		return nil
	}
}

func Resolve(cedula string) (*map[string]interface{}, error) {
	return actas.Retry(func() (*map[string]interface{}, error) {
		return fetch(cedula)
	}, 5, 1000)
}

func fetch(cedula string) (*map[string]interface{}, error) {
	r, err := http.Get(fmt.Sprintf("%s?cdi=V%s", apiEndpoint, cedula))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if string(body) == "{}" {
		return nil, errors.Errorf("empty response for cedula %s", cedula)
	}

	resp := map[string]interface{}{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if _, ok := resp["error"]; ok {
		return nil, errors.Errorf("error response for cedula %s", cedula)
	}

	return &resp, nil
}
