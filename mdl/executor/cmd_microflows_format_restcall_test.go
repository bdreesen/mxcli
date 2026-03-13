// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"testing"

	"github.com/mendixlabs/mxcli/sdk/microflows"
)

// =============================================================================
// formatRestCallAction
// =============================================================================

func TestFormatRestCallAction_GET(t *testing.T) {
	e := newTestExecutor()
	action := &microflows.RestCallAction{
		HttpConfiguration: &microflows.HttpConfiguration{
			HttpMethod:       microflows.HttpMethodGet,
			LocationTemplate: "https://api.example.com/orders",
		},
		ResultHandling: &microflows.ResultHandlingString{VariableName: "Response"},
	}
	got := e.formatRestCallAction(action)
	if got == "" {
		t.Fatal("expected non-empty output")
	}
	assertContains(t, got, "REST CALL GET")
	assertContains(t, got, "'https://api.example.com/orders'")
	assertContains(t, got, "$Response = ")
	assertContains(t, got, "RETURNS String")
}

func TestFormatRestCallAction_POST_CustomBody(t *testing.T) {
	e := newTestExecutor()
	action := &microflows.RestCallAction{
		HttpConfiguration: &microflows.HttpConfiguration{
			HttpMethod:       microflows.HttpMethodPost,
			LocationTemplate: "https://api.example.com/orders",
		},
		RequestHandling: &microflows.CustomRequestHandling{
			Template: `{"name": "test"}`,
		},
		ResultHandling: &microflows.ResultHandlingNone{},
	}
	got := e.formatRestCallAction(action)
	assertContains(t, got, "REST CALL POST")
	assertContains(t, got, "BODY '{\"name\": \"test\"}'")
	assertContains(t, got, "RETURNS Nothing")
}

func TestFormatRestCallAction_WithHeaders(t *testing.T) {
	e := newTestExecutor()
	action := &microflows.RestCallAction{
		HttpConfiguration: &microflows.HttpConfiguration{
			HttpMethod:       microflows.HttpMethodGet,
			LocationTemplate: "https://api.example.com",
			CustomHeaders: []*microflows.HttpHeader{
				{Name: "Authorization", Value: "'Bearer ' + $Token"},
			},
		},
		ResultHandling: &microflows.ResultHandlingString{VariableName: "Resp"},
	}
	got := e.formatRestCallAction(action)
	assertContains(t, got, "HEADER 'Authorization' = 'Bearer ' + $Token")
}

func TestFormatRestCallAction_WithAuth(t *testing.T) {
	e := newTestExecutor()
	action := &microflows.RestCallAction{
		HttpConfiguration: &microflows.HttpConfiguration{
			HttpMethod:        microflows.HttpMethodGet,
			LocationTemplate:  "https://api.example.com",
			UseAuthentication: true,
			Username:          "'admin'",
			Password:          "'secret'",
		},
		ResultHandling: &microflows.ResultHandlingString{},
	}
	got := e.formatRestCallAction(action)
	assertContains(t, got, "AUTH BASIC 'admin' PASSWORD 'secret'")
}

func TestFormatRestCallAction_WithTimeout(t *testing.T) {
	e := newTestExecutor()
	action := &microflows.RestCallAction{
		HttpConfiguration: &microflows.HttpConfiguration{
			HttpMethod:       microflows.HttpMethodGet,
			LocationTemplate: "https://api.example.com",
		},
		TimeoutExpression: "30",
		ResultHandling:    &microflows.ResultHandlingString{},
	}
	got := e.formatRestCallAction(action)
	assertContains(t, got, "TIMEOUT 30")
}
