// SPDX-License-Identifier: Apache-2.0

package mpk

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"
)

// findTestMPK finds a ComboBox .mpk file in the test projects directory.
func findTestMPK(t *testing.T) string {
	t.Helper()
	// Try multiple known locations
	candidates := []string{
		filepath.Join("..", "..", "..", "mx-test-projects", "template-app-116", "widgets", "com.mendix.widget.web.Combobox.mpk"),
		filepath.Join("..", "..", "..", "mx-test-projects", "LatoProductInventory", "widgets", "com.mendix.widget.web.Combobox.mpk"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	t.Skip("No test .mpk file found")
	return ""
}

func findTestProjectDir(t *testing.T) string {
	t.Helper()
	candidates := []string{
		filepath.Join("..", "..", "..", "mx-test-projects", "template-app-116"),
		filepath.Join("..", "..", "..", "mx-test-projects", "LatoProductInventory"),
	}
	for _, c := range candidates {
		widgetsDir := filepath.Join(c, "widgets")
		if _, err := os.Stat(widgetsDir); err == nil {
			return c
		}
	}
	t.Skip("No test project directory found")
	return ""
}

func TestParseMPK(t *testing.T) {
	mpkPath := findTestMPK(t)
	ClearCache()

	def, err := ParseMPK(mpkPath)
	if err != nil {
		t.Fatalf("ParseMPK failed: %v", err)
	}

	if def.ID != "com.mendix.widget.web.combobox.Combobox" {
		t.Errorf("unexpected widget ID: %s", def.ID)
	}

	if def.Name != "Combo box" {
		t.Errorf("unexpected widget name: %s", def.Name)
	}

	if def.Version == "" {
		t.Error("version should not be empty")
	}

	if len(def.Properties) == 0 {
		t.Error("expected at least one property")
	}

	// Check that we found some known properties
	keys := def.PropertyKeys()
	expectedKeys := []string{"source", "attributeEnumeration", "clearable", "filterType"}
	for _, k := range expectedKeys {
		if !keys[k] {
			t.Errorf("expected property key %q not found", k)
		}
	}

	// Check system properties
	if len(def.SystemProps) == 0 {
		t.Error("expected at least one system property")
	}
	sysKeys := def.SystemPropertyKeys()
	if !sysKeys["Label"] {
		t.Error("expected system property 'Label'")
	}

	// Check property types
	sourceProp := def.FindProperty("source")
	if sourceProp == nil {
		t.Fatal("source property not found")
	}
	if sourceProp.Type != "enumeration" {
		t.Errorf("expected source type 'enumeration', got %q", sourceProp.Type)
	}
	if sourceProp.DefaultValue != "context" {
		t.Errorf("expected source default 'context', got %q", sourceProp.DefaultValue)
	}

	// Check category tracking
	if sourceProp.Category == "" {
		t.Error("expected source to have a category")
	}
}

func TestParseMPK_Cached(t *testing.T) {
	mpkPath := findTestMPK(t)
	ClearCache()

	def1, err := ParseMPK(mpkPath)
	if err != nil {
		t.Fatalf("first ParseMPK failed: %v", err)
	}

	def2, err := ParseMPK(mpkPath)
	if err != nil {
		t.Fatalf("second ParseMPK failed: %v", err)
	}

	// Should return same pointer from cache
	if def1 != def2 {
		t.Error("expected cached result to return same pointer")
	}
}

func TestFindMPK(t *testing.T) {
	projectDir := findTestProjectDir(t)
	ClearCache()

	// Find ComboBox
	mpkPath, err := FindMPK(projectDir, "com.mendix.widget.web.combobox.Combobox")
	if err != nil {
		t.Fatalf("FindMPK failed: %v", err)
	}
	if mpkPath == "" {
		t.Fatal("expected to find ComboBox .mpk")
	}
	if !filepath.IsAbs(mpkPath) && !fileExists(mpkPath) {
		t.Errorf("mpk path should exist: %s", mpkPath)
	}

	// Find non-existent widget
	mpkPath2, err := FindMPK(projectDir, "com.example.nonexistent.Widget")
	if err != nil {
		t.Fatalf("FindMPK for nonexistent should not error: %v", err)
	}
	if mpkPath2 != "" {
		t.Errorf("expected empty path for nonexistent widget, got: %s", mpkPath2)
	}
}

func TestFindMPK_Cached(t *testing.T) {
	projectDir := findTestProjectDir(t)
	ClearCache()

	// First call scans directory
	path1, _ := FindMPK(projectDir, "com.mendix.widget.web.combobox.Combobox")
	// Second call should use cache
	path2, _ := FindMPK(projectDir, "com.mendix.widget.web.combobox.Combobox")

	if path1 != path2 {
		t.Error("expected same path from cache")
	}
}

func TestNormalizeType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"attribute", "attribute"},
		{"Attribute", "attribute"},
		{"ATTRIBUTE", "attribute"},
		{"textTemplate", "textTemplate"},
		{"TextTemplate", "textTemplate"},
		{"datasource", "datasource"},
		{"unknownType", "unknownType"},
	}

	for _, tt := range tests {
		result := NormalizeType(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeType(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestObjectListExtraction_GroupedShape covers the
// <properties><propertyGroup><property>...</property></propertyGroup></properties>
// nesting (Accordion, DataGrid).
func TestObjectListExtraction_GroupedShape(t *testing.T) {
	xmlSrc := `<widget id="com.example.Acc" pluginWidget="true">
		<properties>
			<propertyGroup caption="Top">
				<property key="groups" type="object" isList="true">
					<caption>Groups</caption>
					<properties>
						<propertyGroup caption="Inner">
							<property key="headerText" type="textTemplate"><caption>Header</caption></property>
							<property key="content" type="widgets"><caption>Content</caption></property>
						</propertyGroup>
					</properties>
				</property>
			</propertyGroup>
		</properties>
	</widget>`

	var w xmlWidget
	if err := xml.Unmarshal([]byte(xmlSrc), &w); err != nil {
		t.Fatalf("xml.Unmarshal: %v", err)
	}
	def := &WidgetDefinition{ID: "com.example.Acc"}
	for _, pg := range w.PropertyGroups {
		walkPropertyGroup(pg, "", def)
	}

	if len(def.Properties) != 1 {
		t.Fatalf("Properties count = %d, want 1", len(def.Properties))
	}
	groups := def.Properties[0]
	if groups.Key != "groups" || groups.Type != "object" || !groups.IsList {
		t.Errorf("groups property has unexpected fields: %+v", groups)
	}
	if len(groups.Children) != 2 {
		t.Fatalf("groups.Children count = %d, want 2", len(groups.Children))
	}
	wantChildren := map[string]string{"headerText": "textTemplate", "content": "widgets"}
	for _, c := range groups.Children {
		if want, ok := wantChildren[c.Key]; !ok {
			t.Errorf("unexpected child key %q", c.Key)
		} else if c.Type != want {
			t.Errorf("child %q Type = %q, want %q", c.Key, c.Type, want)
		}
	}
}

// TestObjectListExtraction_FlatShape covers the
// <properties><property>...</property></properties> nesting (PopupMenu, Maps).
// Without the NestedDirectProps field, the children would be silently dropped.
func TestObjectListExtraction_FlatShape(t *testing.T) {
	xmlSrc := `<widget id="com.example.Menu" pluginWidget="true">
		<properties>
			<propertyGroup caption="Top">
				<property key="basicItems" type="object" isList="true">
					<caption>Items</caption>
					<properties>
						<property key="caption" type="textTemplate"><caption>Caption</caption></property>
						<property key="action" type="action"><caption>Action</caption></property>
					</properties>
				</property>
			</propertyGroup>
		</properties>
	</widget>`

	var w xmlWidget
	if err := xml.Unmarshal([]byte(xmlSrc), &w); err != nil {
		t.Fatalf("xml.Unmarshal: %v", err)
	}
	def := &WidgetDefinition{ID: "com.example.Menu"}
	for _, pg := range w.PropertyGroups {
		walkPropertyGroup(pg, "", def)
	}

	if len(def.Properties) != 1 {
		t.Fatalf("Properties count = %d, want 1", len(def.Properties))
	}
	items := def.Properties[0]
	if items.Key != "basicItems" || !items.IsList {
		t.Errorf("basicItems property has unexpected fields: %+v", items)
	}
	if len(items.Children) != 2 {
		t.Fatalf("basicItems.Children count = %d, want 2 (children dropped — flat XML shape unsupported?)", len(items.Children))
	}
	wantChildren := map[string]string{"caption": "textTemplate", "action": "action"}
	for _, c := range items.Children {
		if want, ok := wantChildren[c.Key]; !ok {
			t.Errorf("unexpected child key %q", c.Key)
		} else if c.Type != want {
			t.Errorf("child %q Type = %q, want %q", c.Key, c.Type, want)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
