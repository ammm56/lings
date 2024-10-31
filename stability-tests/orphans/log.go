package main

import (
	"github.com/ammm56/lings/infrastructure/logger"
	"github.com/ammm56/lings/util/panics"
)

var (
	backendLog = logger.NewBackend()
	log        = backendLog.Logger("ORPH")
	spawn      = panics.GoroutineWrapperFunc(log)
)