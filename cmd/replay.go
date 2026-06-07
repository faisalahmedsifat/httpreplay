package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"httpreplay/internal/replay"
	"httpreplay/internal/store"
)

func RunReplay() {
	replayCmd := flag.NewFlagSet("replay", flag.ExitOnError)
	targetFlag := replayCmd.String("target", "", "The new destination URL to send traffic to")
	fileFlag := replayCmd.String("file", "requests.ndjson", "The file to replay from")

	replayCmd.Parse(os.Args[2:])

	var newTarget *url.URL
	var err error

	if *targetFlag != "" {
		newTarget, err = url.Parse(*targetFlag)
		if err != nil || newTarget.Scheme == "" || newTarget.Host == "" {
			log.Fatal("Invalid target URL: ", *targetFlag)
		}
		fmt.Printf("Replaying traffic to %s...\n", newTarget.String())
	} else {
		fmt.Println("No target URL provided, using the same target as the capture")
	}

	records, err := store.ReadAll(*fileFlag)
	if err != nil {
		log.Fatal("Failed to read file: ", err)
	}
	
	fmt.Printf("Read %d records from %s...\n", len(records), *fileFlag)

	engine := replay.NewReplayer(newTarget)

	err = engine.Playback(records)
	if err != nil {
		log.Fatal("Failed to playback records: ", err)
	}
}
