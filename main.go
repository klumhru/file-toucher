package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	graceful "gopkg.in/tylerb/graceful.v1"
)

var Options = struct {
	Host string `long:"host" short:"h" description:"hostname" default:"localhost" env:"HOST"`
	Port int    `long:"port" short:"p" description:"port" default:"8000" env:"PORT"`
}{}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) != "put" {
			json.NewEncoder(w).Encode(map[string]string{"status": "400", "error": "method not allowed"})
			w.WriteHeader(http.StatusBadRequest)
		} else {
			body := struct {
				Filename string
			}{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				json.NewEncoder(w).Encode(map[string]string{"status": "400", "error": fmt.Sprintf("could not parse request: %v", err)})
				w.WriteHeader(http.StatusBadRequest)
			} else {
				var fi os.FileInfo
				if fi, err = os.Lstat(body.Filename); err != nil {
					json.NewEncoder(w).Encode(map[string]string{"status": "400", "error": fmt.Sprintf("could find file: %v", err)})
					w.WriteHeader(http.StatusNotFound)
				} else {
					os.Chtimes(fi.Name(), time.Now(), time.Now())
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	})
	graceful.Run(fmt.Sprintf("%s:%d", Options.Host, Options.Port), 10*time.Second, m)
}
