package recorder

import (
	"fmt"
	"os"
)

// Stdout is the struct that user can instantiate to execute
// recording operations
type Stdout struct {
	recorder
}

// Start will start the record process
func (s *Stdout) Start() error {
	var err error
	// backup OS default stdout
	s.originalOutputStream = os.Stdout
	s.readStream, s.writeStream, err = os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create pipe: %v\n", err)
		return err
	}
	// replace OS default stdout with pipe's
	os.Stdout = s.writeStream
	// go recorder go
	s.startRecording()
	return err
}
