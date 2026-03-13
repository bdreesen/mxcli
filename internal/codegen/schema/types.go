// SPDX-License-Identifier: Apache-2.0

// Package schema defines types for parsing Mendix reflection data JSON files.
package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ReflectionData is the root structure mapping qualified type names to their definitions.
type ReflectionData map[string]*TypeDefinition

// TypeDefinition represents a single type in the Mendix metamodel.
type TypeDefinition struct {
	QualifiedName      string                  `json:"qualifiedName"`      // e.g., "DomainModels$Entity"
	StorageName        string                  `json:"storageName"`        // Storage format name
	AllCompatibleTypes []string                `json:"allCompatibleTypes"` // All concrete implementations
	Abstract           bool                    `json:"abstract"`           // True if cannot be instantiated
	Namespace          bool                    `json:"namespace"`          // True if this is a namespace container
	Type               ElementType             `json:"type"`               // ELEMENT, MODEL_UNIT, STRUCTURAL_UNIT
	Properties         map[string]*PropertyDef `json:"properties"`         // Property definitions
	DefaultSettings    map[string]any          `json:"defaultSettings"`    // Default values for properties
}

// ElementType represents the category of a type definition.
type ElementType string

const (
	ElementTypeElement        ElementType = "ELEMENT"
	ElementTypeModelUnit      ElementType = "MODEL_UNIT"
	ElementTypeStructuralUnit ElementType = "STRUCTURAL_UNIT"
)

// PropertyDef represents a property definition within a type.
type PropertyDef struct {
	Name            string   `json:"name"`            // Property name
	StorageName     string   `json:"storageName"`     // How it's stored in BSON/JSON
	List            bool     `json:"list"`            // Is this a list/array?
	StorageListType int      `json:"storageListType"` // List type (1=simple, 2=contained)
	Public          bool     `json:"public"`          // Is it public/user-editable?
	Required        bool     `json:"required"`        // Must have a value?
	TypeInfo        TypeInfo `json:"typeInfo"`        // Type information
}

// TypeInfo describes the type of a property.
type TypeInfo struct {
	Type          TypeInfoType  `json:"type"`                    // PRIMITIVE, ENUMERATION, ELEMENT, UNIT
	PrimitiveType PrimitiveType `json:"primitiveType,omitempty"` // For PRIMITIVE
	Values        []string      `json:"values,omitempty"`        // For ENUMERATION
	ElementType   string        `json:"elementType,omitempty"`   // For ELEMENT - qualified type name
	Kind          ReferenceKind `json:"kind,omitempty"`          // For ELEMENT - PART, BY_ID_REFERENCE, etc.
	UnitType      string        `json:"unitType,omitempty"`      // For UNIT - qualified type name
}

// TypeInfoType represents the category of type information.
type TypeInfoType string

const (
	TypeInfoPrimitive   TypeInfoType = "PRIMITIVE"
	TypeInfoEnumeration TypeInfoType = "ENUMERATION"
	TypeInfoElement     TypeInfoType = "ELEMENT"
	TypeInfoUnit        TypeInfoType = "UNIT"
)

// PrimitiveType represents primitive data types.
type PrimitiveType string

const (
	PrimitiveString   PrimitiveType = "STRING"
	PrimitiveInteger  PrimitiveType = "INTEGER"
	PrimitiveLong     PrimitiveType = "LONG"
	PrimitiveDouble   PrimitiveType = "DOUBLE"
	PrimitiveBoolean  PrimitiveType = "BOOLEAN"
	PrimitiveGUID     PrimitiveType = "GUID"
	PrimitiveDateTime PrimitiveType = "DATE_TIME"
	PrimitivePoint    PrimitiveType = "POINT"
	PrimitiveSize     PrimitiveType = "SIZE"
	PrimitiveColor    PrimitiveType = "COLOR"
	PrimitiveBlob     PrimitiveType = "BLOB"
	PrimitiveUnknown  PrimitiveType = "Unknown"
)

// ReferenceKind represents how an element reference is stored.
type ReferenceKind string

const (
	ReferenceKindPart        ReferenceKind = "PART"                    // Contained/owned element
	ReferenceKindByID        ReferenceKind = "BY_ID_REFERENCE"         // Reference by ID
	ReferenceKindByName      ReferenceKind = "BY_NAME_REFERENCE"       // Reference by qualified name
	ReferenceKindLocalByName ReferenceKind = "LOCAL_BY_NAME_REFERENCE" // Reference by local name
)

// Load reads and parses a reflection data JSON file.
func Load(inputDir, version string) (ReflectionData, error) {
	filename := filepath.Join(inputDir, fmt.Sprintf("%s-structures.json", version))

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read reflection data: %w", err)
	}

	var result ReflectionData
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse reflection data: %w", err)
	}

	return result, nil
}

// GetNamespaces returns a sorted list of unique namespaces in the reflection data.
func (rd ReflectionData) GetNamespaces() []string {
	namespaces := make(map[string]bool)
	for qualifiedName := range rd {
		ns := extractNamespace(qualifiedName)
		if ns != "" {
			namespaces[ns] = true
		}
	}

	result := make([]string, 0, len(namespaces))
	for ns := range namespaces {
		result = append(result, ns)
	}
	return result
}

// GetTypesByNamespace returns all types belonging to a specific namespace.
func (rd ReflectionData) GetTypesByNamespace(namespace string) []*TypeDefinition {
	var result []*TypeDefinition
	for qualifiedName, typeDef := range rd {
		if extractNamespace(qualifiedName) == namespace {
			result = append(result, typeDef)
		}
	}
	return result
}

// extractNamespace extracts the namespace from a qualified name.
// e.g., "DomainModels$Entity" -> "DomainModels"
func extractNamespace(qualifiedName string) string {
	for i, c := range qualifiedName {
		if c == '$' {
			return qualifiedName[:i]
		}
	}
	return ""
}
