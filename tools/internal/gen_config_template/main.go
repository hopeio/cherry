package main

import (
	"fmt"
	"github.com/hopeio/cherry/utils/io/fs"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {

	packageName := "D:/code/hopeio/cherry/initialize"
	fs.RangeDir(packageName, func(dir string, entries []os.DirEntry) ([]os.DirEntry, error) {
		var recursion []os.DirEntry
		for _, entry := range entries {
			if entry.IsDir() {
				recursion = append(recursion, entry)
				getConfigStructs(dir + fs.PathSeparator + entry.Name())
			}

		}
		return recursion, nil
	})

}

// go递归读取某个包下的所有struct
func getStructs(packageName string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packageName, nil, parser.ParseComments)
	if err != nil {
		return
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					_, ok = typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					fmt.Printf("Struct %s found in file %s\n", typeSpec.Name.Name, fset.File(file.Pos()).Name())
				}
			}
		}
	}
}

func getConfigStructs(packageName string) {

	var targetMethod = "Init"
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packageName, nil, parser.ParseComments)
	if err != nil {
		return
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if x, ok := n.(*ast.FuncDecl); ok {
					if x.Recv != nil && len(x.Recv.List) > 0 {
						recvType := x.Recv.List[0].Type
						if star, ok := recvType.(*ast.StarExpr); ok {
							if ident, ok := star.X.(*ast.Ident); ok {
								if x.Name.Name == targetMethod {
									fmt.Printf("Struct %s implements method %s\n", ident.Name, targetMethod)
								}
							}

						}
					}
				}
				return true
			})
		}
	}
}
