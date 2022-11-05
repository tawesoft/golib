package opengraph_test

import (
    "os"
    "time"

    "github.com/tawesoft/golib/v2/meta/opengraph"
    "github.com/tawesoft/golib/v2/must"
)

func Example_article() {
    article := opengraph.Object{
        SiteName:     "My Site",
        Title:        "Top 10 reasons why I love my cat",
        Description:  "My cat can even eat a whole watermelon.",
        Url:          "https://www.example.org/articles/my-cat",
        Locale:       "en-GB",
        Media:        []opengraph.Media{
            {
                Type:   opengraph.MediaTypeImage,
                Url:    "https://www.example.org/media/cat-photos/cat1.jpg",
                Mime:   "image/jpeg",
                Width:  1024,
                Height: 768,
            },
            {
                Type:   opengraph.MediaTypeImage,
                Url:    "https://www.example.org/media/cat-photos/cat2.jpg",
                Mime:   "image/jpeg",
                Width:  1024,
                Height: 768,
            },
            {
                Type: opengraph.MediaTypeAudio,
                Url:  "https://www.example.org/media/cat-photos/purr.ogg",
                Mime: "audio/ogg",
            },
            {
                Type: opengraph.MediaTypeVideo,
                Url:  "https://www.example.org/media/cat-photos/hunting-toy.ogv",
                Mime: "video/ogg",
            },
        },
        Type: opengraph.ObjectTypeArticle,
        Article: opengraph.ObjectArticle{
            Published: time.Date(2022, 10, 31, 13, 17, 0, 0, must.Result(time.LoadLocation("Europe/London"))),
            Authors: []opengraph.Profile{
                {
                    Url:       "https://www.example.org/authors/hopperg",
                    FirstName: "Grace",
                    LastName:  "Hopper",
                },
            },
            Section: "opinion",
            Tags: []string{"cats", "pets", "cute"},
        },
    }

    must.Result(os.Stdout.WriteString(`<!doctype html>
<html lang="en-gb">
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <title>Open Graph example</title>
`))
    must.Check(opengraph.HTML(os.Stdout, article))
    must.Result(os.Stdout.WriteString(`</head><body>Test!</body></html>`))
}
