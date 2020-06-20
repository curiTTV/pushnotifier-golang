package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	pn := New(os.Getenv("PN_PACKAGENAME"), os.Getenv("PN_APITOKEN"))
	if err := pn.Login(os.Getenv("PN_USERNAME"), os.Getenv("PN_PASSWORD")); err != nil {
		log.Fatal(err)
	}

	devices, err := pn.Devices()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		fmt.Printf("%+v", device)
	}

	if err := pn.Refresh(); err != nil {
		log.Fatal(err)
	}

	if err := pn.Text(nil, "das ist ein test"); err != nil {
		log.Fatal(err)
	}
}
