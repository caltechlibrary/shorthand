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

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	shorthand "github.com/caltechlibrary/shorthand"
)

type expressionList []string

var (
	usage = `USAGE: %s [OPTIONS] [FILES_TO_PROCESS]`

	description = `
SYNOPSIS

%s is a command line utility to expand labels based on their assigned definitions. 
The render output is the transformed text and without the shorthand definitions 
themselves. shorthand reads from standard input and writes to standard output.

ASSIGNMENTS AND EXPANSIONS

Shorthand is a simple label expansion utility. It is based on a simple key 
value substitution.  It supports this following types of definitions

+ Assign a string to a LABEL
+ Assign the contents of a file to a LABEL
+ Assign the output of a Bash shell expression to a LABEL
+ Assign the output of a shorthand expansion to a LABEL
+ Read a file of shorthand assignments and assign any expansions to the LABEL
+ Output a LABEL value to a file
+ Output all LABEL values to a file
+ Output a LABEL assignment statement to a file
+ Output all assignment statements to a file

shorthand replaces the LABEL with the value assigned to it whereever it is 
encountered in the text being passed. The assignment statement is 
not written to stdout output.

| operator                  | meaning                                  | example                                                          |
|---------------------------|------------------------------------------|------------------------------------------------------------------|
|:label:                    | Assign String                            | {{name}} :label: Freda                                           |
|:import-text:              | Assign the contents of a file            | {{content}} :import-text: myfile.txt                             |
|:import-shorthand:         | Get assignments from a file              | _ :import-shorthand: myfile.shorthand                            |
|:expand:                   | Assign an expansion                      | {{reportTitle}} :expand: Report: @title for @date                |
|:expand-expansion:         | Assign expanded expansion                | {{reportHeading}} :expand-expansion: @reportTitle                |
|:import:                   | Include a file, procesisng the shorthand | {{nav}} :import: mynav.shorthand                                 |
|:bash:                     | Assign Shell output                      | {{date}} :bash: date +%Y-%m-%d                                   |
|:expand-and-bash:          | Assign Expand then gete Shell output     | {{entry}} :expand-and-bash: cat header.txt @filename footer.txt  |
|:markdown:                 | Assign Markdown processed text           | {{div}} :markdown: # My h1 for a Div                             |
|:expand-markdown:          | Assign Expanded Markdown                 | {{div}} :expand-markdown: Greetings **@name**                    |
|:import-markdown:          | Include Markdown processed text          | {{nav}} :import-markdown: mynav.md                               |
|:import-expanded-markdown: | Include Expanded Markdown processed text | {{nav}} :import-expanded-markdown: mynav.md                      |
|:export:                   | Output a label's value to a file         | {{content}} :export: content.txt                                 |
|:export-all:               | Output all assigned Expansions           | _ :export-all: contents.txt                                      |
|:export-label:             | Output Assignment                        | {{content}} :export-label: content.shorthand                     |
|:export-all-labels:        | Output all Assignments                   | _ :export-all-labels: contents.shorthand                         |
|:exit:                     | Exit the shorthand repl                  | :exit:                                                           |

Notes: Using an underscore as a LABEL means the label will be ignored. 
There are no guarantees of order when writing values or assignment 
statements to a file.

The spaces surrounding 

   " :label: ", " :import-text: ", " :bash: ", " :expand: ", " :export: ", etc. 
   
are required.

PROCESSING MARKDOWN PAGES

shorthand also provides a markdown processor. It uses the [blackfriday](https://github.com/russross/blackfriday) 
markdown library.  This is both a convience and also allows you to treat 
markdown with shorthand assignments as a template that renders HTML or 
HTML with shorthand ready for expansion. It is a poorman's text 
rendering engine.

In this example we'll build a HTML page with shorthand labels from 
markdown text. Then we will use the render HTML as a template for a blog 
page entry.

Our markdown file serving as a template will be call "post-template.md". It 
should contain the outline of the structure of the page plus some shorthand labels 
we'll expand later.

    # @blogTitle

    ## @pageTitle

    ### @dateString

    @contentBlocks

For the purposes of this exercise we'll use _shorthand_ as a repl and just enter 
the assignments sequencly.  Also rather than use the output of shorthand directly 
we'll build up the content for the page in a label and use shorthand itself to write 
the final page out.

The steps we'll follow will be to 

1. Read in our markdown file page.md and turn it into an HTML with embedded shorthand labels
2. Assign some values to the labels
3. Expand the labels in the HTML and assign to a new label
4. Write the new label out to are page call "page.html"

Start the repl with this version of the shorthand command:

    shorthand -p "? "

The _-p_ option tells _shorthand_ to use the value "? " as the prompt. When _shorthand_ starts 
it will display "? " to indicate it is ready for an assignment or expansion.

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

	examples = `
EXAMPLE

In this example a file containing the text of pre-amble is assigned to the 
label @PREAMBLE, the time 3:30 is assigned to the label {{NOW}}.

    {{PREAMBLE}} :import-text: /home/me/preamble.text
    {{NOW}} :label: 3:30

    At {{NOW}} I will be reading the {{PREAMBLE}} until everyone falls asleep.

If the file preamble.txt contained the phrase "Hello World" (including the 
quotes but without any carriage return or line feed) the output 
after processing the shorthand would look like - 

    At 3:30 I will be reading the "Hello World" until everyone falls asleep.

Notice the lines containing the assignments are not included in the output and 
that no carriage returns or line feeds are added the the 
substituted labels.

+ Assign shorthand expansions to a LABEL
    + LABEL :expand: SHORTHAND_TO_BE_EXPANDED
    + @content@ :expand: @report_name@ @report_date@
        + this would concatenate report name and date
`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App Specific Options
	expression              expressionList
	prompt                  string
	noprompt                bool
	vm                      *shorthand.VirtualMachine
	lineNo                  int
	postProcessWithMarkdown bool
)

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
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")

	// App Specific Options
	flag.Var(&expression, "e", "The shorthand notation(s) you wish at add")
	flag.StringVar(&prompt, "p", "=> ", "Output a prompt for interactive processing")
	flag.BoolVar(&noprompt, "n", false, "Turn off the prompt for interactive processing")
	flag.BoolVar(&postProcessWithMarkdown, "m", false, "Run final output through markdown processor")
	flag.BoolVar(&postProcessWithMarkdown, "markdown", false, "Run final output through markdown processor")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	// Configuration and command line interation
	cfg := cli.New(appName, appName, fmt.Sprintf(shorthand.LicenseText, appName, shorthand.Version), shorthand.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName)
	cfg.ExampleText = fmt.Sprintf(examples, appName)

	if showHelp == true {
		fmt.Println(cfg.Usage())
		os.Exit(0)
	}

	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}

	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	vm = shorthand.New()
	vm.RegisterOp(":exit:", exitShorthand, "Exit shorthand repl")
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
