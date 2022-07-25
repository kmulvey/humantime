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

func TestParseTimeString(t *testing.T) {
	t.Parallel()

	var st, err = NewString2Time(today.Location())
	assert.NoError(t, err)

	result, err := st.parseTimeString(today, "  at  3pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day(), 15, 0, 0, 0, today.Location()), result)

	result, err = st.parseTimeString(today, "at 12pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day(), 12, 0, 0, 0, today.Location()), result)

	result, err = st.parseTimeString(today, "at 12am")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()), result)
}

func TestParseDatePhrase(t *testing.T) {
	t.Parallel()

	var st, err = NewString2Time(today.Location())
	assert.NoError(t, err)

	result, err := st.parseDatePhrase("  1pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day(), 13, 0, 0, 0, today.Location()), result)

	result, err = st.parseDatePhrase("yesterday  at  3pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day()-1, 15, 0, 0, 0, today.Location()), result)

	result, err = st.parseDatePhrase("last tuesday   at   3pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Tuesday)-7, 15, 0, 0, 0, today.Location()), result)

	result, err = st.parseDatePhrase("next sunday   at   12pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Sunday)+7, 12, 0, 0, 0, today.Location()), result)

	result, err = st.parseDatePhrase("this sunday   at   12am")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Sunday), 0, 0, 0, 0, today.Location()), result)
}
