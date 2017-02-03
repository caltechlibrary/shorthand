//
// Package shorthand provides shorthand definition and expansion.
//
// shorthand.go - A simple definition and expansion notation to use
// as shorthand when a template language is too much.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package shorthand

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	// 3rd Party packages
	"github.com/russross/blackfriday"
)

// Version nummber of library and utility
const (
	Version = "v0.0.9"

	LicenseText = `
%s %s

Copyright (c) 2016, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`
)

// HowItWorks is a help text describing shorthand.
var HowItWorks = `

ASSIGNMENTS AND EXPANSIONS

Shorthand is a simple label expansion utility. It is based on a simple key value substitution.  It supports this following types of definitions

+ Assign a string to a LABEL
+ Assign the contents of a file to a LABEL
+ Assign the output of a Bash shell expression to a LABEL
+ Assign the output of a shorthand expansion to a LABEL
+ Read a file of shorthand assignments and assign any expansions to the LABEL
+ Output a LABEL value to a file
+ Output all LABEL values to a file
+ Output a LABEL assignment statement to a file
+ Output all assignment statements to a file

*shorthand* replaces the LABEL with the value assigned to it whereever it is encountered in the text being passed. The assignment statement is 
not written to stdout output.

operator                    | meaning                                  | example
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :label:                    | Assign String                            | {{name}} :label: Freda
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-text:              | Assign the contents of a file            | {{content}} :import-text: myfile.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-shorthand:         | Get assignments from a file              | _ :import-shorthand: myfile.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand:                   | Assign an expansion                      | {{reportTitle}} :expand: Report: @title for @date
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-expansion:         | Assign expanded expansion                | {{reportHeading}} :expand-expansion: @reportTitle
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import:                   | Include a file, procesisng the shorthand | {{nav}} :import: mynav.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :bash:                     | Assign Shell output                      | {{date}} :bash: date +%Y-%m-%d
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-and-bash:          | Assign Expand then gete Shell output     | {{entry}} :expand-and-bash: cat header.txt @filename footer.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :markdown:                 | Assign Markdown processed text           | {{div}} :markdown: # My h1 for a Div
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-markdown:          | Assign Expanded Markdown                 | {{div}} :expand-markdown: Greetings **@name**
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-markdown:          | Include Markdown processed text          | {{nav}} :import-markdown: mynav.md
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-expanded-markdown: | Include Expanded Markdown processed text | {{nav}} :import-expanded-markdown: mynav.md
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export:                   | Output a label's value to a file         | {{content}} :export: content.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-all:               | Output all assigned Expansions           | _ :export-all: contents.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-label:             | Output Assignment                        | {{content}} :export-label: content.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-all-labels:        | Output all Assignments                   | _ :export-all-labels: contents.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :exit:                     | Exit the shorthand repl                  | :exit:
----------------------------|------------------------------------------|---------------------------------------------------------------------



Notes: Using an underscore as a LABEL means the label will be ignored. There are no guarantees of order when writing values or assignment 
statements to a file.

The spaces surrounding " :label: ", " :import-text: ", " :bash: ", " :expand: ", " :export: ", etc. are required.


EXAMPLE

In this example a file containing the text of pre-amble is assigned to the label @PREAMBLE, the time 3:30 is assigned to the label {{NOW}}.

    {{PREAMBLE}} :import-text: /home/me/preamble.text
    {{NOW}} :label: 3:30

    At {{NOW}} I will be reading the {{PREAMBLE}} until everyone falls asleep.


If the file preamble.txt contained the phrase "Hello World" (including the quotes but without any carriage return or line feed) the output 
after processing the shorthand would look like - 

    At 3:30 I will be reading the "Hello World" until everyone falls asleep.

Notice the lines containing the assignments are not included in the output and that no carriage returns or line feeds are added the the 
substituted labels.

+ Assign shorthand expansions to a LABEL
    + LABEL :expand: SHORTHAND_TO_BE_EXPANDED
    + @content@ :expand: @report_name@ @report_date@
        + this would concatenate report name and date


PROCESSING MARKDOWN PAGES

_shorthand_ also provides a markdown processor. It uses the [blackfriday](https://github.com/russross/blackfriday) markdown library. 
This is both a convience and also allows you to treat markdown with shorthand assignments as a template that renders HTML or HTML with 
shorthand ready for expansion. It is a poorman's text rendering engine.

In this example we'll build a HTML page with shorthand labels from markdown text. Then
we will use the render HTML as a template for a blog page entry.

Our markdown file serving as a template will be call "post-template.md". It should contain
the outline of the structure of the page plus some shorthand labels we'll expand later.


    # @blogTitle

    ## @pageTitle

    ### @dateString

    @contentBlocks


For the purposes of this exercise we'll use _shorthand_ as a repl and just enter the assignments sequencly.  Also rather than use 
the output of shorthand directly we'll build up the content for the page in a label and use shorthand itself to write the final page out.

The steps we'll follow will be to 

1. Read in our markdown file page.md and turn it into an HTML with embedded shorthand labels
2. Assign some values to the labels
3. Expand the labels in the HTML and assign to a new label
4. Write the new label out to are page call "page.html"

Start the repl with this version of the shorthand command:

    shorthand -p "? "

The _-p_ option tells _shorthand_ to use the value "? " as the prompt. When _shorthand_ starts it will display "? " to indicate it is 
ready for an assignment or expansion.

The following assumes you are in the _shorthand_ repl.

Load the mardkown file and transform it into HTML with embedded shorthand labels

    @doctype :bash: echo "<!DOCTYPE html>"
    @headBlock :label: <head><title>@pageTitle</title>
    @pageTemplate :import-markdown: post-template.md
    @dateString :bash: date
    @blogTitle :label:  My Blog
    @pageTitle :label A Post
    @contentBlock :import-markdown: a-post.md
    @output :expand-expansion: @doctype<html>@headBlock<body>@pageTemplate</body></html>
    @output :export: post.html

`

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

