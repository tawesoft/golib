//go:build windows

package dialog

import (
    "fmt"
    "image/color"
    "os"
    "path/filepath"
    "strings"
    "time"
    "unicode/utf16"
    "unsafe"

    "github.com/tawesoft/golib/v2/must"
    "github.com/tawesoft/golib/v2/operator"
    "golang.org/x/sys/windows"
    "golang.org/x/text/unicode/bidi"
)

func osInit() error {
    return initComctl32()
}

func initComctl32() error {
    libcomctl32 := windows.NewLazySystemDLL("comctl32.dll")
    if err := libcomctl32.Load(); err != nil {
        return fmt.Errorf("comctl32.dll load error: %w", err)
    }

    commonControlsEx := libcomctl32.NewProc("InitCommonControlsEx")
    if err := commonControlsEx.Find(); err != nil {
        return fmt.Errorf("could not find InitCommonControlsEx in comctl32.dll: %w", err)
    }

    type CommonControlsExArgs struct {
        DwSize, DwICC uint32
    }

    ctrls := &CommonControlsExArgs{
        DwSize: 8,
        DwICC: 0 |
            0x000000FF | // windows 95+
            0x00002000 | // fonts
            0x00004000, // buttons, etc.
    }

    ret, _, err := commonControlsEx.Call(uintptr(unsafe.Pointer(ctrls)))

    if ret != 1 { // (int true)
        if err != nil {
            return fmt.Errorf("comctl32.dll InitCommonControlsEx error: %w", err)
        }
    }

    return nil
}

var (
    dllComdlg32 = windows.NewLazySystemDLL("comdlg32.dll")

    procGetOpenFileName      = dllComdlg32.NewProc("GetOpenFileNameW")
    procGetSaveFileName      = dllComdlg32.NewProc("GetSaveFileNameW")
    procCommDlgExtendedError = dllComdlg32.NewProc("CommDlgExtendedError")
)

// https://docs.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-openfilenamew
type openFileNameW struct {
    lStructSize uint32 // unsafe.SizeOf
    hwndOwner   uintptr // leave as nil
    hInstance   uintptr // leave as nil

    // filetype filters.
    // Pairs of (label, filter1;filter2).
    // Nul delimited and double-nul terminated.
    // Remove whitespace between filters
    lpstrFilter uintptr
    lpstrCustomFilter uintptr // buffer for filter selected by user (can be nil)
    nMaxCustFilter uint32 // size in characters of lpstrCustomFilter
    nFilterIndex uint32 // index (get or set) of the filter (1-indexed, or 0 for custom filter)

    // initial or selected (get or set) file name; can be zero-terminated empty
    // string. For multi-select, nul delimited and double-nul terminated. On
    // success, contains the drive designator, path, file name, and extension of
    // selected file.
    lpstrFile uintptr
    nMaxFile uint32 // size in characters of above

    // file name and extension (without path information) of the selected file;
    // can be nil. Always leave this as nil because we use OFN_NOCHANGEDIR.
    // Get path from lpstrFile instead.
    lpstrFileTitle uintptr
    nMaxFileTitle uint32 // size in characters of above
    lpstrInitialDir uintptr // initial directory
    lpstrTitle uintptr // titlebar title; leave as nil to use locale-appropriate system default
    flags uint32

    // returned zero-based character offset of file name in lpstrFile
    // useless when selecting multiple files
    nFileOffset int32
    nFileExtension int32 // as above but for extension, useless for selecting multiple files

    // default extension added to the file name if the user fails to type an extension.
    // Is it better to ignore this and let client code add extension if desired?
    lpstrDefExt uintptr // leave as nil

    lCustData uintptr // leave as nill
    lpfnHook uintptr // leave as nil
    lpTemplateName uintptr // leave as nil
    pvReserved uintptr
    dwReserved uint32
    flagsEx uint32
}

func errSys(err error) error {
    return fmt.Errorf("windows syscall error: %w", err)
}

// splits a slice up to a null terminator, and the rest
func nullTerminated[T ~uint16 | rune](buf []T) ([]T, []T, bool) {
    for i := 0; i < len(buf); i++ {
        if buf[i] == 0 {
            if i + 1 < len(buf) {
                return buf[:i], buf[i+1:], true
            } else {
                return buf[:i], buf[0:0], true
            }
        }
    }
    return nil, nil, false
}

