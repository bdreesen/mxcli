// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"database/sql"
	"fmt"
	"sync"
)

// Connection holds an active database connection.
type Connection struct {
	DB     *sql.DB
	Driver DriverName
	Alias  string
}

// ConnectionInfo is a safe-to-display summary (no DSN — credential isolation).
type ConnectionInfo struct {
	Alias  string
	Driver DriverName
}

// Manager manages named database connections.
type Manager struct {
	mu    sync.Mutex
	conns map[string]*Connection
}

// NewManager creates a new connection manager.
func NewManager() *Manager {
	return &Manager{
		conns: make(map[string]*Connection),
	}
}

// Connect opens a new connection and registers it under the given alias.
func (m *Manager) Connect(driver DriverName, dsn, alias string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Close existing connection with same alias
	if old, ok := m.conns[alias]; ok {
		old.DB.Close()
	}

	db, err := sql.Open(driver.GoDriverName(), dsn)
	if err != nil {
		return fmt.Errorf("failed to open %s connection: %w", driver, err)
	}

	// Verify the connection works
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping %s database: %w", driver, err)
	}

	m.conns[alias] = &Connection{
		DB:     db,
		Driver: driver,
		Alias:  alias,
	}
	return nil
}

// Disconnect closes and removes the connection with the given alias.
func (m *Manager) Disconnect(alias string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, ok := m.conns[alias]
	if !ok {
		return fmt.Errorf("no connection with alias %q", alias)
	}

	err := conn.DB.Close()
	delete(m.conns, alias)
	return err
}

// Get returns the connection with the given alias.
func (m *Manager) Get(alias string) (*Connection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, ok := m.conns[alias]
	if !ok {
		return nil, fmt.Errorf("no connection with alias %q (use SQL CONNECT to connect first)", alias)
	}
	return conn, nil
}

// List returns connection info for all active connections (no DSN — credential isolation).
func (m *Manager) List() []ConnectionInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	infos := make([]ConnectionInfo, 0, len(m.conns))
	for _, conn := range m.conns {
		infos = append(infos, ConnectionInfo{
			Alias:  conn.Alias,
			Driver: conn.Driver,
		})
	}
	return infos
}

// CloseAll closes all active connections.
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for alias, conn := range m.conns {
		conn.DB.Close()
		delete(m.conns, alias)
	}
}
