package customerrors

import (
	"errors"
	"net/http"
)

func ErrorHandler(err error, w http.ResponseWriter) {
	switch {
	case errors.Is(err, ErrUnsupportedType):
		http.Error(w, err.Error(), http.StatusNotImplemented)
	case errors.Is(err, ErrUnknownGaugeName) || errors.Is(err, ErrUnknownCounterName):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidValue):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
