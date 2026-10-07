// Package platform contains small runtime ports shared by service modules.
// It does not import service code or concrete infrastructure adapters.
package platform

import "time"

// Clock supplies wall time to application code. The function form keeps the
// existing Server.now injection seam source-compatible and easy to fake.
type Clock func() time.Time

// IDGenerator creates ordinary resource/database identifiers. Do not use it
// for bearer tokens, setup codes, nonces, or cryptographic keys.
type IDGenerator interface {
	NewID() string
}

// IDGeneratorFunc adapts an ordinary ID function to the shared port.
type IDGeneratorFunc func() string

func (f IDGeneratorFunc) NewID() string { return f() }

var _ IDGenerator = IDGeneratorFunc(nil)
