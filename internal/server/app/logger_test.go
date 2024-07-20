package app

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, err := newLogger()
	require.NoError(t, err)
	require.NotNil(t, logger)
}
