package main

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func convertToFarenheit(m *discordgo.MessageCreate) ([]string, error) {
	var messages []string
	celsius, farenheit, err := getFarenheit(m.Content)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error converting temps to Farenheit: %v", err))
		log.Error(err)
		return messages, err
	}

	for i := range farenheit {
		messages = append(messages, fmt.Sprintf("%v° Celsius is %v° Farenheit", celsius[i], farenheit[i]))
	}

	return messages, nil
}

func convertToCelsius(m *discordgo.MessageCreate) ([]string, error) {
	var messages []string
	farenheit, celsius, err := getCelcius(m.Content)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error converting temps to Celsius: %v", err))
		log.Error(err)
		return messages, err
	}

	for i := range celsius {
		messages = append(messages, fmt.Sprintf("%v° Farenheit is %v° Celsius", farenheit[i], celsius[i]))
	}

	return messages, nil
}
