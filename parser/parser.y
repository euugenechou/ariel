%{
package parser

import (
    "ariel/ast"
    "ariel/color"
    "ariel/misc"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
	"text/scanner"
)

type Token struct {
	Literal string
    Int     int64
    Float   float64
    Bool    bool
}
%}

%union{
    token Token
    Program ast.Program
    DeclList []ast.Statement
    Decl ast.Statement
    Type ast.Type
    FuncDecl ast.FuncDecl
    VarDecl ast.VarDecl
    ParamList []ast.Param
    Param ast.Param
    Block ast.Block
    StmtList []ast.Statement
    Stmt ast.Statement
    While ast.While
    For ast.For
    IfElse ast.IfElse
    Return ast.Return
    ExprStmt ast.ExprStmt
    Expr ast.Expression
    Call ast.Call
    Array ast.Array
    ExprList []ast.Expression
    Id ast.Identifier
}

%token<token> CHAR INT FLOAT STRING BOOL VOID
%token<token> WHILE FOR IF ELSE RETURN
%token<token> '+' '-' '*' '/' '%' '&' '^' '|' '=' '!' LT LE EQ GE GT AND OR
%token<token> ADD SUB MUL DIV MOD RSHIFT LSHIFT
%token<token> ADDS SUBS MULS DIVS MODS LSHIFTS RSHIFTS
%token<token> ID CHARCON INTCON, STRINGCON, FLOATCON, TRUE, FALSE
%token<token> '(' ')' '{' '}' '[' ']' ';'

%type<Program> Program
%type<DeclList> DeclList
%type<Decl> Decl
%type<Type> Type
%type<FuncDecl> FuncDecl
%type<VarDecl> VarDecl
%type<ParamList> ParamList
%type<Param> Param
%type<Block> Block
%type<StmtList> StmtList
%type<Stmt> Stmt
%type<While> While
%type<For> For
%type<IfElse> IfElse
%type<Return> Return
%type<ExprStmt> ExprStmt
%type<Expr> Expr
%type<Call> Call
%type<Array> Array
%type<ExprList> ExprList
%type<Id> Id

%nonassoc '{' '}'
%right IF ELSE
%right '=' ADDS SUBS MULS DIVS MODS ANDS XORS ORS LSHIFTS RSHIFTS
%left OR
%left AND
%left '|'
%left '^'
%left '&'
%left EQ NE
%left '<' LE '>' GE
%left LSHIFT RSHIFT
%left '+' '-'
%left '*' '/' '%'
%right NEG POS NOT TILDE
%left '(' ')' '[' ']'

%start Program

%%
Program
    : DeclList { $$ = ast.Program{Statements: $1}; yylex.(*Lexer).result = $$ }
    ;

DeclList
    : Decl          { $$ = []ast.Statement{$1} }
    | DeclList Decl { $$ = append($1, $2) }
    ;

Decl
    : Stmt      { $$ = $1 }
    | FuncDecl  { $$ = $1 }
    ;

Type
    : CHAR    { $$ = ast.Type{Value: $1.Literal} }
    | INT     { $$ = ast.Type{Value: $1.Literal} }
    | FLOAT   { $$ = ast.Type{Value: $1.Literal} }
    | STRING  { $$ = ast.Type{Value: $1.Literal} }
    | BOOL    { $$ = ast.Type{Value: $1.Literal} }
    | VOID    { $$ = ast.Type{Value: $1.Literal} }
    ;

FuncDecl
    : Type Id '(' ')' Block {
        $$ = ast.FuncDecl{
            Type: $1,
            Ident: $2,
            Parameters: make([]ast.Param, 0),
            Body: $5,
        }
    }
    | Type Id '(' ParamList ')' Block {
        $$ = ast.FuncDecl{
            Type: $1,
            Ident: $2,
            Parameters: $4,
            Body: $6,
        }
    }
    ;

VarDecl
    : Type Id ';' {
        $$ = ast.VarDecl{
            Type: $1,
            Ident: $2,
            Initialized: false,
        }
    }
    | Type Id '=' Expr ';' {
        $$ = ast.VarDecl{
            Type: $1,
            Ident: $2,
            Value: $4,
            Initialized: true,
        }
    }
    | Type Id '[' Expr ']' ';' {
        $$ = ast.VarDecl {
            Type: ast.Type{Value: $1.Value + "arr"},
            Ident: $2,
            Value: $4,
            Initialized: false,
        }
    }
    | Type Id '[' ']' '=' Array ';' {
        $$ = ast.VarDecl {
            Type: ast.Type{Value: $1.Value + "arr"},
            Ident: $2,
            Value: $6,
            Initialized: true,
        }
    }
    ;

