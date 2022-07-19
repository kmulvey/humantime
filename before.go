package humantime

import (
	"errors"
	"strings"
	"time"
)

// Before takes a string starting with the word before
// and parses the remainder as time.Duration, examples:
// before 3/15/2022
// before May 8, 2009 5:57:51 PM
// before 2am
// before tomorrow
// before tomorrow at 4pm
// before tomorrow at 13:34:32
func (st *Humantime) Before(input string) (*TimeRange, error) {
	var tr = new(TimeRange)
	tr.From = time.Now().In(st.Location)

	if len(strings.Fields(input)) < 2 {
		return nil, errors.New("input must have two fields")
	}
	if !strings.HasPrefix(input, "before") {
		return nil, errors.New("input does not start with 'before'")
	}

	var err error
	tr.To, err = st.parseDatePhrase(strings.ReplaceAll(input, "before ", ""))
	return tr, err
}
