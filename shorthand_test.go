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
// shorthand_test.go - tests for short package for handling shorthand
// definition and expansion.
//
package shorthand

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	// 3rd Party packages
	"github.com/russross/blackfriday"
)

// ok is similar to assert true, calls testing.T.Errorf if expression is false
func ok(t *testing.T, expression bool, message string) bool {
	if expression == true {
		return true
	}
	t.Errorf("not OK: %s\n", message)
	return false
}

// TestParseGlyphs
func TestParseGlyphs(t *testing.T) {
	validData := map[string]SourceMap{
		"@now1 :=: $(date)": {
			Label:    "@now1",
			Op:       " :=: ",
			Source:   "$(date)",
			Expanded: "",
			LineNo:   1,
		},
		"this :=: a valid assignment": {
			Label:    "this",
			Op:       " :=: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   2,
		},
		"this; :=: a valid assignment": {
			Label:    "this;",
			Op:       " :=: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   3,
		},
		`now; :=: $(date +%H:%M);`: {
			Label:    "now;",
			Op:       " :=: ",
			Source:   `$(date +%H:%M);`,
			Expanded: "",
			LineNo:   4,
		},
		"@now2 :=: Fred\n": {
			Label:    "@now2",
			Op:       " :=: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   5,
		},
		"@file :=<: file.txt": {
			Label:    "@file",
			Op:       " :=<: ",
			Source:   "file.txt",
			Expanded: "",
			LineNo:   6,
		},
		"@now3 :!: date": {
			Label:    "@now3",
			Op:       " :!: ",
			Source:   "date",
			Expanded: "",
			LineNo:   7,
		},
		"@now4 :{: @one @two": {
			Label:    "@now4",
			Op:       " :{: ",
			Source:   "@one @two",
			Expanded: "",
			LineNo:   8,
		},
		"@now5 :}<: test.shorthand": {
			Label:    "@now5",
			Op:       " :}<: ",
			Source:   "test.shorthand",
			Expanded: "",
			LineNo:   9,
		},
		"@now6 :[: **strong words**": {
			Label:    "@now6",
			Op:       " :[: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   10,
		},
		"@now7 :[<: test.md": {
			Label:    "@now7",
			Op:       " :[<: ",
			Source:   "test.md",
			Expanded: "",
			LineNo:   11,
		},
		"@label0 :>: label0.txt": {
			Label:    "@label0",
			Op:       " :>: ",
			Source:   "label0.txt",
			Expanded: "",
			LineNo:   12,
		},
		"@label1 :@>: label1.txt": {
			Label:    "@label1",
			Op:       " :@>: ",
			Source:   "label1.txt",
			Expanded: "",
			LineNo:   13,
		},
		"@label2 :}>: label2.txt": {
			Label:    "@label2",
			Op:       " :}>: ",
			Source:   "label2.txt",
			Expanded: "",
			LineNo:   14,
		},
		"@label3 :@}>: label3.txt": {
			Label:    "@label3",
			Op:       " :@}>: ",
			Source:   "label3.txt",
			Expanded: "",
			LineNo:   15,
		},
		`@label4 :++ something`: {
			Label:    "",
			Op:       "",
			Source:   "@label4 :++ something",
			Expanded: "",
			LineNo:   16,
		},
		"This should have @label4 and other things.": {
			Label:    "",
			Op:       "",
			Source:   "This should have @label4 and other things.",
			Expanded: "",
			LineNo:   17,
		},
		"{{pageTitle}} :=: Hello World": {
			Label:    "{{pageTitle}}",
			Op:       " :=: ",
			Source:   "Hello World",
			Expanded: "",
			LineNo:   18,
		},
		"{{year}} :!: echo -n $(date +%Y)": {
			Label:    "{{year}}",
			Op:       " :!: ",
			Source:   "echo -n $(date +%Y)",
			Expanded: "",
			LineNo:   19,
		},
		"{fred} :=: Fred": {
			Label:    "{fred}",
			Op:       " :=: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   20,
		},
		"{{strong}} :[: **strong words**": {
			Label:    "{{strong}}",
			Op:       " :[: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   21,
		},
		"{one} :=: 1": {
			Label:    "{one}",
			Op:       " :=: ",
			Source:   "1",
			Expanded: "",
			LineNo:   22,
		},
		"{two} :=: 2": {
			Label:    "{two}",
			Op:       " :=: ",
			Source:   "2",
			Expanded: "",
			LineNo:   23,
		},
		"{it} :{: {one} {two}": {
			Label:    "{it}",
			Op:       " :{: ",
			Source:   "{one} {two}",
			Expanded: "",
			LineNo:   24,
		},
		"{{html}} :[<: testdata/test.md": {
			Label:    "{{html}}",
			Op:       " :[<: ",
			Source:   "testdata/test.md",
			Expanded: "",
			LineNo:   25,
		},
		"{helloWorldTxT} :=<: testdata/helloworld.txt": {
			Label:    "{helloWorldTxT}",
			Op:       " :=<: ",
			Source:   "testdata/helloworld.txt",
			Expanded: "",
			LineNo:   26,
		},
	}

	vm := New()

	i := 1
	for s, ex := range validData {
		sm := vm.Parse(s, i)
		ok(t, strings.Compare(sm.Label, ex.Label) == 0, fmt.Sprintf("%d %q Label should match glyph %q ? %q", i, s, sm.Label, ex.Label))
		ok(t, strings.Compare(sm.Op, ex.Op) == 0, fmt.Sprintf("%d %q Op should match glyph %q ? %q", i, s, sm.Op, ex.Op))
		ok(t, strings.Compare(sm.Source, ex.Source) == 0, fmt.Sprintf("%d %q Source should match glyph %q ? %q", i, s, sm.Source, ex.Source))
		ok(t, strings.Compare(sm.Expanded, ex.Expanded) == 0, fmt.Sprintf("%d %q Source should match glyph %q ? %q", i, s, sm.Source, ex.Source))
		i++
	}
}

