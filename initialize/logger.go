package initialize

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger() {
	// logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create zap logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
}
