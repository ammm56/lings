package rpcclient

import (
	"github.com/ammm56/lings/infrastructure/logger"
	"github.com/ammm56/lings/util/panics"
)

var log = logger.RegisterSubSystem("RPCC")
var spawn = panics.GoroutineWrapperFunc(log)
