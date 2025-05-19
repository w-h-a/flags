package updateflags

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server"
	"github.com/w-h-a/flags/internal/server/clients/file"
	mockfile "github.com/w-h-a/flags/internal/server/clients/file/mock"
	mockmessage "github.com/w-h-a/flags/internal/server/clients/message/mock"
	"github.com/w-h-a/flags/internal/server/services/cache"
	"github.com/w-h-a/flags/internal/server/services/notify"
	"github.com/w-h-a/flags/tests/unit"
)

func TestUpdateFlags_NoChange(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Log("SKIPPING UNIT TEST")
		return
	}

	unit.SetLogger()

	fileClient := mockfile.NewFileClient(
		file.WithDir("any"),
		file.WithFiles("any"),
		mockfile.WithInitialFlags(
			map[string]*flags.Flag{
				"hello": {},
			},
		),
		mockfile.WithUpdatedFlags(
			map[string]*flags.Flag{
				"hello": {},
			},
		),
	)

	cacheService := cache.New(fileClient)

	old, new, err := cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 0, len(old))
	require.Equal(t, 1, len(new))

	f := fileClient.(*mockfile.Client)
	callCount := f.CallCount()
	require.Equal(t, 1, callCount)

	messageClient := mockmessage.NewMessageClient()

	notifyService := notify.New(messageClient)

	errCh := make(chan error, 1)
	updateStop := make(chan struct{})

	go func() {
		errCh <- server.UpdateCache(cacheService, notifyService, updateStop, 10*time.Millisecond)
	}()

	time.Sleep(500 * time.Millisecond)

	_, new, err = cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 1, len(new))

	m := messageClient.(*mockmessage.Client)
	wasCalled := m.WasCalled()
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

	unit.SetLogger()

	fileClient := mockfile.NewFileClient(
		file.WithDir("any"),
		file.WithFiles("any"),
		mockfile.WithInitialFlags(
			map[string]*flags.Flag{
				"flag1": {
					Disabled: unit.Bool(true),
				},
			},
		),
		mockfile.WithUpdatedFlags(
			map[string]*flags.Flag{
				"flag1": {
					Disabled: unit.Bool(false),
				},
				"flag2": {},
			},
		),
	)

	cacheService := cache.New(fileClient)

	old, new, err := cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 0, len(old))
	require.Equal(t, 1, len(new))

	f := fileClient.(*mockfile.Client)
	callCount := f.CallCount()
	require.Equal(t, 1, callCount)

	messageClient := mockmessage.NewMessageClient()

	notifyService := notify.New(messageClient)

	errCh := make(chan error, 1)
	updateStop := make(chan struct{})

	go func() {
		errCh <- server.UpdateCache(cacheService, notifyService, updateStop, 10*time.Millisecond)
	}()

	time.Sleep(500 * time.Millisecond)

	_, new, err = cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 2, len(new))

	m := messageClient.(*mockmessage.Client)
	wasCalled := m.WasCalled()
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

	unit.SetLogger()

	fileClient := mockfile.NewFileClient(
		file.WithDir("any"),
		file.WithFiles("any"),
		mockfile.WithInitialFlags(
			map[string]*flags.Flag{},
		),
		mockfile.WithUpdatedFlags(
			map[string]*flags.Flag{
				"flag1": {
					Disabled: unit.Bool(false),
				},
				"flag2": {},
			},
		),
	)

	cacheService := cache.New(fileClient)

	old, new, err := cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 0, len(old))
	require.Equal(t, 0, len(new))

	f := fileClient.(*mockfile.Client)
	callCount := f.CallCount()
	require.Equal(t, 1, callCount)

	messageClient := mockmessage.NewMessageClient()

	notifyService := notify.New(messageClient)

	errCh := make(chan error, 1)
	updateStop := make(chan struct{})

	go func() {
		errCh <- server.UpdateCache(cacheService, notifyService, updateStop, 10*time.Millisecond)
	}()

	time.Sleep(500 * time.Millisecond)

	_, new, err = cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 2, len(new))

	m := messageClient.(*mockmessage.Client)
	wasCalled := m.WasCalled()
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

	unit.SetLogger()

	fileClient := mockfile.NewFileClient(
		file.WithDir("any"),
		file.WithFiles("any"),
		mockfile.WithInitialFlags(
			map[string]*flags.Flag{
				"flag1": {
					Disabled: unit.Bool(false),
				},
				"flag2": {},
			},
		),
		mockfile.WithUpdatedFlags(
			map[string]*flags.Flag{
				"flag2": {},
			},
		),
	)

	cacheService := cache.New(fileClient)

	old, new, err := cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 0, len(old))
	require.Equal(t, 2, len(new))

	f := fileClient.(*mockfile.Client)
	callCount := f.CallCount()
	require.Equal(t, 1, callCount)

	messageClient := mockmessage.NewMessageClient()

	notifyService := notify.New(messageClient)

	errCh := make(chan error, 1)
	updateStop := make(chan struct{})

	go func() {
		errCh <- server.UpdateCache(cacheService, notifyService, updateStop, 10*time.Millisecond)
	}()

	time.Sleep(500 * time.Millisecond)

	_, new, err = cacheService.RetrieveFlags()
	require.NoError(t, err)
	require.Equal(t, 1, len(new))

	m := messageClient.(*mockmessage.Client)
	wasCalled := m.WasCalled()
	require.True(t, wasCalled)

	close(updateStop)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(30 * time.Second):
	}
}
