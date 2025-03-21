package metrics_test

import (
	"testing"

	"github.com/mobypolo/ya-41go/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func TestGaugeProcessor_ValidateName(t *testing.T) {
	p := metrics.NewGaugeProcessor()

	assert.NoError(t, p.ValidateName("temperature"))
	assert.Error(t, p.ValidateName("unknown"))
}

func TestGaugeProcessor_ParseValue(t *testing.T) {
	p := metrics.NewGaugeProcessor()

	val, err := p.ParseValue("42.5")
	assert.NoError(t, err)
	assert.Equal(t, 42.5, val)

	_, err = p.ParseValue("abc")
	assert.Error(t, err)
}
