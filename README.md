# pushnotifier-golang

This library implements all endpoints of [pushnotifier.de](pushnotifier.de)

## Installation

```bash
go get github.com/curiTTV/pushnotifier-golang
```

## Usage

```go
package main

import (
    "fmt"
    "log"
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

    // An app token expires after one year
    if err := pn.Refresh(); err != nil {
        log.Fatal(err)
    }

    if err := pn.Text(nil, "this is a test"); err != nil {
        log.Fatal(err)
    }

    // if err := pn.URL(nil, "https://github.com/curiTTV"); err != nil {
    //  log.Fatal(err)
    // }
}
```
