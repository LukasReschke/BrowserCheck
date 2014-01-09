// Implements a very simple example of BrowserCheck
// Can be run with `go run main.go` and accessed via `localhost:8080` with a browser
package main

import (
	"encoding/json"
	"fmt"
	"github.com/LukasReschke/BrowserCheck"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", startHandler)
	http.HandleFunc("/check", checkHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("example.html"))
	tpl.Execute(w, nil)
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	result := browsercheck.Check(r.UserAgent() + r.FormValue("plugins"))

	// Return an empty array instead of null if the scan results are empty
	if result == nil {
		fmt.Fprintln(w, "[]")
		return
	}

	json.NewEncoder(w).Encode(result)
}
