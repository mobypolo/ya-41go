//go:build agent
// +build agent

package agent_test

import (
	"errors"
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSource struct {
	name  string
	typ   storage.MetricType
	value interface{}
	err   error
}

func (m mockSource) Name() string                  { return m.name }
func (m mockSource) Type() storage.MetricType      { return m.typ }
func (m mockSource) Collect() (interface{}, error) { return m.value, m.err }

func TestRegisterAndCollectAll(t *testing.T) {
	// Сбросим глобальный список
	agent.ResetSourcesForTest()

	src1 := mockSource{name: "m1", typ: storage.GaugeType, value: 123.45}
	src2 := mockSource{name: "m2", typ: storage.CounterType, value: 5}
	srcErr := mockSource{name: "bad", typ: storage.GaugeType, err: errors.New("fail")}

	agent.Register(src1)
	agent.Register(src2)
	agent.Register(srcErr)

	metrics, err := agent.CollectAll()
	assert.NoError(t, err)
	assert.Len(t, metrics, 2)

	assert.Equal(t, "m1", metrics[0].Name)
	assert.Equal(t, storage.GaugeType, metrics[0].Type)
	assert.Equal(t, 123.45, metrics[0].Value)

	assert.Equal(t, "m2", metrics[1].Name)
	assert.Equal(t, storage.CounterType, metrics[1].Type)
	assert.Equal(t, 5, metrics[1].Value)
}
