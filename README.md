# manuscript

## Introduction

manuscript is a simple command-line tool to convert works of fiction
from a plain-text format to one of the supported output formats.  This
allows you to compose your story in any text-editor of your choice and
store it in small files that you can easily store in a version-control
system (complete with text diffs in between versions).  The target
output (for page-formatted renderers like PDF) is the manuscript
format for [short stories](http://www.shunn.net/format/story.html) and
[novels](http://www.shunn.net/format/novel.html).

## Quick Start

You can find complete examples of valid manuscript files in the
[/example/](https://github.com/bieber/manuscript/tree/master/example)
directory of this project.  To compile the novel example to a PDF
file, you would use the following command:

```
manuscript -o output_file.pdf example/novel.man
```

## Manuscript File Format

A manuscript file has 3 parts, one of which is optional.  It begins
with basic information about the author and story, then an optional
notes section which you may use to keep notes about your story, and
finally the story itself.

### Basic Info

At the very beginning of your file, you may specify as much
information about it and yourself as you want.  At a minimum, you
should include the story's type and title and the author's byline.
Information is specified by writing a directive name at the beginning
of a line followed by the relevant information on the rest of the line
or following lines until the next directive.  This is an example of
the information section of a story:

```
@type             novel
@title            A Most Curious Happening
@shortTitle       Curious Happening
@authorByline     John Anderson
@authorName       Mortimer Willis
@authorShortName  Willis
@authorAddress    1324 W Elm Ave. #1
                  North Boynton, CT 100
@authorEmail      mortimer@willis.com
@authorOrgs       Member, PWA
                  Member, Another Writing Org
```

You may use any of the following directives:

- `@type`: This specifies the type of story.  It must be either
  "novel" or "shortStory" for novels and short stories, respectively.

- `@title`: The full title of the story, as displayed on the title
  page.

- `@shortTitle`: A shortened version of the title to use in the header
  of each page for page-formatted output such as PDF.

- `@authorByline`: The author's name as displayed on the title page.
  If you are writing under a pen name, you should put it here.

- `@authorName`: The author's full, legal name.

- `@authorShortName`: A shortened version of the author's name
  (generally your last name) to use in the header of each page for
  page-formatted output such as PDF.

- `@authorAddress`: The author's mailing address.  This directive may
  span multiple lines.

- `@authorEmail`: The author's contact email address.

- `@authorOrgs`: Professional organizations the author is a member of
  and wishes to display on the title page.

### Notes

After your information section, you may optionally include notes in
your manuscript.  This is simply a section where you can type any
information you want to keep track of that isn't part of the story
itself.  The notes section will never appear in any output.  To use
it, write the directive `@notes` at the beginning of a line after the
information section.  All text from there until the `@begin` directive
(explained in the next section) will be ignored.

### Story

After your information and/or notes, you can begin your actual story
by writing the `@begin` directive at the beginning of a line.  All
text from that point on will be read as part of the story, according
to the following formatting rules:

- `@part`: The part directive specifies the beginning of a new part in
  your story.  It should go on a line by itself, which may optionally
  include a name for the part.

- `@prologue`: The prologue directive specifies the beginning of a
  prologue.  It should go on a line by itself, which may optionally
  include a name for the prologue.

- `@chapter`: The chapter directive specifies the beginning of a
  chapter.  It should go on a line by itself, which may optionally
  include a name for the chapter.

- `@scene`: The scene marks the end of one scene and beginning of
  another.  It should go on a line by itself.

- `@note`: The note directive marks a line as a note.  Anything you
  put on the same line as the note directive will not appear in the
  output.  You can use this to leave notes for yourself within your
  story.

- Paragraphs: All text on contiguous lines is combined into the same
  paragraph.  This means that you can break up long lines of text into
  as many shorter lines as you wish in your text editor, as long as
  all the lines come one after the another.  To start a new paragraph,
  simply leave a line empty.

- Text Styles: You can bold or italicize text by putting it in between
  asterisks.  One asterisk for italic, two asterisks for bold, three
  for bold italic.  For example `*word*` would render "word"
  italicized in the output.

- Escaping: If you need to include an asterisk in the text of your
  story that you're not using for formatting, put a backslash in front
  of it.  You can also put a backslash in front of the `@` symbol to
  include the actual text of a directive in your story.

## The `manuscript` Executable

Once you've written your story, you can use the `manuscript` program
to transform the plain text into your desired output format.  Its
syntax is as follows:

```
manuscript [options] input_file
```

Where `options` is a set of command-line options, and `input_file` is
the path to the input file you want to use.

### Command-line Options

You can use the following command-line options:

- `-h`/`--help`: Display the program's usage text.

- `-o`/`--output`: Specify the file to write the output to.  This
  option is required.

- `-r`/`--renderer`: Sets the renderer to format your story with.  The
  default is pdf, but the following section will explain the renderer
  options in more detail.

### Renderers

You may select a renderer for your story by providing a `-r` or
`--renderer` argument.  The value you provide must be one of the
available renderers.  You may also include options for the renderer by
writing them in between parentheses after the renderer name,
separating option names from values with the equal sign and putting a
comma in between options.  For instance, to use the HTML renderer with
both author information and a table of contents turned on, you might
use a command-line like the following:

```
manuscript -o my_output.html -r 'html(authorInfo=true, includeTOC=true)' my_input
```

The available renderers are as follows:

- `pdf`: This is the default renderer, which writes your story out to a
  PDF file in manuscript format.  It accepts the following options:

  -`pageSize`: Sets the page size of the PDF file.  It defaults to
   `Letter`, other valid options are `A3`, `A4`, `A5`, and `Legal`.

  - `pageOrientation`: Sets the orientation of the page.  Must be
    either `P` or `Portrait` for portrait orientation, or `L` or
    `Landscape` for landscape orientation.  Defaults to portrait.

- `html`: Renders your story to an HTML file.  It accepts the
  following options:

  - `styleSheet`: Allows you to set a path in the HTML file to a
    custom style sheet.  If set, the default style will not be written
    to the output file and your custom style sheet will be used
    instead.


  - `authorInfo`: Set this to `true` or `yes` to include author info,
    which is normally excluded from HTML output.

  - `includeTOC`: Set this to `true` or `yes` to include a table of
    contents in the HTML output.

## Installation

If you have the Go language set up on your computer, you can simply
run

```
go get github.com/bieber/manuscript
```

to install the `manuscript` executable.  Binaries are also available
for 32-bit and 64-bit Linux, OSX and Windows at the
[releases](https://github.com/bieber/manuscript/releases) page.

## Questions

I won't try to pretend these are frequently asked, but they're
questions I imagine people might ask.

> I found a bug, will you fix it?

I'll certainly try to.  File an issue and I'll have a look, but I
can't guarantee any response time.  Pull requests and patches with
bugfixes are of course even better than bug reports.

> I have a feature request, will you build it?

Probably not, unless it's a very small amount of work and I really
want to see it happen.  In particular, it's very unlikely that I'll
add another renderer now that PDF and HTML output are both complete.
Feel free to submit patches or pull requests with feature additions
yourself though.


> Do you accept pull requests and patches?

But of course.  If your PR/patch is for a new feature, you might want
to check with me ahead of time to make sure I'm amenable to your
proposed feature.  I'll try to review requests within a couple days to
a couple weeks after receiving them, but I can't make any guarantees
on response time.  I ask that you follow the following standards when
submitting code:

- Make sure that `gofmt` and `golint` are both clean.  You don't have
  to religiously follow the linter if you have a good reason for
  ignoring it, but if feasible I would really like to keep the
  project's `golint` output empty.

- Please split long lines to respect the 80 column boundary with a tab
  width of 4 columns.  I know people will debate this endlessly, but
  for my own projects I stick to 80 because it fits my preferred
  editing environment well.

> Can I use this for non-fiction/will you add non-fiction support?

You really shouldn't and I definitely won't add any features
specifically for non-fiction.  For works of non-fiction LaTeX already
works excellently and you should just use that.
