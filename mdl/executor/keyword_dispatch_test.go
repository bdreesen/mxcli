// SPDX-License-Identifier: Apache-2.0

package executor

import "testing"

func TestResolveKeyword(t *testing.T) {
	tests := []struct {
		name       string
		keyword    string
		version    string
		wantOK     bool
		wantKind   string
		wantWidget string
		wantDeprec string
	}{
		{
			name:       "DATAGRID on 11.9 → pluggable",
			keyword:    "DATAGRID",
			version:    "11.9.0",
			wantOK:     true,
			wantKind:   "pluggable",
			wantWidget: "com.mendix.widget.web.datagrid.Datagrid",
		},
		{
			name:       "DATAGRID on 10.18 → pluggable (Datagrid 2.x default)",
			keyword:    "DATAGRID",
			version:    "10.18.0",
			wantOK:     true,
			wantKind:   "pluggable",
			wantWidget: "com.mendix.widget.web.datagrid.Datagrid",
		},
		{
			name:       "DATAGRID lowercase still resolves",
			keyword:    "datagrid",
			version:    "11.9.0",
			wantOK:     true,
			wantKind:   "pluggable",
			wantWidget: "com.mendix.widget.web.datagrid.Datagrid",
		},
		{
			name:    "Unknown keyword returns false",
			keyword: "MYCUSTOM",
			version: "11.9.0",
			wantOK:  false,
		},
		{
			name:    "DATAGRID below 9.0 → no match (not in any binding range)",
			keyword: "DATAGRID",
			version: "8.18.0",
			wantOK:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, ok := ResolveKeyword(tc.keyword, tc.version)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if res.Kind != tc.wantKind {
				t.Errorf("Kind = %q, want %q", res.Kind, tc.wantKind)
			}
			if res.WidgetID != tc.wantWidget {
				t.Errorf("WidgetID = %q, want %q", res.WidgetID, tc.wantWidget)
			}
			if res.DeprecatedFrom != tc.wantDeprec {
				t.Errorf("DeprecatedFrom = %q, want %q", res.DeprecatedFrom, tc.wantDeprec)
			}
		})
	}
}

func TestCompareVersion(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"11.0.0", "10.99.99", 1},
		{"10.18.0", "10.24.0", -1},
		{"11.9.0", "11.9.0", 0},
		{"11", "11.0.0", 0},
		{"11.9", "11.9.0", 0},
		{"9.0.0", "10.0.0", -1},
		{"abc", "1.0.0", -1}, // non-numeric → 0, less than 1
	}
	for _, tc := range tests {
		got := compareVersion(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("compareVersion(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestVersionInRange(t *testing.T) {
	tests := []struct {
		version, min, max string
		want              bool
	}{
		{"11.9.0", "11.0.0", "", true},
		{"11.9.0", "9.0.0", "10.99.99", false},
		{"10.24.0", "9.0.0", "10.99.99", true},
		{"10.99.99", "9.0.0", "10.99.99", true},   // exact upper bound is inclusive
		{"10.99.100", "9.0.0", "10.99.99", false}, // patch beyond bound excluded
		{"11.0.0", "11.0.0", "", true},
		{"11.0.0", "11.0.1", "", false},
		{"", "", "", true},
		{"", "9.0.0", "", false}, // empty version excluded when bounds set
	}
	for _, tc := range tests {
		got := versionInRange(tc.version, tc.min, tc.max)
		if got != tc.want {
			t.Errorf("versionInRange(%q, %q, %q) = %v, want %v",
				tc.version, tc.min, tc.max, got, tc.want)
		}
	}
}
