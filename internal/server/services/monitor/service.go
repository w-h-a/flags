package monitor

import "github.com/w-h-a/flags/internal/server/config"

type Service struct {
}

func (m *Service) Info() Info {
	return Info{
		Env:     config.Env(),
		Name:    config.Name(),
		Version: config.Version(),
		Status:  "up",
	}
}

func New() *Service {
	return &Service{}
}
