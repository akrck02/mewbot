package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/akrck02/mewbot/util"
	"github.com/bwmarrin/discordgo"
)

// Dall-e mini API response
type DalleMiniResponse struct {
	Images []string `json:"images"`
}

const DALLE_MINI_GENERATE_ENDPOINT = "https://bf.dallemini.ai/generate"

// Generate image discord commands
func GenerateImageCommand(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	
  // Get command options
  options := interactionCreate.ApplicationCommandData().Options
  prompt := options[0].StringValue()
	log.Printf("Generating images with prompt: %s", prompt)

	// Marshall input parameters
	jsonData, err := json.Marshal(map[string]string{
		"prompt": prompt,
	})

	if nil != err {
		log.Fatal(err)
	}

	// Create a new interaction and send images to discord chat
	currentInteraction := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Generating images for description '%s'", prompt),
		},
	}
	session.InteractionRespond(interactionCreate.Interaction, &currentInteraction)

	// Call dall-e mini api
	result := callGenerateImageEndpoint(jsonData)
	if nil == result {
		log.Fatal("Error occurred on Dall-e mini API call, please try again later.")
	}

	// Create the new data for the interaction
  newMessage := fmt.Sprintf("Generated images for '%s'", prompt)
  newContent := discordgo.WebhookEdit{
		Content: &newMessage,
		Files:  []*discordgo.File{},
	}

	// Convert images from base 64 to png
	for index, base64Image := range result.Images {
		
    filename := fmt.Sprintf("image_%s_%d.png", interactionCreate.Interaction.ID, index)
    filepath := fmt.Sprintf("tmp/%s", filename)
		image := util.GeneratePngFromBase64(filepath, base64Image)

    // if the image cannot be created, ignore
		if nil == image {
			break
		}
  
    // Get image data as reader and remove the temporal file
    var reader io.Reader = bytes.NewReader(image)

    // Convert reader to discord file and add it to interaction
	  discordFile := convertPngImageToDiscordFile(filename, *&reader)
  	newContent.Files = append(newContent.Files, discordFile)
  
  }

	// Send the images to discord chat
  _, err = session.InteractionResponseEdit(interactionCreate.Interaction, &newContent)
  if nil != err {
    log.Println(err.Error())
  }
}

// Call dall-e mini api generate endpoint
func callGenerateImageEndpoint(content []byte) *DalleMiniResponse {

	response, err := http.Post(DALLE_MINI_GENERATE_ENDPOINT, "application/json", bytes.NewBuffer(content))
  if nil != err {
		print(err.Error())
		return nil
	}
	defer response.Body.Close()

	// if response if not OK
	if 200 != response.StatusCode {
		v, _ := io.ReadAll(response.Body)
		println(string(v))
		return nil
	}

	// Decode the response to struct
	var result DalleMiniResponse
	json.NewDecoder(response.Body).Decode(&result)

	return &result
}

// Create a discord file from png image
func convertPngImageToDiscordFile(filename string, reader io.Reader) *discordgo.File {
	return &discordgo.File{
		Name:        filename,
		Reader:      reader,
		ContentType: "image/png",
	}
}
