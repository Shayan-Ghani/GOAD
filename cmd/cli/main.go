package main

import (
    "log"
    "os"

    "github.com/Shayan-Ghani/GOAD/internal/delivery/cli"
)

func main() {
    client := cli.NewClient()
    if err := client.Run(os.Args[1:]); err != nil {
        log.Fatal(err)
    }
}