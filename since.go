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
// since 2am
// since yesterday
// since yesterday at 4pm
// since yesterday at 13:34:32
func (st *String2Time) Since(input string) (*TimeRange, error) {
	var tr = new(TimeRange)
	tr.To = time.Now().In(st.Location)

	if len(strings.Fields(input)) < 2 {
		return nil, errors.New("input must have two fields")
	}
	if !strings.HasPrefix(input, "since") {
		return nil, errors.New("input does not start with 'since'")
	}

	var err error
	tr.From, err = st.parseDatePhrase(strings.ReplaceAll(input, "since ", ""))
	return tr, err
}

func (st *String2Time) parseDatePhrase(input string) (time.Time, error) {
	var tr = new(TimeRange)

	// is the whole thing a date?
	if date, err := dateparse.ParseIn(input, st.Location, dateparse.RetryAmbiguousDateWithSwap(true)); err == nil {
		return date, nil
	}

	var nextEleIsTime bool
	var inputArr = strings.Fields(input)
	for i := 0; i < len(inputArr); i++ {
		if nextEleIsTime {
			var err = st.parseTimeOrDateString(tr, inputArr[i])
			if err != nil {
				return time.Time{}, err
			}
			return tr.From, nil
		} else if syn, found := TimeSynonyms[inputArr[i]]; found {
			tr.From = syn(st.Location)
		} else if inputArr[i] == "at" {
			nextEleIsTime = true
		} else if len(inputArr) == 1 {
			// this block is time only and assumes the time is for today e.g. "2am"
			var now = time.Now().In(st.Location)
			tr.From = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, st.Location)
			var err = st.parseTimeOrDateString(tr, inputArr[i])
			if err != nil {
				return time.Time{}, err
			}
			return tr.From, nil
		}
	}

	return tr.From, nil
}
