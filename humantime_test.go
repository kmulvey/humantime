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
	result, err := st.parseDatePhrase("next tomorrow")
	assert.Equal(t, "could not parse next tomorrow", err.Error())
	assert.Equal(t, time.Time{}, result)
}

func TestCLI(t *testing.T) {
	t.Parallel()

	var st, err = NewString2Time(time.UTC)
	assert.NoError(t, err)

	result, err := st.FromTo("from 1/1/2021 to 2/2/2022")
	assert.NoError(t, err)
	assert.Equal(t, "From: 01 Jan 21 00:00 UTC, To: 02 Feb 22 00:00 UTC", result.String())

	err = result.Set("from 1/1/2001 to 2/2/2002 in America/NoExist")
	assert.Equal(t, "unknown time zone America/NoExist", err.Error())

	err = result.Set("from 1 to 2 in America/Denver")
	assert.Equal(t, "error parsingDatePhrase: could not parse 1", err.Error())

	err = result.Set("from 1/1/2001 to 2/2/2002 in America/Denver")
	assert.NoError(t, err)

	var v = result.Get()
	location, err := time.LoadLocation("America/Denver")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2001, time.January, 1, 0, 0, 0, 0, location), v.From.Round(time.Minute))
	assert.Equal(t, time.Date(2002, time.February, 2, 0, 0, 0, 0, location), v.To.Round(time.Minute))
}

func TestParse(t *testing.T) {
	t.Parallel()

	var st, err = NewString2Time(today.Location())
	assert.NoError(t, err)

	result, err := st.Parse("since yesterday")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day()-1, 0, 0, 0, 0, today.Location()), result.From)

	result, err = st.Parse("2 days ago")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day()-2, 0, 0, 0, 0, today.Location()).Day(), result.From.Day())

	result, err = st.Parse("til 3/15/2026 at 00:00:00")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2026, time.Month(3), 15, 0, 0, 0, 0, today.Location()), result.To)

	result, err = st.Parse("before 3/15/2026 at 00:00:00")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2026, time.Month(3), 15, 0, 0, 0, 0, today.Location()), result.To)

	result, err = st.Parse("after 3/15/2006 at 00:00:00")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2006, time.Month(3), 15, 0, 0, 0, 0, today.Location()), result.From)

	result, err = st.Parse("apples")
	assert.Equal(t, "unsupported format: apples", err.Error())
	assert.Nil(t, result)
}
