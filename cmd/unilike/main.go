package main

import (
	"log"

	"github.com/michaeldebetaz/unilike/internal/env"
	"github.com/michaeldebetaz/unilike/internal/logger"
	"github.com/michaeldebetaz/unilike/internal/router"
)

func main() {
	logger.Init()
	env.Load()

	err := router.Start()
	if err != nil {
		log.Fatalf("Failed to start router: %v", err)
	}
}