ParamList
    : Param                 { $$ = []ast.Param{$1} }
    | ParamList ',' Param   { $$ = append($1, $3) }
    ;

Param
    : Type Id           { $$ = ast.Param{Type: $1, Ident: $2, Array: false} }
    | Type Id '[' ']'   { $$ = ast.Param{Type: $1, Ident: $2, Array: true} }
    ;

Block
    : '{' '}' {
        $$ = ast.Block{Statements: make([]ast.Statement, 0)}
    }
    | '{' StmtList '}' {
        $$ = ast.Block{Statements: $2}
    }
    ;

StmtList
    : Stmt           { $$ = []ast.Statement{$1} }
    | StmtList Stmt  { $$ = append($1, $2) }
    ;

Stmt
    : VarDecl   { $$ = $1 }
    | Block     { $$ = $1 }
    | While     { $$ = $1 }
    | For       { $$ = $1 }
    | IfElse    { $$ = $1 }
    | Return    { $$ = $1 }
    | ExprStmt  { $$ = $1 }
    ;

While
    : WHILE '(' Expr ')' Stmt {
        $$ = ast.While{
            Condition: $3,
            Body: $5,
        }
    }
    ;

For
    : FOR '(' Expr ';' Expr ';' Expr ')' Block {
        $$ = ast.For {
            Init: $3,
            Condition: $5,
            Increment: $7,
            Body: $9,
        }
    }
    | FOR '(' Type Id '=' Expr ';' Expr ';' Expr ')' Block {
        $$ = ast.For {
            VarDecl: true,
            Type: $3,
            Ident: $4,
            Value: $6,
            Condition: $8,
            Increment: $10,
            Body: $12,
        }
    }
    ;

IfElse
    : IF '(' Expr ')' Stmt %prec IF {
        $$ = ast.IfElse{
            Condition: $3,
            Consequence: $5,
            HasAlternative: false,
        }
    }
    | IF '(' Expr ')' Stmt ELSE Stmt {
        $$ = ast.IfElse{
            Condition: $3,
            Consequence: $5,
            Alternative: $7,
            HasAlternative: true,
        }
    }
    ;

Return
    : RETURN ';'      { $$ = ast.Return{Void: true} }
    | RETURN Expr ';' { $$ = ast.Return{Value: $2, Void: false} }
    ;

ExprStmt
    : Expr ';' { $$ = ast.ExprStmt{Expression: $1} }
    ;

