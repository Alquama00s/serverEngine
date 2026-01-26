package autoconfigure

import (
	"os"
	"sync"

	autoConfigModel "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/model"
	"github.com/Alquama00s/serverEngine/lib/logging/loggerFactory"
	"github.com/rs/zerolog"
)

type AppContextBuilder struct {
	rootPath string
	parsers  map[string]func(*autoConfigModel.ScannedElement, *AppContext) string
}

var (
	_appContextBuilderOnce   sync.Once
	_appContextBuilder       *AppContextBuilder
	_appContextBuilderLogger *zerolog.Logger
)

func GetAppContextBuilder() *AppContextBuilder {
	return _appContextBuilder
}

func InitAppContextBuilder(rootPath string) *AppContextBuilder {
	_appContextBuilderOnce.Do(func() {
		_appContextBuilderLogger = loggerFactory.GetLogger("AppContextBuilder")
		_appContextBuilder = &AppContextBuilder{
			rootPath: rootPath,
			parsers:  make(map[string]func(*autoConfigModel.ScannedElement, *AppContext) string),
		}
	})
	return _appContextBuilder
}

func (c *AppContextBuilder) RegisterParser(name string, parser func(*autoConfigModel.ScannedElement, *AppContext) string) {
	_appContextBuilderLogger.Debug().Msg("registering parser" + name)
	_, exist := c.parsers[name]
	if exist {
		panic("parser already exist " + name)
	}
	c.parsers[name] = parser
}

func (c *AppContextBuilder) RootPath(rootPath string) {
	c.rootPath = rootPath

}
func (c *AppContextBuilder) BootStrap() {
	err := os.RemoveAll(c.rootPath + "/generated")
	if err != nil {
		panic(err)
	}

	_appContextBuilderLogger.Debug().Msg("registered root path" + c.rootPath)
	ctx := InitAppContext(c.rootPath)
	for k, v := range c.parsers {
		scannedElements := Scan(k, c.rootPath)
		for _, se := range scannedElements {
			file := BuildFile(se)
			file.Contents = v(se, ctx)
			err = os.MkdirAll(c.rootPath+"/generated/"+file.GetPath(), 0755)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile(c.rootPath+"/generated"+file.GetPath()+"/"+file.FileName, []byte(file.Contents), 0644)
			if err != nil {
				panic(err)
			}
		}
	}
}
