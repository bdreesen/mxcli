// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mendixlabs/mxcli/mdl/ast"
	"github.com/mendixlabs/mxcli/model"
	"github.com/mendixlabs/mxcli/sdk/domainmodel"
)

// domainModelELKData is the JSON output schema for the domain model ELK diagram.
type domainModelELKData struct {
	Format          string                         `json:"format"`
	Type            string                         `json:"type"`
	ModuleName      string                         `json:"moduleName"`
	FocusEntity     string                         `json:"focusEntity,omitempty"`
	Entities        []domainModelELKEntity         `json:"entities"`
	Associations    []domainModelELKAssoc          `json:"associations"`
	Generalizations []domainModelELKGeneralization `json:"generalizations"`
	MdlSource       string                         `json:"mdlSource,omitempty"`
	SourceMap       map[string]elkSourceRange      `json:"sourceMap,omitempty"`
}

type domainModelELKEntity struct {
	ID         string                    `json:"id"`
	Name       string                    `json:"name"`
	Category   string                    `json:"category"`
	IsFocus    bool                      `json:"isFocus,omitempty"`
	Attributes []domainModelELKAttribute `json:"attributes"`
	Width      float64                   `json:"width"`
	Height     float64                   `json:"height"`
}

type domainModelELKAttribute struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type domainModelELKAssoc struct {
	ID       string `json:"id"`
	SourceID string `json:"sourceId"`
	TargetID string `json:"targetId"`
	Name     string `json:"name"`
	Type     string `json:"type"` // "reference" or "referenceSet"
}

type domainModelELKGeneralization struct {
	ChildID    string `json:"childId"`
	ParentID   string `json:"parentId"`
	ParentName string `json:"parentName"`
}

// Sizing constants for ELK node dimension calculation.
const (
	elkCharWidth      = 7.5
	elkHeaderHeight   = 28.0
	elkAttrLineHeight = 18.0
	elkHPadding       = 24.0
	elkMinWidth       = 100.0
)

// DomainModelELK generates a JSON graph of a module's domain model for rendering with ELK.js.
// If name contains a dot (e.g. "Module.Entity"), it delegates to EntityFocusELK for a focused view.
func (e *Executor) DomainModelELK(name string) error {
	if e.reader == nil {
		return fmt.Errorf("not connected to a project")
	}

	// If name is qualified (Module.Entity), render focused entity diagram
	if strings.Contains(name, ".") {
		return e.EntityFocusELK(name)
	}

	moduleName := name
	module, err := e.findModule(moduleName)
	if err != nil {
		return err
	}

	dm, err := e.reader.GetDomainModel(module.ID)
	if err != nil {
		return fmt.Errorf("failed to get domain model: %w", err)
	}

	allEntityNames, _ := e.buildAllEntityNames()

	// Track which entity IDs are in the current module
	moduleEntityIDs := make(map[model.ID]bool)
	for _, entity := range dm.Entities {
		moduleEntityIDs[entity.ID] = true
	}

	ghostEntities := make(map[string]*domainModelELKEntity)

	// Build entity nodes
	var entities []domainModelELKEntity
	for _, entity := range dm.Entities {
		entities = append(entities, buildELKEntity(entity))
	}

	// Build associations
	var associations []domainModelELKAssoc
	for i, assoc := range dm.Associations {
		addGhostIfNeeded(assoc.ParentID, moduleEntityIDs, allEntityNames, ghostEntities)
		addGhostIfNeeded(assoc.ChildID, moduleEntityIDs, allEntityNames, ghostEntities)

		associations = append(associations, domainModelELKAssoc{
			ID:       fmt.Sprintf("assoc-%d", i),
			SourceID: "entity-" + string(assoc.ChildID),
			TargetID: "entity-" + string(assoc.ParentID),
			Name:     assoc.Name,
			Type:     assocTypeStr(assoc.Type),
		})
	}

	// Build generalizations
	var generalizations []domainModelELKGeneralization
	for _, entity := range dm.Entities {
		if entity.GeneralizationRef == "" {
			continue
		}
		gen, parentID := buildGeneralization(entity, moduleEntityIDs, allEntityNames, ghostEntities)
		_ = parentID
		generalizations = append(generalizations, gen)
	}

	// Append ghost entities
	for _, ghost := range ghostEntities {
		if ghost.Width < elkMinWidth {
			ghost.Width = elkMinWidth
		}
		entities = append(entities, *ghost)
	}

	// Generate MDL source with source map
	mdlSource, sourceMap := e.buildDomainModelMdlSource(dm.Entities, moduleName)

	return e.emitDomainModelELK(domainModelELKData{
		Format:          "elk",
		Type:            "domainmodel",
		ModuleName:      moduleName,
		Entities:        entities,
		Associations:    associations,
		Generalizations: generalizations,
		MdlSource:       mdlSource,
		SourceMap:       sourceMap,
	})
}

