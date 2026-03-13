// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debug_schema <mpr-path>")
		os.Exit(1)
	}

	mprPath := os.Args[1]
	db, err := sql.Open("sqlite", mprPath+"?mode=ro")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Check what tables exist
	rows0, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		panic(err)
	}
	fmt.Println("=== Tables ===")
	for rows0.Next() {
		var name string
		rows0.Scan(&name)
		fmt.Printf("  %s\n", name)
	}
	rows0.Close()

	// Check table schema
	rows, err := db.Query("PRAGMA table_info(Unit)")
	if err != nil {
		panic(err)
	}
	fmt.Println("\n=== Unit Table Schema ===")
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull int
		var dfltValue any
		var pk int
		rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk)
		fmt.Printf("  %s: %s\n", name, typ)
	}
	rows.Close()

	// Check sample data
	rows2, err := db.Query("SELECT UnitID, ContentsHash FROM Unit LIMIT 3")
	if err != nil {
		fmt.Printf("Error querying Unit: %v\n", err)
		return
	}
	fmt.Println("\n=== Sample Unit Data ===")
	contentsDir := filepath.Join(filepath.Dir(mprPath), "mprcontents")
	for rows2.Next() {
		var unitID []byte
		var contentsHash string
		rows2.Scan(&unitID, &contentsHash)
		fmt.Printf("UnitID: %s (len=%d)\n", hex.EncodeToString(unitID), len(unitID))
		fmt.Printf("ContentsHash (raw): %s (len=%d)\n", contentsHash, len(contentsHash))

		// Try to read file and compute hash
		uuid := blobToUUIDSwapped(unitID)
		path := filepath.Join(contentsDir, uuid[0:2], uuid[2:4], uuid+".mxunit")
		if data, err := os.ReadFile(path); err == nil {
			// Compute SHA256 of contents and base64 encode
			hash := sha256.Sum256(data)
			hashBase64 := base64.StdEncoding.EncodeToString(hash[:])
			fmt.Printf("Computed SHA256 (base64): %s\n", hashBase64)
			fmt.Printf("Match: %v\n", hashBase64 == contentsHash)
		}
		fmt.Println()
	}
	rows2.Close()
}

func blobToUUIDSwapped(blob []byte) string {
	if len(blob) != 16 {
		return ""
	}
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		blob[3], blob[2], blob[1], blob[0],
		blob[5], blob[4],
		blob[7], blob[6],
		blob[8], blob[9],
		blob[10], blob[11], blob[12], blob[13], blob[14], blob[15])
}
