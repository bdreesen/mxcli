// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/mendixlabs/mxcli/mdl/ast"
	"github.com/mendixlabs/mxcli/mdl/visitor"
	sqllib "github.com/mendixlabs/mxcli/sql"
)

// ensureSQLManager lazily initializes the SQL connection manager.
func (e *Executor) ensureSQLManager() *sqllib.Manager {
	if e.sqlMgr == nil {
		e.sqlMgr = sqllib.NewManager()
	}
	return e.sqlMgr
}

// execSQLConnect handles SQL CONNECT <driver> '<dsn>' AS <alias>
func (e *Executor) execSQLConnect(s *ast.SQLConnectStmt) error {
	driver, err := sqllib.ParseDriver(s.Driver)
	if err != nil {
		return err
	}

	mgr := e.ensureSQLManager()
	if err := mgr.Connect(driver, s.DSN, s.Alias); err != nil {
		return err
	}

	fmt.Fprintf(e.output, "Connected to %s database as '%s'\n", driver, s.Alias)
	return nil
}

// execSQLDisconnect handles SQL DISCONNECT <alias>
func (e *Executor) execSQLDisconnect(s *ast.SQLDisconnectStmt) error {
	mgr := e.ensureSQLManager()
	if err := mgr.Disconnect(s.Alias); err != nil {
		return err
	}

	fmt.Fprintf(e.output, "Disconnected '%s'\n", s.Alias)
	return nil
}

// execSQLConnections handles SQL CONNECTIONS
func (e *Executor) execSQLConnections() error {
	mgr := e.ensureSQLManager()
	infos := mgr.List()

	if len(infos) == 0 {
		fmt.Fprintln(e.output, "No active SQL connections")
		return nil
	}

	// Sort by alias for stable output
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Alias < infos[j].Alias
	})

	result := &sqllib.QueryResult{
		Columns: []string{"Alias", "Driver"},
	}
	for _, info := range infos {
		result.Rows = append(result.Rows, []any{info.Alias, string(info.Driver)})
	}
	sqllib.FormatTable(e.output, result)
	return nil
}

// execSQLQuery handles SQL <alias> <raw-sql>
func (e *Executor) execSQLQuery(s *ast.SQLQueryStmt) error {
	mgr := e.ensureSQLManager()
	conn, err := mgr.Get(s.Alias)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := sqllib.Execute(ctx, conn, s.Query)
	if err != nil {
		return err
	}

	sqllib.FormatTable(e.output, result)
	fmt.Fprintf(e.output, "(%d rows)\n", len(result.Rows))
	return nil
}

// execSQLShowTables handles SQL <alias> SHOW TABLES
func (e *Executor) execSQLShowTables(s *ast.SQLShowTablesStmt) error {
	mgr := e.ensureSQLManager()
	conn, err := mgr.Get(s.Alias)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := sqllib.ShowTables(ctx, conn)
	if err != nil {
		return err
	}

	sqllib.FormatTable(e.output, result)
	fmt.Fprintf(e.output, "(%d tables)\n", len(result.Rows))
	return nil
}

// execSQLShowViews handles SQL <alias> SHOW VIEWS
func (e *Executor) execSQLShowViews(s *ast.SQLShowViewsStmt) error {
	mgr := e.ensureSQLManager()
	conn, err := mgr.Get(s.Alias)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := sqllib.ShowViews(ctx, conn)
	if err != nil {
		return err
	}

	sqllib.FormatTable(e.output, result)
	fmt.Fprintf(e.output, "(%d views)\n", len(result.Rows))
	return nil
}

// execSQLShowFunctions handles SQL <alias> SHOW FUNCTIONS
func (e *Executor) execSQLShowFunctions(s *ast.SQLShowFunctionsStmt) error {
	mgr := e.ensureSQLManager()
	conn, err := mgr.Get(s.Alias)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := sqllib.ShowFunctions(ctx, conn)
	if err != nil {
		return err
	}

	sqllib.FormatTable(e.output, result)
	fmt.Fprintf(e.output, "(%d functions)\n", len(result.Rows))
	return nil
}

// execSQLGenerateConnector handles SQL <alias> GENERATE CONNECTOR INTO <module> [TABLES (...)] [VIEWS (...)] [EXEC]
func (e *Executor) execSQLGenerateConnector(s *ast.SQLGenerateConnectorStmt) error {
	mgr := e.ensureSQLManager()
	conn, err := mgr.Get(s.Alias)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cfg := &sqllib.GenerateConfig{
		Conn:   conn,
		Module: s.Module,
		Alias:  s.Alias,
		Tables: s.Tables,
		Views:  s.Views,
	}

	result, err := sqllib.GenerateConnector(ctx, cfg)
	if err != nil {
		return err
	}

	// Report skipped columns
	for _, skip := range result.SkippedCols {
		fmt.Fprintf(e.output, "-- WARNING: skipped unmappable column: %s\n", skip)
	}

	if s.Exec {
		// Execute constants + entities (parseable by mxcli)
		fmt.Fprintf(e.output, "Generating connector (%d tables, %d views)...\n",
			result.TableCount, result.ViewCount)
		if err := e.executeGeneratedMDL(result.ExecutableMDL); err != nil {
			return err
		}
		// Print DATABASE CONNECTION as reference (not yet executable)
		fmt.Fprintf(e.output, "\n-- Database Connection definition (configure in Studio Pro with Database Connector module):\n")
		fmt.Fprint(e.output, result.ConnectionMDL)
		return nil
	}

	// Print complete MDL to output
	fmt.Fprint(e.output, result.MDL)
	fmt.Fprintf(e.output, "\n-- Generated: %d tables, %d views\n", result.TableCount, result.ViewCount)
	return nil
}

// executeGeneratedMDL parses and executes MDL text as if it were a script.
func (e *Executor) executeGeneratedMDL(mdl string) error {
	prog, errs := visitor.Build(mdl)
	if len(errs) > 0 {
		return fmt.Errorf("failed to parse generated MDL: %v", errs[0])
	}
	return e.ExecuteProgram(prog)
}

// execSQLDescribeTable handles SQL <alias> DESCRIBE <table>
func (e *Executor) execSQLDescribeTable(s *ast.SQLDescribeTableStmt) error {
	mgr := e.ensureSQLManager()
	conn, err := mgr.Get(s.Alias)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := sqllib.DescribeTable(ctx, conn, s.Table)
	if err != nil {
		return err
	}

	sqllib.FormatTable(e.output, result)
	fmt.Fprintf(e.output, "(%d columns)\n", len(result.Rows))
	return nil
}
