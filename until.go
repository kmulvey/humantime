package humantime

import (
	"fmt"
	"strings"
	"time"
)

// Until takes a string starting with the words until or til
// and parses the remainder as time.Time, examples:
// until 3/15/2022
// until May 8, 2009 5:57:51 PM
// until 2am
// until tomorrow
// until tomorrow at 4pm
// until tomorrow at 13:34:32
func (st *Humantime) Until(input string) (*TimeRange, error) {
	var tr = new(TimeRange)
	tr.From = time.Now().In(st.Location)

	if len(strings.Fields(input)) < 2 {
		return nil, fmt.Errorf("input must have at least two fields: %s", input)
	}
	if !strings.HasPrefix(input, "until") {
		return nil, fmt.Errorf("input does not start with 'until': %s", input)
	}

	var err error
	tr.To, err = st.parseDatePhrase(strings.ReplaceAll(input, "until ", ""))
	return tr, err
}
