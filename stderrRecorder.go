package recorder

import (
	"fmt"
	"os"
)

// Stderr is the struct that user can instantiate to execute
// recording operations
type Stderr struct {
	recorder
}

// Start will start the record process
func (s *Stderr) Start() error {
	var err error
	// backup OS default stderr
	s.originalOutputStream = os.Stderr
	s.readStream, s.writeStream, err = os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create pipe: %v\n", err)
		return err
	}
	// replace OS default stderr with pipe's
	os.Stderr = s.writeStream
	// go recorder go
	s.startRecording()
	return err
}
