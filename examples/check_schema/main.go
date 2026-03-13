// SPDX-License-Identifier: Apache-2.0

package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: check_schema <mpr-path>")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", os.Args[1]+"?mode=ro")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Get table info for Unit
	rows, err := db.Query("PRAGMA table_info(Unit)")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Println("Unit table columns:")
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dflt_value any
		rows.Scan(&cid, &name, &typ, &notnull, &dflt_value, &pk)
		fmt.Printf("  %d: %s (%s)\n", cid, name, typ)
	}

	// Show all units and check for matching files
	fmt.Println("\nSearching for matching files...")
	rows2, err := db.Query("SELECT UnitID FROM Unit")
	if err != nil {
		panic(err)
	}
	defer rows2.Close()

	mprContentsDir := "/workspaces/ModelSDKGo/mx-test-projects/test1-go-app/mprcontents"
	found := 0
	notFound := 0

	for rows2.Next() {
		var unitID []byte
		rows2.Scan(&unitID)

		uuid := blobToUUID(unitID)
		uuidSwapped := blobToUUIDSwapped(unitID)

		// Try standard format
		path1 := fmt.Sprintf("%s/%s/%s/%s.mxunit", mprContentsDir, uuid[0:2], uuid[2:4], uuid)
		// Try swapped format
		path2 := fmt.Sprintf("%s/%s/%s/%s.mxunit", mprContentsDir, uuidSwapped[0:2], uuidSwapped[2:4], uuidSwapped)

		if _, err := os.Stat(path1); err == nil {
			found++
			if found <= 3 {
				fmt.Printf("  Found (std): %s\n", path1)
			}
		} else if _, err := os.Stat(path2); err == nil {
			found++
			if found <= 3 {
				fmt.Printf("  Found (swapped): %s\n", path2)
			}
		} else {
			notFound++
			if notFound <= 3 {
				fmt.Printf("  NOT found: %s or %s\n", uuid, uuidSwapped)
			}
		}
	}

	fmt.Printf("\nTotal found: %d, not found: %d\n", found, notFound)

	// Also check the reverse: find a file and see if it's in DB
	fmt.Println("\nChecking existing file against DB...")
	testFile := "61340cb4-c994-469f-8595-528f14a300af"
	fmt.Printf("  File: %s.mxunit\n", testFile)
}

func blobToUUID(blob []byte) string {
	if len(blob) != 16 {
		return hex.EncodeToString(blob)
	}
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		blob[0], blob[1], blob[2], blob[3],
		blob[4], blob[5],
		blob[6], blob[7],
		blob[8], blob[9],
		blob[10], blob[11], blob[12], blob[13], blob[14], blob[15])
}

// blobToUUIDSwapped converts with little-endian first 3 groups (Microsoft GUID style)
func blobToUUIDSwapped(blob []byte) string {
	if len(blob) != 16 {
		return hex.EncodeToString(blob)
	}
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		blob[3], blob[2], blob[1], blob[0],
		blob[5], blob[4],
		blob[7], blob[6],
		blob[8], blob[9],
		blob[10], blob[11], blob[12], blob[13], blob[14], blob[15])
}
