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

package ast

type Node interface{}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

type Program struct {
	Statements []Statement
}

type Type struct {
	Value string
}

type FuncDecl struct {
	Type       Type
	Ident      Identifier
	Parameters []Param
	Body       Block
}

type VarDecl struct {
	Type        Type
	Ident       Identifier
	Value       Expression
	Initialized bool
}

type Param struct {
	Type  Type
	Ident Identifier
	Array bool
}

type Block struct {
	Statements []Statement
}

type While struct {
	Condition Expression
	Body      Statement
}

type For struct {
	Init      Expression
	Condition Expression
	Increment Expression
	Body      Statement
	VarDecl   bool
	Type      Type
	Ident     Identifier
	Value     Expression
}

type IfElse struct {
	Condition      Expression
	Consequence    Statement
	Alternative    Statement
	HasAlternative bool
}

type Return struct {
	Value Expression
	Void  bool
}

type ExprStmt struct {
	Expression Expression
}

type PrefixExpr struct {
	Op    string
	Right Expression
}

type InfixExpr struct {
	Left  Expression
	Op    string
	Right Expression
}

type Assign struct {
	Ident Identifier
	Value Expression
}

type AssignExpr struct {
	Ident Identifier
	Op    string
	Value Expression
}

type Call struct {
	Function  Identifier
	Arguments []Expression
	Void      bool
}

type Identifier struct {
	Name string
}

type CharCon struct {
	Value string
}

type IntCon struct {
	Value int64
}

type FloatCon struct {
	Value float64
}

type StringCon struct {
	Value string
}

type Bool struct {
	Value bool
}

type Array struct {
	Elements []Expression
}

type IndexExpr struct {
	Ident Identifier
	Index Expression
}

type AssignIndexExpr struct {
	Ident Identifier
	Index Expression
	Value Expression
}

type AssignExprIndexExpr struct {
	Ident Identifier
	Index Expression
	Op    string
	Value Expression
}

func (fd FuncDecl) statement()               {}
func (vd VarDecl) statement()                {}
func (bs Block) statement()                  {}
func (w While) statement()                   {}
func (f For) statement()                     {}
func (is IfElse) statement()                 {}
func (r Return) statement()                  {}
func (es ExprStmt) statement()               {}
func (pe PrefixExpr) expression()            {}
func (ie InfixExpr) expression()             {}
func (a Assign) expression()                 {}
func (ae AssignExpr) expression()            {}
func (ce Call) expression()                  {}
func (i Identifier) expression()             {}
func (cc CharCon) expression()               {}
func (ic IntCon) expression()                {}
func (fc FloatCon) expression()              {}
func (sc StringCon) expression()             {}
func (b Bool) expression()                   {}
func (a Array) expression()                  {}
func (ie IndexExpr) expression()             {}
func (aie AssignIndexExpr) expression()      {}
func (aeie AssignExprIndexExpr) expression() {}
