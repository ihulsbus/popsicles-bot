package main

import (
	"os"
	"os/signal"
	"popsicles-bot/internal/config"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//DC is a DiscordClient shortcut for readability
var DC = config.Configuration.DiscordClient

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
