package main

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func setHeight(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var message string
	var uid int
	var err error

	heightInput := i.ApplicationCommandData().Options[0].IntValue()

	if i.User != nil {
		uid, err = convertStrToInt(i.User.ID)
	} else {
		uid, err = convertStrToInt(i.Member.User.ID)
	}

	if err != nil {
		message = fmt.Sprintf("Unable to convert author ID to int: %v", err)
	} else {
		setHeightInStore(uid, int(heightInput))
		message = fmt.Sprintf("height set to %v Centimeters", heightInput)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func getHeight(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var message string

	user := i.ApplicationCommandData().Options[0].UserValue(nil)

	uid, err := convertStrToInt(user.ID)
	if err != nil {
		log.Errorf("Unable to convert author ID to int: %v", err)
	}

	height, err := getStoredHeight(uid)
	if err != nil {
		if err == sql.ErrNoRows {
			message = "This user did not set their height yet."
		} else {
			message = "Something went wrong. Please try again later."
		}
	}

	if height == 0 {
		message = fmt.Sprintf("%v is %v centimeters tall. RIP %v.", user.Mention(), height, user.Mention())

	} else if user.ID == "271075764156235777" {
		message = fmt.Sprintf("%v is %v Centimeters tall. Just Perfect", user.Mention(), height)

	} else if (user.ID == "288046134361063424") && (height >= 220) {
		message = fmt.Sprintf("%v is %v Centimeters tall. I Guess your circumference is larger than that of the sun", user.Mention(), height)

	} else {
		message = fmt.Sprintf("%v is %v Centimeters tall.", user.Mention(), height)

	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
