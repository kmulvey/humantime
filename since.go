package string2time

import (
	"errors"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// Since takes a string starting with the word since
// and parses the remainder as time.Duration, examples:
// since 3/15/2022
// since May 8, 2009 5:57:51 PM
// since yesterday
// since yesterday at 4pm
// since yesterday at 13:34:32
func (st *String2Time) Since(input string) (*TimeRange, error) {
	var tr = new(TimeRange)

	var inputArr = strings.Fields(input)
	if len(inputArr) < 2 {
		return tr, errors.New("input must have two fields")
	}
	if inputArr[0] != "since" {
		return tr, errors.New("input does not start with 'since'")
	}

	var date time.Time
	var err error

	// is the whole thing a date?
	if date, err = dateparse.ParseIn(strings.Join(inputArr[1:], " "), st.Location, dateparse.RetryAmbiguousDateWithSwap(true)); err == nil {
		tr.From = date
		return tr, nil
	}

	var nextEleIsTime bool
	for i := 1; i < len(inputArr); i++ {
		if nextEleIsTime {
			var err = st.parseTimeOrDateString(tr, inputArr[i])
			if err != nil {
				return tr, err
			}
			return tr, nil
		} else if syn, found := TimeSynonyms[inputArr[i]]; found {
			tr.From = syn(st.Location)
			tr.To = time.Now().In(st.Location)
		} else if inputArr[i] == "at" {
			nextEleIsTime = true
		}
	}

	return tr, nil
}
