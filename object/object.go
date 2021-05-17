// Credit to Thorsten Ball for:
//  - Idea of interfaces for AST and objects.
//      - Actual AST structures/object enumeration by me.
//  - Idea of object system to pass evaluated AST results.
//  - Variable argument error messages.
//      - Handling of all syntax/runtime errors decided by me.
//  - Applying expressions to function calls.
//      - Function scoping/reuse of identifier errors done by me.
//  - Built-in functions interface.
//  - Using a map to hold the contents of a state.

package object

import (
	"ariel/ast"
	"bytes"
	"fmt"
)

type ObjectType int

type Object interface {
	Type() ObjectType
	Eval() string
}

const (
	ErrorObj ObjectType = iota
	CharObj
	IntObj
	FloatObj
	StringObj
	BoolObj
	ArrObj
	ReturnObj
	FuncDeclObj
	BuiltInObj
)

func ObjString(obj Object) string {
	switch obj.(type) {
	case Char:
		return "char"
	case Int:
		return "int"
	case Float:
		return "float"
	case String:
		return "string"
	case Bool:
		return "bool"
	case Array:
		return "array"
	case Return:
		return "return"
	case FuncDecl:
		return "funcdecl"
	case BuiltIn:
		return "builtin"
	default:
		return ""
	}
}

type Error struct {
	Message string
}

func (e Error) Type() ObjectType { return ErrorObj }
func (e Error) Eval() string     { return e.Message }

type Char struct {
	Value string
}

func (c Char) Type() ObjectType { return CharObj }
func (c Char) Eval() string     { return c.Value }

type Int struct {
	Value int64
}

func (i Int) Type() ObjectType { return IntObj }
func (i Int) Eval() string     { return fmt.Sprintf("%d", i.Value) }

type Float struct {
	Value float64
}

func (f Float) Type() ObjectType { return FloatObj }
func (f Float) Eval() string     { return fmt.Sprintf("%f", f.Value) }

type String struct {
	Value string
}

func (s String) Type() ObjectType { return StringObj }
func (s String) Eval() string     { return s.Value }

type Bool struct {
	Value bool
}

func (b Bool) Type() ObjectType { return BoolObj }
func (b Bool) Eval() string     { return fmt.Sprintf("%t", b.Value) }

type Array struct {
	ElementType string
	Elements    []Object
}

func (a Array) Type() ObjectType { return ArrObj }
func (a Array) Eval() string {
	var out bytes.Buffer
	if len(a.Elements) > 0 {
		out.WriteString("{ ")
		for i, element := range a.Elements {
			out.WriteString(element.Eval())
			if i+1 != len(a.Elements) {
				out.WriteString(", ")
			}
		}
		out.WriteString(" }")
	} else {
		out.WriteString("{}")
	}
	return out.String()
}

type Return struct {
	Value Object
}

func (r Return) Type() ObjectType { return ReturnObj }
func (r Return) Eval() string     { return r.Value.Eval() }

type FuncDecl struct {
	ReturnType ast.Type
	Ident      ast.Identifier
	Parameters []ast.Param
	Body       ast.Block
	State      *State
}

func (fd FuncDecl) Type() ObjectType { return FuncDeclObj }
func (fd FuncDecl) Eval() string     { return "function" }

type BuiltInFunc func(args ...Object) Object

type BuiltIn struct {
	Function BuiltInFunc
}

func (b BuiltIn) Type() ObjectType { return BuiltInObj }
func (b BuiltIn) Eval() string     { return "builtin" }
