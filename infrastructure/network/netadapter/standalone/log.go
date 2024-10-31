package standalone

import (
	"github.com/ammm56/lings/infrastructure/logger"
	"github.com/ammm56/lings/util/panics"
)

var log = logger.RegisterSubSystem("NTAR")
var spawn = panics.GoroutineWrapperFunc(log)
