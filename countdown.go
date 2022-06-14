package main

import (
	"fmt"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
)

func countdown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var message string

	dateFormat := "2006-01-02"
	t, _ := time.Parse(dateFormat, "2022-06-12")
	duration := time.Until(t)
	roundedDuration := int64(math.RoundToEven(duration.Hours() / 24))

	if roundedDuration <= 0 {
		message = "hooray! it's happening!"
	} else if roundedDuration < 69 && roundedDuration > 59 {
		roundedDuration = 69
		message = fmt.Sprintf("No. no. We don't go lower than %v", roundedDuration)
	} else if roundedDuration == 42 {
		message = fmt.Sprintf("%v. The answer to the ultimate question of life, the universe, and everything", roundedDuration)
	} else {
		message = fmt.Sprintf("There are %v days remaining", roundedDuration)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
