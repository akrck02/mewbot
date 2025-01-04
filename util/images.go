package util

import (
	"encoding/base64"
	"log"
)

// Generate png image from base 64
func GeneratePngFromBase64(filename string, b64 string) []byte {

  // Decode base 64 to byte array
  unbased, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Println("Cannot decode b64.")
		return nil
	}

  return unbased
}
