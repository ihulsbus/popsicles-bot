package main

import (
	"os"
	"os/signal"
	"popsicles-bot/internal/config"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// DC is a DiscordClient shortcut for readability
// DS is a DataStoreClient shortcut for readability
var (
	DC = config.Configuration.DiscordClient
	DS = config.Configuration.DataStore.Client

	integerOptionMinValue float64 = -99999999
	heightOptionMinValue  float64 = 50
	sizeOptionMinValue    float64 = 0.001

	// All commands and options must have a description
	// Commands/options without description will fail the registration
	// of the command.
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "source",
			Description: "Get a link to the bot's source code",
		},
		{
			Name:        "farenheit",
			Description: "Convert celsius values to farenheit",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "celsius",
					Description: "temperature value",
					MinValue:    &integerOptionMinValue,
					MaxValue:    99999999,
					Required:    true,
				},
			},
		},
		{
			Name:        "celsius",
			Description: "Convert farenheit to celsius",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "farenheit",
					Description: "temperature value",
					MinValue:    &integerOptionMinValue,
					MaxValue:    99999999,
					Required:    true,
				},
			},
		},
		{
			Name:        "height",
			Description: "Get the height for the given user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user you want to get the height from",
					Required:    true,
				},
			},
		},
		{
			Name:        "setheight",
			Description: "set your height",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "height",
					Description: "height in Centimeters",
					MinValue:    &heightOptionMinValue,
					MaxValue:    230,
					Required:    true,
				},
			},
		},
		{
			Name:        "girth",
			Description: "Get the ...",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user you want to get the girth from",
					Required:    true,
				},
			},
		},
		{
			Name:        "setgirth",
			Description: "Well we don't need to explain this do we?",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "size",
					Description: "size in millimeters",
					MinValue:    &sizeOptionMinValue,
					MaxValue:    1,
					Required:    true,
				},
			},
		},
		{
			Name:        "countdown",
			Description: "Get the days remaining till X",
		},
		{
			Name:        "shitlords",
			Description: "Get the current shitlord or crown a new one",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"source":    source,
		"farenheit": convertToFarenheit,
		"celsius":   convertToCelsius,
		"height":    getHeight,
		"setheight": setHeight,
		"girth":     getGirth,
		"setgirth":  setGirth,
		"countdown": countdown,
		"shitlords": shitlord,
	}
)

func main() {
	log.Info("Starting main loop")

	// Make sure we close the datastore on exit
	if err := setupDatastore(); err != nil {
		log.Fatalf("Exiting as the datastore is not functional: %v", err)
	}

	defer DS.Close()

	// Register handler for incoming messages
	DC.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	DC.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := DC.Open()
	if err != nil {
		log.Fatalf("Error opening discord session: %v", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := DC.ApplicationCommandCreate(DC.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	defer DC.Close()

	log.Infoln("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Infoln("Stopping bot.")

	log.Println("Removing commands...")

	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.
	registeredCommands, err = DC.ApplicationCommands(DC.State.User.ID, "")
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := DC.ApplicationCommandDelete(DC.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	// Cleanly close the discord session
	log.Infoln("All cleaned up and done. Bye")
}
