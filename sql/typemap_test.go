// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"testing"
)

func TestMapSQLType(t *testing.T) {
	tests := []struct {
		driver   DriverName
		sqlType  string
		isPK     bool
		wantType string
		wantLen  int
		wantNil  bool
	}{
		// Integer types
		{DriverPostgres, "integer", false, "Integer", 0, false},
		{DriverPostgres, "integer", true, "AutoNumber", 0, false},
		{DriverPostgres, "int4", false, "Integer", 0, false},
		{DriverPostgres, "smallint", false, "Integer", 0, false},
		{DriverPostgres, "serial", true, "AutoNumber", 0, false},
		{DriverPostgres, "serial", false, "Integer", 0, false},

		// Long
		{DriverPostgres, "bigint", false, "Long", 0, false},
		{DriverPostgres, "bigint", true, "AutoNumber", 0, false},

		// Decimal
		{DriverPostgres, "numeric", false, "Decimal", 0, false},
		{DriverPostgres, "real", false, "Decimal", 0, false},
		{DriverPostgres, "double precision", false, "Decimal", 0, false},

		// Oracle NUMBER
		{DriverOracle, "NUMBER(5)", false, "Integer", 0, false},
		{DriverOracle, "NUMBER(15)", false, "Long", 0, false},
		{DriverOracle, "NUMBER", false, "Decimal", 0, false},

		// String with length
		{DriverPostgres, "varchar(100)", false, "String", 100, false},
		{DriverPostgres, "character varying(200)", false, "String", 200, false},
		{DriverOracle, "VARCHAR2(50)", false, "String", 50, false},
		{DriverSQLServer, "nvarchar(255)", false, "String", 255, false},
		{DriverPostgres, "char(1)", false, "String", 1, false},

		// Unlimited string
		{DriverPostgres, "text", false, "String", 0, false},
		{DriverOracle, "clob", false, "String", 0, false},
		{DriverPostgres, "json", false, "String", 0, false},
		{DriverPostgres, "jsonb", false, "String", 0, false},

		// Boolean
		{DriverPostgres, "boolean", false, "Boolean", 0, false},
		{DriverSQLServer, "bit", false, "Boolean", 0, false},

		// DateTime
		{DriverPostgres, "timestamp", false, "DateTime", 0, false},
		{DriverPostgres, "date", false, "DateTime", 0, false},
		{DriverSQLServer, "datetime2", false, "DateTime", 0, false},
		{DriverOracle, "date", false, "DateTime", 0, false},

		// UUID
		{DriverPostgres, "uuid", false, "String", 36, false},
		{DriverSQLServer, "uniqueidentifier", false, "String", 36, false},

		// Unmappable → nil
		{DriverPostgres, "bytea", false, "", 0, true},
		{DriverOracle, "blob", false, "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.sqlType, func(t *testing.T) {
			got := MapSQLType(tt.driver, tt.sqlType, tt.isPK)
			if tt.wantNil {
				if got != nil {
					t.Errorf("MapSQLType(%q) = %v, want nil", tt.sqlType, got)
				}
				return
			}
			if got == nil {
				t.Fatalf("MapSQLType(%q) = nil, want %s", tt.sqlType, tt.wantType)
			}
			if got.TypeName != tt.wantType {
				t.Errorf("MapSQLType(%q).TypeName = %q, want %q", tt.sqlType, got.TypeName, tt.wantType)
			}
			if got.Length != tt.wantLen {
				t.Errorf("MapSQLType(%q).Length = %d, want %d", tt.sqlType, got.Length, tt.wantLen)
			}
		})
	}
}

func TestGoDriverDSNToJDBC(t *testing.T) {
	tests := []struct {
		driver DriverName
		dsn    string
		want   string
	}{
		{DriverPostgres, "postgres://user:pass@host:5432/mydb", "jdbc:postgresql://host:5432/mydb"},
		{DriverPostgres, "postgres://user@host/mydb", "jdbc:postgresql://host:5432/mydb"},
		{DriverOracle, "oracle://user:pass@host:1521/service", "jdbc:oracle:thin:@//host:1521/service"},
		{DriverSQLServer, "sqlserver://user:pass@host:1433?database=mydb", "jdbc:sqlserver://host:1433;databaseName=mydb"},
	}

	for _, tt := range tests {
		t.Run(string(tt.driver), func(t *testing.T) {
			got, err := GoDriverDSNToJDBC(tt.driver, tt.dsn)
			if err != nil {
				t.Fatalf("GoDriverDSNToJDBC() error: %v", err)
			}
			if got != tt.want {
				t.Errorf("GoDriverDSNToJDBC() = %q, want %q", got, tt.want)
			}
		})
	}
}
