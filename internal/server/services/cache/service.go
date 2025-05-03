package cache

import (
	"context"
	"maps"
	"sync"
	"time"

	"github.com/w-h-a/flags/internal/server/clients/file"
)

type Service struct {
	fileClient   file.Client
	store        map[string]file.Flag
	latestUpdate time.Time
	mtx          sync.RWMutex
}

func (s *Service) Flags() AllFlags {
	flags := map[string]file.Flag{}

	s.mtx.RLock()
	maps.Copy(flags, s.store)
	s.mtx.RUnlock()

	allFlags := NewAllFlags()

	for k, flag := range flags {
		flagValue, resolutionDetails := flag.Evaluate(file.EvaluationContext{
			DefaultValue: nil,
		})

		allFlags.AddFlag(k, FlagState{
			Value:     flagValue,
			Variant:   resolutionDetails.Variant,
			Reason:    resolutionDetails.Reason,
			Timestamp: time.Now().Unix(),
		})
	}

	return allFlags
}

func (s *Service) RetrieveFlags() (map[string]file.Flag, map[string]file.Flag, error) {
	flags, err := s.fileClient.Read(context.TODO())
	if err != nil {
		return nil, nil, err
	}

	var old map[string]file.Flag

	s.mtx.Lock()
	old = s.store
	s.store = flags
	s.latestUpdate = time.Now()
	s.mtx.Unlock()

	return old, flags, nil
}

func (s *Service) LatestUpdate() time.Time {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.latestUpdate
}

func New(fileClient file.Client) *Service {
	return &Service{
		fileClient: fileClient,
		store:      map[string]file.Flag{},
		mtx:        sync.RWMutex{},
	}
}
