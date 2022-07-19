package humantime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Ago takes a string starting with the word since
// and parses the remainder as time.Duration, examples:
// 3 hours ago
// 8 days and three hours ago
func (st *String2Time) Ago(input string) (*TimeRange, error) {
	var tr = new(TimeRange)

	var inputArr = strings.Fields(input)
	if len(inputArr) < 3 {
		return tr, errors.New("input must have three fields")
	}

	var multiple, err = strconv.Atoi(inputArr[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %s, err: %w", input, err)
	}

	var unit, found = DurationWords[inputArr[1]]
	if !found {
		return nil, fmt.Errorf("could not parse %s", input)
	}

	tr.From = time.Now().Add(time.Duration(multiple*-1) * unit)
	tr.To = time.Now().In(st.Location)

	return tr, nil
}
