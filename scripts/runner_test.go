package scripts

import (
	"os/exec"
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
