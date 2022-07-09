package dialog

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "strings"

    "github.com/alessio/shellescape"
)

type whiptail struct {
    xterm    string
    whiptail string
}

func (x whiptail) getString(title string, label string, placeholder string) (string, bool, error) {

    f, err := os.CreateTemp("", "dialog")
    if err != nil {
        return "", false, fmt.Errorf("error creating temporary communication file %w", err)
    }
    defer os.Remove(f.Name())

    clean := func(x string) string {
        return shellescape.Quote(x)
    }

    cmd := exec.Command(x.xterm,
        "-T", title,
        "-e", "/bin/sh", "-l", "-c",
        strings.Join([]string{
            clean(x.whiptail),
                "--inputbox",
                    clean(label), "8", "70", clean(placeholder),
                "--title",
                    clean(title),
        }, " ") + " 2> " + clean(f.Name()),
    )
    if err := cmd.Run(); err != nil {
        return "", false, fmt.Errorf("xterm/whiptail error: %v", err)
    }

    if in, err := io.ReadAll(f); err != nil {
        return "", false, fmt.Errorf("I/O error: %v", err)
    } else {
        return string(in), len(in) != 0, nil
    }
}

func (x whiptail) open(m FilePicker) (string, bool, error) {
    path := m.Path
    if path == "" {
        cwd, err := os.Getwd()
        if err != nil {
            return "", false, fmt.Errorf("error getting working directory: %v", err)
        }
        path = cwd
    }

    return x.getString(m.Title, "Path:", path)
}
