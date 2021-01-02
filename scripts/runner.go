package scripts

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	locker        *uint32
	Script        string
	Errors        []string
	stopRequested bool
	Enqueued      []string
	Successful    []string
	Failed        []string
	Skipped       []string
	Canceled      []string
	Done          bool
}

func (r *Runner) Stop() {
	r.stopRequested = true
}

func (r *Runner) IdentifyTargets(directory string, extension string, sample string,
	excluded map[string]bool) ([]string, error) {
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
		if excluded[path] {
			r.Skipped = append(r.Skipped, path)
			return nil
		}

		r.Enqueued = append(r.Enqueued, path)
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

func (r *Runner) ReadExcluded(path string) ([]string, error) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening %q:  %s", path, err)
	}
	var result []string

	for _, file := range bytes.Split(text, []byte("\r\n")) {
		result = append(result, string(file))
	}

	return result, nil
}

func (r *Runner) RunScript(degreeOfParallelism int, script string, targets []string, argv []string,
	bookkeepingFile string, commander Commander) error {
	err := r.tryLock()
	if err != nil {
		return err
	}

	defer r.unlock()

	sink, err := os.OpenFile(bookkeepingFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening book-keeping file %q:  %s", bookkeepingFile, err)
	}

	var wg sync.WaitGroup
	wg.Add(len(targets))

	targetQ := make(chan string, len(targets))
	doneQ := make(chan scriptResult, len(targets))

	for _, target := range targets {
		targetQ <- target
	}

	for i := 0; i < degreeOfParallelism; i++ {
		go r.runScriptForChannel(targetQ, doneQ, script, argv, commander)
	}

	var innerErr error

	go r.processScriptResults(doneQ, &innerErr, sink, &wg)

	wg.Wait()
	r.Done = true

	return innerErr
}

func (r *Runner) tryLock() error {
	if !atomic.CompareAndSwapUint32(r.locker, 0, 1) {
		return errors.New("another script is already running")
	}
	return nil
}

func (r *Runner) unlock() {
	atomic.StoreUint32(r.locker, 0)
}

func (r *Runner) runScriptForChannel(targetQ <-chan string, doneQ chan<- scriptResult, script string, argv []string,
	commander Commander) {
	for target := range targetQ {
		if r.stopRequested {
			doneQ <- scriptResult{
				Path:    target,
				Success: false,
				Skipped: true,
				Err:     nil,
			}
			return
		}
		args := append([]string{script, target}, argv...)
		//cmd := exec.Command("cmd.exe", args...)
		cmd := commander.Command("cmd.exe", args...)
		err := cmd.Run()
		if err != nil {
			doneQ <- scriptResult{
				Path:    target,
				Success: false,
				Skipped: false,
				Err:     err,
			}
			return
		}
		doneQ <- scriptResult{
			Path:    target,
			Success: true,
			Skipped: false,
			Err:     nil,
		}
	}
}

func (r *Runner) processScriptResults(doneQ <-chan scriptResult, innerErr *error, sink *os.File, wg *sync.WaitGroup) {
	for result := range doneQ {
		if result.Skipped {
			r.Skipped = append(r.Skipped, result.Path)
		} else if !result.Success {
			r.Failed = append(r.Failed, result.Path)
		} else {
			r.Successful = append(r.Successful, result.Path)
			_, err := sink.WriteString(result.Path + "\r\n")
			innerErr = &err
		}
		wg.Done()
	}
}
