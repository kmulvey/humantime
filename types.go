package humantime

import (
	"regexp"
	"time"
)

type Humantime struct {
	*time.Location
	AMOrPMRegex    *regexp.Regexp
	ExactTimeRegex *regexp.Regexp
	SynonymRegex   *regexp.Regexp
	AtTimeRegex    *regexp.Regexp
	WeekdayRegex   *regexp.Regexp
}

type TimeRange struct {
	From time.Time
	To   time.Time
}

// all text is passed through strings.ToLower() before these regexs are evaluated
const ExactTime = `\d{1,2}:\d{1,2}(:\d{1,2})?`                                                // one or two digits, ':', one or two digits, optional: [':' one or two digits]
const AMOrPM = `(\d{1,2}am)|(\d{1,2}pm)`                                                      // one or two digits, followed by 'am' OR [same for pm]
const Synonyms = `(yesterday|today|tomorrow)`                                                 // any of these three words
const AtTime = `(at)?\s*(\d{1,2}am)|(at)?\s*(\d{1,2}pm)|(at)?\s*(\d{1,2}:\d{1,2}(:\d{1,2})?)` // [optional 'at'], any amout of spcace, one or two digits, 'am' OR [same for pm] OR [similar for 00:11:22]
const Weekdays = `(next|last|this)\s*((mon|tues|wed(nes)?|thur(s)?|fri|sat(ur)?|sun)(day)?)`  // any of these three words, any amount of space, any day of the week with optional abbreviation

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

var StringToWeekdays = map[string]time.Weekday{
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"webnesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
	"sunday":    time.Sunday,
}
