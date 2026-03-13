// SPDX-License-Identifier: Apache-2.0

// Example: Inspecting MPR database schema
package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: inspect_mpr <path-to-mpr-file>")
		os.Exit(1)
	}

	mprPath := os.Args[1]

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?mode=ro", mprPath))
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("=== Tables ===")
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		fmt.Printf("Error querying tables: %v\n", err)
		os.Exit(1)
	}

	var tables []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		tables = append(tables, name)
		fmt.Printf("  %s\n", name)
	}
	rows.Close()

	// Show schema for each table
	fmt.Println("\n=== Table Schemas ===")
	for _, table := range tables {
		fmt.Printf("\n--- %s ---\n", table)
		rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		for rows.Next() {
			var cid int
			var name, colType string
			var notNull, pk int
			var dfltValue sql.NullString
			rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
			pkMark := ""
			if pk > 0 {
				pkMark = " (PK)"
			}
			fmt.Printf("  %s: %s%s\n", name, colType, pkMark)
		}
		rows.Close()

		// Show sample data
		rows, err = db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 3", table))
		if err != nil {
			continue
		}

		cols, _ := rows.Columns()
		fmt.Printf("  Columns: %v\n", cols)

		// Count rows
		var count int
		db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		fmt.Printf("  Row count: %d\n", count)

		rows.Close()
	}

	// Check for Unit table specifically
	fmt.Println("\n=== Unit Table Sample ===")
	rows, err = db.Query("SELECT * FROM Unit LIMIT 5")
	if err != nil {
		fmt.Printf("Unit table not found or error: %v\n", err)
	} else {
		cols, _ := rows.Columns()
		fmt.Printf("Columns: %v\n", cols)

		for rows.Next() {
			values := make([]any, len(cols))
			valuePtrs := make([]any, len(cols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)

			fmt.Println("Row:")
			for i, col := range cols {
				val := values[i]
				switch v := val.(type) {
				case []byte:
					if len(v) > 50 {
						fmt.Printf("  %s: [%d bytes]\n", col, len(v))
					} else {
						fmt.Printf("  %s: %v\n", col, v)
					}
				default:
					fmt.Printf("  %s: %v\n", col, val)
				}
			}
		}
		rows.Close()
	}
}
