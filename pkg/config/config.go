package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	PlaylistsRoot     string `env:"PLAYLISTS_ROOT, default=./"`
	DestinationFolder string `env:"DESTINATION_FOLDER, default=./"`

	// add more fields here
}

func NewConfig(ctx context.Context) (*Config, error) {
	cfg := Config{}
	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
