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

type State struct {
	store map[string]Object
}

func NewState() *State {
	store := make(map[string]Object)
	return &State{store: store}
}

func NewCopiedState(s *State) *State {
	store := make(map[string]Object)
	for key, val := range s.store {
		store[key] = val
	}
	return &State{store: store}
}

func UpdateState(s *State, t *State) {
	for key, _ := range s.store {
		s.store[key] = t.store[key]
	}
}

func CopyFunctions(s *State, t *State) {
	for key, val := range t.store {
		if val.Type() == FuncDeclObj {
			s.store[key] = t.store[key]
		}
	}
}

func (s *State) Get(id string) (Object, bool) {
	obj, ok := s.store[id]
	return obj, ok
}

func (s *State) Set(id string, val Object) Object {
	s.store[id] = val
	return val
}
