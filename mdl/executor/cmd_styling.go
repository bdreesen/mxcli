// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/mendixlabs/mxcli/mdl/ast"
	"github.com/mendixlabs/mxcli/model"
	"github.com/mendixlabs/mxcli/sdk/pages"
)

// ============================================================================
// SHOW DESIGN PROPERTIES
// ============================================================================

func (e *Executor) execShowDesignProperties(s *ast.ShowDesignPropertiesStmt) error {
	if e.reader == nil {
		return fmt.Errorf("not connected to a project")
	}

	projectDir := filepath.Dir(e.mprPath)
	registry, err := loadThemeRegistry(projectDir)
	if err != nil {
		return fmt.Errorf("failed to load theme registry: %w", err)
	}

	if len(registry.WidgetProperties) == 0 {
		fmt.Fprintln(e.output, "No design properties found. Check that themesource/*/web/design-properties.json exists in the project directory.")
		return nil
	}

	if s.WidgetType != "" {
		// Show properties for a specific widget type
		dpKey := resolveDesignPropsKey(s.WidgetType)
		props := registry.GetPropertiesForWidget(dpKey)
		if len(props) == 0 {
			fmt.Fprintf(e.output, "No design properties found for widget type %s (%s)\n", s.WidgetType, dpKey)
			return nil
		}
		fmt.Fprintf(e.output, "Design Properties for %s:\n\n", s.WidgetType)
		e.printDesignProperties(registry, dpKey)
	} else {
		// Show all widget types and their properties
		keys := make([]string, 0, len(registry.WidgetProperties))
		for k := range registry.WidgetProperties {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			props := registry.WidgetProperties[key]
			if len(props) == 0 {
				continue
			}
			fmt.Fprintf(e.output, "=== %s ===\n", key)
			for _, p := range props {
				e.printOneProperty(p)
			}
			fmt.Fprintln(e.output)
		}
	}

	return nil
}

// printDesignProperties prints properties for a widget type, showing inherited "Widget" props separately.
func (e *Executor) printDesignProperties(registry *ThemeRegistry, dpKey string) {
	// Print inherited Widget properties
	if widgetProps, ok := registry.WidgetProperties["Widget"]; ok && len(widgetProps) > 0 {
		fmt.Fprintf(e.output, "From: Widget (inherited)\n")
		for _, p := range widgetProps {
			e.printOneProperty(p)
		}
	}

	// Print type-specific properties
	if dpKey != "Widget" {
		if typeProps, ok := registry.WidgetProperties[dpKey]; ok && len(typeProps) > 0 {
			fmt.Fprintf(e.output, "From: %s\n", dpKey)
			for _, p := range typeProps {
				e.printOneProperty(p)
			}
		}
	}
}

// printOneProperty prints a single design property in a readable format.
func (e *Executor) printOneProperty(p ThemeProperty) {
	switch p.Type {
	case "Toggle":
		fmt.Fprintf(e.output, "  %-24s Toggle      class: %s\n", p.Name, p.Class)
	case "Dropdown", "ColorPicker", "ToggleButtonGroup":
		options := make([]string, 0, len(p.Options))
		for _, o := range p.Options {
			options = append(options, o.Name)
		}
		fmt.Fprintf(e.output, "  %-24s %-11s [%s]\n", p.Name, p.Type, strings.Join(options, ", "))
	default:
		fmt.Fprintf(e.output, "  %-24s %s\n", p.Name, p.Type)
	}
}

// ============================================================================
// DESCRIBE STYLING
// ============================================================================

