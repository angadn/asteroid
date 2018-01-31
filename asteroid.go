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
	readCloser io.ReadCloser
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
	}
	return ret
}

// SetAsteriskPath to the Asterisk executable
func (watch *SIPHeaderWatch) SetAsteriskPath(path string) {
	watch.asteriskPath = path
}

// Start watching for headers.
func (watch *SIPHeaderWatch) Start() {
	cmd := exec.Command(watch.asteriskPath, "-r")

	var err error
	if err = cmd.Run(); err == nil {
		if watch.readCloser, err = cmd.StdoutPipe(); err == nil {
			go func() {
				s := bufio.NewScanner(watch.readCloser)
				var (
					callID       string
					headerValues map[string]string
				)

				headerValues = make(map[string]string)
				for s.Scan() {
					line := s.Text()
					if len(line) > 0 {
						if len(callID) == 0 {
							fmt.Sscanf(line, "Call-ID: %s@", &callID)
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
							// Fire callback
							watch.cb(callID, headerValues)

							// Reset our state
							callID = ""
							headerValues = make(map[string]string)
						}
					}
				}
			}()
		}
	}
}
