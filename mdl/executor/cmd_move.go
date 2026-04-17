// SPDX-License-Identifier: Apache-2.0

// Package executor - MOVE command
package executor

import (
	"fmt"

	"github.com/mendixlabs/mxcli/mdl/ast"
	mdlerrors "github.com/mendixlabs/mxcli/mdl/errors"
	"github.com/mendixlabs/mxcli/model"
	"github.com/mendixlabs/mxcli/sdk/domainmodel"
)

// execMove handles MOVE PAGE/MICROFLOW/SNIPPET/NANOFLOW/ENTITY/ENUMERATION statements.
func (e *Executor) execMove(s *ast.MoveStmt) error {
	if e.writer == nil {
		return mdlerrors.NewNotConnected()
	}

	// Find the source module
	sourceModule, err := e.findModule(s.Name.Module)
	if err != nil {
		return mdlerrors.NewBackend("find source module", err)
	}

	// Determine target module
	var targetModule *model.Module
	isCrossModuleMove := false
	if s.TargetModule != "" {
		targetModule, err = e.findModule(s.TargetModule)
		if err != nil {
			return mdlerrors.NewBackend("find target module", err)
		}
		isCrossModuleMove = targetModule.ID != sourceModule.ID
	} else {
		targetModule = sourceModule
	}

	// Entity moves are handled specially (entities are embedded in domain models, not top-level units)
	if s.DocumentType == ast.DocumentTypeEntity {
		return e.moveEntity(s.Name, sourceModule, targetModule)
	}

	// Resolve target container (folder or module root)
	var targetContainerID model.ID
	if s.Folder != "" {
		targetContainerID, err = e.resolveFolder(targetModule.ID, s.Folder)
		if err != nil {
			return mdlerrors.NewBackend("resolve target folder", err)
		}
	} else {
		targetContainerID = targetModule.ID
	}

	// Execute move based on document type
	switch s.DocumentType {
	case ast.DocumentTypePage:
		if err := e.movePage(s.Name, targetContainerID); err != nil {
			return err
		}
	case ast.DocumentTypeMicroflow:
		if err := e.moveMicroflow(s.Name, targetContainerID); err != nil {
			return err
		}
	case ast.DocumentTypeSnippet:
		if err := e.moveSnippet(s.Name, targetContainerID); err != nil {
			return err
		}
	case ast.DocumentTypeNanoflow:
		if err := e.moveNanoflow(s.Name, targetContainerID); err != nil {
			return err
		}
	case ast.DocumentTypeEnumeration:
		return e.moveEnumeration(s.Name, targetContainerID, targetModule.Name)
	case ast.DocumentTypeConstant:
		if err := e.moveConstant(s.Name, targetContainerID); err != nil {
			return err
		}
	case ast.DocumentTypeDatabaseConnection:
		if err := e.moveDatabaseConnection(s.Name, targetContainerID); err != nil {
			return err
		}
	default:
		return mdlerrors.NewUnsupported("unsupported document type: " + string(s.DocumentType))
	}

	// For cross-module moves, update all BY_NAME references throughout the project
	if isCrossModuleMove {
		if err := e.updateQualifiedNameRefs(s.Name, targetModule.Name); err != nil {
			return err
		}
	}

	return nil
}

// updateQualifiedNameRefs updates all BY_NAME references to an element after a cross-module move.
func (e *Executor) updateQualifiedNameRefs(name ast.QualifiedName, newModule string) error {
	oldQN := name.String()               // "OldModule.ElementName"
	newQN := newModule + "." + name.Name // "NewModule.ElementName"
	updated, err := e.writer.UpdateQualifiedNameInAllUnits(oldQN, newQN)
	if err != nil {
		return mdlerrors.NewBackend("update references", err)
	}
	if updated > 0 {
		fmt.Fprintf(e.output, "Updated references in %d document(s): %s → %s\n", updated, oldQN, newQN)
	}
	return nil
}

// movePage moves a page to a new container.
func (e *Executor) movePage(name ast.QualifiedName, targetContainerID model.ID) error {
	// Find the page
	pages, err := e.reader.ListPages()
	if err != nil {
		return mdlerrors.NewBackend("list pages", err)
	}

	h, err := e.getHierarchy()
	if err != nil {
		return mdlerrors.NewBackend("build hierarchy", err)
	}

	for _, p := range pages {
		modID := h.FindModuleID(p.ContainerID)
		modName := h.GetModuleName(modID)
		if modName == name.Module && p.Name == name.Name {
			// Update container ID and move the unit
			p.ContainerID = targetContainerID
			if err := e.writer.MovePage(p); err != nil {
				return mdlerrors.NewBackend("move page", err)
			}
			fmt.Fprintf(e.output, "Moved page %s to new location\n", name.String())
			return nil
		}
	}

	return mdlerrors.NewNotFound("page", name.String())
}

