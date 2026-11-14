// Command http-client drives the Web API using only the Go standard library.
// It covers the full request chain used by the Web UI:
// generate -> status -> result -> names -> explore.
//
//	go run ./cmd/server -addr :18080     # in another terminal
//	go run ./examples/http-client
//
// Override the server address with FATE_BASE_URL.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type response map[string]any

func call(client *http.Client, method, url string, payload any) (int, response) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("marshal request: %v", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalf("build request: %v", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("%s %s: %v (is the server running?)", method, req.URL.Path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read response: %v", err)
	}

	var parsed response
	if err := json.Unmarshal(raw, &parsed); err != nil {
		log.Fatalf("%s %s: invalid JSON (status %d): %s", method, req.URL.Path, resp.StatusCode, raw)
	}
	fmt.Printf("%s %s -> %d (%d bytes)\n", method, req.URL.Path, resp.StatusCode, len(raw))
	return resp.StatusCode, parsed
}

func mustStatus(step string, got, want int) {
	if got != want {
		log.Fatalf("%s: status %d, want %d", step, got, want)
	}
}

func main() {
	base := os.Getenv("FATE_BASE_URL")
	if base == "" {
		base = "http://127.0.0.1:18080"
	}
	// /api/v1/generate is synchronous and may take a while on a cold cache.
	client := &http.Client{Timeout: 5 * time.Minute}

	status, generated := call(client, http.MethodPost, base+"/api/v1/generate", map[string]any{
		"surname":        "张",
		"born":           "2024/06/15 10:30",
		"sex":            "boy",
		"poetry_mode":    0,
		"xi_yong_method": "balance",
		// The API takes comma-separated strings, not JSON arrays.
		"avoid_chars":   "病,死",
		"require_chars": "",
		"filter_type":   0,
	})
	mustStatus("generate", status, http.StatusOK)

	taskID, _ := generated["task_id"].(string)
	if taskID == "" {
		log.Fatalf("generate response has no task_id: %v", generated)
	}
	fmt.Printf("task_id=%s total=%v state=%v\n", taskID, generated["total"], generated["state"])

	status, state := call(client, http.MethodGet, base+"/api/v1/status?task_id="+taskID, nil)
	mustStatus("status", status, http.StatusOK)
	fmt.Printf("state=%v\n", state["state"])

	status, result := call(client, http.MethodGet, base+"/api/v1/result?task_id="+taskID, nil)
	mustStatus("result", status, http.StatusOK)
	topNames, ok := result["top_names"].([]any)
	if !ok {
		log.Fatalf("top_names must be an array, got %T", result["top_names"])
	}
	fmt.Printf("result: top_names=%d total=%v\n", len(topNames), result["total"])

	status, paged := call(client, http.MethodGet,
		base+"/api/v1/names?task_id="+taskID+"&page=1&size=5", nil)
	mustStatus("names", status, http.StatusOK)
	pageNames, ok := paged["names"].([]any)
	if !ok {
		log.Fatalf("names must be an array, got %T", paged["names"])
	}
	fmt.Printf("names: page=%v size=%v got=%d total=%v\n",
		paged["page"], paged["size"], len(pageNames), paged["total"])

	status, explored := call(client, http.MethodGet,
		base+"/api/v1/explore?task_id="+taskID+"&count=3", nil)
	mustStatus("explore", status, http.StatusOK)
	// Regression point from v4.3.0: an exhausted pool must be [] and not null.
	exploredNames, ok := explored["names"].([]any)
	if !ok {
		log.Fatalf("explore names must be a JSON array, got %T (%v)",
			explored["names"], explored["names"])
	}
	fmt.Printf("explore: names is a JSON array of %d entries, shown_count=%v\n",
		len(exploredNames), explored["shown_count"])

	fmt.Println("http-client complete: 5 endpoints OK")
}
