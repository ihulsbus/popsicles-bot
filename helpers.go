package main

import (
	"strconv"
	"time"
)

func convertStrToInt(input string) (int, error) {
	output, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	return output, nil
}

func getHeight(message string) (int, error) {
	match := numberRegex.FindString(message)
	i, err := convertStrToInt(match)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
