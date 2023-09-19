package model

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vision-cli/vision/common/cases"
)

const (
	PersistenceDb   = "db"
	PersistenceNone = "none"
)

type Entity struct {
	Name        string  `yaml:"entity"`
	Persistence string  `yaml:"persistence"`
	Fields      []Field `yaml:"fields"`
}

func (e Entity) isValid() (bool, error) {
	if e.Name == "" {
		return false, fmt.Errorf("Entity name is required")
	}

	if e.Persistence != PersistenceDb && e.Persistence != PersistenceNone {
		return false, fmt.Errorf("Persisitence for %s is required and must be either '%s' or '%s'", e.Name, PersistenceDb, PersistenceNone)
	}

	if len(e.Fields) == 0 {
		return false, fmt.Errorf("Entity %s must have at least one field", e.Name)
	}

	for _, field := range e.Fields {
		if ok, err := field.isValid(); !ok {
			return false, err
		}
	}

	return true, nil
}

func (e Entity) GoAstTypeWithName(name string, target int) *ast.GenDecl {
	// Create the new field declarations
	fList := []*ast.Field{}
	for _, f := range e.Fields {
		field := &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(cases.Pascal(f.Name)),
			},
			Type: ast.NewIdent(f.GoType(target)),
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: f.Tag,
			},
		}
		fList = append(fList, field)
	}

	// Create the field list and add the new field to it
	fieldList := &ast.FieldList{
		List: fList,
	}

	// Create the type declaration for the Project struct
	structType := &ast.StructType{
		Fields: fieldList,
	}

	// Create the type spec for the Project struct
	typeSpec := &ast.TypeSpec{
		Name: ast.NewIdent(cases.Pascal(name)),
		Type: structType,
	}

	// Create the type declaration for the Entitie's type struct
	return &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{typeSpec},
	}
}

func (e Entity) GoAstType(target int) *ast.GenDecl {
	return e.GoAstTypeWithName(e.Name, target)
}

func (e Entity) GetFields() *ast.ExprStmt {
	elts := []ast.Expr{}
	for _, f := range e.Fields {
		switch f.Type {
		case TypeTimestamp:
			elts = append(elts, &ast.KeyValueExpr{
				Key:   ast.NewIdent(cases.Pascal(f.Name)),
				Value: ast.NewIdent("req." + cases.Pascal(f.Name) + ".AsTime()"),
			})
		case TypeEnum:
			elts = append(elts, &ast.KeyValueExpr{
				Key:   ast.NewIdent(cases.Pascal(f.Name)),
				Value: ast.NewIdent("req." + cases.Pascal(f.Name) + ".Enum().String()"),
			})
		case TypeId:
			elts = append(elts, &ast.KeyValueExpr{
				Key: ast.NewIdent(cases.Pascal(f.Name)),
				Value: &ast.CallExpr{
					Fun: ast.NewIdent("uuid.MustParse("),
					Args: []ast.Expr{
						ast.NewIdent("req." + cases.Pascal(f.Name) + ")"),
					},
				},
			})
		default:
			elts = append(elts, &ast.KeyValueExpr{
				Key:   ast.NewIdent(cases.Pascal(f.Name)),
				Value: ast.NewIdent("req." + cases.Pascal(f.Name)),
			})
		}
	}
	return &ast.ExprStmt{
		X: &ast.CompositeLit{
			Elts: elts,
		},
	}
}

func (e Entity) GoAstPbReturn(respType string) *ast.ReturnStmt {
	elts := *e.ConvertToReturnTypes()
	return &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("pb"),
						Sel: ast.NewIdent(cases.Pascal(respType)),
					},
					Elts: append([]ast.Expr{ast.NewIdent("Id: req.Id")}, elts...),
				},
			},
		},
	}
}

func (e Entity) GoAstListPbReturn(respType string) *ast.ExprStmt {
	elts := *e.ConvertToReturnTypes()
	return &ast.ExprStmt{
		X: &ast.UnaryExpr{
			Op: token.AND,
			X: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("pb"),
					Sel: ast.NewIdent(cases.Pascal(respType)),
				},
				Elts: elts,
			},
		},
	}
}

