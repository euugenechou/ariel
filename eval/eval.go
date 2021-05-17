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

package eval

import (
	"ariel/ast"
	"ariel/color"
	"ariel/misc"
	"ariel/object"
	"fmt"
	"math"
)

func errorObj(format string, a ...interface{}) object.Error {
	msg := color.Red
	msg += misc.Flounder("error: " + fmt.Sprintf(format, a...))
	msg += color.Reset
	return object.Error{Message: msg}
}

func IsError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	}
	return false
}

func Eval(n ast.Node, s *object.State) object.Object {
	switch n := n.(type) {
	case ast.Program:
		return evalProgram(n, s)
	case ast.FuncDecl:
		return evalFuncDecl(n, s)
	case ast.VarDecl:
		return evalVarDecl(n, s)
	case ast.Block:
		return evalBlock(n, s)
	case ast.While:
		return evalWhile(n, s)
	case ast.For:
		return evalFor(n, s)
	case ast.IfElse:
		return evalIfElse(n, s)
	case ast.Return:
		return evalReturn(n, s)
	case ast.ExprStmt:
		return evalExprStmt(n, s)
	case ast.PrefixExpr:
		return evalPrefixExpr(n, s)
	case ast.InfixExpr:
		return evalInfixExpr(n, s)
	case ast.Assign:
		return evalAssign(n, s)
	case ast.AssignExpr:
		return evalAssignExpr(n, s)
	case ast.IndexExpr:
		return evalIndexExpr(n, s)
	case ast.AssignIndexExpr:
		return evalAssignIndexExpr(n, s)
	case ast.AssignExprIndexExpr:
		return evalAssignExprIndexExpr(n, s)
	case ast.Call:
		return evalCall(n, s)
	case ast.Array:
		return evalArray(n, s)
	case ast.CharCon:
		return evalCharCon(n)
	case ast.IntCon:
		return evalIntCon(n)
	case ast.FloatCon:
		return evalFloatCon(n)
	case ast.StringCon:
		return evalStringCon(n)
	case ast.Bool:
		return evalBool(n)
	case ast.Identifier:
		return evalIdent(n, s)
	default:
		return nil
	}
}

func evalProgram(p ast.Program, s *object.State) object.Object {
	var result object.Object
	for _, stmt := range p.Statements {
		result = Eval(stmt, s)
		switch result := result.(type) {
		case object.Return:
			return result.Value
		case object.Error:
			return result
		}
	}
	return result
}

func evalFuncDecl(fd ast.FuncDecl, s *object.State) object.Object {
	if _, ok := s.Get(fd.Ident.Name); ok {
		return errorObj("%s already declared", fd.Ident.Name)
	}

	function := object.FuncDecl{
		ReturnType: fd.Type,
		Ident:      fd.Ident,
		Parameters: fd.Parameters,
		Body:       fd.Body,
		State:      s,
	}

	s.Set(fd.Ident.Name, function)
	return nil
}

