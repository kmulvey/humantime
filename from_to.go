package string2time

import (
	"fmt"
	"strings"
)

// Ago takes a string starting with the word since
// and parses the remainder as time.Duration, examples:
// 3 hours ago
// 8 days and three hours ago
func (st *String2Time) FromTo(input string) (*TimeRange, error) {
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
