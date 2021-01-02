package scripts

import (
	"bufio"
	"bytes"
	"os/exec"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

type TestCommander struct {
	Targets sync.Map
}

func (t *TestCommander) Command(name string, arg ...string) *exec.Cmd {
	t.Targets.Store(arg[1], arg)
	cmd := exec.Command("cmd.exe", "/C", "echo", "hello")
	return cmd
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
	excluded := make(map[string]bool)
	for i := 0; i < 50; i++ {
		targets = append(targets, strconv.Itoa(i))
		if i > 25 {
			excluded[strconv.Itoa(i)] = true
		}
	}
	err := runner.runScriptWithCommander(16, "c:\\users\\my user\\script.cmd", targets,
		[]string{"a", "b", "c"}, excluded, writer, &commander)

	if err != nil {
		t.Errorf("failed running script:  %s", err)
		return
	}

	for i := 0; i < 50; i++ {
		str := strconv.Itoa(i)
		_, ok := commander.Targets.Load(str)
		if i > 25 == ok {
			t.Errorf("unexpected inclusion value %t for %q", ok, str)
		}
	}
}

func TestRunner_runScriptWithCommanderOutputWithoutParallelism(t *testing.T) {
	commander := TestCommander{}
	runner := Runner{}
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	var targets []string
	excluded := make(map[string]bool)
	for i := 0; i < 10; i++ {
		targets = append(targets, strconv.Itoa(i))
		if i > 5 {
			excluded[strconv.Itoa(i)] = true
		}
	}
	err := runner.runScriptWithCommander(1, "c:\\users\\my user\\script.cmd", targets,
		[]string{"a", "b", "c"}, excluded, writer, &commander)

	if err != nil {
		t.Errorf("failed running script:  %s", err)
		return
	}

	for i := 0; i < 10; i++ {
		str := strconv.Itoa(i)
		_, ok := commander.Targets.Load(str)
		if i > 5 == ok {
			t.Errorf("unexpected inclusion value %t for %q", ok, str)
		}
	}

	expectedSink := "0\r\n1\r\n2\r\n3\r\n4\r\n5\r\n"
	writer.Flush()
	resultSink := b.String()

	if resultSink != expectedSink {
		t.Errorf("Unexpected bookkeeping file result:  \r\n%v\r\n\r\nExpected:  \r\n%v", resultSink, expectedSink)
	}

	expectedSuccessful := []string{"0", "1", "2", "3", "4", "5"}

	if !reflect.DeepEqual(runner.Successful, expectedSuccessful) {
		t.Errorf("Expected successful records %v but found %v", expectedSuccessful, runner.Successful)
	}

	expectedSkipped := []string{"6", "7", "8", "9"}

	if !reflect.DeepEqual(runner.Skipped, expectedSkipped) {
		t.Errorf("Expected skipped records %v but found %v", expectedSkipped, runner.Skipped)
	}
}

func TestRunner_runScriptWithCommanderStopRequested(t *testing.T) {
	commander := TestCommander{}
	runner := Runner{}
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	var targets []string
	excluded := make(map[string]bool)
	for i := 0; i < 10; i++ {
		targets = append(targets, strconv.Itoa(i))
		if i > 5 {
			excluded[strconv.Itoa(i)] = true
		}
	}

	runner.Stop()
	if !runner.stopRequested {
		t.Error("failed to stop")
		return
	}

	err := runner.runScriptWithCommander(1, "c:\\users\\my user\\script.cmd", targets,
		[]string{"a", "b", "c"}, excluded, writer, &commander)

	if err != nil {
		t.Errorf("failed running script:  %s", err)
		return
	}

	for i := 0; i < 10; i++ {
		str := strconv.Itoa(i)
		_, ok := commander.Targets.Load(str)
		if ok {
			t.Errorf("unexpected inclusion value %t for %q", ok, str)
		}
	}

	expectedSink := ""
	writer.Flush()
	resultSink := b.String()

	if resultSink != expectedSink {
		t.Errorf("Unexpected bookkeeping file result:  \r\n%v\r\n\r\nExpected:  \r\n%v", resultSink, expectedSink)
	}

	expectedSkipped := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	if !reflect.DeepEqual(runner.Skipped, expectedSkipped) {
		t.Errorf("Expected skipped records %v but found %v", expectedSkipped, runner.Skipped)
	}
}
