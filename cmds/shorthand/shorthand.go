//
// shorthand.go - command line utility to process shorthand definitions
// and render output with the transformed text and without any
// shorthand definitions.
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
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	// Caltech Library packages
	shorthand "github.com/caltechlibrary/shorthand"
)

type expressionList []string

var (
	help                    bool
	version                 bool
	expression              expressionList
	prompt                  string
	noprompt                bool
	vm                      *shorthand.VirtualMachine
	lineNo                  int
	postProcessWithMarkdown bool
)

var usage = func(exit_code int, msg string) {
	var fh = os.Stderr
	if exit_code == 0 {
		fh = os.Stdout
	}
	cmdName := path.Base(os.Args[0])

	fmt.Fprintf(fh, `%s
USAGE %s [options] [FILES_TO_PROCESS]

%s is a command line utility to expand labels based on their
assigned definitions. The render output is the transformed text 
and without the shorthand definitions themselves. %s reads 
from standard input and writes to standard output. 

OPTIONS

`, msg, cmdName, cmdName)

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(fh, "\t-%s\t\t%s\n", f.Name, f.Usage)
	})

	fmt.Fprintf(fh, `

HOW IT WORKS

%s

Version: %s

copyright (c) 2015 all rights reserved.
Released under the BSD 2-Clause license.
See: http://opensource.org/licenses/BSD-2-Clause
`, shorthand.HowItWorks, shorthand.Version)
	os.Exit(exit_code)
}

var welcome = `
  Welcome to shorthand the simple label expander and markdown processor.
  Use ':exit:' to quit the repl, ':help:' to get a list of supported operators.

`

var helpShorthand = func(vm *shorthand.VirtualMachine, sm shorthand.SourceMap) (shorthand.SourceMap, error) {
	fmt.Printf(`
The following operators are supported in shorthand:

`)
	for op, msg := range vm.Help {
		fmt.Printf("\t%s\t%s\n", op, msg)
	}
	fmt.Printf("\nshorthand %s\n\n", shorthand.Version)
	return shorthand.SourceMap{Label: "", Op: ":help:", Source: "", Expanded: ""}, nil
}

//exitShorthand - call os.Exit() with appropriate value and exit the repl
var exitShorthand = func(vm *shorthand.VirtualMachine, sm shorthand.SourceMap) (shorthand.SourceMap, error) {
	if sm.Source == "" {
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, sm.Source)
	os.Exit(1)
	return shorthand.SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: ""}, nil
}

func revision() {
	fmt.Printf("%s %s\n", filepath.Base(os.Args[0]), shorthand.Version)
	os.Exit(0)
}

func (e *expressionList) String() string {
	return fmt.Sprintf("%s", *e)
}

func (e *expressionList) Set(value string) error {
	lineNo++
	out, err := vm.Eval(value, lineNo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR (%d): %s\n", lineNo, err)
		return err
	}
	if out != "" {
		fmt.Fprintf(os.Stdout, "%s\n", out)
	}
	return nil
}

func init() {
	vm = shorthand.New()
	vm.RegisterOp(":exit:", exitShorthand, "Exit shorthand repl")
}

func main() {
	flag.Var(&expression, "e", "The shorthand notation(s) you wish at add")
	flag.StringVar(&prompt, "p", "=> ", "Output a prompt for interactive processing")
	flag.BoolVar(&noprompt, "n", false, "Turn off the prompt for interactive processing")
	flag.BoolVar(&help, "h", false, "Display this help document")
	flag.BoolVar(&help, "help", false, "Display this help document")
	flag.BoolVar(&version, "v", false, "Version information")
	flag.BoolVar(&version, "version", false, "Version information")
	flag.BoolVar(&postProcessWithMarkdown, "m", false, "Run final output through markdown processor")
	flag.BoolVar(&postProcessWithMarkdown, "markdown", false, "Run final output through markdown processor")
	flag.Parse()
	args := flag.Args()
	if help == true {
		usage(0, "")
	}
	if version == true {
		revision()
	}

	if noprompt == true {
		prompt = ""
	}
	vm.SetPrompt(prompt)

	// If a filename is provided on the command line use it instead of standard input.
	if len(args) > 0 {
		vm.SetPrompt("")
		for _, arg := range args {
			fp, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			defer fp.Close()
			reader := bufio.NewReader(fp)
			vm.Run(reader, postProcessWithMarkdown)
		}
	} else {
		// Run as repl
		vm.RegisterOp(":help:", helpShorthand, "This help message")
		if prompt != "" {
			fmt.Println(welcome)
		}
		reader := bufio.NewReader(os.Stdin)
		vm.Run(reader, false)
	}
}
