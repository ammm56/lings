package rpchandlers

import (
	"github.com/ammm56/lings/infrastructure/logger"
	"github.com/ammm56/lings/util/panics"
)

var log = logger.RegisterSubSystem("RPCS")
var spawn = panics.GoroutineWrapperFunc(log)
