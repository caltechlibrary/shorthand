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
// operators - assign a function with the func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) signature
// and use RegisterOp (e.g. in the New() function) to add support to Shorthand.
//
package shorthand

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	// 3rd party packages
	"github.com/russross/blackfriday"
)

//AssignString take the Source and copy to Expanded
var AssignString = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := sm.Source
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

//AssignInclude read a file using Source as filename and put the results in Expanded
var AssignInclude = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: "", LineNo: sm.LineNo},
			fmt.Errorf("Cannot read %s: %s\n", sm.Source, err)
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// ImportAssignments evaluates the file for assignment operations
var ImportAssignments = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	var output []string
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: "", LineNo: sm.LineNo}, err
	}
	lineNo := 1
	for _, src := range strings.Split(string(buf), "\n") {
		s, err := vm.Eval(src, lineNo)
		if err != nil {
			return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: "", LineNo: sm.LineNo},
				fmt.Errorf("ERROR (%s %d): %s", sm.Source, lineNo, err)
		}
		if s != "" {
			output = append(output, s)
		}
	}
	expanded := strings.Join(output, "\n")
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignExpansion expands Source and copy to Expanded
var AssignExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := vm.Expand(sm.Source)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignExpandExpansion expand an expanded Source and copy to Expanded
var AssignExpandExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	tmp := vm.Expand(sm.Source)
	expanded := vm.Expand(tmp)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// IncludeExpansion include the filename from Source, expand and copy to Expanded
var IncludeExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	expanded := vm.Expand(string(buf))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignShell pass Source to shell and copy stdout to Expanded
var AssignShell = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := exec.Command("bash", "-c", sm.Source).Output()
	if err != nil {
		return sm, err
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignExpandShell expand Source, pass to Bash and assign output to Expanded
var AssignExpandShell = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := exec.Command("bash", "-c", vm.Expand(sm.Source)).Output()
	if err != nil {
		return sm, err
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignMarkdown process Source with Blackfriday and copy
var AssignMarkdown = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := string(blackfriday.MarkdownCommon([]byte(sm.Source)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignExpandMarkdown process source, expand witi BlackFriday and copy
var AssignExpandMarkdown = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	tmp := vm.Expand(sm.Source)
	expanded := string(blackfriday.MarkdownCommon([]byte(tmp)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// IncludeMarkdown run through markdown then assign
var IncludeMarkdown = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	tmp := string(buf)
	expanded := string(blackfriday.MarkdownCommon([]byte(tmp)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// IncludeExpandMarkdown expand then include the markdown processed source
var IncludeExpandMarkdown = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	tmp := vm.Expand(string(buf))
	expanded := string(blackfriday.MarkdownCommon([]byte(tmp)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// OutputExpansion expanded and write content to file
var OutputExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	oSM := vm.Symbols.GetSymbol(sm.Label)
	out := oSM.Expanded
	fname := sm.Source
	err := ioutil.WriteFile(fname, []byte(out), 0666)
	if err != nil {
		return sm, fmt.Errorf("%d Write error %s: %s", sm.LineNo, fname, err)
	}
	return oSM, nil
}

// OutputExpansions write the expanded content out
var OutputExpansions = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	fp, err := os.Create(sm.Source)
	if err != nil {
		return sm, fmt.Errorf("%d Create error %s: %s", sm.LineNo, sm.Source, err)
	}
	defer fp.Close()
	symbols := vm.Symbols.GetSymbols()
	for _, oSM := range symbols {
		fmt.Fprintln(fp, vm.Expand(oSM.Expanded))
	}
	return sm, nil
}

// ExportAssignment write the assignment to a file
var ExportAssignment = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	oSM := vm.Symbols.GetSymbol(sm.Label)
	out := fmt.Sprintf("%s%s%s", oSM.Label, oSM.Op, oSM.Source)
	fname := sm.Source
	err := ioutil.WriteFile(fname, []byte(out), 0666)
	if err != nil {
		return sm, fmt.Errorf("%d Write error %s: %s", sm.LineNo, fname, err)
	}
	return oSM, nil
}

// ExportAssignments write multiple assignments to a file
var ExportAssignments = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	fp, err := os.Create(sm.Source)
	if err != nil {
		return sm, fmt.Errorf("%d Create error %s: %s", sm.LineNo, sm.Source, err)
	}
	defer fp.Close()
	symbols := vm.Symbols.GetSymbols()
	for _, oSM := range symbols {
		fmt.Fprintf(fp, "%s%s%s\n", oSM.Label, oSM.Op, oSM.Source)
	}
	return sm, nil
}
