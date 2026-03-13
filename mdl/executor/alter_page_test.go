// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mendixlabs/mxcli/mdl/ast"
)

// Helper to build a minimal raw BSON page structure for testing.
func makeRawPage(widgets ...map[string]any) map[string]any {
	widgetArr := []any{int32(2)} // type marker
	for _, w := range widgets {
		widgetArr = append(widgetArr, w)
	}
	return map[string]any{
		"FormCall": map[string]any{
			"Arguments": []any{
				int32(2), // type marker
				map[string]any{
					"Widgets": widgetArr,
				},
			},
		},
	}
}

func makeWidget(name string, typeName string) map[string]any {
	return map[string]any{
		"$Type": typeName,
		"Name":  name,
	}
}

func makeContainerWidget(name string, children ...map[string]any) map[string]any {
	childArr := []any{int32(2)} // type marker
	for _, c := range children {
		childArr = append(childArr, c)
	}
	return map[string]any{
		"$Type":   "Pages$DivContainer",
		"Name":    name,
		"Widgets": childArr,
	}
}

func TestFindBsonWidget_TopLevel(t *testing.T) {
	w1 := makeWidget("txtName", "Pages$TextBox")
	w2 := makeWidget("txtEmail", "Pages$TextBox")
	rawData := makeRawPage(w1, w2)

	result := findBsonWidget(rawData, "txtName")
	if result == nil {
		t.Fatal("Expected to find txtName")
	}
	name, _ := result.widget["Name"].(string)
	if name != "txtName" {
		t.Errorf("Expected name 'txtName', got %q", name)
	}
	if result.index != 0 {
		t.Errorf("Expected index 0, got %d", result.index)
	}
}

func TestFindBsonWidget_Nested(t *testing.T) {
	inner := makeWidget("txtInner", "Pages$TextBox")
	container := makeContainerWidget("ctn1", inner)
	rawData := makeRawPage(container)

	result := findBsonWidget(rawData, "txtInner")
	if result == nil {
		t.Fatal("Expected to find txtInner inside container")
	}
	name, _ := result.widget["Name"].(string)
	if name != "txtInner" {
		t.Errorf("Expected name 'txtInner', got %q", name)
	}
}

func TestFindBsonWidget_NotFound(t *testing.T) {
	w1 := makeWidget("txtName", "Pages$TextBox")
	rawData := makeRawPage(w1)

	result := findBsonWidget(rawData, "nonexistent")
	if result != nil {
		t.Error("Expected nil for nonexistent widget")
	}
}

func TestApplyDropWidget_Single(t *testing.T) {
	w1 := makeWidget("txtName", "Pages$TextBox")
	w2 := makeWidget("txtEmail", "Pages$TextBox")
	w3 := makeWidget("txtPhone", "Pages$TextBox")
	rawData := makeRawPage(w1, w2, w3)

	op := &ast.DropWidgetOp{WidgetNames: []string{"txtEmail"}}
	if err := applyDropWidget(rawData, op); err != nil {
		t.Fatalf("applyDropWidget failed: %v", err)
	}

	// Verify txtEmail was removed
	formCall := rawData["FormCall"].(map[string]any)
	args := getBsonArrayElements(formCall["Arguments"])
	argMap := args[0].(map[string]any)
	widgets := getBsonArrayElements(argMap["Widgets"])

	if len(widgets) != 2 {
		t.Fatalf("Expected 2 widgets after drop, got %d", len(widgets))
	}

	name0, _ := widgets[0].(map[string]any)["Name"].(string)
	name1, _ := widgets[1].(map[string]any)["Name"].(string)
	if name0 != "txtName" {
		t.Errorf("Expected first widget 'txtName', got %q", name0)
	}
	if name1 != "txtPhone" {
		t.Errorf("Expected second widget 'txtPhone', got %q", name1)
	}
}

