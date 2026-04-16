// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFetchODataMetadata_LocalFile(t *testing.T) {
	// Create a temporary metadata file
	tmpDir := t.TempDir()
	metadataContent := `<?xml version="1.0"?><edmx:Edmx xmlns:edmx="http://docs.oasis-open.org/odata/ns/edmx" Version="4.0"><edmx:DataServices><Schema xmlns="http://docs.oasis-open.org/odata/ns/edm" Namespace="Test"><EntityType Name="Product"><Key><PropertyRef Name="ID"/></Key><Property Name="ID" Type="Edm.Int32"/></EntityType></Schema></edmx:DataServices></edmx:Edmx>`
	metadataPath := filepath.Join(tmpDir, "metadata.xml")
	if err := os.WriteFile(metadataPath, []byte(metadataContent), 0644); err != nil {
		t.Fatalf("Failed to create test metadata file: %v", err)
	}

	tests := []struct {
		name        string
		url         string
		mprDir      string
		wantErr     bool
		errContains string
	}{
		{
			name:    "absolute file:// URL",
			url:     "file://" + metadataPath,
			mprDir:  "",
			wantErr: false,
		},
		{
			name:    "absolute path without file://",
			url:     metadataPath,
			mprDir:  "",
			wantErr: false,
		},
		{
			name:    "relative path with mprDir",
			url:     "metadata.xml",
			mprDir:  tmpDir,
			wantErr: false,
		},
		{
			name:    "relative path with ./ prefix",
			url:     "./metadata.xml",
			mprDir:  tmpDir,
			wantErr: false,
		},
		{
			name:        "nonexistent file",
			url:         "file:///nonexistent/metadata.xml",
			mprDir:      "",
			wantErr:     true,
			errContains: "failed to read local metadata file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, hash, err := fetchODataMetadata(tt.url, tt.mprDir)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if metadata != metadataContent {
				t.Errorf("Metadata content mismatch.\nGot: %q\nWant: %q", metadata, metadataContent)
			}

			if hash == "" {
				t.Errorf("Expected non-empty hash")
			}

			// Hash should be consistent
			_, hash2, _ := fetchODataMetadata(tt.url, tt.mprDir)
			if hash != hash2 {
				t.Errorf("Hash inconsistent between calls: %q vs %q", hash, hash2)
			}
		})
	}
}

func TestNormalizeMetadataUrl(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		input      string
		mprDir     string
		wantPrefix string
		wantErr    bool
	}{
		{
			name:       "HTTP URL unchanged",
			input:      "https://api.example.com/$metadata",
			mprDir:     "",
			wantPrefix: "https://",
			wantErr:    false,
		},
		{
			name:       "HTTPS URL unchanged",
			input:      "http://localhost:8080/odata/$metadata",
			mprDir:     "",
			wantPrefix: "http://",
			wantErr:    false,
		},
		{
			name:       "Absolute file:// unchanged",
			input:      "file:///tmp/metadata.xml",
			mprDir:     "",
			wantPrefix: "file://",
			wantErr:    false,
		},
		{
			name:       "Relative path normalized to file://",
			input:      "./metadata.xml",
			mprDir:     tmpDir,
			wantPrefix: "file://",
			wantErr:    false,
		},
		{
			name:       "Bare relative path normalized to file://",
			input:      "metadata.xml",
			mprDir:     tmpDir,
			wantPrefix: "file://",
			wantErr:    false,
		},
		{
			name:       "Absolute path normalized to file://",
			input:      "/tmp/metadata.xml",
			mprDir:     "",
			wantPrefix: "file://",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeMetadataUrl(tt.input, tt.mprDir)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !strings.HasPrefix(result, tt.wantPrefix) {
				t.Errorf("Result %q does not start with %q", result, tt.wantPrefix)
			}

			// Verify file:// URLs are absolute
			if strings.HasPrefix(result, "file://") {
				path := strings.TrimPrefix(result, "file://")
				if !filepath.IsAbs(path) {
					t.Errorf("file:// URL contains relative path: %q", result)
				}
			}
		})
	}
}

func TestFetchODataMetadata_RelativePathWithoutProject(t *testing.T) {
	// Create metadata file in current directory
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	metadataContent := `<?xml version="1.0"?><edmx:Edmx xmlns:edmx="http://docs.oasis-open.org/odata/ns/edmx" Version="4.0"></edmx:Edmx>`
	if err := os.WriteFile("local.xml", []byte(metadataContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with empty mprDir (should resolve against cwd)
	metadata, hash, err := fetchODataMetadata("local.xml", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if metadata != metadataContent {
		t.Errorf("Metadata content mismatch")
	}
	if hash == "" {
		t.Errorf("Expected non-empty hash")
	}
}