// TestParseReadable
func TestParseReadable(t *testing.T) {
	validData := map[string]SourceMap{
		"@now1 :label: $(date)": {
			Label:    "@now1",
			Op:       " :label: ",
			Source:   "$(date)",
			Expanded: "",
			LineNo:   1,
		},
		"this :label: a valid assignment": {
			Label:    "this",
			Op:       " :label: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   2,
		},
		"this; :label: a valid assignment": {
			Label:    "this;",
			Op:       " :label: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   3,
		},
		`now; :label: $(date +%H:%M);`: {
			Label:    "now;",
			Op:       " :label: ",
			Source:   `$(date +%H:%M);`,
			Expanded: "",
			LineNo:   4,
		},
		"@now2 :label: Fred\n": {
			Label:    "@now2",
			Op:       " :label: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   5,
		},
		"@file :import-text: file.txt": {
			Label:    "@file",
			Op:       " :import-text: ",
			Source:   "file.txt",
			Expanded: "",
			LineNo:   6,
		},
		"@now3 :bash: date": {
			Label:    "@now3",
			Op:       " :bash: ",
			Source:   "date",
			Expanded: "",
			LineNo:   7,
		},
		"@now4 :expand: @one @two": {
			Label:    "@now4",
			Op:       " :expand: ",
			Source:   "@one @two",
			Expanded: "",
			LineNo:   8,
		},
		"@now5 :import: test.shorthand": {
			Label:    "@now5",
			Op:       " :import: ",
			Source:   "test.shorthand",
			Expanded: "",
			LineNo:   9,
		},
		"@now6 :markdown: **strong words**": {
			Label:    "@now6",
			Op:       " :markdown: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   10,
		},
		"@now7 :import-markdown: test.md": {
			Label:    "@now7",
			Op:       " :import-markdown: ",
			Source:   "test.md",
			Expanded: "",
			LineNo:   11,
		},
		"@label0 :export: label0.txt": {
			Label:    "@label0",
			Op:       " :export: ",
			Source:   "label0.txt",
			Expanded: "",
			LineNo:   12,
		},
		"@label1 :export-all: label1.txt": {
			Label:    "@label1",
			Op:       " :export-all: ",
			Source:   "label1.txt",
			Expanded: "",
			LineNo:   13,
		},
		"@label2 :export: label2.txt": {
			Label:    "@label2",
			Op:       " :export: ",
			Source:   "label2.txt",
			Expanded: "",
			LineNo:   14,
		},
		"@label3 :export-all: label3.txt": {
			Label:    "@label3",
			Op:       " :export-all: ",
			Source:   "label3.txt",
			Expanded: "",
			LineNo:   15,
		},
		`@label4 :++ something`: {
			Label:    "",
			Op:       "",
			Source:   "@label4 :++ something",
			Expanded: "",
			LineNo:   16,
		},
		"This should have @label4 and other things.": {
			Label:    "",
			Op:       "",
			Source:   "This should have @label4 and other things.",
			Expanded: "",
			LineNo:   17,
		},
		"{{pageTitle}} :label: Hello World": {
			Label:    "{{pageTitle}}",
			Op:       " :label: ",
			Source:   "Hello World",
			Expanded: "",
			LineNo:   18,
		},
		"{{year}} :bash: echo -n $(date +%Y)": {
			Label:    "{{year}}",
			Op:       " :bash: ",
			Source:   "echo -n $(date +%Y)",
			Expanded: "",
			LineNo:   19,
		},
		"{fred} :label: Fred": {
			Label:    "{fred}",
			Op:       " :label: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   20,
		},
		"{{strong}} :markdown: **strong words**": {
			Label:    "{{strong}}",
			Op:       " :markdown: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   21,
		},
		"{one} :label: 1": {
			Label:    "{one}",
			Op:       " :label: ",
			Source:   "1",
			Expanded: "",
			LineNo:   22,
		},
		"{two} :label: 2": {
			Label:    "{two}",
			Op:       " :label: ",
			Source:   "2",
			Expanded: "",
			LineNo:   23,
		},
		"{it} :expand: {one} {two}": {
			Label:    "{it}",
			Op:       " :expand: ",
			Source:   "{one} {two}",
			Expanded: "",
			LineNo:   24,
		},
		"{{html}} :import-markdown: testdata/test.md": {
			Label:    "{{html}}",
			Op:       " :import-markdown: ",
			Source:   "testdata/test.md",
			Expanded: "",
			LineNo:   25,
		},
		"{helloWorldTxT} :import-text: testdata/helloworld.txt": {
			Label:    "{helloWorldTxT}",
			Op:       " :import-text: ",
			Source:   "testdata/helloworld.txt",
			Expanded: "",
			LineNo:   26,
		},
	}

	vm := New()

	i := 1
	for s, ex := range validData {
		sm := vm.Parse(s, i)
		ok(t, strings.Compare(sm.Label, ex.Label) == 0, fmt.Sprintf("%d %q Label should match %q ? %q", i, s, sm.Label, ex.Label))
		ok(t, strings.Compare(sm.Op, ex.Op) == 0, fmt.Sprintf("%d %q Op should match %q ? %q", i, s, sm.Op, ex.Op))
		ok(t, strings.Compare(sm.Source, ex.Source) == 0, fmt.Sprintf("%d %q Source should match %q ? %q", i, s, sm.Source, ex.Source))
		ok(t, strings.Compare(sm.Expanded, ex.Expanded) == 0, fmt.Sprintf("%d %q Source should match %q ? %q", i, s, sm.Source, ex.Source))
		i++
	}
}

