package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var TestAfterTestCases = map[string]time.Time{
	"after 3/15/2022":              time.Date(2022, time.Month(3), 15, 0, 0, 0, 0, today.Location()),
	"after May 8, 2009 5:57:51 PM": time.Date(2009, time.Month(5), 8, 17, 57, 51, 0, today.Location()),
	"after yesterday":              time.Date(today.Year(), today.Month(), today.Day()-1, 0, 0, 0, 0, today.Location()),
	"after yesterday at 4pm":       time.Date(today.Year(), today.Month(), today.Day()-1, 16, 0, 0, 0, today.Location()),
	"after yesterday at 13:34:32":  time.Date(today.Year(), today.Month(), today.Day()-1, 13, 34, 32, 0, today.Location()),
	"after 2am":                    time.Date(today.Year(), today.Month(), today.Day(), 02, 00, 00, 0, today.Location()),
}

func TestAfter(t *testing.T) {
	t.Parallel()

	var now = time.Now()
	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)

	for input, expected := range TestAfterTestCases {
		result, err := st.After(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result.From)
		assert.Equal(t, now.Round(time.Second), result.To.Round(time.Second))
	}

	// error cases
	result, err := st.After("after ")
	assert.Equal(t, "input must have at least two fields", err.Error())
	assert.Nil(t, result)

	result, err = st.After("before 3pm ")
	assert.Equal(t, "input does not start with 'after'", err.Error())
	assert.Nil(t, result)

}