Expr
    : Expr '+' Expr            { $$ = ast.InfixExpr{Left: $1, Op: "+", Right: $3} }
    | Expr '-' Expr            { $$ = ast.InfixExpr{Left: $1, Op: "-", Right: $3} }
    | Expr '*' Expr            { $$ = ast.InfixExpr{Left: $1, Op: "*", Right: $3} }
    | Expr '/' Expr            { $$ = ast.InfixExpr{Left: $1, Op: "/", Right: $3} }
    | Expr '%' Expr            { $$ = ast.InfixExpr{Left: $1, Op: "%", Right: $3} }
    | Expr '&' Expr            { $$ = ast.InfixExpr{Left: $1, Op: "&", Right: $3} }
    | Expr '^' Expr            { $$ = ast.InfixExpr{Left: $1, Op: "^", Right: $3} }
    | Expr '|' Expr            { $$ = ast.InfixExpr{Left: $1, Op: "|", Right: $3} }
    | Expr LSHIFT Expr         { $$ = ast.InfixExpr{Left: $1, Op: "<<", Right: $3} }
    | Expr RSHIFT Expr         { $$ = ast.InfixExpr{Left: $1, Op: ">>", Right: $3} }
    | Expr '<' Expr            { $$ = ast.InfixExpr{Left: $1, Op: "<", Right: $3} }
    | Expr LE Expr             { $$ = ast.InfixExpr{Left: $1, Op: "<=", Right: $3} }
    | Expr EQ Expr             { $$ = ast.InfixExpr{Left: $1, Op: "==", Right: $3} }
    | Expr NE Expr             { $$ = ast.InfixExpr{Left: $1, Op: "!=", Right: $3} }
    | Expr GE Expr             { $$ = ast.InfixExpr{Left: $1, Op: ">=", Right: $3} }
    | Expr '>' Expr            { $$ = ast.InfixExpr{Left: $1, Op: ">", Right: $3} }
    | Expr AND Expr            { $$ = ast.InfixExpr{Left: $1, Op: "&&", Right: $3} }
    | Expr OR Expr             { $$ = ast.InfixExpr{Left: $1, Op: "||", Right: $3} }
    | Id '=' Expr              { $$ = ast.Assign{Ident: $1, Value: $3} }
    | Id ADDS Expr             { $$ = ast.AssignExpr{Ident: $1, Op: "+=", Value: $3} }
    | Id SUBS Expr             { $$ = ast.AssignExpr{Ident: $1, Op: "-=", Value: $3} }
    | Id MULS Expr             { $$ = ast.AssignExpr{Ident: $1, Op: "*=", Value: $3} }
    | Id DIVS Expr             { $$ = ast.AssignExpr{Ident: $1, Op: "/=", Value: $3} }
    | Id MODS Expr             { $$ = ast.AssignExpr{Ident: $1, Op: "%=", Value: $3} }
    | Id ANDS Expr             { $$ = ast.AssignExpr{Ident: $1, Op: "&=", Value: $3} }
    | Id XORS Expr             { $$ = ast.AssignExpr{Ident: $1, Op: "^=", Value: $3} }
    | Id ORS Expr              { $$ = ast.AssignExpr{Ident: $1, Op: "|=", Value: $3} }
    | Id LSHIFTS Expr          { $$ = ast.AssignExpr{Ident: $1, Op: "<<=", Value: $3} }
    | Id RSHIFTS Expr          { $$ = ast.AssignExpr{Ident: $1, Op: ">>=", Value: $3} }
    | '-' Expr %prec NEG       { $$ = ast.PrefixExpr{Op: "-", Right: $2} }
    | '+' Expr %prec POS       { $$ = ast.PrefixExpr{Op: "+", Right: $2} }
    | '!' Expr %prec NOT       { $$ = ast.PrefixExpr{Op: "!", Right: $2} }
    | '~' Expr %prec TILDE     { $$ = ast.PrefixExpr{Op: "~", Right: $2} }
    | '(' Expr ')'             { $$ = $2 }
    | Call                     { $$ = $1 }
    | Id                       { $$ = $1 }
    | CHARCON                  { $$ = ast.CharCon{Value: $1.Literal} }
    | INTCON                   { $$ = ast.IntCon{Value: $1.Int} }
    | FLOATCON                 { $$ = ast.FloatCon{Value: $1.Float} }
    | STRINGCON                { $$ = ast.StringCon{Value: $1.Literal} }
    | TRUE                     { $$ = ast.Bool{Value: $1.Bool} }
    | FALSE                    { $$ = ast.Bool{Value: $1.Bool} }
    | Id '[' Expr ']'          { $$ = ast.IndexExpr{Ident: $1, Index: $3} }
    | Id '[' Expr ']' '=' Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' ADDS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "+=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' SUBS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "-=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' MULS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "*=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' DIVS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "/=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' MODS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "%=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' ANDS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "&=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' XORS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "^=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' ORS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "|=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' LSHIFTS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: "<<=",
            Value: $6,
        }
    }
    | Id '[' Expr ']' RSHIFTS Expr {
        $$ = ast.AssignExprIndexExpr{
            Ident: $1,
            Index: $3,
            Op: ">>=",
            Value: $6,
        }
    }
    ;

Call
    : Id '(' ')' {
        $$ = ast.Call{Function: $1, Void: true}
    }
    | Id '(' ExprList ')' {
        $$ = ast.Call{Function: $1, Arguments: $3, Void: false}
    }
    ;

Array
    : '{' ExprList '}' { $$ = ast.Array{Elements: $2} }
    ;

ExprList
    : Expr                  { $$ = []ast.Expression{$1} }
    | ExprList ',' Expr     { $$ = append($1, $3) }
    ;

Id
    : ID { $$ = ast.Identifier{Name: $1.Literal} }
    ;
%%

type Lexer struct {
	scanner.Scanner
	result ast.Program
    debug bool
}

