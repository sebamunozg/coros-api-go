package helpers

import (
	"strconv"
	"time"
)

func FormatActivityDate(date int) string {
	parsed, err := time.Parse("20060102", strconv.Itoa(date))
	if err != nil {
		return strconv.Itoa(date)
	}

	return parsed.Format("2006-01-02")
}
