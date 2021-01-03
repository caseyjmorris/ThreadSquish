package main

import (
	"encoding/json"
	"fmt"
	"github.com/caseyjmorris/threadsquish/scripts"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var runner = scripts.Runner{}

func parseProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	query := r.URL.Query()
	filePath := query["filePath"][0]
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("%q does not exist", filePath), http.StatusNotFound)
		return
	}
	parsed, err := scripts.ParseINIFile(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing file:  %q", err), http.StatusInternalServerError)
		return
	}
	jsonText, _ := json.Marshal(parsed)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonText)
}

func status(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonText, _ := json.Marshal(runner)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonText)
}

func stop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	runner.Stop()
	w.WriteHeader(http.StatusOK)
}

func start(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var cmdReq scripts.CommandRequest
	bodyText, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading response:  %s", err), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bodyText, &cmdReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing request:  %s", err), http.StatusBadRequest)
		return
	}

	if cmdReq.DegreeOfParallelism < 1 {
		http.Error(w, "degree of parallelism must be 1 or greater", http.StatusBadRequest)
		return
	}

	_, err = os.Stat(cmdReq.Script)
	if os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("%q does not exist", cmdReq.Script), http.StatusNotFound)
		return
	}
	parsed, err := scripts.ParseINIFile(cmdReq.Script)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing file:  %q", err), http.StatusInternalServerError)
		return
	}

	files, err := runner.IdentifyTargets(cmdReq.Directory, filepath.Ext(parsed.Example), parsed.Example)
	if err != nil {
		http.Error(w, fmt.Sprintf("error finding targets:  %s", err), http.StatusBadRequest)
		return
	}

	go runner.RunScript(cmdReq.DegreeOfParallelism, cmdReq.Script, files, cmdReq.Arguments,
		filepath.Join(cmdReq.Directory, "_THREADSQUISH.TXT"))

	w.WriteHeader(http.StatusOK)
}

func terminate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if runner.Started && !runner.Done {
		http.Error(w, "cannot terminate the application while tasks are in progress", http.StatusForbidden)
		return
	}

	os.Exit(0)
}

func serveStatic(w http.ResponseWriter, location string, contentType string) {
	body, err := ioutil.ReadFile(location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("Error opening %q:  %s", location, err), http.StatusInternalServerError)
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
	http.HandleFunc("/status", status)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/terminate", terminate)
	http.HandleFunc("/start", start)
	go func() {
		time.Sleep(time.Second)
		cmd := exec.Command("powershell.exe", "-command", "start http://localhost:9090")
		cmd.Start()
	}()
	err := http.ListenAndServe(":9090", nil)
	log.Fatal("ListenAndServe: ", err)
}
