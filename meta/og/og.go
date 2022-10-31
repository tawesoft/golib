// Package og implements a way to repressent data in the Open Graph protocol as
// typed Go structs, and render it as HTML meta tags. Typically, this is used
// to customise how a page looks when shared on Facebook or LinkedIn.
//
// This is not a parser.
//
// Note: currently only the Article and Website types are implemented.
// Generally, besides the use-case of sharing on Facebook or Linkedin, you
// should avoid relying on the format for anything complicated. It has poor
// support for internationalization, the protocol itself is archived, and the
// Open Web Foundation is moribund.
package og

import (
    "time"
)

type ObjectType string
const (
    ObjectTypeArticle           = "article"
    ObjectTypeWebsite           = "website"
)

type Object struct {
    // If your object is part of a larger website, the name which should be
    // displayed for the overall site. e.g., "IMDb".
    SiteName string

    // The title of your object as it should appear within the graph, e.g.,
    // "The Rock".
    Title string

    // A one to two sentence description of your object.
    Description string

    // The canonical URL of your object that will be used as its permanent ID in
    // the graph, e.g., "https://www.imdb.com/title/tt0117500/".
    Url string

    // The locale these tags are marked up in. Of the format language_TERRITORY.
    // Default is en_US.
    Locale string

    // An array of other locales this page is available in, depending on a
    // query parameter. Real-world support seems to be lacking. Instead, have
    // different pages for each language variant. Use hreflang meta tags
    // instead.
    // LocaleAlternates []string

    // An array of media files that represent the content.
    Media []Media

    // The type of your object, e.g., "video.movie". Depending on the type you
    // specify, other properties may also be required. If blank, defaults to
    // "website".
    Type ObjectType

    // Discriminated by Type field
    Article ObjectArticle
}

type ObjectArticle struct {
    Published time.Time // When the article was first published.
    Modified  time.Time // When the article was last changed.
    Expires   time.Time // When the article is out of date after.
    Authors   []Profile // Writers of the article. Also consider setting meta name="author".
    Section   string    // A high-level section name. E.g. Technology
    Tags []string // e.g. []string{"environment", "science"}...
}

type Profile struct {
    Url       string
    FirstName string
    LastName  string
}

type MediaType string
const (
    MediaTypeAudio = MediaType("audio")
    MediaTypeImage = MediaType("image")
    MediaTypeVideo = MediaType("video")
)

type Media struct {
    Type MediaType
    Url string // should always start with https://
    Mime string // e.g. image/jpeg, audio/mpeg, etc.

    // Alternate text (for accessibility reasons). May not work how you expect.
    // Leave blank if the parent object's title is a good description.
    Alt string

    // If Type field is MediaTypeImage or MediaTypeVideo
    Width int  // in pixels
    Height int // in pixels
}
