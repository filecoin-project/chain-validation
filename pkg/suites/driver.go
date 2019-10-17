package suites

import (
	"github.com/filecoin-project/chain-validation/pkg/chain"
	"github.com/filecoin-project/chain-validation/pkg/state"
)

// Driver wraps up all the implementation-specific integration points.
type Driver interface {
	state.Factory
	chain.MessageFactory
	chain.Applier
}