func TestApplyDropWidget_Multiple(t *testing.T) {
	w1 := makeWidget("a", "Pages$TextBox")
	w2 := makeWidget("b", "Pages$TextBox")
	w3 := makeWidget("c", "Pages$TextBox")
	rawData := makeRawPage(w1, w2, w3)

	op := &ast.DropWidgetOp{WidgetNames: []string{"a", "c"}}
	if err := applyDropWidget(rawData, op); err != nil {
		t.Fatalf("applyDropWidget failed: %v", err)
	}

	formCall := rawData["FormCall"].(map[string]any)
	args := getBsonArrayElements(formCall["Arguments"])
	argMap := args[0].(map[string]any)
	widgets := getBsonArrayElements(argMap["Widgets"])

	if len(widgets) != 1 {
		t.Fatalf("Expected 1 widget after dropping a and c, got %d", len(widgets))
	}

	name, _ := widgets[0].(map[string]any)["Name"].(string)
	if name != "b" {
		t.Errorf("Expected remaining widget 'b', got %q", name)
	}
}

func TestApplyDropWidget_NotFound(t *testing.T) {
	w1 := makeWidget("txtName", "Pages$TextBox")
	rawData := makeRawPage(w1)

	op := &ast.DropWidgetOp{WidgetNames: []string{"nonexistent"}}
	err := applyDropWidget(rawData, op)
	if err == nil {
		t.Fatal("Expected error for nonexistent widget")
	}
}

func TestApplyDropWidget_Nested(t *testing.T) {
	inner1 := makeWidget("txtInner1", "Pages$TextBox")
	inner2 := makeWidget("txtInner2", "Pages$TextBox")
	container := makeContainerWidget("ctn1", inner1, inner2)
	rawData := makeRawPage(container)

	op := &ast.DropWidgetOp{WidgetNames: []string{"txtInner1"}}
	if err := applyDropWidget(rawData, op); err != nil {
		t.Fatalf("applyDropWidget failed: %v", err)
	}

	// Verify txtInner1 was removed from container
	result := findBsonWidget(rawData, "txtInner1")
	if result != nil {
		t.Error("txtInner1 should have been removed")
	}

	// txtInner2 should still exist
	result = findBsonWidget(rawData, "txtInner2")
	if result == nil {
		t.Error("txtInner2 should still exist")
	}
}

func TestApplySetProperty_Name(t *testing.T) {
	w1 := makeWidget("txtOld", "Pages$TextBox")
	rawData := makeRawPage(w1)

	op := &ast.SetPropertyOp{
		WidgetName: "txtOld",
		Properties: map[string]interface{}{
			"Name": "txtNew",
		},
	}
	if err := applySetProperty(rawData, op); err != nil {
		t.Fatalf("applySetProperty failed: %v", err)
	}

	// Verify name was changed
	result := findBsonWidget(rawData, "txtNew")
	if result == nil {
		t.Fatal("Expected to find renamed widget 'txtNew'")
	}
}

func TestApplySetProperty_ButtonStyle(t *testing.T) {
	w1 := map[string]any{
		"$Type":       "Pages$ActionButton",
		"Name":        "btnSave",
		"ButtonStyle": "Default",
	}
	rawData := makeRawPage(w1)

	op := &ast.SetPropertyOp{
		WidgetName: "btnSave",
		Properties: map[string]interface{}{
			"ButtonStyle": "Success",
		},
	}
	if err := applySetProperty(rawData, op); err != nil {
		t.Fatalf("applySetProperty failed: %v", err)
	}

	result := findBsonWidget(rawData, "btnSave")
	if result == nil {
		t.Fatal("Expected to find btnSave")
	}
	if result.widget["ButtonStyle"] != "Success" {
		t.Errorf("Expected ButtonStyle='Success', got %v", result.widget["ButtonStyle"])
	}
}

func TestApplySetProperty_WidgetNotFound(t *testing.T) {
	w1 := makeWidget("txtName", "Pages$TextBox")
	rawData := makeRawPage(w1)

	op := &ast.SetPropertyOp{
		WidgetName: "nonexistent",
		Properties: map[string]interface{}{
			"Name": "new",
		},
	}
	err := applySetProperty(rawData, op)
	if err == nil {
		t.Fatal("Expected error for nonexistent widget")
	}
}

