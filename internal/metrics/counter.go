package metrics

import (
	"github.com/mobypolo/ya-41go/internal/customerrors"
	"strconv"

	"github.com/mobypolo/ya-41go/internal/storage"
)

type counterProcessor struct{}

func (c counterProcessor) ValidateName(name string) error {
	_, err := storage.ParseCounterMetric(name)
	return err
}

func (c counterProcessor) ParseValue(value string) (any, error) {
	return strconv.ParseInt(value, 10, 64)
}

func (c counterProcessor) Update(s storage.Storage, name string, value any) error {
	v, ok := value.(int64)
	if !ok {
		return customerrors.ErrInvalidValue
	}
	return s.UpdateCounter(name, v)
}
