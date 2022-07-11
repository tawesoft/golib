//go:build (linux || unix)

package dialog

import (
    "fmt"
    "image/color"
    "os"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
    "time"
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
    args := []string{
        "--question",
        "--no-markup",
        "--default-cancel",
        "--width",       "480",
        "--window-icon", "question",
        "--text",        message,
    }

    if m.Title != "Question" {
        args = append(args, "--title", m.Title)
    }

    cmd := exec.Command(x.path, args...)

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
    args := []string{
        "--"+x.iconString(m.Icon),
        "--no-markup",
        "--width",      "480",
        "--window-icon", x.iconString(m.Icon),
        "--text",        message,
    }

    if (m.Title != x.iconString(m.Icon)) && (m.Title != "Message")  {
        args = append(args, "--title", m.Title)
    }

    cmd := exec.Command(x.path, args...)

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("zenity error: %w", err)
    }

    return nil
}

var zenityColorPickerRE = regexp.MustCompile(`^rgb\((?P<red>\d+),(?P<green>\d+),(?P<blue>\d+)\)$`)

func (x zenity) color(m ColorPicker) (color.Color, bool, error) {
    var zero color.Color

    args := []string{
        "--color-selection",
        "--title", m.Title,
    }

    if m.Palette {
        args = append(args, "--show-palette")
    }

    if m.Initial != zero {
        r, g, b, _ := m.Initial.RGBA()
        r /= 256; g /= 256; b /= 256;
        hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)
        args = append(args, "--color", hex)
    }

    var sb strings.Builder
    cmd := exec.Command(x.path, args...)
    cmd.Stdout = &sb

    if err := cmd.Run(); err != nil {
        if ExitError, ok := err.(*exec.ExitError); ok && (ExitError.ExitCode() == 1) {
            return zero, false, nil // cancel
        } else {
            return zero, false, fmt.Errorf("zenity error: %w", err)
        }
    }

    result := strings.TrimSpace(sb.String())
    matches := zenityColorPickerRE.FindStringSubmatch(result)
    if len(matches) != 4 {
        return zero, false, fmt.Errorf("zenity --color-selection parse error parsing %q", result)
    }
    cr, errR := strconv.Atoi(matches[1])
    cg, errG := strconv.Atoi(matches[2])
    cb, errB := strconv.Atoi(matches[3])
    if (errR != nil) || (errG != nil) || (errB != nil) {
        return zero, false, fmt.Errorf("zenity --color-selection parse error parsing %q", result)
    }

    return color.RGBA{
        R: uint8(cr & 255),
        G: uint8(cg & 255),
        B: uint8(cb & 255),
        A: 255,
    }, true, nil
}

func (x zenity) date(m DatePicker) (time.Time, bool, error) {
    var zero time.Time

    args := []string{
        "--calendar",
        "--day",         fmt.Sprintf("%d", m.Initial.Day()),
        "--month",       fmt.Sprintf("%d", m.Initial.Month()),
        "--year",        fmt.Sprintf("%d", m.Initial.Year()),
        "--date-format", "%Y%m%d",
    }

    if len(m.Title) > 0 {
        args = append(args, "--title", m.Title)
    }

    if len(m.LongTitle) > 0 {
        args = append(args, "--text", m.LongTitle)
    }

    var sb strings.Builder
    cmd := exec.Command(x.path, args...)
    cmd.Stdout = &sb

    if err := cmd.Run(); err != nil {
        if ExitError, ok := err.(*exec.ExitError); ok && (ExitError.ExitCode() == 1) {
            return zero, false, nil // cancel
        } else {
            return zero, false, fmt.Errorf("zenity error: %w", err)
        }
    }

    ymd := strings.TrimSpace(sb.String()) // YYYYMMDD
    if len(ymd) != 8 {
        return zero, false, fmt.Errorf("zenity --calendar format error")
    }
    dy, errY := strconv.Atoi(ymd[0:4])
    dm, errM := strconv.Atoi(ymd[4:6])
    dd, errD := strconv.Atoi(ymd[6:8])
    if (errY != nil) || (errM != nil) || (errD != nil) {
        return zero, false, fmt.Errorf("zenity --calendar format error")
    }

    return time.Date(dy, time.Month(dm), dd, 0, 0, 0, 0, m.Location), true, nil
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
    if err := cmd.Run(); err != nil {
        if ExitError, ok := err.(*exec.ExitError); ok && (ExitError.ExitCode() == 1) {
            return nil, false, nil // closed/cancelled
        } else {
            return nil, false, fmt.Errorf("zenity error: %w", err)
        }
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
