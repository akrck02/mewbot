package commands

import "github.com/bwmarrin/discordgo"

// Ping discord command implementation
func PingCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong! :ping_pong:",
		},
	})
}
