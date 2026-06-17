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
				Name:  "capture",
				Usage: "Capture HTTP traffic",
				Flags: []cli.Flag{
					cli.IntFlag{Name: "port", Value: 8080, Usage: "port to capture from"},
					cli.StringFlag{Name: "target", Value: "http://api.example.com", Usage: "target to capture from"},
					cli.StringFlag{Name: "output", Value: "requests.ndjson", Usage: "output file"},
				},
				Action: cmd.RunCapture,
			},
			{
				Name:  "replay",
				Usage: "Replay captured HTTP traffic",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "target", Usage: "the new destination URL to send traffic to"},
					cli.StringFlag{Name: "file", Value: "requests.ndjson", Usage: "the file to replay from"},
				},
				Action: cmd.RunReplay,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
