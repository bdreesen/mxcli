// SPDX-License-Identifier: Apache-2.0

// Package mpr - Reader methods for agent-editor CustomBlobDocuments.
//
// Covers the four document types created by the Studio Pro Agent Editor
// extension: Agent, Model, Knowledge Base, Consumed MCP Service. Each
// shares the outer CustomBlobDocument BSON wrapper and is discriminated
// by CustomDocumentType. This file currently implements Model only; the
// other three will follow the same pattern.
package mpr

import (
	"fmt"

	"github.com/mendixlabs/mxcli/sdk/agenteditor"
)

// ListAgentEditorModels returns all agent-editor Model documents in the
// project (CustomDocumentType == "agenteditor.model").
func (r *Reader) ListAgentEditorModels() ([]*agenteditor.Model, error) {
	units, err := r.listUnitsByType(customBlobDocType)
	if err != nil {
		return nil, err
	}

	var result []*agenteditor.Model
	for _, u := range units {
		wrap, err := parseCustomBlobWrapper(u.Contents)
		if err != nil {
			// Skip units we can't decode; log to error list if useful later.
			continue
		}
		if wrap.CustomDocumentType != agenteditor.CustomTypeModel {
			continue
		}
		m, err := r.parseAgentEditorModel(u.ID, u.ContainerID, u.Contents)
		if err != nil {
			return nil, fmt.Errorf("failed to parse agent-editor model %s: %w", u.ID, err)
		}
		result = append(result, m)
	}
	return result, nil
}
