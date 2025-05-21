package notify

import (
	"context"
	"sync"

	"github.com/google/go-cmp/cmp"
	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/notifier"
	"github.com/w-h-a/pkg/telemetry/log"
)

type Service struct {
	notifyClient notifier.Notifier
	waitGroup    *sync.WaitGroup
}

func (s *Service) Notify(old, new map[string]*flags.Flag) {
	diff := s.diff(old, new)

	if !diff.HasDiff() {
		return
	}

	s.waitGroup.Add(1)

	go func() {
		defer s.waitGroup.Done()

		err := s.notifyClient.Notify(context.TODO(), diff)
		if err != nil {
			log.Errorf("notify service failed to send message: %v", err)
		}
	}()
}

func (s *Service) Close() {
	s.waitGroup.Wait()
}

func (s *Service) diff(old, new map[string]*flags.Flag) notifier.Diff {
	diff := notifier.Diff{
		Deleted: map[string]*flags.Flag{},
		Added:   map[string]*flags.Flag{},
		Updated: map[string]notifier.DiffUpdated{},
	}

	for k := range old {
		nf, ok := new[k]
		of := old[k]

		// if it's not in new, it needs to be shown as deleted
		if !ok {
			diff.Deleted[k] = of
			continue
		}

		// if it's not equal, it needs to be shown as updated
		if !cmp.Equal(of, nf) {
			diff.Updated[k] = notifier.DiffUpdated{
				Before: of,
				After:  nf,
			}
		}
	}

	for k := range new {
		// if not in old, it needs to be shown as added
		if _, ok := old[k]; !ok {
			f := new[k]
			diff.Added[k] = f
		}
	}

	return diff
}

func New(notifyClient notifier.Notifier) *Service {
	return &Service{
		notifyClient: notifyClient,
		waitGroup:    &sync.WaitGroup{},
	}
}
