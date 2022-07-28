package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var TestAgoTestCases = map[string]time.Time{
	"3 days ago":   time.Date(today.Year(), today.Month(), today.Day()-3, today.Hour(), today.Minute(), today.Second(), 0, today.Location()),
	"14 years ago": time.Date(today.Year()-14, today.Month(), today.Day(), today.Hour(), today.Minute(), today.Second(), 0, today.Location()),
	"1 year 2 months 3 days 4 hours 5 minutes 6 seconds ago": time.Date(today.Year()-1, today.Month()-time.Month(2), today.Day()-3, today.Hour()-4, today.Minute()-5, today.Second()-6, 0, today.Location()),
}

func TestAgo(t *testing.T) {
	t.Parallel()

	var now = time.Now()
	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)

	for input, expected := range TestAgoTestCases {
		result, err := st.Ago(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result.From)
		assert.Equal(t, now.Round(time.Second), result.To.Round(time.Second))
	}

	result, err := st.Ago(" ago")
	assert.Equal(t, "input must have at least three fields:  ago", err.Error())
	assert.Nil(t, result)

	result, err = st.Ago("2 years days ago")
	assert.Equal(t, "number of input fields must be even: 2 years days ago", err.Error())
	assert.Nil(t, result)

	result, err = st.Ago("DD years ago")
	assert.Equal(t, "error parsing large units: DD years ago, err: error parsing time: DD years, err: strconv.Atoi: parsing \"DD\": invalid syntax", err.Error())
	assert.Nil(t, result)

	result, err = st.Ago("DD seconds ago")
	assert.Equal(t, "error parsing small units: DD seconds ago, err: time: invalid duration \"s\"", err.Error())
	assert.Nil(t, result)

	result, err = st.Ago("14 years before")
	assert.Equal(t, "input does not end with 'ago'", err.Error())
	assert.Nil(t, result)
}

func TestParseLargeUnits(t *testing.T) {
	t.Parallel()

	var today = time.Now()
	var expected = time.Date(today.Year()-1, today.Month()-time.Month(2), today.Day()-3, today.Hour(), today.Minute(), today.Second(), 0, today.Location())
	var result, nextString, err = parseLargeUnits("1 year 2 months 3 days 4 hours 5 minutes 6 seconds", today.Location())
	assert.NoError(t, err)
	assert.Equal(t, &expected, result)
	assert.Equal(t, "4 hours 5 minutes 6 seconds", nextString)
}

func TestParseSmallUnits(t *testing.T) {
	t.Parallel()

	var result, err = parseSmallUnits("4 hours 5 minutes 6 seconds")
	assert.NoError(t, err)
	d, err := time.ParseDuration("4h5m6s")
	assert.NoError(t, err)
	assert.Equal(t, d, result)
}
