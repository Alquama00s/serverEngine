package autoconfigure

import (
	"sync"

	"github.com/Alquama00s/serverEngine/lib/logging/loggerFactory"
	"github.com/rs/zerolog"
)

type AppContext struct {
	moduleName string
	parsers    map[string]func() string
}

var (
	_appContextOnce   sync.Once
	_appContext       *AppContext
	_appContextLogger *zerolog.Logger
)

func GetAppContext() *AppContext {
	return _appContext
}

func InitAppContext(root string) *AppContext {
	_appContextOnce.Do(func() {
		_appContextLogger = loggerFactory.GetLogger("AppContext")
		m := getModuleName(root)
		_appContext = &AppContext{
			moduleName: m,
		}
	})
	return _appContext
}

func (c *AppContext) GetModuleName() string {
	return c.moduleName
}

func (c *AppContext) RegisterComponent(name string, component any) {
	_appContextLogger.Debug().Msg("registering " + name)
}
