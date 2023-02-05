package twittercard_test

import (
    "os"

    "github.com/tawesoft/golib/v2/html/meta/twittercard"
    "github.com/tawesoft/golib/v2/must"
)

func Example_summary() {
    card := twittercard.Card{
        Site:        twittercard.Account{Username: "tawesoft"},
        Title:       "Top 10 reasons why I love my cat",
        Description: "My cat can even eat a whole watermelon.",
        Type:        twittercard.CardTypeSummaryLargeImage,
        SummaryLargeImage: twittercard.CardSummaryLargeImage{
            Image:   twittercard.Image{
                Url: "https://www.example.org/media/cat-photos/cat1.jpg",
                Alt: "A black and white cat (looking very cute) sitting on a blanket with a soft toy in mid-air.",
            },
            Creator: twittercard.Account{ID: "8574052"}, // or @golightlyb
        },
    }

    must.Result(os.Stdout.WriteString(`<!doctype html>
<html lang="en-gb">
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <title>Twitter Card example</title>
`))
    must.Check(card.Write(os.Stdout))
    must.Result(os.Stdout.WriteString(`</head><body>Test!</body></html>`))
}

func Example_player() {
    video := twittercard.Video{
        Url:    "https://www.example.org/media/cat-photos/video-player.html",
        Width:  440,
        Height: 800,
        Streams: []twittercard.Media{
            {
                Url: "https://www.example.org/media/cat-photos/video.mp4",
                Type: "video/mp4",
            },
        },
    }

    card := twittercard.Card{
        Site:        twittercard.Account{Username: "tawesoft"},
        Title:       "Top 10 reasons why I love my cat",
        Description: "My cat can even eat a whole watermelon.",
        Type:        twittercard.CardTypePlayer,
        Player:      twittercard.CardPlayer{
            Video:   video,
            Image:   twittercard.Image{
                Url: "https://www.example.org/media/cat-photos/cat1.jpg",
                Alt: "A black and white cat (looking very cute) sitting on a blanket with a soft toy in mid-air.",
            },
        },
    }

    must.Result(os.Stdout.WriteString(`<!doctype html>
<html lang="en-gb">
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>Twitter Card video player example</title>
`))
    must.Check(card.Write(os.Stdout))
    must.Result(os.Stdout.WriteString(`</head><body>`))
    must.Check(video.Write(os.Stdout))
    must.Result(os.Stdout.WriteString(`</body></html>`))
}

func Example_app() {
    app := twittercard.Card{
        Site:               twittercard.Account{Username: "tawesoft"},
        Title:              "Get the Kittens Game app!",
        Description:        "My cat can even eat a whole watermelon.",
        Type:               twittercard.CardTypeApp,
        App:                twittercard.CardApp{
            Country:        "GB",
            Apps:           []twittercard.App{
                {
                    Store:  twittercard.AppStoreIPad,
                    Name:   "Kittens Game",
                    ID:     "1198099725",
                    Url:    "kittens-game://home",
                },
                {
                    Store:  twittercard.AppStoreGooglePlay,
                    Name:   "Kittens Game",
                    ID:     "com.nuclearunicorn.kittensgame",
                },
            },
        },
    }

    must.Result(os.Stdout.WriteString(`<!doctype html>
<html lang="en-gb">
    <head>
        <title>Twitter Card app example</title>
`))
    must.Check(app.Write(os.Stdout))
    must.Result(os.Stdout.WriteString(`</head><body>Test!</body></html>`))
}