func (l *Lexer) Lex(lval *yySymType) int {
    var ttype int
	token := l.Scan()
    lit := l.TokenText()
    tok := int(token)

	switch tok {
    case scanner.Ident:
        ttype = ID
    case scanner.Char:
        ttype = CHARCON
	case scanner.Int:
		ttype = INTCON
    case scanner.Float:
        ttype = FLOATCON
    case scanner.String:
        ttype = STRINGCON
    default:
        ttype = tok
    }

    if strings.Contains("<!=>+-*/%", lit) {
        if l.Peek() == '=' {
            l.Next()
            lit += "="
        }
    }

    if lit == "<" {
        if l.Peek() == '<' {
            l.Next()
            lit += "<"
            if l.Peek() == '=' {
                l.Next()
                lit += "="
            }
        }
    }

    if lit == ">" {
        if l.Peek() == '>' {
            l.Next()
            lit += ">"
            if l.Peek() == '=' {
                l.Next()
                lit += "="
            }
        }
    }

    if lit == "&" {
        if l.Peek() == '&' {
            l.Next()
            lit += "&"
        } else if l.Peek() == '=' {
            l.Next()
            lit += "="
        }
    }

    if lit == "^" {
        if l.Peek() == '=' {
            l.Next()
            lit += "="
        }
    }

    if lit == "|" {
        if l.Peek() == '|' {
            l.Next()
            lit += "|"
        } else if l.Peek() == '=' {
            l.Next()
            lit += "="
        }
    }

    var reserved = map[string]int{
        "char":   CHAR,
        "int":    INT,
        "float":  FLOAT,
        "string": STRING,
        "bool":   BOOL,
        "void":   VOID,
        "while":  WHILE,
        "for":    FOR,
        "if":     IF,
        "else":   ELSE,
        "return": RETURN,
        "true":   TRUE,
        "false":  FALSE,
        "<=":     LE,
        "==":     EQ,
        "!=":     NE,
        ">=":     GE,
        "+=":     ADDS,
        "-=":     SUBS,
        "*=":     MULS,
        "/=":     DIVS,
        "%=":     MODS,
        "&=":     ANDS,
        "^=":     XORS,
        "|=":     ORS,
        "&&":     AND,
        "||":     OR,
        "<<":     LSHIFT,
        ">>":     RSHIFT,
        "<<=":    LSHIFTS,
        ">>=":    RSHIFTS,
    }

    if t, ok := reserved[lit]; ok {
        ttype = t
    }

    lval.token = Token{Literal: lit}

    switch ttype {
    case INTCON:
        if i, err := strconv.ParseInt(lit, 10, 64); err == nil {
            lval.token.Int = i
        }
    case FLOATCON:
        if f, err := strconv.ParseFloat(lit, 64); err == nil {
            lval.token.Float = f
        }
    case STRINGCON:
        lval.token.Literal = strings.TrimPrefix(lval.token.Literal, "\"")
        lval.token.Literal = strings.TrimSuffix(lval.token.Literal, "\"")
    case CHARCON:
        lval.token.Literal = strings.TrimPrefix(lval.token.Literal, "'")
        lval.token.Literal = strings.TrimSuffix(lval.token.Literal, "'")
    case TRUE:
        lval.token.Bool = true
    case FALSE:
        lval.token.Bool = false
    }

    if l.debug && len(strings.TrimSpace(lit)) > 0 {
        fmt.Println(lit, "\t\t\t", l.Position)
    }

	return ttype
}

func (l *Lexer) Error(e string) {
    err := fmt.Sprintf("%s: line %d, column %d",
        e, l.Position.Line, l.Position.Column)
    fmt.Fprintln(os.Stderr, color.Red + misc.Flounder(err) + color.Reset)
}

func ParseProgram(input io.Reader, debug bool) ast.Program {
	l := new(Lexer)
    l.debug = debug
	l.Init(input)
    l.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars
    l.Mode |= scanner.ScanStrings | scanner.SkipComments
	yyParse(l)
    return l.result
}

func ParseProgramString(input string, debug bool) ast.Program {
	l := new(Lexer)
    l.debug = debug
	l.Init(strings.NewReader(input))
    l.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars
    l.Mode |= scanner.ScanStrings | scanner.SkipComments
	yyParse(l)
    return l.result
}
