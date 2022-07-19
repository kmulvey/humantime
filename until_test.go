package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUntil(t *testing.T) {
	t.Parallel()

	var now = time.Now()
	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)

	result, err := st.Until("until 3/15/2026")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(2026, time.Month(3), 15, 0, 0, 0, 0, now.Location()), result.To)

	result, err = st.Until("until May 8, 2009 5:57:51 PM")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(2009, time.Month(5), 8, 17, 57, 51, 0, now.Location()), result.To)

	result, err = st.Until("until tomorrow")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()), result.To)

	result, err = st.Until("until tomorrow at 4pm")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()+1, 16, 0, 0, 0, now.Location()), result.To)

	result, err = st.Until("until tomorrow at 13:34:32")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()+1, 13, 34, 32, 0, now.Location()), result.To)

	result, err = st.Until("until 2pm")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 14, 00, 00, 0, now.Location()), result.To)
}
