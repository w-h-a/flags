package export

import (
	"context"
	"sync"

	"github.com/w-h-a/flags/internal/server/clients/report"
	"github.com/w-h-a/pkg/telemetry/log"
)

type Service struct {
	reportClient report.Client
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

	records := []report.Record{}

	for _, event := range s.store {
		record := report.Record{
			CreationDate: event.CreationDate,
			Key:          event.Key,
			Value:        event.Value,
			Variant:      event.Variant,
			Reason:       event.Reason,
			ErrorCode:    event.ErrorCode,
		}

		records = append(records, record)
	}

	if err := s.reportClient.Create(context.TODO(), records); err != nil {
		log.Warnf("failed to export evaluation event: %v", err)
		return
	}

	s.store = make([]Event, 0)
}

func (s *Service) Close() {
	s.Flush()
}

func New(reportClient report.Client) *Service {
	return &Service{
		reportClient: reportClient,
		store:        []Event{},
		mtx:          sync.RWMutex{},
	}
}
