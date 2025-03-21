package storage_test

import (
	"errors"
	"github.com/mobypolo/ya-41go/internal/customerrors"
	"github.com/mobypolo/ya-41go/internal/repositories"
	"github.com/mobypolo/ya-41go/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

type dummyProcessor struct{}

func (d dummyProcessor) ValidateName(name string) error {
	if name == "valid" {
		return nil
	}
	return errors.New("invalid name")
}

func (d dummyProcessor) ParseValue(value string) (any, error) {
	return value, nil
}

func (d dummyProcessor) Update(_ repositories.MetricsRepository, _ string, _ any) error {
	return nil
}

func TestRegisterAndGetProcessor(t *testing.T) {
	storage.RegisterProcessor("dummy", dummyProcessor{})

	p, err := storage.GetProcessor("dummy")
	assert.NoError(t, err)
	assert.NotNil(t, p)

	assert.Equal(t, nil, p.ValidateName("valid"))
	assert.Error(t, p.ValidateName("invalid"))
}

func TestGetProcessor_NotRegistered(t *testing.T) {
	_, err := storage.GetProcessor("nonexistent")
	assert.ErrorIs(t, err, customerrors.ErrUnsupportedType)
}