// TestEvalGlyph
func TestEvalGlyph(t *testing.T) {
	testData := []string{
		"@now :=: $(date)",                        // 0
		"this :=: a valid assignment",             // 1
		"this; :=: is a valid assignment",         // 2
		"now; :=: $(date +\"%H:%M\");",            // 3
		"@here :=<: testdata/testme.md",           // 4
		"@there :{<: testdata/testme.md",          // 5
		"{here} :=<: testdata/testme.md",          // 6
		"{{here}} :=<: testdata/testme.md",        // 7
		"This is not an assignment",               // 8
		"this:=: is not a valid assignment",       // 9
		"nor :=:is this a valid assignment",       // 10
		"and not : =:is this a valid assignment",  // 11
		"also not := :is this a valid assignment", // 12
	}
	vm := New()
	for i, src := range testData {
		eSM := vm.Parse(src, i)
		s, err := vm.Eval(src, i)
		sm := vm.Symbols.GetSymbol(eSM.Label)
		ok(t, err == nil, fmt.Sprintf("err should be nill: %s", err))
		if eSM.Label == "" && eSM.Op == "" {
			ok(t, s != "", "glyph not an assignment so should be non-empty: "+s)
			ok(t, s == vm.Expand(s), "glyph expected '"+vm.Expand(s)+"' got "+s)
		} else {
			ok(t, eSM.LineNo == sm.LineNo, fmt.Sprintf("glyph expected line no. %d got %d", eSM.LineNo, sm.LineNo))
			ok(t, eSM.Source == sm.Source, "glyph expected source '"+eSM.Source+"' got "+sm.Source)
			ok(t, eSM.Label == sm.Label, "glyph expected label '"+eSM.Label+"' got "+sm.Label)
			ok(t, eSM.Op == sm.Op, "glyph expected op '"+eSM.Op+"' got "+sm.Op)
			ok(t, sm.Expanded != "", "glyph expected some expansion, got "+sm.Expanded)
		}
		i++
	}
}

