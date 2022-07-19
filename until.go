package humantime

import (
	"errors"
	"strings"
	"time"
)

// Until takes a string starting with the words until or til
// and parses the remainder as time.Duration, examples:
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
		return nil, errors.New("input must have two fields")
	}
	if !strings.HasPrefix(input, "until") {
		return nil, errors.New("input does not start with 'until'")
	}

	var err error
	tr.To, err = st.parseDatePhrase(strings.ReplaceAll(input, "until ", ""))
	return tr, err
}
