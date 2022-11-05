// Package twittercard implements a way to represent Twitter Cards Markup data
// as typed Go structs, and render it as HTML meta tags, to customise how a
// page looks when shared on Twitter.
//
// This is not a parser.
package twittercard

type CardType string
const (
    CardTypeSummary = "summary"
    CardTypeSummaryLargeImage = "summary_large_image"
    CardTypePlayer = "player"
    CardTypeApp = "app"
)

// Account represents a Twitter account. Either Username or ID should be
// specified. Presumably, ID is constant even if a username may change.
type Account struct {
    Username string // e.g. "@username"
    ID string       // e.g. "1234567890"
}

type Image struct {
    // URL to the image. May have aspect ratio and/or size requirements. Must be
    // less than 5MB in file size. Only JPG, PNG, WEBP and GIF formats are
    // supported. Only the first frame of an animated GIF will be used.
    Url string

    // A text description of the image conveying the essential nature of the
    // image to users who are visually impaired. Maximum 420 characters.
    Alt string
}

type Card struct {
    Site Account // the website that the card should be attributed to
    Title string // max 70 chars
    Description string // max 200 chars
    Type CardType

    // Discriminated by Type
    Summary CardSummary
    SummaryLargeImage CardSummaryLargeImage
    Player CardPlayer
    App CardApp
}

type CardSummary struct {
    Image Image // Aspect ratio 1:1. From 144x144 to 4096x4096 pixels.
}

type CardSummaryLargeImage struct {
    Image Image // Aspect ratio 2:1. From 300x157 to 4096x4096 pixels.
    Creator Account // the content creator, optional
}

type Video struct {
    // HTTPS URL of player frame (see https://github.com/twitterdev/cards-player-samples).
    // You can generate the body of this frame with the [Video.Write] method.
    Url string

    Width int // of frame, in pixels
    Height int // of frame, in pixels

    Streams []Media // Optional. URL to raw video or audio streams.
}

type Media struct {
    Type string // MIME type
    Url  string
}

type CardPlayer struct {
    Video Video
    Image Image // Fallback. Same aspect ratio as video. At least 68,600 pixels.
}

// App represents a specific app on a specific app store.
type App struct {
    Store AppStore // e.g. AppStoreIPad
    Name string // e.g. "My App"
    ID   string // ID on App store e.g. "1234567890" or "org.example.myapp"
    Url string  // Optional. Your app’s custom URL scheme for App "deep links".
}

type AppStore string
const (
    AppStoreIPad       = AppStore("ipad")
    AppStoreIPhone     = AppStore("iphone")
    AppStoreGooglePlay = AppStore("googleplay")
)

type CardApp struct {
    // If your application is not available in the US App Store, you must set this
    // value to the two-letter country code for the App Store that contains your
    // application.
    Country string // e.g. "US".
    Apps []App
}