// EntityFocusELK generates a focused ELK diagram showing only the selected entity
// and entities directly connected to it via associations or generalization.
func (e *Executor) EntityFocusELK(qualifiedName string) error {
	if e.reader == nil {
		return fmt.Errorf("not connected to a project")
	}

	parts := strings.SplitN(qualifiedName, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("expected qualified name Module.Entity, got: %s", qualifiedName)
	}
	moduleName, entityName := parts[0], parts[1]

	module, err := e.findModule(moduleName)
	if err != nil {
		return err
	}

	dm, err := e.reader.GetDomainModel(module.ID)
	if err != nil {
		return fmt.Errorf("failed to get domain model: %w", err)
	}

	// Find the focus entity
	var focusEntity *domainmodel.Entity
	for _, entity := range dm.Entities {
		if entity.Name == entityName {
			focusEntity = entity
			break
		}
	}
	if focusEntity == nil {
		return fmt.Errorf("entity not found: %s", qualifiedName)
	}

	// If this is a view entity with an OQL query, render query plan instead
	if classifyEntity(focusEntity) == "view" && focusEntity.OqlQuery != "" {
		return e.OqlQueryPlanELK(qualifiedName, focusEntity)
	}

	allEntityNames, _ := e.buildAllEntityNames()

	// Build map of all entities in this module by ID for quick lookup
	moduleEntitiesByID := make(map[model.ID]*domainmodel.Entity)
	for _, entity := range dm.Entities {
		moduleEntitiesByID[entity.ID] = entity
	}

	// Collect the set of entity IDs that should appear in the diagram
	includedIDs := make(map[model.ID]bool)
	includedIDs[focusEntity.ID] = true

	// Find associations touching the focus entity
	var relevantAssocs []*domainmodel.Association
	for _, assoc := range dm.Associations {
		if assoc.ParentID == focusEntity.ID || assoc.ChildID == focusEntity.ID {
			relevantAssocs = append(relevantAssocs, assoc)
			includedIDs[assoc.ParentID] = true
			includedIDs[assoc.ChildID] = true
		}
	}

	// Also scan all domain models for cross-module associations referencing this entity
	allDMs, _ := e.reader.ListDomainModels()
	for _, otherDM := range allDMs {
		if otherDM.ID == dm.ID {
			continue
		}
		for _, assoc := range otherDM.Associations {
			if assoc.ParentID == focusEntity.ID || assoc.ChildID == focusEntity.ID {
				relevantAssocs = append(relevantAssocs, assoc)
				includedIDs[assoc.ParentID] = true
				includedIDs[assoc.ChildID] = true
			}
		}
	}

	// Include generalization parent if present
	if focusEntity.GeneralizationID != "" {
		includedIDs[focusEntity.GeneralizationID] = true
	}

	// Also include any entities that generalize TO the focus entity
	for _, entity := range dm.Entities {
		if entity.GeneralizationID == focusEntity.ID {
			includedIDs[entity.ID] = true
		}
	}

	// Build entity nodes
	ghostEntities := make(map[string]*domainModelELKEntity)
	var entities []domainModelELKEntity

	for id := range includedIDs {
		if ent, ok := moduleEntitiesByID[id]; ok {
			elkEnt := buildELKEntity(ent)
			if id == focusEntity.ID {
				elkEnt.IsFocus = true
			}
			entities = append(entities, elkEnt)
		} else {
			// Entity from another module — add as ghost
			ghostID := "entity-" + string(id)
			if _, exists := ghostEntities[ghostID]; !exists {
				name := "Unknown"
				if qn, ok := allEntityNames[id]; ok {
					name = qn
				}
				ghost := makeGhostEntity(ghostID, name)
				ghostEntities[ghostID] = &ghost
			}
		}
	}

	// Build associations (only those involving included entities)
	var associations []domainModelELKAssoc
	for i, assoc := range relevantAssocs {
		associations = append(associations, domainModelELKAssoc{
			ID:       fmt.Sprintf("assoc-%d", i),
			SourceID: "entity-" + string(assoc.ChildID),
			TargetID: "entity-" + string(assoc.ParentID),
			Name:     assoc.Name,
			Type:     assocTypeStr(assoc.Type),
		})
	}

	// Build generalizations for included entities
	var generalizations []domainModelELKGeneralization
	// Focus entity's own generalization
	if focusEntity.GeneralizationRef != "" {
		gen, _ := buildGeneralization(focusEntity, includedIDs, allEntityNames, ghostEntities)
		generalizations = append(generalizations, gen)
	}
	// Entities that generalize to the focus entity
	for _, entity := range dm.Entities {
		if entity.GeneralizationID == focusEntity.ID && entity.ID != focusEntity.ID {
			gen, _ := buildGeneralization(entity, includedIDs, allEntityNames, ghostEntities)
			generalizations = append(generalizations, gen)
		}
	}

	// Append ghost entities
	for _, ghost := range ghostEntities {
		if ghost.Width < elkMinWidth {
			ghost.Width = elkMinWidth
		}
		entities = append(entities, *ghost)
	}

	// Generate MDL source with source map for focus entity
	mdlSource, sourceMap := e.buildDomainModelMdlSource([]*domainmodel.Entity{focusEntity}, moduleName)

	return e.emitDomainModelELK(domainModelELKData{
		Format:          "elk",
		Type:            "domainmodel",
		ModuleName:      moduleName,
		FocusEntity:     entityName,
		Entities:        entities,
		Associations:    associations,
		Generalizations: generalizations,
		MdlSource:       mdlSource,
		SourceMap:       sourceMap,
	})
}

