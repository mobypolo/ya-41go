package customerrors

import "errors"

var (
	ErrUnsupportedType = errors.New("unsupported metric type")
	ErrInvalidValue    = errors.New("invalid metric value")
)

var (
	ErrUnknownGaugeName   = errors.New("unknown gauge metric name")
	ErrUnknownCounterName = errors.New("unknown counter metric name")
)

var ErrNotFound = errors.New("err Not Found")
var ErrNotEnoughParts = errors.New("not enough path parts")
