package main

import (
    "fmt"

    "github.com/tawesoft/golib/v2/dialog"
    "github.com/tawesoft/golib/v2/ks"
)

func main() {
    fmt.Printf("Supported features: %+v\n", ks.Must(dialog.Supported()))

    {
        name, ok, err := dialog.FilePicker{
            Title:             "",
            Path:              "",
            FileTypes:         nil,
            DefaultFileType:   0,
            AlwaysShowHidden:  false,
            AddToRecent:       false,
        }.Open()

        fmt.Printf("Got %v, %v, %v\n", name, ok, err)
    }
    return

    // For windows, enable modern styles. Does nothing on other platforms.
    osInit()

    dialog.Raise("Hello %s. Here's some Unicode: £¹²³€½¾", "world")
    dialog.Raise(`
Here's a really long string, to show that word-wrapping works.

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum id lorem id ligula molestie gravida ac nec justo. Suspendisse orci massa, rutrum non mauris at, malesuada molestie velit. Etiam vel velit magna. Sed eget urna tristique, pulvinar ligula at, aliquet dolor. Proin vestibulum nulla vel nisl dignissim, vitae interdum urna suscipit. Nullam non dui in purus eleifend congue vitae eget neque. In hac habitasse platea dictumst. Ut tincidunt urna placerat viverra congue. Sed sed pretium augue. Fusce nec molestie ante, vitae consequat leo. Vestibulum in diam sed eros pulvinar sagittis vel at erat. Maecenas sed mi turpis. Nam placerat ex risus, vitae venenatis erat convallis eu.
`)

    // right-to-left
    dialog.Message{
        Title: "Right-to-left writing system",
        // see https://en.wikipedia.org/wiki/Bidirectional_text#Marks
        Format: "قرأ Wikipedia™‎ طوال اليوم.",
    }.Raise()

    dialog.Message{
        Title:  "Did you know?",
        Format: "There are %d lights.",
        Icon:   dialog.IconWarning,
    }.WithArgs(4).Raise()

    qYesNo, _ := dialog.Message{
        Title:  "Let's test you",
        Format: "There were four lights, correct?",
    }.Ask()

    if qYesNo {
        dialog.Message{
            Format: "You are obedient. Well done!",
        }.Raise()
    } else {
        dialog.Message{
            Title:  "Wrong!!!",
            Format: "You are insubordinate!!! (angry face)",
            Icon:   dialog.IconError,
        }.Raise()
    }

    dialog.Alert("Now please pick any file (I won't do anything with it)")

    name, ok, err := dialog.FilePicker{
        Title:             "",
        Path:              "",
        FileTypes:         nil,
        DefaultFileType:   0,
        AlwaysShowHidden:  false,
        AddToRecent:       false,
    }.Open()

    if err != nil {
        dialog.Alert("Got an error: %v", err)
    } else  if ok {
        dialog.Alert("You selected: %s", name)
    } else {
        dialog.Alert("You didn't pick anything. But that's okay!")
    }
}
