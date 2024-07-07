package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

const refreshHashInSeconds = 60

func downloadAndHashImage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching image: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response bytes: %v", err)
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func downloadAndHashImages(config gjson.Result) {
	for {
		state.mu.Lock()

		state.EnImageHash, _ = downloadAndHashImage(config.Get("secrets.en_image").String())
		state.EnPImageHash, _ = downloadAndHashImage(config.Get("secrets.en_image_p").String())
		state.EsImageHash, _ = downloadAndHashImage(config.Get("secrets.es_image").String())
		state.EsPImageHash, _ = downloadAndHashImage(config.Get("secrets.es_image_p").String())
		state.FrImageHash, _ = downloadAndHashImage(config.Get("secrets.fr_image").String())
		state.PoImageHash, _ = downloadAndHashImage(config.Get("secrets.po_image").String())
		state.ItImageHash, _ = downloadAndHashImage(config.Get("secrets.it_image").String())
		state.DeImageHash, _ = downloadAndHashImage(config.Get("secrets.de_image").String())

		state.mu.Unlock()

		time.Sleep(time.Duration(refreshHashInSeconds) * time.Second)
	}
}

func getHash(w http.ResponseWriter, _ *http.Request, hash *string) {
	state.mu.Lock()
	defer state.mu.Unlock()

	fmt.Fprint(w, *hash)
}