// TestEvalReadable
func TestEvalReadable(t *testing.T) {
	testData := []string{
		"@now :label: $(date)",                        // 0
		"this :label: a valid assignment",             // 1
		"this; :label: is a valid assignment",         // 2
		"now; :label: $(date +\"%H:%M\");",            // 3
		"@here :import-text: testdata/testme.md",      // 4
		"@there :import: testdata/testme.md",          // 5
		"{here} :import-text: testdata/testme.md",     // 6
		"{{here}} :import-text: testdata/testme.md",   // 7
		"This is not an assignment",                   // 8
		"this:label: is not a valid assignment",       // 9
		"nor :label:is this a valid assignment",       // 10
		"and not : label:is this a valid assignment",  // 11
		"also not :label :is this a valid assignment", // 12
	}
	vm := New()
	for i, src := range testData {
		eSM := vm.Parse(src, i)
		s, err := vm.Eval(src, i)
		sm := vm.Symbols.GetSymbol(eSM.Label)
		ok(t, err == nil, fmt.Sprintf("err should be nill: %s", err))
		if eSM.Label == "" && eSM.Op == "" {
			ok(t, s != "", "Not an assignment so should be non-empty: "+s)
			ok(t, s == vm.Expand(s), "expected '"+vm.Expand(s)+"' got "+s)
		} else {
			ok(t, eSM.LineNo == sm.LineNo, fmt.Sprintf("expected line no. %d got %d", eSM.LineNo, sm.LineNo))
			ok(t, eSM.Source == sm.Source, "expected source '"+eSM.Source+"' got "+sm.Source)
			ok(t, eSM.Label == sm.Label, "expected label '"+eSM.Label+"' got "+sm.Label)
			ok(t, eSM.Op == sm.Op, "expected op '"+eSM.Op+"' got "+sm.Op)
			ok(t, sm.Expanded != "", "expected some expansion, got "+sm.Expanded)
		}
		i++
	}
}
func TestSymbolTable(t *testing.T) {
	vm := New()
	st := new(SymbolTable)
	ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok(t, len(st.labels) == 0, "st.labels should be zero too")

	sm1 := st.GetSymbol("@missing")
	ok(t, sm1.LineNo == -1, "Should fail with an empty symbol table")
	sm1 = vm.Parse("@now :=: This is now.", 1)
	i := st.SetSymbol(sm1)
	ok(t, i == 0, "Expected i to be zero as first element in symbol table")
	sm2 := st.GetSymbol("@now")
	ok(t, sm1.Label == sm2.Label, "expected label '"+sm1.Label+"' got "+sm2.Label)
	ok(t, sm1.Op == sm2.Op, "expected op '"+sm1.Op+"' got "+sm2.Op)
	ok(t, sm1.Source == sm2.Source, "expected source '"+sm1.Source+"' got "+sm2.Source)
	ok(t, sm1.Expanded == sm2.Expanded, "expected expanded '"+sm1.Expanded+"' got "+sm2.Expanded)
	ok(t, sm1.LineNo == sm2.LineNo, "expected expanded '"+sm1.Expanded+"' got "+sm2.Expanded)

	vm.Eval("@now :=: This is now.", 1)
	resultText := vm.Expand("This is '@now'")
	ok(t, resultText == "This is 'This is now.'", "Should have an expansion. ["+resultText+"]")
}

