package replay

import (
	"bytes"
	"fmt"
	"httpreplay/internal/store"
	"log"
	"net/http"
	"net/url"
)

type Replayer struct {
	Target *url.URL
	Client *http.Client
}

func (rp *Replayer) TranslateURL(originalURL string) (string, error) {
	if rp.Target == nil {
		return originalURL, nil
	}

	origURL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	origURL.Host = rp.Target.Host
	origURL.Scheme = rp.Target.Scheme
	return origURL.String(), nil
}

func NewReplayer(newTarget *url.URL) *Replayer {
	return &Replayer{
		Target: newTarget,
		Client: &http.Client{},
	}
}

func (rp *Replayer) Playback(records []store.Record) error {
	for _, record := range records {
		targetURL, err := rp.TranslateURL(record.URL)
		if err != nil {
			log.Println("Failed to translate URL: ", err)
			continue
		}

		req, err := http.NewRequest(record.Method, targetURL, bytes.NewReader(record.ReqBody))
		if err != nil {
			log.Println("Failed to create request: ", err)
			continue
		}

		for key, values := range record.ReqHeaders {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		log.Println("Sending request: ", req.URL.String())

		resp, err := rp.Client.Do(req)
		if err != nil {
			log.Println("Failed network transmission: ", err)
			continue
		}
		resp.Body.Close()
		// Print your victory lap!
		fmt.Printf("[%s] Replayed %s -> New Status: %d (Recorded was: %d)\n",
			record.Method, req.URL.Path, resp.StatusCode, record.StatusCode)
	}
	return nil
}
