package cache

import (
	"context"
	"errors"
	"maps"
	"sync"
	"time"

	"github.com/w-h-a/flags/internal/server/clients/file"
)

var (
	ErrNotFound = errors.New("flag not found")
)

type Service struct {
	fileClient file.Client
	store      map[string]*file.Flag
	lastUpdate time.Time
	mtx        sync.RWMutex
}

func (s *Service) Flags() AllFlags {
	flags := map[string]*file.Flag{}

	s.mtx.RLock()
	maps.Copy(flags, s.store)
	s.mtx.RUnlock()

	allFlags := NewAllFlags()

	for k, flag := range flags {
		flagValue, resolutionDetails := flag.Evaluate()

		allFlags.AddFlag(k, FlagState{
			Value:   flagValue,
			Variant: resolutionDetails.Variant,
			Reason:  resolutionDetails.Reason,
		})
	}

	return allFlags
}

func (s *Service) Flag(flagKey string) (FlagState, error) {
	var flag *file.Flag
	var ok bool

	s.mtx.RLock()
	flag, ok = s.store[flagKey]
	s.mtx.RUnlock()

	if !ok {
		result := FlagState{
			Value:  nil,
			Reason: file.ReasonError,
		}

		return result, ErrNotFound
	}

	flagValue, resolutionDetails := flag.Evaluate()

	result := FlagState{
		Value:   flagValue,
		Variant: resolutionDetails.Variant,
		Reason:  resolutionDetails.Reason,
	}

	return result, nil
}

func (s *Service) RetrieveFlags() (map[string]*file.Flag, map[string]*file.Flag, error) {
	flags, err := s.fileClient.Read(context.TODO())
	if err != nil {
		return nil, nil, err
	}

	var old map[string]*file.Flag

	s.mtx.Lock()
	old = s.store
	s.store = flags
	s.lastUpdate = time.Now()
	s.mtx.Unlock()

	return old, flags, nil
}

func (s *Service) LastUpdate() time.Time {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.lastUpdate
}

func New(fileClient file.Client) *Service {
	return &Service{
		fileClient: fileClient,
		store:      map[string]*file.Flag{},
		mtx:        sync.RWMutex{},
	}
}
