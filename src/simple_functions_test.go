package golisp

import (
	"testing"
	//"bytes"
	"fmt"
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

func getTokens(src string) (ts []string) {
	ts = make([]string, 0)
	for _, ch := range src {
		if ch == ' ' {
			ts = append(ts, "")
		} else if ch == '(' || ch == ')' || ch == '[' || ch == ']' {
			if len(ts) > 0 && ts[len(ts)-1] == "" {
				ts = ts[0:len(ts)-1]
			}
			ts = append(ts, string(ch))
			ts = append(ts, "")
		} else {
			ts[len(ts)-1] += string(ch)
		}
	}
	return
}

func ParseDefn(tokens []string, state *State) {
	fnName := tokens[2]
	paramStart, paramEnd := -1, -1
	for i, tok := range tokens {
		if tok == "[" {
			paramStart = i
		} else if tok == "]" {
			paramEnd = i
		}
	}
	if paramStart < 0 || paramEnd < 0 {
		state.Error = "Missing params"
	}
	state.Fns[fnName] = InitFunction()
}

func Compile(src string, state *State) {
	tokens := getTokens(src)
	if tokens[len(tokens)-1] != ")" {
		state.Error = "Unbalanced parens"
	}
	if tokens[1] == "defn" {
		ParseDefn(tokens, state)
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

func TestErrorIfNoParams(t *testing.T) {
	state := InitState()
	Compile("(defn one )", &state)
	if state.Error != "Missing params" {
		t.Errorf("Should have identified missing params.")
	}
}

func TestGetsParensAsTokens(t *testing.T) {
	tokens := getTokens("()")
	if len(tokens) < 2 || tokens[0] != "(" || tokens[1] != ")" {
		t.Errorf(fmt.Sprintf("Tokenization failed simple parens. Expected '(' and ')', but got %v", tokens))
	}
}

func TestGetsTextAsSeparateToken(t *testing.T) {
	tokens := getTokens("(this)")
	if len(tokens) < 3 || tokens[0] != "(" || tokens[2] != ")" || tokens[1] != "this" {
		t.Errorf(fmt.Sprintf("Tokenization failed simple parens. Expected '(', 'this', and ')', but got %v", tokens))
	}
}

func TestGetsSecondString(t *testing.T) {
	tokens := getTokens("(defn fnname [])")
	if len(tokens) < 5 || tokens[0] != "(" || tokens[1] != "defn" || tokens[2] != "fnname" || tokens[3] != "[" || tokens[4] != "]" || tokens[5] != ")" {
		t.Errorf(fmt.Sprintf("Expected '(', 'defn', 'fnname', '[', ']', ')', but got %v", tokens))
	}
}

func TestGetsBracketsAsSeparateTokens(t *testing.T) {
	tokens := getTokens("(this [])")
	if len(tokens) < 5 || tokens[0] != "(" || tokens[1] != "this" || tokens[2] != "[" || tokens[3] != "]" || tokens[4] != ")" {
		t.Errorf(fmt.Sprintf("Expected '(', 'this', '[', ']', ')', but got %v", tokens))
	}
}

/*
func TestCanDetermineFnParamNames(t *testing.T) {
	state := InitState()
	Compile("(defn f [a b c])", &state)
	fn := state.Fns["f"]
	if len(fn.Params) != 3 || ("a" != fn.Params[0] && "b" != fn.Params[1] && "c" != fn.Params[2]) {
		t.Errorf(fmt.Sprintf("Failed to determine parameters. Instead of [a b c], got %v\n", fn.Params))
	}
}
*/
