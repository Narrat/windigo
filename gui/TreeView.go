/**
 * Part of Wingows - Win32 API layer for Go
 * https://github.com/rodrigocfd/wingows
 * This library is released under the MIT license.
 */

package gui

import (
	"wingows/co"
	"wingows/win"
)

// Native tree view control.
type TreeView struct {
	_ControlNativeBase
}

// Adds a new root item; returns the newly inserted item.
func (me *TreeView) AddRoot(text string) *TreeViewItem {
	return me.Item(0).AddChild(text)
}

// Calls CreateWindowEx(). This is a basic method: no styles are provided by
// default, you must inform all of them. Position and size will be adjusted to
// the current system DPI.
func (me *TreeView) Create(
	parent Window, ctrlId, x, y int32, width, height uint32,
	exStyles co.WS_EX, styles co.WS,
	tvExStyles co.TVS_EX, tvStyles co.TVS) *TreeView {

	x, y, width, height = _globalDpi.multiply(x, y, width, height)

	me._ControlNativeBase.create(exStyles,
		"SysTreeView32", "", styles|co.WS(tvStyles),
		x, y, width, height, parent, ctrlId)

	if tvExStyles != co.TVS_EX(0) {
		me.SetExtendedStyle(tvExStyles, tvExStyles)
	}
	return me
}

// Calls CreateWindowEx(). Tree view control will have lines and buttons.
// Position and size will be adjusted to the current system DPI.
func (me *TreeView) CreateSimple(
	parent Window, ctrlId, x, y int32, width, height uint32) *TreeView {

	return me.Create(parent, ctrlId, x, y, width, height,
		co.WS_EX_CLIENTEDGE,
		co.WS_CHILD|co.WS_GROUP|co.WS_TABSTOP|co.WS_VISIBLE,
		co.TVS_EX(0),
		co.TVS_HASLINES|co.TVS_LINESATROOT|co.TVS_SHOWSELALWAYS|co.TVS_HASBUTTONS)
}

// Retrieves extended styles with TVM_GETEXTENDEDSTYLE.
func (me *TreeView) ExtendedStyle() co.TVS_EX {
	return co.TVS_EX(me.sendTvmMessage(co.TVM_GETEXTENDEDSTYLE, 0, 0))
}

// Sends TVM_GETNEXTITEM with TVGN_ROOT, returns nil if none.
func (me *TreeView) FirstRoot() *TreeViewItem {
	return me.Item(0).NextItem(co.TVGN_ROOT)
}

// Returns nil if none.
func (me *TreeView) FirstVisible() *TreeViewItem {
	return me.Item(0).NextItem(co.TVGN_FIRSTVISIBLE)
}

// Returns the item of the given HTREEITEM.
func (me *TreeView) Item(hTreeItem win.HTREEITEM) *TreeViewItem {
	return &TreeViewItem{
		owner:     me,
		hTreeItem: hTreeItem,
	}
}

// Returns the number of items with TVM_GETCOUNT.
func (me *TreeView) ItemCount() uint32 {
	ret := me.sendTvmMessage(co.TVM_GETCOUNT, 0, 0)
	if ret < 0 {
		panic("TVM_GETCOUNT failed.")
	}
	return uint32(ret)
}

// Sends TVM_GETNEXTITEM, returns nil if none found.
func (me *TreeView) NextItem(flags co.TVGN) *TreeViewItem {
	return me.Item(0).NextItem(flags)
}

func (me *TreeView) Roots() []TreeViewItem {
	return me.Item(0).Children()
}

// Sends TVM_GETNEXTITEM with TVGN_CARET, returns nil if none.
func (me *TreeView) SelectedItem() *TreeViewItem {
	return me.Item(0).NextItem(co.TVGN_CARET)
}

// Sends TVM_SETEXTENDEDSTYLE.
func (me *TreeView) SetExtendedStyle(mask, exStyle co.TVS_EX) *TreeView {
	ret := me.sendTvmMessage(co.TVM_SETEXTENDEDSTYLE,
		win.WPARAM((mask)), win.LPARAM(exStyle))

	lerr := co.ERROR(ret)
	if lerr != co.ERROR_S_OK {
		panic(lerr.Format("TVM_SETEXTENDEDSTYLE failed."))
	}
	return me
}

// Sends WM_SETREDRAW to enable or disable UI updates.
func (me *TreeView) SetRedraw(allowRedraw bool) *TreeView {
	wp := 0
	if allowRedraw {
		wp = 1
	}
	me.hwnd.SendMessage(co.WM_SETREDRAW, win.WPARAM(wp), 0)
	return me
}

// Simple wrapper.
func (me *TreeView) sendTvmMessage(msg co.TVM,
	wParam win.WPARAM, lParam win.LPARAM) uintptr {

	return me.Hwnd().SendMessage(co.WM(msg), wParam, lParam)
}
