package autoconfigure

import (
	"os"
	"strings"

	generate "github.com/Alquama00s/serverEngine/lib/goGen"
)

func WriteService(root string) {
	err := os.RemoveAll(root + "/generated")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(root+"/generated", 0755)
	if err != nil {
		panic(err)
	}
	generate.Write(root+"/generated/genService.gen.go", BuildString(root))
}

func BuildString(root string) string {

	serviceList := Scan("@service", root)
	res := strings.Builder{}
	moduleName := GetAppContext().GetModuleName()
	res.WriteString(`
package generatedCmd

import autoconfigure "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure"
`)
	for _, s := range serviceList {
		res.WriteString(s.GetImportLine(moduleName) + "\n")
	}

	res.WriteString(`
func RegisterService() {
	ctx := autoconfigure.GetAppContext()
	`)
	for _, s := range serviceList {
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

	return res.String()
}
