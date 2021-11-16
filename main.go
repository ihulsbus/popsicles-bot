package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"os/signal"
	"popsicles-bot/internal/config"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// DC is a DiscordClient shortcut for readability
// DS is a DataStoreClient shortcut for readability
var (
	DC                     = config.Configuration.DiscordClient
	DS                     = config.Configuration.DataStore.Client
	prefix                 = config.Configuration.Global.Prefix
	numberRegex            = regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	disableDave       bool = false
	shtlrds_timestamp time.Time
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
			message := "Source code can be found at https://github.com/ihulsbus/popsicles-bot"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
			return
		}

		// convert to farenheit
		if strings.HasPrefix(m.Content, prefix+"farenheit") {
			celsius, farenheit, err := getFarenheit(m.Content)
			if err != nil {
				log.Error("Error converting temps to Farenheit: %v", err)
				_, err := s.ChannelMessageSend(m.ChannelID, err.Error())
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			}
			for index := range farenheit {
				message := fmt.Sprintf("%v° Celsius is %v° Farenheit", celsius[index], farenheit[index])
				_, err := s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
			}
			return
		}

		//convert to celsius
		if strings.HasPrefix(m.Content, prefix+"celsius") {
			farenheit, celsius, err := getCelcius(m.Content)
			if err != nil {
				log.Error("Error converting temps to Celsius: %v", err)
				_, err := s.ChannelMessageSend(m.ChannelID, err.Error())
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			}
			for index := range celsius {
				message := fmt.Sprintf("%v° Farenheit is %v° Celsius", farenheit[index], celsius[index])
				_, err := s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
			}
			return
		}

		// set height
		if strings.HasPrefix(m.Content, prefix+"setheight") {
			height, err := getHeight(m.Content)
			if err != nil {
				log.Errorf("Unable to get the height from the message: %v", err)
				_, err = s.ChannelMessageSend(m.ChannelID, `Unable to get height from message. Do "setheight <length in Centimeters>`)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			}
			uid, err := convertStrToInt(m.Author.ID)
			if err != nil {
				log.Errorf("Unable to convert author ID to int: %v", err)
				_, err = s.ChannelMessageSend(m.ChannelID, "Something went wrong. Please try again later.")
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			}
			if (height >= 50) && (height <= 230) {
				log.Debugf("executing setHeight with values %v, %v", m.Author.ID, height)
				setHeight(uid, height)
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("height set to %v Centimeters", height))
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			} else {
				_, err = s.ChannelMessageSend(m.ChannelID, "Invalid height given. Height needs to be between 50 and 230 centimeters.")
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			}
		}

		// Get Height
		if strings.HasPrefix(m.Content, prefix+"height") {
			uid, err := convertStrToInt(m.Mentions[len(m.Mentions)-1].ID)
			if err != nil {
				log.Errorf("Unable to convert author ID to int: %v", err)
			}
			height, err := getStoredHeight(uid)
			if err != nil {
				if err == sql.ErrNoRows {
					_, err = s.ChannelMessageSend(m.ChannelID, "This user did not set their height yet.")
					if err != nil {
						log.Errorf("Error sending message: %v", err)
					}
					return
				} else {
					_, err = s.ChannelMessageSend(m.ChannelID, "Something went wrong. Please try again later.")
					if err != nil {
						log.Errorf("Error sending message: %v", err)
					}
					return
				}

			}
			if m.Mentions[len(m.Mentions)-1].ID == "271075764156235777" {
				message := m.Mentions[len(m.Mentions)-1].Mention() + fmt.Sprintf(" is %v Centimeters tall. Just Perfect", height)
				_, err = s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			} else if (m.Mentions[len(m.Mentions)-1].ID == "288046134361063424") && (height >= 220) {
				message := m.Mentions[len(m.Mentions)-1].Mention() + fmt.Sprintf(" is %v Centimeters tall. I Guess your circumference is larger than that of the sun", height)
				_, err = s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			} else {
				message := m.Mentions[len(m.Mentions)-1].Mention() + fmt.Sprintf(" is %v Centimeters tall.", height)
				_, err = s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				return
			}
		}

		// dirty boi
		if strings.HasPrefix(m.Content, prefix+"girth") || strings.HasPrefix(m.Content, prefix+"setgirth") {
			message := "You wish ya dirty wanker"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
			return
		}

		// oh boi
		if strings.HasPrefix(m.Content, prefix+"countdown") {
			timeFormat := "2006-01-02"
			t, _ := time.Parse(timeFormat, "2022-06-12")
			duration := time.Until(t)
			roundedDuration := int64(math.RoundToEven(duration.Hours() / 24))

			if roundedDuration < 69 {
				roundedDuration = 69
			}
			message := fmt.Sprintf("There are %v days remaining", roundedDuration)
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			}
			return
		}

		// shitlords
		if strings.HasPrefix(m.Content, prefix+"shitlords") {

			if m.ChannelID == "796541275549073488" || m.ChannelID == "871809242305273896" {

				if DateEqual(shtlrds_timestamp, time.Now()) {
					message := "This command is on cooldown. Try again tomorrow."
					_, err := s.ChannelMessageSend(m.ChannelID, message)
					if err != nil {
						log.Errorf("Error sending message: %v", err)
					}
					return
				}

				message := "It was fun while it lasted. For now, the last shitlord has been relieved from his duties. But don't be sad, the next shitlord will be appointed soon enough.."
				_, err := s.ChannelMessageSend(m.ChannelID, message)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}

				shtlrds_timestamp = time.Now()

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
