package exportevals

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/flags/internal/server/clients/report"
	"github.com/w-h-a/flags/internal/server/clients/report/mock"
	"github.com/w-h-a/flags/internal/server/services/export"
	"github.com/w-h-a/flags/tests/unit"
)

func TestExportEvals_FlushWithTime(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	unit.SetLogger()

	reportClient := mock.NewReportClient(
		report.WithDir("any"),
	)

	exportService := export.New(reportClient)

	errCh := make(chan error, 1)
	exportStop := make(chan struct{})

	go func() {
		errCh <- exportReports(exportService, exportStop, 10*time.Millisecond)
	}()

	want := []export.Event{
		{
			CreationDate: time.Now().Unix(),
			Key:          "random-key",
			Value:        "YO",
			Variant:      "default",
			Reason:       "DEFAULT",
		},
	}

	for _, event := range want {
		exportService.Add(event)
	}

	time.Sleep(500 * time.Millisecond)

	c := reportClient.(*mock.Client)

	got := c.Records()

	for i, event := range want {
		require.Equal(t, event.CreationDate, got[i].CreationDate)
		require.Equal(t, event.Key, got[i].Key)
	}

	close(exportStop)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(30 * time.Second):
	}
}

func TestExportEvals_FlushWithClose(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	unit.SetLogger()

	reportClient := mock.NewReportClient(
		report.WithDir("any"),
	)

	exportService := export.New(reportClient)

	errCh := make(chan error, 1)
	exportStop := make(chan struct{})

	go func() {
		errCh <- exportReports(exportService, exportStop, 10*time.Minute)
	}()

	want := []export.Event{
		{
			CreationDate: time.Now().Unix(),
			Key:          "random-key",
			Value:        "YO",
			Variant:      "default",
			Reason:       "DEFAULT",
		},
	}

	for _, event := range want {
		exportService.Add(event)
	}

	time.Sleep(500 * time.Millisecond)

	c := reportClient.(*mock.Client)

	got := c.Records()

	require.Equal(t, 0, len(got))

	close(exportStop)

	<-time.After(500 * time.Millisecond)

	got = c.Records()

	for i, event := range want {
		require.Equal(t, event.CreationDate, got[i].CreationDate)
		require.Equal(t, event.Key, got[i].Key)
	}

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(30 * time.Second):
	}
}

func exportReports(
	exportService *export.Service,
	stop chan struct{},
	dur time.Duration,
) error {
	ticker := time.NewTicker(dur)

	for {
		select {
		case <-ticker.C:
			exportService.Flush()
		case <-stop:
			ticker.Stop()
			exportService.Close()
			return nil
		}
	}
}
