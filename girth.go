package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func getGirth(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var requestingUser *discordgo.User

	if i.User != nil {
		requestingUser = i.User
	} else {
		requestingUser = i.Member.User
	}
	requestedUser := i.ApplicationCommandData().Options[0].UserValue(nil)

	message := fmt.Sprintf("%v is fancying %v :O #MeToo", requestingUser.Mention(), requestedUser.Mention())

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func setGirth(s *discordgo.Session, i *discordgo.InteractionCreate) {
	message := "You wish ya dirty wanker"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

}
