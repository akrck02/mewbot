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

func GenerateImageCommand(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	options := interactionCreate.ApplicationCommandData().Options
	generateImages(session, interactionCreate, options[0].StringValue())
}

// Generate images with dall-e mini
func generateImages(session *discordgo.Session, interactionCreate *discordgo.InteractionCreate, prompt string) {

	log.Printf("Generating images with prompt: %s", prompt)

	// Marshall input parameters
	jsonData, err := json.Marshal(map[string]string{
		"prompt": prompt,
	})

	if nil != err {
		log.Fatal(err)
	}

	// Call dall-e mini api
	result := callGenerateImageEndpoint(jsonData)
	if nil == result {
		log.Fatal("Error occurred on Dall-e mini API call, please try again later.")
	}

	// Convert images from base 64 to png
	files := []*discordgo.File{}
	for index, base64Image := range result.Images {
		println(base64Image)
		filename := fmt.Sprintf("image_%d.png", index)
		image := util.GeneratePngFromBase64(filename, base64Image)

		if nil == image {
			break
		}

		discordFile := convertPngImageToDiscordFile(filename, *image)
		files = append(files, discordFile)
	}

	// Create a new interaction and send images to discord chat
	currentInteraction := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Generated images for description '%s'", prompt),
			Files:   files,
		},
	}
	session.InteractionRespond(interactionCreate.Interaction, &currentInteraction)
}

// Call dall-e mini api generate endpoint
func callGenerateImageEndpoint(content []byte) *DalleMiniResponse {

	response, err := http.Post(DALLE_MINI_GENERATE_ENDPOINT, "application/json", bytes.NewBuffer(content))
	if nil != err {
		err.Error()
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
