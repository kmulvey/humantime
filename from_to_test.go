package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromTo(t *testing.T) {
	t.Parallel()

	var now = time.Now()
	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)
	result, err := st.FromTo("from yesterday to today")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location()), result.From)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), result.To)

	result, err = st.FromTo("from 8am to 6pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location()), result.From)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location()), result.To)

	result, err = st.FromTo("from May 8, 2009 5:57:51 PM to Sep 12, 2021 3:21:22 PM")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2009, time.Month(5), 8, 17, 57, 51, 0, now.Location()), result.From)
	assert.Equal(t, time.Date(2021, time.Month(9), 12, 15, 21, 22, 0, now.Location()), result.To)
}
