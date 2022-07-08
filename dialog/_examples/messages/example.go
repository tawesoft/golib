package main

import (
    "fmt"

    "github.com/tawesoft/golib/v2/dialog"
)

func main() {
    osInit()

    // dialog.Open("hello.txt")
    name, ok, err := dialog.FilePicker{
        Title:             "",
        Path:              "/home/ben/Desktop/test.txt",
        FileTypes:         nil,
        DefaultFileType:   0,
        AlwaysShowHidden:  false,
        AddToRecent:       false,
    }.OpenMultiple()

    if ok { fmt.Println(name) }
    if err != nil { fmt.Println(err.Error()) }

    return

    dialog.Alert("Hello %s.", "world")

    dialog.Message{
        Title:  "Did you know?",
        Format: "There are %d lights.",
        Icon:   dialog.IconWarning,
    }.WithArgs(4).Raise()

    qYesNo := dialog.Message{
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

}
