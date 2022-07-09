// Package dialog implements native, cross-platform message boxes,
// yes/no/okay/cancel confirmation prompts, file pickers, and others.
//
// This is a light-weight implementation (without using cgo or gtk bindings
// etc.) for developers who just need these basic features with a basic API.
//
// There are more complete options. Here are some:
//
//   - [sqweek/dialog]
//   - [ncruces/zenity]
//
// [sqweek/dialog]: https://github.com/sqweek/dialog
// [ncruces/zenity]:https://github.com/ncruces/zenity
//
// On Windows, this package uses the native windows dialogs, converts Go
// strings into Windows UTF-16, handles null terminators, and (TODO) uses
// right-to-left display mode when using a RTL writing system such as Arabic or
// Hebrew. This package uses the older (pre-Vista) APIs because the new APIs
// make an awkward mix with Go's concurrency model and we don't want that to
// complicate our API just for simple features.
//
// On other systems (Linux, etc), this package uses (in order of priority) one
// or more of:
//
//   - zenity
//   - xmessage
//   - whiptail in an xterm
//   - osascript (Apple script) (TODO) (Note: not tested!)
//
// Feature support:
//
//   Platform/software | Message.Raise | Message.Ask | FilePicker
//   ------------------------------------------------------------
//   Windows           | Yes           | Yes         | Yes
//   zenity            | Yes           | Yes         | Yes
//   xmessage          | Yes           | Yes         |  X
//   whiptail + xterm  | Yes           | Yes         | Yes
//
// Additional feature support:
//
//   Platform/software | ColorPicker | DatePicker
//   ----------------------------------------------------------
//   Windows           |  X          |  X
//   zenity            | Yes         | Yes
//   xmessage          |  X          |  X
//   whiptail + xterm  | Yes         | Yes
//
// Please note: the message box appears on the local system. If you are writing
// a web application and want a message to appear in a client's web browser,
// output HTML such as "<script>alert('Hello');</script>" instead!
package dialog

import (
    "fmt"
)

type Support struct {
    MessageRaise  bool // Can use Message.Raise?
    MessageAsk    bool // Can use Message.Ask?
    FilePicker    bool // Can use FileSelect.Open, FileSelect.Save, etc?
    ColorPicker   bool // Can use ColorPicker.Pick?
    DatePicker    bool // Can use DatePicker.Pick?
}

// Supported returns a [Support] struct detailing what features are available
// on the current system. Using a feature that isn't supported will silently
// proceed as documented.
func Supported() (Support, error) {
    return supported()
}

// Alert is like [Raise], but doesn't return any error message on failure.
//
// Deprecated. This is here for legacy reasons.
func Alert(message string, args...interface{}) {
    Message{
        Title:  "Alert",
        Format: message,
        Args:   args,
        Icon:   IconInfo,
    }.Raise()
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
func Raise(message string, args...interface{}) error {
    return Message{
        Title:  "Alert",
        Format: message,
        Args:   args,
        Icon:   IconInfo,
    }.Raise()
}

// Open is a convenience function to display a file picker dialog to select a
// single file. It blocks until an option is picked, then returns the selected
// path as an absolute path, and true, or an empty string and false if no file
// was selected (i.e. the user selected the cancel option).
//
// Where not supported, immediately returns ("", false, nil) without blocking
// (see [Supported]).
//
// Note that this does not actually read the file or open it for writing,
// but merely selects a path.
func Open(file string) (string, bool, error) {
    return FilePicker{
    }.Open()
}

// OpenMultiple is like [Open], but allows multiple files to be selected. Each
// returned path is still an absolute path.
func OpenMultiple(file string) ([]string, bool, error) {
    return FilePicker{
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
    // may be used to filter visible files. May be nil.
    //
    //   FileTypes: [][2]string{
    //       {"Text Document", "*.txt; *.rtf"},
    //       {"All Documents, ".*."},
    //       ...
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
    if m.Title == "" { m.Title = "Open File..." }
    return m.open()
}

// OpenMultiple is like [FilePicker.Open], but allows multiple files to be
// selected. Each returned path is still an absolute path.
func (m FilePicker) OpenMultiple() ([]string, bool, error) {
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
