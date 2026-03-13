// SPDX-License-Identifier: Apache-2.0

package ast

// ============================================================================
// Connection Statements
// ============================================================================

// ConnectStmt represents: CONNECT LOCAL 'path' or CONNECT TO FILESYSTEM 'path'
type ConnectStmt struct {
	Path string
}

func (s *ConnectStmt) isStatement() {}

// DisconnectStmt represents: DISCONNECT
type DisconnectStmt struct{}

func (s *DisconnectStmt) isStatement() {}

// StatusStmt represents: STATUS or SHOW STATUS
type StatusStmt struct{}

func (s *StatusStmt) isStatement() {}
