package golisp

import (
	"testing"
	//"bytes"
	"fmt"
	"strings"
)

type Function struct {
	Params []string
}

func InitFunction() (fn Function) {
	fn.Params = make([]string, 1)
	return 
}

type State struct {
	Fns map[string]Function
	Error string
}

func InitState() State {
	var state State
	state.Fns = make(map[string]Function)
	return state
}

func Compile(src string, state *State) {
	state.Error = "Unbalanced parens"
	for _, ch := range src {
		if ch == ')' {
			state.Error = ""
		}
	}
	tokens := strings.Split(src, " ")
	if tokens[0] == "(defn" {
		fnName := tokens[1]
		state.Fns[fnName] = InitFunction()
	}
}

func TestAddsANewFunctionToState(t *testing.T) {
	state := InitState()
	Compile("(defn f [])", &state)
	if _, hasFn := state.Fns["f"]; !hasFn {
		t.Errorf("Function 'f' not added, but should have been.")
	}
}

func TestAddsADifferentlyNamedFunctionToState(t *testing.T) {
	state := InitState()
	Compile("(defn a [])", &state)
	if _, hasFn := state.Fns["a"]; !hasFn {
		t.Errorf("Function 'a' not added, but should have been.")
	}
}

func TestPlusDoesNotCreateAFunction(t *testing.T) {
	state := InitState()
	Compile("(+ 1 1)", &state)
	if len(state.Fns) != 0 {
		t.Errorf("No functions should have been added.")
	}
}

func TestErrorIfUnbalancedParens(t *testing.T) {
	state := InitState()
	Compile("(defn one []", &state)
	if state.Error != "Unbalanced parens" {
		t.Errorf("Should have error'd on unbalanced parens.")
	}
}

func TestCanDetermineFnParamNames(t *testing.T) {
	state := InitState()
	Compile("(defn f [a b c])", &state)
	fn := state.Fns["f"]
	if len(fn.Params) != 3 || ("a" != fn.Params[0] && "b" != fn.Params[1] && "c" != fn.Params[2]) {
		t.Errorf(fmt.Sprintf("Failed to determine parameters. Instead of [a b c], got %v\n", fn.Params))
	}
}