// Test Expand
func TestExpand(t *testing.T) {
	vm := New()
	ok(t, len(vm.Symbols.entries) == 0, "vm.Symbols.entries should be zero")
	ok(t, len(vm.Symbols.labels) == 0, "vm.Symbols.labels should be zero too")

	text := `
		   @me

		   This is some line that should not change.

		   8:00 - @now; some stuff

		   This "now" should not change. This "me" should not change.`

	expected := `
		   Fred

		   This is some line that should not change.

		   8:00 - 9:00; some stuff

		   This "now" should not change. This "me" should not change.`

	vm.Eval("@me :=: Fred", 1)
	vm.Eval("@now :=: 9:00", 2)
	result := vm.Expand(text)
	if result != expected {
		t.Fatalf("Expected:\n\n" + expected + "\n\nReceived:\n\n" + result)
	}
}

// Test include file
func TestInclude(t *testing.T) {
	vm := New()
	ok(t, len(vm.Symbols.entries) == 0, "vm.Symbols.entries should be zero")
	ok(t, len(vm.Symbols.labels) == 0, "vm.Symbols.labels should be zero too")

	buf, err := ioutil.ReadFile("testdata/testme.md")
	ok(t, err == nil, fmt.Sprintf("Should be able to read testdata/testme.md: %s", err))

	text := string(buf)
	_, err = vm.Eval("@TESTME :=<: testdata/testme.md", 1)
	ok(t, err == nil, "Should not get error of Eval assignment")
	resultText, err := vm.Eval("@TESTME", 1)
	ok(t, err == nil, "Should not get error on eval expand")
	ok(t, strings.Compare(text, resultText) == 0, "Should get same text for @TESTME")

	l := len(text)
	ok(t, len(resultText) >= l, fmt.Sprintf("Should %d have got %d results: %s", l, len(resultText), resultText))
	ok(t, strings.Contains(resultText, "A nimble webserver"), fmt.Sprintf("Should have 'A nimble webserver' in %s", resultText))
	ok(t, strings.Contains(resultText, "JSON"), fmt.Sprintf("Should have 'JSON' in %s", resultText))
}

func TestShellAssignment(t *testing.T) {
	vm := New()
	s, err := vm.Eval("@ECHO :!: echo -n 'Hello World!'", 1)
	ok(t, err == nil, fmt.Sprintf("assignment should not have an error: %s", err))
	ok(t, s == "", "Assignment should yield an empty string "+s)

	s, err = vm.Eval("@ECHO", 2)
	ok(t, err == nil, fmt.Sprintf("Expansion should not have an error: %s", err))
	ok(t, s == "Hello World!", "Should have @ECHO assignment: "+s)
}