// GetSymbol finds the symbol entry and returns the SourceMap
func (st *SymbolTable) GetSymbol(sym string) SourceMap {
	i, ok := st.labels[sym]
	if ok == true {
		return st.entries[i]
	}
	return SourceMap{Label: "", Op: "", Source: "", Expanded: "", LineNo: -1}
}

// GetSymbols returns a list of all symbols defined by labels as an array of SourceMaps
func (st *SymbolTable) GetSymbols() []SourceMap {
	var symbols []SourceMap

	for _, i := range st.labels {
		symbols = append(symbols, st.entries[i])
	}
	return symbols
}

// SetSymbol adds a SourceMap to entries and points the labels at the most recent definition.
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

// VirtualMachine defines the structure which holds symbols, operator map,
// ops and current prompt setting for a shorthand instance.
type VirtualMachine struct {
	prompt    string
	Symbols   *SymbolTable
	Operators OperatorMap
	Ops       []string
	Help      map[string]string
}

// New returns a VirtualMachine struct and registers all Operators
func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.Symbols = new(SymbolTable)
	vm.Operators = make(OperatorMap)
	vm.Help = make(map[string]string)

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

	// Register the built-in operators (glyph versions, these should be depreciated at somepoint)
	vm.RegisterOp(" :=: ", AssignString, "Assign a string to label")
	vm.RegisterOp(" :=<: ", AssignInclude, "Include content and assign to label")
	vm.RegisterOp(" :}<: ", ImportAssignments, "Import assignments from a shorthand file")

	vm.RegisterOp(" :{: ", AssignExpansion, "Expand and assign to label")
	vm.RegisterOp(" :{{: ", AssignExpandExpansion, "Expand and expansion and assign to label")
	vm.RegisterOp(" :{<: ", IncludeExpansion, "Include a file, expand and assign to label")

	vm.RegisterOp(" :!: ", AssignShell, "Assign the output of a Bash command to label")
	vm.RegisterOp(" :{!: ", AssignExpandShell, "Expand and then assign the results of a Bash command to label")

	vm.RegisterOp(" :[: ", AssignMarkdown, "Convert markdown and assign to label")
	vm.RegisterOp(" :{[: ", AssignExpandMarkdown, "Expand and convert markdown and assign to label")
	vm.RegisterOp(" :[<: ", IncludeMarkdown, "Include and convert markdown and assign to label")
	vm.RegisterOp(" :{[<: ", IncludeExpandMarkdown, "Include an expansion, convert with Markdown and assign to label")

	vm.RegisterOp(" :>: ", OutputExpansion, "Write an expansion to a file")
	vm.RegisterOp(" :@>: ", OutputExpansions, "Write all expansions to a file (order not guaranteed)")

	vm.RegisterOp(" :}>: ", ExportAssignment, "Export assignment to a file")
	vm.RegisterOp(" :@}>: ", ExportAssignments, "Expand all assignments (order not guaranteed)")

	// Register the built-in operators (readable versions)
	vm.RegisterOp(" :label: ", AssignString, "Assign a string to label")
	vm.RegisterOp(" :import-text: ", AssignInclude, "Include content and assign to label")
	vm.RegisterOp(" :import-shorthand: ", ImportAssignments, "Import assignments from a shorthand file")

	vm.RegisterOp(" :expand: ", AssignExpansion, "Expand and assign to label")
	vm.RegisterOp(" :expand-expansion: ", AssignExpandExpansion, "Expand and expansion and assign to label")
	vm.RegisterOp(" :import: ", IncludeExpansion, "Include a file, evaluate the shorthand")

	vm.RegisterOp(" :bash: ", AssignShell, "Assign the output of a Bash command to label")
	vm.RegisterOp(" :expand-and-bash: ", AssignExpandShell, "Expand and then assign the results of a Bash command to label")

	vm.RegisterOp(" :markdown: ", AssignMarkdown, "Convert markdown and assign to label")
	vm.RegisterOp(" :expand-markdown: ", AssignExpandMarkdown, "Expand and convert markdown and assign to label")
	vm.RegisterOp(" :import-markdown: ", IncludeMarkdown, "Include and convert markdown and assign to label")
	vm.RegisterOp(" :import-expand-markdown: ", IncludeExpandMarkdown, "Include an expansion, convert with Markdown and assign to label")

	vm.RegisterOp(" :export: ", OutputExpansion, "Write an the contents of an label to a file")
	vm.RegisterOp(" :export-all: ", OutputExpansions, "Write all label contents to a (order not guaranteed)")

	vm.RegisterOp(" :export-shorthand: ", ExportAssignment, "Export assignment to a file")
	vm.RegisterOp(" :export-all-shorthand: ", ExportAssignments, "Expand all assignments (order not guaranteed)")

	return vm
}

