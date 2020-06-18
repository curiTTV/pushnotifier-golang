package main

import (
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	pn, err := New(os.Getenv("PN_PACKAGENAME"), os.Getenv("PN_APITOKEN"))
	if err != nil {
		t.Error(err)
	}

	err = pn.Login(os.Getenv("PN_USERNAME"), os.Getenv("PN_PASSWORD"))
	if err != nil {
		t.Error(err)
	}
}