func evalVarDecl(vd ast.VarDecl, s *object.State) object.Object {
	if _, ok := s.Get(vd.Ident.Name); ok {
		return errorObj("%s already declared", vd.Ident.Name)
	}

	var val object.Object

	if vd.Initialized {
		val = Eval(vd.Value, s)
		if IsError(val) {
			return val
		}

		switch vd.Type.Value {
		case "char":
			if val.Type() != object.CharObj {
				return errorObj("mismatched types: char %s = %s",
					vd.Ident.Name, object.ObjString(val))
			}
		case "int":
			if val.Type() != object.IntObj {
				return errorObj("mismatched types: int %s = %s",
					vd.Ident.Name, object.ObjString(val))
			}
		case "float":
			if val.Type() != object.FloatObj {
				return errorObj("mismatched types: float %s = %s",
					vd.Ident.Name, object.ObjString(val))
			}
		case "string":
			if val.Type() != object.StringObj {
				return errorObj("mismatched types: string %s = %s",
					vd.Ident.Name, object.ObjString(val))
			}
		case "bool":
			if val.Type() != object.BoolObj {
				return errorObj("mismatched types: bool %s = %s",
					vd.Ident.Name, object.ObjString(val))
			}
		case "chararr":
			val = Eval(vd.Value, s)
			arr := val.(object.Array)
			elements := val.(object.Array).Elements
			val = object.Array{ElementType: "char", Elements: elements}
			if !isHeterogeneous(arr) {
				return errorObj("heterogeneous array typings: %s",
					vd.Ident.Name)
			}
			if arr.Elements[0].Type() != object.CharObj {
				return errorObj("illegal type in char array: %s",
					object.ObjString(arr.Elements[0]))
			}
		case "intarr":
			val = Eval(vd.Value, s)
			arr := val.(object.Array)
			elements := val.(object.Array).Elements
			val = object.Array{ElementType: "int", Elements: elements}
			if !isHeterogeneous(arr) {
				return errorObj("heterogeneous array typings: %s",
					vd.Ident.Name)
			}
			if arr.Elements[0].Type() != object.IntObj {
				return errorObj("illegal type in int array: %s",
					object.ObjString(arr.Elements[0]))
			}
		case "floatarr":
			val = Eval(vd.Value, s)
			arr := val.(object.Array)
			elements := val.(object.Array).Elements
			val = object.Array{ElementType: "float", Elements: elements}
			if !isHeterogeneous(arr) {
				return errorObj("heterogeneous array typings: %s",
					vd.Ident.Name)
			}
			if arr.Elements[0].Type() != object.FloatObj {
				return errorObj("illegal type in float array: %s",
					object.ObjString(arr.Elements[0]))
			}
		case "stringarr":
			val = Eval(vd.Value, s)
			arr := val.(object.Array)
			elements := val.(object.Array).Elements
			val = object.Array{ElementType: "string", Elements: elements}
			if !isHeterogeneous(arr) {
				return errorObj("heterogeneous array typings: %s",
					vd.Ident.Name)
			}
			if arr.Elements[0].Type() != object.StringObj {
				return errorObj("illegal type in string array: %s",
					object.ObjString(arr.Elements[0]))
			}
		case "boolarr":
			val = Eval(vd.Value, s)
			arr := val.(object.Array)
			elements := val.(object.Array).Elements
			val = object.Array{ElementType: "bool", Elements: elements}
			if !isHeterogeneous(arr) {
				return errorObj("heterogeneous array typings: %s",
					vd.Ident.Name)
			}
			if arr.Elements[0].Type() != object.BoolObj {
				return errorObj("illegal type in bool array: %s",
					object.ObjString(arr.Elements[0]))
			}
		default:
			return errorObj("invalid declaration type: %s %s = %s",
				vd.Type.Value, vd.Ident.Name, object.ObjString(val))
		}
	} else {
		switch vd.Type.Value {
		case "char":
			val = object.Char{Value: ""}
		case "int":
			val = object.Int{Value: 0}
		case "float":
			val = object.Float{Value: 0.0}
		case "string":
			val = object.String{Value: ""}
		case "bool":
			val = object.Bool{Value: false}
		case "chararr":
			elements := Eval(vd.Value, s)
			if IsError(elements) {
				return elements
			}
			if elements.Type() != object.IntObj {
				return errorObj("array size must be integer")
			}
			numElements := elements.(object.Int).Value
			arr := make([]object.Object, numElements)
			for i := 0; i < int(numElements); i++ {
				arr[i] = object.Char{Value: ""}
			}
			val = object.Array{ElementType: "char", Elements: arr}
		case "intarr":
			elements := Eval(vd.Value, s)
			if IsError(elements) {
				return elements
			}
			if elements.Type() != object.IntObj {
				return errorObj("array size must be integer")
			}
			numElements := elements.(object.Int).Value
			arr := make([]object.Object, numElements)
			for i := 0; i < int(numElements); i++ {
				arr[i] = object.Int{Value: 0}
			}
			val = object.Array{ElementType: "int", Elements: arr}
		case "floatarr":
			elements := Eval(vd.Value, s)
			if IsError(elements) {
				return elements
			}
			if elements.Type() != object.IntObj {
				return errorObj("array size must be integer")
			}
			numElements := elements.(object.Int).Value
			arr := make([]object.Object, numElements)
			for i := 0; i < int(numElements); i++ {
				arr[i] = object.Float{Value: 0.0}
			}
			val = object.Array{ElementType: "float", Elements: arr}
		case "stringarr":
			elements := Eval(vd.Value, s)
			if IsError(elements) {
				return elements
			}
			if elements.Type() != object.IntObj {
				return errorObj("array size must be integer")
			}
			numElements := elements.(object.Int).Value
			arr := make([]object.Object, numElements)
			for i := 0; i < int(numElements); i++ {
				arr[i] = object.String{Value: ""}
			}
			val = object.Array{ElementType: "string", Elements: arr}
		case "boolarr":
			elements := Eval(vd.Value, s)
			if IsError(elements) {
				return elements
			}
			if elements.Type() != object.IntObj {
				return errorObj("array size must be integer")
			}
			numElements := elements.(object.Int).Value
			arr := make([]object.Object, numElements)
			for i := 0; i < int(numElements); i++ {
				arr[i] = object.Bool{Value: false}
			}
			val = object.Array{ElementType: "bool", Elements: arr}
		default:
			return errorObj("invalid declaration type: %s %s = %s",
				vd.Type.Value, vd.Ident.Name, object.ObjString(val))
		}
	}

	s.Set(vd.Ident.Name, val)
	return nil
}

