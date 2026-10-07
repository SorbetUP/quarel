// Package ids adapts Quarel's existing CSPRNG-backed ordinary IDs to the
// platform IDGenerator port.
package ids

import (
	"github.com/anlekg/quarel/internal/platform"
	"github.com/anlekg/quarel/internal/secret"
)

type Generator struct{}

func (Generator) NewID() string { return secret.NewID() }

var _ platform.IDGenerator = Generator{}
