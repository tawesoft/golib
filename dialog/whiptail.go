//go:build (linux || unix)

package dialog

import (
    "fmt"
    "image/color"
    "io"
    "os"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/tawesoft/golib/v2/operator"
)

type whiptail struct {
    shell    string
    xterm    string
    whiptail string
}

func tryGetGeometry(shell string) (int, int) {
    xwininfo, err := find("xwininfo")
    if err != nil { return 0, 0 }

    grep, err := find("grep")
    if err != nil { return 0, 0 }

    xargs, err := find("xargs")
    if err != nil { return 0, 0 }

    cut, err := find("cut")
    if err != nil { return 0, 0 }

    var sb strings.Builder
    cmd := exec.Command(shell, "-c",
        fmt.Sprintf(`%s -root | %s -E "\-geometry" | %s | %s -f2 -d " " | %s -f1 -d "+"`,
            clean(xwininfo),
            clean(grep),
            clean(xargs),
            clean(cut),
            clean(cut),
        ),
    )
    cmd.Stdout = &sb
    err = cmd.Run()
    if err != nil { return 0, 0 }

    geom := strings.TrimSpace(sb.String()) // e.g. 2560x1440
    w, h, found := strings.Cut(geom, "x")
    if !found { return 0, 0 }

    pw, err := strconv.Atoi(w)
    if err != nil { return 0, 0 }
    ph, err := strconv.Atoi(h)
    if err != nil { return 0, 0 }

    return pw, ph
}

func (x whiptail) getString(title string, label string, placeholder string) (string, bool, error) {

    f, err := os.CreateTemp("", "dialog")
    if err != nil {
        return "", false, fmt.Errorf("error creating temporary communication file %w", err)
    }
    defer os.Remove(f.Name())

    width, height := tryGetGeometry(x.shell)
    if width == 0  { width  = 10 } else { width = (width / 2) - 242 }
    if height == 0 { height = 10 } else { height = (height / 2) - 158 }

    cmd := exec.Command(x.xterm,
        "-geometry", fmt.Sprintf("80x24+%d+%d", width, height),
        "-T", title,
        "-e", x.shell, "-l", "-c",
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
        return strings.TrimSpace(string(in)), len(in) != 0, nil
    }
}

var whiptailColorPickerRE = regexp.MustCompile(`^#(?P<red>[[:xdigit:]]{2})(?P<green>[[:xdigit:]]{2})(?P<blue>[[:xdigit:]]{2})$`)

func (x whiptail) color(m ColorPicker) (color.Color, bool, error) {
    zero := operator.Zero[color.Color]()

    r, g, b, _ := m.Initial.RGBA()
    r /= 256; g /= 256; b /= 256;
    initial := fmt.Sprintf("#%02x%02x%02x", r, g, b)

    result, ok, err := x.getString(m.Title, "Color (hexadecimal #RRGGBB):", initial)
    if !ok { return zero, ok, nil }
    if err != nil {
        return zero, ok, fmt.Errorf("error picking color: %v", err)
    }

    matches := whiptailColorPickerRE.FindStringSubmatch(result)
    if len(matches) != 4 {
        return zero, false, fmt.Errorf("whiptail color selection parse error parsing %q", result)
    }
    cr, errR := strconv.ParseInt(matches[1], 16, 16)
    cg, errG := strconv.ParseInt(matches[2], 16, 16)
    cb, errB := strconv.ParseInt(matches[3], 16, 16)
    if (errR != nil) || (errG != nil) || (errB != nil) {
        return zero, false, fmt.Errorf("whiptail color selection parse error parsing %q", result)
    }

    return color.RGBA{
        R: uint8(cr & 255),
        G: uint8(cg & 255),
        B: uint8(cb & 255),
        A: 255,
    }, true, nil
}

var whiptailDatePickerRE = regexp.MustCompile(`^(?P<year>\d+)/(?P<month>\d+)/(?P<day>\d+)$`)

func (x whiptail) date(m DatePicker) (time.Time, bool, error) {
    if m.Title == "" {
        m.Title = "Select Date"
    }

    if m.LongTitle == "" {
        m.LongTitle = "Date: (YYYY/MM/DD)"
    } else {
        m.LongTitle += " (YYYY/MM/DD)"
    }

    zero := operator.Zero[time.Time]()

    initial := fmt.Sprintf("%d/%d/%d",
        m.Initial.Year(), m.Initial.Month(), m.Initial.Day())

    result, ok, err := x.getString(m.Title, m.LongTitle, initial)
    if !ok { return zero, ok, nil }
    if err != nil {
        return zero, ok, fmt.Errorf("error picking date: %v", err)
    }

    matches := whiptailDatePickerRE.FindStringSubmatch(result)
    if len(matches) != 4 {
        return zero, false, fmt.Errorf("whiptail date selection parse error parsing %q", result)
    }
    dy, errY := strconv.Atoi(matches[1])
    dm, errM := strconv.Atoi(matches[2])
    dd, errD := strconv.Atoi(matches[3])
    if (errY != nil) || (errM != nil) || (errD != nil) {
        return zero, false, fmt.Errorf("whiptail color selection parse error parsing %q", result)
    }

    return time.Date(dy, time.Month(dm), dd, 0, 0, 0, 0, m.Location), true, nil
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

    return x.getString(m.Title, m.Title, path)
}

func (x whiptail) save(m FilePicker) (string, bool, error) {
    path := m.Path
    if path == "" {
        cwd, err := os.Getwd()
        if err != nil {
            return "", false, fmt.Errorf("error getting working directory: %v", err)
        }
        path = cwd
    }

    return x.getString(m.Title, m.Title, path)
}
