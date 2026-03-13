// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"fmt"
	"strings"
)

// DriverName identifies a supported database driver.
type DriverName string

const (
	DriverPostgres  DriverName = "postgres"
	DriverOracle    DriverName = "oracle"
	DriverSQLServer DriverName = "sqlserver"
)

// driverGoName maps our driver names to Go sql.Open driver names.
var driverGoName = map[DriverName]string{
	DriverPostgres:  "pgx",
	DriverOracle:    "oracle",
	DriverSQLServer: "sqlserver",
}

// GoDriverName returns the Go database/sql driver name for this driver.
func (d DriverName) GoDriverName() string {
	return driverGoName[d]
}

// String returns the driver name as a string.
func (d DriverName) String() string {
	return string(d)
}

// ParseDriver parses a driver name string (case-insensitive).
func ParseDriver(s string) (DriverName, error) {
	switch strings.ToLower(s) {
	case "postgres", "postgresql", "pg":
		return DriverPostgres, nil
	case "oracle", "ora":
		return DriverOracle, nil
	case "sqlserver", "mssql":
		return DriverSQLServer, nil
	default:
		return "", fmt.Errorf("unsupported driver: %q (supported: postgres, oracle, sqlserver)", s)
	}
}
