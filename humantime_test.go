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
	"  at  3pm":  time.Date(today.Year(), today.Month(), today.Day(), 15, 0, 0, 0, today.Location()),
	"at  12pm":   time.Date(today.Year(), today.Month(), today.Day(), 12, 0, 0, 0, today.Location()),
	"  at  12am": time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()),
}

var TestParseDatePhraseTestCases = map[string]time.Time{
	"    1pm":                 time.Date(today.Year(), today.Month(), today.Day(), 13, 0, 0, 0, today.Location()),
	"yesterday  at  3pm":      time.Date(today.Year(), today.Month(), today.Day()-1, 15, 0, 0, 0, today.Location()),
	"last tuesday   at   3pm": time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Tuesday)-7, 15, 0, 0, 0, today.Location()),
	"next sunday   at   12pm": time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Sunday)+7, 12, 0, 0, 0, today.Location()),
	"this sunday   at   12am": time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Sunday), 0, 0, 0, 0, today.Location()),
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
}
