//go:build windows

package ui

import (
	"unsafe"

	"github.com/rodrigocfd/windigo/internal/util"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
)

// Native [toolbar] control.
//
// [toolbar]: https://learn.microsoft.com/en-us/windows/win32/controls/toolbar-controls-overview
type Toolbar interface {
	AnyNativeControl
	implToolbar() // prevent public implementation

	// Exposes all the [Toolbar notifications] the can be handled.
	//
	// Panics if called after the control was created.
	//
	// [Toolbar notifications]: https://learn.microsoft.com/en-us/windows/win32/controls/bumper-toolbar-control-reference-notifications
	On() *_ToolbarEvents

	AutoSize()                                                   // Sends a TB_AUTOSIZE message to resize the toolbar.
	ExtendedStyle() co.TBSTYLE_EX                                // Retrieves the extended style flags.
	Buttons() *_ToolbarButtons                                   // Button methods.
	SetExtendedStyle(doSet bool, styles co.TBSTYLE_EX)           // Sets or unsets extended style flags.
	SetImageList(index int, himgl win.HIMAGELIST) win.HIMAGELIST // Sets the nth image list for the control.
	SetParent(hNewParent win.HWND) win.HWND                      // Sets the window to which the toolbar control sends notification messages.
}

//------------------------------------------------------------------------------

type _Toolbar struct {
	_NativeControlBase
	events  _ToolbarEvents
	buttons _ToolbarButtons
}

// Creates a new Toolbar. Call ui.ToolbarOpts() to define the options to be
// passed to the underlying CreateWindowEx().
//
// # Example
//
//	var owner ui.AnyParent // initialized somewhere
//
//	myBar := ui.NewToolbar(
//		owner,
//		ui.ToolbarOpts().
//			CtrlStyles(co.TBSTYLE_BUTTON|co.TBSTYLE_LIST|co.TBSTYLE_FLAT),
//		),
//	)
func NewToolbar(parent AnyParent, opts *_ToolbarO) Toolbar {
	if opts == nil {
		opts = ToolbarOpts()
	}
	opts.lateDefaults()

	me := &_Toolbar{}
	me._NativeControlBase.new(parent, opts.ctrlId)
	me.events.new(&me._NativeControlBase)
	me.buttons.new(me)

	parent.internalOn().addMsgNoRet(_CreateOrInitDialog(parent), func(_ wm.Any) {
		me._NativeControlBase.createWindow(opts.wndExStyles,
			win.ClassNameStr("ToolbarWindow32"), win.StrOptNone(),
			opts.wndStyles|co.WS(opts.ctrlStyles),
			win.POINT{}, win.SIZE{}, win.HMENU(opts.ctrlId))

		me.Hwnd().SendMessage(co.TB_BUTTONSTRUCTSIZE,
			win.WPARAM(unsafe.Sizeof(win.TBBUTTON{})), 0)
		me.Hwnd().SendMessage(co.CCM_SETVERSION, 5, 0)

		if opts.ctrlExStyles != co.TBSTYLE_EX_NONE {
			me.SetExtendedStyle(true, opts.ctrlExStyles)
		}
	})

	return me
}

// Implements Toolbar.
func (*_Toolbar) implToolbar() {}

func (me *_Toolbar) On() *_ToolbarEvents {
	if me.Hwnd() != 0 {
		panic("Cannot add event handling after the Toolbar is created.")
	}
	return &me.events
}

func (me *_Toolbar) Buttons() *_ToolbarButtons {
	return &me.buttons
}

func (me *_Toolbar) AutoSize() {
	me.Hwnd().SendMessage(co.TB_AUTOSIZE, 0, 0)
}

func (me *_Toolbar) ExtendedStyle() co.TBSTYLE_EX {
	return co.TBSTYLE_EX(
		me.Hwnd().SendMessage(co.TB_GETEXTENDEDSTYLE, 0, 0),
	)
}

func (me *_Toolbar) SetExtendedStyle(doSet bool, styles co.TBSTYLE_EX) {
	curStyles := me.ExtendedStyle()
	newStyles := util.Iif(doSet,
		curStyles|styles, curStyles & ^styles).(co.TBSTYLE_EX)

	me.Hwnd().SendMessage(co.TB_SETEXTENDEDSTYLE, 0, win.LPARAM(newStyles))
}

