package humantime

import (
	"fmt"
	"strings"
)

// FromTo takes a string in the format from [date phrase] to [date phrase]
// and parses the remainder as time.Time, examples:
// from yesterday to today
// from May 8, 2009 5:57:51 PM to Sep 12, 2021 3:21:22 PM
func (st *Humantime) FromTo(input string) (*TimeRange, error) {
	var tr = new(TimeRange)

	if !strings.HasPrefix(input, "from ") {
		return nil, fmt.Errorf("first arg must be 'from': %s", input)
	}
	if !strings.Contains(input, " to ") {
		return nil, fmt.Errorf("input must contain 'to': %s", input)
	}

	var fromDateStr, toDateStr, found = strings.Cut(strings.ReplaceAll(input, "from ", ""), " to ")
	if !found {
		return nil, fmt.Errorf("input must contain 'to': %s", input)
	}

	var err error
	tr.From, err = st.parseDatePhrase(fromDateStr)
	if err != nil {
		return nil, err
	}

	tr.To, err = st.parseDatePhrase(toDateStr)
	if err != nil {
		return nil, err
	}

	return tr, nil
}
