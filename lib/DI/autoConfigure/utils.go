package autoconfigure

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	autoConfigModel "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/model"
)

func getModuleName(root string) string {
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_appContextLogger.Debug().Msg("running in path: " + p)
	mod, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(mod), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			_appContextLogger.Debug().Msg("parsed module: " + moduleName)
			return moduleName
		}
	}
	panic("could not parse module")
}

func Scan(comment string, root string) []*autoConfigModel.ScannedElement {
	res := []*autoConfigModel.ScannedElement{}

	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Doc == nil {
				continue
			}

			for _, c := range fn.Doc.List {
				if strings.Contains(c.Text, comment) {
					se := autoConfigModel.NewFuncElement(fn.Name.Name).
						Doc(c.Text).PackageName(file.Name.Name)

					splittedPath := strings.Split(path, "/")
					se.FileName(splittedPath[len(splittedPath)-1])
					splittedPath = splittedPath[:len(splittedPath)-1]

					se.Path(strings.Join(splittedPath, "/"))

					res = append(res, se)
				}
			}
		}
		return nil
	})

	return res
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

func BuildFile(se *autoConfigModel.ScannedElement) *autoConfigModel.GeneratedFile {
	gf := autoConfigModel.GetNewGeneratedFile("/" + se.GetPackageName())
	gf.FileName = se.GetName() + ".gen.go"
	return gf
}
