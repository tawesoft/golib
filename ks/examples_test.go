package ks_test

import (
    "errors"
    "fmt"

    "github.com/tawesoft/golib/v2/ks"
)

func Example_Range() {
    {
        fmt.Println("String:")
        str := "Hello, world!"

        k, err := ks.Range[int, rune](func (k int, v rune) error {
            fmt.Printf("%d: %c (%x)\n", k, v, v)
            if v == ',' {
                return errors.New("oops")
            }
            return nil
        }, str)

        fmt.Printf("key at error: %d; error: %v\n", k, err)
    }

    {
        fmt.Println("Slice:")
        slice := []int{100, 200, 300, 400, 500}

        k, err := ks.Range[int, int](func (k int, v int) error {
            fmt.Printf("%d: %d\n", k, v)
            if k == 2 {
                return errors.New("oops")
            }
            return nil
        }, slice)

        fmt.Printf("key at error: %d; error: %v\n", k, err)
    }

    {
        fmt.Println("Map:")
        m := map[string]string{
            "cat": "meow",
            "cow": "moo",
            "dog": "woof",
            "lion": "roar",
        }

        k, err := ks.Range[string, string](func (k string, v string) error {
            if k == "lion" {
                return errors.New("scary")
            }
            return nil
        }, m)

        fmt.Printf("key at error: %s; error: %v\n", k, err)
    }

    // Output:
    // String:
    // 0: H (48)
    // 1: e (65)
    // 2: l (6c)
    // 3: l (6c)
    // 4: o (6f)
    // 5: , (2c)
    // key at error: 5; error: oops
    // Slice:
    // 0: 100
    // 1: 200
    // 2: 300
    // key at error: 2; error: oops
    // Map:
    // key at error: lion; error: scary
}
