package twittercard

import (
    _ "embed"
    "html/template"
    "io"

    "github.com/tawesoft/golib/v2/must"
)

//go:embed card.gohtml
var rawCardTemplate string

//go:embed video.gohtml
var rawVideoTemplate string

var gCardTemplate *template.Template
var gVideoTemplate *template.Template

func init() {
    gCardTemplate = must.Result(template.New("").Parse(rawCardTemplate))
    gVideoTemplate = must.Result(template.New("").Parse(rawVideoTemplate))
}

// Write renders a Twitter Card as HTML. This should be done in the <head>
// of the document.
//
// If an error occurs executing the template or writing its output, execution
// stops, but partial results may already have been written to the output
// writer.
func (c Card) Write(wr io.Writer) error {
    return gCardTemplate.Execute(wr, c)
}

// Write renders a video player frame as HTML. This should be done in the
// <body> of the document.
//
// Note that this frame does not have to appear in the same document as the
// Twitter Card. It is the frame that appears at the [Video.Url] location.
//
// If an error occurs executing the template or writing its output, execution
// stops, but partial results may already have been written to the output
// writer.
func (c Video) Write(wr io.Writer) error {
    return gVideoTemplate.Execute(wr, c)
}
