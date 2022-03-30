package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func setHeight(m *discordgo.MessageCreate) (string, error) {
	var message string
	height, err := getHeightFromMessage(m.Content)
	if err != nil {
		err = errors.New(fmt.Sprint("Unable to get the height from the message: %v", err))
		log.Error(err)
		return message, err
	}
	uid, err := convertStrToInt(m.Author.ID)
	if err != nil {
		err = errors.New(fmt.Sprint("Unable to convert author ID to int: %v", err))
		log.Error(err)
		return message, err
	}
	if (height >= 50) && (height <= 230) {
		log.Debugf("executing setHeight with values %v, %v", m.Author.ID, height)
		setHeightInStore(uid, height)
		message = fmt.Sprintf("height set to %v Centimeters", height)

		return message, nil
	} else {
		err = errors.New("Invalid height given. Height needs to be between 50 and 230 centimeters.")
		return message, err
	}
}

func getHeight(m *discordgo.MessageCreate) (string, error) {
	var message string

	uid, err := convertStrToInt(m.Mentions[len(m.Mentions)-1].ID)
	if err != nil {
		log.Errorf("Unable to convert author ID to int: %v", err)
	}

	height, err := getStoredHeight(uid)
	if err != nil {
		if err == sql.ErrNoRows {
			message = "This user did not set their height yet."
			return message, nil
		} else {
			err = errors.New("Something went wrong. Please try again later.")
			return message, err
		}
	}

	if m.Mentions[len(m.Mentions)-1].ID == "271075764156235777" {

		message = m.Mentions[len(m.Mentions)-1].Mention() + fmt.Sprintf(" is %v Centimeters tall. Just Perfect", height)

		return message, nil
	} else if (m.Mentions[len(m.Mentions)-1].ID == "288046134361063424") && (height >= 220) {

		message = m.Mentions[len(m.Mentions)-1].Mention() + fmt.Sprintf(" is %v Centimeters tall. I Guess your circumference is larger than that of the sun", height)

		return message, nil
	} else {

		message = m.Mentions[len(m.Mentions)-1].Mention() + fmt.Sprintf(" is %v Centimeters tall.", height)

		return message, nil
	}
}
