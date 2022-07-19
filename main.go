package string2time

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 10 days ago
// since yesterday
// since last week
// since

type String2Time struct {
	*time.Location
	AMRegex        *regexp.Regexp
	PMRegex        *regexp.Regexp
	DateSlashRegex *regexp.Regexp
	DateDashRegex  *regexp.Regexp
	ExactTimeRegex *regexp.Regexp
}

type TimeRange struct {
	From time.Time
	To   time.Time
}

const AM = `^\dam`
const PM = `^\dpm`
const DateSlash = `\d{1,2}/\d{1,2}/\d{2,4}`
const DateDash = `\d{1,2}-\d{1,2}-\d{2,4}`
const ExactTime = `\d{1,2}:\d{1,2}(:\d{1,2})?` // can detect optional seconds

/*
var Words = []string{
	"since",
	"ago",
	"until",
	"til",
	"after",
	"before",
	"from",
	"to",
}
*/

var DurationWords = map[string]time.Duration{
	"second":  time.Second,
	"seconds": time.Second,
	"minute":  time.Minute,
	"minutes": time.Minute,
	"hour":    time.Hour,
	"hours":   time.Hour,
	"day":     time.Hour * 24,
	"days":    time.Hour * 24,
	"week":    time.Hour * 24 * 7,
	"weeks":   time.Hour * 24 * 7,
	"month":   time.Hour * 24 * 7 * 30, // TODO 30 is probaby wrong here
	"months":  time.Hour * 24 * 7 * 30, // TODO 30 is probaby wrong here
	"year":    time.Second * 31536000,  // TODO 30 is probaby wrong here
	"years":   time.Second * 31536000,  // TODO 30 is probaby wrong here
}

var DurationStringToMilli = map[string]int{
	"second":  time.Now().Second(),
	"seconds": time.Now().Second(),
	"minute":  time.Now().Minute(),
	"minutes": time.Now().Minute(),
	"hour":    time.Now().Hour(),
	"hours":   time.Now().Hour(),
	"day":     time.Now().Day(),
	"days":    time.Now().Day(),
	"week":    time.Now().Day() * 7,
	"weeks":   time.Now().Day() * 7,
	"month":   int(time.Now().Month()),
	"months":  int(time.Now().Month()),
	"year":    time.Now().Year(),
	"years":   time.Now().Year(),
}

var TimeSynonyms = map[string]func(*time.Location) time.Time{
	"yesterday": func(loc *time.Location) time.Time {
		var now = time.Now().Add(time.Hour * -24)
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	},
	"today": func(loc *time.Location) time.Time {
		var now = time.Now()
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	},
	"tomorrow": func(loc *time.Location) time.Time {
		var now = time.Now().Add(time.Hour * 24)
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	},
}

func (v TimeRange) String() string {
	return fmt.Sprintf("From: %s, To: %s", v.From.Format(time.RFC822), v.To.Format(time.RFC822))
}

func (v *TimeRange) Get(s string) TimeRange {
	return *v
}

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

func (st *String2Time) Parse(input string) (*TimeRange, error) {

	input = strings.ToLower(input)

	if strings.Contains(input, "since") {
		return st.Since(input)
	} else if strings.Contains(input, "ago") {
		return st.Ago(input)
	}

	return nil, nil
}

func (st *String2Time) parseTimeOrDateString(tr *TimeRange, input string) error {
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
