// SPDX-License-Identifier: Apache-2.0

package catalog

import (
	"github.com/mendixlabs/mxcli/sdk/microflows"
)

// buildStrings extracts string literals from documents into the FTS5 strings table.
// Only runs in full mode.
func (b *Builder) buildStrings() error {
	if !b.fullMode {
		return nil
	}

	stmt, err := b.tx.Prepare(`
		INSERT INTO strings (QualifiedName, ObjectType, StringValue, StringContext, ModuleName)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	insert := func(qn, objType, value, ctx, module string) {
		if value == "" {
			return
		}
		stmt.Exec(qn, objType, value, ctx, module)
		count++
	}

	// Extract from pages (title, URL) — using cached list
	pageList, err := b.cachedPages()
	if err == nil {
		for _, pg := range pageList {
			moduleID := b.hierarchy.findModuleID(pg.ContainerID)
			moduleName := b.hierarchy.getModuleName(moduleID)
			qn := moduleName + "." + pg.Name

			// Page title translations
			if pg.Title != nil && pg.Title.Translations != nil {
				for _, t := range pg.Title.Translations {
					insert(qn, "PAGE", t, "page_title", moduleName)
				}
			}

			// Page URL
			if pg.URL != "" {
				insert(qn, "PAGE", pg.URL, "page_url", moduleName)
			}
		}
	}

	// Extract from microflows — using cached list
	mfList, err := b.cachedMicroflows()
	if err == nil {
		for _, mf := range mfList {
			moduleID := b.hierarchy.findModuleID(mf.ContainerID)
			moduleName := b.hierarchy.getModuleName(moduleID)
			qn := moduleName + "." + mf.Name

			// Documentation
			if mf.Documentation != "" {
				insert(qn, "MICROFLOW", mf.Documentation, "documentation", moduleName)
			}

			// Extract strings from activities
			extractActivityStrings(mf.ObjectCollection, qn, "MICROFLOW", moduleName, insert)
		}
	}

	// Extract from enumerations (value captions) — using cached list
	enums, err := b.cachedEnumerations()
	if err == nil {
		for _, enum := range enums {
			moduleID := b.hierarchy.findModuleID(enum.ContainerID)
			moduleName := b.hierarchy.getModuleName(moduleID)
			qn := moduleName + "." + enum.Name

			for _, val := range enum.Values {
				if val.Caption != nil && val.Caption.Translations != nil {
					for _, t := range val.Caption.Translations {
						insert(qn, "ENUMERATION", t, "enum_caption", moduleName)
					}
				}
			}
		}
	}

	b.report("strings", count)
	return nil
}

// extractActivityStrings extracts string literals from microflow/nanoflow activities.
func extractActivityStrings(oc *microflows.MicroflowObjectCollection, qn, objType, moduleName string, insert func(string, string, string, string, string)) {
	if oc == nil {
		return
	}

	for _, obj := range oc.Objects {
		act, ok := obj.(*microflows.ActionActivity)
		if !ok || act.Action == nil {
			continue
		}

		switch a := act.Action.(type) {
		case *microflows.LogMessageAction:
			if a.MessageTemplate != nil && a.MessageTemplate.Translations != nil {
				for _, t := range a.MessageTemplate.Translations {
					insert(qn, objType, t, "log_message", moduleName)
				}
			}
			if a.LogNodeName != "" {
				insert(qn, objType, a.LogNodeName, "log_node", moduleName)
			}
		case *microflows.ShowMessageAction:
			if a.Template != nil && a.Template.Translations != nil {
				for _, t := range a.Template.Translations {
					insert(qn, objType, t, "show_message", moduleName)
				}
			}
		case *microflows.ValidationFeedbackAction:
			if a.Template != nil && a.Template.Translations != nil {
				for _, t := range a.Template.Translations {
					insert(qn, objType, t, "validation_message", moduleName)
				}
			}
		}
	}
}
