package main

import (
	"strconv"
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
