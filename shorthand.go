//
// Package shorthand provides shorthand definition and expansion.
//
// shorthand.go - A simple definition and expansion notation to use
// as shorthand when a template language is too much.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the BSD 2-Clause license
// See: http://opensource.org/licenses/BSD-2-Clause
//
package shorthand

import (
	"bufio"
	"fmt"
	//"github.com/russross/blackfriday"
	//"io/ioutil"
	"os"
	//"os/exec"
	"strings"
)

// The version nummber of library and utility
const Version = "v0.0.4-next"

//
// An Op is built from a multi character symbol
// Each element in the symbol has meaning
// " :" is the start of a glyph band the end ": " is a colon and trailing space
// = the source value is next, this is basic assignment of a string value to a symbol
// < is input from a file
// { expand (resolved label values)
// } assign a statement (i.e. label, op, value)
// ! input from a shell expression
// [ is a markdown expansion
// > is to output an to a file
// @ operate on whole symbol table
//

// SourceMap holds the source and value of an assignment
type SourceMap struct {
	Label    string // Label is the symbol to be replace based on Op and Source
	Op       string // Op is the type of assignment being made (if is an empty string if not an assignment)
	Source   string // Source is argument to the right of Op
	Expanded string // Expanded is the value calculated based on Label, Op and Source
	LineNo   int
}

// SymbolTable holds the exressions, values and other errata of parsing assignments making expansions
type SymbolTable struct {
	entries []SourceMap
	labels  map[string]int
}

func (st *SymbolTable) GetSymbol(sym string) SourceMap {
	i, ok := st.labels[sym]
	if ok == true {
		return st.entries[i]
	}
	return SourceMap{Label: "", Op: "", Source: "", Expanded: "", LineNo: -1}
}

func (st *SymbolTable) GetSymbols() []SourceMap {
	var symbols []SourceMap

	for _, i := range st.labels {
		symbols = append(symbols, st.entries[i])
	}
	return symbols
}

func (st *SymbolTable) SetSymbol(sm SourceMap) int {
	st.entries = append(st.entries, sm)
	if st.labels == nil {
		st.labels = make(map[string]int)
	}
	i := len(st.entries) - 1
	st.labels[sm.Label] = i
	return st.labels[sm.Label]
}

// OperatorMap is a map of operator testings and their related functions
type OperatorMap map[string]func(*VirtualMachine, SourceMap) (SourceMap, error)

type VirtualMachine struct {
	Symbols   *SymbolTable
	Operators OperatorMap
	Ops       []string
}

// VirtualMachine binds methods to shared Symbol and Operator structure.
// It is responsible for registering all supported Operators
func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.Symbols = new(SymbolTable)
	vm.Operators = make(OperatorMap)
	//
	// An Op is built from a multi character symbol
	// Each element in the symbol has meaning
	// " :" is the start of a glyph band the end ": " is a colon and trailing space
	// = the source value is next, this is basic assignment of a string value to a symbol
	// < is input from a file
	// { expand (resolved label values)
	// } assign a statement (i.e. label, op, value)
	// ! input from a shell expression
	// [ is a markdown expansion
	// > is to output an to a file
	// @ operate on whole symbol table
	//
	// Now register the built-in operators
	vm.RegisterOp(" :=: ", AssignString)
	vm.RegisterOp(" :=<: ", AssignInclude)
	vm.RegisterOp(" :}<: ", ImportAssignments)
	vm.RegisterOp(" :{: ", AssignExpansion)
	vm.RegisterOp(" :{{: ", AssignExpandExpansion)
	vm.RegisterOp(" :{<: ", IncludeExpansion)
	vm.RegisterOp(" :!: ", AssignShell)
	vm.RegisterOp(" :{!: ", AssignExpandShell)
	vm.RegisterOp(" :[: ", AssignMarkdown)
	vm.RegisterOp(" :{[: ", AssignExpandMarkdown)
	vm.RegisterOp(" :[<: ", IncludeMarkdown)
	vm.RegisterOp(" :{[<: ", IncludeExpandMarkdown)
	vm.RegisterOp(" :>: ", OutputExpansion)
	vm.RegisterOp(" :@>: ", OutputExpansions)
	vm.RegisterOp(" :}>: ", ExportAssignment)
	vm.RegisterOp(" :@}>: ", ExportAssignments)
	vm.RegisterOp(" :exit: ", ExitShorthand)
	vm.RegisterOp(" :quit: ", ExitShorthand)
	return vm
}

// RegisterOp associate a operation and function
func (vm *VirtualMachine) RegisterOp(op string, callback func(*VirtualMachine, SourceMap) (SourceMap, error)) error {
	_, ok := vm.Operators[op]
	if ok == true {
		return fmt.Errorf("Cannot redefine function %s\n", op)
	}
	vm.Operators[op] = callback
	vm.Ops = append(vm.Ops, op)
	return nil
}

// Parse a string, return a source map. Takes advantage of the internal ops list.
// If no valid op is found then return a source map with Label and Op set to an empty string
// while Source is set the the string that was parsed.  Expanded should always be an empty string
// at the parse stage.
func (vm *VirtualMachine) Parse(s string, lineNo int) SourceMap {
	for _, op := range vm.Ops {
		if strings.Index(s, op) != -1 {
			parts := strings.SplitN(strings.TrimSpace(s), op, 2)
			return SourceMap{Label: parts[0], Op: op, Source: parts[1], LineNo: lineNo, Expanded: ""}
		}
	}
	return SourceMap{Label: "", Op: "", Source: s, LineNo: lineNo, Expanded: ""}
}

// Expand takes some text and expands all labels to their values
func (vm *VirtualMachine) Expand(text string) string {
	// labels hash should also point at the last known state of
	// the label
	result := text
	symbols := vm.Symbols.GetSymbols()
	for _, sm := range symbols {
		if strings.Contains(text, sm.Label) {
			tmp := strings.Replace(result, sm.Label, sm.Expanded, -1)
			result = tmp
		}
	}
	return result
}

// Eval stores a shorthand assignment or expands and writes the content to stdout
// Returns the expanded  and any error
func (vm *VirtualMachine) Eval(s string, lineNo int) (string, error) {
	sm := vm.Parse(s, lineNo)
	// If not an assignment Expand and return the expansion
	if sm.Label == "" && sm.Op == "" {
		return fmt.Sprintf("%s", vm.Expand(s)), nil
	}

	callback, ok := vm.Operators[sm.Op]
	if ok == false {
		return "", fmt.Errorf("ERROR (%d): %s is not a supported assignment.\n", lineNo, s)
	}

	// Make the associated assignment and save the symbol to the symbol table.
	newSM, err := callback(vm, sm)
	if err != nil {
		return "", err
	}

	vm.Symbols.SetSymbol(newSM)
	return "", nil
}

// Run takes a reader (e.g. os.Stdin), and two writers (e.g. os.Stdout and os.Stderr)
// It reads until EOF, :exit:, or :quit: operation is encountered
// returns the number of lines processed.
func (vm *VirtualMachine) Run(in *bufio.Reader) int {
	lineNo := 0
	for {
		src, rErr := in.ReadString('\n')
		if rErr != nil {
			break
		}
		lineNo += 1
		if strings.Contains(src, ":exit:") || strings.Contains(src, ":quit:") {
			break
		}
		out, err := vm.Eval(src, lineNo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR (%d): %s\n", lineNo, err)
		}
		if out != "" {
			fmt.Fprintf(os.Stdout, out)
		}
	}
	return lineNo
}
