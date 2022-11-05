package opengraph

import (
    _ "embed"
    "html/template"
    "io"
    "time"

    "github.com/tawesoft/golib/v2/must"
)

//go:embed og.gohtml
var rawTemplate string

var gTemplate *template.Template

func init() {
    fm := map[string]any{}

    fm["ISO8601"] = func(dt time.Time) string {
        if dt.IsZero() { return "" }
        return dt.UTC().Format("2006-01-02T15:04Z")
    }

    gTemplate = must.Result(template.New("").Funcs(fm).Parse(rawTemplate))
}

// HTML renders an Open Graph object as HTML.
//
// If an error occurs executing the template or writing its output, execution
// stops, but partial results may already have been written to the output
// writer.
func HTML(wr io.Writer, object Object) error {
    return gTemplate.Execute(wr, object)
}
