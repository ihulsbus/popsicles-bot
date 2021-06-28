package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"popsicles-bot/internal/config"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//DC is a DiscordClient shortcut for readability
var (
	DC = config.Configuration.DiscordClient
	// WAC         = config.Configuration.WolframAlphaClient
	numberRegex = regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
)

type (
	// or centigrade
	Celsius    float64
	Fahrenheit float64
)

func getTemperatureValue(message string) []string {
	submatchall := numberRegex.FindAllString(message, -1)
	return submatchall
}

func Celsius2Fahrenheit(c float64) Fahrenheit {
	return Fahrenheit(math.Round((c*9/5+32)*100) / 100)
}

func Fahrenheit2Celsius(f float64) Celsius {
	return Celsius(math.Round(((f-32)*5/9)*100) / 100)
}

func getFarenheit(message string) ([]string, []Fahrenheit, error) {
	var calculations []Fahrenheit
	temperatures := getTemperatureValue(message)
	if len(temperatures) == 0 {
		err := errors.New("no temperatures found in the message")
		return temperatures, calculations, err
	}
	for _, temperature := range temperatures {
		temp, err := strconv.ParseFloat(temperature, 64)
		if err != nil {
			return temperatures, calculations, err
		}
		calculations = append(calculations, Celsius2Fahrenheit(temp))

	}
	return temperatures, calculations, nil
}

func getCelcius(message string) ([]string, []Celsius, error) {
	var calculations []Celsius
	temperatures := getTemperatureValue(message)
	if len(temperatures) == 0 {
		err := errors.New("no temperatures found in the message")
		return temperatures, calculations, err
	}
	for _, temperature := range temperatures {
		temp, err := strconv.ParseFloat(temperature, 64)
		if err != nil {
			return temperatures, calculations, err
		}
		calculations = append(calculations, Fahrenheit2Celsius(temp))

	}
	return temperatures, calculations, nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the message is our own so we can ignore those (something with loops)
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the message is "!hey"
	if strings.HasPrefix(m.Content, "!help") {
		message := "Somebody has been lazy..."
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			log.Errorln("Error sending message: %v", err)
		}
		return
	}

	if strings.HasPrefix(m.Content, "!source") {
		message := "Source code can be found at https://github.com/ihulsbus/popsicles-bot"
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			log.Errorln("Error sending message: %v", err)
		}
		return
	}

	if strings.HasPrefix(m.Content, "!farenheit") {

		celsius, farenheit, err := getFarenheit(m.Content)
		if err != nil {
			log.Error("Error converting temps to Farenheit: %v", err)
			_, err := s.ChannelMessageSend(m.ChannelID, err.Error())
			if err != nil {
				log.Errorln("Error sending message: %v", err)
			}
			return
		}
		for index := range farenheit {
			message := fmt.Sprintf("%v° Celsius is %v° Farenheit", celsius[index], farenheit[index])
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorln("Error sending message: %v", err)
			}
		}
		return
	}
	if strings.HasPrefix(m.Content, "!celsius") {
		farenheit, celsius, err := getCelcius(m.Content)
		if err != nil {
			log.Error("Error converting temps to Celsius: %v", err)
			_, err := s.ChannelMessageSend(m.ChannelID, err.Error())
			if err != nil {
				log.Errorln("Error sending message: %v", err)
			}
			return
		}
		for index := range celsius {
			message := fmt.Sprintf("%v° Farenheit is %v° Celsius", farenheit[index], celsius[index])
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			if err != nil {
				log.Errorln("Error sending message: %v", err)
			}
		}
		return
	}
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}
	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, err := s.ChannelMessageSend(channel.ID, "Popsicles for sale! Get 'em here! Use !help for the available commands")
			if err != nil {
				log.Errorln("Error sending message: %v", err)
			}
			return
		}
	}
}

func main() {
	log.Info("Starting main loop")

	// Register handler for incoming messages
	DC.AddHandler(messageCreate)

	// Register handler for server join messages
	DC.AddHandler(guildCreate)

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
