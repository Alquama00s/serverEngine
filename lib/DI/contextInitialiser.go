package DI

import (
	autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
	autoConfigParsers "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/parsers"
)

func init() {}

func InitialiseContextBuilder(rootPath string) {
	ctxBuilder := autoconfigure.InitAppContextBuilder(rootPath)
	// ctxBuilder.RegisterParser("@service", autoConfigParsers.ParseService)
	// ctxBuilder.ImportLine("generatedCmd \"serverEngineTests/generated/service\"")
	// ctxBuilder.InitLine("generatedCmd.RegisterService()")

	ctxBuilder.RegisterParser("@Controller", autoConfigParsers.ParseController)
	ctxBuilder.ImportLine("generatedController \"serverEngineTests/generated/controller\"")
	ctxBuilder.InitLine("generatedController.RegisterControllers()")

	ctxBuilder.GenerateCode()
}
