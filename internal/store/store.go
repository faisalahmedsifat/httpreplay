package store

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"bufio"
)

type Record struct {
	Timestamp  time.Time     `json:"timestamp"`
	Method     string        `json:"method"`
	URL        string        `json:"url"`
	ReqHeaders http.Header   `json:"req_headers"`
	ReqBody    []byte        `json:"req_body"`
	StatusCode int           `json:"status_code"` // 200, 400, 500, etc.
	ResHeaders http.Header   `json:"res_headers"`
	ResBody    []byte        `json:"res_body"`
	Duration   time.Duration `json:"duration"`
}

func NewWriter(path string) (*os.File, *json.Encoder, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}
	return f, json.NewEncoder(f), nil
}

func ReadAll(path string) ([]Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var records []Record

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		var r Record

		err := json.Unmarshal(lineBytes, &r)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
