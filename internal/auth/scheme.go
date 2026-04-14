// SPDX-License-Identifier: Apache-2.0

package auth

// hostSchemes maps Mendix platform hostnames to the auth scheme they require.
//
// Add a host here when wiring a new platform API consumer. If a request
// targets an unlisted host, the client returns an error rather than silently
// sending a token to the wrong service.
var hostSchemes = map[string]Scheme{
	"appstore.home.mendix.com": SchemePAT,
}

// SchemeForHost returns the auth scheme required by the given hostname.
// Returns false if the host is not a known Mendix platform endpoint.
func SchemeForHost(host string) (Scheme, bool) {
	s, ok := hostSchemes[host]
	return s, ok
}
