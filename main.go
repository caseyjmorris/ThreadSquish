package main

import (
	"encoding/json"
	"fmt"
	"github.com/caseyjmorris/threadsquish/scripts"
	"log"
	"net/http"
	"os"
)

func parseProfile(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filePath := query["filePath"][0]
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, fmt.Sprintf("%v does not exist", filePath), http.StatusNotFound)
		return
	}
	parsed, err := scripts.ParseINIFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("Error parsing file:  %v", err), http.StatusInternalServerError)
		return
	}
	jsonText, _ := json.Marshal(parsed)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonText)
}

func main() {
	//http.HandleFunc("/", sayhelloName) // setting router rule
	http.HandleFunc("/profile", parseProfile)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
