// SPDX-License-Identifier: Apache-2.0

package mpr

import (
	"testing"

	"github.com/mendixlabs/mxcli/model"
	"github.com/mendixlabs/mxcli/sdk/microflows"
	"go.mongodb.org/mongo-driver/bson"
)

// TestShowPageAction_Roundtrip verifies that a ShowPageAction with parameters
// survives BSON serialization/deserialization.
func TestShowPageAction_Roundtrip(t *testing.T) {
	// Build a ShowPageAction with parameter mappings
	action := &microflows.ShowPageAction{
		BaseElement: model.BaseElement{ID: "test-action-id"},
		PageName:    "Sales.Product_NewEdit",
		PageParameterMappings: []*microflows.PageParameterMapping{
			{
				BaseElement: model.BaseElement{ID: "test-mapping-id"},
				Parameter:   "Sales.Product_NewEdit.Product",
				Argument:    "$Product",
			},
		},
	}

	// Serialize to BSON using the writer
	doc := serializeMicroflowAction(action)

	// Marshal to bytes (simulates writing to MPR)
	data, err := bson.Marshal(doc)
	if err != nil {
		t.Fatalf("failed to marshal BSON: %v", err)
	}

	// Unmarshal back to map (simulates reading from MPR)
	var raw map[string]any
	if err := bson.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal BSON: %v", err)
	}

	// Parse using the parser
	parsed := parseShowPageAction(raw)

	// Verify page name
	if parsed.PageName != "Sales.Product_NewEdit" {
		t.Errorf("PageName = %q, want %q", parsed.PageName, "Sales.Product_NewEdit")
	}

	// Verify parameter mappings
	if len(parsed.PageParameterMappings) != 1 {
		t.Fatalf("PageParameterMappings count = %d, want 1", len(parsed.PageParameterMappings))
	}
	pm := parsed.PageParameterMappings[0]
	if pm.Parameter != "Sales.Product_NewEdit.Product" {
		t.Errorf("Parameter = %q, want %q", pm.Parameter, "Sales.Product_NewEdit.Product")
	}
	if pm.Argument != "$Product" {
		t.Errorf("Argument = %q, want %q", pm.Argument, "$Product")
	}
}

// TestShowPageAction_RoundtripNoParams verifies that a ShowPageAction without parameters
// survives BSON serialization/deserialization.
func TestShowPageAction_RoundtripNoParams(t *testing.T) {
	action := &microflows.ShowPageAction{
		BaseElement: model.BaseElement{ID: "test-action-id"},
		PageName:    "Sales.Customer_Overview",
	}

	doc := serializeMicroflowAction(action)
	data, err := bson.Marshal(doc)
	if err != nil {
		t.Fatalf("failed to marshal BSON: %v", err)
	}

	var raw map[string]any
	if err := bson.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal BSON: %v", err)
	}

	parsed := parseShowPageAction(raw)

	if parsed.PageName != "Sales.Customer_Overview" {
		t.Errorf("PageName = %q, want %q", parsed.PageName, "Sales.Customer_Overview")
	}
	if len(parsed.PageParameterMappings) != 0 {
		t.Errorf("PageParameterMappings count = %d, want 0", len(parsed.PageParameterMappings))
	}
}

// TestShowPageAction_RoundtripMultipleParams verifies multiple parameter mappings survive roundtrip.
func TestShowPageAction_RoundtripMultipleParams(t *testing.T) {
	action := &microflows.ShowPageAction{
		BaseElement: model.BaseElement{ID: "test-action-id"},
		PageName:    "Sales.Order_Detail",
		PageParameterMappings: []*microflows.PageParameterMapping{
			{
				BaseElement: model.BaseElement{ID: "mapping-1"},
				Parameter:   "Sales.Order_Detail.Order",
				Argument:    "$Order",
			},
			{
				BaseElement: model.BaseElement{ID: "mapping-2"},
				Parameter:   "Sales.Order_Detail.Customer",
				Argument:    "$Customer",
			},
		},
	}

	doc := serializeMicroflowAction(action)
	data, err := bson.Marshal(doc)
	if err != nil {
		t.Fatalf("failed to marshal BSON: %v", err)
	}

	var raw map[string]any
	if err := bson.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal BSON: %v", err)
	}

	parsed := parseShowPageAction(raw)

	if parsed.PageName != "Sales.Order_Detail" {
		t.Errorf("PageName = %q, want %q", parsed.PageName, "Sales.Order_Detail")
	}
	if len(parsed.PageParameterMappings) != 2 {
		t.Fatalf("PageParameterMappings count = %d, want 2", len(parsed.PageParameterMappings))
	}
	if parsed.PageParameterMappings[0].Argument != "$Order" {
		t.Errorf("first Argument = %q, want %q", parsed.PageParameterMappings[0].Argument, "$Order")
	}
	if parsed.PageParameterMappings[1].Argument != "$Customer" {
		t.Errorf("second Argument = %q, want %q", parsed.PageParameterMappings[1].Argument, "$Customer")
	}
}
