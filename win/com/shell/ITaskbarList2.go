package shell

import (
	"syscall"
	"unsafe"

	"github.com/rodrigocfd/windigo/internal/util"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/errco"
)

type _ITaskbarList2Vtbl struct {
	_ITaskbarListVtbl
	MarkFullscreenWindow uintptr
}

//------------------------------------------------------------------------------

// 📑 https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nn-shobjidl_core-itaskbarlist2
type ITaskbarList2 struct{ ITaskbarList }

// Constructs a COM object from a pointer to its COM virtual table.
//
// ⚠️ You must defer ITaskbarList2.Release().
//
// Example:
//
//  taskbl2 := shell.NewITaskbarList2(
//      win.CoCreateInstance(
//          shellco.CLSID_TaskbarList, nil,
//          co.CLSCTX_INPROC_SERVER,
//          shellco.IID_ITaskbarList2),
//  )
//  defer taskbl2.Release()
func NewITaskbarList2(ptr win.IUnknownPtr) ITaskbarList2 {
	return ITaskbarList2{
		ITaskbarList: NewITaskbarList(ptr),
	}
}

// 📑 https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-itaskbarlist2-markfullscreenwindow
func (me *ITaskbarList2) MarkFullscreenWindow(hwnd win.HWND, fullScreen bool) {
	ret, _, _ := syscall.Syscall(
		(*_ITaskbarList2Vtbl)(unsafe.Pointer(*me.Ptr())).MarkFullscreenWindow, 3,
		uintptr(unsafe.Pointer(me.Ptr())),
		uintptr(hwnd), util.BoolToUintptr(fullScreen))

	if hr := errco.ERROR(ret); hr != errco.S_OK {
		panic(hr)
	}
}
