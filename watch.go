package asteroid

import (
	"bufio"
	"io"
	"os/exec"
)

// Watch watches for SIP Debug messages and fires a callback for all calls
// observed.
type Watch struct {
	// Config
	asteriskPath string
	onLine       func(line string)

	// Internals
	reader    io.Reader
	cmd       *exec.Cmd
	isStopped bool
}

// NewWatch for a set of headers to watch for, and a callback to run when a call
// is observed.
func NewWatch() Watch {
	ret := Watch{
		asteriskPath: "/usr/sbin/asterisk",
		isStopped:    false,
	}
	return ret
}

// SetAsteriskPath to the Asterisk executable.
func (watch *Watch) SetAsteriskPath(path string) {
	watch.asteriskPath = path
}

// SetReader instead of Asterisk executable.
func (watch *Watch) SetReader(reader io.Reader) {
	watch.reader = reader
}

// Start watching for headers.
func (watch *Watch) Start() error {
	var err error
	if watch.reader == nil {
		watch.cmd = exec.Command(watch.asteriskPath, "-r")
		watch.reader, err = watch.cmd.StdoutPipe()
		if err == nil {
			err = watch.cmd.Start()
		}
	}

	if err == nil {
		go func() {
			s := bufio.NewScanner(watch.reader)
			for !watch.isStopped {
				if s.Scan() {
					line := s.Text()
					watch.onLine(line)
				}
			}
		}()
	}

	return err
}

// Stop kills our underlying `asterisk -r`.
func (watch *Watch) Stop() {
	watch.isStopped = true
	if watch.cmd != nil {
		watch.cmd.Process.Kill()
	}
}