func TestApplySetProperty_PluggableWidget(t *testing.T) {
	// Pluggable widget properties are identified by TypePointer referencing
	// a PropertyType entry in Type.ObjectType.PropertyTypes, NOT by a "Key" field.
	propTypeID := primitive.Binary{Subtype: 0x04, Data: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}}
	w1 := map[string]any{
		"$Type": "CustomWidgets$CustomWidget",
		"Name":  "cb1",
		"Type": map[string]any{
			"$Type": "CustomWidgets$CustomWidgetType",
			"ObjectType": map[string]any{
				"PropertyTypes": []any{
					int32(2), // type marker
					map[string]any{
						"$ID":         propTypeID,
						"PropertyKey": "showLabel",
					},
				},
			},
		},
		"Object": map[string]any{
			"Properties": []any{
				int32(2), // type marker
				map[string]any{
					"TypePointer": propTypeID,
					"Value": map[string]any{
						"PrimitiveValue": "yes",
					},
				},
			},
		},
	}
	rawData := makeRawPage(w1)

	op := &ast.SetPropertyOp{
		WidgetName: "cb1",
		Properties: map[string]interface{}{
			"showLabel": false,
		},
	}
	if err := applySetProperty(rawData, op); err != nil {
		t.Fatalf("applySetProperty failed: %v", err)
	}

	result := findBsonWidget(rawData, "cb1")
	if result == nil {
		t.Fatal("Expected to find cb1")
	}
	obj := result.widget["Object"].(map[string]any)
	props := getBsonArrayElements(obj["Properties"])
	propMap := props[0].(map[string]any)
	valMap := propMap["Value"].(map[string]any)
	if valMap["PrimitiveValue"] != "no" {
		t.Errorf("Expected PrimitiveValue='no', got %v", valMap["PrimitiveValue"])
	}
}

func TestSetBsonArray_PreservesMarker(t *testing.T) {
	parent := map[string]any{
		"Widgets": []any{int32(2), "a", "b"},
	}
	setBsonArray(parent, "Widgets", []any{"x", "y"})

	result := parent["Widgets"].([]any)
	if len(result) != 3 {
		t.Fatalf("Expected 3 elements (marker + 2), got %d", len(result))
	}
	if result[0] != int32(2) {
		t.Errorf("Expected marker int32(2), got %v", result[0])
	}
	if result[1] != "x" || result[2] != "y" {
		t.Errorf("Expected [x, y], got %v", result[1:])
	}
}

func TestSetBsonArray_NoMarker(t *testing.T) {
	parent := map[string]any{
		"Widgets": []any{"a", "b"},
	}
	setBsonArray(parent, "Widgets", []any{"x"})

	result := parent["Widgets"].([]any)
	if len(result) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(result))
	}
	if result[0] != "x" {
		t.Errorf("Expected [x], got %v", result)
	}
}