func isHeterogeneous(arr object.Array) bool {
	if len(arr.Elements) > 0 {
		for _, element := range arr.Elements {
			if element.Type() != arr.Elements[0].Type() {
				return false
			}
		}
	}
	return true
}

func evalBlock(b ast.Block, s *object.State) object.Object {
	var result object.Object
	copied := object.NewCopiedState(s)

	for _, stmt := range b.Statements {
		result = Eval(stmt, copied)
		if result != nil {
			switch result := result.(type) {
			case object.Return:
				return result
			case object.Error:
				return result
			}
		}
	}

	object.UpdateState(s, copied)
	return result
}

func evalWhile(w ast.While, s *object.State) object.Object {
	var result object.Object

	cond := Eval(w.Condition, s)
	if IsError(cond) {
		return cond
	}

	if cond.Type() != object.BoolObj {
		return errorObj("improper while condition type: %s",
			object.ObjString(cond))
	}

	for cond.(object.Bool).Value {
		result := Eval(w.Body, s)
		switch result := result.(type) {
		case object.Return:
			return result.Value
		case object.Error:
			return result
		}
		cond = Eval(w.Condition, s)
	}

	return result
}

func evalFor(f ast.For, s *object.State) object.Object {
	var result object.Object
	copied := object.NewCopiedState(s)

	if f.VarDecl {
		vardecl := ast.VarDecl{
			Type:        f.Type,
			Ident:       f.Ident,
			Value:       f.Value,
			Initialized: true,
		}
		declared := Eval(vardecl, copied)
		if IsError(declared) {
			return declared
		}
	} else {
		init := Eval(f.Init, copied)
		if IsError(init) {
			return init
		}
	}

	cond := Eval(f.Condition, copied)
	if IsError(cond) {
		return cond
	}

	if cond.Type() != object.BoolObj {
		return errorObj("improper for condition type: %s",
			object.ObjString(cond))
	}

	for cond.(object.Bool).Value {
		result := Eval(f.Body, copied)
		switch result := result.(type) {
		case object.Return:
			return result.Value
		case object.Error:
			return result
		}

		increment := Eval(f.Increment, copied)
		if IsError(increment) {
			return increment
		}

		cond = Eval(f.Condition, copied)
		if IsError(cond) {
			return cond
		}
	}

	object.UpdateState(s, copied)
	return result
}

func evalIfElse(ie ast.IfElse, s *object.State) object.Object {
	cond := Eval(ie.Condition, s)
	if IsError(cond) {
		return cond
	}

	if cond.Type() != object.BoolObj {
		return errorObj("improper if condition type: %s",
			object.ObjString(cond))
	}

	if cond.(object.Bool).Value {
		return Eval(ie.Consequence, s)
	} else if ie.HasAlternative {
		return Eval(ie.Alternative, s)
	} else {
		return nil
	}
}

func evalReturn(ie ast.Return, s *object.State) object.Object {
	if !ie.Void {
		val := Eval(ie.Value, s)
		if IsError(val) {
			return val
		}
		return object.Return{Value: val}
	}
	return nil
}

func evalExprStmt(es ast.ExprStmt, s *object.State) object.Object {
	return Eval(es.Expression, s)
}

func evalPrefixExpr(pe ast.PrefixExpr, s *object.State) object.Object {
	right := Eval(pe.Right, s)
	if IsError(right) {
		return right
	}

	switch pe.Op {
	case "!":
		return evalNotOp(right)
	case "-":
		return evalNegOp(right)
	case "+":
		return evalPosOp(right)
	case "~":
		return evalTildeOp(right)
	default:
		return errorObj("unknown prefix operator: %s", pe.Op)
	}
}

