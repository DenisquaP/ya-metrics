package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	yametrics "github.com/DenisquaP/ya-metrics/internal/server/yaMetrics"
)

func TestInitHandlers(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	defer logger.Sync()

	suggared := logger.Sugar()
	mem := yametrics.NewMemStorage("mem.json")

	r := NewRouterWithMiddlewares(suggared, mem)

	assert.NotEmpty(t, r)
}
