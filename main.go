package main

import (
	"encoding/json"
	"fmt"
	"github.com/caseyjmorris/threadsquish/scripts"
	"io/ioutil"
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

func serveStatic(w http.ResponseWriter, location string, contentType string) {
	body, err := ioutil.ReadFile(location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("Error opening %v:  %v", location, err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", contentType)
	_, _ = w.Write(body)
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		serveStatic(writer, "static/index.html", "text/html")
	})
	http.HandleFunc("/threadsquish.js", func(writer http.ResponseWriter, request *http.Request) {
		serveStatic(writer, "static/threadsquish.js", "text/javascript")
	})
	http.HandleFunc("/threadsquish.css", func(writer http.ResponseWriter, request *http.Request) {
		serveStatic(writer, "static/threadsquish.css", "text/css")
	})
	http.HandleFunc("/profile", parseProfile)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