func evalNotOp(expr object.Object) object.Object {
	if expr.Type() != object.BoolObj {
		return errorObj("illegal operation: !%s", object.ObjString(expr))
	}
	if expr.(object.Bool).Value {
		return object.Bool{Value: false}
	}
	return object.Bool{Value: true}
}

func evalNegOp(expr object.Object) object.Object {
	switch expr := expr.(type) {
	case object.Int:
		return object.Int{Value: -expr.Value}
	case object.Float:
		return object.Float{Value: -expr.Value}
	default:
		return errorObj("illegal operation: -%s", object.ObjString(expr))
	}
}

func evalPosOp(expr object.Object) object.Object {
	switch expr := expr.(type) {
	case object.Int:
		abs := expr.Value
		if abs < 0 {
			abs = -abs
		}
		return object.Int{Value: abs}
	case object.Float:
		return object.Float{Value: math.Abs(expr.Value)}
	default:
		return errorObj("illegal operation: +%s", object.ObjString(expr))
	}
}

func evalTildeOp(expr object.Object) object.Object {
	switch expr := expr.(type) {
	case object.Int:
		return object.Int{Value: ^expr.Value}
	default:
		return errorObj("illegal operation: ~%s", object.ObjString(expr))
	}
}

func evalInfixExpr(ie ast.InfixExpr, s *object.State) object.Object {
	left := Eval(ie.Left, s)
	if IsError(left) {
		return left
	}

	right := Eval(ie.Right, s)
	if IsError(right) {
		return right
	}

	if left.Type() != right.Type() {
		return errorObj("mismatched types: %s %s %s",
			object.ObjString(left), ie.Op, object.ObjString(right))
	}

	switch right.Type() {
	case object.CharObj:
		return evalInfixExprChar(ie.Op, left.(object.Char), right.(object.Char))
	case object.IntObj:
		return evalInfixExprInt(ie.Op, left.(object.Int), right.(object.Int))
	case object.FloatObj:
		return evalInfixExprFloat(ie.Op, left.(object.Float), right.(object.Float))
	case object.StringObj:
		return evalInfixExprString(ie.Op, left.(object.String), right.(object.String))
	case object.BoolObj:
		return evalInfixExprBool(ie.Op, left.(object.Bool), right.(object.Bool))
	default:
		return errorObj("invalid expression types: %s %s %s",
			object.ObjString(left), ie.Op, object.ObjString(right))
	}
}

func evalInfixExprChar(op string, left, right object.Char) object.Object {
	switch op {
	case "<":
		return object.Bool{Value: left.Value < right.Value}
	case "<=":
		return object.Bool{Value: left.Value <= right.Value}
	case "==":
		return object.Bool{Value: left.Value == right.Value}
	case "!=":
		return object.Bool{Value: left.Value != right.Value}
	case ">=":
		return object.Bool{Value: left.Value >= right.Value}
	case ">":
		return object.Bool{Value: left.Value > right.Value}
	case "+":
		return object.String{Value: left.Value + right.Value}
	default:
		return errorObj("illegal operator: %s %s %s",
			object.ObjString(left), op, object.ObjString(right))
	}
}

func evalInfixExprInt(op string, left, right object.Int) object.Object {
	switch op {
	case "<":
		return object.Bool{Value: left.Value < right.Value}
	case "<=":
		return object.Bool{Value: left.Value <= right.Value}
	case "==":
		return object.Bool{Value: left.Value == right.Value}
	case "!=":
		return object.Bool{Value: left.Value != right.Value}
	case ">=":
		return object.Bool{Value: left.Value >= right.Value}
	case ">":
		return object.Bool{Value: left.Value > right.Value}
	case "+":
		return object.Int{Value: left.Value + right.Value}
	case "-":
		return object.Int{Value: left.Value - right.Value}
	case "*":
		return object.Int{Value: left.Value * right.Value}
	case "/":
		if right.Value == 0 {
			return errorObj("divide by zero error")
		}
		return object.Int{Value: left.Value / right.Value}
	case "%":
		return object.Int{Value: left.Value % right.Value}
	case "&":
		return object.Int{Value: left.Value & right.Value}
	case "^":
		return object.Int{Value: left.Value ^ right.Value}
	case "|":
		return object.Int{Value: left.Value | right.Value}
	case "<<":
		return object.Int{Value: left.Value << right.Value}
	case ">>":
		return object.Int{Value: left.Value >> right.Value}
	default:
		return errorObj("illegal operator: %s %s %s",
			object.ObjString(left), op, object.ObjString(right))
	}
}

