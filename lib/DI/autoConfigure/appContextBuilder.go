package autoconfigure

import (
	"os"
	"sync"

	autoConfigModel "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/model"
	"github.com/Alquama00s/serverEngine/lib/logging/loggerFactory"
	"github.com/rs/zerolog"
)

type AppContextBuilder struct {
	rootPath        string
	initImportLines []string
	initLines       []string
	parsers         map[string]func([]*autoConfigModel.ScannedElement, *AppContext) []*autoConfigModel.GeneratedFile
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
		err := os.Chdir(rootPath)
		if err != nil {
			panic(err)
		}
		_appContextBuilderLogger = loggerFactory.GetLogger("AppContextBuilder")
		_appContextBuilder = &AppContextBuilder{
			rootPath: ".",
			parsers:  make(map[string]func([]*autoConfigModel.ScannedElement, *AppContext) []*autoConfigModel.GeneratedFile),
		}
	})
	return _appContextBuilder
}

func (c *AppContextBuilder) RegisterParser(name string, parser func([]*autoConfigModel.ScannedElement, *AppContext) []*autoConfigModel.GeneratedFile) {
	_appContextBuilderLogger.Debug().Msg("registering parser" + name)
	_, exist := c.parsers[name]
	if exist {
		panic("parser already exist " + name)
	}
	c.parsers[name] = parser
}

func (c *AppContextBuilder) ImportLine(il string) *AppContextBuilder {
	if c.initImportLines == nil {
		c.initImportLines = []string{}
	}
	c.initImportLines = append(c.initImportLines, il)
	return c
}

func (c *AppContextBuilder) InitLine(il string) *AppContextBuilder {
	if c.initLines == nil {
		c.initLines = []string{}
	}
	c.initLines = append(c.initLines, il)
	return c
}

func (c *AppContextBuilder) RootPath(rootPath string) {
	c.rootPath = rootPath

}

func (c *AppContextBuilder) GenerateCode() {
	shoudGenerateMainFile := false
	err := os.RemoveAll(c.rootPath + "/generated")
	if err != nil {
		panic(err)
	}

	_appContextBuilderLogger.Debug().Msg("registered root path" + c.rootPath)
	ctx := InitAppContext(c.rootPath)
	for k, v := range c.parsers {
		scannedElements := Scan(k, c.rootPath)
		if len(scannedElements) > 0 {
			shoudGenerateMainFile = true
			files := v(scannedElements, ctx)
			for _, f := range files {
				WriteFile(c, f)
			}
		}
	}
	if shoudGenerateMainFile {
		file := BuildInitFile(c.initImportLines, c.initLines)
		WriteFile(c, file)
	}
}
