package cmd

import (
	"fmt"
	"log"
	"net/url"

	"httpreplay/internal/replay"
	"httpreplay/internal/store"

	"github.com/urfave/cli"
)

func RunReplay(c *cli.Context) error {
	targetFlag := c.String("target")
	fileFlag := c.String("file")

	var newTarget *url.URL
	var err error

	if targetFlag != "" {
		newTarget, err = url.Parse(targetFlag)
		if err != nil || newTarget.Scheme == "" || newTarget.Host == "" {
			log.Fatal("Invalid target URL: ", targetFlag)
		}
		fmt.Printf("Replaying traffic to %s...\n", newTarget.String())
	} else {
		fmt.Println("No target URL provided, using the same target as the capture")
	}

	records, err := store.ReadAll(fileFlag)
	if err != nil {
		log.Fatal("Failed to read file: ", err)
	}

	fmt.Printf("Read %d records from %s...\n", len(records), fileFlag)

	engine := replay.NewReplayer(newTarget)

	err = engine.Playback(records)
	if err != nil {
		log.Fatal("Failed to playback records: ", err)
	}
	return nil
}
