package shell

import (
	"syscall"
	"unsafe"

	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/com/shell/shellco"
	"github.com/rodrigocfd/windigo/win/com/shell/shellvt"
	"github.com/rodrigocfd/windigo/win/errco"
)

// 📑 https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nn-shobjidl_core-ishellitemarray
type IShellItemArray struct{ win.IUnknown }

// Constructs a COM object from a pointer to its COM virtual table.
//
// ⚠️ You must defer IShellItemArray.Release().
func NewIShellItemArray(ptr win.IUnknownPtr) IShellItemArray {
	return IShellItemArray{
		IUnknown: win.NewIUnknown(ptr),
	}
}

// Prefer using IShellItemArray.GetDisplayNames(), which directly retrieves all
// full paths at once.
//
// ⚠️ You must defer IShellItem.Release().
//
// 📑 https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-ishellitemarray-getitemat
func (me *IShellItemArray) GetItemAt(index int) IShellItem {
	var ppvQueried win.IUnknownPtr
	ret, _, _ := syscall.Syscall(
		(*shellvt.IShellItemArray)(unsafe.Pointer(*me.Ptr())).GetItemAt, 3,
		uintptr(unsafe.Pointer(me.Ptr())),
		uintptr(index), uintptr(unsafe.Pointer(&ppvQueried)))

	if hr := errco.ERROR(ret); hr == errco.S_OK {
		return NewIShellItem(ppvQueried)
	} else {
		panic(hr)
	}
}

// 📑 https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-ishellitemarray-getcount
func (me *IShellItemArray) GetCount() int {
	var count uint32
	ret, _, _ := syscall.Syscall(
		(*shellvt.IShellItemArray)(unsafe.Pointer(*me.Ptr())).GetCount, 3,
		uintptr(unsafe.Pointer(me.Ptr())),
		uintptr(unsafe.Pointer(&count)), 0)

	if hr := errco.ERROR(ret); hr == errco.S_OK {
		return int(count)
	} else {
		panic(hr)
	}
}

// Calls GetDisplayName() on each IShellItem, retrieving the names as strings.
//
// Example:
//
//  var shia shell.IShellItemArray // initialized somewhere
//
//  fullPaths := shia.GetDisplayNames(shellco.SIGDN_FILESYSPATH)
func (me *IShellItemArray) GetDisplayNames(sigdnName shellco.SIGDN) []string {
	count := me.GetCount()
	files := make([]string, 0, count)

	for i := 0; i < count; i++ {
		shellItem := me.GetItemAt(i)
		defer shellItem.Release() // will pile up at the end of the function, but that's fine
		files = append(files, shellItem.GetDisplayName(sigdnName))
	}

	return files
}
