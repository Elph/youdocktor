package main

import (
	"fmt"
	"regexp"
	"time"
)

// rangeDate is useful to iterate between to dates
func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

// getRegexGroups returns matching named groups for a regular expression evaluation
func getRegexGroups(input string, expression *regexp.Regexp) (result map[string]string, err error) {
	match := expression.FindStringSubmatch(input)
	result = make(map[string]string)
	if len(match) == 0 {
		err = fmt.Errorf("No groups found on expression: %s", input)
		return
	}
	for i, name := range expression.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return
}

// Returns a valid representation of a date to be sent to Youtrack API
func makeTimestamp(v time.Time) int64 {
	return v.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
