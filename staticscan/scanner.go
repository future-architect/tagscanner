package staticscan

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Field struct {
	Name     string
	Type     string
	Tag      string
	Comment  string
	FullPath string
}

type Struct struct {
	PackageName string
	StructName  string
	Fields      []Field
	Comment     string
	FullPath    string
}

func parseType(e ast.Expr) (string, error) {
	switch ft := e.(type) {
	case *ast.Ident:
		return ft.Name, nil
	case *ast.SelectorExpr:
		fldPkg, err := parseType(ft.X)
		if err != nil {
			return "", err
		}
		return fldPkg + "." + ft.Sel.Name, nil
	case *ast.StarExpr:
		elmType, err := parseType(ft.X)
		if err != nil {
			return "", err
		}
		return "*" + elmType, nil
	case *ast.ArrayType:
		elmType, err := parseType(ft.Elt)
		if err != nil {
			return "", err
		}
		return "[]" + elmType, nil
	case *ast.MapType:
		keyType, err := parseType(ft.Key)
		if err != nil {
			return "", err
		}
		valueType, err := parseType(ft.Value)
		if err != nil {
			return "", err
		}
		return "map[" + keyType + "]" + valueType, nil
	default:
		t := reflect.TypeOf(ft)
		return "", fmt.Errorf("Unsupported field type: *%v\n", t.Elem().String())
	}

}

func Scan(rootPath, tagName string) ([]Struct, error) {
	var results []Struct
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			cleanPath := filepath.Clean(path)
			fSet := token.NewFileSet()
			astFile, err := parser.ParseFile(fSet, cleanPath, nil, parser.ParseComments)
			if err != nil {
				return fmt.Errorf("parser.ParseFile() failed: %w", err)
			}
			ast.Inspect(astFile, func(n ast.Node) bool {
				nt, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}
				expr := nt.Type
				name := nt.Name.Name
				if expr == nil {
					return true
				}

				st, ok := expr.(*ast.StructType)
				if !ok {
					return true
				}
				var comments []string
				cmap := ast.NewCommentMap(fSet, st, astFile.Comments)
				for _, comment := range cmap[st] {
					comments = append(comments, comment.Text())
				}
				result := Struct{
					PackageName: astFile.Name.Name,
					StructName:  name,
					Comment:     strings.Join(comments, ""),
					FullPath:    cleanPath,
				}
				for _, f := range st.Fields.List {
					if f.Tag == nil {
						continue
					}
					tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
					tv, ok := tag.Lookup(tagName)
					if !ok {
						continue
					}
					typeName, err := parseType(f.Type)
					if err != nil {
						panic(err)
					}
					for _, name := range f.Names {
						result.Fields = append(result.Fields, Field{
							Name:     name.Name,
							Type:     typeName,
							Tag:      tv,
							FullPath: cleanPath,
						})
					}
				}
				if len(result.Fields) > 0 {
					results = append(results, result)
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}
