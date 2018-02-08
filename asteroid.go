package asteroid

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

// SIPHeaderWatch watches for SIP Debug messages and fires a callback for all calls
// observed.
type SIPHeaderWatch struct {
	// Config
	asteriskPath string
	headers      []string
	cb           func(callID string, headers map[string]string)

	// Internals
	reader    io.Reader
	cmd       *exec.Cmd
	isStopped bool
}

// NewSIPHeaderWatch for a set of headers to watch for, and a callback to run when a call
// is observed.
func NewSIPHeaderWatch(
	headersToWatch []string,
	callback func(callID string, headers map[string]string),
) SIPHeaderWatch {
	ret := SIPHeaderWatch{
		asteriskPath: "/usr/sbin/asterisk",
		headers:      headersToWatch,
		cb:           callback,
		isStopped:    false,
	}
	return ret
}

// SetAsteriskPath to the Asterisk executable.
func (watch *SIPHeaderWatch) SetAsteriskPath(path string) {
	watch.asteriskPath = path
}

// SetReader instead of Asterisk executable.
func (watch *SIPHeaderWatch) SetReader(reader io.Reader) {
	watch.reader = reader
}

// Start watching for headers.
func (watch *SIPHeaderWatch) Start() error {
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
			var (
				callID       string
				headerValues map[string]string
			)

			headerValues = make(map[string]string)
			for !watch.isStopped {
				if s.Scan() {
					line := s.Text()
					if len(line) > 0 {
						if len(callID) == 0 {
							fmt.Sscanf(line, "Call-ID: %s", &callID)
						} else {
							for _, key := range watch.headers {
								var val string
								if fmt.Sscanf(line, key+": %s", &val); len(val) > 0 {
									headerValues[key] = val
									break
								}
							}
						}
					} else {
						if len(callID) > 0 {
							// Fire callback if all headers are present
							allHeadersPresent := true
							for _, h := range watch.headers {
								if len(headerValues[h]) <= 0 {
									allHeadersPresent = false
								}
							}

							if allHeadersPresent {
								watch.cb(callID, headerValues)
							}

							// Reset our state
							callID = ""
							headerValues = make(map[string]string)
						}
					}
				}
			}
		}()
	}

	return err
}

// Stop kills our underlying `asterisk -r`.
func (watch *SIPHeaderWatch) Stop() {
	watch.isStopped = true
	if watch.cmd != nil {
		watch.cmd.Process.Kill()
	}
}
