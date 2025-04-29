package monitor

import "time"

type Info struct {
	Env          string    `json:"env"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	Status       string    `json:"up"`
	CacheRefresh time.Time `json:"cacheRefresh" example:"2022-06-13T11:22:55.941628+02:00"`
}
