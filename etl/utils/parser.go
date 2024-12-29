package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	datePattern = regexp.MustCompile(`^\d{2}/\d{2}/\d{4}$`) // Matches YYYY-MM-DD
)

func ParseDate(value string) time.Time {
	value = strings.TrimSpace(value)

	if value == "" || !datePattern.MatchString(value) {
		return time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	}

	layout := "02/01/2006"

	date, err := time.Parse(layout, value)
	if err != nil {
		log.Printf("Erro ao converter a data '%s': %v", value, err)
		return time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	}

	return date
}

func ParseInt(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}

	// Parse as float to handle scientific notation
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error parsing number '%s': %v", value, err)
		return 0
	}
	return int(f)
}
