// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// MendixType represents a Mendix attribute type mapped from a SQL column.
type MendixType struct {
	TypeName string // "String", "Integer", "Decimal", "Boolean", "DateTime", "Long", "AutoNumber"
	Length   int    // for String(n); 0 = unlimited
}

// String returns the MDL representation, e.g. "String(100)" or "Integer".
func (t *MendixType) String() string {
	if t.TypeName == "String" && t.Length > 0 {
		return fmt.Sprintf("String(%d)", t.Length)
	}
	return t.TypeName
}

// reParenNum matches a type name with a parenthesized number, e.g. "varchar(100)".
var reParenNum = regexp.MustCompile(`\((\d+)`)

// MapSQLType maps a SQL data_type string to a Mendix type.
// Returns nil if the type cannot be mapped.
func MapSQLType(driver DriverName, sqlType string, isPK bool) *MendixType {
	lower := strings.ToLower(strings.TrimSpace(sqlType))

	// Extract length from parenthesized types like varchar(100)
	length := 0
	if m := reParenNum.FindStringSubmatch(lower); len(m) > 1 {
		length, _ = strconv.Atoi(m[1])
	}

	// Strip parenthesized portion for matching
	base := lower
	if idx := strings.Index(base, "("); idx >= 0 {
		base = strings.TrimSpace(base[:idx])
	}

	switch base {
	// Integer types
	case "integer", "int", "int4", "int2", "smallint", "tinyint", "mediumint":
		if isPK {
			return &MendixType{TypeName: "AutoNumber"}
		}
		return &MendixType{TypeName: "Integer"}
	case "serial", "smallserial":
		if isPK {
			return &MendixType{TypeName: "AutoNumber"}
		}
		return &MendixType{TypeName: "Integer"}

	// Long types
	case "bigint", "int8", "bigserial":
		if isPK {
			return &MendixType{TypeName: "AutoNumber"}
		}
		return &MendixType{TypeName: "Long"}

	// Oracle NUMBER: NUMBER(p,0) with p <= 10 → Integer, else Long
	case "number":
		if length > 0 && length <= 10 {
			if isPK {
				return &MendixType{TypeName: "AutoNumber"}
			}
			return &MendixType{TypeName: "Integer"}
		}
		if length > 10 {
			if isPK {
				return &MendixType{TypeName: "AutoNumber"}
			}
			return &MendixType{TypeName: "Long"}
		}
		return &MendixType{TypeName: "Decimal"}

	// Decimal/float types
	case "numeric", "decimal", "money", "smallmoney",
		"real", "float", "float4", "float8", "double", "double precision",
		"binary_float", "binary_double":
		return &MendixType{TypeName: "Decimal"}

	// String types with length
	case "varchar", "character varying", "varchar2", "nvarchar", "nvarchar2":
		if length > 0 {
			return &MendixType{TypeName: "String", Length: length}
		}
		return &MendixType{TypeName: "String"}
	case "char", "character", "nchar":
		if length > 0 {
			return &MendixType{TypeName: "String", Length: length}
		}
		return &MendixType{TypeName: "String", Length: 1}

	// Unlimited string types
	case "text", "clob", "nclob", "ntext", "xml", "json", "jsonb", "long":
		return &MendixType{TypeName: "String"}

	// Boolean
	case "boolean", "bool", "bit":
		return &MendixType{TypeName: "Boolean"}

	// DateTime types
	case "date", "timestamp", "timestamptz",
		"timestamp without time zone", "timestamp with time zone",
		"datetime", "datetime2", "smalldatetime", "datetimeoffset",
		"time", "interval":
		return &MendixType{TypeName: "DateTime"}

	// UUID → String(36)
	case "uuid", "uniqueidentifier":
		return &MendixType{TypeName: "String", Length: 36}

	// Identity column (SQL Server)
	case "identity":
		return &MendixType{TypeName: "AutoNumber"}

	default:
		// Check for "max" length variants like nvarchar(max)
		if strings.Contains(lower, "max") &&
			(strings.HasPrefix(base, "varchar") || strings.HasPrefix(base, "nvarchar") || strings.HasPrefix(base, "varbinary")) {
			if strings.HasPrefix(base, "varbinary") {
				return nil // binary → skip
			}
			return &MendixType{TypeName: "String"}
		}
		return nil
	}
}

// GoDriverDSNToJDBC converts a Go driver DSN to a JDBC URL.
func GoDriverDSNToJDBC(driver DriverName, dsn string) (string, error) {
	switch driver {
	case DriverPostgres:
		u, err := url.Parse(dsn)
		if err != nil {
			return "", fmt.Errorf("failed to parse postgres DSN: %w", err)
		}
		host := u.Hostname()
		port := u.Port()
		if port == "" {
			port = "5432"
		}
		db := strings.TrimPrefix(u.Path, "/")
		return fmt.Sprintf("jdbc:postgresql://%s:%s/%s", host, port, db), nil

	case DriverOracle:
		u, err := url.Parse(dsn)
		if err != nil {
			return "", fmt.Errorf("failed to parse oracle DSN: %w", err)
		}
		host := u.Hostname()
		port := u.Port()
		if port == "" {
			port = "1521"
		}
		service := strings.TrimPrefix(u.Path, "/")
		return fmt.Sprintf("jdbc:oracle:thin:@//%s:%s/%s", host, port, service), nil

	case DriverSQLServer:
		u, err := url.Parse(dsn)
		if err != nil {
			return "", fmt.Errorf("failed to parse sqlserver DSN: %w", err)
		}
		host := u.Hostname()
		port := u.Port()
		if port == "" {
			port = "1433"
		}
		db := u.Query().Get("database")
		if db == "" {
			db = strings.TrimPrefix(u.Path, "/")
		}
		jdbc := fmt.Sprintf("jdbc:sqlserver://%s:%s", host, port)
		if db != "" {
			jdbc += ";databaseName=" + db
		}
		return jdbc, nil

	default:
		return dsn, nil
	}
}