func TestExpandedAssignment(t *testing.T) {
	vm := New()

	dateFormat := "2006-01-02"
	now := time.Now()
	// Date will generate a LF so the text will also contain it. So we'll test against a Trim later.
	vm.Eval(`@now :!: date +%Y-%m-%d`, 1)
	vm.Eval("@title :{: This is a title with date: @now", 2)
	resultText, err := vm.Eval("@title", 3)
	ok(t, err == nil, fmt.Sprintf("Expanded title should not have an error %s\n", err))
	expectedText := fmt.Sprintf("This is a title with date: %s\n", now.Format(dateFormat))
	ok(t, resultText == expectedText, "expected '"+expectedText+"' got '"+resultText+"'")

	// Now test a label that holds multiple lines that need expanding.
	text := `
			@one this is a line
			@two this is also a line
			@three this is the last line`

	vm.Eval("@one :=: 1", 1)
	vm.Eval("@two :=: 2", 2)
	vm.Eval("@three :=: 3", 3)
	vm.Eval("@text :=: "+text, 4)
	resultText = vm.Expand("@text")

	ok(t, strings.Contains(resultText, "@one this is a line"), "Should have line @one ["+resultText+"]")
	ok(t, strings.Contains(resultText, "@two this is also a line"), "Should have line @two "+resultText+"]")
	ok(t, strings.Contains(resultText, "@three this is the last line"), "Should have line @three "+resultText+"]")

	vm.Eval("@out :{{: @text", 5)
	resultText = vm.Expand("@out")

	ok(t, strings.Contains(resultText, "1 this is a line"), "Should have line 1 "+resultText)
	ok(t, strings.Contains(resultText, "2 this is also a line"), "Should have line 2 "+resultText)
	ok(t, strings.Contains(resultText, "3 this is the last line"), "Should have line 3 "+resultText)
}

func TestImport(t *testing.T) {
	vm := New()

	dateFormat := "2006-01-02"
	now := time.Now()
	atNow := fmt.Sprintf("%s\n", now.Format(dateFormat))
	atTitle := "Today's date report."
	atGreeting := "Hello World"

	// Now test evaluating a shorthand file
	_, err := vm.Eval("#T# :}<: testdata/test1.shorthand", 1)
	ok(t, err == nil, "Should be able to import testdata/test1.shorthand")

	result := vm.Expand("@now")
	ok(t, result == atNow, "Should have @now "+result)
	result = vm.Expand("@title")
	ok(t, result == atTitle, "Should have @title "+result)
	result = vm.Expand("@greeting")
	ok(t, result == atGreeting, "Should have @greeting "+result)

	shorthandText := `
		@title @now
		@greeting
		`

	result = vm.Expand(shorthandText)
	ok(t, strings.Contains(result, atTitle), "Should have title: "+atTitle)
	ok(t, strings.Contains(result, atNow), "Should have name: "+atNow)
	ok(t, strings.Contains(result, atGreeting), "Should have greeting: "+atGreeting)
}

func TestExpandingSourcesToFile(t *testing.T) {
	testFiles := map[string]string{
		"testdata/expansion1.txt": "Hello World",
		"testdata/expansion2.txt": "Hello World\nHello Max\n",
	}

	// Clean up testdata area
	for fname := range testFiles {
		if _, err := os.Stat(fname); err != nil {
			os.Remove(fname)
		}
	}

	//Generate test data and verify output
	vm := New()

	testData := []string{
		"@hello_world :=: Hello World",
		"@max :!: echo -n 'Hello Max'",
		"@hello_world :>: testdata/expansion1.txt",
		"_ :@>: testdata/expansion2.txt",
	}

	for i, data := range testData {
		_, err := vm.Eval(data, i)
		ok(t, err == nil, fmt.Sprintf("error for %s -> %s", data, err))
	}

	for fname, expectedText := range testFiles {
		buf, err := ioutil.ReadFile(fname)
		ok(t, err == nil, fmt.Sprintf("%s error: %s", fname, err))
		resultText := string(buf)
		terms := strings.Split(expectedText, "\n")
		for _, term := range terms {
			ok(t, strings.Contains(resultText, term),
				fmt.Sprintf("%s expected '%s' got '%s'", fname, term, resultText))
		}
	}
}

func TestExportAssignments(t *testing.T) {
	testFiles := map[string]string{
		"testdata/assigned1.txt": "@hello_world :=: Hello World",
		"testdata/assigned2.txt": "@hello_world :=: Hello World\n@max :!: echo -n 'Hello Max'\n",
	}
	for fname := range testFiles {
		if _, err := os.Stat(fname); err != nil {
			os.Remove(fname)
		}
	}

	vm := New()

	testData := []string{
		`@hello_world :=: Hello World`,
		`@max :!: echo -n 'Hello Max'`,
		`@hello_world :}>: testdata/assigned1.txt`,
		`_ :@}>: testdata/assigned2.txt`,
	}

	for i, src := range testData {
		_, err := vm.Eval(src, i)
		ok(t, err == nil, fmt.Sprintf("%d %s error: %s", i, src, err))
	}

	for fname, text := range testFiles {
		buf, err := ioutil.ReadFile(fname)
		ok(t, err == nil, "Should beable to read "+fname)
		resultText := string(buf)
		terms := strings.Split(text, "\n")
		for _, term := range terms {
			ok(t, strings.Contains(resultText, term),
				fmt.Sprintf("%s expected '%s' got '%s'", fname, text, resultText))
		}
	}
}