func (me *_Toolbar) SetImageList(
	index int, himgl win.HIMAGELIST) win.HIMAGELIST {

	return win.HIMAGELIST(
		me.Hwnd().SendMessage(co.TB_SETIMAGELIST,
			win.WPARAM(index), win.LPARAM(himgl)),
	)
}

func (me *_Toolbar) SetParent(hNewParent win.HWND) win.HWND {
	return win.HWND(me.Hwnd().
		SendMessage(co.TB_SETPARENT, win.WPARAM(hNewParent), 0))
}

//------------------------------------------------------------------------------

type _ToolbarO struct {
	ctrlId int

	ctrlStyles   co.TBSTYLE
	ctrlExStyles co.TBSTYLE_EX
	wndStyles    co.WS
	wndExStyles  co.WS_EX
}

// Control ID.
//
// Defaults to an auto-generated ID.
func (o *_ToolbarO) CtrlId(i int) *_ToolbarO { o.ctrlId = i; return o }

// Toolbar control styles, passed to CreateWindowEx().
//
// Defaults to TBSTYLE_BUTTON | TBSTYLE_FLAT
func (o *_ToolbarO) CtrlStyles(s co.TBSTYLE) *_ToolbarO { o.ctrlStyles = s; return o }

// Toolbar extended control styles, passed to CreateWindowEx().
//
// Defaults to TBSTYLE_EX_NONE.
func (o *_ToolbarO) CtrlExStyles(s co.TBSTYLE_EX) *_ToolbarO { o.ctrlExStyles = s; return o }

// Window styles, passed to CreateWindowEx().
//
// Defaults to co.WS_CHILD | co.WS_VISIBLE.
func (o *_ToolbarO) WndStyles(s co.WS) *_ToolbarO { o.wndStyles = s; return o }

// Extended window styles, passed to CreateWindowEx().
//
// Defaults to WS_EX_NONE.
func (o *_ToolbarO) WndExStyles(s co.WS_EX) *_ToolbarO { o.wndExStyles = s; return o }

func (o *_ToolbarO) lateDefaults() {
	if o.ctrlId == 0 {
		o.ctrlId = _NextCtrlId()
	}
}

// Options for NewToolbar().
func ToolbarOpts() *_ToolbarO {
	return &_ToolbarO{
		ctrlStyles: co.TBSTYLE_BUTTON | co.TBSTYLE_FLAT,
		wndStyles:  co.WS_CHILD | co.WS_VISIBLE,
	}
}

//------------------------------------------------------------------------------

// Toolbar control notifications.
type _ToolbarEvents struct {
	ctrlId int
	events *_EventsWmNfy
}

func (me *_ToolbarEvents) new(ctrl *_NativeControlBase) {
	me.ctrlId = ctrl.CtrlId()
	me.events = ctrl.Parent().On()
}

// [TBN_BEGINADJUST] message handler.
//
// [TBN_BEGINADJUST]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-beginadjust
func (me *_ToolbarEvents) TbnBeginAdjust(userFunc func()) {
	me.events.addNfyZero(me.ctrlId, co.TBN_BEGINADJUST, func(_ unsafe.Pointer) {
		userFunc()
	})
}

// [TBN_BEGINDRAG] message handler.
//
// [TBN_BEGINDRAG]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-begindrag
func (me *_ToolbarEvents) TbnBeginDrag(userFunc func(p *win.NMTOOLBAR)) {
	me.events.addNfyZero(me.ctrlId, co.TBN_BEGINDRAG, func(p unsafe.Pointer) {
		userFunc((*win.NMTOOLBAR)(p))
	})
}

// [TBN_CUSTHELP] message handler.
//
// [TBN_CUSTHELP]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-custhelp
func (me *_ToolbarEvents) TbnCustHelp(userFunc func()) {
	me.events.addNfyZero(me.ctrlId, co.TBN_CUSTHELP, func(_ unsafe.Pointer) {
		userFunc()
	})
}

// [TBN_DELETINGBUTTON] message handler.
//
// [TBN_DELETINGBUTTON]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-deletingbutton
func (me *_ToolbarEvents) TbnDeletingButton(userFunc func(p *win.NMTOOLBAR)) {
	me.events.addNfyZero(me.ctrlId, co.TBN_DELETINGBUTTON, func(p unsafe.Pointer) {
		userFunc((*win.NMTOOLBAR)(p))
	})
}

