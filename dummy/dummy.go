package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
)

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func main() {
	rb, err := ioutil.ReadFile("001.jpg")
	if err != nil {
		fmt.Println(err)
	}

	var b string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(rb)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		b += "data:image/jpeg;base64,"
	case "image/png":
		b += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	b += toBase64(rb)
	// fmt.Println(b)

	unbased, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		panic("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	im, err := jpeg.Decode(r)
	if err != nil {
		panic("Bad jpeg")
	}

	f, err := os.OpenFile("example.jpeg", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}
	var opt jpeg.Options
	opt.Quality = 100
	jpeg.Encode(f, im, &opt)

}
