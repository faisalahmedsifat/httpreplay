package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"httpreplay/internal/proxy"
	"httpreplay/internal/store"
)

func RunCapture() {
	// httpreplay capture --port=8080 --target=http://api.production.com --output=requests.ndjson

	fs := flag.NewFlagSet("capture", flag.ExitOnError)
	port := fs.Int("port", 8080, "port to capture from")
	target := fs.String("target", "http://api.example.com", "target to capture from")
	output := fs.String("output", "requests.ndjson", "output file")

	fs.Parse(os.Args[2:])

	if *target == "" {
		log.Fatal("Missing required flag: --target")
	}

	if *port <= 0 || *port > 65535 {
		log.Fatal("Invalid port: ", *port)
	}

	fmt.Printf("Starting httpreplay proxy on port %d -> forwarding to %s\n", *port, *target)
	fmt.Printf("Capturing requests to %s...\n", *output)

	file, encoder, err := store.NewWriter(*output)
	if err != nil {
		log.Fatal("Failed to create store: ", err)
	}
	defer file.Close()

	proxyHandler := proxy.NewReverseProxy(*target, encoder)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, proxyHandler))
}