// moveMicroflow moves a microflow to a new container.
func (e *Executor) moveMicroflow(name ast.QualifiedName, targetContainerID model.ID) error {
	// Find the microflow
	mfs, err := e.reader.ListMicroflows()
	if err != nil {
		return mdlerrors.NewBackend("list microflows", err)
	}

	h, err := e.getHierarchy()
	if err != nil {
		return mdlerrors.NewBackend("build hierarchy", err)
	}

	for _, mf := range mfs {
		modID := h.FindModuleID(mf.ContainerID)
		modName := h.GetModuleName(modID)
		if modName == name.Module && mf.Name == name.Name {
			// Update container ID and move the unit
			mf.ContainerID = targetContainerID
			if err := e.writer.MoveMicroflow(mf); err != nil {
				return mdlerrors.NewBackend("move microflow", err)
			}
			fmt.Fprintf(e.output, "Moved microflow %s to new location\n", name.String())
			return nil
		}
	}

	return mdlerrors.NewNotFound("microflow", name.String())
}

// moveSnippet moves a snippet to a new container.
func (e *Executor) moveSnippet(name ast.QualifiedName, targetContainerID model.ID) error {
	// Find the snippet
	snippets, err := e.reader.ListSnippets()
	if err != nil {
		return mdlerrors.NewBackend("list snippets", err)
	}

	h, err := e.getHierarchy()
	if err != nil {
		return mdlerrors.NewBackend("build hierarchy", err)
	}

	for _, s := range snippets {
		modID := h.FindModuleID(s.ContainerID)
		modName := h.GetModuleName(modID)
		if modName == name.Module && s.Name == name.Name {
			// Update container ID and move the unit
			s.ContainerID = targetContainerID
			if err := e.writer.MoveSnippet(s); err != nil {
				return mdlerrors.NewBackend("move snippet", err)
			}
			fmt.Fprintf(e.output, "Moved snippet %s to new location\n", name.String())
			return nil
		}
	}

	return mdlerrors.NewNotFound("snippet", name.String())
}

// moveNanoflow moves a nanoflow to a new container.
func (e *Executor) moveNanoflow(name ast.QualifiedName, targetContainerID model.ID) error {
	// Find the nanoflow
	nfs, err := e.reader.ListNanoflows()
	if err != nil {
		return mdlerrors.NewBackend("list nanoflows", err)
	}

	h, err := e.getHierarchy()
	if err != nil {
		return mdlerrors.NewBackend("build hierarchy", err)
	}

	for _, nf := range nfs {
		modID := h.FindModuleID(nf.ContainerID)
		modName := h.GetModuleName(modID)
		if modName == name.Module && nf.Name == name.Name {
			// Update container ID and move the unit
			nf.ContainerID = targetContainerID
			if err := e.writer.MoveNanoflow(nf); err != nil {
				return mdlerrors.NewBackend("move nanoflow", err)
			}
			fmt.Fprintf(e.output, "Moved nanoflow %s to new location\n", name.String())
			return nil
		}
	}

	return mdlerrors.NewNotFound("nanoflow", name.String())
}

// moveEntity moves an entity from one domain model to another.
// Entities are embedded inside DomainModel documents, so we must remove from source DM and add to target DM.
// Associations referencing the entity are converted to CrossAssociations.
// ViewEntitySourceDocuments for view entities are also moved.
func (e *Executor) moveEntity(name ast.QualifiedName, sourceModule, targetModule *model.Module) error {
	// Get source domain model
	sourceDM, err := e.reader.GetDomainModel(sourceModule.ID)
	if err != nil {
		return mdlerrors.NewBackend("get source domain model", err)
	}

	// Find the entity in the source domain model
	var entity *domainmodel.Entity
	for _, ent := range sourceDM.Entities {
		if ent.Name == name.Name {
			entity = ent
			break
		}
	}
	if entity == nil {
		return mdlerrors.NewNotFound("entity", name.String())
	}

	// Get target domain model
	targetDM, err := e.reader.GetDomainModel(targetModule.ID)
	if err != nil {
		return mdlerrors.NewBackend("get target domain model", err)
	}

	// Move entity via writer (converts associations to CrossAssociations, updates validation rule refs)
	convertedAssocs, err := e.writer.MoveEntity(entity, sourceDM.ID, targetDM.ID, sourceModule.Name, targetModule.Name)
	if err != nil {
		return mdlerrors.NewBackend("move entity", err)
	}

	// Move ViewEntitySourceDocument for view entities
	if entity.Source == "DomainModels$OqlViewEntitySource" && entity.SourceDocumentRef != "" {
		// The SourceDocumentRef was already updated by MoveEntity to use the new module name.
		// Extract the original doc name (before the module prefix was changed).
		docName := name.Name // ViewEntitySourceDocument name matches the entity name
		if err := e.writer.MoveViewEntitySourceDocument(sourceModule.Name, targetModule.ID, docName); err != nil {
			fmt.Fprintf(e.output, "Warning: Could not move ViewEntitySourceDocument: %v\n", err)
		}
	}

	// Update OQL queries in all ViewEntitySourceDocuments that reference the moved entity
	oldQualifiedName := name.String()                       // e.g., "DmTest.Customer"
	newQualifiedName := targetModule.Name + "." + name.Name // e.g., "DmTest2.Customer"
	if oqlUpdated, err := e.writer.UpdateOqlQueriesForMovedEntity(oldQualifiedName, newQualifiedName); err != nil {
		fmt.Fprintf(e.output, "Warning: Could not update OQL queries: %v\n", err)
	} else if oqlUpdated > 0 {
		fmt.Fprintf(e.output, "Updated %d OQL query(ies) referencing %s\n", oqlUpdated, oldQualifiedName)
	}

	fmt.Fprintf(e.output, "Moved entity %s to %s\n", name.String(), targetModule.Name)
	if len(convertedAssocs) > 0 {
		fmt.Fprintf(e.output, "Converted %d association(s) to cross-module associations:\n", len(convertedAssocs))
		for _, assocName := range convertedAssocs {
			fmt.Fprintf(e.output, "  - %s\n", assocName)
		}
	}
	return nil
}

