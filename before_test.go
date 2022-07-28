package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var TestBeforeTestCases = map[string]time.Time{
	"before 3/15/2026":                time.Date(2026, time.Month(3), 15, 0, 0, 0, 0, today.Location()),
	"before May 8, 2009 5:57:51 PM":   time.Date(2009, time.Month(5), 8, 17, 57, 51, 0, today.Location()),
	"before tomorrow":                 time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, today.Location()),
	"before tomorrow at 4pm":          time.Date(today.Year(), today.Month(), today.Day()+1, 16, 0, 0, 0, today.Location()),
	"before tomorrow at 13:34:32":     time.Date(today.Year(), today.Month(), today.Day()+1, 13, 34, 32, 0, today.Location()),
	"before 2pm":                      time.Date(today.Year(), today.Month(), today.Day(), 14, 00, 00, 0, today.Location()),
	"before next tuesday at 05:23:43": time.Date(today.Year(), today.Month(), today.Day()-int(today.Weekday()-time.Tuesday)+7, 5, 23, 43, 0, today.Location()),
}

func TestBefore(t *testing.T) {
	t.Parallel()

	var now = time.Now()
	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)

	for input, expected := range TestBeforeTestCases {
		result, err := st.Before(input)
		assert.NoError(t, err)
		assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
		assert.Equal(t, expected, result.To)
	}
	result, err := st.Before("before")
	assert.Equal(t, "input must have at least two fields", err.Error())
	assert.Nil(t, result)

	result, err = st.Before("after 2pm")
	assert.Equal(t, "input does not start with 'before'", err.Error())
	assert.Nil(t, result)
}