func widePermitNul(input string) *uint16 {
    s := utf16.Encode([]rune(input + "\x00"))
    return &s[0]
}

func pwidePermitNul(input string) uintptr {
    return uintptr(unsafe.Pointer(widePermitNul(input)))
}

func wide(input string) *uint16 {
    if result, err := windows.UTF16PtrFromString(input); err != nil {
        return must.Result(windows.UTF16PtrFromString("Unicode error"))
    } else {
        return result
    }
}

func pwide(input string) uintptr {
    return uintptr(unsafe.Pointer(wide(input)))
}

func (i IconType) iconFlag() uint32 {
    switch i {
        case IconInfo:    return windows.MB_ICONINFORMATION
        case IconWarning: return windows.MB_ICONWARNING
        case IconError:   return windows.MB_ICONERROR
    }

    must.Never()
    return 0
}

func supported() (Support, error) {
    return Support{
        MessageRaise:    true,
        MessageAsk:      true,
        FilePicker:      true,
        MultiFilePicker: true,
    }, nil
}

func (m ColorPicker) pick() (color.Color, bool, error) {
    return operator.Zero[color.Color](), false, nil
}

func (m DatePicker) pick() (time.Time, bool, error) {
    return operator.Zero[time.Time](), false, nil
}

func (m FilePicker) pick(
    mode rune, // (o)pen, (m)ultiple, (s)ave
) ([]string, bool, error) {
    var err error

    if err := procGetOpenFileName.Find(); err != nil {
        err = fmt.Errorf("missing requried comdlg32.dll procedure GetOpenFileNameW: %w", err)
        return nil, false, errSys(err)
    }

    if err := procGetSaveFileName.Find(); err != nil {
        err = fmt.Errorf("missing requried comdlg32.dll procedure GetSaveFileNameW: %w", err)
        return nil, false, errSys(err)
    }

    if err := procGetOpenFileName.Find(); err != nil {
        err = fmt.Errorf("missing requried comdlg32.dll procedure CommDlgExtendedError: %w", err)
        return nil, false, errSys(err)
    }

    var flags uint32 = 0 |
        0x00000008 | // OFN_NOCHANGEDIR  - don't let selecting a file change the program working directory
        0x00000004 | // OFN_HIDEREADONLY - don't show a readonly checkbox (very old!)
        0

    if (mode == 'o') || (mode == 'm') {
        flags |=
            0x00000800 | // OFN_PATHMUSTEXIST
            0x00001000 | // OFN_FILEMUSTEXIST
            0
    }

    if mode == 'm' {
        flags |=
            0x00000200 | // OFN_ALLOWMULTISELECT
            0x00080000 | // OFN_EXPLORER force more modern multiselect
            0
    }

    if mode == 's' {
        flags |=
            0x00008000 | // OFN_NOREADONLYRETURN
            0x00010000 | // OFN_NOTESTFILECREATE
            0
    }

    if !m.AddToRecent {
        flags |= 0x02000000 // OFN_DONTADDTORECENT
    }

    if m.AlwaysShowHidden {
        flags |= 0x10000000 // OFN_FORCESHOWHIDDEN
    }

    const bufSize = 16*1024 // UTF16 chars
    buf := make([]uint16, bufSize)
    /*
    base := filepath.Base(m.Path)
    if base != "" {
        if lpstrBase, err := windows.UTF16FromString(base); err == nil {
            copy(buf, lpstrBase)
        } else {
            return nil, false, fmt.Errorf("Unicode error: %w", err)
        }
    }*/

    initialDir := filepath.Dir(m.Path)
    if initialDir == "" {
        initialDir, err = os.Getwd()
        if err != nil {
            return nil, false, fmt.Errorf("error getting current working directory: %w", err)
        }
    }

    var filters strings.Builder
    if len(m.FileTypes) == 0 { must.Never() }
    for _, f := range m.FileTypes {
        name, patterns := f[0], strings.Split(f[1], " ")
        filters.WriteString(fmt.Sprintf("%s (%s)", name, strings.Join(patterns, ", ")))
        filters.WriteByte(0) // null terminate
        filters.WriteString(strings.Join(patterns, ";"))
        filters.WriteByte(0) // null terminate
    }

    // double null terminate
    filters.WriteByte(0)

    pOpenFileNameW := openFileNameW{
        lStructSize:       152,
        lpstrFilter:       pwidePermitNul(filters.String()),
        nFilterIndex:      uint32(m.DefaultFileType + 1),
        lpstrFile:         uintptr(unsafe.Pointer(&buf[0])),
        nMaxFile:          bufSize,
        lpstrInitialDir:   pwide(initialDir),
        flags:             flags,
    }

    if mode != 's' {
        ret, _, err := procGetOpenFileName.Call(
            uintptr(unsafe.Pointer(&pOpenFileNameW)),
        )

        if ret == 0 {
            ret, _, err = procCommDlgExtendedError.Call()
            if ret != 0 {
                err = fmt.Errorf("error in comdlg32.dll procedure GetOpenFileNameW: CommDlgExtendedError returned 0x%x: %w", ret, err)
                return nil, false, errSys(err)
            } else {
                return nil, false, nil // closed / cancelled
            }
        }
    } else {
        ret, _, err := procGetSaveFileName.Call(
            uintptr(unsafe.Pointer(&pOpenFileNameW)),
        )

        if ret == 0 {
            ret, _, err = procCommDlgExtendedError.Call()
            if ret != 0 {
                err = fmt.Errorf("error in comdlg32.dll procedure GetSaveFileNameW: CommDlgExtendedError returned 0x%x: %w", ret, err)
                return nil, false, errSys(err)
            } else {
                return nil, false, nil // closed / cancelled
            }
        }
    }

    // multiple items
    if mode == 'm' {
        results := make([]string, 0)
        for {
            result, rest, ok := nullTerminated(buf)
            if len(result) == 0 { break }
            if !ok { break }
            results = append(results, string(utf16.Decode(result)))
            buf = rest
        }

        if len(results) <= 1 {
            // single item
            return results, true, nil
        } else {
            // multiple items so the first is the directory
            // and the remaining are relative.
            root, rest := results[0], results[1:]
            for i := 0; i < len(rest); i++ {
                rest[i] = filepath.Join(root, rest[i])
            }
            return rest, true, nil
        }
    }

    // single item
    if result, _, ok := nullTerminated(buf); !ok {
        err = fmt.Errorf("error in comdlg32.dll procedure CommDlgExtendedError: expected null-terminated result")
        return nil, false, errSys(err)
    } else {
        resultStr := string(utf16.Decode(result))
        return []string{resultStr}, true, nil
    }
}

