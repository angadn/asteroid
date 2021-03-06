package asteroid

import (
	"fmt"
	"strings"
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

	// TODO: Modify to watch CSeq: 1* BYE, or maybe Reason?
	watch.onLine = func(line string) {
		if strings.Contains(line, "destroying SIP dialog") {
			var callID string
			if fmt.Sscanf(
				line, "%s destroying SIP dialog %s Method:", &callID, &callID,
			); len(callID) > 0 {
				callback(callID[1 : len(callID)-1])
			}
		}
	}

	return watch
}
