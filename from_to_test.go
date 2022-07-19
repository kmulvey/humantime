package string2time

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
}
