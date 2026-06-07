package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"httpreplay/internal/store"
)

// --- Types & Constants ---

type ctxKey string

const recordDataKey ctxKey = "recordData"

// requestTracker passes timing and payload data from ServeHTTP to ModifyResponse
type requestTracker struct {
	startTime time.Time
	reqBody   []byte
}

// ProxyHandler implements http.Handler to intercept and log traffic
type ProxyHandler struct {
	targetURL *url.URL
	proxy     *httputil.ReverseProxy
	encoder   *json.Encoder
}

// --- Public Constructor ---

// NewReverseProxy builds and configures the reverse proxy pipeline
func NewReverseProxy(target string, encoder *json.Encoder) http.Handler {
	parsedURL, err := url.Parse(target)
	log.Println("Parsed target URL: ", parsedURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		panic("Invalid target URL: " + target)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(parsedURL)

	p := &ProxyHandler{
		targetURL: parsedURL,
		proxy:     reverseProxy,
		encoder:   encoder,
	}

	// Configure outbound request transformations
	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(r *http.Request) {
		originalDirector(r)
		p.configureDirector(r)
	}

	// Configure inbound response interception
	reverseProxy.ModifyResponse = func(resp *http.Response) error {
		return p.handleResponseInterception(resp)
	}

	return p
}

// --- Core Middleware Routing ---

// ServeHTTP handles the live lifecycle of every incoming request
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqBodyBytes := p.extractRequestBody(r)

	tracker := &requestTracker{
		startTime: time.Now(),
		reqBody:   reqBodyBytes,
	}

	// Attach data tracking backpack to request context
	ctx := context.WithValue(r.Context(), recordDataKey, tracker)
	r = r.WithContext(ctx)

	// Forward to the proxy engine
	p.proxy.ServeHTTP(w, r)
}

// --- Isolated Helper Methods ---

// configureDirector rewrites headers to satisfy Cloudflare security requirements
func (p *ProxyHandler) configureDirector(r *http.Request) {
	r.Host = p.targetURL.Host
	r.URL.Host = p.targetURL.Host
	r.URL.Scheme = p.targetURL.Scheme

	// Prevent server-side compression so we capture readable plain text/JSON
	r.Header.Del("Accept-Encoding")
}

// extractRequestBody copies incoming body bytes without draining the browser stream
func (p *ProxyHandler) extractRequestBody(r *http.Request) []byte {
	if r.Body == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body: ", err)
		return nil
	}

	// Re-seal the body stream for standard library routing
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

// handleResponseInterception reads response streams, calculates timing, and writes logs
func (p *ProxyHandler) handleResponseInterception(resp *http.Response) error {
	log.Printf("<<< Intercepted Response Status: %d", resp.StatusCode)

	// 1. Unpack context tracker data
	var duration time.Duration
	var reqBodyBytes []byte
	if tracker, ok := resp.Request.Context().Value(recordDataKey).(*requestTracker); ok {
		duration = time.Since(tracker.startTime)
		reqBodyBytes = tracker.reqBody
	}

	// 2. Clone response body bytes safely
	var resBodyBytes []byte
	if resp.Body != nil {
		var err error
		resBodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read response body: ", err)
		}
		// Re-seal the body stream for delivery to user browser
		resp.Body = io.NopCloser(bytes.NewBuffer(resBodyBytes))
	}

	// 3. Assemble and persist the transaction log line
	record := store.Record{
		Timestamp:  time.Now().UTC(),
		Method:     resp.Request.Method,
		URL:        resp.Request.URL.String(),
		ReqHeaders: resp.Request.Header,
		ReqBody:    reqBodyBytes,
		StatusCode: resp.StatusCode,
		ResHeaders: resp.Header,
		ResBody:    resBodyBytes,
		Duration:   duration,
	}

	if err := p.encoder.Encode(record); err != nil {
		log.Println("Failed to encode record: ", err)
	}

	return nil
}
