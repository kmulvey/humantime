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

var TimeSynonyms = map[string]func() time.Time{
	"yesterday": func() time.Time {
		var now = time.Now().Add(time.Hour * 24)
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	},
	"tomorrow": func() time.Time {
		var now = time.Now()
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	},
}

func NewString2Time() (*String2Time, error) {

	var err error
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

	return new(String2Time), nil
}

func main() {
	fmt.Println("vim-go")
}

func Parse(input string) error {
	var inputArr = strings.Fields(input)
	if len(inputArr) < 2 {
		return errors.New("input must have two fields")
	}
	fmt.Println(inputArr)

	return nil
}

/*
func LintSince(inputArr []string) (bool, error) {
	if len(inputArr) < 2 {
		return false, errors.New("since must have at least one field after it")
	}
}
*/

// Since takes a string starting with the word since
// and parses the remainder as time.Duration
// since yesterday
// since 3/15/2022
// since yesterday at 4pm
func Since(input string) (time.Duration, error) {
	var inputArr = strings.Fields(input)
	if len(inputArr) < 2 {
		return 0, errors.New("input must have two fields")
	}
	if inputArr[0] != "since" {
		return 0, errors.New("input does not start with 'since'")
	}

	var date time.Time
	var nextEleIsTime bool
	for i := 1; i < len(inputArr); i++ {
		if syn, found := TimeSynonyms[inputArr[i]]; found {
			date = syn()
		} else if inputArr[i] == "am" || inputArr[i] == "pm" {
		}
	}

	if len(inputArr) == 2 {
		parsedTime, err := dateparse.ParseAny(inputArr[1])
		if err != nil {
			return 0, fmt.Errorf("error parsing since [date]: %w", err)
		}
		return time.Since(parsedTime), nil
	}

	return 0, nil // TODO
}

func (st *String2Time) parseTimeOrDateString(input string) (time.Time, error) {
	if AMRegex.MatchString(input) {
		var hourString = strings.ReplaceAll(input, "am", "")
		var hourNum, err = strconv.Atoi(hourString)
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}

		return time.Date(0, 0, 0, hourNum, 0, 0, 0, st.Location), nil
	} else if PMRegex.MatchString(input) {
		var hourString = strings.ReplaceAll(input, "pm", "")
		var hourNum, err = strconv.Atoi(hourString)
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}

		return time.Date(0, 0, 0, hourNum, 0, 0, 0, st.Location), nil
	} else if ExactTimeRegex.MatchString(input) {
		var timeArr = strings.Split(input, ":")

		var err error
		var hour int
		var minute int
		var second int
		hour, err = strconv.Atoi(strings.ReplaceAll(timeArr[0], ":", ""))
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}
		minute, err = strconv.Atoi(strings.ReplaceAll(timeArr[1], ":", ""))
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
		}
		if len(timeArr) == 3 {
			hour, err = strconv.Atoi(strings.ReplaceAll(timeArr[2], ":", ""))
			if err != nil {
				return time.Time{}, fmt.Errorf("error parsing time: %s, err: %w", input, err)
			}
		}
		return time.Date(0, 0, 0, hour, minute, second, 0, st.Location), nil
	} else if DateDashRegex.MatchString(input) {
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

	return time.Time{}, errors.New("unable to parse date: " + input)
}
