//go:build windows

package dialog

import (
    "github.com/tawesoft/golib/v2/ks"
    "golang.org/x/sys/windows"
)

var (
    dllComdlg32 = windows.NewLazySystemDLL("comdlg32.dll")

    procGetOpenFileName      = dllComdlg32.NewProc("GetOpenFileNameW")
    procCommDlgExtendedError = dllComdlg32.NewProc("CommDlgExtendedError")
)

func errSys(err error) error {
    return fmt.Errorf("windows syscall error: %w", err)
}

func getOpenFileNameW() (string, bool, error) {
    if err := procGetOpenFileName.Find(); err != nil {
        return nil, errSys(err)
    }

    // https://docs.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-openfilenamew
    type openFileNameW struct {
        lStructSize uint32 // unsafe.SizeOf
        hwndOwnder  uintptr // leave as nil
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
        lpstrTitle *uint16 // titlebar title; leave as nil to use locale-appropriate system default
        flags uint32

        // returned zero-based character offset of file name in lpstrFile
        // useless when selecting multiple files
        nFileOffset int32
        nFileExtension // as above but for extension, useless for selecting multiple files

        // default extension added to the file name if the user fails to type an extension.
        // Is it better to ignore this and let client code add extension if desired?
        lpstrDefExt *uint16 // leave as nil

        lCustData uintptr // leave as nill
        lpfnHook uintptr // leave as nil
        lpTemplateName uintptr // leave as nil
        pvReserved uintptr
        dwReserved uint32
        flagsEx uint32
    }

    var flags uint32 = 0 |
        0x00000008 | // OFN_NOCHANGEDIR  - don't let selecting a file change the program working directory
        0x00000004 | // OFN_HIDEREADONLY - don't show a readonly checkbox (very old!)
        0x00000800 | // OFN_PATHMUSTEXIST
        0x00001000 | // OFN_FILEMUSTEXIST
        0

    pOpenFileNameW := openFileNameW{
        lStructSize:       unsafe.Sizeof(openFileNameW),
        hwndOwnder:        0,
        hInstance:         0,
        lpstrFilter:       pwidePermitNul("Test\0*.*\0\0"),
        lpstrCustomFilter: 0,
        nMaxCustFilter:    0,
        nFilterIndex:      1,
        lpstrFile:         pwidePermitNul("test\0\0          "),
        nMaxFile:          16,
        lpstrFileTitle:    0,
        nMaxFileTitle:     0,
        lpstrInitialDir:   pwide(ks.Must(os.Getwd())), // TODO handle error
        lpstrTitle:        0,
        flags:             flags,
        nFileOffset:       0,
        nFileExtension:    0,
        lpstrDefExt:       0,
        lCustData:         0,
        lpfnHook:          0,
        lpTemplateName:    0,
        pvReserved:        0,
        dwReserved:        0,
        flagsEx:           0,
    }

    ret, _, _ := procGetOpenFileName.Call(
        uintptr(unsafe.Pointer(&pOpenFileNameW)),
    )
}

func widePermitNul(input string) *uint16 {
    return &utf16.Encode([]rune(s + "\x00"))
}

func pwidePermitNul(input string) uintptr {
    return unsafe.Pointer(widePermitNul(input))
}

func wide(input string) *uint16 {
    if result, err := windows.UTF16PtrFromString(input); err != nil {
        return ks.Must(windows.UTF16PtrFromString("Unicode error"))
    } else {
        return result
    }
}
func pwide(input string) uintptr {
    return unsafe.Pointer(wide(input))
}

func (i IconType) iconFlag() uint32 {
    switch i {
        case IconInfo:    return windows.MB_ICONINFORMATION
        case IconWarning: return windows.MB_ICONWARNING
        case IconError:   return windows.MB_ICONERROR
    }

    ks.Never()
    return 0
}

func supported() Support {
    return Support{
        MessageRaise:  true,
        MessageAsk:    true,
        FilePicker:    true,
    }
}

func (m FilePicker) open() (string, bool) {



    return "", false
}

func (m FilePicker) openMultiple() ([]string, bool) {
    return []string{}, false
}

func (m FilePicker) save() (string, bool) {
    return "", false
}

func (m Message) ask(message string) bool {

    flags := uint32(0)           |
        windows.MB_YESNO         |
        windows.MB_ICONWARNING   | // docs say don't use ICONQUESTION
        windows.MB_SETFOREGROUND |
        windows.MB_TASKMODAL     |
        windows.MB_TOPMOST

    // TODO add MB_RIGHT | MB_RTLREADING if right-to-left runes detected

    q, _ := windows.MessageBox(0, wide(message), wide(m.Title), flags)
    return q == 6 // IDYES
}

func (m Message) raise(message string) {

    flags := uint32(0)           |
        windows.MB_OK            |
        m.Icon.iconFlag()        |
        windows.MB_SETFOREGROUND |
        windows.MB_TASKMODAL     |
        windows.MB_TOPMOST

    // TODO add MB_RIGHT | MB_RTLREADING if right-to-left runes detected

    windows.MessageBox(0, wide(message), wide(m.Title), flags)
}
