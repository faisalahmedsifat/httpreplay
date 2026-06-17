package main

import (
	"httpreplay/cmd"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "HttpReplay",
		Usage: "A CLI tool for capturing and replaying HTTP traffic",
		Commands: []cli.Command{
			{
				Name:   "capture",
				Usage:  "Capture HTTP traffic",
				Action: cmd.RunCapture,
			},
			{
				Name:   "replay",
				Usage:  "Replay captured HTTP traffic",
				Action: cmd.RunReplay,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
