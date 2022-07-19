package humantime

import (
	"regexp"
	"time"
)

type Humantime struct {
	*time.Location
	AMRegex        *regexp.Regexp
	PMRegex        *regexp.Regexp
	DateSlashRegex *regexp.Regexp
	DateDashRegex  *regexp.Regexp
	ExactTimeRegex *regexp.Regexp
}

type TimeRange struct {
	From time.Time
	To   time.Time
}

const AM = `^\dam`
const PM = `^\dpm`
const DateSlash = `\d{1,2}/\d{1,2}/\d{2,4}`
const DateDash = `\d{1,2}-\d{1,2}-\d{2,4}`
const ExactTime = `\d{1,2}:\d{1,2}(:\d{1,2})?` // can detect optional seconds

var DurationWords = map[string]time.Duration{
	"second":  time.Second,
	"seconds": time.Second,
	"minute":  time.Minute,
	"minutes": time.Minute,
	"hour":    time.Hour,
	"hours":   time.Hour,
	"day":     time.Hour * 24,
	"days":    time.Hour * 24,
	"week":    time.Hour * 24 * 7,
	"weeks":   time.Hour * 24 * 7,
	"month":   time.Hour * 24 * 7 * 30,
	"months":  time.Hour * 24 * 7 * 30,
	"year":    time.Second * 31536000,
	"years":   time.Second * 31536000,
}

var DurationStringToMilli = map[string]int{
	"second":  time.Now().Second(),
	"seconds": time.Now().Second(),
	"minute":  time.Now().Minute(),
	"minutes": time.Now().Minute(),
	"hour":    time.Now().Hour(),
	"hours":   time.Now().Hour(),
	"day":     time.Now().Day(),
	"days":    time.Now().Day(),
	"week":    time.Now().Day() * 7,
	"weeks":   time.Now().Day() * 7,
	"month":   int(time.Now().Month()),
	"months":  int(time.Now().Month()),
	"year":    time.Now().Year(),
	"years":   time.Now().Year(),
}

var TimeSynonyms = map[string]func(*time.Location) time.Time{
	"yesterday": func(loc *time.Location) time.Time {
		var now = time.Now().Add(time.Hour * -24)
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	},
	"today": func(loc *time.Location) time.Time {
		var now = time.Now()
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	},
	"tomorrow": func(loc *time.Location) time.Time {
		var now = time.Now().Add(time.Hour * 24)
		var y, m, d = now.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, loc)
	},
}
