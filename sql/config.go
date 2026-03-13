// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ResolveOptions controls how DSN resolution works.
type ResolveOptions struct {
	// Explicit DSN provided directly (highest priority).
	DSN string
	// Alias to look up in env vars and config file.
	Alias string
	// Driver name to look up in env vars.
	Driver DriverName
	// ConfigDir is the directory to search for .mxcli/connections.yaml.
	// If empty, the current working directory is used.
	ConfigDir string
}

// ResolvedConnection holds the result of DSN resolution, including the driver.
type ResolvedConnection struct {
	DSN    string
	Driver DriverName
}

// connectionsConfig represents the .mxcli/connections.yaml file.
// Supports two formats:
//
//	# List format:
//	connections:
//	  - alias: mydb
//	    driver: postgres
//	    dsn: postgres://user:pass@localhost:5432/db
//
//	# Map format:
//	connections:
//	  mydb:
//	    driver: postgres
//	    dsn: postgres://user:pass@localhost:5432/db
type connectionsConfig struct {
	Connections yaml.Node `yaml:"connections"`
}

type connectionEntry struct {
	Alias  string `yaml:"alias"`
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

// ResolveDSN resolves a DSN using the following priority:
//  1. Explicit DSN (opts.DSN)
//  2. Environment variable MXCLI_SQL_<ALIAS>_DSN
//  3. Environment variable MXCLI_SQL_<DRIVER>_DSN
//  4. .mxcli/connections.yaml file
func ResolveDSN(opts ResolveOptions) (string, error) {
	rc, err := ResolveConnection(opts)
	if err != nil {
		return "", err
	}
	return rc.DSN, nil
}

// ResolveConnection resolves both DSN and driver using the standard priority order.
// The driver is determined from: opts.Driver (if set), config file driver field,
// or auto-detected from the DSN scheme.
func ResolveConnection(opts ResolveOptions) (*ResolvedConnection, error) {
	result := &ResolvedConnection{Driver: opts.Driver}

	// 1. Explicit DSN
	if opts.DSN != "" {
		result.DSN = opts.DSN
		if result.Driver == "" {
			result.Driver = DetectDriverFromDSN(opts.DSN)
		}
		return result, nil
	}

	// 2. Env var by alias: MXCLI_SQL_<ALIAS>_DSN
	if opts.Alias != "" {
		envKey := "MXCLI_SQL_" + strings.ToUpper(opts.Alias) + "_DSN"
		if dsn := os.Getenv(envKey); dsn != "" {
			result.DSN = dsn
			if result.Driver == "" {
				result.Driver = DetectDriverFromDSN(dsn)
			}
			return result, nil
		}
	}

	// 3. Env var by driver: MXCLI_SQL_<DRIVER>_DSN
	if opts.Driver != "" {
		envKey := "MXCLI_SQL_" + strings.ToUpper(string(opts.Driver)) + "_DSN"
		if dsn := os.Getenv(envKey); dsn != "" {
			result.DSN = dsn
			return result, nil
		}
	}

	// 4. Config file: .mxcli/connections.yaml
	if opts.Alias != "" {
		configDir := opts.ConfigDir
		if configDir == "" {
			configDir, _ = os.Getwd()
		}
		if configDir != "" {
			entry, err := loadEntryFromConfig(filepath.Join(configDir, ".mxcli", "connections.yaml"), opts.Alias)
			if err == nil && entry.DSN != "" {
				result.DSN = entry.DSN
				if result.Driver == "" && entry.Driver != "" {
					d, _ := ParseDriver(entry.Driver)
					result.Driver = d
				}
				if result.Driver == "" {
					result.Driver = DetectDriverFromDSN(entry.DSN)
				}
				return result, nil
			}
		}
	}

	return nil, fmt.Errorf("could not resolve DSN: set --dsn flag, MXCLI_SQL_%s_DSN env var, or add to .mxcli/connections.yaml",
		strings.ToUpper(opts.Alias))
}

// DetectDriverFromDSN infers the database driver from the DSN scheme.
// Returns empty string if the scheme is not recognized.
func DetectDriverFromDSN(dsn string) DriverName {
	lower := strings.ToLower(dsn)
	switch {
	case strings.HasPrefix(lower, "postgres://") || strings.HasPrefix(lower, "postgresql://"):
		return DriverPostgres
	case strings.HasPrefix(lower, "oracle://"):
		return DriverOracle
	case strings.HasPrefix(lower, "sqlserver://"):
		return DriverSQLServer
	default:
		return ""
	}
}

// loadEntryFromConfig reads a connection entry from a connections.yaml file.
// Supports both list format and map format.
func loadEntryFromConfig(path, alias string) (*connectionEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg connectionsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	if cfg.Connections.Kind == 0 {
		return nil, fmt.Errorf("no 'connections' key in %s", path)
	}

	switch cfg.Connections.Kind {
	case yaml.SequenceNode:
		// List format: connections: [{alias: x, driver: y, dsn: z}, ...]
		return loadFromList(&cfg.Connections, alias, path)
	case yaml.MappingNode:
		// Map format: connections: {alias: {driver: y, dsn: z}, ...}
		return loadFromMap(&cfg.Connections, alias, path)
	default:
		return nil, fmt.Errorf("'connections' in %s must be a list or map", path)
	}
}

func loadFromList(node *yaml.Node, alias, path string) (*connectionEntry, error) {
	for _, item := range node.Content {
		var entry connectionEntry
		if err := item.Decode(&entry); err != nil {
			continue
		}
		if strings.EqualFold(entry.Alias, alias) {
			return &entry, nil
		}
	}
	return nil, fmt.Errorf("alias %q not found in %s", alias, path)
}

func loadFromMap(node *yaml.Node, alias, path string) (*connectionEntry, error) {
	// MappingNode has alternating key/value content nodes
	for i := 0; i+1 < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		if strings.EqualFold(keyNode.Value, alias) {
			var entry connectionEntry
			if err := valNode.Decode(&entry); err != nil {
				return nil, fmt.Errorf("failed to parse entry for %q in %s: %w", alias, path, err)
			}
			entry.Alias = keyNode.Value
			return &entry, nil
		}
	}
	return nil, fmt.Errorf("alias %q not found in %s", alias, path)
}
