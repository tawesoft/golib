//go:build (linux || unix)

package dialog

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "strconv"
    "strings"
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
