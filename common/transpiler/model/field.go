package model

import (
	"fmt"

	"github.com/oshothebig/pbast"

	"github.com/vision-cli/vision/common/cases"
)

const (
	GoTargetModel    int = 0
	GoTargetResolver int = 1
)

const (
	TypeString    = "string"
	TypeTimestamp = "timestamp"
	TypeUnsigned  = "unsigned"
	TypeIntger    = "integer"
	TypeFloat     = "float"
	TypeId        = "id"
	TypeBoolean   = "boolean"
	TypeEnum      = "enum"
)

const (
	TypeGqlString  = "String"
	TypeGqlInt     = "Int"
	TypeGqlBoolean = "Boolean"
	TypeGqlID      = "ID"
)

const (
	TypeGoTime      = "time.Time"
	TypeGoString    = "string"
	TypeGoInt       = "int32"
	TypeGoUnsigned  = "uint32"
	TypeGoBoolean   = "bool"
	TypeGoID        = "uuid.UUID"
	TypeGoGraphqlID = "graphql.ID"
)

type Field struct {
	Name         string `yaml:"name"`
	Type         string `yaml:"type"`
	Tag          string `yaml:"tag"`
	IsArray      bool   `yaml:"array"`
	IsNullable   bool   `yaml:"required"`
	IsSearchable bool   `yaml:"searchable"`
}

func (f Field) isValid() (bool, error) {
	if f.Name == "" {
		return false, fmt.Errorf("Field name is required")
	}

	if f.Type != TypeString &&
		f.Type != TypeUnsigned &&
		f.Type != TypeIntger &&
		f.Type != TypeFloat &&
		f.Type != TypeBoolean &&
		f.Type != TypeTimestamp &&
		f.Type != TypeId &&
		f.Type != TypeEnum {
		return false, fmt.Errorf("Field type must be one of '%s', '%s', '%s', '%s', '%s', '%s', '%s' or '%s'", TypeString, TypeUnsigned, TypeIntger, TypeFloat, TypeBoolean, TypeTimestamp, TypeId, TypeEnum)
	}

	return true, nil
}

func (f Field) GoType(target int) string {
	var t string
	switch f.Type {
	case TypeTimestamp:
		if target == GoTargetResolver {
			t = TypeGoString // graphql doesn't support timestamp so we're using string
			break
		}
		if target == GoTargetModel {
			t = TypeGoTime
			break
		}
	case TypeUnsigned:
		if target == GoTargetResolver {
			t = TypeGoInt // graphql doesn't support uint so we're converting to int
			break
		}
		if target == GoTargetModel {
			t = TypeGoUnsigned
			break
		}
	case TypeIntger:
		t = TypeGoInt
	case TypeId:
		if target == GoTargetResolver {
			t = TypeGoGraphqlID
			break
		}
		if target == GoTargetModel {
			t = TypeGoID
			break
		}
	case TypeBoolean:
		t = TypeGoBoolean
	case TypeString:
		t = TypeGoString
	case TypeEnum:
		t = TypeGoString
	default:
		t = f.Type
	}
	if target == GoTargetResolver && f.IsNullable {
		t = "*" + t
	}
	if f.IsArray {
		return "[]" + t
	}
	return t
}

func (f Field) FilterGoType(target int) string {
	var t string
	switch f.Type {
	case TypeTimestamp:
		if target == GoTargetResolver {
			t = TypeGoString // graphql doesn't support timestamp so we're using string
			break
		}
		if target == GoTargetModel {
			t = TypeGoTime
			break
		}
	case TypeUnsigned:
		if target == GoTargetResolver {
			t = TypeGoInt // graphql doesn't support uint so we're converting to int
			break
		}
		if target == GoTargetModel {
			t = TypeGoUnsigned
			break
		}
	case TypeIntger:
		t = TypeGoInt
	case TypeId:
		if target == GoTargetResolver {
			t = TypeGoGraphqlID
			break
		}
		if target == GoTargetModel {
			t = TypeGoID
			break
		}
	case TypeBoolean:
		t = TypeGoBoolean
	case TypeString:
		t = TypeGoString
	case TypeEnum:
		t = cases.Pascal(f.Name)
	default:
		t = f.Type
	}
	t = "*" + t
	if f.IsArray {
		return "[]" + t
	}
	return t
}

func (f Field) PbType() pbast.BuiltinType {
	var t pbast.BuiltinType
	switch f.Type {
	case TypeString:
		t = pbast.String
	case TypeTimestamp:
		t = pbast.BuiltinType(pbast.Timestamp)
	case TypeUnsigned:
		t = pbast.UInt32
	case TypeIntger:
		t = pbast.Int32
	case TypeId:
		t = pbast.String
	case TypeBoolean:
		t = pbast.Bool
	default:
		return pbast.BuiltinType(f.Type)
	}
	return t
}

func (f Field) GqlType() string {
	var t string
	switch f.Type {
	case TypeString:
		t = TypeGqlString
	case TypeTimestamp:
		t = TypeGqlString
	case TypeUnsigned:
		t = TypeGqlInt
	case TypeIntger:
		t = TypeGqlInt
	case TypeId:
		t = TypeGqlID
	case TypeBoolean:
		t = TypeGqlBoolean
	case TypeEnum:
		t = TypeGqlString
	default:
		panic("Dont know how to convert " + f.Type + " to pbast.BuiltinType")
	}
	if f.IsArray {
		t = "[]" + t
	}
	if !f.IsNullable {
		t = t + "!"
	}
	return t
}

func (f Field) FilterGqlType() string {
	var t string
	switch f.Type {
	case TypeString:
		t = TypeGqlString
	case TypeTimestamp:
		t = TypeGqlString
	case TypeUnsigned:
		t = TypeGqlInt
	case TypeIntger:
		t = TypeGqlInt
	case TypeId:
		t = TypeGqlID
	case TypeBoolean:
		t = TypeGqlBoolean
	case TypeEnum:
		t = cases.Pascal(f.Name)
	default:
		panic("Dont know how to convert " + f.Type + " to pbast.BuiltinType")
	}
	if f.IsArray {
		t = "[]" + t
	}
	return t
}
