package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/server/helpers"
	"github.com/mobypolo/ya-41go/internal/server/storage"
	"github.com/mobypolo/ya-41go/internal/shared/dto"
	"github.com/mobypolo/ya-41go/internal/shared/logger"
	"github.com/mobypolo/ya-41go/internal/shared/utils"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

import _ "github.com/mobypolo/ya-41go/internal/server/metrics"

type metricUpdateHandler interface {
	Update(metricType storage.MetricType, name, value string) error
}

type UpdateHandler struct {
	svc metricUpdateHandler
}

func NewMetricUpdateHandler(svc metricUpdateHandler) *UpdateHandler {
	return &UpdateHandler{svc: svc}
}

func (h *UpdateHandler) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := helpers.SplitStrToChunks(r.URL.Path)

		metricType, metricName, metricValue := parts[1], parts[2], parts[3]

		if err := h.svc.Update(storage.MetricType(metricType), metricName, metricValue); err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		_, err := fmt.Fprintf(w, "Metric %s/%s updated with value %s\n", metricType, metricName, metricValue)
		if err != nil {
			log.Println(customerrors.ErrNotFound)
		}
	}
}

type metricUpdateJSONHandler interface {
	UpdateFromDTO(m dto.Metrics) error
	GetAsDTO(mType storage.MetricType, id string) (dto.Metrics, error)
}

type UpdateJSONHandler struct {
	svc metricUpdateJSONHandler
}

func NewMetricUpdateJSONHandler(svc metricUpdateJSONHandler) *UpdateJSONHandler {
	return &UpdateJSONHandler{svc: svc}
}

func (h *UpdateJSONHandler) UpdateJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m dto.Metrics
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if err := h.svc.UpdateFromDTO(m); err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		actual, err := h.svc.GetAsDTO(m.MType, m.ID)
		if err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(actual)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
	}
}

type metricUpdateJSONHandlerBatch interface {
	UpdateFromDTO(m dto.Metrics) error
}

type UpdateJSONHandlerBatch struct {
	svc metricUpdateJSONHandlerBatch
}

func NewMetricUpdateJSONHandlerBath(svc metricUpdateJSONHandlerBatch) *UpdateJSONHandlerBatch {
	return &UpdateJSONHandlerBatch{svc: svc}
}

func (h *UpdateJSONHandlerBatch) UpdateJSONHandlerBatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			logger.L().Info("cannot read body", zap.Any("body", body))
			return
		}

		var batch []dto.Metrics
		if err := json.Unmarshal(body, &batch); err != nil {
			http.Error(w, "invalid JSON format", http.StatusBadRequest)
			logger.L().Info("invalid JSON format", zap.Any("body", body))
			return
		}

		for _, metric := range batch {
			err := utils.RetryWithBackoff(r.Context(), 3, func() error {
				if err := h.svc.UpdateFromDTO(metric); err != nil {
					if isRetriablePgError(err) {
						return err
					}
					return nil
				}
				return nil
			})

			if err != nil {
				http.Error(w, fmt.Sprintf("failed to update metric %s: %v", metric.ID, err), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func isRetriablePgError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ConnectionException,
			pgerrcode.ConnectionDoesNotExist,
			pgerrcode.ConnectionFailure,
			pgerrcode.SQLClientUnableToEstablishSQLConnection,
			pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
			pgerrcode.TransactionResolutionUnknown:
			return true
		}
	}
	return false
}
