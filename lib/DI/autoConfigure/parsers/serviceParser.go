package autoConfigParsers

import (
	"strings"

	autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
	autoConfigModel "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/model"
)

func ParseService(se []*autoConfigModel.ScannedElement, ac *autoconfigure.AppContext) []*autoConfigModel.GeneratedFile {
	res := strings.Builder{}
	importsMap := make(map[string]struct{})
	moduleName := ac.GetModuleName()
	res.WriteString(`
package generatedCmd

import autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
`)
	for _, s := range se {
		_, exists := importsMap[s.GetImportLine(moduleName)]
		if !exists {
			res.WriteString(s.GetImportLine(moduleName) + "\n")
			importsMap[s.GetImportLine(moduleName)] = struct{}{}
		}
	}

	res.WriteString(`
func RegisterService() {
	`)
	if len(se) > 0 {
		res.WriteString("ctx := autoconfigure.GetAppContext()\n")
	}
	for _, s := range se {
		temp := `ctx.Register("@service"+`
		temp += "\"" + s.GetName() + "\""
		temp += ","
		temp += s.GetPackageName() + "." + s.GetName()
		temp += "())"
		temp += "\n"
		res.WriteString(temp)
	}
	res.WriteString(`
}
	`)

	gf := autoConfigModel.GetNewGeneratedFile("/service")
	gf.FileName = "allServices.gen.go"
	gf.Contents = res.String()

	return []*autoConfigModel.GeneratedFile{gf}
}
