package cache

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sort"
	"sync"
	"time"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/flags/internal/server/config"
)

var (
	ErrNotFound = errors.New("flag not found")
)

type Service struct {
	readClient reader.Reader
	store      map[string]*flags.Flag
	lastUpdate time.Time
	mtx        sync.RWMutex
}

func (s *Service) EvaluateFlag(ctx context.Context, flagKey string, evalCtx map[string]any) (FlagState, error) {
	var flag *flags.Flag
	var ok bool

	s.mtx.RLock()
	flag, ok = s.store[flagKey]
	s.mtx.RUnlock()

	if !ok {
		result := FlagState{
			Key:          flagKey,
			ErrorCode:    flags.ErrorNotFound,
			ErrorMessage: fmt.Sprintf("flag for key '%s' does not exist", flagKey),
		}

		return result, ErrNotFound
	}

	flagValue, resolutionDetails := flag.Evaluate(evalCtx)

	result := FlagState{
		Key:     flagKey,
		Value:   flagValue,
		Variant: resolutionDetails.Variant,
		Reason:  resolutionDetails.Reason,
	}

	return result, nil
}

func (s *Service) EvaluateFlags(ctx context.Context) AllFlags {
	flags := map[string]*flags.Flag{}

	s.mtx.RLock()
	maps.Copy(flags, s.store)
	s.mtx.RUnlock()

	allFlags := NewAllFlags()

	for k, flag := range flags {
		flagValue, resolutionDetails := flag.Evaluate(map[string]any{})

		allFlags.AddFlag(FlagState{
			Key:     k,
			Value:   flagValue,
			Variant: resolutionDetails.Variant,
			Reason:  resolutionDetails.Reason,
		})
	}

	sort.Slice(allFlags.Flags, func(i, j int) bool {
		return allFlags.Flags[i].Key < allFlags.Flags[j].Key
	})

	return allFlags
}

func (s *Service) RetrieveFlags() (map[string]*flags.Flag, map[string]*flags.Flag, error) {
	bs, err := s.readClient.Read(context.TODO())
	if err != nil {
		return nil, nil, err
	}

	new, err := flags.Factory(bs, config.FlagFormat())
	if err != nil {
		return nil, nil, err
	}

	var old map[string]*flags.Flag

	s.mtx.Lock()
	old = s.store
	s.store = new
	s.lastUpdate = time.Now()
	s.mtx.Unlock()

	return old, new, nil
}

func (s *Service) LastUpdate() time.Time {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.lastUpdate
}

func New(readClient reader.Reader) *Service {
	return &Service{
		readClient: readClient,
		store:      map[string]*flags.Flag{},
		mtx:        sync.RWMutex{},
	}
}
