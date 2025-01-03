package util

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"io"
	"log"
	"os"
)

// Generate png image from base 64
func GeneratePngFromBase64(filename string, b64 string) *io.Reader {
	unbased, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Println("Cannot decode b64")
		return nil
	}

	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
		log.Println("Bad png")
		return nil
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println("Cannot open file")
		return nil
	}

	png.Encode(f, im)
  var reader io.Reader = (*os.File)(f)
  return &reader
}
