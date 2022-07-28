package humantime

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Ago takes a string starting with the word since
// and parses the remainder as time.Time, examples:
// 3 hours ago
// 8 days and three hours ago
// 1 year 2 months 3 days 4 hours 5 minutes 6 seconds ago
func (st *Humantime) Ago(input string) (*TimeRange, error) {
	var tr = new(TimeRange)
	var inputCopy = input

	// lint the input
	if len(strings.Fields(input)) < 3 {
		return nil, fmt.Errorf("input must have at least three fields: %s", input)
	}
	if !strings.HasSuffix(input, "ago") {
		return nil, errors.New("input does not end with 'ago'")
	}

	// remove stop words
	inputCopy = strings.ReplaceAll(input, "ago", "")
	inputCopy = strings.ReplaceAll(inputCopy, "and", "")
	inputCopy = strings.ReplaceAll(inputCopy, ",", "")
	inputCopy = strings.TrimSpace(inputCopy)

	// more linting
	var inputArr = strings.Fields(inputCopy)
	if len(inputArr)%2 != 0 {
		return nil, fmt.Errorf("number of input fields must be even: %s", input)
	}

	var baseTime, nextString, err = parseLargeUnits(inputCopy, st.Location)
	if err != nil {
		return nil, fmt.Errorf("error parsing large units: %s, err: %w", input, err)
	}

	duration, err := parseSmallUnits(nextString)
	if err != nil {
		return nil, fmt.Errorf("error parsing small units: %s, err: %w", input, err)
	}

	tr.From = *baseTime
	tr.From = tr.From.Add(duration * -1)
	tr.To = time.Now().In(st.Location)

	return tr, nil
}

// parseLargeUnits takes the input string and parses out years, months and days.
// It returns a Time pointer corresponding to now - [result of parsing].
// It also returns the rest of the input string i.e. (full string - the part parsed).
// The string should then be passed to parseSmallUnits().
func parseLargeUnits(input string, loc *time.Location) (*time.Time, string, error) {
	var year int
	var month int
	var day int
	var tempNum int
	var lastIndex int
	var inputArr = strings.Fields(input)

	var matched, err = regexp.MatchString(`year|month|day`, input)
	if err != nil {
		return nil, "", err
	}
	if matched {
		for i, word := range strings.Fields(input) {
			//fmt.Printf("i: %d, tempNum: %d, year: %d, month: %d, day: %d, word: %s \n", i, tempNum, year, month, day, word)
			if i%2 == 0 {
				tempNum, err = strconv.Atoi(word)
				if err != nil {
					return nil, "", fmt.Errorf("error parsing time: %s, err: %w", input, err)
				}
			} else {
				if found := strings.Contains(word, "year"); found {
					year = tempNum
					lastIndex = i
				} else if found := strings.Contains(word, "month"); found {
					month = tempNum
					lastIndex = i
				} else if found := strings.Contains(word, "day"); found {
					day = tempNum
					lastIndex = i
				}
			}
		}
	}

	var today = time.Now()
	var t = time.Date(today.Year()-year, today.Month()-time.Month(month), today.Day()-day, today.Hour(), today.Minute(), today.Second(), 0, loc)
	return &t, strings.Join(inputArr[lastIndex+1:], " "), nil
}

// parseSmallUnits takes the input string and parses out hours, minutes and seconds.
// It does this by converting the string to the time.Duration string format: DhDmDs
func parseSmallUnits(input string) (time.Duration, error) {

	var matched, err = regexp.MatchString(`hour|minute|second`, input)
	if err != nil {
		return 0, err
	}
	if !matched {
		return 0, nil
	}

	hr, err := regexp.Compile(`hour(s)?`)
	if err != nil {
		return 0, err
	}
	input = hr.ReplaceAllString(input, "h")

	mr, err := regexp.Compile(`minute(s)?`)
	if err != nil {
		return 0, err
	}
	input = mr.ReplaceAllString(input, "m")

	sr, err := regexp.Compile(`second(s)?`)
	if err != nil {
		return 0, err
	}
	input = sr.ReplaceAllString(input, "s")

	return time.ParseDuration(strings.ReplaceAll(input, " ", ""))
}
