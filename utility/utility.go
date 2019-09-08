package utility

import (
	"time"
)

const (
	FORMAT_YYYYMMDDHHMMSS = "20060102150405"
	FORMAT_YYYYMMDD       = "20060102"
)

func GetCurrentDateTime(format string) string {
	return time.Now().Format(format)
}