// SetPrompt sets the string value of the prompt for a VirtualMachine instance
func (vm *VirtualMachine) SetPrompt(s string) {
	vm.prompt = s
}

// RegisterOp associate a operation and function
func (vm *VirtualMachine) RegisterOp(op string, callback func(*VirtualMachine, SourceMap) (SourceMap, error), help string) error {
	_, ok := vm.Operators[op]
	if ok == true {
		return fmt.Errorf("Cannot redefine function %s\n", op)
	}
	vm.Operators[op] = callback
	vm.Ops = append(vm.Ops, op)
	vm.Help[op] = help
	return nil
}

// Parse a string, return a source map. Takes advantage of the internal ops list.
// If no valid op is found then return a source map with Label and Op set to an empty string
// while Source is set the the string that was parsed.  Expanded should always be an empty string
// at the parse stage.
func (vm *VirtualMachine) Parse(s string, lineNo int) SourceMap {
	for _, op := range vm.Ops {
		if strings.Contains(s, op) {
			parts := strings.SplitN(strings.TrimSpace(s), op, 2)
			if len(parts) == 2 {
				return SourceMap{Label: parts[0], Op: op, Source: parts[1], LineNo: lineNo, Expanded: ""}
			}
			if len(parts) == 1 {
				return SourceMap{Label: parts[0], Op: op, Source: "", LineNo: lineNo, Expanded: ""}
			}
			return SourceMap{Label: "", Op: op, Source: "", LineNo: lineNo, Expanded: ""}
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
func (vm *VirtualMachine) Run(in *bufio.Reader, postProcessWithMarkdown bool) int {
	lineNo := 0
	for {
		if vm.prompt != "" {
			fmt.Fprint(os.Stdout, vm.prompt)
		}
		src, rErr := in.ReadString('\n')
		if rErr != nil {
			break
		}
		lineNo++
		if strings.Contains(src, ":exit:") || strings.Contains(src, ":quit:") {
			break
		}
		out, err := vm.Eval(src, lineNo)
		if postProcessWithMarkdown == true {
			out = string(blackfriday.MarkdownCommon([]byte(out)))
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR (%d): %s\n", lineNo, err)
		}
		if out != "" {
			fmt.Fprint(os.Stdout, out)
		}
	}
	return lineNo
}
