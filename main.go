package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"log"
	"net/http"
)
import _ "go.uber.org/automaxprocs"

type BackendConfig struct {
	Backends []*WeightedURL `toml:"backends"`
}

type WeightedURL struct {
	URL    string
	Weight int
}

var (
	backends    []*WeightedURL
	nextBackend = 0
)

func main() {
	var config BackendConfig
	_, err := toml.DecodeFile("./backends.toml", &config)
	if err != nil {
		fmt.Println(err)
	}

	var targets []*WeightedURL
	for _, backend := range config.Backends {
		targets = append(targets, &WeightedURL{
			URL:    backend.URL,
			Weight: backend.Weight,
		})
	}
	backends = targets

	// start server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request:", w)
		u := targetURL()
		resp, err := http.Get(u + r.RequestURI)
		if err != nil {
			_, _ = fmt.Fprintln(w, "Error:", err)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			_, _ = fmt.Fprintln(w, "Error:", err)
		} else {
			_, _ = fmt.Fprintf(w, "%s\n", body)
		}
		return
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	log.Printf("Starting server")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func targetURL() string {
	b := backends[nextBackend]
	nextBackend = (nextBackend + 1) % len(backends)
	log.Println("target: ", b.URL)
	return b.URL
}
