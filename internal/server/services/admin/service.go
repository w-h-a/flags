package admin

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/w-h-a/flags/internal/flags"
	"github.com/w-h-a/flags/internal/server/clients/reader"
	"github.com/w-h-a/flags/internal/server/clients/writer"
	"github.com/w-h-a/flags/internal/server/config"
	"gopkg.in/yaml.v3"
)

type Service struct {
	writeClient writer.Writer
	readClient  reader.Reader
}

func (s *Service) RetrieveFlag(ctx context.Context, key string) (map[string]*flags.Flag, error) {
	bs, err := s.readClient.ReadByKey(ctx, key)
	if err != nil && errors.Is(err, reader.ErrRecordNotFound) {
		return nil, flags.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return flags.Factory(bs, config.FlagFormat())
}

func (s *Service) RetrieveFlags(ctx context.Context) (map[string]*flags.Flag, error) {
	bs, err := s.readClient.Read(ctx)
	if err != nil {
		return nil, err
	}

	return flags.Factory(bs, config.FlagFormat())
}

func (s *Service) UpsertFlag(ctx context.Context, key string, flag map[string]*flags.Flag) (map[string]*flags.Flag, error) {
	var bs []byte
	var err error

	switch strings.ToLower(config.FlagFormat()) {
	case "json":
		bs, err = json.Marshal(flag)
	default:
		bs, err = yaml.Marshal(flag)
	}

	if err != nil {
		return nil, err
	}

	if err := s.writeClient.Write(ctx, key, bs); err != nil {
		return nil, err
	}

	return flag, nil
}

func New(writeClient writer.Writer, readClient reader.Reader) *Service {
	return &Service{
		writeClient: writeClient,
		readClient:  readClient,
	}
}
