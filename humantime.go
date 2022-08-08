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
func (v *TimeRange) Get() TimeRange {
	return *v
}

// Set fulfils the flag.Value interface https://pkg.go.dev/flag#Value
// must end in the format " in [timezone]" e.g. "3pm in America/New_York"
func (v *TimeRange) Set(s string) error {
	var inputArr = strings.Fields(s)
	if len(inputArr) < 3 {
		return fmt.Errorf("input must have at least three fields: %s", s)
	}

	var location = time.Local
	var err error
	if index := strings.Index(s, " in "); index > -1 {
		location, err = time.LoadLocation(inputArr[len(inputArr)-1])
		if err != nil {
			return err
		}
		s = s[:index]
	}

	st, err := NewString2Time(location)
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
func NewString2Time(loc *time.Location) (*Humantime, error) {

	var st = new(Humantime)
	st.Location = loc

	// init regexs
	st.ExactTimeRegex = regexp.MustCompile(exactTime)
	st.SynonymRegex = regexp.MustCompile(synonyms)
	st.AtTimeRegex = regexp.MustCompile(atTime)
	st.WeekdayRegex = regexp.MustCompile(weekdays)
	st.AMOrPMRegex = regexp.MustCompile(amORpm)

	return st, nil
}

// Parse is the entry point for parsing English input and performs the
// switching between different phrase types
func (st *Humantime) Parse(input string) (*TimeRange, error) {

	input = strings.ToLower(input)

	switch {
	case strings.Contains(input, "since"):
		return st.Since(input)
	case strings.Contains(input, "ago"):
		return st.Ago(input)
	case strings.Contains(input, "til"):
		return st.Until(input)
	case strings.Contains(input, "before"):
		return st.Before(input)
	case strings.Contains(input, "after"):
		return st.After(input)
	case strings.Contains(input, "from") && strings.Contains(input, "to"):
		return st.FromTo(input)
	}

	return nil, fmt.Errorf("unsupported format: %s", input)
}

// parseTimeString reads phrases only containing time, examples:
// 2am
// 7pm
// 04:12:43 -- this format assumes 24h i.e. no a/pm
func (st *Humantime) parseTimeString(timestamp time.Time, input string) (time.Time, error) {
	input = strings.ReplaceAll(input, "at", "")
	input = strings.TrimSpace(input)

	if result := st.AMOrPMRegex.FindString(input); result != "" {
		var period = result[len(result)-2:]
		var hourNum, err = strconv.Atoi(result[:len(result)-2])
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing hour (%s) in: %s, err: %w", result[:len(result)-2], input, err)
		}

		switch {
		case hourNum > 12:
			return time.Time{}, fmt.Errorf("error parsing hour (%d) in: %s, err: hour cannot be > 12", hourNum, input)
		case result == "12am":
			return timestamp, nil
		case period == "am" || result == "12pm": // have to check for noon
			return timestamp.Add(time.Duration(hourNum) * time.Hour), nil
		default:
			return timestamp.Add(time.Duration(hourNum+12) * time.Hour), nil
		}

	} else if st.ExactTimeRegex.MatchString(input) {
		var timeArr = strings.Split(input, ":")

		var err error
		var hour int
		var minute int
		var second int
		hour, err = strconv.Atoi(strings.ReplaceAll(timeArr[0], ":", ""))
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing hour in: %s, err: %w", input, err)
		} else if hour > 23 {
			return time.Time{}, fmt.Errorf("error parsing hour (%d) in: %s, err: hour cannot be > 23", hour, input)
		}

		minute, err = strconv.Atoi(strings.ReplaceAll(timeArr[1], ":", ""))
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing minute in: %s, err: %w", input, err)
		} else if minute > 59 {
			return time.Time{}, fmt.Errorf("error parsing minute (%d) in: %s, err: minute cannot be > 59", minute, input)
		}
		if len(timeArr) == 3 {
			second, err = strconv.Atoi(strings.ReplaceAll(timeArr[2], ":", ""))
			if err != nil {
				return time.Time{}, fmt.Errorf("error parsing second in: %s, err: %w", input, err)
			} else if second > 59 {
				return time.Time{}, fmt.Errorf("error parsing second (%d) in: %s, err: second cannot be > 59", second, input)
			}
		}

		return timestamp.Add(time.Duration(hour) * time.Hour).Add(time.Duration(minute) * time.Minute).Add(time.Duration(second) * time.Second), nil
	}
	return time.Time{}, errors.New("unable to parse date: " + input)
}

// parseDatePhrase parses dates, examples:
// yesterday
// yesterday at 3pmp
// May 8, 2009 5:57:51 PM
// 3/15/2022
// next tuesday at 12am
func (ht *Humantime) parseDatePhrase(input string) (time.Time, error) {

	if date, err := dateparse.ParseIn(input, ht.Location, dateparse.RetryAmbiguousDateWithSwap(true)); err == nil {
		return date, nil
	}

	var inputCopy = strings.TrimSpace(input) // so we can use the original in errors
	var now = time.Now().In(ht.Location)
	var nilTime = time.Time{} // used for if() testing
	var timestamp time.Time   // this is the return val that we incrementally add to each time through the loop
	var i int                 // count iterations to prevent infinitely looping
	for inputCopy != "" {
		if result := ht.WeekdayRegex.FindString(inputCopy); result != "" {
			var resultArr = strings.Fields(result)
			if len(resultArr) != 2 {
				return time.Time{}, fmt.Errorf("could not parse weekday: %s in input: %s", resultArr[1], input)
			}

			var weekday, found = StringToWeekdays[resultArr[1]]
			if !found {
				return time.Time{}, fmt.Errorf("could not parse weekday: %s in input: %s", resultArr[1], input)
			}

			switch resultArr[0] {
			case "last":
				timestamp = time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()-weekday)-7, 0, 0, 0, 0, ht.Location)
			case "this":
				timestamp = time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()-weekday), 0, 0, 0, 0, ht.Location)
			case "next":
				timestamp = time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()-weekday)+7, 0, 0, 0, 0, ht.Location)
			default:
				return time.Time{}, fmt.Errorf("could not parse weekday: %s in input: %s", resultArr[1], input)
			}

			inputCopy = strings.Replace(inputCopy, result, "", 1)
		} else if result := ht.SynonymRegex.FindString(inputCopy); result != "" {
			var syn, _ = TimeSynonyms[result] // ignore second return val as how could it not be found?
			inputCopy = strings.Replace(inputCopy, result, "", 1)
			timestamp = syn(ht.Location)
		} else if result := ht.AtTimeRegex.FindString(inputCopy); result != "" {
			var err error
			if timestamp == nilTime { // no day specified, assume today e.g. "3pm"
				timestamp = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, ht.Location)
			}
			timestamp, err = ht.parseTimeString(timestamp, result)
			if err != nil {
				return time.Time{}, err
			}
			inputCopy = strings.Replace(inputCopy, result, "", 1)
		} else if i == 5 { // catch all so we dont loop forever
			return time.Time{}, fmt.Errorf("could not parse %s", input)
		}
		inputCopy = strings.TrimSpace(inputCopy)
		i++
	}
	return timestamp, nil
}
