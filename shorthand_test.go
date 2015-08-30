//
// Package stn is a library for processing Simple Timesheet Notation.
//
// shorthand_test.go - tests for short package for handling shorthand
// definition and expansion.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the BSD 2-Clause license
// See: http://opensource.org/licenses/BSD-2-Clause
//
package shorthand

import (
	"fmt"
	"github.com/rsdoiel/ok"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

// Test IsAssignment
func TestIsAssignment(t *testing.T) {
	validAssignments := []string{
		"@now := $(date)",
		"this := a valid assignment",
		"this; := is a valid assignment",
		"now; := $(date +\"%H:%M\");",
		"@here :< testdata/testme.md",
	}

	invalidAssignments := []string{
		"This is not an assignment",
		"this:=  is not a valid assignment",
		"nor :=is this a valid assignment",
	}

	for i := range validAssignments {
		if IsAssignment(validAssignments[i]) == false {
			t.Fatalf(validAssignments[i] + " should be a valid assignment.")
		}
	}

	for i := range invalidAssignments {
		if IsAssignment(invalidAssignments[i]) == true {
			t.Fatalf(invalidAssignments[i] + " should be an invalid assignment.")
		}
	}
}

// Test Assign
func TestAssign(t *testing.T) {
	Clear()
	validAssignments := []string{
		"@now := $(date)",
		"this := a valid assignment",
		"this; := is a valid assignment",
		"now; := $(date +\"%H:%M\");",
		"@new := Fred\n",
	}
	expectedMap := map[string]string{
		"@now":  "$(date)",
		"this":  "a valid assignment",
		"this;": "is a valid assignment",
		"now;":  "$(date +\"%H:%M\");",
		"@new":  "Fred",
	}

	for i := range validAssignments {
		if Assign(validAssignments[i]) == false {
			t.Fatalf(validAssignments[i] + " should be assigned.")
		}
	}

	for key, value := range expectedMap {
		sm, OK := Abbreviations[key]
		if !OK {
			t.Fatalf("Could not find the shorthand for " + key)
		}
		if value != sm.value {
			t.Fatalf("[" + value + "] != [" + sm.value + "]")
		}
	}
}

// Test Expand
func TestExpand(t *testing.T) {
	Clear()
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

	Assign("@me := Fred\n")
	Assign("@now := 9:00")
	result := Expand(text)
	if result != expected {
		t.Fatalf("Expected:\n\n" + expected + "\n\nReceived:\n\n" + result)
	}
}

// Test include file
func TestInclude(t *testing.T) {
	text := `
Today is @NOW.

Now add the testme.md to this.
-------------------------------
@TESTME
-------------------------------
Did it work?
`
	Assign("@NOW := 2015-07-04")
	expected := true
	results := HasAssignment("@NOW")
	ok.Ok(t, results == expected, "Should have @NOW assignment")
	Assign("@TESTME :< testdata/testme.md")
	results = HasAssignment("@TESTME")
	ok.Ok(t, results == expected, "Should have @TESTME assignment")
	resultText := Expand(text)
	l := len(text)
	ok.Ok(t, len(resultText) > l, "Should have more results: "+resultText)
	ok.Ok(t, strings.Contains(resultText, "A nimble webserver"), fmt.Sprintf("Should have 'A nimble webserver' in %s", resultText))
	ok.Ok(t, strings.Contains(resultText, "JSON"), fmt.Sprintf("Should have 'JSON' in %s", resultText))
}

func TestShellAssignment(t *testing.T) {
	expected := true
	expectedText := "Hello World!"
	Assign("@ECHO :! echo -n 'Hello World!'")
	results := HasAssignment("@ECHO")
	ok.Ok(t, results == expected, "Should have @ECHO assignment")
	resultText := Expand("@ECHO")
	l := len(strings.Trim(resultText, "\n"))
	ok.Ok(t, l == len(expectedText), fmt.Sprintf("Expected length %d got %d for @ECHO", len(expectedText), l))
	ok.Ok(t, strings.Contains(strings.Trim(resultText, "\n"), expectedText), "Should have matching text for @ECHO")
}

func TestExpandedAssignment(t *testing.T) {
	dateFormat := "2006-01-02"
	now := time.Now()
	// Date will generate a LF so the text will also contain it. So we'll test against a Trim later.
	Assign(`@now :! date +%Y-%m-%d`)
	Assign("@title :{ This is a title with date: @now")
	text := `@title`
	expected := true
	results := HasAssignment("@now")
	ok.Ok(t, results == expected, "Should have @now")
	results = HasAssignment("@title")
	ok.Ok(t, results == expected, "Should have @title")
	expectedText := fmt.Sprintf("This is a title with date: %s", now.Format(dateFormat))
	resultText := Expand(text)
	l := len(strings.Trim(resultText, "\n"))
	ok.Ok(t, l == len(expectedText), "Should have expected length for @title")
	ok.Ok(t, strings.Contains(strings.Trim(resultText, "\n"), expectedText), "Should have matching text for @title")
}

func TestExpandingValuesToFile(t *testing.T) {
	a1 := `@hello_world := Hello World`
	a2 := `@max :! echo -n 'Hello Max'`
	e1 := "Hello World"
	e2 := "Hello Max"
	Assign(a1)
	Assign(`@hello_world :> testdata/helloworld.txt`)
	b, err := ioutil.ReadFile("testdata/helloworld.txt")
	ok.Ok(t, err == nil, "Should beable to read testdata/helloworld.txt")
	resultText := string(b)
	ok.Ok(t, resultText == e1, "Shoud have Hello World from file.")
	Assign(a2)
	Assign(`@hello_world :=> testdata/helloworld.txt`)
	b, err = ioutil.ReadFile("testdata/helloworld.txt")
	ok.Ok(t, err == nil, "Should be able to read testdata/helloworld.txt")
	resultText = string(b)
	ok.Ok(t, strings.Contains(resultText, e1), "Should find "+e1)
	ok.Ok(t, strings.Contains(resultText, e2), "Should find "+e2)
}

func TestExpandingAssignmentsToFile(t *testing.T) {
	a1 := `@hello_world := Hello World`
	a2 := `@max :! echo -n 'Hello Max'`
	Assign(a1)
	Assign(`@hello_world :} testdata/assigned.txt`)
	b, err := ioutil.ReadFile("testdata/assigned.txt")
	ok.Ok(t, err == nil, "Should beable to read testdata/assigned.txt")
	resultText := string(b)
	ok.Ok(t, resultText == a1, "Shoud have @hello_world assignment in file.")
	Assign(a2)
	Assign(`@hello_world :=} testdata/assigned.txt`)
	b, err = ioutil.ReadFile("testdata/assigned.txt")
	ok.Ok(t, err == nil, "Should have all assigments in file.")
	resultText = string(b)
	ok.Ok(t, strings.Contains(resultText, a1), "Should find "+a1)
	ok.Ok(t, strings.Contains(resultText, a2), "Should find "+a2)
}
