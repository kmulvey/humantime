package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var today = time.Now()

func init() {
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
}

var TestParseTimeStringTestCases = map[string]time.Time{
	"  at  3pm":   time.Date(today.Year(), today.Month(), today.Day(), 15, 0, 0, 0, today.Location()),
	"    5am":     time.Date(today.Year(), today.Month(), today.Day(), 5, 0, 0, 0, today.Location()),
	"at  12pm":    time.Date(today.Year(), today.Month(), today.Day(), 12, 0, 0, 0, today.Location()),
	"  at  12am":  time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()),
	"at 15:50:09": time.Date(today.Year(), today.Month(), today.Day(), 15, 50, 9, 0, today.Location()),
	"  00:00:00 ": time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()),
	"  23:59:59 ": time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, today.Location()),
	"  3:9 ":      time.Date(today.Year(), today.Month(), today.Day(), 3, 9, 0, 0, today.Location()),
}

var TestParseDatePhraseTestCases = map[string]time.Time{
	"    1pm":                     time.Date(today.Year(), today.Month(), today.Day(), 13, 0, 0, 0, today.Location()),
	"at 15:50:09":                 time.Date(today.Year(), today.Month(), today.Day(), 15, 50, 9, 0, today.Location()),
	"yesterday  at  3pm":          time.Date(today.Year(), today.Month(), today.Day()-1, 15, 0, 0, 0, today.Location()),
	"tomorrow  at  23:59:59 ":     time.Date(today.Year(), today.Month(), today.Day()+1, 23, 59, 59, 0, today.Location()),
	"today  at  12:33:42 ":        time.Date(today.Year(), today.Month(), today.Day(), 12, 33, 42, 0, today.Location()),
	"last tuesday   at   3pm":     time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Tuesday)-7, 15, 0, 0, 0, today.Location()),
	"next sunday   at   12pm":     time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Sunday)+7, 12, 0, 0, 0, today.Location()),
	"this sunday   at   12:33:42": time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Sunday), 12, 33, 42, 0, today.Location()),
}

func TestParseTimeString(t *testing.T) {
	t.Parallel()

	var st, err = NewString2Time(today.Location())
	assert.NoError(t, err)

	for input, expected := range TestParseTimeStringTestCases {
		result, err := st.parseTimeString(today, input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	}

	// error cases
	result, err := st.parseTimeString(today, "23pm")
	assert.Equal(t, "error parsing hour (23) in: 23pm, err: hour cannot be > 12", err.Error())
	assert.Equal(t, time.Time{}, result)

	result, err = st.parseTimeString(today, "33:23")
	assert.Equal(t, "error parsing hour (33) in: 33:23, err: hour cannot be > 23", err.Error())
	assert.Equal(t, time.Time{}, result)

	result, err = st.parseTimeString(today, "3:73:12")
	assert.Equal(t, "error parsing minute (73) in: 3:73:12, err: minute cannot be > 59", err.Error())
	assert.Equal(t, time.Time{}, result)

	result, err = st.parseTimeString(today, "3:3:82")
	assert.Equal(t, "error parsing second (82) in: 3:3:82, err: second cannot be > 59", err.Error())
	assert.Equal(t, time.Time{}, result)
}

func TestParseDatePhrase(t *testing.T) {
	t.Parallel()

	var st, err = NewString2Time(today.Location())
	assert.NoError(t, err)

	for input, expected := range TestParseDatePhraseTestCases {
		result, err := st.parseDatePhrase(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	}

	// error cases
	result, err := st.parseTimeString(today, "next tomorrow")
	assert.Equal(t, "unable to parse date: next tomorrow", err.Error())
	assert.Equal(t, time.Time{}, result)
}
