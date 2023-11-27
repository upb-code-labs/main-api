package infrastructure

import "time"

func ParseISODate(date string) (time.Time, error) {
	layout := "2006-01-02T15:04"
	return time.Parse(layout, date)
}
