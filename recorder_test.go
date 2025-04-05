package recorder_test

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/amalfra/recorder/v4"
)

var wg sync.WaitGroup

func TestStdout(t *testing.T) {
	stdoutText := `This should be recorded.
	second line that should be recorded~~~
	end output`
	stderrText := "this goes to stderr"
	stdoutRecorder := new(recorder.Stdout)
	fmt.Println("before recording")
	err := stdoutRecorder.Start()
	if err != nil {
		t.Fatalf("Error initializing recorder: %s", err)
	}
	fmt.Println(stdoutText)
	fmt.Fprintf(os.Stderr, "%s", stderrText)
	stdoutRecorder.Stop()
	recordedOutput := strings.TrimSpace(stdoutRecorder.GetOutput())
	if stdoutText != recordedOutput {
		t.Fatalf("Wrong output recorded: %s", recordedOutput)
	}
}

func TestStderr(t *testing.T) {
	stderrText := `This should be recorded.
	second line that should be recorded~~~
	end output`
	stdoutText := "this goes to stdout"
	stderrRecorder := new(recorder.Stderr)
	fmt.Println("before recording")
	err := stderrRecorder.Start()
	if err != nil {
		t.Fatalf("Error initializing recorder: %s", err)
	}
	fmt.Println(stdoutText)
	fmt.Fprintf(os.Stderr, "%s", stderrText)
	stderrRecorder.Stop()
	recordedOutput := strings.TrimSpace(stderrRecorder.GetOutput())
	if stderrText != recordedOutput {
		t.Fatalf("Wrong output recorded: %s", recordedOutput)
	}
}

func TestStdoutGoRoutine(t *testing.T) {
	stdoutText := `This should be recorded.
	second line that should be recorded~~~
	end output`
	stderrText := "this goes to stderr"
	stdoutRecorder := new(recorder.Stdout)
	fmt.Println("before recording")
	err := stdoutRecorder.Start()
	if err != nil {
		t.Fatalf("Error initializing recorder: %s", err)
	}
	fmt.Println(stdoutText)
	wg.Add(1)
	go goRoutine()
	wg.Wait()
	fmt.Fprintf(os.Stderr, "%s", stderrText)
	stdoutRecorder.Stop()
	recordedOutput := strings.TrimSpace(stdoutRecorder.GetOutput())
	stdoutText += "\nin goRoutine"
	stdoutText += "\nout goRoutine"
	if stdoutText != recordedOutput {
		t.Fatalf("Wrong output recorded: %s", recordedOutput)
	}
}

func goRoutine() {
	defer wg.Done()
	fmt.Println("in goRoutine")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("out goRoutine")
}
