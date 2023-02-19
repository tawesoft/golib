package algorithmic_test

import (
    "fmt"

    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/text/number/algorithmic"
)

func ExampleFormat() {
    print := func(system string, n int64) {
        if s, err := algorithmic.Format(system, n); err == nil {
            fmt.Printf("algorithmic.Format(%q, %d): %q\n", system, n, s)
        } else {
            fmt.Printf("algorithmic.Format(%q, %d): error: %v\n", system, n, err)
        }
    }

    print("roman-upper", 2023)

    // Output:
    // algorithmic.Format("roman-upper", 2023): "MMXXIII"
}

func Example_rulesetNames() {
    const v = int64(1234)
    fmt.Printf("Supported rulesets and example(%d):\n", v)
    for _, name := range algorithmic.RulesetNames {
        str := must.Result(algorithmic.Format(name, v))
        fmt.Printf("> %s; %s\n", name, str)
    }

    // Output:
    // Supported rulesets and example(1234):
    // > armenian-lower; ռմլդ
    // > armenian-upper; ՌՄԼԴ
    // > cyrillic-lower; ҂асл҃д
    // > ethiopic; ፩፻፪፻፴፬
    // > georgian; შსლდ
    // > greek-lower; ͵ασλδ´
    // > greek-upper; ͵ΑΣΛΔ´
    // > hebrew; א׳רל״ד
    // > hebrew-item; תתתלד
    // > roman-lower; mccxxxiv
    // > roman-upper; MCCXXXIV
}
