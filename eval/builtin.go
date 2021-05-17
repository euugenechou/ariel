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
	"ariel/object"
	"fmt"
	"math/rand"
)

var builtins = map[string]object.BuiltIn{
	"println": object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Eval())
			}
			fmt.Println()
			return nil
		},
	},
	"print": object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Eval())
			}
			return nil
		},
	},
	"rand": object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return errorObj("too many arguments to rand()")
			}
			return object.Int{Value: rand.Int63()}
		},
	},
}
