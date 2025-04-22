package storage

import (
	"encoding/json"
	"github.com/mobypolo/ya-41go/internal/shared/logger"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

type PersistentStorage struct {
	*MemStorage
	filePath      string
	autoStore     bool
	storeInterval time.Duration
	quitChan      chan struct{}
	once          sync.Once
}

func NewPersistentStorage(filePath string, storeInterval time.Duration, restore bool) *PersistentStorage {
	ps := &PersistentStorage{
		MemStorage:    NewMemStorage(),
		filePath:      filePath,
		autoStore:     storeInterval > 0,
		storeInterval: storeInterval,
		quitChan:      make(chan struct{}),
	}
	if restore {
		_ = ps.LoadFromDisk()
	}
	if ps.autoStore {
		go ps.startAutoSave()
	}
	return ps
}

func (s *PersistentStorage) startAutoSave() {
	if s.storeInterval == 0 {
		return
	}

	ticker := time.NewTicker(s.storeInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := s.SaveToDisk()
			if err != nil {
				logger.L().Info("SaveToDisk error", zap.Error(err))
				return
			}
		case <-s.quitChan:
			return
		}
	}
}

func (s *PersistentStorage) Stop() {
	s.once.Do(func() {
		close(s.quitChan)
		err := s.SaveToDisk()
		if err != nil {
			logger.L().Info("file save of exit store close error", zap.Error(err))
			return
		}
	})
}

func (s *PersistentStorage) SaveToDisk() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := map[string]interface{}{
		"gauges":   s.Gauges,
		"counters": s.Counters,
	}
	f, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.L().Info("file store close error", zap.Error(err))
		}
	}(f)
	return json.NewEncoder(f).Encode(data)
}

func (s *PersistentStorage) LoadFromDisk() error {
	f, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.L().Info("file load store close error", zap.Error(err))
		}
	}(f)

	var data struct {
		Gauges   map[string]float64 `json:"gauges"`
		Counters map[string]int64   `json:"counters"`
	}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.Gauges = data.Gauges
	s.Counters = data.Counters
	return nil
}

func (s *PersistentStorage) UpdateGauge(name string, value float64) error {
	err := s.MemStorage.UpdateGauge(name, value)
	if err != nil {
		return err
	}
	if s.storeInterval == 0 {
		return s.SaveToDisk()
	}
	return nil
}

func (s *PersistentStorage) UpdateCounter(name string, delta int64) error {
	err := s.MemStorage.UpdateCounter(name, delta)
	if err != nil {
		return err
	}
	if s.storeInterval == 0 {
		return s.SaveToDisk()
	}
	return nil
}
