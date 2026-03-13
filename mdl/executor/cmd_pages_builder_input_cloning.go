// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"github.com/mendixlabs/mxcli/sdk/mpr"
	"github.com/mendixlabs/mxcli/sdk/pages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// cloneDataGrid2ObjectWithDatasourceOnly clones a template Object, only updating the datasource.
// This is for testing to isolate whether column building is the issue.
func (pb *pageBuilder) cloneDataGrid2ObjectWithDatasourceOnly(templateObject bson.D, propertyTypeIDs map[string]pages.PropertyTypeIDEntry, datasource pages.DataSource) bson.D {
	result := make(bson.D, 0, len(templateObject))

	for _, elem := range templateObject {
		if elem.Key == "$ID" {
			// Generate new ID for the object
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "Properties" {
			// Update only datasource property
			if propsArr, ok := elem.Value.(bson.A); ok {
				updatedProps := pb.updateOnlyDatasource(propsArr, propertyTypeIDs, datasource)
				result = append(result, bson.E{Key: "Properties", Value: updatedProps})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}

	return result
}

// updateOnlyDatasource only updates the datasource property, keeping everything else as-is.
func (pb *pageBuilder) updateOnlyDatasource(props bson.A, propertyTypeIDs map[string]pages.PropertyTypeIDEntry, datasource pages.DataSource) bson.A {
	result := bson.A{int32(2)} // Version marker
	datasourceEntry := propertyTypeIDs["datasource"]

	for _, propVal := range props {
		if _, ok := propVal.(int32); ok {
			continue // Skip version markers
		}
		propMap, ok := propVal.(bson.D)
		if !ok {
			continue
		}

		typePointer := pb.getTypePointerFromProperty(propMap)
		if typePointer == datasourceEntry.PropertyTypeID {
			// Replace datasource
			result = append(result, pb.buildDataGrid2Property(datasourceEntry, datasource, "", ""))
		} else {
			// Keep as-is but with new IDs
			result = append(result, pb.clonePropertyWithNewIDs(propMap))
		}
	}

	return result
}

// getTypePointerFromProperty extracts the TypePointer ID from a WidgetProperty.
func (pb *pageBuilder) getTypePointerFromProperty(prop bson.D) string {
	for _, elem := range prop {
		if elem.Key == "TypePointer" {
			switch v := elem.Value.(type) {
			case primitive.Binary:
				return mpr.BsonBinaryToID(v)
			case []byte:
				// When loaded from JSON template, binary is []byte instead of primitive.Binary
				return mpr.BlobToUUID(v)
			}
		}
	}
	return ""
}

// clonePropertyWithNewIDs clones a WidgetProperty with new IDs.
func (pb *pageBuilder) clonePropertyWithNewIDs(prop bson.D) bson.D {
	result := make(bson.D, 0, len(prop))
	for _, elem := range prop {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "Value" {
			if valMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "Value", Value: pb.cloneValueWithNewIDs(valMap)})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// cloneValueWithNewIDs clones a WidgetValue with new IDs.
func (pb *pageBuilder) cloneValueWithNewIDs(val bson.D) bson.D {
	result := make(bson.D, 0, len(val))
	for _, elem := range val {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "Action" {
			// Clone nested action with new ID
			if actionMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "Action", Value: pb.cloneWithNewID(actionMap)})
			} else {
				result = append(result, elem)
			}
		} else if elem.Key == "TextTemplate" {
			// Clone nested TextTemplate with new IDs
			if ttMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "TextTemplate", Value: pb.cloneTextTemplateWithNewIDs(ttMap)})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// clonePropertyWithPrimitiveValue clones a WidgetProperty with new IDs and an updated PrimitiveValue.
// This preserves the template's exact structure (TextTemplate, Objects, etc.) while only changing the value.
func (pb *pageBuilder) clonePropertyWithPrimitiveValue(prop bson.D, newValue string) bson.D {
	result := make(bson.D, 0, len(prop))
	for _, elem := range prop {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "Value" {
			if valMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "Value", Value: pb.cloneValueWithUpdatedPrimitive(valMap, newValue)})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// cloneValueWithUpdatedPrimitive clones a WidgetValue with new IDs and an updated PrimitiveValue.
func (pb *pageBuilder) cloneValueWithUpdatedPrimitive(val bson.D, newValue string) bson.D {
	result := make(bson.D, 0, len(val))
	for _, elem := range val {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "PrimitiveValue" {
			result = append(result, bson.E{Key: "PrimitiveValue", Value: newValue})
		} else if elem.Key == "Action" {
			if actionMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "Action", Value: pb.cloneWithNewID(actionMap)})
			} else {
				result = append(result, elem)
			}
		} else if elem.Key == "TextTemplate" {
			if ttMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "TextTemplate", Value: pb.cloneTextTemplateWithNewIDs(ttMap)})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// clonePropertyClearingTextTemplate clones a WidgetProperty with new IDs but sets TextTemplate to nil.
// Used for mode-dependent properties where TextTemplate should not be present.
func (pb *pageBuilder) clonePropertyClearingTextTemplate(prop bson.D) bson.D {
	result := make(bson.D, 0, len(prop))
	for _, elem := range prop {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "Value" {
			if valMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "Value", Value: pb.cloneValueClearingTextTemplate(valMap)})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// cloneValueClearingTextTemplate clones a WidgetValue with new IDs and TextTemplate set to nil.
func (pb *pageBuilder) cloneValueClearingTextTemplate(val bson.D) bson.D {
	result := make(bson.D, 0, len(val))
	for _, elem := range val {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "TextTemplate" {
			result = append(result, bson.E{Key: "TextTemplate", Value: nil})
		} else if elem.Key == "Action" {
			if actionMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "Action", Value: pb.cloneWithNewID(actionMap)})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// cloneWithNewID clones a BSON document replacing only the $ID field.
func (pb *pageBuilder) cloneWithNewID(doc bson.D) bson.D {
	result := make(bson.D, 0, len(doc))
	for _, elem := range doc {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// cloneTextTemplateWithNewIDs clones a Forms$ClientTemplate with new IDs.
func (pb *pageBuilder) cloneTextTemplateWithNewIDs(tt bson.D) bson.D {
	result := make(bson.D, 0, len(tt))
	for _, elem := range tt {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "Fallback" || elem.Key == "Template" {
			if textMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: elem.Key, Value: pb.cloneWithNewID(textMap)})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// clonePropertyWithExpression clones a WidgetProperty with new IDs and an updated Expression.
// Same as clonePropertyWithPrimitiveValue but replaces the Expression field instead.
func (pb *pageBuilder) clonePropertyWithExpression(prop bson.D, newExpr string) bson.D {
	result := make(bson.D, 0, len(prop))
	for _, elem := range prop {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "Value" {
			if valMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "Value", Value: pb.cloneValueWithUpdatedExpression(valMap, newExpr)})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}
	return result
}

// cloneValueWithUpdatedExpression clones a WidgetValue with new IDs and an updated Expression.
func (pb *pageBuilder) cloneValueWithUpdatedExpression(val bson.D, newExpr string) bson.D {
	result := make(bson.D, 0, len(val))
	for _, elem := range val {
		if elem.Key == "$ID" {
			result = append(result, bson.E{Key: "$ID", Value: mpr.IDToBsonBinary(mpr.GenerateID())})
		} else if elem.Key == "Expression" {
			result = append(result, bson.E{Key: "Expression", Value: newExpr})
		} else if elem.Key == "Action" {
			if actionMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "Action", Value: pb.cloneWithNewID(actionMap)})
			} else {
				result = append(result, elem)
			}
		} else if elem.Key == "TextTemplate" {
			if ttMap, ok := elem.Value.(bson.D); ok {
				result = append(result, bson.E{Key: "TextTemplate", Value: pb.cloneTextTemplateWithNewIDs(ttMap)})
			} else {
				result = append(result, elem)
			}
		} else {
			result = append(result, elem)
		}
	}
	return result
}
