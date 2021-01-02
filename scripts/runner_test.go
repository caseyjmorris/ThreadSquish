package scripts

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"sync"
	"testing"
)

type TestCommander struct {
	Targets sync.Map
}

func (t *TestCommander) Command(name string, arg ...string) *exec.Cmd {
	t.Targets.Store(arg[0], arg)
	cmd := exec.Cmd{}
	return &cmd
}

var _ Commander = (*TestCommander)(nil)

func TestRunner_tryLock(t *testing.T) {
	runner := Runner{}
	err := runner.tryLock()
	if err != nil {
		t.Errorf("failed to get first lock:  %s", err)
		return
	}
	err = runner.tryLock()
	if err == nil {
		t.Error("allowed multiple lockings")
		return
	}
	runner.unlock()
	err = runner.tryLock()
	if err != nil {
		t.Errorf("failed to get second lock:  %s", err)
	}
}

func TestRunner_runScriptWithCommander(t *testing.T) {
	commander := TestCommander{}
	runner := Runner{}
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	var targets []string
	var excluded []string
	for i := 0; i < 50; i++ {
		targets = append(targets, strconv.Itoa(i))
		if i > 25 {
			excluded = append(excluded, strconv.Itoa(i))
		}
	}
	err := runner.runScriptWithCommander(16, "c:\\users\\my user\\script.cmd", targets, []string{"a", "b", "c"}, writer, &commander)

	if err != nil {
		t.Errorf("failed running script:  %s", err)
		return
	}

	for i := 0; i < 50; i++ {
		str := strconv.Itoa(i)
		_, ok := commander.Targets.Load(str)
		if i > 25 == ok {
			t.Errorf("unexpected inclusion value %t for %q", ok, str)
			return
		}
	}
}
