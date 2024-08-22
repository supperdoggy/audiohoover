package main

import (
	"context"

	"github.com/supperdoggy/audiohoover/pkg/config"
	"github.com/supperdoggy/audiohoover/pkg/service"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewDevelopment()
	ctx := context.Background()

	cfg, err := config.NewConfig(ctx)
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	count, err := service.RunApp(log, cfg)
	if err != nil {
		log.Fatal("failed to run app", zap.Error(err))
	}

	log.Info("done", zap.Int("count", count))

}
