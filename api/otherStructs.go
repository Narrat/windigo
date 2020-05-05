package api

import (
	c "winffi/consts"
)

type ACCEL struct {
	FVirt uint8
	Key   uint16
	Cmd   uint16
}

type CREATESTRUCT struct {
	LpCreateParams uintptr
	HInstance      HINSTANCE
	HMenu          HMENU
	HwndParent     HWND
	Cy, Cx, Y, X   int32
	Style          c.WS
	LpszName       *uint16
	LpszClass      *uint16
	ExStyle        c.WS_EX
}

type LVCOLUMN struct {
	Mask       c.LVCF
	Fmt        int32
	Cx         int32
	PszText    *uint16
	CchTextMax int32
	ISubItem   int32
	IImage     int32
	IOrder     int32
	CxMin      int32
	CxDefault  int32
	CxIdeal    int32
}

type LVITEM struct {
	Mask       c.LVIF
	IItem      int32
	ISubItem   int32
	State      uint32
	StateMask  uint32
	PszText    *uint16
	CchTextMax int32
	IImage     int32
	LParam     uintptr
	IIndent    int32
	IGroupId   int32
	CColumns   uint32
	PuColumns  *uint32
	PiColFmt   *int32
	IGroup     int32
}

type MONITORINFOEX struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	Flags     uint32
	SzDevice  [32]uint16 // CCHDEVICENAME
}

type NMHDR struct {
	HWndFrom HWND
	IdFrom   uint32
	Code     uint32
}

type NONCLIENTMETRICS struct {
	CbSize             uint32
	IBorderWidth       int32
	IScrollWidth       int32
	IScrollHeight      int32
	ICaptionWidth      int32
	ICaptionHeight     int32
	LfCaptionFont      LOGFONT
	ISmCaptionWidth    int32
	ISmCaptionHeight   int32
	LfSmCaptionFont    LOGFONT
	IMenuWidth         int32
	IMenuHeight        int32
	LfMenuFont         LOGFONT
	LfStatusFont       LOGFONT
	LfMessageFont      LOGFONT
	IPaddedBorderWidth int32
}

type POINT struct {
	X, Y int32
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

type SIZE struct {
	Cx, Cy int32
}
