package cmd

import (
	"fmt"
	"log"
	"net/http"

	"httpreplay/internal/proxy"
	"httpreplay/internal/store"

	"github.com/urfave/cli"
)

func RunCapture(c *cli.Context) error {
	port := c.Int("port")
	target := c.String("target")
	output := c.String("output")

	if target == "" {
		log.Fatal("Missing required flag: --target")
	}

	if port <= 0 || port > 65535 {
		log.Fatal("Invalid port: ", port)
	}

	fmt.Printf("Starting httpreplay proxy on port %d -> forwarding to %s\n", port, target)
	fmt.Printf("Capturing requests to %s...\n", output)

	file, encoder, err := store.NewWriter(output)
	if err != nil {
		log.Fatal("Failed to create store: ", err)
	}
	defer file.Close()

	proxyHandler := proxy.NewReverseProxy(target, encoder)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, proxyHandler))
	return nil
}
