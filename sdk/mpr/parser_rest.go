// SPDX-License-Identifier: Apache-2.0

package mpr

import (
	"fmt"

	"github.com/mendixlabs/mxcli/model"

	"go.mongodb.org/mongo-driver/bson"
)

// parsePublishedRestService parses a published REST service from BSON.
func (r *Reader) parsePublishedRestService(unitID, containerID string, contents []byte) (*model.PublishedRestService, error) {
	contents, err := r.resolveContents(unitID, contents)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := bson.Unmarshal(contents, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal BSON: %w", err)
	}

	svc := &model.PublishedRestService{}
	svc.ID = model.ID(unitID)
	svc.TypeName = "Rest$PublishedRestService"
	svc.ContainerID = model.ID(containerID)

	svc.Name = extractString(raw["Name"])
	svc.Path = extractString(raw["Path"])
	svc.Version = extractString(raw["Version"])
	svc.ServiceName = extractString(raw["ServiceName"])
	svc.Excluded = extractBool(raw["Excluded"], false)

	// Parse resources
	resources := extractBsonArray(raw["Resources"])
	for _, res := range resources {
		if resMap, ok := res.(map[string]any); ok {
			resource := &model.PublishedRestResource{}
			resource.ID = model.ID(extractBsonID(resMap["$ID"]))
			resource.TypeName = extractString(resMap["$Type"])
			resource.Name = extractString(resMap["Name"])

			// Parse operations
			ops := extractBsonArray(resMap["Operations"])
			for _, op := range ops {
				if opMap, ok := op.(map[string]any); ok {
					operation := &model.PublishedRestOperation{}
					operation.ID = model.ID(extractBsonID(opMap["$ID"]))
					operation.TypeName = extractString(opMap["$Type"])
					operation.Path = extractString(opMap["Path"])
					operation.HTTPMethod = extractString(opMap["HttpMethod"])
					operation.Summary = extractString(opMap["Summary"])
					operation.Microflow = extractString(opMap["Microflow"])
					operation.Deprecated = extractBool(opMap["Deprecated"], false)
					resource.Operations = append(resource.Operations, operation)
				}
			}

			svc.Resources = append(svc.Resources, resource)
		}
	}

	return svc, nil
}
