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
	err := os.Chdir(root)
	if err != nil {
		panic(err)
	}
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
			// Handle function declarations
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Doc != nil {
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
				continue
			}

			// Handle struct declarations
			gd, ok := decl.(*ast.GenDecl)
			if ok && gd.Tok == token.TYPE && gd.Doc != nil {
				for _, c := range gd.Doc.List {
					if strings.Contains(c.Text, comment) {
						for _, spec := range gd.Specs {
							ts, ok := spec.(*ast.TypeSpec)
							if !ok {
								continue
							}
							// Only process struct types
							if _, ok := ts.Type.(*ast.StructType); !ok {
								continue
							}

							se := autoConfigModel.NewStructElement(ts.Name.Name).
								Doc(c.Text).PackageName(file.Name.Name)

							splittedPath := strings.Split(path, "/")
							se.FileName(splittedPath[len(splittedPath)-1])
							splittedPath = splittedPath[:len(splittedPath)-1]

							se.Path(strings.Join(splittedPath, "/"))

							res = append(res, se)
						}
					}
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

func BuildInitFile(importLines []string, funCalls []string) *autoConfigModel.GeneratedFile {
	res := strings.Builder{}
	res.WriteString(`
package generated

import (
	`)
	for _, il := range importLines {
		res.WriteString(il + "\n")
	}

	res.WriteString(`
	"github.com/Alquama00s/serverEngine"
	"github.com/Alquama00s/serverEngine/lib/routing/registrar"
)

func init() {
	`)

	for _, fc := range funCalls {
		res.WriteString(fc + "\n")
	}

	res.WriteString(`
}

func InitialiseServer() *registrar.DefaultRegistrar {
	return serverEngine.Registrar()
}
	`)

	gc := autoConfigModel.GetNewGeneratedFile("")
	gc.FileName = "init.gen.go"
	gc.Contents = res.String()
	return gc
}

func WriteFile(c *AppContextBuilder, file *autoConfigModel.GeneratedFile) {

	err := os.MkdirAll(c.rootPath+"/generated/"+file.GetPath(), 0755)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(c.rootPath+"/generated"+file.GetPath()+"/"+file.FileName, []byte(file.Contents), 0644)
	if err != nil {
		panic(err)
	}

}
