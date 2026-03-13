// SPDX-License-Identifier: Apache-2.0

package mpr

import (
	"fmt"

	"github.com/mendixlabs/mxcli/model"
	"github.com/mendixlabs/mxcli/sdk/microflows"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *Reader) parseNanoflow(unitID, containerID string, contents []byte) (*microflows.Nanoflow, error) {
	contents, err := r.resolveContents(unitID, contents)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := bson.Unmarshal(contents, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal BSON: %w", err)
	}

	nf := &microflows.Nanoflow{}
	nf.ID = model.ID(unitID)
	nf.TypeName = "Microflows$Nanoflow"
	nf.ContainerID = model.ID(containerID)

	if name, ok := raw["Name"].(string); ok {
		nf.Name = name
	}
	if doc, ok := raw["Documentation"].(string); ok {
		nf.Documentation = doc
	}
	if markAsUsed, ok := raw["MarkAsUsed"].(bool); ok {
		nf.MarkAsUsed = markAsUsed
	}
	if excluded, ok := raw["Excluded"].(bool); ok {
		nf.Excluded = excluded
	}

	// Parse parameters
	if params, ok := raw["Parameters"].([]any); ok {
		for _, p := range params {
			if paramMap, ok := p.(map[string]any); ok {
				param := parseMicroflowParameter(paramMap)
				nf.Parameters = append(nf.Parameters, param)
			}
		}
	}

	return nf, nil
}

// parsePage parses page contents from BSON.
