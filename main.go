package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
	graceful "gopkg.in/tylerb/graceful.v1"
)

var Options = struct {
	Host string `long:"host" short:"h" description:"hostname" default:"localhost" env:"HOST"`
	Port int    `long:"port" short:"p" description:"port" default:"8000" env:"PORT"`
}{}

func main() {
	flags.Parse(&Options)
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) != "put" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"status": "400", "error": "method not allowed"})
		} else {
			body := struct {
				Filename string `json:"filename"`
			}{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"status": "400", "error": fmt.Sprintf("could not parse request: %v", err)})
			} else {
				if _, err = os.Lstat(body.Filename); err != nil {
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(map[string]string{"status": "400", "error": fmt.Sprintf("could find file: %v", err)})
				} else {
					log.Printf(body.Filename)
					err := os.Chtimes(body.Filename, time.Now(), time.Now())
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]string{"status": "400", "error": fmt.Sprintf("could touch file: %v", err)})
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}
			}
		}
	})
	log.Printf("starting file-toucher on %s:%d", Options.Host, Options.Port)
	graceful.Run(fmt.Sprintf("%s:%d", Options.Host, Options.Port), 10*time.Second, m)
}
