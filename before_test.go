package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBefore(t *testing.T) {
	t.Parallel()

	var now = time.Now()
	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)

	result, err := st.Before("before 3/15/2026")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(2026, time.Month(3), 15, 0, 0, 0, 0, now.Location()), result.To)

	result, err = st.Before("before May 8, 2009 5:57:51 PM")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(2009, time.Month(5), 8, 17, 57, 51, 0, now.Location()), result.To)

	result, err = st.Before("before tomorrow")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()), result.To)

	result, err = st.Before("before tomorrow at 4pm")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()+1, 16, 0, 0, 0, now.Location()), result.To)

	result, err = st.Before("before tomorrow at 13:34:32")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day()+1, 13, 34, 32, 0, now.Location()), result.To)

	result, err = st.Before("before 2pm")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), 14, 00, 00, 0, now.Location()), result.To)

	result, err = st.Before("before next tuesday at 05:23:43")
	assert.NoError(t, err)
	assert.Equal(t, now.Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, time.Date(now.Year(), now.Month(), today.Day()-int(today.Weekday()-time.Tuesday)+7, 5, 23, 43, 0, now.Location()), result.To)
}