// [TBN_DRAGOUT] message handler.
//
// [TBN_DRAGOUT]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-dragout
func (me *_ToolbarEvents) TbnDragOut(userFunc func(p *win.NMTOOLBAR)) {
	me.events.addNfyZero(me.ctrlId, co.TBN_DRAGOUT, func(p unsafe.Pointer) {
		userFunc((*win.NMTOOLBAR)(p))
	})
}

// [TBN_DRAGOVER] message handler.
//
// [TBN_DRAGOVER]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-dragover
func (me *_ToolbarEvents) TbnDragOver(userFunc func(p *win.NMTBHOTITEM) bool) {
	me.events.addNfyRet(me.ctrlId, co.TBN_DRAGOVER, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMTBHOTITEM)(p)))
	})
}

// [TBN_DROPDOWN] message handler.
//
// [TBN_DROPDOWN]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-dropdown
func (me *_ToolbarEvents) TbnDropDown(userFunc func(p *win.NMTOOLBAR) co.TBDDRET) {
	me.events.addNfyRet(me.ctrlId, co.TBN_DROPDOWN, func(p unsafe.Pointer) uintptr {
		return uintptr(userFunc((*win.NMTOOLBAR)(p)))
	})
}

// [TBN_DUPACCELERATOR] message handler.
//
// [TBN_DUPACCELERATOR]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-dupaccelerator
func (me *_ToolbarEvents) TbnDupAccelerator(userFunc func(p *win.NMTBDUPACCELERATOR) bool) {
	me.events.addNfyRet(me.ctrlId, co.TBN_DUPACCELERATOR, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMTBDUPACCELERATOR)(p)))
	})
}

// [TBN_ENDADJUST] message handler.
//
// [TBN_ENDADJUST]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-endadjust
func (me *_ToolbarEvents) TbnEndAdjust(userFunc func()) {
	me.events.addNfyZero(me.ctrlId, co.TBN_ENDADJUST, func(_ unsafe.Pointer) {
		userFunc()
	})
}

// [TBN_ENDDRAG] message handler.
//
// [TBN_ENDDRAG]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-enddrag
func (me *_ToolbarEvents) TbnEndDrag(userFunc func(p *win.NMTOOLBAR)) {
	me.events.addNfyZero(me.ctrlId, co.TBN_ENDDRAG, func(p unsafe.Pointer) {
		userFunc((*win.NMTOOLBAR)(p))
	})
}

// [TBN_GETBUTTONINFO] message handler.
//
// [TBN_GETBUTTONINFO]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-getbuttoninfo
func (me *_ToolbarEvents) TbnGetButtonInfo(userFunc func(p *win.NMTOOLBAR) bool) {
	me.events.addNfyRet(me.ctrlId, co.TBN_GETBUTTONINFO, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMTOOLBAR)(p)))
	})
}

// [TBN_GETDISPINFO] message handler.
//
// [TBN_GETDISPINFO]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-getdispinfo
func (me *_ToolbarEvents) TbnGetDispInfo(userFunc func(p *win.NMTBDISPINFO)) {
	me.events.addNfyZero(me.ctrlId, co.TBN_GETDISPINFO, func(p unsafe.Pointer) {
		userFunc((*win.NMTBDISPINFO)(p))
	})
}

// [TBN_GETINFOTIP] message handler.
//
// [TBN_GETINFOTIP]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-getinfotip
func (me *_ToolbarEvents) TbnGetInfoTip(userFunc func(p *win.NMTBGETINFOTIP)) {
	me.events.addNfyZero(me.ctrlId, co.TBN_GETINFOTIP, func(p unsafe.Pointer) {
		userFunc((*win.NMTBGETINFOTIP)(p))
	})
}

// [TBN_GETOBJECT] message handler.
//
// [TBN_GETOBJECT]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-getobject
func (me *_ToolbarEvents) TbnGetObject(userFunc func(p *win.NMOBJECTNOTIFY)) {
	me.events.addNfyZero(me.ctrlId, co.TBN_GETOBJECT, func(p unsafe.Pointer) {
		userFunc((*win.NMOBJECTNOTIFY)(p))
	})
}

