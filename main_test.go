package string2time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSince(t *testing.T) {
	t.Parallel()

	var now = time.Now()
	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)
	result, err := st.Since("since 3/15/2022")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2022, time.Month(3), 15, 0, 0, 0, 0, now.Location()), result.From)
	assert.Equal(t, time.Time{}, result.To)

	result, err = st.Since("since May 8, 2009 5:57:51 PM")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2009, time.Month(5), 8, 17, 57, 51, 0, now.Location()), result.From)
	assert.Equal(t, time.Time{}, result.To)

	result, err = st.Since("since yesterday")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location()), result.From)
	assert.Equal(t, time.Time{}, result.To)

	result, err = st.Since("since yesterday at 4pm")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()-1, 16, 0, 0, 0, now.Location()), result.From)
	assert.Equal(t, time.Time{}, result.To)

	result, err = st.Since("since yesterday at 13:34:32")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()-1, 13, 34, 32, 0, now.Location()), result.From)
	assert.Equal(t, time.Time{}, result.To)
}
