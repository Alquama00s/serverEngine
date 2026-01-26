package autoConfigParsers

import (
	"strings"

	autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
	autoConfigModel "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/model"
)

func ParseService(se *autoConfigModel.ScannedElement, ac *autoconfigure.AppContext) string {
	res := strings.Builder{}
	moduleName := ac.GetModuleName()
	res.WriteString(`
package generatedCmd

import autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
`)
	res.WriteString(se.GetImportLine(moduleName) + "\n")

	res.WriteString(`
func RegisterService() {
	ctx := autoconfigure.GetAppContext()
	`)
	temp := `ctx.Register("@service"+`
	temp += "\"" + se.GetName() + "\""
	temp += ","
	temp += se.GetPackageName() + "." + se.GetName()
	temp += "())"
	temp += "\n"
	res.WriteString(temp)
	res.WriteString(`
}
	`)

	return res.String()
}
