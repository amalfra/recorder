package recorder

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

// recorder is the base struct for recorders to inherit. It has references to
// all streams and core functions
type recorder struct {
	output               string
	originalOutputStream *os.File
	readStream           *os.File
	writeStream          *os.File
	wg                   sync.WaitGroup // wg will prevent the main process
	// from exiting before stream copy completes
}

// GetOutput will return the captured output
func (r *recorder) GetOutput() string {
	return r.output
}

// ClearOutput will clear the captured output
func (r *recorder) ClearOutput() {
	r.output = ""
}

// Stop will terminate recoding
func (r *recorder) Stop() {
	r.cleanup()
	// we need to wait till the goroutine which copies stream is completed
	r.wg.Wait()
}

// cleanup will reset the streams
func (r *recorder) cleanup() {
	// close streams and reset to OS default output stream
	r.writeStream.Close()
	os.Stdout = r.originalOutputStream
}

// startRecording will initiate copying between streams
func (r *recorder) startRecording() {
	// we need to wait till the goroutine which copies stream is completed
	r.wg.Add(1)
	// do copy process in a separate non blocking goroutine so that
	// printing won't block indefinitely
	go r.copyStream()
}

// copyStream does the operation of copying outputs between streams
func (r *recorder) copyStream() {
	// we need to let know that once copying is complete its free to release wait
	defer r.wg.Done()
	var buf bytes.Buffer
	// copy the pipe stream to both the output buffer as well as OS default
	// output stream so that code output behaviour won't change
	w := io.MultiWriter(&buf, r.originalOutputStream)
	_, err := io.Copy(w, r.readStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to copy from pipe: %v\n", err)
		return
	}
	// save the stream output to our string buffer
	r.output = buf.String()
}