func TestFindBsonWidget_LayoutGrid(t *testing.T) {
	inner := makeWidget("txtInGrid", "Pages$TextBox")
	rawData := map[string]any{
		"FormCall": map[string]any{
			"Arguments": []any{
				int32(2),
				map[string]any{
					"Widgets": []any{
						int32(2),
						map[string]any{
							"$Type": "Pages$LayoutGrid",
							"Name":  "lg1",
							"Rows": []any{
								int32(2),
								map[string]any{
									"Columns": []any{
										int32(2),
										map[string]any{
											"Widgets": []any{int32(2), inner},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	result := findBsonWidget(rawData, "txtInGrid")
	if result == nil {
		t.Fatal("Expected to find txtInGrid inside LayoutGrid")
	}
}

// ============================================================================
// Snippet BSON tests
// ============================================================================

// Helper to build a minimal raw BSON snippet structure (Studio Pro format).
func makeRawSnippet(widgets ...map[string]any) map[string]any {
	widgetArr := []any{int32(2)} // type marker
	for _, w := range widgets {
		widgetArr = append(widgetArr, w)
	}
	return map[string]any{
		"Widgets": widgetArr,
	}
}

// Helper to build a minimal raw BSON snippet structure (mxcli format).
func makeRawSnippetMxcli(widgets ...map[string]any) map[string]any {
	widgetArr := []any{int32(2)} // type marker
	for _, w := range widgets {
		widgetArr = append(widgetArr, w)
	}
	return map[string]any{
		"Widget": map[string]any{
			"Widgets": widgetArr,
		},
	}
}

func TestFindBsonWidgetInSnippet_TopLevel(t *testing.T) {
	w1 := makeWidget("txtName", "Pages$TextBox")
	w2 := makeWidget("txtEmail", "Pages$TextBox")
	rawData := makeRawSnippet(w1, w2)

	result := findBsonWidgetInSnippet(rawData, "txtName")
	if result == nil {
		t.Fatal("Expected to find txtName in snippet")
	}
	name, _ := result.widget["Name"].(string)
	if name != "txtName" {
		t.Errorf("Expected 'txtName', got %q", name)
	}
}

func TestFindBsonWidgetInSnippet_MxcliFormat(t *testing.T) {
	w1 := makeWidget("txtName", "Pages$TextBox")
	rawData := makeRawSnippetMxcli(w1)

	result := findBsonWidgetInSnippet(rawData, "txtName")
	if result == nil {
		t.Fatal("Expected to find txtName in mxcli-format snippet")
	}
}

func TestFindBsonWidgetInSnippet_Nested(t *testing.T) {
	inner := makeWidget("txtInner", "Pages$TextBox")
	container := makeContainerWidget("ctn1", inner)
	rawData := makeRawSnippet(container)

	result := findBsonWidgetInSnippet(rawData, "txtInner")
	if result == nil {
		t.Fatal("Expected to find txtInner nested in snippet")
	}
}

func TestFindBsonWidgetInSnippet_NotFound(t *testing.T) {
	w1 := makeWidget("txtName", "Pages$TextBox")
	rawData := makeRawSnippet(w1)

	result := findBsonWidgetInSnippet(rawData, "nonexistent")
	if result != nil {
		t.Error("Expected nil for nonexistent widget in snippet")
	}
}

func TestApplyDropWidget_Snippet(t *testing.T) {
	w1 := makeWidget("txtName", "Pages$TextBox")
	w2 := makeWidget("txtEmail", "Pages$TextBox")
	rawData := makeRawSnippet(w1, w2)

	op := &ast.DropWidgetOp{WidgetNames: []string{"txtEmail"}}
	if err := applyDropWidgetWith(rawData, op, findBsonWidgetInSnippet); err != nil {
		t.Fatalf("applyDropWidgetWith failed: %v", err)
	}

	// Verify txtEmail was removed
	widgets := getBsonArrayElements(rawData["Widgets"])
	if len(widgets) != 1 {
		t.Fatalf("Expected 1 widget after drop, got %d", len(widgets))
	}
	name, _ := widgets[0].(map[string]any)["Name"].(string)
	if name != "txtName" {
		t.Errorf("Expected remaining widget 'txtName', got %q", name)
	}
}

func TestApplySetProperty_Snippet(t *testing.T) {
	w1 := map[string]any{
		"$Type":       "Pages$ActionButton",
		"Name":        "btnAction",
		"ButtonStyle": "Default",
	}
	rawData := makeRawSnippet(w1)

	op := &ast.SetPropertyOp{
		WidgetName: "btnAction",
		Properties: map[string]interface{}{
			"ButtonStyle": "Danger",
		},
	}
	if err := applySetPropertyWith(rawData, op, findBsonWidgetInSnippet); err != nil {
		t.Fatalf("applySetPropertyWith failed: %v", err)
	}

	result := findBsonWidgetInSnippet(rawData, "btnAction")
	if result == nil {
		t.Fatal("Expected to find btnAction")
	}
	if result.widget["ButtonStyle"] != "Danger" {
		t.Errorf("Expected ButtonStyle='Danger', got %v", result.widget["ButtonStyle"])
	}
}

func TestFindBsonWidget_DataViewFooter(t *testing.T) {
	footer := makeWidget("btnFooter", "Pages$ActionButton")
	rawData := map[string]any{
		"FormCall": map[string]any{
			"Arguments": []any{
				int32(2),
				map[string]any{
					"Widgets": []any{
						int32(2),
						map[string]any{
							"$Type":         "Pages$DataView",
							"Name":          "dv1",
							"Widgets":       []any{int32(2)},
							"FooterWidgets": []any{int32(2), footer},
						},
					},
				},
			},
		},
	}

	result := findBsonWidget(rawData, "btnFooter")
	if result == nil {
		t.Fatal("Expected to find btnFooter in DataView FooterWidgets")
	}
}
