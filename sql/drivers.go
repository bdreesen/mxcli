// SPDX-License-Identifier: Apache-2.0

// Package sql provides database connectivity for mxcli.
//
// Import this package for side-effect driver registration.
package sql

import (
	_ "github.com/jackc/pgx/v5/stdlib"  // registers "pgx" driver
	_ "github.com/microsoft/go-mssqldb" // registers "sqlserver" driver
	_ "github.com/sijms/go-ora/v2"      // registers "oracle" driver
)
