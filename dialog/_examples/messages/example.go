package main

import (
    "github.com/tawesoft/golib/v2/dialog"
)

func main() {
    osInit()

    dialog.Open("hello.txt")

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
