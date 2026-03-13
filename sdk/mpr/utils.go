// SPDX-License-Identifier: Apache-2.0

package mpr

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateID generates a new unique ID for model elements.
func GenerateID() string {
	return generateUUID()
}

// BlobToUUID converts a binary ID blob to a UUID string.
func BlobToUUID(data []byte) string {
	return blobToUUID(data)
}

// IDToBsonBinary converts a UUID string to a BSON binary value.
func IDToBsonBinary(id string) primitive.Binary {
	return idToBsonBinary(id)
}

// BsonBinaryToID converts a BSON binary value to a UUID string.
func BsonBinaryToID(bin primitive.Binary) string {
	return BlobToUUID(bin.Data)
}

// Hash computes a hash for content (used for content deduplication).
func Hash(content []byte) string {
	// Simple hash for now - could use crypto/sha256 for better hashing
	var sum uint64
	for i, b := range content {
		sum += uint64(b) * uint64(i+1)
	}
	return fmt.Sprintf("%016x", sum)
}

// ValidateID checks if an ID is valid.
func ValidateID(id string) bool {
	if len(id) != 36 {
		return false
	}
	// Check UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	for i, c := range id {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		} else {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}
