// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"fmt"
	"strings"
)

// tableQuery maps driver to the query that lists user tables.
var tableQuery = map[DriverName]string{
	DriverPostgres: `SELECT table_schema, table_name, table_type
FROM information_schema.tables
WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
ORDER BY table_schema, table_name`,
	DriverOracle: `SELECT owner AS table_schema, table_name, 'BASE TABLE' AS table_type
FROM all_tables
WHERE owner NOT IN ('SYS','SYSTEM','XDB','CTXSYS','MDSYS','OLAPSYS','WMSYS','DBSNMP','APPQOSSYS','DBSFWUSER','OUTLN','LBACSYS','ORDDATA','ORDSYS','GSMADMIN_INTERNAL')
ORDER BY owner, table_name`,
	DriverSQLServer: `SELECT TABLE_SCHEMA AS table_schema, TABLE_NAME AS table_name, TABLE_TYPE AS table_type
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_TYPE = 'BASE TABLE'
ORDER BY TABLE_SCHEMA, TABLE_NAME`,
}

// describeQuery maps driver to the query that describes a table's columns.
// The placeholder %s will be replaced with the table name.
var describeQuery = map[DriverName]string{
	DriverPostgres: `SELECT column_name, data_type,
  CASE WHEN is_nullable = 'YES' THEN 'YES' ELSE 'NO' END AS nullable,
  column_default
FROM information_schema.columns
WHERE table_name = '%s'
ORDER BY ordinal_position`,
	DriverOracle: `SELECT column_name, data_type, data_length, nullable
FROM all_tab_columns
WHERE %s
ORDER BY column_id`,
	DriverSQLServer: `SELECT COLUMN_NAME AS column_name, DATA_TYPE AS data_type,
  IS_NULLABLE AS nullable,
  COLUMN_DEFAULT AS column_default
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_NAME = '%s'
ORDER BY ORDINAL_POSITION`,
}

// viewQuery maps driver to the query that lists user views.
var viewQuery = map[DriverName]string{
	DriverPostgres: `SELECT schemaname AS view_schema, viewname AS view_name, definition
FROM pg_views
WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
ORDER BY schemaname, viewname`,
	DriverOracle: `SELECT owner AS view_schema, view_name, text AS definition
FROM all_views
WHERE owner NOT IN ('SYS','SYSTEM','XDB','CTXSYS','MDSYS','OLAPSYS','WMSYS','DBSNMP','APPQOSSYS','DBSFWUSER','OUTLN','LBACSYS','ORDDATA','ORDSYS','GSMADMIN_INTERNAL')
ORDER BY owner, view_name`,
	DriverSQLServer: `SELECT s.name AS view_schema, v.name AS view_name, m.definition
FROM sys.views v
JOIN sys.schemas s ON s.schema_id = v.schema_id
LEFT JOIN sys.sql_modules m ON m.object_id = v.object_id
ORDER BY s.name, v.name`,
}

// functionQuery maps driver to the query that lists user functions and procedures.
var functionQuery = map[DriverName]string{
	DriverPostgres: `SELECT n.nspname AS schema, p.proname AS name,
  CASE p.prokind WHEN 'f' THEN 'function' WHEN 'p' THEN 'procedure' WHEN 'a' THEN 'aggregate' WHEN 'w' THEN 'window' END AS kind,
  pg_get_function_result(p.oid) AS return_type,
  pg_get_function_arguments(p.oid) AS arguments
FROM pg_proc p
JOIN pg_namespace n ON n.oid = p.pronamespace
WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
ORDER BY n.nspname, p.proname`,
	DriverOracle: `SELECT owner AS schema, object_name AS name, object_type AS kind
FROM all_procedures
WHERE owner NOT IN ('SYS','SYSTEM','XDB','CTXSYS','MDSYS','OLAPSYS','WMSYS','DBSNMP','APPQOSSYS','DBSFWUSER','OUTLN','LBACSYS','ORDDATA','ORDSYS','GSMADMIN_INTERNAL')
  AND object_type IN ('FUNCTION','PROCEDURE','PACKAGE')
ORDER BY owner, object_type, object_name`,
	DriverSQLServer: `SELECT s.name AS [schema], o.name, o.type_desc AS kind,
  m.definition
FROM sys.objects o
JOIN sys.schemas s ON s.schema_id = o.schema_id
LEFT JOIN sys.sql_modules m ON m.object_id = o.object_id
WHERE o.type IN ('FN','IF','TF','P','AF')
  AND o.is_ms_shipped = 0
ORDER BY s.name, o.type_desc, o.name`,
}

// ShowTables returns a list of tables for the given connection.
func ShowTables(ctx context.Context, conn *Connection) (*QueryResult, error) {
	query, ok := tableQuery[conn.Driver]
	if !ok {
		return nil, fmt.Errorf("SHOW TABLES not supported for driver %s", conn.Driver)
	}
	return Execute(ctx, conn, query)
}

// DescribeTable returns the column definitions for a table or view.
func DescribeTable(ctx context.Context, conn *Connection, table string) (*QueryResult, error) {
	tmpl, ok := describeQuery[conn.Driver]
	if !ok {
		return nil, fmt.Errorf("DESCRIBE TABLE not supported for driver %s", conn.Driver)
	}

	var query string
	if conn.Driver == DriverOracle {
		// Oracle: split SCHEMA.TABLE into owner + table_name filter
		query = fmt.Sprintf(tmpl, oracleTableFilter(table))
	} else {
		query = fmt.Sprintf(tmpl, table)
	}
	return Execute(ctx, conn, query)
}

// oracleTableFilter builds a WHERE clause for Oracle's all_tab_columns.
// Supports both "TABLE" and "SCHEMA.TABLE" formats.
func oracleTableFilter(table string) string {
	if i := strings.IndexByte(table, '.'); i >= 0 {
		owner := strings.ToUpper(table[:i])
		name := strings.ToUpper(table[i+1:])
		return fmt.Sprintf("owner = '%s' AND table_name = '%s'", owner, name)
	}
	return fmt.Sprintf("table_name = UPPER('%s')", table)
}

// ShowViews returns a list of views for the given connection.
func ShowViews(ctx context.Context, conn *Connection) (*QueryResult, error) {
	query, ok := viewQuery[conn.Driver]
	if !ok {
		return nil, fmt.Errorf("SHOW VIEWS not supported for driver %s", conn.Driver)
	}
	return Execute(ctx, conn, query)
}

// ShowFunctions returns a list of functions and procedures for the given connection.
func ShowFunctions(ctx context.Context, conn *Connection) (*QueryResult, error) {
	query, ok := functionQuery[conn.Driver]
	if !ok {
		return nil, fmt.Errorf("SHOW FUNCTIONS not supported for driver %s", conn.Driver)
	}
	return Execute(ctx, conn, query)
}
