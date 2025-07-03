package export

import (
	"context"
	"log/slog"
	"sync"

	"github.com/w-h-a/flags/internal/server/clients/exporter"
)

type Service struct {
	exportClient exporter.Exporter
	store        []Event
	mtx          sync.RWMutex
}

func (s *Service) Add(event Event) {
	if int64(len(s.store)) >= MaxEventsInMemory {
		s.Flush()
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.store = append(s.store, event)
}

func (s *Service) Flush() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if len(s.store) == 0 {
		return
	}

	records := []exporter.Record{}

	for _, event := range s.store {
		record := exporter.Record{
			CreationDate: event.CreationDate,
			Key:          event.Key,
			Value:        event.Value,
			Variant:      event.Variant,
			Reason:       event.Reason,
			ErrorCode:    event.ErrorCode,
		}

		records = append(records, record)
	}

	if err := s.exportClient.Export(context.TODO(), records); err != nil {
		slog.WarnContext(context.TODO(), "failed to export evaluation event", "error", err)
		return
	}

	s.store = make([]Event, 0)
}

func (s *Service) Close() {
	s.Flush()
}

func New(exportClient exporter.Exporter) *Service {
	return &Service{
		exportClient: exportClient,
		store:        []Event{},
		mtx:          sync.RWMutex{},
	}
}
