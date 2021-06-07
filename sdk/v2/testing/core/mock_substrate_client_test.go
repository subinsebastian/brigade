package core

import (
	"testing"

	"github.com/brigadecore/brigade/sdk/v2/core"
	"github.com/stretchr/testify/require"
)

func TestMockSubstrateClient(t *testing.T) {
	require.Implements(t, (*core.SubstrateClient)(nil), &MockSubstrateClient{})
}
