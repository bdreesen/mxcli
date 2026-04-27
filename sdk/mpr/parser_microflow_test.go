// SPDX-License-Identifier: Apache-2.0

package mpr

import (
	"testing"

	"github.com/mendixlabs/mxcli/model"
	"github.com/mendixlabs/mxcli/sdk/microflows"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestParseSequenceFlow_NewCaseValueEnumerationCase(t *testing.T) {
	flow := parseSequenceFlow(map[string]any{
		"$ID":                        "flow-1",
		"OriginPointer":              "start-1",
		"DestinationPointer":         "dest-1",
		"OriginConnectionIndex":      int32(1),
		"DestinationConnectionIndex": int32(2),
		"NewCaseValue": primitive.D{
			{Key: "$ID", Value: "case-1"},
			{Key: "$Type", Value: "Microflows$EnumerationCase"},
			{Key: "Value", Value: "true"},
		},
	})

	got, ok := flow.CaseValue.(*microflows.EnumerationCase)
	if !ok {
		t.Fatalf("expected *EnumerationCase, got %T", flow.CaseValue)
	}
	if got.Value != "true" {
		t.Fatalf("expected true branch, got %q", got.Value)
	}
}

func TestParseSequenceFlow_NewCaseValueNoCase(t *testing.T) {
	flow := parseSequenceFlow(map[string]any{
		"$ID":                "flow-1",
		"OriginPointer":      "start-1",
		"DestinationPointer": "dest-1",
		"NewCaseValue": primitive.D{
			{Key: "$ID", Value: "case-1"},
			{Key: "$Type", Value: "Microflows$NoCase"},
		},
	})

	if _, ok := flow.CaseValue.(*microflows.NoCase); !ok {
		t.Fatalf("expected *NoCase, got %T", flow.CaseValue)
	}
}

func TestParseCommitAction_ErrorHandlingTypeExplicit(t *testing.T) {
	action := parseCommitAction(map[string]any{
		"$ID":                "commit-1",
		"CommitVariableName": "Order",
		"WithEvents":         true,
		"RefreshInClient":    false,
		"ErrorHandlingType":  "Continue",
	})

	if action.ErrorHandlingType != microflows.ErrorHandlingTypeContinue {
		t.Errorf("expected Continue, got %q", action.ErrorHandlingType)
	}
	if action.CommitVariable != "Order" {
		t.Errorf("expected CommitVariable Order, got %q", action.CommitVariable)
	}
}

func TestParseCommitAction_ErrorHandlingTypeDefaultsToRollback(t *testing.T) {
	// When ErrorHandlingType is absent from BSON, the describer must still
	// emit "on error rollback" — matching Mendix Studio Pro's default.
	// Without this default, describe → exec → describe drops the suffix
	// because the writer omits the field when it equals Rollback.
	action := parseCommitAction(map[string]any{
		"$ID":                "commit-1",
		"CommitVariableName": "Order",
		"WithEvents":         false,
		"RefreshInClient":    false,
	})

	if action.ErrorHandlingType != microflows.ErrorHandlingTypeRollback {
		t.Errorf("expected default Rollback, got %q", action.ErrorHandlingType)
	}
}

func TestParseResultHandlingMappingUsesRangeForSingleObject(t *testing.T) {
	got := parseResultHandling(map[string]any{
		"$ID":                "result-handling-1",
		"ResultVariableName": "RemoteApp",
		"ImportMappingCall": map[string]any{
			"ReturnValueMapping":    "SampleRuntimeApi.IMM_RemoteApp",
			"ForceSingleOccurrence": false,
			"Range": map[string]any{
				"SingleObject": true,
			},
		},
		"VariableType": map[string]any{
			"$Type":  "DataTypes$ObjectType",
			"Entity": "SampleRuntimeApi.RemoteApp",
		},
	}, "Mapping")

	rh, ok := got.(*microflows.ResultHandlingMapping)
	if !ok {
		t.Fatalf("got %T, want *microflows.ResultHandlingMapping", got)
	}
	if !rh.SingleObject {
		t.Fatal("Range.SingleObject=true must make the result object-valued")
	}
	if rh.ForceSingleOccurrence == nil || *rh.ForceSingleOccurrence {
		t.Fatalf("ForceSingleOccurrence = %v, want explicit false", rh.ForceSingleOccurrence)
	}
}

func TestSerializeRestResultHandlingPreservesForceSingleOccurrenceSeparately(t *testing.T) {
	forceSingleOccurrence := false
	doc := serializeRestResultHandling(&microflows.ResultHandlingMapping{
		BaseElement:           model.BaseElement{ID: model.ID("result-handling-1")},
		MappingID:             model.ID("SampleRuntimeApi.IMM_RemoteApp"),
		ResultEntityID:        model.ID("SampleRuntimeApi.RemoteApp"),
		ResultVariable:        "RemoteApp",
		SingleObject:          true,
		ForceSingleOccurrence: &forceSingleOccurrence,
	}, "RemoteApp")

	importCall, ok := bsonDMap(doc)["ImportMappingCall"].(primitive.D)
	if !ok {
		t.Fatalf("ImportMappingCall missing or wrong type: %T", bsonDMap(doc)["ImportMappingCall"])
	}
	callFields := bsonDMap(importCall)
	if got := callFields["ForceSingleOccurrence"]; got != false {
		t.Fatalf("ForceSingleOccurrence = %v, want false", got)
	}
	rangeDoc, ok := callFields["Range"].(primitive.D)
	if !ok {
		t.Fatalf("Range missing or wrong type: %T", callFields["Range"])
	}
	if got := bsonDMap(rangeDoc)["SingleObject"]; got != true {
		t.Fatalf("Range.SingleObject = %v, want true", got)
	}
	varType, ok := bsonDMap(doc)["VariableType"].(primitive.D)
	if !ok {
		t.Fatalf("VariableType missing or wrong type: %T", bsonDMap(doc)["VariableType"])
	}
	if got := bsonDMap(varType)["$Type"]; got != "DataTypes$ObjectType" {
		t.Fatalf("VariableType.$Type = %v, want DataTypes$ObjectType", got)
	}
}

func TestSerializeImportXmlActionPreservesSingleObjectRange(t *testing.T) {
	forceSingleOccurrence := false
	doc := serializeImportXmlAction(&microflows.ImportXmlAction{
		BaseElement: model.BaseElement{ID: model.ID("import-action-1")},
		ResultHandling: &microflows.ResultHandlingMapping{
			BaseElement:           model.BaseElement{ID: model.ID("result-handling-1")},
			MappingID:             model.ID("SampleRest.IMM_ErrorResponse"),
			ResultEntityID:        model.ID("SampleRest.Error"),
			ResultVariable:        "ErrorResponse",
			SingleObject:          true,
			ForceSingleOccurrence: &forceSingleOccurrence,
		},
		XmlDocumentVariable: "LatestHttpResponse",
	})

	resultHandling, ok := bsonDMap(doc)["ResultHandling"].(primitive.D)
	if !ok {
		t.Fatalf("ResultHandling missing or wrong type: %T", bsonDMap(doc)["ResultHandling"])
	}
	importCall, ok := bsonDMap(resultHandling)["ImportMappingCall"].(primitive.D)
	if !ok {
		t.Fatalf("ImportMappingCall missing or wrong type: %T", bsonDMap(resultHandling)["ImportMappingCall"])
	}
	callFields := bsonDMap(importCall)
	if got := callFields["ForceSingleOccurrence"]; got != false {
		t.Fatalf("ForceSingleOccurrence = %v, want false", got)
	}
	rangeDoc, ok := callFields["Range"].(primitive.D)
	if !ok {
		t.Fatalf("Range missing or wrong type: %T", callFields["Range"])
	}
	if got := bsonDMap(rangeDoc)["SingleObject"]; got != true {
		t.Fatalf("Range.SingleObject = %v, want true", got)
	}
}

func bsonDMap(doc primitive.D) map[string]any {
	out := make(map[string]any, len(doc))
	for _, elem := range doc {
		out[elem.Key] = elem.Value
	}
	return out
}
