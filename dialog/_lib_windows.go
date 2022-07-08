// This file is the start of an implementation of the newer Common Item Dialog.
// However, due to Go's concurrency model and Windows' concurrency model, this
// hasn't proceeded as it would need the caller to coordinate CoInitializeEx,
// use of runtime.LockOSThread, etc. and we don't want to expose that in the
// simple dialog API. People who need the extra features can use a library
// built for that - golib dialog is a minimal package!

package dialog

import (
    "fmt"
    "unsafe"

    "golang.org/x/sys/windows"

    "github.com/tawesoft/golib/v2/ks"
)

var (
    dllOle32 = windows.NewLazySystemDLL("ole32.dll")

    procClsidFromString = dllOle32.NewProc("CLSIDFromString")
    procCoCreateInstance = dllOle32.NewProc("CoCreateInstance")

    clsidFileOpenDialog = ks.Must(clsidFromString("{DC1C5A9C-E88A-4dde-A5A1-60F82A20AEF7}"))
    clsidFileSaveDialog = ks.Must(clsidFromString("{C0B4E2F3-BA21-4773-8DBA-335EC946EB8B}"))
)

type guid struct {
    Data1   uint32
    Data2   uint16
    Data3   uint16
    Data [8] uint8
}

func errSys(err error) error {
    return fmt.Errorf("windows syscall error: %w", err)
}

func clsidFromString(id string) (*guid, error) {
    if err := procClsidFromString.Find(); err != nil {
        return nil, errSys(err)
    }

    wid := wide(id)
    dest := guid{}
    ret, _, _ := procClsidFromString.Call(
        uintptr(unsafe.Pointer(wid)),
        uintptr(unsafe.Pointer(&dest)),
    )
    if ret != 0 {
        return nil, errSys(fmt.Errorf("CLSIDFromString returned 0x%x", ret))
    }

    return &dest, nil
}

/*

int WINAPI wWinMain(HINSTANCE hInstance, HINSTANCE, PWSTR pCmdLine, int nCmdShow)
{
    HRESULT hr = CoInitializeEx(NULL, COINIT_APARTMENTTHREADED |
        COINIT_DISABLE_OLE1DDE);
    if (SUCCEEDED(hr))
    {
        IFileOpenDialog *pFileOpen;

        // Create the FileOpenDialog object.
        hr = CoCreateInstance(CLSID_FileOpenDialog, NULL, CLSCTX_ALL,
                IID_IFileOpenDialog, reinterpret_cast<void**>(&pFileOpen));

        if (SUCCEEDED(hr))
        {
            // Show the Open dialog box.
            hr = pFileOpen->Show(NULL);

            // Get the file name from the dialog box.
            if (SUCCEEDED(hr))
            {
                IShellItem *pItem;
                hr = pFileOpen->GetResult(&pItem);
                if (SUCCEEDED(hr))
                {
                    PWSTR pszFilePath;
                    hr = pItem->GetDisplayName(SIGDN_FILESYSPATH, &pszFilePath);

                    // Display the file name to the user.
                    if (SUCCEEDED(hr))
                    {
                        MessageBoxW(NULL, pszFilePath, L"File Path", MB_OK);
                        CoTaskMemFree(pszFilePath);
                    }
                    pItem->Release();
                }
            }
            pFileOpen->Release();
        }
        CoUninitialize();
    }
    return 0;
}

*/
