// SPDX-License-Identifier: Apache-2.0

package mpr

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// idToBsonBinary converts a UUID string to BSON Binary format.
// Mendix stores IDs as Binary with Subtype 0.
func idToBsonBinary(id string) primitive.Binary {
	blob := uuidToBlob(id)
	if blob == nil || len(blob) != 16 {
		// Generate a new UUID if the provided one is invalid
		blob = uuidToBlob(generateUUID())
	}
	return primitive.Binary{
		Subtype: 0x00,
		Data:    blob,
	}
}

// Writer provides methods to write Mendix project files.
type Writer struct {
	reader *Reader
}

// NewWriter creates a new writer from a reader opened in read-write mode.
func NewWriter(path string) (*Writer, error) {
	reader, err := OpenWithOptions(path, OpenOptions{ReadOnly: false})
	if err != nil {
		return nil, err
	}
	return &Writer{reader: reader}, nil
}

// Close closes the writer.
func (w *Writer) Close() error {
	return w.reader.Close()
}

// Reader returns the underlying reader.
func (w *Writer) Reader() *Reader {
	return w.reader
}

// Transaction support

// Transaction represents a database transaction.
type Transaction struct {
	tx     *sql.Tx
	writer *Writer
}

// BeginTransaction starts a new transaction.
func (w *Writer) BeginTransaction() (*Transaction, error) {
	tx, err := w.reader.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx, writer: w}, nil
}

// Commit commits the transaction.
func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction.
func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

// WriteTransaction provides atomic write operations for MPR v2 format.
// It coordinates database and file system changes to ensure consistency.
type WriteTransaction struct {
	tx           *sql.Tx
	writer       *Writer
	pendingFiles []pendingFile
	committed    bool
}

type pendingFile struct {
	tempPath  string
	finalPath string
}

// BeginWriteTransaction starts a new write transaction.
// For v2 format, this coordinates both database and file writes.
func (w *Writer) BeginWriteTransaction() (*WriteTransaction, error) {
	tx, err := w.reader.db.Begin()
	if err != nil {
		return nil, err
	}
	return &WriteTransaction{
		tx:           tx,
		writer:       w,
		pendingFiles: make([]pendingFile, 0),
	}, nil
}

// WriteUnit writes a unit within the transaction.
// The actual file write is deferred until Commit.
func (wt *WriteTransaction) WriteUnit(unitID string, contents []byte) error {
	unitIDBlob := uuidToBlob(unitID)

	if wt.writer.reader.version == MPRVersionV2 {
		swappedUUID := blobToUUIDSwapped(unitIDBlob)

		// Create directory if needed
		dir := filepath.Join(wt.writer.reader.contentsDir, swappedUUID[0:2], swappedUUID[2:4])
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Write to temp file first
		finalPath := filepath.Join(dir, swappedUUID+".mxunit")
		tempPath := finalPath + ".tmp"

		if err := os.WriteFile(tempPath, contents, 0644); err != nil {
			return fmt.Errorf("failed to write temp file: %w", err)
		}

		wt.pendingFiles = append(wt.pendingFiles, pendingFile{
			tempPath:  tempPath,
			finalPath: finalPath,
		})

		// Update hash in DB
		hash := sha256.Sum256(contents)
		contentsHash := base64.StdEncoding.EncodeToString(hash[:])
		_, err := wt.tx.Exec(`
			UPDATE Unit SET ContentsHash = ? WHERE UnitID = ?
		`, contentsHash, unitIDBlob)
		return err
	}

	// V1: Update in database directly
	_, err := wt.tx.Exec(`
		UPDATE Unit SET Contents = ? WHERE UnitID = ?
	`, contents, unitIDBlob)
	return err
}

// Commit commits the transaction.
// For v2, this first commits the database, then finalizes file writes.
func (wt *WriteTransaction) Commit() error {
	if wt.committed {
		return fmt.Errorf("transaction already committed")
	}

	// Commit database transaction first
	if err := wt.tx.Commit(); err != nil {
		// Clean up temp files
		wt.cleanupTempFiles()
		return err
	}

	// Finalize file writes by renaming temp files to final paths
	for _, pf := range wt.pendingFiles {
		if err := os.Rename(pf.tempPath, pf.finalPath); err != nil {
			// Log error but continue - DB is already committed
			// This could leave some files in inconsistent state
			fmt.Printf("Warning: failed to finalize file %s: %v\n", pf.finalPath, err)
		}
	}

	wt.committed = true
	return nil
}

// Rollback rolls back the transaction and cleans up temp files.
func (wt *WriteTransaction) Rollback() error {
	if wt.committed {
		return fmt.Errorf("transaction already committed")
	}

	// Clean up temp files
	wt.cleanupTempFiles()

	// Rollback database
	return wt.tx.Rollback()
}

func (wt *WriteTransaction) cleanupTempFiles() {
	for _, pf := range wt.pendingFiles {
		os.Remove(pf.tempPath)
	}
}

// generateUUID generates a new UUID v4 for model elements.
// Returns format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant is 10

	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		b[0], b[1], b[2], b[3],
		b[4], b[5],
		b[6], b[7],
		b[8], b[9],
		b[10], b[11], b[12], b[13], b[14], b[15])
}

// uuidToBlob converts a UUID string to a 16-byte blob in Microsoft GUID format.
// UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// Microsoft GUID format byte-swaps the first 3 groups (little-endian):
// - First 4 bytes: reversed
// - Next 2 bytes: reversed
// - Next 2 bytes: reversed
// - Last 8 bytes: unchanged
func uuidToBlob(uuid string) []byte {
	if uuid == "" {
		return nil
	}
	// Remove dashes
	var clean strings.Builder
	for _, c := range uuid {
		if c != '-' {
			clean.WriteString(string(c))
		}
	}
	// Decode hex to bytes
	decoded, err := hex.DecodeString(clean.String())
	if err != nil || len(decoded) != 16 {
		return nil
	}
	// Swap bytes to Microsoft GUID format
	blob := make([]byte, 16)
	// First 4 bytes: reversed
	blob[0] = decoded[3]
	blob[1] = decoded[2]
	blob[2] = decoded[1]
	blob[3] = decoded[0]
	// Next 2 bytes: reversed
	blob[4] = decoded[5]
	blob[5] = decoded[4]
	// Next 2 bytes: reversed
	blob[6] = decoded[7]
	blob[7] = decoded[6]
	// Last 8 bytes: unchanged
	copy(blob[8:], decoded[8:])
	return blob
}
