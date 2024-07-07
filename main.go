package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

func loadConfig(filename string) (gjson.Result, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("unable to read config file: %v", err)
	}

	config := gjson.ParseBytes(data)
	return config, nil
}

func main() {
	config, err := loadConfig("Config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	go downloadAndHashImages(config)

	http.HandleFunc("/en", RateLimited(func(w http.ResponseWriter, r *http.Request) {
		getHash(w, r, &state.EnImageHash)
	}))
	http.HandleFunc("/en_p", RateLimited(func(w http.ResponseWriter, r *http.Request) {
		getHash(w, r, &state.EnPImageHash)
	}))
	http.HandleFunc("/es", RateLimited(func(w http.ResponseWriter, r *http.Request) {
		getHash(w, r, &state.EsImageHash)
	}))
	http.HandleFunc("/es_p", RateLimited(func(w http.ResponseWriter, r *http.Request) {
		getHash(w, r, &state.EsPImageHash)
	}))
	http.HandleFunc("/fr", RateLimited(func(w http.ResponseWriter, r *http.Request) {
		getHash(w, r, &state.FrImageHash)
	}))
	http.HandleFunc("/po", RateLimited(func(w http.ResponseWriter, r *http.Request) {
		getHash(w, r, &state.PoImageHash)
	}))
	http.HandleFunc("/it", RateLimited(func(w http.ResponseWriter, r *http.Request) {
		getHash(w, r, &state.ItImageHash)
	}))
	http.HandleFunc("/de", RateLimited(func(w http.ResponseWriter, r *http.Request) {
		getHash(w, r, &state.DeImageHash)
	}))

	log.Println("Starting server on :9191")
	log.Fatal(http.ListenAndServe(":9191", nil))
}