func (e *Executor) execDescribeStyling(s *ast.DescribeStylingStmt) error {
	if e.reader == nil {
		return fmt.Errorf("not connected to a project")
	}

	h, err := e.getHierarchy()
	if err != nil {
		return fmt.Errorf("failed to build hierarchy: %w", err)
	}

	var rawWidgets []rawWidget

	if s.ContainerType == "PAGE" {
		// Find page
		allPages, err := e.reader.ListPages()
		if err != nil {
			return fmt.Errorf("failed to list pages: %w", err)
		}

		var foundPage *pages.Page
		for _, p := range allPages {
			modID := h.FindModuleID(p.ContainerID)
			modName := h.GetModuleName(modID)
			if p.Name == s.ContainerName.Name && (s.ContainerName.Module == "" || modName == s.ContainerName.Module) {
				foundPage = p
				break
			}
		}
		if foundPage == nil {
			return fmt.Errorf("page %s not found", s.ContainerName.String())
		}
		rawWidgets = e.getPageWidgetsFromRaw(foundPage.ID)
	} else if s.ContainerType == "SNIPPET" {
		// Find snippet
		allSnippets, err := e.reader.ListSnippets()
		if err != nil {
			return fmt.Errorf("failed to list snippets: %w", err)
		}

		var foundSnippet *pages.Snippet
		for _, sn := range allSnippets {
			modID := h.FindModuleID(sn.ContainerID)
			modName := h.GetModuleName(modID)
			if sn.Name == s.ContainerName.Name && (s.ContainerName.Module == "" || modName == s.ContainerName.Module) {
				foundSnippet = sn
				break
			}
		}
		if foundSnippet == nil {
			return fmt.Errorf("snippet %s not found", s.ContainerName.String())
		}
		rawWidgets = e.getSnippetWidgetsFromRaw(foundSnippet.ID)
	}

	if len(rawWidgets) == 0 {
		fmt.Fprintf(e.output, "No widgets found in %s %s\n", s.ContainerType, s.ContainerName.String())
		return nil
	}

	// Collect styled widgets
	styledWidgets := collectStyledWidgets(rawWidgets, s.WidgetName)

	if len(styledWidgets) == 0 {
		if s.WidgetName != "" {
			return fmt.Errorf("widget %q not found in %s %s", s.WidgetName, s.ContainerType, s.ContainerName.String())
		}
		fmt.Fprintf(e.output, "No styled widgets found in %s %s\n", s.ContainerType, s.ContainerName.String())
		return nil
	}

	// Output
	for i, w := range styledWidgets {
		if i > 0 {
			fmt.Fprintln(e.output)
		}
		displayName := getWidgetDisplayName(w.Type)
		fmt.Fprintf(e.output, "WIDGET %s (%s)\n", w.Name, displayName)
		if w.Class != "" {
			fmt.Fprintf(e.output, "  Class: '%s'\n", w.Class)
		}
		if w.Style != "" {
			fmt.Fprintf(e.output, "  Style: '%s'\n", w.Style)
		}
		if len(w.DesignProperties) > 0 {
			fmt.Fprintf(e.output, "  DesignProperties: [")
			for j, dp := range w.DesignProperties {
				if j > 0 {
					fmt.Fprint(e.output, ", ")
				}
				if dp.ValueType == "toggle" {
					fmt.Fprintf(e.output, "'%s': ON", dp.Key)
				} else {
					fmt.Fprintf(e.output, "'%s': '%s'", dp.Key, dp.Option)
				}
			}
			fmt.Fprintln(e.output, "]")
		}
	}

	return nil
}

// collectStyledWidgets walks rawWidget tree and collects widgets that have styling.
// If widgetName is set, only returns the widget matching that name.
func collectStyledWidgets(widgets []rawWidget, widgetName string) []rawWidget {
	var result []rawWidget
	var walk func(ws []rawWidget)
	walk = func(ws []rawWidget) {
		for _, w := range ws {
			if widgetName != "" {
				// Looking for specific widget
				if w.Name == widgetName {
					result = append(result, w)
					return // Found it
				}
			} else {
				// Collect all widgets with any styling
				if w.Class != "" || w.Style != "" || len(w.DesignProperties) > 0 {
					result = append(result, w)
				}
			}
			// Walk children
			walk(w.Children)
			// Walk rows (for LayoutGrid)
			for _, row := range w.Rows {
				for _, col := range row.Columns {
					walk(col.Widgets)
				}
			}
			// Walk filter/controlbar widgets
			walk(w.FilterWidgets)
			walk(w.ControlBar)
		}
	}
	walk(widgets)
	return result
}

// ============================================================================
// ALTER STYLING
// ============================================================================

