package main

import (
	"log"

	"go.uber.org/zap"
)

func initLogger() *zap.Logger {
	var logger *zap.Logger
	var err error
	if *develMode {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal("cannot init zap", err)
	}

	return logger
}