// --- helpers ---

// buildAllEntityNames loads all entities across all modules.
// Returns ID -> "Module.Entity" map and ID -> module name map.
func (e *Executor) buildAllEntityNames() (map[model.ID]string, map[model.ID]string) {
	allEntityNames := make(map[model.ID]string)
	allEntityModules := make(map[model.ID]string)
	h, err := e.getHierarchy()
	if err != nil {
		return allEntityNames, allEntityModules
	}
	domainModels, _ := e.reader.ListDomainModels()
	for _, otherDM := range domainModels {
		modName := h.GetModuleName(otherDM.ContainerID)
		for _, entity := range otherDM.Entities {
			allEntityNames[entity.ID] = modName + "." + entity.Name
			allEntityModules[entity.ID] = modName
		}
	}
	return allEntityNames, allEntityModules
}

// buildELKEntity converts a domain model entity to an ELK node with calculated dimensions.
func buildELKEntity(entity *domainmodel.Entity) domainModelELKEntity {
	cat := classifyEntity(entity)
	var attrs []domainModelELKAttribute
	maxTextLen := float64(len(entity.Name))

	for _, attr := range entity.Attributes {
		typeName := attr.Type.GetTypeName()
		attrs = append(attrs, domainModelELKAttribute{
			Name: attr.Name,
			Type: typeName,
		})
		lineLen := float64(len(typeName) + 1 + len(attr.Name))
		if lineLen > maxTextLen {
			maxTextLen = lineLen
		}
	}

	width := maxTextLen*elkCharWidth + elkHPadding
	if width < elkMinWidth {
		width = elkMinWidth
	}
	height := elkHeaderHeight + float64(len(attrs))*elkAttrLineHeight
	if len(attrs) == 0 {
		height = elkHeaderHeight + elkAttrLineHeight
	}

	return domainModelELKEntity{
		ID:         "entity-" + string(entity.ID),
		Name:       entity.Name,
		Category:   cat,
		Attributes: attrs,
		Width:      width,
		Height:     height,
	}
}

