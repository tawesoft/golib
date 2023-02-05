// Package dialog implements native, cross-platform message boxes,
// yes/no/okay/cancel confirmation prompts, file pickers, and others.
//
// This is a light-weight implementation (without using cgo or gtk bindings
// etc.) for developers who just need these basic features with a basic API.
//
// All dialogs will default to using localised text for window titles, buttons,
// etc. where possible, but may default to using English in places depending on
// the implementation.
//
// ## Alternatives
//
// There are more complete options. Here are some:
//
//   - [sqweek/dialog]
//   - [ncruces/zenity]
//
// [sqweek/dialog]: https://github.com/sqweek/dialog
// [ncruces/zenity]: https://github.com/ncruces/zenity
//
// ## Please note!
//
// The message box appears on the local system. If you are writing a web
// application and want a message to appear in a client's web browser, output
// HTML such as "<script>alert('Hello');</script>" instead!
//
// ## Windows
//
// On Windows, this package uses the native windows dialogs, converts Go
// strings into Windows UTF-16, handles null terminators, and uses the
// right-to-left display mode when using a RTL writing system such as Arabic or
// Hebrew. This package uses the older (pre-Vista) APIs because the new APIs
// make an awkward mix with Go's concurrency model and we don't want that to
// complicate our API just for simple features.
//
// For modern/prettier buttons on Windows using Common Control Styles ([Visual
// Styles]), see the examples folder for creating a manifest, (cross)compiling
// the manifest with a resource compiler ([akavel/rsrc]), and initialising
// comctl32.dll.
//
// [Visual Styles]: https://github.com/MicrosoftDocs/win32/blob/docs/desktop-src/Controls/visual-styles-overview.md
// [akavel/rsrc]: https://github.com/akavel/rsrc
//
// ## Other platforms
//
// On other systems (Linux, etc), this package uses (in order of priority) one
// or more of:
//
//   - zenity
//   - xmessage
//   - whiptail in an xterm
//   - osascript (Apple script) (TODO)
//
// ## Feature support
//
//   Platform/software | Message.Raise | Message.Ask | FilePicker | ColorPicker | DatePicker
//   ---------------------------------------------------------------------------------------
//   Windows           | Yes           | Yes         | Yes        |  No         |  No
//   zenity            | Yes           | Yes         | Yes        | Yes         | Yes
//   xmessage          | Yes           | Yes         |  No        |  No         |  No
//   whiptail + xterm  |  No           |  No         | Yes        | Yes         | Yes
//   osascript         | TODO          | TODO        | TODO       | TODO        | TODO
//
package dialog

import (
    "fmt"
    "image/color"
    "time"
)

// Init performs optional per-platform initialisation. On Windows, this loads
// modern visual styles (also requires the program is compiled with a matching
// manifest that enables visual styles). An error is not fatal. On platforms
// other than Windows, this currently doesn't do anything.
func Init() error {
    return osInit()
}

type Support struct {
    MessageRaise    bool // Can use Message.Raise?
    MessageAsk      bool // Can use Message.Ask?
    FilePicker      bool // Can use FilePicker.Open, FilePicker.Save?
    MultiFilePicker bool // Can use FilePicker.OpenMultiple?
    ColorPicker     bool // Can use ColorPicker.Pick?
    DatePicker      bool // Can use DatePicker.Pick?
}

// Supported returns a [Support] struct detailing what features are available
// on the current system. Using a feature that isn't supported will silently
// proceed as documented.
func Supported() (Support, error) {
    return supported()
}

// Alert is like [Raise] and the other convenience methods, but doesn't return
// any error message on failure.
//
// Deprecated. This is here for legacy reasons. Use [Raise], [Info], [Warning]
// etc. instead and just ignore the returned error if you don't care about it.
// Will be removed in golib/v3 but always available in golib/v2.
func Alert(message string, args...interface{}) {
    _ = Warning(message, args...)
}

// Ask is a convenience function to display a modal message box asking a
// question. The message string can be a printf-style format string for an
// optional sequence of additional arguments of any type. It blocks until an
// option is picked. Where not supported, immediately returns true without
// blocking.
func Ask(message string, args...interface{}) (bool, error) {
    return Message{
        Format: message,
        Args:   args,
    }.Ask()
}

// Raise is a convenience function to display a modal message box with a
// message. The message string can be a printf-style format string for an
// optional sequence of additional arguments of any type. It blocks until an
// option is picked. Where not supported, immediately returns without blocking.
//
// Raise uses the default icon and title bar, which is currently the same as
// [Info].
func Raise(message string, args...interface{}) error {
    return Info(message, args...)
}

// Info is like [Raise], but uses an information icon and title bar.
func Info(message string, args...interface{}) error {
    return Message{
        Format: message,
        Args:   args,
        Icon:   IconInfo,
    }.Raise()
}

// Warning is like [Raise], but uses a warning icon and title bar.
func Warning(message string, args...interface{}) error {
    return Message{
        Title:  "Warning", // may be localised
        Format: message,
        Args:   args,
        Icon:   IconWarning,
    }.Raise()
}

