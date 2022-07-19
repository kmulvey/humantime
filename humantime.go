package humantime

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// String fulfils the flag.Value interface https://pkg.go.dev/flag#Value
func (v TimeRange) String() string {
	return fmt.Sprintf("From: %s, To: %s", v.From.Format(time.RFC822), v.To.Format(time.RFC822))
}

// Get fulfils the flag.Getter interface https://pkg.go.dev/flag#Getter
func (v *TimeRange) Get(s string) TimeRange {
	return *v
}

// Set fulfils the flag.Value interface https://pkg.go.dev/flag#Value
func (v *TimeRange) Set(s string) error {
	var st, err = NewString2Time(time.UTC) // TODO not always utc
	if err != nil {
		return err
	}

	if r, err := st.Parse(s); err != nil {
		return err
	} else {
		v.To = r.To
		v.From = r.From
	}
	return nil
}

// NewString2Time is just a constructor
func NewString2Time(loc *time.Location) (*String2Time, error) {

	var err error
	var st = new(String2Time)
	st.Location = loc

	// init regexs
	st.AMRegex, err = regexp.Compile(AM)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", AM, err)
	}
	st.PMRegex, err = regexp.Compile(PM)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", PM, err)
	}
	st.DateSlashRegex, err = regexp.Compile(DateSlash)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", DateSlash, err)
	}
	st.DateDashRegex, err = regexp.Compile(DateDash)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", DateDash, err)
	}
	st.ExactTimeRegex, err = regexp.Compile(ExactTime)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", ExactTime, err)
	}

	return st, nil
}

// Parse is the entry point for parsing English input and performs the
// switching between different phrase types
func (st *String2Time) Parse(input string) (*TimeRange, error) {

	input = strings.ToLower(input)

	if strings.Contains(input, "since") {
		return st.Since(input)
	} else if strings.Contains(input, "ago") {
		return st.Ago(input)
	} else if strings.Contains(input, "from") && strings.Contains(input, "to") {
		return st.FromTo(input)
	}

	return nil, nil
}

func (st *String2Time) parseTimeString(tr *TimeRange, input string) error {
	if st.AMRegex.MatchString(input) {
		var hourString = strings.ReplaceAll(input, "am", "")
		var hourNum, err = strconv.Atoi(hourString)
		if err != nil {
			return fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}

		tr.From = tr.From.Add(time.Duration(hourNum) * time.Hour)
		return nil
	} else if st.PMRegex.MatchString(input) {
		var hourString = strings.ReplaceAll(input, "pm", "")
		var hourNum, err = strconv.Atoi(hourString)
		if err != nil {
			return fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}

		tr.From = tr.From.Add(time.Duration(hourNum+12) * time.Hour)
		return nil
	} else if st.ExactTimeRegex.MatchString(input) {
		var timeArr = strings.Split(input, ":")

		var err error
		var hour int
		var minute int
		var second int
		hour, err = strconv.Atoi(strings.ReplaceAll(timeArr[0], ":", ""))
		if err != nil {
			return fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}
		minute, err = strconv.Atoi(strings.ReplaceAll(timeArr[1], ":", ""))
		if err != nil {
			return fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}
		if len(timeArr) == 3 {
			second, err = strconv.Atoi(strings.ReplaceAll(timeArr[2], ":", ""))
			if err != nil {
				return fmt.Errorf("error parsing time: %s, err: %w", input, err)
			}
		}

		tr.From = tr.From.Add(time.Duration(hour) * time.Hour).Add(time.Duration(minute) * time.Minute).Add(time.Duration(second) * time.Second)
		return nil
	}
	return errors.New("unable to parse date: " + input)
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
			var err = st.parseTimeString(tr, inputArr[i])
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
			var err = st.parseTimeString(tr, inputArr[i])
			if err != nil {
				return time.Time{}, err
			}
			return tr.From, nil
		}
	}

	return tr.From, nil
}
