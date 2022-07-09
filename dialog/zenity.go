//go:build (linux || unix)

package dialog

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type zenity struct {
    path string
}

func (x zenity) iconString(i IconType) string {
    switch i {
        case IconInfo: return "info"
        case IconWarning: return "warning"
        case IconError: return "error"
        default: return "info"
    }
}

func (x zenity) ask(m Message, message string) (bool, error) {
cmd := exec.Command(x.path,
        "--question",
        "--no-markup",
        "--default-cancel",
        "--width",       "480",
        "--window-icon", "question",
        "--title",       m.Title,
        "--text",        message,
    )

    if err := cmd.Run(); err != nil {
        if ExitError, ok := err.(*exec.ExitError); ok && (ExitError.ExitCode() == 1) {
            return false, nil // "no"
        } else {
            return false, fmt.Errorf("zenity error: %w", err)
        }
    }

    return true, nil // "yes"
}

func (x zenity) raise(m Message, message string) error {
    cmd := exec.Command(x.path,
        "--"+x.iconString(m.Icon),
        "--no-markup",
        "--width",      "480",
        "--window-icon", x.iconString(m.Icon),
        "--title",       m.Title,
        "--text",        message,
    )

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("zenity error: %w", err)
    }

    return nil
}

func (x zenity) pickFile(
    m FilePicker,
    mode rune, // (o)pen, (m)ultiple, (s)ave
) ([]string, bool, error) {
    path := m.Path
    if path == "" {
        cwd, err := os.Getwd()
        if err != nil {
            return nil, false, fmt.Errorf("error getting working directory: %v", err)
        }
        path = cwd + "/"
    }

    const sep = ","

    args := []string{
        "--file-selection",
        "--filename", path,
        "--separator", ",",
    }

    if mode == 'm' {
        args = append(args, "--multiple")
    } else if mode == 's' {
        args = append(args, "--save")
    }

    // append filters...
    for _, f := range m.FileTypes {
        name, patterns := f[0], f[1]

        // name can't contain a stave, because zenity uses that
        // TODO we should clean the string rather than skip
        if strings.ContainsRune(name, '|') { continue }

        var filter strings.Builder
        filter.WriteString(fmt.Sprintf("%s (%s)",
            name,
            strings.Join(strings.Split(patterns, " "), ", "),
        ))
        filter.WriteString(" | ")
        filter.WriteString(patterns) // space separated is fine as-is
        args = append(args, "--file-filter", filter.String())
    }

    var sb strings.Builder
    cmd := exec.Command(x.path, args...)
    cmd.Stdout = &sb
    fmt.Println(cmd)
    if err := cmd.Run(); err != nil {
        return nil, false, fmt.Errorf("zenity error: %w", err)
    } else {
        f := sb.String()
        if len(f) == 0 {
            return nil, false, nil
        }
        if mode == 'm' {
            return strings.Split(f, sep), true, nil
        } else {
            return []string{f}, true, nil
        }
    }
}

func (x zenity) open(m FilePicker) (string, bool, error) {
    if xs, ok, err := x.pickFile(m, 'o'); err != nil {
        return "", false, fmt.Errorf("error opening file picker: %w", err)
    } else if !ok {
        return "", false, nil
    } else {
        if len(xs) > 0 {
            return xs[0], true, nil
        } else {
            return "", false, nil
        }
    }
}

func (x zenity) openMultiple(m FilePicker) ([]string, bool, error) {
    if xs, ok, err := x.pickFile(m, 'm'); err != nil {
        return []string{}, false, fmt.Errorf("error opening save file picker: %w", err)
    } else if !ok {
        return []string{}, false, nil
    } else {
        if len(xs) > 0 {
            return xs, true, nil
        } else {
            return []string{}, false, nil
        }
    }
}

func (x zenity) save(m FilePicker) (string, bool, error) {
    if xs, ok, err := x.pickFile(m, 's'); err != nil {
        return "", false, fmt.Errorf("error opening save file picker: %w", err)
    } else if !ok {
        return "", false, nil
    } else {
        if len(xs) > 0 {
            return xs[0], true, nil
        } else {
            return "", false, nil
        }
    }
}
