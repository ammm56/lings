package prefixmanager

import (
	"github.com/ammm56/lings/infrastructure/logger"
	"github.com/ammm56/lings/util/panics"
)

var log = logger.RegisterSubSystem("PRFX")
var spawn = panics.GoroutineWrapperFunc(log)