// [TBN_HOTITEMCHANGE] message handler.
//
// [TBN_HOTITEMCHANGE]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-hotitemchange
func (me *_ToolbarEvents) TbnHotItemChange(userFunc func(*win.NMTBHOTITEM) int) {
	me.events.addNfyRet(me.ctrlId, co.TBN_HOTITEMCHANGE, func(p unsafe.Pointer) uintptr {
		return uintptr(userFunc((*win.NMTBHOTITEM)(p)))
	})
}

// [TBN_INITCUSTOMIZE] message handler.
//
// [TBN_INITCUSTOMIZE]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-initcustomize
func (me *_ToolbarEvents) TbnInitCustomize(userFunc func() co.TBNRF) {
	me.events.addNfyRet(me.ctrlId, co.TBN_INITCUSTOMIZE, func(p unsafe.Pointer) uintptr {
		return uintptr(userFunc())
	})
}

// [TBN_MAPACCELERATOR] message handler.
//
// [TBN_MAPACCELERATOR]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-mapaccelerator
func (me *_ToolbarEvents) TbnMapAccelerator(userFunc func(p *win.NMCHAR) bool) {
	me.events.addNfyRet(me.ctrlId, co.TBN_MAPACCELERATOR, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMCHAR)(p)))
	})
}

// [TBN_QUERYDELETE] message handler.
//
// [TBN_QUERYDELETE]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-querydelete
func (me *_ToolbarEvents) TbnQueryDelete(userFunc func(p *win.NMTOOLBAR) bool) {
	me.events.addNfyRet(me.ctrlId, co.TBN_QUERYDELETE, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMTOOLBAR)(p)))
	})
}

// [TBN_QUERYINSERT] message handler.
//
// [TBN_QUERYINSERT]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-queryinsert
func (me *_ToolbarEvents) TbnQueryInsert(userFunc func(p *win.NMTOOLBAR) bool) {
	me.events.addNfyRet(me.ctrlId, co.TBN_QUERYINSERT, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMTOOLBAR)(p)))
	})
}

// [TBN_RESET] message handler.
//
// [TBN_RESET]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-reset
func (me *_ToolbarEvents) TbnReset(userFunc func() co.TBNRF) {
	me.events.addNfyRet(me.ctrlId, co.TBN_RESET, func(p unsafe.Pointer) uintptr {
		return uintptr(userFunc())
	})
}

// [TBN_RESTORE] message handler.
//
// [TBN_RESTORE]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-restore
func (me *_ToolbarEvents) TbnRestore(userFunc func(p *win.NMTBRESTORE) int) {
	me.events.addNfyRet(me.ctrlId, co.TBN_RESTORE, func(p unsafe.Pointer) uintptr {
		return uintptr(userFunc((*win.NMTBRESTORE)(p)))
	})
}

// [TBN_SAVE] message handler.
//
// [TBN_SAVE]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-save
func (me *_ToolbarEvents) TbnSave(userFunc func(p *win.NMTBSAVE)) {
	me.events.addNfyZero(me.ctrlId, co.TBN_SAVE, func(p unsafe.Pointer) {
		userFunc((*win.NMTBSAVE)(p))
	})
}

// [TBN_TOOLBARCHANGE] message handler.
//
// [TBN_TOOLBARCHANGE]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-toolbarchange
func (me *_ToolbarEvents) TbnToolbarChange(userFunc func()) {
	me.events.addNfyZero(me.ctrlId, co.TBN_TOOLBARCHANGE, func(p unsafe.Pointer) {
		userFunc()
	})
}

// [TBN_WRAPACCELERATOR] message handler.
//
// [TBN_WRAPACCELERATOR]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-wrapaccelerator
func (me *_ToolbarEvents) TbnWrapAccelerator(userFunc func(p *win.NMTBWRAPACCELERATOR) bool) {
	me.events.addNfyRet(me.ctrlId, co.TBN_WRAPACCELERATOR, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMTBWRAPACCELERATOR)(p)))
	})
}

