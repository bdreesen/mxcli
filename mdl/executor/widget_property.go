// SPDX-License-Identifier: Apache-2.0

// Package executor - Widget property navigation for UPDATE WIDGETS
package executor

import (
	"fmt"
	"reflect"
	"strings"

	mdlerrors "github.com/mendixlabs/mxcli/mdl/errors"
	"github.com/mendixlabs/mxcli/model"
	"github.com/mendixlabs/mxcli/sdk/pages"
	"go.mongodb.org/mongo-driver/bson"
)

// getWidgetID extracts the ID from a widget.
func getWidgetID(widget any) string {
	if widget == nil {
		return ""
	}

	// Use reflection to get the ID field
	v := reflect.ValueOf(widget)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}

	// Try to get BaseWidget.ID
	if baseWidget := v.FieldByName("BaseWidget"); baseWidget.IsValid() {
		if idField := baseWidget.FieldByName("ID"); idField.IsValid() {
			return fmt.Sprintf("%v", idField.Interface())
		}
	}

	// Direct ID field
	if idField := v.FieldByName("ID"); idField.IsValid() {
		return fmt.Sprintf("%v", idField.Interface())
	}

	return ""
}

// setWidgetProperty sets a property value on a widget by path.
func setWidgetProperty(widget any, path string, value any) error {
	if widget == nil {
		return mdlerrors.NewValidation("widget is nil")
	}

	// Handle CustomWidget specifically
	if cw, ok := widget.(*pages.CustomWidget); ok {
		return setCustomWidgetProperty(cw, path, value)
	}

	// For other widget types, use reflection to set simple fields
	return setWidgetFieldByReflection(widget, path, value)
}

// setCustomWidgetProperty sets a property on a CustomWidget.
func setCustomWidgetProperty(cw *pages.CustomWidget, path string, value any) error {
	// CustomWidget can have properties in RawObject (BSON) or WidgetObject (structured)
	if cw.RawObject != nil {
		return updateBSONWidgetProperty(cw.RawObject, path, value)
	}

	if cw.WidgetObject != nil {
		return updateStructuredWidgetProperty(cw.WidgetObject, cw.PropertyTypeIDMap, path, value)
	}

	return mdlerrors.NewValidation("widget has no property data")
}

// updateBSONWidgetProperty updates a property in a BSON document.
func updateBSONWidgetProperty(doc bson.D, path string, value any) error {
	// Find the "Object" field which contains widget properties
	for i := range doc {
		if doc[i].Key == "Object" {
			if objDoc, ok := doc[i].Value.(bson.D); ok {
				// Find Properties array
				for j := range objDoc {
					if objDoc[j].Key == "Properties" {
						if props, ok := objDoc[j].Value.(bson.A); ok {
							return updateBSONPropertyByKey(props, path, value)
						}
					}
				}
			}
		}
	}
	return mdlerrors.NewNotFound("property", path)
}

// updateBSONPropertyByKey finds and updates a property by key in a BSON array.
func updateBSONPropertyByKey(props bson.A, path string, value any) error {
	for i, prop := range props {
		if propDoc, ok := prop.(bson.D); ok {
			// Find the Key field
			for _, field := range propDoc {
				if field.Key == "Key" {
					if key, ok := field.Value.(string); ok && key == path {
						// Found the property, now update its Value
						return updateBSONPropertyValueAtIndex(props, i, value)
					}
				}
			}
		}
	}
	return mdlerrors.NewNotFound("property", path)
}

// updateBSONPropertyValueAtIndex updates the Value field of a property at the given index.
func updateBSONPropertyValueAtIndex(props bson.A, index int, newValue any) error {
	propDoc, ok := props[index].(bson.D)
	if !ok {
		return mdlerrors.NewValidation("property is not a BSON document")
	}

	for i := range propDoc {
		if propDoc[i].Key == "Value" {
			if valueDoc, ok := propDoc[i].Value.(bson.D); ok {
				// Find PrimitiveValue in the Value document
				for j := range valueDoc {
					if valueDoc[j].Key == "PrimitiveValue" {
						valueDoc[j].Value = convertToBSONValue(newValue)
						propDoc[i].Value = valueDoc
						props[index] = propDoc
						return nil
					}
				}
				// PrimitiveValue not found, try to add it
				valueDoc = append(valueDoc, bson.E{Key: "PrimitiveValue", Value: convertToBSONValue(newValue)})
				propDoc[i].Value = valueDoc
				props[index] = propDoc
				return nil
			}
		}
	}

	return mdlerrors.NewValidation("Value field not found in property")
}

