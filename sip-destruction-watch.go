package asteroid

import (
	"fmt"
)

// SIPDestructionWatch watches for SIP Dialogs being 'really destroyed'.
type SIPDestructionWatch struct {
	Watch
}

// NewSIPDestructionWatch constructs a SIPDestructionWatch.
func NewSIPDestructionWatch(callback func(callID string)) SIPDestructionWatch {
	watch := SIPDestructionWatch{
		Watch: NewWatch(),
	}

	watch.onLine = func(line string) {
		var callID string
		if fmt.Sscanf(
			line, "Really destroying SIP dialog %s Method: INVITE", &callID,
		); len(callID) > 0 {
			callback(callID[1 : len(callID)-1])
		}

	}

	return watch
}
