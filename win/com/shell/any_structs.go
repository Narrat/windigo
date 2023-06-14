//go:build windows

package shell

import (
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/com/shell/shellco"
)

// [COMDLG_FILTERSPEC] struct.
//
// [COMDLG_FILTERSPEC]: https://docs.microsoft.com/en-us/windows/win32/api/shtypes/ns-shtypes-comdlg_filterspec
type COMDLG_FILTERSPEC struct {
	PszName *uint16
	PszSpec *uint16
}

// [THUMBBUTTON] struct.
//
// [THUMBBUTTON]: https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/ns-shobjidl_core-thumbbutton
type THUMBBUTTON struct {
	DwMask  shellco.THB
	IId     uint32
	IBitmap uint32
	HIcon   win.HICON
	szTip   [260]uint16
	DwFlags shellco.THBF
}

func (tb *THUMBBUTTON) SzTip() string { return win.Str.FromNativeSlice(tb.szTip[:]) }
func (tb *THUMBBUTTON) SetSzTip(val string) {
	copy(tb.szTip[:], win.Str.ToNativeSlice(win.Str.Substr(val, 0, len(tb.szTip)-1)))
}
