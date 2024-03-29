package humantime

import (
	"fmt"
	"strings"
	"time"
)

// Since takes a string starting with the word since
// and parses the remainder as time.Time, examples:
// since 3/15/2022
// since May 8, 2009 5:57:51 PM
// since 2am
// since yesterday
// since yesterday at 4pm
// since yesterday at 13:34:32
func (st *Humantime) Since(input string) (*TimeRange, error) {
	var tr = new(TimeRange)
	tr.To = time.Now().In(st.Location)

	if len(strings.Fields(input)) < 2 {
		return nil, fmt.Errorf("input must have at least two fields: %s", input)
	}
	if !strings.HasPrefix(input, "since") {
		return nil, fmt.Errorf("input does not start with 'since': %s", input)
	}

	var err error
	tr.From, err = st.parseDatePhrase(strings.ReplaceAll(input, "since ", ""))
	return tr, err
}