func evalInfixExprFloat(op string, left, right object.Float) object.Object {
	switch op {
	case "<":
		return object.Bool{Value: left.Value < right.Value}
	case "<=":
		return object.Bool{Value: left.Value <= right.Value}
	case "==":
		return object.Bool{Value: left.Value == right.Value}
	case "!=":
		return object.Bool{Value: left.Value != right.Value}
	case ">=":
		return object.Bool{Value: left.Value >= right.Value}
	case ">":
		return object.Bool{Value: left.Value > right.Value}
	case "+":
		return object.Float{Value: left.Value + right.Value}
	case "-":
		return object.Float{Value: left.Value - right.Value}
	case "*":
		return object.Float{Value: left.Value * right.Value}
	case "/":
		if right.Value == 0.0 {
			return errorObj("divide by zero error")
		}
		return object.Float{Value: left.Value / right.Value}
	default:
		return errorObj("illegal operator: %s %s %s",
			object.ObjString(left), op, object.ObjString(right))
	}
}

func evalInfixExprString(op string, left, right object.String) object.Object {
	switch op {
	case "<":
		return object.Bool{Value: left.Value < right.Value}
	case "<=":
		return object.Bool{Value: left.Value <= right.Value}
	case "==":
		return object.Bool{Value: left.Value == right.Value}
	case "!=":
		return object.Bool{Value: left.Value != right.Value}
	case ">=":
		return object.Bool{Value: left.Value >= right.Value}
	case ">":
		return object.Bool{Value: left.Value > right.Value}
	case "+":
		return object.String{Value: left.Value + right.Value}
	default:
		return errorObj("illegal operator: %s %s %s",
			object.ObjString(left), op, object.ObjString(right))
	}
}

func evalInfixExprBool(op string, left, right object.Bool) object.Object {
	switch op {
	case "==":
		return object.Bool{Value: left.Value == right.Value}
	case "!=":
		return object.Bool{Value: left.Value != right.Value}
	case "&&":
		return object.Bool{Value: left.Value && right.Value}
	case "||":
		return object.Bool{Value: left.Value || right.Value}
	default:
		return errorObj("illegal operator: %s %s %s",
			object.ObjString(left), op, object.ObjString(right))
	}
}

func evalAssign(a ast.Assign, s *object.State) object.Object {
	ident := Eval(a.Ident, s)
	if IsError(ident) {
		return ident
	}

	val := Eval(a.Value, s)
	if IsError(val) {
		return val
	}

	if ident.Type() != val.Type() {
		return errorObj("assignment type mismatch: %s and %s",
			object.ObjString(ident), object.ObjString(val))
	}

	s.Set(a.Ident.Name, val)
	return nil
}

func evalAssignExpr(ae ast.AssignExpr, s *object.State) object.Object {
	ident := Eval(ae.Ident, s)
	if IsError(ident) {
		return ident
	}

	val := Eval(ae.Value, s)
	if IsError(val) {
		return val
	}

	self, _ := s.Get(ae.Ident.Name)
	if self.Type() != val.Type() {
		return errorObj("mismatched types: %s %s %s",
			object.ObjString(self), ae.Op, object.ObjString(val))
	}

	var op string
	var newVal object.Object

	switch ae.Op {
	case "+=":
		op = "+"
	case "-=":
		op = "-"
	case "*=":
		op = "*"
	case "/=":
		op = "/"
	case "%=":
		op = "%"
	case "&=":
		op = "&"
	case "^=":
		op = "^"
	case "|=":
		op = "|"
	case "<<=":
		op = "<<"
	case ">>=":
		op = ">>"
	default:
		return errorObj("illegal operator: %s %s %s",
			object.ObjString(self), ae.Op, object.ObjString(val))
	}

	switch val.Type() {
	case object.IntObj:
		newVal = evalInfixExprInt(op, self.(object.Int), val.(object.Int))
	case object.FloatObj:
		newVal = evalInfixExprFloat(op, self.(object.Float), val.(object.Float))
	case object.StringObj:
		newVal = evalInfixExprString(op, self.(object.String), val.(object.String))
	default:
		return errorObj("illegal assignment: %s %s %s",
			object.ObjString(self), ae.Op, object.ObjString(val))
	}

	if IsError(newVal) {
		return newVal
	}

	s.Set(ae.Ident.Name, newVal)
	return nil
}