func (e *Executor) execAlterStyling(s *ast.AlterStylingStmt) error {
	if e.reader == nil {
		return fmt.Errorf("not connected to a project")
	}
	if e.writer == nil {
		return fmt.Errorf("project not opened for writing")
	}

	h, err := e.getHierarchy()
	if err != nil {
		return fmt.Errorf("failed to build hierarchy: %w", err)
	}

	if s.ContainerType == "PAGE" {
		return e.alterStylingOnPage(s, h)
	} else if s.ContainerType == "SNIPPET" {
		return e.alterStylingOnSnippet(s, h)
	}

	return fmt.Errorf("unsupported container type: %s", s.ContainerType)
}

func (e *Executor) alterStylingOnPage(s *ast.AlterStylingStmt, h *ContainerHierarchy) error {
	// Find page
	allPages, err := e.reader.ListPages()
	if err != nil {
		return fmt.Errorf("failed to list pages: %w", err)
	}

	var page *pages.Page
	for _, p := range allPages {
		modID := h.FindModuleID(p.ContainerID)
		modName := h.GetModuleName(modID)
		if p.Name == s.ContainerName.Name && (s.ContainerName.Module == "" || modName == s.ContainerName.Module) {
			page = p
			break
		}
	}
	if page == nil {
		return fmt.Errorf("page %s not found", s.ContainerName.String())
	}

	// Walk the page to find the widget by name
	found := false
	err = walkPageWidgets(page, func(widget any) error {
		name := getWidgetName(widget)
		if name != s.WidgetName {
			return nil
		}
		found = true
		return applyStylingAssignments(widget, s.Assignments, s.ClearDesignProps)
	})
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("widget %q not found in page %s", s.WidgetName, s.ContainerName.String())
	}

	// Save the page
	if err := e.writer.UpdatePage(page); err != nil {
		return fmt.Errorf("failed to save page: %w", err)
	}

	fmt.Fprintf(e.output, "Updated styling on widget %q in page %s\n", s.WidgetName, s.ContainerName.String())
	return nil
}

func (e *Executor) alterStylingOnSnippet(s *ast.AlterStylingStmt, h *ContainerHierarchy) error {
	// Find snippet
	allSnippets, err := e.reader.ListSnippets()
	if err != nil {
		return fmt.Errorf("failed to list snippets: %w", err)
	}

	var snippet *pages.Snippet
	for _, sn := range allSnippets {
		modID := h.FindModuleID(sn.ContainerID)
		modName := h.GetModuleName(modID)
		if sn.Name == s.ContainerName.Name && (s.ContainerName.Module == "" || modName == s.ContainerName.Module) {
			snippet = sn
			break
		}
	}
	if snippet == nil {
		return fmt.Errorf("snippet %s not found", s.ContainerName.String())
	}

	// Walk the snippet to find the widget by name
	found := false
	err = walkSnippetWidgets(snippet, func(widget any) error {
		name := getWidgetName(widget)
		if name != s.WidgetName {
			return nil
		}
		found = true
		return applyStylingAssignments(widget, s.Assignments, s.ClearDesignProps)
	})
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("widget %q not found in snippet %s", s.WidgetName, s.ContainerName.String())
	}

	// Save the snippet
	if err := e.writer.UpdateSnippet(snippet); err != nil {
		return fmt.Errorf("failed to save snippet: %w", err)
	}

	fmt.Fprintf(e.output, "Updated styling on widget %q in snippet %s\n", s.WidgetName, s.ContainerName.String())
	return nil
}

// getWidgetName extracts the Name from a widget using reflection.
func getWidgetName(widget any) string {
	if widget == nil {
		return ""
	}

	// Try Widget interface first
	if w, ok := widget.(pages.Widget); ok {
		return w.GetName()
	}

	// Fall back to reflection
	v := reflect.ValueOf(widget)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}

	// Try BaseWidget.Name
	if baseWidget := v.FieldByName("BaseWidget"); baseWidget.IsValid() {
		if nameField := baseWidget.FieldByName("Name"); nameField.IsValid() {
			return nameField.String()
		}
	}

	// Direct Name field
	if nameField := v.FieldByName("Name"); nameField.IsValid() {
		return nameField.String()
	}

	return ""
}

