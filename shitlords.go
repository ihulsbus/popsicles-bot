package main

import "github.com/bwmarrin/discordgo"

func shitlord(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// var message string

	message := "It was fun while it lasted. For now, the last shitlord has been relieved from his duties. But don't be sad, the next shitlord will be appointed soon enough.."
	// //Get last shitlord to check timestamp
	// lastShitlord, err := getLastShitlord()
	// if err != nil {
	// 	log.Errorf("Error retrieving last shitlord: %v", err)
	// }

	// if DateEqual(lastShitlord.Timestamp, time.Now()) {
	// 	message = "Today's shitlord has already been crowned. scroll up."
	// 	return message, nil
	// }

	// if m.ChannelID == "796541275549073488" {
	// 	user.ID = "222510172705259521"
	// } else if m.ChannelID == "871809242305273896" {
	// 	user.ID = "188032617793323008"
	// } else {
	// 	user.ID = "859070011846688808"
	// }
	// m.Mentions = append(m.Mentions, &user)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
