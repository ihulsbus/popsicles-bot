package main

import (
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
	"github.com/martinlindhe/unit"
)

func convertToFarenheit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var message string

	input := float64(i.ApplicationCommandData().Options[0].IntValue())
	temp := unit.FromCelsius(input)
	temp1 := math.Round((temp.Fahrenheit() * 100) / 100)

	message = fmt.Sprintf("%v° Celsius is %v° Farenheit", input, temp1)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func convertToCelsius(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var message string

	input := float64(i.ApplicationCommandData().Options[0].IntValue())
	temp := unit.FromFahrenheit(input)
	temp1 := math.Round((temp.Celsius() * 100) / 100)

	message = fmt.Sprintf("%v° Farenheit is %v° Celsius", input, temp1)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
