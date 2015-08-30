shorthand
=========

A simple label expander. Example use cases

+ label or abbreviation expansion in Markdown files
+ label expansion in comma or tab delimited files
+ Very lightweight ad-hoc template engine for generating HTML or LaTeX

The supported command line options can be listed using the _--help_
options.

```shell
    shorthand --help
```

## Tutorial

### Timestamp in Markdown file

If the content of the markdown file _testdata/report.md_ was

```markdown

    Report Date: @now

    # Topic: The current local time.

    This report highlights the current local time of rendering this document

    The current local time is @now

```

From the command line you can do something like this

```shell
    shorthand -e "@now :! date" <  testdata/report.md | markdown > testdata/report.html
```

What this command does is launch the _shorthand_ interpreter and it
replaces all occurrences of "@now" in the markdown document with the
output from the Unix program _date_. 

The output would look something like

```html
    <p>Report Date: Sat Aug 29 11:25:48 PDT 2015</p>

    <h1>Topic: The current local time.</h1>

    <p>This report highlights the current local time of rendering this document</p>

    <p>The current local time is Sat Aug 29 11:25:48 PDT 2015</p>
```

Notice that both "@now" are replace with the same date information.

### embedding shorthand definitions

You could also embed the shorthand definitions command straight in the
markdown itself. with something like

```markdown
    @now :! date

    Report Date: @now

    # Topic: The current local time.

    This report highlights the current local time of rendering this document

    The current local time is @now

```

That makes the command line a little shorter

```shell
    shorthand < testdata/report.md | markdown > testdata/report.html
```

For all types of label assignments support see [shorthand.md](shorthand.md)

