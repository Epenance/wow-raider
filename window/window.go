package window

import (
	"syscall"
	"unsafe"
)

var (
	user32, _                = syscall.LoadLibrary("user32.dll")
	findWindowW, _           = syscall.GetProcAddress(user32, "FindWindowW")
	dwmapi, _                = syscall.LoadLibrary("dwmapi.dll")
	dwmGetWindowAttribute, _ = syscall.GetProcAddress(dwmapi, "DwmGetWindowAttribute")
)

const DWMWA_EXTENDED_FRAME_BOUNDS = 9

type HWND uintptr

type RECT struct {
	Left, Top, Right, Bottom int32
}

func FindWindowByTitle(title string) HWND {
	ret, _, _ := syscall.Syscall(
		findWindowW,
		2,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		0,
	)
	return HWND(ret)
}

func DwmGetWindowBounds(hwnd HWND) (*RECT, error) {
	var rect RECT
	r1, _, err := syscall.Syscall6(
		dwmGetWindowAttribute,
		4,
		uintptr(hwnd),
		uintptr(DWMWA_EXTENDED_FRAME_BOUNDS),
		uintptr(unsafe.Pointer(&rect)),
		uintptr(unsafe.Sizeof(rect)),
		0,
		0,
	)
	if r1 != 0 {
		return nil, err
	}
	return &rect, nil
}
