package config

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
	
	"github.com/spf13/viper"
)

var Configuration Config

type Config struct {
	Global        GlobalConfig
	Discord       DiscordConfig
	DiscordClient *discordgo.Session
	DataStore     DataStoreConfig
}

type GlobalConfig struct {
	Debug  bool
	Prefix string
}

type DataStoreConfig struct {
	Path   string
	Client *sql.DB
}

type DiscordConfig struct {
	Token string
}

func initViper() error {
	log.Debug("Reading config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/popsicles")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found: %v", err)
		} else {
			log.Fatalf(" Unknown error occured while reading config. error: %v", err)
		}
	}
	err := viper.Unmarshal(&Configuration)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	viper.WatchConfig()

	log.Infof("Using config file found at %v", viper.GetViper().ConfigFileUsed())

	return err
}

func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	if Configuration.Global.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Enabled DEBUG logging level")
	}
}

func initDiscord() (*discordgo.Session, error) {
	log.Debug("Configuring Discord client")
	discord, err := discordgo.New("Bot " + Configuration.Discord.Token)
	if err != nil {
		log.Error("Error creating discord client")
		return discord, err
	}
	return discord, err
}

func initDatastore() {
	db, err := sql.Open("sqlite3", Configuration.DataStore.Path)
	if err != nil {
		log.Fatalf("Unable to open the datastore: %v", err)
	}
	Configuration.DataStore.Client = db
}

func init() {

	// Build config
	err := initViper()
	if err != nil {
		log.Fatal("Unable to init config. Bye.")
	}
	if len(Configuration.Global.Prefix) != 1 {
		log.Fatal("Please check the configured prefix.")
	}

	// Configure logger
	initLogging()

	// Init datastore
	initDatastore()

	// Configure Discord Client
	Configuration.DiscordClient, err = initDiscord()
	if err != nil {
		log.Fatal("No discord client could be created. The bot cannot function. Exiting..")
	}

}
