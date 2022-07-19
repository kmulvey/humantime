package humantime

import (
	"errors"
	"strings"
	"time"
)

// After takes a string starting with the word after
// and parses the remainder as time.Duration, examples:
// after 3/15/2022
// after May 8, 2009 5:57:51 PM
// after 2am
// after yesterday
// after yesterday at 4pm
// after yesterday at 13:34:32
func (st *Humantime) After(input string) (*TimeRange, error) {
	var tr = new(TimeRange)
	tr.To = time.Now().In(st.Location)

	if len(strings.Fields(input)) < 2 {
		return nil, errors.New("input must have two fields")
	}
	if !strings.HasPrefix(input, "after") {
		return nil, errors.New("input does not start with 'after'")
	}

	var err error
	tr.From, err = st.parseDatePhrase(strings.ReplaceAll(input, "after ", ""))
	return tr, err
}
