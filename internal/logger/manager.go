package logger

import (
	"sync"

	"github.com/hashicorp/go-hclog"
)

type InterceptManager struct {
	mu             sync.RWMutex // Intercepts can register and deregister sinks, so we need to lock access to them for thread safety
	SyncIntercept  *hclog.Logger
	AsyncIntercept *hclog.Logger
}

func NewInterceptManager() *InterceptManager {
	return &InterceptManager{
		mu: sync.RWMutex{},
	}
}

func (im *InterceptManager) WithSyncIntercept(logger *hclog.Logger) *InterceptManager {
	im.mu.Lock()
	defer im.mu.Unlock()
	return &InterceptManager{
		mu:             sync.RWMutex{},
		SyncIntercept:  logger,
		AsyncIntercept: im.AsyncIntercept,
	}
}

func (im *InterceptManager) WithAsyncIntercept(logger *hclog.Logger) *InterceptManager {
	im.mu.Lock()
	defer im.mu.Unlock()
	return &InterceptManager{
		mu:             sync.RWMutex{},
		SyncIntercept:  im.SyncIntercept,
		AsyncIntercept: logger,
	}
}
