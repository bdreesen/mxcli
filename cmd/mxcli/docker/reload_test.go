// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReload_CSSOnly(t *testing.T) {
	var receivedAction string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		receivedAction = body["action"].(string)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":0,"feedback":{}}`))
	}))
	defer server.Close()

	host, port := parseTestServerAddr(t, server.URL)

	var stdout bytes.Buffer
	opts := ReloadOptions{
		CSSOnly: true,
		Host:    host,
		Port:    port,
		Token:   "testpass",
		Direct:  true,
		Stdout:  &stdout,
	}

	err := Reload(opts)
	if err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if receivedAction != "update_styling" {
		t.Errorf("expected update_styling action, got %q", receivedAction)
	}
	if !strings.Contains(stdout.String(), "Styling updated") {
		t.Errorf("expected 'Styling updated' in output, got: %s", stdout.String())
	}
}

func TestReload_ModelOnly(t *testing.T) {
	var receivedAction string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		receivedAction = body["action"].(string)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":0,"feedback":{"startup_metrics":{"duration":98}}}`))
	}))
	defer server.Close()

	host, port := parseTestServerAddr(t, server.URL)

	var stdout bytes.Buffer
	opts := ReloadOptions{
		SkipBuild: true,
		Host:      host,
		Port:      port,
		Token:     "testpass",
		Direct:    true,
		Stdout:    &stdout,
	}

	err := Reload(opts)
	if err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if receivedAction != "reload_model" {
		t.Errorf("expected reload_model action, got %q", receivedAction)
	}
	if !strings.Contains(stdout.String(), "Model reloaded") {
		t.Errorf("expected 'Model reloaded' in output, got: %s", stdout.String())
	}
}

func TestReload_ParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		feedback map[string]any
		want     string
	}{
		{
			name:     "nil feedback",
			feedback: nil,
			want:     "",
		},
		{
			name:     "no startup_metrics",
			feedback: map[string]any{},
			want:     "",
		},
		{
			name: "duration in ms",
			feedback: map[string]any{
				"startup_metrics": map[string]any{
					"duration": float64(98),
				},
			},
			want: "98ms",
		},
		{
			name: "duration in seconds",
			feedback: map[string]any{
				"startup_metrics": map[string]any{
					"duration": float64(2500),
				},
			},
			want: "2.5s",
		},
		{
			name: "duration as string",
			feedback: map[string]any{
				"startup_metrics": map[string]any{
					"duration": "150ms",
				},
			},
			want: "150ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractReloadDuration(tt.feedback)
			if got != tt.want {
				t.Errorf("extractReloadDuration: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReload_ModelOnly_WithDuration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":0,"feedback":{"startup_metrics":{"duration":98}}}`))
	}))
	defer server.Close()

	host, port := parseTestServerAddr(t, server.URL)

	var stdout bytes.Buffer
	opts := ReloadOptions{
		SkipBuild: true,
		Host:      host,
		Port:      port,
		Token:     "testpass",
		Direct:    true,
		Stdout:    &stdout,
	}

	err := Reload(opts)
	if err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if !strings.Contains(stdout.String(), "98ms") {
		t.Errorf("expected duration in output, got: %s", stdout.String())
	}
}

func TestReload_CSSOnly_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"cause":"runtime not running"}`))
	}))
	defer server.Close()

	host, port := parseTestServerAddr(t, server.URL)

	var stdout bytes.Buffer
	opts := ReloadOptions{
		CSSOnly: true,
		Host:    host,
		Port:    port,
		Token:   "testpass",
		Direct:  true,
		Stdout:  &stdout,
	}

	err := Reload(opts)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "runtime not running") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReload_ModelOnly_ReloadError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"cause":"model contains errors"}`))
	}))
	defer server.Close()

	host, port := parseTestServerAddr(t, server.URL)

	var stdout bytes.Buffer
	opts := ReloadOptions{
		SkipBuild: true,
		Host:      host,
		Port:      port,
		Token:     "testpass",
		Direct:    true,
		Stdout:    &stdout,
	}

	err := Reload(opts)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "model contains errors") {
		t.Errorf("unexpected error: %v", err)
	}
}
