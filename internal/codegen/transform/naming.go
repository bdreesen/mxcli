// SPDX-License-Identifier: Apache-2.0

// Package transform provides utilities for transforming Mendix metamodel types to Go types.
package transform

import (
	"strings"
	"unicode"
)

// ExtractNamespace extracts the namespace from a qualified name.
// e.g., "DomainModels$Entity" -> "DomainModels"
func ExtractNamespace(qualifiedName string) string {
	before, _, ok := strings.Cut(qualifiedName, "$")
	if !ok {
		return ""
	}
	return before
}

// ExtractTypeName extracts the type name from a qualified name.
// e.g., "DomainModels$Entity" -> "Entity"
func ExtractTypeName(qualifiedName string) string {
	_, after, ok := strings.Cut(qualifiedName, "$")
	if !ok {
		return qualifiedName
	}
	return after
}

// ToGoPackage converts a namespace to a Go package name.
// e.g., "DomainModels" -> "domainmodels"
func ToGoPackage(namespace string) string {
	return strings.ToLower(namespace)
}

// ToGoTypeName ensures a type name is a valid Go exported identifier.
// e.g., "Entity" -> "Entity" (no change needed)
// e.g., "XMLSchema" -> "XMLSchema"
func ToGoTypeName(typeName string) string {
	// Handle some special cases
	if typeName == "" {
		return ""
	}

	// Ensure first character is uppercase for export
	runes := []rune(typeName)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// StorageToJSONTag converts a storage name to a JSON tag (camelCase).
// e.g., "MyFieldName" -> "myFieldName"
// e.g., "$ID" -> "$ID" (special fields unchanged)
func StorageToJSONTag(storageName string) string {
	if storageName == "" {
		return ""
	}

	// Special Mendix fields keep their format
	if strings.HasPrefix(storageName, "$") {
		return storageName
	}

	// Convert to camelCase
	runes := []rune(storageName)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// ToGoFieldName converts a property name to a Go field name.
// e.g., "name" -> "Name"
// e.g., "entityId" -> "EntityID" (special handling for Id suffix)
func ToGoFieldName(propertyName string) string {
	if propertyName == "" {
		return ""
	}

	result := strings.Builder{}
	runes := []rune(propertyName)

	// First character uppercase
	result.WriteRune(unicode.ToUpper(runes[0]))

	// Rest of the string
	for i := 1; i < len(runes); i++ {
		result.WriteRune(runes[i])
	}

	name := result.String()

	// Handle common suffixes for Go conventions
	name = strings.ReplaceAll(name, "Id", "ID")
	name = strings.ReplaceAll(name, "Url", "URL")
	name = strings.ReplaceAll(name, "Uri", "URI")
	name = strings.ReplaceAll(name, "Xml", "XML")
	name = strings.ReplaceAll(name, "Json", "JSON")
	name = strings.ReplaceAll(name, "Html", "HTML")
	name = strings.ReplaceAll(name, "Http", "HTTP")
	name = strings.ReplaceAll(name, "Uuid", "UUID")
	name = strings.ReplaceAll(name, "Api", "API")

	return name
}

// ToMarkerMethodName creates a marker method name for an interface.
// e.g., "AttributeType" -> "isAttributeType"
func ToMarkerMethodName(typeName string) string {
	if typeName == "" {
		return ""
	}

	runes := []rune(typeName)
	runes[0] = unicode.ToLower(runes[0])
	return "is" + string(runes)
}

// ToEnumTypeName creates an enum type name from property context.
// e.g., namespace="DomainModels", property="deletingBehavior" -> "DeletingBehavior"
func ToEnumTypeName(namespace, propertyName string) string {
	return ToGoFieldName(propertyName)
}

// ToEnumValueName creates a valid Go constant name for an enum value.
// e.g., enumType="DeletingBehavior", value="DeleteMeAndReferences" -> "DeletingBehaviorDeleteMeAndReferences"
func ToEnumValueName(enumType, value string) string {
	// Clean up the value for Go naming
	cleaned := strings.ReplaceAll(value, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	return enumType + ToGoFieldName(cleaned)
}

// SanitizeGoIdentifier ensures a string is a valid Go identifier.
func SanitizeGoIdentifier(name string) string {
	if name == "" {
		return "_"
	}

	result := strings.Builder{}
	for i, r := range name {
		if i == 0 {
			if unicode.IsLetter(r) || r == '_' {
				result.WriteRune(r)
			} else {
				result.WriteRune('_')
			}
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				result.WriteRune(r)
			}
		}
	}

	// Check for Go reserved words
	sanitized := result.String()
	if isGoKeyword(sanitized) {
		return sanitized + "_"
	}

	return sanitized
}

// isGoKeyword returns true if the string is a Go reserved keyword.
func isGoKeyword(s string) bool {
	keywords := map[string]bool{
		"break": true, "case": true, "chan": true, "const": true, "continue": true,
		"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
		"func": true, "go": true, "goto": true, "if": true, "import": true,
		"interface": true, "map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true, "var": true,
	}
	return keywords[s]
}

// QualifiedNameToImport converts a qualified type name to its import path and local type.
// e.g., "DomainModels$Entity" with current namespace "Microflows"
// -> import "domainmodels", type "domainmodels.Entity"
func QualifiedNameToImport(qualifiedName, currentNamespace string) (importPath, typeName string, needsImport bool) {
	targetNamespace := ExtractNamespace(qualifiedName)
	targetType := ExtractTypeName(qualifiedName)

	if targetNamespace == currentNamespace || targetNamespace == "" {
		return "", ToGoTypeName(targetType), false
	}

	pkg := ToGoPackage(targetNamespace)
	return pkg, pkg + "." + ToGoTypeName(targetType), true
}
