package lang

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode/utf8"
)

func FixedString(str string, size int, ellipse string) string {
	s := []rune(str)
	l := utf8.RuneCountInString(str)
	el := utf8.RuneCountInString(ellipse)
	if l == size {
		return str
	} else if l > size {
		return string(s[:size-el]) + ellipse
	} else {
		return str + strings.Repeat(" ", size-l)
	}
}

const (
	secsMinute = float64(60)
	secsHour   = float64(60) * secsMinute
	secsDay    = float64(24) * secsHour
	secsYear   = float64(365) * secsDay
)

func FormatDuration(dur time.Duration) string {
	seconds := dur.Seconds()
	years := math.Floor(seconds / secsYear)
	seconds -= years * secsYear
	days := math.Floor(seconds / secsDay)
	seconds -= days * secsDay
	hours := math.Floor(seconds / secsHour)
	seconds -= hours * secsHour
	minutes := math.Floor(seconds / secsMinute)
	seconds -= minutes * secsMinute
	parts := make([]string, 0, 5)
	parts = appendFormat(parts, int(years), "Jahre")
	parts = appendFormat(parts, int(days), "Tage")
	parts = appendFormat(parts, int(hours), "Stunden")
	parts = appendFormat(parts, int(minutes), "Minuten")
	parts = appendFormat(parts, int(seconds), "Sekunden")
	return strings.Join(parts, ", ")
}

func appendFormat(parts []string, i int, unit string) []string {
	if i > 0 {
		if i == 1 {
			unit = unit[:len(unit)-1]
		}
		parts = append(parts, fmt.Sprintf("%d %s", i, unit))
	}
	return parts
}
