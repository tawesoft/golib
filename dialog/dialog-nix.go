//go:build (linux || unix)

package dialog

import (
    "fmt"
    "image/color"
    "os/exec"
    "strings"
    "time"

    "github.com/alessio/shellescape"
    "github.com/tawesoft/golib/v2/operator"
    "golang.org/x/sys/execabs"
)

// Compile-time constants to enable/disable implementations, even when they're
// found at runtime.
const (
    enableWhiptail = true
    enableXMessage = true
    enableZenity   = true
)

var (
    gfuns  funs  // global init of functions to implementations
    giniterr error
)

type funs struct {
    messageRaise         func(m Message, message string) error
    messageAsk           func(m Message, message string) (bool, error)
    filePickOpen         func(m FilePicker)  (string, bool, error)
    filePickOpenMultiple func(m FilePicker)  ([]string, bool, error)
    filePickSave         func(m FilePicker)  (string, bool, error)
    colorPick            func(m ColorPicker) (color.Color, bool, error)
    datePick             func(m DatePicker)  (time.Time, bool, error)
}

func clean(x string) string {
    return shellescape.Quote(x)
}

func osInit() error {
    return nil
}

func init() {
    type paths struct {
        shell    string
        whiptail string
        xmessage string
        xterm    string
        zenity   string
    }

    stash := func(dest *string, bin string) error {
        path, err := find(bin)
        *dest = path
        return err
    }

    var err error
    p := &paths{}
    if err == nil { err = stash(&p.shell,    "sh") }
    if err == nil { err = stash(&p.whiptail, "whiptail") }
    if err == nil { err = stash(&p.xmessage, "xmessage") }
    if err == nil { err = stash(&p.xterm,    "xterm") }
    if err == nil { err = stash(&p.zenity,   "zenity"  ) }
    giniterr = err

    if (p.zenity != "") && enableZenity {
        z := zenity{p.zenity}
        if gfuns.messageAsk == nil {
            gfuns.messageAsk = z.ask
        }
        if gfuns.messageRaise == nil {
            gfuns.messageRaise = z.raise
        }
        if gfuns.filePickOpen == nil {
            gfuns.filePickOpen = z.open
        }
        if gfuns.filePickOpenMultiple == nil {
            gfuns.filePickOpenMultiple = z.openMultiple
        }
        if gfuns.filePickSave == nil {
            gfuns.filePickSave = z.save
        }
        if gfuns.colorPick == nil {
            gfuns.colorPick = z.color
        }
        if gfuns.datePick == nil {
            gfuns.datePick = z.date
        }
    }

    if (p.xmessage != "") && enableXMessage {
        x := xmessage{p.xmessage}
        if gfuns.messageAsk == nil {
            gfuns.messageAsk = x.ask
        }
        if gfuns.messageRaise == nil {
            gfuns.messageRaise = x.raise
        }
    }

    if (p.shell != "") && (p.xterm != "") && (p.whiptail != "") && enableWhiptail {
        w := whiptail{
            shell:    p.shell,
            xterm:    p.xterm,
            whiptail: p.whiptail,
        }
        if gfuns.filePickOpen == nil {
            gfuns.filePickOpen = w.open
        }
        if gfuns.filePickSave == nil {
            gfuns.filePickSave = w.save
        }
        if gfuns.colorPick == nil {
            gfuns.colorPick = w.color
        }
        if gfuns.datePick == nil {
            gfuns.datePick = w.date
        }
    }
}

func supported() (Support, error) {
    return Support{
        MessageRaise:    gfuns.messageRaise != nil,
        MessageAsk:      gfuns.messageAsk   != nil,
        FilePicker:      gfuns.filePickOpen != nil,
        MultiFilePicker: gfuns.filePickOpenMultiple != nil,
        DatePicker:      gfuns.datePick     != nil,
        ColorPicker:     gfuns.colorPick    != nil,
    }, giniterr
}

func find(bin string) (string, error) {
    // "In older versions of Go, LookPath could return a path relative to the
    // current directory. As of Go 1.19, LookPath will instead return that path
    // along with an error satisfying errors.Is(err, ErrDot). See the package
    // documentation for more details."

    // Due to build constraints, we don't have to care about the
    // “what to do about PATH lookups on Windows” "saga", but some unix users
    // might have this set insecurely, too.

    // We use execabs.Command, which should be safe on old versions, too.

    var buf strings.Builder
    cmd := execabs.Command("which", bin)
    cmd.Stdout = &buf

    if err := cmd.Run(); err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            if 1 == exitError.ExitCode() {
                return "", nil
            }
        }
        return "", fmt.Errorf("error running 'which' command: %v", err)
    }

    return strings.TrimSpace(buf.String()), nil
}

func (m ColorPicker) pick() (color.Color, bool, error) {
    if gfuns.colorPick == nil { return operator.Zero[color.Color](), false, nil }
    return gfuns.colorPick(m)
}

func (m DatePicker) pick() (time.Time, bool, error) {
    if gfuns.datePick == nil { return operator.Zero[time.Time](), false, nil }
    return gfuns.datePick(m)
}

func (m FilePicker) open() (string, bool, error) {
    if gfuns.filePickOpen == nil { return "", false, nil }
    return gfuns.filePickOpen(m)
}

func (m FilePicker) openMultiple() ([]string, bool, error) {
    if gfuns.filePickSave == nil { return []string{}, false, nil }
    return gfuns.filePickOpenMultiple(m)
}

func (m FilePicker) save() (string, bool, error) {
    if gfuns.filePickSave == nil { return "", false, nil }
    return gfuns.filePickSave(m)
}

func (m Message) ask(message string) (bool, error) {
    if gfuns.messageAsk == nil { return false, nil }
    return gfuns.messageAsk(m, message)
}

func (m Message) raise(message string) error {
    if gfuns.messageRaise == nil { return nil }
    return gfuns.messageRaise(m, message)
}
