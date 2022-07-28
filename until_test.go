package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var TestUntilTestCases = map[string]time.Time{
	"until 3/15/2026":              time.Date(2026, time.Month(3), 15, 0, 0, 0, 0, today.Location()),
	"until May 8, 2009 5:57:51 PM": time.Date(2009, time.Month(5), 8, 17, 57, 51, 0, today.Location()),
	"until tomorrow":               time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, today.Location()),
	"until tomorrow at 4pm":        time.Date(today.Year(), today.Month(), today.Day()+1, 16, 0, 0, 0, today.Location()),
	"until tomorrow at 13:34:32":   time.Date(today.Year(), today.Month(), today.Day()+1, 13, 34, 32, 0, today.Location()),
	"until 2pm":                    time.Date(today.Year(), today.Month(), today.Day(), 14, 00, 00, 0, today.Location()),
}

func TestUntil(t *testing.T) {
	t.Parallel()

	var now = time.Now()
	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)

	for input, expected := range TestUntilTestCases {
		result, err := st.Until(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result.To)
		assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	}

	result, err := st.Until("ago ")
	assert.Equal(t, "input must have at least two fields: ago ", err.Error())
	assert.Nil(t, result)

	result, err = st.Until("tomorrow until ")
	assert.Equal(t, "input does not start with 'until': tomorrow until ", err.Error())
	assert.Nil(t, result)
}
