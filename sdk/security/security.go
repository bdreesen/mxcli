// SPDX-License-Identifier: Apache-2.0

// Package security provides types for Mendix project security configuration.
package security

import (
	"github.com/mendixlabs/mxcli/model"
)

// ProjectSecurity represents the project-level security configuration.
type ProjectSecurity struct {
	model.BaseElement
	SecurityLevel      string          `json:"securityLevel"`
	AdminUserName      string          `json:"adminUserName"`
	AdminPassword      string          `json:"adminPassword"`
	AdminUserRole      string          `json:"adminUserRole"`
	CheckSecurity      bool            `json:"checkSecurity"`
	StrictMode         bool            `json:"strictMode"`
	StrictPageUrlCheck bool            `json:"strictPageUrlCheck"`
	EnableDemoUsers    bool            `json:"enableDemoUsers"`
	EnableGuestAccess  bool            `json:"enableGuestAccess"`
	GuestUserRole      string          `json:"guestUserRole,omitempty"`
	UserRoles          []*UserRole     `json:"userRoles,omitempty"`
	DemoUsers          []*DemoUser     `json:"demoUsers,omitempty"`
	PasswordPolicy     *PasswordPolicy `json:"passwordPolicy,omitempty"`
}

// UserRole represents an application-level user role that combines module roles.
type UserRole struct {
	model.BaseElement
	Name                    string   `json:"name"`
	Description             string   `json:"description,omitempty"`
	ModuleRoles             []string `json:"moduleRoles,omitempty"`
	ManageAllRoles          bool     `json:"manageAllRoles"`
	ManageUsersWithoutRoles bool     `json:"manageUsersWithoutRoles"`
	ManageableRoles         []string `json:"manageableRoles,omitempty"`
	CheckSecurity           bool     `json:"checkSecurity"`
}

// DemoUser represents a demo user for development/testing.
type DemoUser struct {
	model.BaseElement
	UserName  string   `json:"userName"`
	Password  string   `json:"password"`
	Entity    string   `json:"entity"`
	UserRoles []string `json:"userRoles,omitempty"`
}

// PasswordPolicy represents the password policy settings.
type PasswordPolicy struct {
	model.BaseElement
	MinimumLength    int  `json:"minimumLength"`
	RequireDigit     bool `json:"requireDigit"`
	RequireMixedCase bool `json:"requireMixedCase"`
	RequireSymbol    bool `json:"requireSymbol"`
}

// ModuleSecurity represents the security configuration for a module.
type ModuleSecurity struct {
	model.BaseElement
	ContainerID model.ID      `json:"containerId"`
	ModuleRoles []*ModuleRole `json:"moduleRoles,omitempty"`
}

// GetContainerID returns the ID of the containing module.
func (ms *ModuleSecurity) GetContainerID() model.ID {
	return ms.ContainerID
}

// ModuleRole represents a module-level security role.
type ModuleRole struct {
	model.BaseElement
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// GetName returns the module role's name.
func (mr *ModuleRole) GetName() string {
	return mr.Name
}

// SecurityLevel constants matching BSON SecurityLevel enum values.
const (
	SecurityLevelOff        = "CheckNothing"
	SecurityLevelPrototype  = "CheckFormsAndMicroflows"
	SecurityLevelProduction = "CheckEverything"
)

// SecurityLevelDisplay returns a human-friendly name for a security level.
func SecurityLevelDisplay(level string) string {
	switch level {
	case SecurityLevelOff:
		return "Off"
	case SecurityLevelPrototype:
		return "Prototype / demo"
	case SecurityLevelProduction:
		return "Production"
	default:
		return level
	}
}