// [TBN_WRAPHOTITEM] message handler.
//
// [TBN_WRAPHOTITEM]: https://learn.microsoft.com/en-us/windows/win32/controls/tbn-wraphotitem
func (me *_ToolbarEvents) TbnWrapHotItem(userFunc func(p *win.NMTBWRAPHOTITEM) bool) {
	me.events.addNfyRet(me.ctrlId, co.TBN_WRAPHOTITEM, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMTBWRAPHOTITEM)(p)))
	})
}

// [NM_CHAR] message handler.
//
// [NM_CHAR]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-char-toolbar
func (me *_ToolbarEvents) NmChar(userFunc func(p *win.NMCHAR) bool) {
	me.events.addNfyRet(me.ctrlId, co.NM_CHAR, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMCHAR)(p)))
	})
}

// [NM_CLICK] message handler.
//
// [NM_CLICK]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-click-toolbar
func (me *_ToolbarEvents) NmClick(userFunc func(p *win.NMMOUSE) bool) {
	me.events.addNfyRet(me.ctrlId, co.NM_CLICK, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMMOUSE)(p)))
	})
}

// [NM_CUSTOMDRAW] message handler.
//
// [NM_CUSTOMDRAW]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-customdraw-toolbar
func (me *_ToolbarEvents) NmCustomDraw(userFunc func(p *win.NMTBCUSTOMDRAW) co.CDRF) {
	me.events.addNfyRet(me.ctrlId, co.NM_CUSTOMDRAW, func(p unsafe.Pointer) uintptr {
		return uintptr(userFunc((*win.NMTBCUSTOMDRAW)(p)))
	})
}

// [NM_DBLCLK] message handler.
//
// [NM_DBLCLK]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-dblclk-toolbar
func (me *_ToolbarEvents) NmDblClk(userFunc func(p *win.NMMOUSE) bool) {
	me.events.addNfyRet(me.ctrlId, co.NM_DBLCLK, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMMOUSE)(p)))
	})
}

// [NM_KEYDOWN] message handler.
//
// [NM_KEYDOWN]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-keydown-toolbar
func (me *_ToolbarEvents) NmKeyDown(userFunc func(p *win.NMKEY) int) {
	me.events.addNfyRet(me.ctrlId, co.NM_KEYDOWN, func(p unsafe.Pointer) uintptr {
		return uintptr(userFunc((*win.NMKEY)(p)))
	})
}

// [NM_LDOWN] message handler.
//
// [NM_LDOWN]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-ldown-toolbar
func (me *_ToolbarEvents) NmLDown(userFunc func(p *win.NMMOUSE) bool) {
	me.events.addNfyRet(me.ctrlId, co.NM_LDOWN, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMMOUSE)(p)))
	})
}

// [NM_RCLICK] message handler.
//
// [NM_RCLICK]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-rclick-toolbar
func (me *_ToolbarEvents) NmRClick(userFunc func(p *win.NMMOUSE) bool) {
	me.events.addNfyRet(me.ctrlId, co.NM_RCLICK, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMMOUSE)(p)))
	})
}

// [NM_RDBLCLK] message handler.
//
// [NM_RDBLCLK]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-rdblclk-toolbar
func (me *_ToolbarEvents) NmRDblClk(userFunc func(p *win.NMMOUSE) bool) {
	me.events.addNfyRet(me.ctrlId, co.NM_RDBLCLK, func(p unsafe.Pointer) uintptr {
		return util.BoolToUintptr(userFunc((*win.NMMOUSE)(p)))
	})
}

// [NM_RELEASEDCAPTURE] message handler.
//
// [NM_RELEASEDCAPTURE]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-releasedcapture-list-view-
func (me *_ToolbarEvents) NmReleasedCapture(userFunc func()) {
	me.events.addNfyZero(me.ctrlId, co.NM_RELEASEDCAPTURE, func(_ unsafe.Pointer) {
		userFunc()
	})
}

// [NM_TOOLTIPSCREATED] message handler.
//
// [NM_TOOLTIPSCREATED]: https://learn.microsoft.com/en-us/windows/win32/controls/nm-tooltipscreated-toolbar-
func (me *_ToolbarEvents) NmTooltipsCreated(userFunc func(p *win.NMTOOLTIPSCREATED)) {
	me.events.addNfyZero(me.ctrlId, co.NM_TOOLTIPSCREATED, func(p unsafe.Pointer) {
		userFunc((*win.NMTOOLTIPSCREATED)(p))
	})
}
