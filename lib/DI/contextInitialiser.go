package DI

import (
	autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
	autoConfigParsers "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/parsers"
)

func InitialiseContextBuilder(rootPath string) {
	ctxBuilder := autoconfigure.InitAppContextBuilder(rootPath)
	ctxBuilder.RegisterParser("@service", autoConfigParsers.ParseService)
}