// moveEnumeration moves an enumeration to a new container.
// For cross-module moves, updates all EnumerationAttributeType references across all domain models.
func (e *Executor) moveEnumeration(name ast.QualifiedName, targetContainerID model.ID, targetModuleName string) error {
	enum := e.findEnumeration(name.Module, name.Name)
	if enum == nil {
		return mdlerrors.NewNotFound("enumeration", name.String())
	}

	oldQualifiedName := name.String() // e.g., "DmTest.Country"
	enum.ContainerID = targetContainerID
	if err := e.writer.MoveEnumeration(enum); err != nil {
		return mdlerrors.NewBackend("move enumeration", err)
	}

	// For cross-module moves, update enumeration references in all domain models
	if targetModuleName != "" && targetModuleName != name.Module {
		newQualifiedName := targetModuleName + "." + name.Name
		if err := e.writer.UpdateEnumerationRefsInAllDomainModels(oldQualifiedName, newQualifiedName); err != nil {
			fmt.Fprintf(e.output, "Warning: Could not update enumeration references: %v\n", err)
		} else {
			fmt.Fprintf(e.output, "Updated enumeration references: %s -> %s\n", oldQualifiedName, newQualifiedName)
		}
	}

	fmt.Fprintf(e.output, "Moved enumeration %s to new location\n", name.String())
	return nil
}

// moveConstant moves a constant to a new container.
func (e *Executor) moveConstant(name ast.QualifiedName, targetContainerID model.ID) error {
	constants, err := e.reader.ListConstants()
	if err != nil {
		return mdlerrors.NewBackend("list constants", err)
	}

	h, err := e.getHierarchy()
	if err != nil {
		return mdlerrors.NewBackend("build hierarchy", err)
	}

	for _, c := range constants {
		modID := h.FindModuleID(c.ContainerID)
		modName := h.GetModuleName(modID)
		if modName == name.Module && c.Name == name.Name {
			c.ContainerID = targetContainerID
			if err := e.writer.MoveConstant(c); err != nil {
				return mdlerrors.NewBackend("move constant", err)
			}
			fmt.Fprintf(e.output, "Moved constant %s to new location\n", name.String())
			return nil
		}
	}

	return mdlerrors.NewNotFound("constant", name.String())
}

// moveDatabaseConnection moves a database connection to a new container.
func (e *Executor) moveDatabaseConnection(name ast.QualifiedName, targetContainerID model.ID) error {
	connections, err := e.reader.ListDatabaseConnections()
	if err != nil {
		return mdlerrors.NewBackend("list database connections", err)
	}

	h, err := e.getHierarchy()
	if err != nil {
		return mdlerrors.NewBackend("build hierarchy", err)
	}

	for _, conn := range connections {
		modID := h.FindModuleID(conn.ContainerID)
		modName := h.GetModuleName(modID)
		if modName == name.Module && conn.Name == name.Name {
			conn.ContainerID = targetContainerID
			if err := e.writer.MoveDatabaseConnection(conn); err != nil {
				return mdlerrors.NewBackend("move database connection", err)
			}
			fmt.Fprintf(e.output, "Moved database connection %s to new location\n", name.String())
			return nil
		}
	}

	return mdlerrors.NewNotFound("database connection", name.String())
}
