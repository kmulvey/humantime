package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var TestToFromTestCases = map[string]TimeRange{
	"from yesterday to today": {
		From: time.Date(today.Year(), today.Month(), today.Day()-1, 0, 0, 0, 0, today.Location()),
		To:   time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()),
	},
	"from 8am to 6pm": {
		From: time.Date(today.Year(), today.Month(), today.Day(), 8, 0, 0, 0, today.Location()),
		To:   time.Date(today.Year(), today.Month(), today.Day(), 18, 0, 0, 0, today.Location()),
	},
	"from May 8, 2009 5:57:51 PM to Sep 12, 2021 3:21:22 PM": {
		From: time.Date(2009, time.Month(5), 8, 17, 57, 51, 0, today.Location()),
		To:   time.Date(2021, time.Month(9), 12, 15, 21, 22, 0, today.Location()),
	},
}

func TestFromTo(t *testing.T) {
	t.Parallel()

	var today = time.Now()
	var st, err = NewString2Time(today.Location())
	assert.NoError(t, err)

	for input, expected := range TestToFromTestCases {
		result, err := st.FromTo(input)
		assert.NoError(t, err)
		assert.Equal(t, expected.From, result.From.Round(time.Second))
		assert.Equal(t, expected.To, result.To)
	}

	result, err := st.FromTo("before yesterday")
	assert.Equal(t, "first arg must be 'from': before yesterday", err.Error())
	assert.Nil(t, result)

	result, err = st.FromTo("from yesterday")
	assert.Equal(t, "input must contain 'to': from yesterday", err.Error())
	assert.Nil(t, result)

	result, err = st.FromTo("from yesterday to no")
	assert.Equal(t, "error parsingDatePhrase: could not parse no", err.Error())
	assert.Nil(t, result)

	result, err = st.FromTo("from yesterday to")
	assert.Equal(t, "input must contain ' to ': from yesterday to", err.Error())
	assert.Nil(t, result)

	result, err = st.FromTo("from nope to tomorrow")
	assert.Equal(t, "error parsingDatePhrase: could not parse nope", err.Error())
	assert.Nil(t, result)
}
