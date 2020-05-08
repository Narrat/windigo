package ui

import (
	"gowinui/api"
	c "gowinui/consts"
	"unsafe"
)

type windowControlSetup struct {
	ClassName  string
	ClassStyle c.CS
	HCursor    api.HCURSOR
	HBrushBg   api.HBRUSH

	Style   c.WS
	ExStyle c.WS_EX // for border, use WS_EX_CLIENTEDGE
}

func makeWindowControlSetup() windowControlSetup {
	return windowControlSetup{
		ClassStyle: c.CS_DBLCLKS,

		Style:   c.WS_CHILD | c.WS_VISIBLE | c.WS_CLIPCHILDREN | c.WS_CLIPSIBLINGS,
		ExStyle: c.WS_EX(0),
	}
}

func (me *windowControlSetup) genWndClassEx(
	hInst api.HINSTANCE) *api.WNDCLASSEX {

	wcx := api.WNDCLASSEX{}

	wcx.CbSize = uint32(unsafe.Sizeof(wcx))
	wcx.HInstance = hInst
	wcx.LpszClassName = api.StrToUtf16PtrBlankIsNil(me.ClassName)
	wcx.Style = me.ClassStyle

	if me.HCursor != 0 {
		wcx.HCursor = me.HCursor
	} else {
		wcx.HCursor = api.HINSTANCE(0).LoadCursor(c.IDC_ARROW)
	}

	if me.HBrushBg != 0 {
		wcx.HbrBackground = me.HBrushBg
	} else {
		wcx.HbrBackground = api.NewBrushFromSysColor(c.COLOR_WINDOW)
	}

	return &wcx
}