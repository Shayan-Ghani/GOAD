package main

import (
	"log"
	"os"

	"github.com/Shayan-Ghani/GOAD/config"
	"github.com/Shayan-Ghani/GOAD/internal/delivery/cli"
)

func main() {
    client := cli.NewClient()
    if err := client.Run(os.Args[1:], config.ItemSvcAddr, config.TagSvcAddr); err != nil {
        log.Fatal(err)
    }
}