package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/akrck02/mewbot/commands"
	"github.com/akrck02/mewbot/configuration"
	"github.com/bwmarrin/discordgo"
)

func main() {
	dg, err := discordgo.New("Bot " + configuration.Global.DiscordBotToken)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(ready)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Mewbot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

// Set ready function
func ready(s *discordgo.Session, event *discordgo.Ready) {

  // Set the playing status.
	s.UpdateGameStatus(1, "meow meow.")
	commands.RegisterCommands(s)
	commands.SetCommands(s)
}
