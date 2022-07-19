package humantime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgo(t *testing.T) {
	t.Parallel()

	var now = time.Now()
	var st, err = NewString2Time(now.Location())
	assert.NoError(t, err)
	result, err := st.Ago("3 days ago")
	assert.NoError(t, err)
	assert.Equal(t, now.Add(time.Hour*24*-3).Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, now.Round(time.Second), result.To.Round(time.Second))

	result, err = st.Ago("14 years ago")
	assert.NoError(t, err)
	assert.Equal(t, now.Add(time.Second*31536000*-14).Round(time.Second), result.From.Round(time.Second))
	assert.Equal(t, now.Round(time.Second), result.To.Round(time.Second))
}
