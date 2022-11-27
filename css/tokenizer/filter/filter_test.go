package filter_test

import (
    "strings"
    "testing"
    "unicode/utf8"

    "github.com/stretchr/testify/assert"
    "github.com/tawesoft/golib/v2/css/internal/filter"
    "golang.org/x/text/transform"
)

func TestFilter(t *testing.T) {
    type row struct {
        input string
        expected string
        err error
    }

    rows := []row{
        {"foo\r\r\n", "foo\n\n", nil},
        {"foo\r\n\r", "foo\n\n", nil},
        {"foo\000foo", "foo\uFFFDfoo", nil},
    }

    f := filter.Transformer()
    for _, r := range rows {
        actual, _, err := transform.String(f, r.input)
        assert.Equal(t, r.err, err)
        if err == r.err {
            assert.Equal(t, r.expected, actual)
        }
    }
}

func FuzzFilter(f *testing.F) {
    testcases := []string{"foo\r\r\n", "foo\ffoo\r", "foo\000"}
    for _, tc := range testcases {
        f.Add(tc)
    }

    filter := filter.Transformer()
    f.Fuzz(func(t *testing.T, orig string) {
        filtered, _, err := transform.String(filter, orig)
        if err != nil { return }

        if strings.ContainsAny(filtered, "\r\f") {
            t.Errorf("Transform failed to filter string %q", filtered)
        }

        if !utf8.ValidString(filtered) {
            t.Errorf("Transform produced invalid UTF-8 string %q", filtered)
        }
    })
}
