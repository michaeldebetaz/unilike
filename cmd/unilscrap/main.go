package main

import (
	"github.com/michaeldebetaz/unilscrap/internal/env"
	"github.com/michaeldebetaz/unilscrap/internal/logger"
	"github.com/michaeldebetaz/unilscrap/internal/scrapper"
)

func main() {
	logger.Init()
	env.Load()

	scrapper.Scrape()
}
