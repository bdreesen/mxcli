// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// MendixDBAlias is the reserved connection alias for the Mendix app database.
// The underscore prefix avoids collisions with user-defined aliases.
const MendixDBAlias = "_mendix"

// BuildMendixDSN constructs a PostgreSQL DSN from Mendix project settings.
//
// Environment variable overrides (useful for devcontainers/Docker):
//
//	MXCLI_DB_TYPE     — overrides DatabaseType (e.g., "PostgreSql")
//	MXCLI_DB_HOST     — overrides host part of DatabaseUrl
//	MXCLI_DB_PORT     — overrides port part of DatabaseUrl
//	MXCLI_DB_NAME     — overrides DatabaseName
//	MXCLI_DB_USER     — overrides DatabaseUserName
//	MXCLI_DB_PASSWORD — overrides DatabasePassword
func BuildMendixDSN(dbType, dbURL, dbName, dbUser, dbPassword string) (string, error) {
	// Apply env var overrides
	if v := os.Getenv("MXCLI_DB_TYPE"); v != "" {
		dbType = v
	}
	if v := os.Getenv("MXCLI_DB_NAME"); v != "" {
		dbName = v
	}
	if v := os.Getenv("MXCLI_DB_USER"); v != "" {
		dbUser = v
	}
	if v := os.Getenv("MXCLI_DB_PASSWORD"); v != "" {
		dbPassword = v
	}

	if !strings.EqualFold(dbType, "POSTGRESQL") && !strings.EqualFold(dbType, "PostgreSql") {
		return "", fmt.Errorf("only PostgreSQL Mendix databases are supported (got %q); set MXCLI_DB_TYPE=PostgreSql to override", dbType)
	}

	if dbURL == "" {
		dbURL = "localhost:5432"
	}

	// dbURL is typically "localhost:5432" or "host:port"
	host := dbURL

	// Apply host/port overrides
	envHost := os.Getenv("MXCLI_DB_HOST")
	envPort := os.Getenv("MXCLI_DB_PORT")

	if envHost != "" || envPort != "" {
		// Parse existing host:port
		origHost := host
		origPort := ""
		if idx := strings.LastIndex(host, ":"); idx >= 0 {
			origHost = host[:idx]
			origPort = host[idx+1:]
		}
		if envHost != "" {
			origHost = envHost
		}
		if envPort != "" {
			origPort = envPort
		}
		if origPort != "" {
			host = origHost + ":" + origPort
		} else {
			host = origHost
		}
	}

	// Build postgres:// DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		url.PathEscape(dbUser),
		url.PathEscape(dbPassword),
		host,
		url.PathEscape(dbName),
	)

	return dsn, nil
}

// EntityToTableName converts a Mendix qualified name to the PostgreSQL table name.
// "MyModule.Customer" → "mymodule$customer"
func EntityToTableName(qualifiedName string) (string, error) {
	parts := strings.SplitN(qualifiedName, ".", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid entity name %q (expected Module.Entity)", qualifiedName)
	}
	return strings.ToLower(parts[0]) + "$" + strings.ToLower(parts[1]), nil
}

// AttributeToColumnName converts a Mendix attribute name to the PostgreSQL column name.
func AttributeToColumnName(attrName string) string {
	return strings.ToLower(attrName)
}
