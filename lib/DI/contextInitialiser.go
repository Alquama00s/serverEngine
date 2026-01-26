package DI

import (
	autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
	autoConfigParsers "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/parsers"
)

func initialiseContext(rootPath string) {
	ctxBuilder := autoconfigure.GetAppContextBuilder(rootPath)
	ctxBuilder.RegisterParser("@service", autoConfigParsers.ParseService)
}