func evalIndexExpr(ie ast.IndexExpr, s *object.State) object.Object {
	ident := evalIdent(ie.Ident, s)
	if IsError(ident) {
		return ident
	}

	index := Eval(ie.Index, s)
	if IsError(index) {
		return index
	}
	if index.Type() != object.IntObj {
		return errorObj("illegal array index: %s", object.ObjString(index))
	}

	array, ok := s.Get(ie.Ident.Name)
	if !ok {
		return errorObj("undeclared array: %s", ie.Ident.Name)
	}
	if array.Type() != object.ArrObj {
		return errorObj("%s is not an array", ie.Ident.Name)
	}

	idx := index.(object.Int).Value
	arr := array.(object.Array).Elements

	if int(idx) < 0 || int(idx) >= len(arr) {
		return errorObj("array index out of bounds: %s[%d]",
			ie.Ident.Name, idx)
	}

	return arr[idx]
}

func evalAssignIndexExpr(aie ast.AssignIndexExpr, s *object.State) object.Object {
	ident := evalIdent(aie.Ident, s)
	if IsError(ident) {
		return ident
	}

	index := Eval(aie.Index, s)
	if IsError(index) {
		return index
	}
	if index.Type() != object.IntObj {
		return errorObj("illegal array index: %s", object.ObjString(index))
	}

	array, ok := s.Get(aie.Ident.Name)
	if !ok {
		return errorObj("undeclared array: %s", aie.Ident.Name)
	}
	if array.Type() != object.ArrObj {
		return errorObj("%s is not an array", aie.Ident.Name)
	}

	idx := index.(object.Int).Value
	arr := array.(object.Array).Elements

	if int(idx) < 0 || int(idx) >= len(arr) {
		return errorObj("array index out of bounds: %s[%d]",
			aie.Ident.Name, idx)
	}

	val := Eval(aie.Value, s)
	if IsError(val) {
		return val
	}

	if arr[idx].Type() != val.Type() {
		return errorObj("assignment type mismatch: %s and %s",
			object.ObjString(arr[idx]), object.ObjString(val))
	}

	arr[idx] = val
	elementType := array.(object.Array).ElementType
	s.Set(aie.Ident.Name, object.Array{ElementType: elementType, Elements: arr})
	return nil
}

func evalAssignExprIndexExpr(aeie ast.AssignExprIndexExpr, s *object.State) object.Object {
	ident := evalIdent(aeie.Ident, s)
	if IsError(ident) {
		return ident
	}

	index := Eval(aeie.Index, s)
	if IsError(index) {
		return index
	}
	if index.Type() != object.IntObj {
		return errorObj("illegal array index: %s", object.ObjString(index))
	}

	array, ok := s.Get(aeie.Ident.Name)
	if !ok {
		return errorObj("undeclared array: %s", aeie.Ident.Name)
	}
	if array.Type() != object.ArrObj {
		return errorObj("%s is not an array", aeie.Ident.Name)
	}

	idx := index.(object.Int).Value
	arr := array.(object.Array).Elements

	if int(idx) < 0 || int(idx) >= len(arr) {
		return errorObj("array index out of bounds: %s[%d]",
			aeie.Ident.Name, idx)
	}

	val := Eval(aeie.Value, s)
	if IsError(val) {
		return val
	}

	if arr[idx].Type() != val.Type() {
		return errorObj("assignment type mismatch: %s and %s",
			object.ObjString(arr[idx]), object.ObjString(val))
	}

	var op string
	var newVal object.Object

	switch aeie.Op {
	case "=":
		op = "="
	case "+=":
		op = "+"
	case "-=":
		op = "-"
	case "*=":
		op = "*"
	case "/=":
		op = "/"
	case "%=":
		op = "%"
	case "&=":
		op = "&"
	case "^=":
		op = "^"
	case "|=":
		op = "|"
	case "<<=":
		op = "<<"
	case ">>=":
		op = ">>"
	default:
		return errorObj("illegal operator: %s %s %s",
			object.ObjString(arr[idx]), aeie.Op, object.ObjString(val))
	}

	if op == "=" {
		newVal = val
	} else {
		switch val.Type() {
		case object.IntObj:
			newVal = evalInfixExprInt(op, arr[idx].(object.Int), val.(object.Int))
		case object.FloatObj:
			newVal = evalInfixExprFloat(op, arr[idx].(object.Float), val.(object.Float))
		case object.StringObj:
			newVal = evalInfixExprString(op, arr[idx].(object.String), val.(object.String))
		default:
			return errorObj("illegal assignment: %s %s %s",
				object.ObjString(arr[idx]), aeie.Op, object.ObjString(val))
		}
	}

	arr[idx] = newVal
	elementType := array.(object.Array).ElementType
	s.Set(aeie.Ident.Name, object.Array{ElementType: elementType, Elements: arr})
	return nil
}

