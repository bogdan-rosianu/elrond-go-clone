package persister

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUint64_WrongTypeShouldReturn0(t *testing.T) {
	t.Parallel()

	value := GetUint64(25.0)

	require.Equal(t, uint64(0), value)
}

func TestGetUint64_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedValue := uint64(100)
	value := GetUint64(uint64(100))

	require.Equal(t, expectedValue, value)
}

func TestGetString_WrongTypeShouldReturnEmptyString(t *testing.T) {
	t.Parallel()

	value := GetString(24)

	require.Equal(t, "", value)
}

func TestGetString_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedValue := "tesstt"
	value := GetString(expectedValue)

	require.Equal(t, expectedValue, value)
}
