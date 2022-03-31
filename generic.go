package main

import "github.com/bwmarrin/discordgo"

func source(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Source code can be found at https://github.com/ihulsbus/popsicles-bot",
		},
	})
}
