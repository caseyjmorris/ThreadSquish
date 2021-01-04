package scripts

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"sync/atomic"
)

type scriptResult struct {
	Path    string
	Success bool
	Skipped bool
	Err     error
}

type Runner struct {
	locker        uint32
	Script        string   `json:"script"`
	Errors        []string `json:"errors"`
	StopRequested bool     `json:"stopRequested"`
	Enqueued      []string `json:"enqueued"`
	Successful    []string `json:"successful"`
	Failed        []string `json:"failed"`
	Skipped       []string `json:"skipped"`
	Started       bool     `json:"started"`
	Done          bool     `json:"done"`
}

func (r *Runner) Stop() {
	r.StopRequested = true
}

func (r *Runner) IdentifyTargets(directory string, extension string, sample string) ([]string, error) {
	var result []string
	foundSample := false
	var innerErr error
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			innerErr = err
		}
		if filepath.Ext(path) != extension {
			return nil
		}
		if info.Name() == sample {
			foundSample = true
		}

		result = append(result, path)
		return nil
	})
	if err != nil {
		return result, fmt.Errorf("error walking directory:  %s", err)
	}
	if innerErr != nil {
		return result, fmt.Errorf("error walking directory:  %s", innerErr)
	}
	if !foundSample {
		return result, fmt.Errorf("did not find sample %q in %q", sample, directory)
	}

	return result, nil
}

func (r *Runner) readExcluded(path string) (map[string]bool, error) {
	whitespace := regexp.MustCompile("^\\s*$")
	text, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return make(map[string]bool), nil
	}
	if err != nil {
		return nil, fmt.Errorf("error opening %q:  %s", path, err)
	}
	result := make(map[string]bool)

	for _, file := range bytes.Split(text, []byte("\r\n")) {
		if whitespace.Match(file) {
			continue
		}
		result[string(file)] = true
	}

	return result, nil
}

func (r *Runner) RunScript(degreeOfParallelism int, script string, targets []string, argv []string,
	bookkeepingFile string) error {
	excluded, err := r.readExcluded(bookkeepingFile)
	if err != nil {
		errTxt := fmt.Sprintf("error opening book-keeping file %q:  %s", bookkeepingFile, err)
		r.Errors = append(r.Errors, errTxt)
		return errors.New(errTxt)
	}
	sink, err := os.OpenFile(bookkeepingFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errTxt := fmt.Sprintf("error opening book-keeping file %q:  %s", bookkeepingFile, err)
		r.Errors = append(r.Errors, errTxt)
		return errors.New(errTxt)
	}
	defer sink.Close()

	return r.runScriptWithCommander(degreeOfParallelism, script, targets, argv, excluded, sink, &StandardCommander{})
}

func (r *Runner) runScriptWithCommander(degreeOfParallelism int, script string, targets []string, argv []string,
	excluded map[string]bool, sink io.Writer, commander Commander) error {
	err := r.tryLock()
	if err != nil {
		return err
	}

	r.Started = true

	defer r.unlock()

	var wg sync.WaitGroup

	wg.Add(len(targets))

	targetQ := make(chan string, len(targets))
	doneQ := make(chan scriptResult, len(targets))
	r.Enqueued = append(r.Enqueued, targets...)

	for _, target := range targets {
		targetQ <- target
	}

	for i := 0; i < degreeOfParallelism; i++ {
		go r.runScriptForChannel(targetQ, doneQ, script, argv, excluded, commander)
	}

	var innerErr error

	go r.processScriptResults(doneQ, &innerErr, sink, &wg)

	wg.Wait()
	r.Done = true

	return innerErr
}

func (r *Runner) tryLock() error {
	if !atomic.CompareAndSwapUint32(&r.locker, 0, 1) {
		return errors.New("another script is already running")
	}
	return nil
}

func (r *Runner) unlock() {
	atomic.StoreUint32(&r.locker, 0)
}

func (r *Runner) runScriptForChannel(targetQ <-chan string, doneQ chan<- scriptResult, script string, argv []string,
	excluded map[string]bool, commander Commander) {
	for target := range targetQ {
		if r.StopRequested || excluded[target] {
			doneQ <- scriptResult{
				Path:    target,
				Success: false,
				Skipped: true,
				Err:     nil,
			}
			continue
		}
		args := append([]string{script, target}, argv...)
		//cmd := exec.Command("cmd.exe", args...)
		cmd := commander.Command("cmd.exe", args...)
		output, err := cmd.Output()
		log.Print(output)
		//err := cmd.Run()
		doneQ <- scriptResult{
			Path:    target,
			Success: err == nil,
			Skipped: false,
			Err:     err,
		}
	}
}

func (r *Runner) processScriptResults(doneQ <-chan scriptResult, innerErr *error, sink io.Writer, wg *sync.WaitGroup) {
	for result := range doneQ {
		if result.Skipped {
			r.Skipped = append(r.Skipped, result.Path)
		} else if !result.Success {
			r.Failed = append(r.Failed, result.Path)
		} else {
			r.Successful = append(r.Successful, result.Path)
			_, err := sink.Write([]byte(result.Path + "\r\n"))
			innerErr = &err
		}
		wg.Done()
	}
}
