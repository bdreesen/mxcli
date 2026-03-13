// SPDX-License-Identifier: Apache-2.0

package rules

import "testing"

func TestIsUnconfiguredImage_StaticNoSource(t *testing.T) {
	w := map[string]any{
		"$Type": "Forms$StaticImageViewer",
		"Name":  "img1",
		"Image": nil,
	}
	if !IsUnconfiguredImage("Forms$StaticImageViewer", w) {
		t.Error("expected unconfigured static image with nil Image")
	}
}

func TestIsUnconfiguredImage_StaticWithSource(t *testing.T) {
	w := map[string]any{
		"$Type": "Forms$StaticImageViewer",
		"Name":  "img1",
		"Image": []byte{0x01, 0x02, 0x03}, // binary image reference
	}
	if IsUnconfiguredImage("Forms$StaticImageViewer", w) {
		t.Error("expected configured static image with Image set")
	}
}

func TestIsUnconfiguredImage_DynamicNoDataSource(t *testing.T) {
	w := map[string]any{
		"$Type":      "Forms$ImageViewer",
		"Name":       "img1",
		"DataSource": nil,
	}
	if !IsUnconfiguredImage("Forms$ImageViewer", w) {
		t.Error("expected unconfigured dynamic image with nil DataSource")
	}
}

func TestIsUnconfiguredImage_DynamicNoEntityRef(t *testing.T) {
	w := map[string]any{
		"$Type": "Forms$ImageViewer",
		"Name":  "img1",
		"DataSource": map[string]any{
			"EntityRef": nil,
		},
	}
	if !IsUnconfiguredImage("Forms$ImageViewer", w) {
		t.Error("expected unconfigured dynamic image with nil EntityRef")
	}
}

func TestIsUnconfiguredImage_DynamicWithSource(t *testing.T) {
	w := map[string]any{
		"$Type": "Forms$ImageViewer",
		"Name":  "img1",
		"DataSource": map[string]any{
			"EntityRef": "MyModule.MyEntity",
		},
	}
	if IsUnconfiguredImage("Forms$ImageViewer", w) {
		t.Error("expected configured dynamic image with EntityRef set")
	}
}

func TestIsUnconfiguredImage_OtherWidgetType(t *testing.T) {
	w := map[string]any{
		"$Type": "Forms$TextBox",
		"Name":  "txt1",
	}
	if IsUnconfiguredImage("Forms$TextBox", w) {
		t.Error("expected non-image widget to return false")
	}
}

func TestFindImageWidgets_Nested(t *testing.T) {
	// Simulate a page with an image widget nested inside a LayoutGrid
	rawData := map[string]any{
		"FormCall": map[string]any{
			"Arguments": []any{
				map[string]any{
					"Widgets": []any{
						map[string]any{
							"$Type": "Forms$LayoutGrid",
							"Name":  "grid1",
							"Rows": []any{
								map[string]any{
									"Columns": []any{
										map[string]any{
											"Widgets": []any{
												map[string]any{
													"$Type": "Forms$StaticImageViewer",
													"Name":  "imgNested",
													"Image": nil,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	widgets := findImageWidgets(rawData)
	if len(widgets) != 1 {
		t.Fatalf("expected 1 image widget, got %d", len(widgets))
	}
	if widgets[0].Name != "imgNested" {
		t.Errorf("expected widget name 'imgNested', got %q", widgets[0].Name)
	}
	if widgets[0].Configured {
		t.Error("expected widget to be unconfigured")
	}
}

func TestFindImageWidgets_Empty(t *testing.T) {
	rawData := map[string]any{
		"FormCall": map[string]any{
			"Arguments": []any{
				map[string]any{
					"Widgets": []any{
						map[string]any{
							"$Type": "Forms$TextBox",
							"Name":  "txt1",
						},
					},
				},
			},
		},
	}

	widgets := findImageWidgets(rawData)
	if len(widgets) != 0 {
		t.Errorf("expected 0 image widgets, got %d", len(widgets))
	}
}

func TestFindImageWidgets_SnippetStructure(t *testing.T) {
	// Snippets have Widgets directly (no FormCall)
	rawData := map[string]any{
		"Widgets": []any{
			map[string]any{
				"$Type": "Forms$StaticImageViewer",
				"Name":  "imgSnippet",
				"Image": nil,
			},
		},
	}

	widgets := findImageWidgets(rawData)
	if len(widgets) != 1 {
		t.Fatalf("expected 1 image widget, got %d", len(widgets))
	}
	if widgets[0].Name != "imgSnippet" {
		t.Errorf("expected widget name 'imgSnippet', got %q", widgets[0].Name)
	}
}

func TestFindImageWidgets_ConfiguredIgnored(t *testing.T) {
	rawData := map[string]any{
		"FormCall": map[string]any{
			"Arguments": []any{
				map[string]any{
					"Widgets": []any{
						map[string]any{
							"$Type": "Forms$StaticImageViewer",
							"Name":  "imgConfigured",
							"Image": []byte{0x01},
						},
					},
				},
			},
		},
	}

	widgets := findImageWidgets(rawData)
	if len(widgets) != 1 {
		t.Fatalf("expected 1 image widget, got %d", len(widgets))
	}
	if !widgets[0].Configured {
		t.Error("expected widget to be configured")
	}
}

func TestDocNameFromQualified(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Module.Page", "Page"},
		{"MyModule.Customer_Edit", "Customer_Edit"},
		{"SimpleName", "SimpleName"},
		{"A.B.C", "C"},
	}
	for _, tt := range tests {
		got := docNameFromQualified(tt.input)
		if got != tt.expected {
			t.Errorf("docNameFromQualified(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestGetBsonArray_TypeIndicator(t *testing.T) {
	// Mendix BSON arrays have a leading int32 type indicator
	arr := []any{int32(5), map[string]any{"$Type": "Forms$TextBox"}}
	result := getBsonArray(arr)
	if len(result) != 1 {
		t.Fatalf("expected 1 element after stripping type indicator, got %d", len(result))
	}
}

func TestGetBsonArray_NoTypeIndicator(t *testing.T) {
	arr := []any{map[string]any{"$Type": "Forms$TextBox"}}
	result := getBsonArray(arr)
	if len(result) != 1 {
		t.Fatalf("expected 1 element, got %d", len(result))
	}
}

func TestGetBsonArray_Nil(t *testing.T) {
	result := getBsonArray(nil)
	if result != nil {
		t.Error("expected nil for nil input")
	}
}
