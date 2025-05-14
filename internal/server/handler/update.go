package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/server/helpers"
	"github.com/mobypolo/ya-41go/internal/server/service"
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

func UpdateHandler(service service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		parts := helpers.SplitStrToChunks(r.URL.Path)

		metricType, metricName, metricValue := parts[1], parts[2], parts[3]

		if err := service.Update(storage.MetricType(metricType), metricName, metricValue); err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		_, err := fmt.Fprintf(w, "Metric %s/%s updated with value %s\n", metricType, metricName, metricValue)
		if err != nil {
			log.Println(customerrors.ErrNotFound)
		}
	}
}

func UpdateJSONHandler(service service.MetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m dto.Metrics
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if err := service.UpdateFromDTO(m); err != nil {
			customerrors.ErrorHandler(err, w)
			return
		}

		actual, err := service.GetAsDTO(m.MType, m.ID)
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

func UpdateJSONHandlerBatch(service service.MetricService) http.HandlerFunc {
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
				if err := service.UpdateFromDTO(metric); err != nil {
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
