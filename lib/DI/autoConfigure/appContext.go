package autoconfigure

import (
	"os"
	"strings"
	"sync"

	"github.com/Alquama00s/serverEngine/lib/logging/loggerFactory"
	"github.com/rs/zerolog"
)

type AppContext struct {
	moduleName string
}

var (
	_appContextOnce sync.Once
	_appContext     *AppContext
	_logger         *zerolog.Logger
)

func GetAppContext() *AppContext {
	_appContextOnce.Do(func() {
		_logger = loggerFactory.GetLogger()
		m := getModuleName()
		_appContext = &AppContext{
			moduleName: m,
		}
	})
	return _appContext
}

func getModuleName() string {
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_logger.Debug().Msg("running in path: " + p)
	mod, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(mod), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			_logger.Debug().Msg("parsed module: " + moduleName)
			return moduleName
		}
	}
	panic("could not parse module")
}

func (c *AppContext) GetModuleName() string {
	return c.moduleName
}

func (c *AppContext) Register(name string, component any) {
	_logger.Debug().Msg("registering " + name)
}
