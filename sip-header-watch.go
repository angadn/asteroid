package asteroid

import (
	"fmt"
)

// SIPHeaderWatch watches for SIP Debug messages and fires a callback for all calls
// observed.
type SIPHeaderWatch struct {
	// Config
	headers []string
	cb      func(callID string, headers map[string]string)

	// Internals
	Watch
}

// NewSIPHeaderWatch for a set of headers to watch for, and a callback to run when a call
// is observed.
func NewSIPHeaderWatch(
	headersToWatch []string,
	callback func(callID string, headers map[string]string),
) SIPHeaderWatch {
	watch := SIPHeaderWatch{
		Watch:   NewWatch(),
		headers: headersToWatch,
		cb:      callback,
	}

	var (
		callID       string
		headerValues map[string]string
	)

	headerValues = make(map[string]string)

	watch.OnLine = func(line string) {
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

	return watch
}
