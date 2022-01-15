package main

import (
	"fmt"
	"os"
)

func man() {
	dir := "../database/sessions/1/2020115001"
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile(dir+"/attendance1_image.txt", []byte("aaditya"), 0777)

}
