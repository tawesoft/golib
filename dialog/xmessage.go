//go:build (linux || unix)

package dialog

import (
    "fmt"
    "os/exec"
    "strings"

    "github.com/tawesoft/golib/v2/ks"
)

type xmessage struct {
    path string
}

func (x xmessage) exec(m Message, buttons []string) *exec.Cmd {
    var btns strings.Builder
    for i, b := range buttons {
        btns.WriteString(fmt.Sprintf("%s:%d", b, i + 1))
        if i + 1 != len(buttons) {
            btns.WriteString(",")
        }
    }

    return exec.Command(x.path,
        "-xrm",     "*international:true", // unicode support
        "-xrm",     "*fontSet:-*-fixed-medium-r-normal-*-20-*-*-*-*-*-*-*,-*-fixed-*-*-*-*-20-*-*-*-*-*-*-*,-*-*-*-*-*-*-20-*-*-*-*-*-*-*",
        "-title",   m.Title,
        "-g",       "600x400", // bigger window
        "-buttons", btns.String(), // return code is 100 + button code (unreliable)
        "-default", buttons[0],   // set focus
        "-center",
        "-print",        // return code isn't working properly!
        "-file",    "-", // read from stdin
    )
}

func (x xmessage) ask(m Message, message string) (bool, error) {
    var output strings.Builder
    message = ks.WrapBlock(message, 56)
    cmd := x.exec(m, []string{"Yes", "No"})
    cmd.Stdin = strings.NewReader(message)
    cmd.Stdout = &output

    if err := cmd.Run(); err == nil {
        // supposed to have a return code, but -print the label works too
        clicked := strings.TrimSpace(output.String())
        if clicked == "Yes" {
            return true, nil
        } else {
            return false, nil
        }
    } else {
        if exitError, ok := err.(*exec.ExitError); ok {
            pressed := exitError.ExitCode()
            if (pressed == 101) || (pressed == 1) {
                return true, nil // Yes
            } else if (pressed == 102) || (pressed == 2) {
                return false, nil // No
            } else {
                // assume X close button press
                return false, nil
            }
        } else {
            return false, err
        }
    }
}

func (x xmessage) raise(m Message, message string) error {
    message = ks.WrapBlock(message, 56)
    cmd := x.exec(m, []string{"OK"})
    cmd.Stdin = strings.NewReader(message)

    if err := cmd.Run(); err == nil {
        return nil // shouldn't happen, but ok
    } else {
        if _, ok := err.(*exec.ExitError); ok {
            // pressed := exitError.ExitCode()
            // assume X close button press
            return nil
        } else {
            return err
        }
    }
}