// makeGhostEntity creates a minimal entity node for cross-module references.
func makeGhostEntity(id, name string) domainModelELKEntity {
	width := float64(len(name))*elkCharWidth + elkHPadding
	if width < elkMinWidth {
		width = elkMinWidth
	}
	return domainModelELKEntity{
		ID:       id,
		Name:     name,
		Category: "external",
		Width:    width,
		Height:   elkHeaderHeight + elkAttrLineHeight,
	}
}

// addGhostIfNeeded adds a ghost entity if the given ID is not in the included set.
func addGhostIfNeeded(id model.ID, includedIDs map[model.ID]bool, allEntityNames map[model.ID]string, ghosts map[string]*domainModelELKEntity) {
	if includedIDs[id] {
		return
	}
	ghostID := "entity-" + string(id)
	if _, exists := ghosts[ghostID]; exists {
		return
	}
	name := "Unknown"
	if qn, ok := allEntityNames[id]; ok {
		name = qn
	}
	ghost := makeGhostEntity(ghostID, name)
	ghosts[ghostID] = &ghost
}

// buildGeneralization builds a generalization record and creates ghost nodes as needed.
func buildGeneralization(entity *domainmodel.Entity, includedIDs map[model.ID]bool, allEntityNames map[model.ID]string, ghosts map[string]*domainModelELKEntity) (domainModelELKGeneralization, string) {
	var parentID string
	if entity.GeneralizationID != "" {
		parentID = "entity-" + string(entity.GeneralizationID)
		if !includedIDs[entity.GeneralizationID] {
			if _, exists := ghosts[parentID]; !exists {
				name := entity.GeneralizationRef
				if qn, ok := allEntityNames[entity.GeneralizationID]; ok {
					name = qn
				}
				ghost := makeGhostEntity(parentID, name)
				ghosts[parentID] = &ghost
			}
		}
	} else {
		syntheticID := "entity-gen-" + strings.ReplaceAll(entity.GeneralizationRef, ".", "-")
		parentID = syntheticID
		if _, exists := ghosts[syntheticID]; !exists {
			ghost := makeGhostEntity(syntheticID, entity.GeneralizationRef)
			ghosts[syntheticID] = &ghost
		}
	}

	return domainModelELKGeneralization{
		ChildID:    "entity-" + string(entity.ID),
		ParentID:   parentID,
		ParentName: entity.GeneralizationRef,
	}, parentID
}

// assocTypeStr returns the string representation of an association type.
func assocTypeStr(t domainmodel.AssociationType) string {
	if t == domainmodel.AssociationTypeReferenceSet {
		return "referenceSet"
	}
	return "reference"
}

// buildDomainModelMdlSource generates combined MDL source for a set of entities
// and returns the source string and a source map mapping entity ELK IDs to line ranges.
func (e *Executor) buildDomainModelMdlSource(entities []*domainmodel.Entity, moduleName string) (string, map[string]elkSourceRange) {
	sourceMap := make(map[string]elkSourceRange)
	var allSource strings.Builder
	lineCount := 0

	for i, entity := range entities {
		qn := ast.QualifiedName{Module: moduleName, Name: entity.Name}
		entityMdl, err := e.describeEntityToString(qn)
		if err != nil {
			continue
		}

		startLine := lineCount
		if i > 0 {
			allSource.WriteString("\n")
			lineCount++
		}
		allSource.WriteString(entityMdl)

		// Count lines in this entity's MDL
		entityLines := strings.Count(entityMdl, "\n")
		endLine := lineCount + entityLines
		if !strings.HasSuffix(entityMdl, "\n") {
			endLine--
		}

		sourceMap["entity-"+string(entity.ID)] = elkSourceRange{StartLine: startLine, EndLine: endLine}
		lineCount = endLine + 1
	}

	return allSource.String(), sourceMap
}

// emitDomainModelELK marshals and writes the domain model ELK data to output.
func (e *Executor) emitDomainModelELK(data domainModelELKData) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Fprint(e.output, string(out))
	return nil
}

// classifyEntity determines the category of an entity for visual styling.
func classifyEntity(entity *domainmodel.Entity) string {
	if strings.Contains(entity.Source, "OqlView") {
		return "view"
	}
	if strings.Contains(entity.Source, "OData") || entity.RemoteSource != "" || string(entity.RemoteSourceDocument) != "" {
		return "external"
	}
	if !entity.Persistable {
		return "nonpersistent"
	}
	return "persistent"
}
