package roundActivation_test

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/process/roundActivation"
	"github.com/stretchr/testify/require"
)

func TestRoundActivation_IsEnabledInRound_IsEnabled(t *testing.T) {
	t.Parallel()

	ra, err := roundActivation.NewRoundActivation(config.RoundConfig{})
	require.Nil(t, err)
	require.False(t, ra.IsEnabledInRound("Fix1", 100))
	require.False(t, ra.IsEnabledInRound("Fix2", 1000))
	require.False(t, ra.IsEnabled("Fix1"))
	require.False(t, ra.IsEnabled("Fix2"))
}

func TestRoundActivation_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	ra, _ := roundActivation.NewRoundActivation(config.RoundConfig{})
	require.False(t, ra.IsInterfaceNil())
}