func evalCall(c ast.Call, s *object.State) object.Object {
	function := Eval(c.Function, s)
	if IsError(function) {
		return function
	}

	isBuiltin := function.Type() == object.BuiltInObj
	isFunction := function.Type() == object.FuncDeclObj
	if !isBuiltin && !isFunction {
		return errorObj("%s is not a declared or built-in function",
			c.Function.Name)
	}

	args := evalExpressions(c.Arguments, s)
	if len(args) == 1 && IsError(args[0]) {
		return args[0]
	}

	if function.Type() == object.FuncDeclObj {
		decl, _ := s.Get(c.Function.Name)
		params := decl.(object.FuncDecl).Parameters
		if len(args) > len(params) {
			return errorObj("too many arguments supplied to %s()",
				c.Function.Name)
		} else if len(args) < len(params) {
			return errorObj("not enough arguments supplied to %s()",
				c.Function.Name)
		}

		for i := 0; i < len(params); i++ {
			if params[i].Array {
				if args[i].Type() != object.ArrObj {
					return errorObj("passed non-array as array parameter")
				}
				elementType := args[i].(object.Array).ElementType
				if elementType != params[i].Type.Value {
					return errorObj("mismatched types for argument %d", i+1)
				}
			} else {
				switch args[i].Type() {
				case object.IntObj:
					if params[i].Type.Value != "int" {
						return errorObj("mismatched types for argument %d", i+1)
					}
				case object.FloatObj:
					if params[i].Type.Value != "float" {
						return errorObj("mismatched types for argument %d", i+1)
					}
				case object.StringObj:
					if params[i].Type.Value != "string" {
						return errorObj("mismatched types for argument %d", i+1)
					}
				case object.BoolObj:
					if params[i].Type.Value != "bool" {
						return errorObj("mismatched types for argument %d", i+1)
					}
				default:
					return errorObj("illegal type for argument %d", i)
				}
			}
		}
	}

	switch function := function.(type) {
	case object.FuncDecl:
		callState := object.NewState()
		for i, param := range function.Parameters {
			callState.Set(param.Ident.Name, args[i])
		}
		object.CopyFunctions(callState, s)

		evaluated := Eval(function.Body, callState)
		if evaluated != nil && evaluated.Type() == object.ReturnObj {
			return evaluated.(object.Return).Value
		}
		return evaluated
	case object.BuiltIn:
		return function.Function(args...)
	default:
		return errorObj("not a function: %s", c.Function.Name)
	}
}

func evalArray(a ast.Array, s *object.State) object.Object {
	elements := evalExpressions(a.Elements, s)
	if len(elements) == 1 && IsError(elements[0]) {
		return elements[0]
	}
	return object.Array{Elements: elements}
}

func evalExpressions(args []ast.Expression, s *object.State) []object.Object {
	var result []object.Object
	for _, arg := range args {
		evaluated := Eval(arg, s)
		if IsError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalCharCon(cc ast.CharCon) object.Object {
	return object.Char{Value: cc.Value}
}

func evalIntCon(ic ast.IntCon) object.Object {
	return object.Int{Value: ic.Value}
}

func evalFloatCon(fc ast.FloatCon) object.Object {
	return object.Float{Value: fc.Value}
}

func evalStringCon(sc ast.StringCon) object.Object {
	return object.String{Value: sc.Value}
}

func evalBool(b ast.Bool) object.Object {
	return object.Bool{Value: b.Value}
}

func evalIdent(i ast.Identifier, s *object.State) object.Object {
	if val, ok := s.Get(i.Name); ok {
		return val
	}

	if function, ok := builtins[i.Name]; ok {
		return function
	}

	return errorObj("identifier %s undeclared", i.Name)
}
