package view_test

import (
    "fmt"
    "net/url"
    "path"
    "sort"

    "github.com/acarl005/stripansi"
    "github.com/tawesoft/golib/v2/view"
)

// Untrusted applies visible runtime tainting of untrusted values (as
// recommended by OWASP). This means we are unable to use an untrusted
// input accidentally, and must first access it using an escape function.
type Untrusted struct {
    value string
}

func (u Untrusted) Raw() string {
    return u.value
}

func (u Untrusted) Escape(esc ... func (x string) string) string {
    result := u.value
    for _, escaper := range esc {
        result = escaper(result)
    }
    return result
}

func UntrustedString(x string) Untrusted {
    return Untrusted{x}
}

func Example_urlQuery() {

    // Assume an input URL query from an intrusted source
    args := url.Values{} // map[string][]string
    args.Set("name", "Ava")
    args.Add("friend", "Jess")
    args.Add("friend", "Sarah")
    args.Add("friend", "Zoe")
    args.Add("filename", "../index.html") // malicious input
    args.Add("fbclid", "nonsense we don't care about")
    // args.Encode() == ...

    recognisedKeys := []string{"name", "friend", "filename"}
    sort.Strings(recognisedKeys)
    onlyRecognised := func(k string, _ []string) bool {
        i := sort.SearchStrings(recognisedKeys, k)
        return (i < len(recognisedKeys)) && (recognisedKeys[i] == k)
    }

    // Construct a view that can read keys and values from the input query,
    // wrapping them in the Untrusted type. Additionally, we filter only
    // the keys we care about.
    //
    // Like the url.Values.Get method, returns only the first value associated
    // with the given key.
    taintedValues := view.Map[string, []string, Untrusted, Untrusted]{
        Filterer:  onlyRecognised,
        ToKey:     func (k string) Untrusted { return UntrustedString(k) },
        FromKey:   func(k Untrusted) string { return k.Raw() },
        ToValue:   func(x []string) Untrusted {
            if len(x) >= 1 { return UntrustedString(x[0]) }
            return Untrusted{}
        },
        FromValue: nil, // omitted as we don't need to map a value back
    }.Bind(args)

    if name, ok := taintedValues.Get(UntrustedString("name")); ok {
        fmt.Printf("Hi %s!\n", name.Escape(stripansi.Strip))
    } else {
        fmt.Printf("Hello anonymous!\n")
    }

    if friend, ok := taintedValues.Get(UntrustedString("friend")); ok {
        fmt.Printf("I see that you're friends with %s!\n", friend.Escape(stripansi.Strip))
    }

    if filename, ok := taintedValues.Get(UntrustedString("filename")); ok {
        // NOTE: this is an example only and is not complete as there are still
        // other ways the path could be unsafe.
        fmt.Printf("Safe filename: %s\n", filename.Escape(path.Clean, path.Base, stripansi.Strip))
    } else {
        fmt.Printf("No file specified.\n")
    }

    if _, ok := taintedValues.Get(UntrustedString("fbclid")); ok {
        panic("Didn't expect to see this!")
    }

    // Output:
    // Hi Ava!
    // I see that you're friends with Jess!
    // Safe filename: index.html
}
