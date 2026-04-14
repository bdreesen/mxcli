// SPDX-License-Identifier: Apache-2.0

// Package mpr - Parsing of CustomBlobDocuments$CustomBlobDocument units.
//
// The agent-editor Studio Pro extension (Mendix 11.9+) stores all of its
// documents — Agent, Model, Knowledge Base, Consumed MCP Service — as
// generic CustomBlobDocument units. They share the same BSON wrapper and
// are discriminated by the CustomDocumentType field. The actual document
// payload lives in a JSON string in the Contents field.
//
// This file provides the generic wrapper decode plus type-specific
// decoders for each inner JSON schema.
package mpr

import (
	"encoding/json"
	"fmt"

	"github.com/mendixlabs/mxcli/model"
	"github.com/mendixlabs/mxcli/sdk/agenteditor"

	"go.mongodb.org/mongo-driver/bson"
)

// customBlobDocType is the BSON $Type of the wrapper.
const customBlobDocType = "CustomBlobDocuments$CustomBlobDocument"

// rawCustomBlobDoc is the decoded BSON wrapper (fields we care about).
type rawCustomBlobDoc struct {
	Name               string
	Documentation      string
	Excluded           bool
	ExportLevel        string
	CustomDocumentType string
	Contents           string // JSON payload
}

// parseCustomBlobWrapper decodes the outer CustomBlobDocument BSON wrapper.
// Returns a rawCustomBlobDoc or an error.
func parseCustomBlobWrapper(contents []byte) (*rawCustomBlobDoc, error) {
	var raw map[string]any
	if err := bson.Unmarshal(contents, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CustomBlobDocument BSON: %w", err)
	}

	out := &rawCustomBlobDoc{}
	if v, ok := raw["Name"].(string); ok {
		out.Name = v
	}
	if v, ok := raw["Documentation"].(string); ok {
		out.Documentation = v
	}
	if v, ok := raw["Excluded"].(bool); ok {
		out.Excluded = v
	}
	if v, ok := raw["ExportLevel"].(string); ok {
		out.ExportLevel = v
	}
	if v, ok := raw["CustomDocumentType"].(string); ok {
		out.CustomDocumentType = v
	}
	if v, ok := raw["Contents"].(string); ok {
		out.Contents = v
	}
	return out, nil
}

// parseAgentEditorModel parses a CustomBlobDocument with
// CustomDocumentType == "agenteditor.model" into an agenteditor.Model.
func (r *Reader) parseAgentEditorModel(unitID, containerID string, contents []byte) (*agenteditor.Model, error) {
	contents, err := r.resolveContents(unitID, contents)
	if err != nil {
		return nil, err
	}

	wrap, err := parseCustomBlobWrapper(contents)
	if err != nil {
		return nil, err
	}
	if wrap.CustomDocumentType != agenteditor.CustomTypeModel {
		return nil, fmt.Errorf("unit %s is not an agent-editor model (CustomDocumentType=%q)",
			unitID, wrap.CustomDocumentType)
	}

	m := &agenteditor.Model{}
	m.ID = model.ID(unitID)
	m.TypeName = customBlobDocType
	m.ContainerID = model.ID(containerID)
	m.Name = wrap.Name
	m.Documentation = wrap.Documentation
	m.Excluded = wrap.Excluded
	m.ExportLevel = wrap.ExportLevel

	// Decode the Contents JSON payload.
	if wrap.Contents != "" {
		var payload struct {
			Type           string `json:"type"`
			Name           string `json:"name"`
			DisplayName    string `json:"displayName"`
			Provider       string `json:"provider"`
			ProviderFields struct {
				Environment  string                    `json:"environment"`
				DeepLinkURL  string                    `json:"deepLinkURL"`
				KeyID        string                    `json:"keyId"`
				KeyName      string                    `json:"keyName"`
				ResourceName string                    `json:"resourceName"`
				Key          *agenteditor.ConstantRef  `json:"key"`
			} `json:"providerFields"`
		}
		if err := json.Unmarshal([]byte(wrap.Contents), &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal agent-editor Model Contents JSON: %w", err)
		}

		m.Type = payload.Type
		m.InnerName = payload.Name
		m.DisplayName = payload.DisplayName
		m.Provider = payload.Provider
		m.Environment = payload.ProviderFields.Environment
		m.DeepLinkURL = payload.ProviderFields.DeepLinkURL
		m.KeyID = payload.ProviderFields.KeyID
		m.KeyName = payload.ProviderFields.KeyName
		m.ResourceName = payload.ProviderFields.ResourceName
		m.Key = payload.ProviderFields.Key
	}

	return m, nil
}
