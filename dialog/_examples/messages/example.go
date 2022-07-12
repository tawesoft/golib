package main

// On windows, built in the parent directory with the command
//
//     `./windows.sh messages`
//
// Needs rsrc provided by https://github.com/akavel/rsrc

import (
    "fmt"
    "image/color"
    "time"

    "github.com/tawesoft/golib/v2/dialog"
    "github.com/tawesoft/golib/v2/ks"
)

func main() {
    supported := ks.Must(dialog.Supported())
    fmt.Printf("Supported features: %+v\n", supported)

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

    dialog.Info("Now please pick any file (I won't do anything with it)")

    name, ok, err := dialog.FilePicker{
        FileTypes: [][2]string{
            {"Text Document", "*.txt *.rtf"},
            {"Image",         "*.png *.jpg *.jpeg *.bmp *.gif"},

            // "All Files", "*.*" appears by default, but you can suppress this
            // and add your own as long as the last item has a filter of
            // exactly "*.*".
            {"Pob Ffeil",      "*.*"}, // e.g. Welsh "All Files".
        },
        DefaultFileType:   0, // default to first one
    }.Open()

    if err != nil {
        dialog.Error("Got an error: %v", err)
    } else  if ok {
        dialog.Info("You selected: %s", name)
    } else {
        dialog.Warning("You didn't pick anything. But that's okay!")
    }

    if supported.DatePicker {
        t, ok, err := dialog.DatePicker{
            Title:     "",
            LongTitle: "Pick your favourite date in the year 2000:",
            Initial:   time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
            Location:  nil,
        }.Pick()

        if err != nil {
            dialog.Error("Got an error: %v", err)
        } else  if ok {
            if t.Year() == 2000 {
                dialog.Info("That %s was my favourite date, too!", t.Weekday())
            } else {
                dialog.Error("I said in the year 2000, not the year %d!", t.Year())
            }
        } else {
            dialog.Warning("You didn't pick anything. But that's okay!")
        }
    } else {
        dialog.Error("Date picker isn't supported for your machine, sorry.")
    }

    if supported.ColorPicker {
        c, ok, err := dialog.ColorPicker{
            Title:     "Favourite Colour?",
            Initial:   color.RGBA{
                R: 182,
                G: 51,
                B: 85,
            }, // "Tawesoft Red", a nice wine colour
            // Palette: true, <- optional, works with Zenity
        }.Pick()

        if err != nil {
            dialog.Error("Got an error: %v", err)
        } else  if ok {
            r, _, _, _ := c.RGBA()
            if float64(r)/0xffff > 0.6 {
                dialog.Info("You like red too, huh?")
            } else {
                dialog.Info("Needs more red. %d", r)
            }
        } else {
            dialog.Warning("You didn't pick anything. But that's okay!")
        }
    } else {
        dialog.Error("Color picker isn't supported for your machine, sorry.")
    }

}
