package main

import (
    "unsafe"

    "golang.org/x/sys/windows"
)

func osInit() {
    libcomctl32 := windows.MustLoadDLL("comctl32.dll")
    defer libcomctl32.Release()
    commonControlsEx := libcomctl32.MustFindProc("InitCommonControlsEx")

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
        if err != nil { panic(err) }
    }
}
