package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	pn, err := New(os.Getenv("PN_PACKAGENAME"), os.Getenv("PN_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	err = pn.Login(os.Getenv("PN_USERNAME"), os.Getenv("PN_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}

	devices, err := pn.Devices()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		fmt.Printf("%+v", device)
	}
}
