package main

import (
	"github.com/michaeldebetaz/unilscrap/internal/env"
	"github.com/michaeldebetaz/unilscrap/internal/logger"
	"github.com/michaeldebetaz/unilscrap/internal/scraper"
)

func main() {
	logger.Init()
	env.Load()

	scraper.Scrape()
}
