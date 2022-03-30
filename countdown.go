package main

import (
	"fmt"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
)

func countdown(m *discordgo.MessageCreate) string {
	dateFormat := "2006-01-02"
	t, _ := time.Parse(dateFormat, "2022-06-12")
	duration := time.Until(t)
	roundedDuration := int64(math.RoundToEven(duration.Hours() / 24))

	if roundedDuration < 69 {
		roundedDuration = 69
	}
	return fmt.Sprintf("There are %v days remaining", roundedDuration)
}
