// SPDX-License-Identifier: Apache-2.0

package mpr

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// =============================================================================
// removeRolesFromAccessRule — unit tests for multi-role handling
// =============================================================================

func makeAccessRule(roleNames ...string) bson.D {
	roles := bson.A{int32(1)} // Mendix array marker
	for _, r := range roleNames {
		roles = append(roles, r)
	}
	return bson.D{
		{Key: "$Type", Value: "DomainModels$AccessRule"},
		{Key: "AllowedModuleRoles", Value: roles},
		{Key: "AllowCreate", Value: true},
	}
}

func getRoleNames(rule bson.D) []string {
	for _, f := range rule {
		if f.Key != "AllowedModuleRoles" {
			continue
		}
		arr, ok := f.Value.(bson.A)
		if !ok {
			return nil
		}
		var names []string
		for _, item := range arr {
			if s, ok := item.(string); ok {
				names = append(names, s)
			}
		}
		return names
	}
	return nil
}

func TestRemoveRolesFromAccessRule_SingleRole_ExactMatch(t *testing.T) {
	rule := makeAccessRule("Mod.RoleA")
	keep, modified := removeRolesFromAccessRule(rule, map[string]bool{"Mod.RoleA": true})
	if keep {
		t.Error("expected rule to be deleted (no roles remaining)")
	}
	if !modified {
		t.Error("expected modified=true")
	}
}

func TestRemoveRolesFromAccessRule_MultiRole_RemoveOne(t *testing.T) {
	rule := makeAccessRule("Mod.User", "Mod.Admin")
	keep, modified := removeRolesFromAccessRule(rule, map[string]bool{"Mod.User": true})
	if !keep {
		t.Error("expected rule to be kept (Admin still present)")
	}
	if !modified {
		t.Error("expected modified=true")
	}
	names := getRoleNames(rule)
	if len(names) != 1 || names[0] != "Mod.Admin" {
		t.Errorf("expected [Mod.Admin], got %v", names)
	}
}

func TestRemoveRolesFromAccessRule_MultiRole_RemoveAll(t *testing.T) {
	rule := makeAccessRule("Mod.User", "Mod.Admin")
	keep, modified := removeRolesFromAccessRule(rule, map[string]bool{"Mod.User": true, "Mod.Admin": true})
	if keep {
		t.Error("expected rule to be deleted (no roles remaining)")
	}
	if !modified {
		t.Error("expected modified=true")
	}
}

func TestRemoveRolesFromAccessRule_NoMatch(t *testing.T) {
	rule := makeAccessRule("Mod.User", "Mod.Admin")
	keep, modified := removeRolesFromAccessRule(rule, map[string]bool{"Mod.Other": true})
	if !keep {
		t.Error("expected rule to be kept")
	}
	if modified {
		t.Error("expected modified=false")
	}
	names := getRoleNames(rule)
	if len(names) != 2 {
		t.Errorf("expected 2 roles unchanged, got %v", names)
	}
}

func TestRemoveRolesFromAccessRule_NoAllowedModuleRoles(t *testing.T) {
	rule := bson.D{
		{Key: "$Type", Value: "DomainModels$AccessRule"},
		{Key: "AllowCreate", Value: true},
	}
	keep, modified := removeRolesFromAccessRule(rule, map[string]bool{"Mod.User": true})
	if !keep {
		t.Error("expected rule to be kept (no AllowedModuleRoles field)")
	}
	if modified {
		t.Error("expected modified=false")
	}
}

func TestRemoveRolesFromAccessRule_ThreeRoles_RemoveMiddle(t *testing.T) {
	rule := makeAccessRule("Mod.A", "Mod.B", "Mod.C")
	keep, modified := removeRolesFromAccessRule(rule, map[string]bool{"Mod.B": true})
	if !keep {
		t.Error("expected rule to be kept")
	}
	if !modified {
		t.Error("expected modified=true")
	}
	names := getRoleNames(rule)
	if len(names) != 2 || names[0] != "Mod.A" || names[1] != "Mod.C" {
		t.Errorf("expected [Mod.A, Mod.C], got %v", names)
	}
}
