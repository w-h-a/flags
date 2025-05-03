package notify

import (
	"sync"

	"github.com/google/go-cmp/cmp"
	"github.com/w-h-a/flags/internal/server/clients/file"
	"github.com/w-h-a/flags/internal/server/clients/message"
	"github.com/w-h-a/pkg/telemetry/log"
)

type Service struct {
	notifiers []message.Client
	waitGroup *sync.WaitGroup
}

func (s *Service) Notify(old, new map[string]*file.Flag) {
	diff := s.diff(old, new)

	if !diff.HasDiff() {
		return
	}

	for _, n := range s.notifiers {
		s.waitGroup.Add(1)

		go func() {
			err := n.Send(diff, s.waitGroup)
			if err != nil {
				log.Errorf("notify service failed to send message: %v", err)
			}
		}()
	}
}

func (s *Service) Close() {
	s.waitGroup.Wait()
}

func (s *Service) diff(old, new map[string]*file.Flag) message.Diff {
	diff := message.Diff{
		Deleted: map[string]*file.Flag{},
		Added:   map[string]*file.Flag{},
		Updated: map[string]message.DiffUpdated{},
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
			diff.Updated[k] = message.DiffUpdated{
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

func New(notifiers ...message.Client) *Service {
	return &Service{
		notifiers: notifiers,
		waitGroup: &sync.WaitGroup{},
	}
}
