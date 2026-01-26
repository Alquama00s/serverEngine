package autoConfigParsers

import (
	"strings"

	autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
	autoConfigModel "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/model"
)

func ParseInitializer(se []*autoConfigModel.ScannedElement, ac *autoconfigure.AppContext) []*autoConfigModel.GeneratedFile {
	res := strings.Builder{}
	importsMap := make(map[string]struct{})
	moduleName := ac.GetModuleName()
	res.WriteString(`
package generatedInitializer

`)

	if len(se) > 0 {
		res.WriteString("import \"github.com/Alquama00s/serverEngine\"\n")
	}
	for _, s := range se {
		_, exists := importsMap[s.GetImportLine(moduleName)]
		if !exists {
			res.WriteString(s.GetImportLine(moduleName) + "\n")
			importsMap[s.GetImportLine(moduleName)] = struct{}{}
		}
	}

	res.WriteString(`
func RegisterInitializers() {
	`)

	if len(se) > 0 {
		res.WriteString("reg := serverEngine.Registrar()\n")
	}
	for _, s := range se {
		temp := "reg.RegisterInitializers(&"
		temp += s.GetPackageName() + "." + s.GetName()
		temp += "{})"
		temp += "\n"
		res.WriteString(temp)
	}
	res.WriteString(`
}
	`)

	gf := autoConfigModel.GetNewGeneratedFile("/initializer")
	gf.FileName = "allInitializers.gen.go"
	gf.Contents = res.String()

	return []*autoConfigModel.GeneratedFile{gf}
}