// convertToBSONValue converts a Go value to appropriate BSON format.
func convertToBSONValue(value any) any {
	switch v := value.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	case string:
		return v
	case int, int64, int32:
		return fmt.Sprintf("%d", v)
	case float64, float32:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// updateStructuredWidgetProperty updates a property in a structured WidgetObject.
func updateStructuredWidgetProperty(obj *pages.WidgetObject, typeMap map[string]pages.PropertyTypeIDEntry, path string, value any) error {
	if obj == nil || obj.Properties == nil {
		return mdlerrors.NewValidation("widget object has no properties")
	}

	bsonValue := convertToBSONValue(value)
	strValue, ok := bsonValue.(string)
	if !ok {
		strValue = fmt.Sprintf("%v", bsonValue)
	}

	// Find the property by PropertyKey
	for i, prop := range obj.Properties {
		if prop.PropertyKey == path {
			// Update the primitive value
			if prop.Value != nil {
				obj.Properties[i].Value.PrimitiveValue = strValue
				return nil
			}
			// Create a value if it doesn't exist
			obj.Properties[i].Value = &pages.WidgetValue{
				PrimitiveValue: strValue,
			}
			return nil
		}
	}

	// Check if we can find it via the type map (case-insensitive)
	if typeMap != nil {
		for key := range typeMap {
			if strings.EqualFold(key, path) {
				// Find the actual property
				for i, prop := range obj.Properties {
					if prop.PropertyKey == key {
						if prop.Value != nil {
							obj.Properties[i].Value.PrimitiveValue = strValue
							return nil
						}
						obj.Properties[i].Value = &pages.WidgetValue{
							PrimitiveValue: strValue,
						}
						return nil
					}
				}
			}
		}
	}

	return mdlerrors.NewNotFound("property", path)
}

// setWidgetFieldByReflection sets a simple field on a widget using reflection.
func setWidgetFieldByReflection(widget any, fieldName string, value any) error {
	v := reflect.ValueOf(widget)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return mdlerrors.NewValidation("widget is not a struct")
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return mdlerrors.NewNotFound("field", fieldName)
	}
	if !field.CanSet() {
		return mdlerrors.NewValidationf("field not settable: %s", fieldName)
	}

	// Convert value to field type
	fieldType := field.Type()
	valueVal := reflect.ValueOf(value)

	if valueVal.Type().ConvertibleTo(fieldType) {
		field.Set(valueVal.Convert(fieldType))
		return nil
	}

	// Handle string conversion for common types
	if fieldType.Kind() == reflect.String {
		field.SetString(fmt.Sprintf("%v", value))
		return nil
	}

	return mdlerrors.NewValidationf("cannot convert %T to %s", value, fieldType)
}

// walkPageWidgets walks all widgets in a page and calls the visitor function.
func walkPageWidgets(page *pages.Page, visitor func(widget any) error) error {
	if page == nil || page.LayoutCall == nil {
		return nil
	}

	// Walk through layout call arguments (each argument has a single widget)
	for _, arg := range page.LayoutCall.Arguments {
		if arg.Widget != nil {
			if err := walkWidget(arg.Widget, visitor); err != nil {
				return err
			}
		}
	}

	return nil
}

// walkSnippetWidgets walks all widgets in a snippet and calls the visitor function.
func walkSnippetWidgets(snippet *pages.Snippet, visitor func(widget any) error) error {
	if snippet == nil {
		return nil
	}

	for _, widget := range snippet.Widgets {
		if err := walkWidget(widget, visitor); err != nil {
			return err
		}
	}

	return nil
}

// walkWidget recursively walks a widget and its children.
func walkWidget(widget pages.Widget, visitor func(widget any) error) error {
	if widget == nil {
		return nil
	}

	// Visit this widget
	if err := visitor(widget); err != nil {
		return err
	}

	// Recursively walk children based on widget type
	switch w := widget.(type) {
	case *pages.LayoutGrid:
		for _, row := range w.Rows {
			for _, col := range row.Columns {
				for _, child := range col.Widgets {
					if err := walkWidget(child, visitor); err != nil {
						return err
					}
				}
			}
		}
	case *pages.DataView:
		for _, child := range w.Widgets {
			if err := walkWidget(child, visitor); err != nil {
				return err
			}
		}
		for _, child := range w.FooterWidgets {
			if err := walkWidget(child, visitor); err != nil {
				return err
			}
		}
	case *pages.ListView:
		for _, child := range w.Widgets {
			if err := walkWidget(child, visitor); err != nil {
				return err
			}
		}
	case *pages.Container:
		for _, child := range w.Widgets {
			if err := walkWidget(child, visitor); err != nil {
				return err
			}
		}
	case *pages.GroupBox:
		for _, child := range w.Widgets {
			if err := walkWidget(child, visitor); err != nil {
				return err
			}
		}
	case *pages.TabContainer:
		for _, pg := range w.TabPages {
			for _, child := range pg.Widgets {
				if err := walkWidget(child, visitor); err != nil {
					return err
				}
			}
		}
	case *pages.ScrollContainer:
		for _, child := range w.Widgets {
			if err := walkWidget(child, visitor); err != nil {
				return err
			}
		}
	case *pages.CustomWidget:
		// Custom widgets may have nested widgets in their value properties
		if w.WidgetObject != nil {
			for _, prop := range w.WidgetObject.Properties {
				if prop.Value != nil {
					for _, child := range prop.Value.Widgets {
						if err := walkWidget(child, visitor); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

// getSnippetByID finds a snippet by ID from the list.
func getSnippetByID(snippets []*pages.Snippet, id model.ID) *pages.Snippet {
	for _, s := range snippets {
		if s.ID == id {
			return s
		}
	}
	return nil
}
