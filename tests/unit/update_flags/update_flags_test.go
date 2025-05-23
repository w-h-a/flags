package updateflags

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server"
	mocknotifier "github.com/w-h-a/flags/internal/server/clients/notifier/mock"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	mockreader "github.com/w-h-a/flags/internal/server/clients/reader/mock"
	"github.com/w-h-a/flags/internal/server/services/cache"
	"github.com/w-h-a/flags/internal/server/services/notify"
	"github.com/w-h-a/flags/tests/unit"
)

func TestUpdateFlags_NoChange(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	unit.SetLogger("update")

	readClient := mockreader.NewReader(
		reader.WithLocation("any"),
		mockreader.WithInitialFlags(
			map[string]*flags.Flag{
				"hello": {
					Variants: map[string]any{
						"default": "default",
					},
				},
			},
		),
		mockreader.WithUpdatedFlags(
			map[string]*flags.Flag{
				"hello": {
					Variants: map[string]any{
						"default": "default",
					},
				},
			},
		),
	)

	cacheService := cache.New(readClient)

	old, new, err := cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 0, len(old))
	require.Equal(t, 1, len(new))

	r := readClient.(*mockreader.Client)
	callCount := r.CallCount()
	require.Equal(t, 1, callCount)

	notifyClient := mocknotifier.NewNotifier()

	notifyService := notify.New(notifyClient)

	errCh := make(chan error, 1)
	updateStop := make(chan struct{})

	go func() {
		errCh <- server.UpdateCache(cacheService, notifyService, updateStop, 10*time.Millisecond)
	}()

	time.Sleep(500 * time.Millisecond)

	_, new, err = cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 1, len(new))

	n := notifyClient.(*mocknotifier.Client)
	wasCalled := n.WasCalled()
	require.False(t, wasCalled)

	close(updateStop)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(30 * time.Second):
	}
}

func TestUpdateFlags_UpdatedFlags(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	unit.SetLogger("update")

	readClient := mockreader.NewReader(
		reader.WithLocation("any"),
		mockreader.WithInitialFlags(
			map[string]*flags.Flag{
				"flag1": {
					Disabled: unit.Bool(true),
					Variants: map[string]any{
						"default": "default",
					},
				},
			},
		),
		mockreader.WithUpdatedFlags(
			map[string]*flags.Flag{
				"flag1": {
					Disabled: unit.Bool(false),
					Variants: map[string]any{
						"default": "default",
					},
				},
				"flag2": {
					Variants: map[string]any{
						"default": "default",
					},
				},
			},
		),
	)

	cacheService := cache.New(readClient)

	old, new, err := cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 0, len(old))
	require.Equal(t, 1, len(new))

	r := readClient.(*mockreader.Client)
	callCount := r.CallCount()
	require.Equal(t, 1, callCount)

	notifyClient := mocknotifier.NewNotifier()

	notifyService := notify.New(notifyClient)

	errCh := make(chan error, 1)
	updateStop := make(chan struct{})

	go func() {
		errCh <- server.UpdateCache(cacheService, notifyService, updateStop, 10*time.Millisecond)
	}()

	time.Sleep(500 * time.Millisecond)

	_, new, err = cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 2, len(new))

	n := notifyClient.(*mocknotifier.Client)
	wasCalled := n.WasCalled()
	require.True(t, wasCalled)

	close(updateStop)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(30 * time.Second):
	}
}

func TestUpdateFlags_NewFlags(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	unit.SetLogger("update")

	readClient := mockreader.NewReader(
		reader.WithLocation("any"),
		mockreader.WithInitialFlags(
			map[string]*flags.Flag{},
		),
		mockreader.WithUpdatedFlags(
			map[string]*flags.Flag{
				"flag1": {
					Disabled: unit.Bool(false),
					Variants: map[string]any{
						"default": "default",
					},
				},
				"flag2": {
					Variants: map[string]any{
						"default": "default",
					},
				},
			},
		),
	)

	cacheService := cache.New(readClient)

	old, new, err := cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 0, len(old))
	require.Equal(t, 0, len(new))

	r := readClient.(*mockreader.Client)
	callCount := r.CallCount()
	require.Equal(t, 1, callCount)

	notifyClient := mocknotifier.NewNotifier()

	notifyService := notify.New(notifyClient)

	errCh := make(chan error, 1)
	updateStop := make(chan struct{})

	go func() {
		errCh <- server.UpdateCache(cacheService, notifyService, updateStop, 10*time.Millisecond)
	}()

	time.Sleep(500 * time.Millisecond)

	_, new, err = cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 2, len(new))

	n := notifyClient.(*mocknotifier.Client)
	wasCalled := n.WasCalled()
	require.True(t, wasCalled)

	close(updateStop)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(30 * time.Second):
	}
}

func TestUpdateFlags_RemoveFlags(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	unit.SetLogger("update")

	readClient := mockreader.NewReader(
		reader.WithLocation("any"),
		mockreader.WithInitialFlags(
			map[string]*flags.Flag{
				"flag1": {
					Disabled: unit.Bool(false),
					Variants: map[string]any{
						"default": "default",
					},
				},
				"flag2": {
					Variants: map[string]any{
						"default": "default",
					},
				},
			},
		),
		mockreader.WithUpdatedFlags(
			map[string]*flags.Flag{
				"flag2": {
					Variants: map[string]any{
						"default": "default",
					},
				},
			},
		),
	)

	cacheService := cache.New(readClient)

	old, new, err := cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 0, len(old))
	require.Equal(t, 2, len(new))

	r := readClient.(*mockreader.Client)
	callCount := r.CallCount()
	require.Equal(t, 1, callCount)

	notifyClient := mocknotifier.NewNotifier()

	notifyService := notify.New(notifyClient)

	errCh := make(chan error, 1)
	updateStop := make(chan struct{})

	go func() {
		errCh <- server.UpdateCache(cacheService, notifyService, updateStop, 10*time.Millisecond)
	}()

	time.Sleep(500 * time.Millisecond)

	_, new, err = cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 1, len(new))

	n := notifyClient.(*mocknotifier.Client)
	wasCalled := n.WasCalled()
	require.True(t, wasCalled)

	close(updateStop)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(30 * time.Second):
	}
}