// applyStylingAssignments applies styling changes to a widget.
func applyStylingAssignments(widget any, assignments []ast.StylingAssignment, clearDesignProps bool) error {
	v := reflect.ValueOf(widget)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("widget is not a struct")
	}

	// Get BaseWidget
	baseWidget := v.FieldByName("BaseWidget")
	if !baseWidget.IsValid() {
		return fmt.Errorf("widget has no BaseWidget field")
	}

	// Clear design properties if requested
	if clearDesignProps {
		dpField := baseWidget.FieldByName("DesignProperties")
		if dpField.IsValid() && dpField.CanSet() {
			dpField.Set(reflect.Zero(dpField.Type()))
		}
	}

	for _, a := range assignments {
		switch a.Property {
		case "Class":
			classField := baseWidget.FieldByName("Class")
			if classField.IsValid() && classField.CanSet() {
				classField.SetString(a.Value)
			}
		case "Style":
			styleField := baseWidget.FieldByName("Style")
			if styleField.IsValid() && styleField.CanSet() {
				styleField.SetString(a.Value)
			}
		default:
			// Design property assignment
			if err := setDesignProperty(baseWidget, a); err != nil {
				return err
			}
		}
	}

	return nil
}

// setDesignProperty sets or updates a design property on the widget's BaseWidget.
func setDesignProperty(baseWidget reflect.Value, a ast.StylingAssignment) error {
	dpField := baseWidget.FieldByName("DesignProperties")
	if !dpField.IsValid() || !dpField.CanSet() {
		return fmt.Errorf("widget does not support design properties")
	}

	// Get existing design properties
	var existing []pages.DesignPropertyValue
	if !dpField.IsNil() {
		existing = dpField.Interface().([]pages.DesignPropertyValue)
	}

	if a.IsToggle && !a.ToggleOn {
		// OFF: remove the design property
		var updated []pages.DesignPropertyValue
		for _, dp := range existing {
			if dp.Key != a.Property {
				updated = append(updated, dp)
			}
		}
		dpField.Set(reflect.ValueOf(updated))
		return nil
	}

	// Update existing or append new
	found := false
	for i, dp := range existing {
		if dp.Key == a.Property {
			if a.IsToggle {
				existing[i].ValueType = "toggle"
				existing[i].Option = ""
			} else {
				existing[i].ValueType = "option"
				existing[i].Option = a.Value
			}
			found = true
			break
		}
	}

	if !found {
		newProp := pages.DesignPropertyValue{
			Key: a.Property,
		}
		if a.IsToggle {
			newProp.ValueType = "toggle"
		} else {
			newProp.ValueType = "option"
			newProp.Option = a.Value
		}
		existing = append(existing, newProp)
	}

	dpField.Set(reflect.ValueOf(existing))
	return nil
}

// findPageByName looks up a page by qualified name.
func (e *Executor) findPageByName(name ast.QualifiedName, h *ContainerHierarchy) (*pages.Page, error) {
	allPages, err := e.reader.ListPages()
	if err != nil {
		return nil, fmt.Errorf("failed to list pages: %w", err)
	}
	for _, p := range allPages {
		modID := h.FindModuleID(p.ContainerID)
		modName := h.GetModuleName(modID)
		if p.Name == name.Name && (name.Module == "" || modName == name.Module) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("page %s not found", name.String())
}

// findSnippetByName looks up a snippet by qualified name.
func (e *Executor) findSnippetByName(name ast.QualifiedName, h *ContainerHierarchy) (*pages.Snippet, model.ID, error) {
	allSnippets, err := e.reader.ListSnippets()
	if err != nil {
		return nil, "", fmt.Errorf("failed to list snippets: %w", err)
	}
	for _, s := range allSnippets {
		modID := h.FindModuleID(s.ContainerID)
		modName := h.GetModuleName(modID)
		if s.Name == name.Name && (name.Module == "" || modName == name.Module) {
			return s, modID, nil
		}
	}
	return nil, "", fmt.Errorf("snippet %s not found", name.String())
}
