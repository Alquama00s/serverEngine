package autoconfigure

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"

	autoConfigModel "github.com/Alquama00s/serverEngine/lib/DI/autoConfigure/model"
)

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
