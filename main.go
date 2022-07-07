package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// 10 days ago
// since yesterday
// since last week
// since

type String2Time struct {
	*time.Location
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
var AMRegex *regexp.Regexp
var PMRegex *regexp.Regexp
var DateSlashRegex *regexp.Regexp
var DateDashRegex *regexp.Regexp
var ExactTimeRegex *regexp.Regexp

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

var DurationWords = map[string]time.Duration{
	"second": time.Second,
	"minute": time.Minute,
	"hour":   time.Hour,
	"day":    time.Hour * 24,
	"week":   time.Hour * 24 * 7,
	"month":  time.Hour * 24 * 7 * 30,      // TODO 30 is probaby wrong here
	"year":   time.Hour * 24 * 7 * 30 * 12, // TODO 30 is probaby wrong here
}

var DurationWordsPlural = map[string]func(time.Duration) time.Duration{
	"seconds": func(multiplier time.Duration) time.Duration { return multiplier * DurationWords["second"] },
	"minutes": func(multiplier time.Duration) time.Duration { return multiplier * DurationWords["minute"] },
	"hours":   func(multiplier time.Duration) time.Duration { return multiplier * DurationWords["hour"] },
	"days":    func(multiplier time.Duration) time.Duration { return multiplier * DurationWords["day"] },
	"weeks":   func(multiplier time.Duration) time.Duration { return multiplier * DurationWords["week"] },
	"months":  func(multiplier time.Duration) time.Duration { return multiplier * DurationWords["month"] },
	"years":   func(multiplier time.Duration) time.Duration { return multiplier * DurationWords["year"] },
}

var TimeSynonyms = map[string]func(*time.Location) time.Time{
	"yesterday": func(loc *time.Location) time.Time {
		var now = time.Now().Add(time.Hour * -24)
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	},
	"tomorrow": func(loc *time.Location) time.Time {
		var now = time.Now().Add(time.Hour * 24)
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	},
}

func NewString2Time(loc *time.Location) (*String2Time, error) {

	var err error
	var st = new(String2Time)
	st.Location = loc

	// init regexs
	AMRegex, err = regexp.Compile(AM)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", AM, err)
	}
	PMRegex, err = regexp.Compile(PM)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", PM, err)
	}
	DateSlashRegex, err = regexp.Compile(DateSlash)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", DateSlash, err)
	}
	DateDashRegex, err = regexp.Compile(DateDash)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", DateDash, err)
	}
	ExactTimeRegex, err = regexp.Compile(ExactTime)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %s, err: %w", ExactTime, err)
	}

	return st, nil
}

func (st *String2Time) Parse(input string) (*TimeRange, error) {
	var inputArr = strings.Fields(input)
	if len(inputArr) < 2 {
		return nil, errors.New("input must have at least two fields")
	}

	if strings.HasPrefix(input, "since") {
		return st.Since(input)
	}

	return nil, nil
}

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
		//fmt.Println(i, inputArr[i])

		if nextEleIsTime {
			var err = st.parseTimeOrDateString(tr, inputArr[i])
			if err != nil {
				return tr, err
			}
			return tr, nil
		} else if syn, found := TimeSynonyms[inputArr[i]]; found {
			tr.From = syn(st.Location)
		} else if inputArr[i] == "at" {
			nextEleIsTime = true
		}
	}

	return tr, nil //errors.New("could not parse: " + input)
}

func (st *String2Time) parseTimeOrDateString(tr *TimeRange, input string) error {
	if AMRegex.MatchString(input) {
		var hourString = strings.ReplaceAll(input, "am", "")
		var hourNum, err = strconv.Atoi(hourString)
		if err != nil {
			return fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}

		tr.From.Add(time.Duration(hourNum) * time.Hour)
		return nil
	} else if PMRegex.MatchString(input) {
		var hourString = strings.ReplaceAll(input, "pm", "")
		var hourNum, err = strconv.Atoi(hourString)
		if err != nil {
			return fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}

		tr.From = tr.From.Add(time.Duration(hourNum+12) * time.Hour)
		return nil
	} else if ExactTimeRegex.MatchString(input) {
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
	/*
		else if DateDashRegex.MatchString(input) {
			var timeArr = strings.Split(input, "-")
			// TODO only works for MM/DD/YY for now
			var err error
			var day int
			var month int
			var year int
			month, err = strconv.Atoi(strings.ReplaceAll(timeArr[0], "-", ""))
			if err != nil {
				return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
			}
			day, err = strconv.Atoi(strings.ReplaceAll(timeArr[1], "-", ""))
			if err != nil {
				return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
			}
			if len(timeArr) == 3 {
				year, err = strconv.Atoi(strings.ReplaceAll(timeArr[2], "-", ""))
				if err != nil {
					return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
				}
			} else {
				year = time.Now().Year()
			}
			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, st.Location), nil
		} else if DateSlashRegex.MatchString(input) {
			var timeArr = strings.Split(input, "/")
			// TODO only works for MM/DD/YY for now
			var err error
			var day int
			var month int
			var year int
			month, err = strconv.Atoi(strings.ReplaceAll(timeArr[0], "/", ""))
			if err != nil {
				return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
			}
			day, err = strconv.Atoi(strings.ReplaceAll(timeArr[1], "/", ""))
			if err != nil {
				return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
			}
			if len(timeArr) == 3 {
				year, err = strconv.Atoi(strings.ReplaceAll(timeArr[2], "/", ""))
				if err != nil {
					return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
				}
			} else {
				year = time.Now().Year()
			}
			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, st.Location), nil
		}
	*/
	return errors.New("unable to parse date: " + input)
}
