package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
)


func PlayStory(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {

  options := interactionCreate.ApplicationCommandData().Options
  prompt := options[0].StringValue()
  lang := "br"

  if 2 == len(options) {
   options[1].StringValue()
  }

	// Create a new interaction and send images to discord chat
  currentInteraction := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Generating audio '%s'", prompt),
		},
	}
	session.InteractionRespond(interactionCreate.Interaction, &currentInteraction)

  // Generate speech to text
  voice := voices.Portuguese
  if("jp" == lang){
    voice = voices.Japanese
  }

  speech := htgotts.Speech{Folder: "audio", Language: voice, Handler: &handlers.Native{}}
  filepath, err := speech.CreateSpeechFile(prompt, fmt.Sprintf("%s-%d", session.State.SessionID, time.Now().UnixNano()))
  if nil != err {
    err.Error()
  }

  fileBytes , err := os.ReadFile(filepath)
  if nil != err {
    err.Error()
  }

  newMessage := "Generated story here :)"
  newContent := discordgo.WebhookEdit{
		Content: &newMessage,
		Files:  []*discordgo.File{
	    { 
		    Name:        "Story.mp3",
		    Reader:      bytes.NewBuffer(fileBytes),
		    ContentType: "audio/mp3",
	    },
    },
	}
  
  // Send the audio to discord chat
  _, err = session.InteractionResponseEdit(interactionCreate.Interaction, &newContent)
  if nil != err {
    log.Println(err.Error())
  }
}