// Error is like [Raise], but uses an error icon and title bar.
func Error(message string, args...interface{}) error {
    return Message{
        Title:  "Error", // may be localised
        Format: message,
        Args:   args,
        Icon:   IconError,
    }.Raise()
}

// Color is a convenience function to display a colour picker. Use the
// [ColorPicker.Pick] method on a configured [ColorPicker] for more options. It
// blocks until an option is picked. Where not supported, immediately returns
// false without blocking.
func Color() (color.Color, bool, error) {
    return ColorPicker{
    }.Pick()
}

// Pick displays a colour selection dialog to select a single colour. It blocks
// until an option is picked, then returns the selected colour and true, or
// false if no colour was selected (i.e. the user selected the cancel option).
//
// Where not supported, immediately returns (zero, false, nil) without blocking
// (see [Supported]).
func (m ColorPicker) Pick() (color.Color, bool, error) {
    if m.Title == "" { m.Title = "Select color" }
    return m.pick()
}

// ColorPicker is a dialog to select a colour.
type ColorPicker struct {
    // Title is the colour picker window title (may be empty) e.g. "Colour".
    // If omitted, defaults to en-US "Color".
    Title string

    // Palette controls the color picker mode. If false, this is a simple
    // colour picker suited for picking a colour as a "one-off" e.g. with
    // sliders or by typing in a value. If true, where supported, this is
    // an extended picker with support for a palette where the user can
    // define favourite colours and easily pick those colours again later.
    Palette bool

    // Initial specifies the colour initially picked by default in the dialog.
    // For example, in a painting tool, you might want to set a default to
    // black. Or, if your application manages its own palette of colours, you
    // might want to launch the color picker with an existing palette colour so
    // that the user can adjust it slightly.
    Initial color.Color
}

// Date is a convenience function to display a date picker. Use the
// [DatePicker.Pick] method on a configured [DatePicker] for more options. It
// blocks until an option is picked. Where not supported, immediately returns
// (zero, false, nil) without blocking.
func Date() (time.Time, bool, error) {
    return DatePicker{
    }.Pick()
}

// Pick displays a (day, month, year( selection dialog to select a single date.
// It blocks until an option is picked, then returns the selected date and
// true, or false if no date was selected (i.e. the user selected the cancel
// option).
//
// The returned date is location and timezone aware. If the [DatePicker]
// Location is nil, this defaults to the user's configured location and
// timezone (i.e. [time.Local]). For example, set this to [time.UTC].
//
// Where not supported, immediately returns (zero, false, nil) without blocking
// (see [Supported]).
func (m DatePicker) Pick() (time.Time, bool, error) {
    // m.Title has sensible defaults on some implementations, so not set here
    if m.Location == nil { m.Location = time.Local }
    if m.Initial.IsZero() { m.Initial = time.Now() }
    return m.pick()
}

// DatePicker is a dialog to select a date (year, month, day).
type DatePicker struct {
    // Title is the date picker window title (may be empty). For example,
    // "Date of publication".
    Title string

    // LongTitle is some extra text (may be empty) that, where supported, can
    // give the user further context for the date picker. For example, "Select
    // date article first published".
    LongTitle string

    // Initial specifies the date initially picked by default in the dialog.
    // If not set, defaults to the current day.
    Initial time.Time

    // Location is used to return a location and timezone-aware time. If nil,
    // defaults to the user's configured location and timezone (i.e. [time.Local]).
    // For example, set this to [time.UTC].
    Location *time.Location
}

// Open is a convenience function to display a file picker dialog to select a
// single file. It blocks until an option is picked, then returns the selected
// path as an absolute path, and true, or an empty string and false if no file
// was selected (i.e. the user selected the cancel option).
//
// The provided file path argument is the initial file selected (if empty,
// defaults to current working directory). To open in a specific directory
// without specifying a file name, use a trailing slash.
//
// Where not supported, immediately returns ("", false, nil) without blocking
// (see [Supported]).
//
// Note that this does not actually read the file or open it for writing,
// but merely selects a path.
func Open(file string) (string, bool, error) {
    return FilePicker{
        Path: file,
    }.Open()
}

// OpenMultiple is like [Open], but allows multiple files to be selected. Each
// returned path is still an absolute path.
func OpenMultiple(file string) ([]string, bool, error) {
    return FilePicker{
        Path: file,
    }.OpenMultiple()
}

// Save is like [Open], but for picking a file to write to. This may change the
// look of the file picker (e.g. to have a button that says "Save" instead of
// "Open").
//
// Note that this does not actually write to the file or open it for writing,
// but merely selects a path.
func Save(file string) (string, bool, error) {
    return FilePicker{
        Path: file,
    }.Save()
}

