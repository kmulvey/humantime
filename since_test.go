package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var TestSinceTestCases = map[string]time.Time{
	"since 3/15/2022":              time.Date(2022, time.Month(3), 15, 0, 0, 0, 0, today.Location()),
	"since May 8, 2009 5:57:51 PM": time.Date(2009, time.Month(5), 8, 17, 57, 51, 0, today.Location()),
	"since yesterday":              time.Date(today.Year(), today.Month(), today.Day()-1, 0, 0, 0, 0, today.Location()),
	"since yesterday at 4pm":       time.Date(today.Year(), today.Month(), today.Day()-1, 16, 0, 0, 0, today.Location()),
	"since yesterday at 13:34:32":  time.Date(today.Year(), today.Month(), today.Day()-1, 13, 34, 32, 0, today.Location()),
	"since 2am":                    time.Date(today.Year(), today.Month(), today.Day(), 02, 00, 00, 0, today.Location()),
}

func TestSince(t *testing.T) {
	t.Parallel()

	var today = time.Now()
	var st, err = NewString2Time(today.Location())
	assert.NoError(t, err)

	for input, expected := range TestSinceTestCases {
		result, err := st.Since(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result.From)
		assert.Equal(t, today.Round(time.Second), result.To.Round(time.Second))
	}

	result, err := st.Since("since")
	assert.Equal(t, "input must have at least two fields: since", err.Error())
	assert.Nil(t, result)

	result, err = st.Since("after 4pm")
	assert.Equal(t, "input does not start with 'since': after 4pm", err.Error())
	assert.Nil(t, result)
}