func (m FilePicker) open() (string, bool, error) {
    if xs, ok, err := m.pick('o'); err != nil {
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

func (m FilePicker) openMultiple() ([]string, bool, error) {
    if xs, ok, err := m.pick('m'); err != nil {
        return nil, false, fmt.Errorf("error opening file picker: %w", err)
    } else {
        return xs, ok, nil
    }
}

func (m FilePicker) save() (string, bool, error) {
    if xs, ok, err := m.pick('s'); err != nil {
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

func (m Message) flags(message string) uint32 {
    flags := uint32(0)           |
        windows.MB_SETFOREGROUND |
        windows.MB_TASKMODAL     |
        windows.MB_TOPMOST

    // right-to-left writing system?
    for _, r := range message {
        prop, _ := bidi.LookupRune(r)
        cls := prop.Class()
        var isRtl bool

        switch cls {
            case bidi.R:   isRtl = true // Strong R-to-L
            case bidi.AL:  isRtl = true // Strong R-to-L
            case bidi.RLO: isRtl = true // Explicit R-to-L
        }
        if isRtl {
            rtlFlags := uint32(0) | windows.MB_RIGHT | windows.MB_RTLREADING
            flags |= rtlFlags
            continue
        }
    }

    return flags
}

func (m Message) ask(message string) (bool, error) {

    // no need to word-wrap message for windows.Messagebox

    flags := uint32(0)           |
        windows.MB_YESNO         |
        windows.MB_ICONWARNING   | // docs say don't use ICONQUESTION
        m.flags(message)

    // TODO add MB_RIGHT | MB_RTLREADING if right-to-left runes detected

    q, _ := windows.MessageBox(0, wide(message), wide(m.Title), flags)
    return q == 6, nil // IDYES
}

func (m Message) raise(message string) error {

    // no need to word-wrap message for windows.Messagebox

    flags := uint32(0)           |
        windows.MB_OK            |
        m.Icon.iconFlag()        |
        m.flags(message)

    windows.MessageBox(0, wide(message), wide(m.Title), flags)

    return nil
}