func (e Entity) ConvertToReturnTypes() *[]ast.Expr {
	elts := []ast.Expr{}
	for _, f := range e.Fields {
		switch f.Type {
		case TypeTimestamp:
			elts = append(elts, &ast.KeyValueExpr{
				Key: ast.NewIdent(cases.Pascal(f.Name)),
				Value: &ast.CallExpr{
					Fun:  ast.NewIdent("timestamppb.New"),
					Args: []ast.Expr{ast.NewIdent("m.Data." + cases.Pascal(f.Name))}},
			})
		case TypeId:
			elts = append(elts, &ast.KeyValueExpr{
				Key:   ast.NewIdent(cases.Pascal(f.Name)),
				Value: ast.NewIdent("m.Data." + cases.Pascal(f.Name) + ".String()"),
			})
		case TypeEnum:
			pascalEnumName := cases.Pascal(f.Name)
			elts = append(elts, &ast.KeyValueExpr{
				Key:   ast.NewIdent(cases.Pascal(f.Name)),
				Value: ast.NewIdent(fmt.Sprintf("pb.%s(pb.%s_value[m.Data.%s])", pascalEnumName, pascalEnumName, pascalEnumName)),
			})
		default:
			elts = append(elts, &ast.KeyValueExpr{
				Key:   ast.NewIdent(cases.Pascal(f.Name)),
				Value: ast.NewIdent("m.Data." + cases.Pascal(f.Name)),
			})
		}
	}
	return &elts
}

func (e Entity) GoAstSearch() *ast.ExprStmt {
	elts := []ast.Expr{}
	for _, e := range e.Fields {
		if e.IsSearchable && (e.Type == "string" || e.Type == "id") {
			elts = append(elts,
				ast.NewIdent(`"`+cases.Snake(e.Name)+`"`),
			)
		}
	}
	return &ast.ExprStmt{
		X: &ast.CompositeLit{
			Type: ast.NewIdent("[]string"),
			Elts: elts,
		},
	}
}

func (e Entity) GoAstFilter() []*ast.ExprStmt {
	elts := []ast.Expr{}
	for _, e := range e.Fields {
		if e.IsSearchable && e.Type == TypeId {
			elts = append(elts,
				ast.NewIdent(`
    if req.FilterBy != nil && req.FilterBy.`+cases.Pascal(e.Name)+` != "" {
        filterValue := req.FilterBy.`+cases.Pascal(e.Name)+`
        if &filterValue != nil {
            filters["`+cases.Snake(e.Name)+`"] = filterValue
        }
    }
                `),
			)
		}
		if e.IsSearchable && e.Type == TypeEnum {
			elts = append(elts,
				ast.NewIdent(`
	if req.FilterBy != nil && req.FilterBy.`+cases.Pascal(e.Name)+` != nil {
		filterValue := req.FilterBy.`+cases.Pascal(e.Name)+`.Get`+cases.Pascal(e.Name)+`()
		if &filterValue != nil {
			filters["`+cases.Snake(e.Name)+`"] = fmt.Sprintf("%v", filterValue)
		}
	}
				`),
			)
		}
		if e.IsSearchable && e.Type == TypeBoolean {
			elts = append(elts,
				ast.NewIdent(`
	if req.FilterBy != nil && req.FilterBy.`+cases.Pascal(e.Name)+` != nil {
		filterValue := req.FilterBy.`+cases.Pascal(e.Name)+`.GetValue()
		if &filterValue != nil {
			filters["`+cases.Snake(e.Name)+`"] = filterValue
		}
	}
				`),
			)
		}
		if e.IsSearchable && e.Type == TypeTimestamp {
			elts = append(elts,
				ast.NewIdent(`
	// TODO implement time stamp filter for`+cases.Pascal(e.Name)+`
				`),
			)
		}
	}
	exprStmt := []*ast.ExprStmt{}
	for _, elt := range elts {
		exprStmt = append(exprStmt, &ast.ExprStmt{X: elt})
	}
	return exprStmt
}
