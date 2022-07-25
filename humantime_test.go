package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var now = time.Now()

func init() {
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

var TestParseTimeStringTestCases = map[string]time.Time{
	"  at  3pm":  time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location()),
	"at  12pm":   time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location()),
	"  at  12am": time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
}

func TestParseTimeString(t *testing.T) {
	t.Parallel()

	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)

	result, err := st.parseTimeString(now, "  at  3pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location()), result)

	result, err = st.parseTimeString(now, "at 12pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location()), result)

	result, err = st.parseTimeString(now, "at 12am")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), result)
}

func TestParseDatePhrase(t *testing.T) {
	t.Parallel()

	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)

	result, err := st.parseDatePhrase("  1pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, now.Location()), result)

	result, err = st.parseDatePhrase("yesterday  at  3pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()-1, 15, 0, 0, 0, now.Location()), result)

	result, err = st.parseDatePhrase("last tuesday   at   3pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()-time.Tuesday)-7, 15, 0, 0, 0, now.Location()), result)

	result, err = st.parseDatePhrase("next sunday   at   12pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()-time.Sunday)+7, 12, 0, 0, 0, now.Location()), result)

	result, err = st.parseDatePhrase("this sunday   at   12am")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()-time.Sunday), 0, 0, 0, 0, now.Location()), result)
}
