package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var Options = struct {
	Host string `long:"host" short:"h" description:"hostname" default:"localhost" env:"HOST"`
	Port int    `long:"port" short:"p" description:"port" default:"8000" env:"PORT"`
}{}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) != "post" {
			return json.NewEncoder(w).Encode(map[string]string{"status": "400", "error": "method not allowed"})
		}
	})
}
