package main

import (
	"fmt"
	"os"
	"os/signal"
	"popsicles-bot/internal/config"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// DC is a DiscordClient shortcut for readability
// DS is a DataStoreClient shortcut for readability
var (
	DC               = config.Configuration.DiscordClient
	DS               = config.Configuration.DataStore.Client
	prefix           = config.Configuration.Global.Prefix
	numberRegex      = regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	disableDave bool = false
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the message is our own so we can ignore those (something with loops)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, prefix) {

		// TODO: make a case switch instead of endless if statements
		if m.Author.ID == "288046134361063424" && disableDave {
			message := "From tato, with love: ( ° ͜ʖ͡°)╭∩╮"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
			return
		}

		// toggleannoy
		if strings.HasPrefix(m.Content, prefix+"toggleannoy") {
			if m.Author.ID == "188032617793323008" {
				if disableDave {
					disableDave = false
					message := "A certain person will now be allowed to use the bot"
					_, err := s.ChannelMessageSend(m.ChannelID, message)
					if err != nil {
						log.Errorf("Error sending message: %v", err)
					}
					return
				} else {
					disableDave = true
					message := "A certain person will now be blocked from using the bot"
					_, err := s.ChannelMessageSend(m.ChannelID, message)
					if err != nil {
						log.Errorf("Error sending message: %v", err)
					}
					return
				}
			} else {
				message := "you are not authorized to perform this action"
				_, err := s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			}
		}

		// Check if the message is "help"
		if strings.HasPrefix(m.Content, prefix+"help") {
			message := "Available commands:\n source, farenheit, celsius toggleannoy"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
			return
		}

		// source
		if strings.HasPrefix(m.Content, prefix+"source") {
			_, err := s.ChannelMessageSend(m.ChannelID, "Source code can be found at https://github.com/ihulsbus/popsicles-bot")
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
			return
		}

		// convert to farenheit
		if strings.HasPrefix(m.Content, prefix+"farenheit") {
			messages, err := convertToFarenheit(m)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			for i := range messages {
				_, err := s.ChannelMessageSend(m.ChannelID, messages[i])
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
			}
			return
		}

		//convert to celsius
		if strings.HasPrefix(m.Content, prefix+"celsius") {
			messages, err := convertToCelsius(m)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			for i := range messages {
				_, err := s.ChannelMessageSend(m.ChannelID, messages[i])
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
			}
			return
		}

		// Set Height
		if strings.HasPrefix(m.Content, prefix+"setheight") {
			message, err := setHeight(m)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			_, err = s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
		}

		// Get Height
		if strings.HasPrefix(m.Content, prefix+"height") {
			response, err := getHeight(m)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}

			_, err = s.ChannelMessageSend(m.ChannelID, response)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
		}

		// dirty boi
		if strings.HasPrefix(m.Content, prefix+"girth") || strings.HasPrefix(m.Content, prefix+"setgirth") {
			message, _ := girth(m)
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
			return
		}

		// oh boi
		if strings.HasPrefix(m.Content, prefix+"countdown") {
			message := countdown(m)
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
			return
		}

		// shitlords
		if strings.HasPrefix(m.Content, prefix+"shitlords") {
			if m.ChannelID == "796541275549073488" || m.ChannelID == "871809242305273896" {
				message, err := shitlord(m)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				}

				_, err = s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return

			} else {
				message := "Who? What? I do not recognise such filth. Ask me something else..."
				_, err := s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
			}
		}

	} else {
		return
	}
}

func main() {
	log.Info(fmt.Sprintf("Bot prefix is set to %v", prefix))
	log.Info("Starting main loop")

	// Make sure we close the datastore on exit
	defer DS.Close()
	if err := setupDatastore(); err != nil {
		log.Fatalf("Exiting as the datastore is not functional: %v", err)
	}

	// Register handler for incoming messages
	DC.AddHandler(messageCreate)

	DC.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	err := DC.Open()
	if err != nil {
		log.Fatalf("Error opening discord session: %v", err)
	}

	log.Infoln("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Infoln("Stopping bot.")

	// Cleanly close the discord session
	DC.Close()
	log.Infoln("All cleaned up and done. Bye")
}