func TestMarkdownSupport(t *testing.T) {
	vm := New()

	testData := map[string]string{
		"[my link](http://example.org)": string(blackfriday.MarkdownCommon([]byte("[my link](http://example.org)"))),
		"**strong**":                    string(blackfriday.MarkdownCommon([]byte("**strong**"))),
	}

	i := 0
	for src, expected := range testData {
		vm.Eval(fmt.Sprintf("@test :[: %s", src), i)
		result := vm.Expand("@test")
		ok(t, expected == result, fmt.Sprintf("%s -> %s", expected, result))
		i++
	}

	vm.Eval("@link :=: my link", i)
	i++
	vm.Eval("@url :=: http://example.com", i)
	i++
	vm.Eval("@html :{[: [@link](@url)", i)
	i++
	expected := "<p><a href=\"http://example.com\">my link</a></p>\n"
	result := vm.Expand("@html")
	ok(t, strings.Compare(expected, result) == 0, fmt.Sprintf("%s != %s", expected, result))

	_, err := vm.Eval("@page :[<: testdata/test.md", i)
	i++
	ok(t, err == nil, fmt.Sprintf("%d testdata/test.md error: %s", i, err))
	result = vm.Expand("@page")
	ok(t, strings.Contains(result, "<h2>Another H2</h2>"), "Should have a h2 from test.md"+result)

	_, err = vm.Eval("H2 :=: heading two element", i)
	i++
	ok(t, err == nil, fmt.Sprintf("Should be able to assign string to 'H2': %s", err))
	result = vm.Expand("H2")
	ok(t, result == "heading two element", fmt.Sprintf("Should be able to expand 'H2': %s", result))

	_, err = vm.Eval("@page :{[<: testdata/test.md", i)
	i++
	ok(t, err == nil, fmt.Sprintf("%d testdata/test.md error: %s", i, err))
	result = vm.Expand("@page")
	ok(t, strings.Contains(result, "<h2>Another heading two element</h2>"), "Should have another heading two element from test.md"+result)

	// Re-read testdata/test.md and process the maarkdown
	_, err = vm.Eval("H2 :=: heading two element", i)
	i++
	ok(t, err == nil, fmt.Sprintf("Should be able to assign string to 'H2': %s", err))
	_, err = vm.Eval("@page :{[<: testdata/test.md", i)
	i++
	ok(t, err == nil, fmt.Sprintf("%d testdata/test.md error: %s", i, err))
	_, err = vm.Eval("@page :>: testdata/test.html", i)
	ok(t, err == nil, fmt.Sprintf("%d write testdata/test.html error: %s", i, err))
	i++

	_, err = os.Stat("testdata/test.html")
	ok(t, err == nil, "testdata/test.html should exist.")

	buf, err := ioutil.ReadFile("testdata/test.html")
	ok(t, err == nil, "should be able to read testdata/test.html.")
	result2 := string(buf)
	ok(t, strings.Compare(result, result2) == 0, fmt.Sprintf("Should have same results: '%s' <> '%s'", result, result2))

}

func TestRun(t *testing.T) {
	vm := New()
	if vm == nil {
		t.Error("vm was not created by New()")
	}
	fp, err := os.Open("testdata/run1.shorthand")
	ok(t, err == nil, "Should be able to open testdata/run1.shorthand")

	reader := bufio.NewReader(fp)

	fmt.Println("Starting vm.Run()")
	cnt := vm.Run(reader, false)
	ok(t, cnt == 3, fmt.Sprintf("Exited int wrong : %d", cnt))
	fmt.Println("Success.")
}