// FilePicker is a dialog to select file(s) to load or save.
type FilePicker struct {
    // Title is the file picker window title (may be empty) e.g. "Open" or
    // "Save". If the implementation has a locale-aware default, then this
    // is ignored and the default is used instead.
    Title string

    // Path is the initial file selected (if empty, defaults to
    // current working directory). To open in a specific directory without
    // specifying a file name, use a trailing slash.
    Path string

    // FileTypes is hint describing known file types and file extensions. It
    // may be used to filter visible files. May be nil. This is a slice of
    // 2-tuples. The first item in the tuple is a human-readable label, the
    // second item in the tuple is a list of patterns delimited by space. Can
    // be left as nil as default. ("All Files", "*.*") is automatically added
    // to the end if the last item doesn't have the exact filter "*.*".
    //
    // For example:
    //
    //   FileTypes: [][2]string{
    //       {"Text Document", "*.txt *.rtf"},
    //       {"Image",         "*.png"},
    //       ...
    //       {"Pob Ffeil",      "*.*"}, // suppress English "All Files"
    //   }
    FileTypes [][2]string

    // DefaultFileType is an index into the FileTypes array identifying the
    // default file type to use.
    DefaultFileType int

    // AlwaysShowHidden, if true, is a hint that hidden files should always be
    // revealed if possible. If false, is a hint that hidden files should be
    // shown, or not shown, as normal depending on the user's settings.
    AlwaysShowHidden bool

    // AddToRecent, if true, is a hint that the opened or saved file should
    // be added to the user's history of recent files. If false, is a hint
    // that the file should not be added to that history.
    AddToRecent bool
}

// clear sets defaults on the file picker (and makes its own copy of the file
// types slice).
func (m FilePicker) clear() FilePicker {
    if m.FileTypes == nil { m.FileTypes = [][2]string{} }
    fileTypes := make([][2]string, 0, len(m.FileTypes) + 1)
    fileTypes = append(fileTypes, m.FileTypes...)

    if (len(fileTypes)) > 0 && (fileTypes[len(fileTypes) - 1][1] != "*.*") {
        fileTypes = append(fileTypes, [2]string{"All Files", "*.*"})
    } else if len(fileTypes) == 0 {
        fileTypes = append(fileTypes, [2]string{"All Files", "*.*"})
    }

    m.FileTypes = fileTypes
    return m
}

// Open displays a file picker dialog to select a single file. It blocks until
// an option is picked, then returns the selected path as an absolute path, and
// true, or an empty string and false if no file was selected (i.e. the user
// selected the cancel option).
//
// Where not supported, immediately returns ("", false, nil) without blocking
// (see [Supported]).
//
// Note that this does not actually read the file or open it for writing,
// but merely selects a path.
func (m FilePicker) Open() (string, bool, error) {
    m = m.clear()
    if m.Title == "" { m.Title = "Open File..." }
    return m.open()
}

// OpenMultiple is like [FilePicker.Open], but allows multiple files to be
// selected. Each returned path is still an absolute path.
func (m FilePicker) OpenMultiple() ([]string, bool, error) {
    m = m.clear()
    if m.Title == "" { m.Title = "Open Files..." }
    return m.openMultiple()
}

// Save is like [FilePicker.Open], but for writing to a file. This may change
// the look of the file picker (e.g. to have a button that says "Save" instead
// of "Open").
//
// Note that this does not actually write to the file or open it for writing,
// but merely selects a path.
func (m FilePicker) Save() (string, bool, error) {
    m = m.clear()
    if m.Title == "" { m.Title = "Save As..." }
    return m.save()
}

type IconType int
const (
    IconInfo IconType = iota
    IconWarning
    IconError
)

// Message is a prompt or question.
type Message struct {
    // Title is the message box window title (may be empty).
    Title string

    // Format is a printf-style format string. This is word-wrapped for you.
    Format string

    // Args are printf-style arguments
    Args []interface{}

    // Icon is a IconInfo, IconWarning, or IconError and is displayed when
    // possible on the message box.
    Icon IconType
}

// WithArgs returns Message with the Args field set.
func (m Message) WithArgs(args ... interface{}) Message {
    m.Args = args
    return m
}

// clear returns a Message with its format and args fields cleared, and its
// message formatted with its args, so that implementations don't incorrectly
// format it themselves. Implementations still control wrapping, where needed.
func (m Message) clear() (Message, string) {
    result := fmt.Sprintf(m.Format, m.Args...)
    m.Format = ""
    m.Args = nil
    return m, result
}

// Ask displays a message as a question and returns true if the affirmative
// option (such as "yes" or "okay") is picked, or false if the negative option
// (such as "no" or "cancel") is picked. It blocks until an option is picked.
//
// Where fully supported by a platform, the options are "yes" and "no" as
// appropriate for the current language and locale. Otherwise, the options may
// default to English-language.
//
// Where not supported, immediately returns true without blocking.
//
// Ask also ignores the supplied icon option and sets an appropriate question
// icon instead.
func (m Message) Ask() (bool, error) {
    n, s := m.clear()
    if n.Title == "" { n.Title = "Question" }
    return n.ask(s)
}

// Raise displays a message. It blocks until the message is acknowledged
// (for example by clicking "Okay").
//
// Where not supported, returns immediately without blocking.
func (m Message) Raise() error {
    n, s := m.clear()
    if n.Title == "" { n.Title = "Message" }
    return n.raise(s)
}
