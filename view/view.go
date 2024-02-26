package view

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type View interface {
	Show()
	DrawControls()
}

type BoundView interface {
	View
	SetRefreshTagsFunc(onRefreshTags func(boundType string, tag []string))
}

func newForm(title string, width, height int32) *vcl.TForm {
	form := vcl.NewForm(nil)
	form.SetCaption(title)
	form.SetWidth(width)
	form.SetHeight(height)
	form.SetPosition(types.PoScreenCenter) // center display
	return form
}

func newMenuItem(parent vcl.IComponent, index *int, name, title string, clickEvent vcl.TNotifyEvent) *vcl.TMenuItem {
	menuItem := vcl.NewMenuItem(parent)
	menuItem.SetTag(*index)
	menuItem.SetName(name)
	menuItem.SetCaption(title)
	menuItem.SetOnClick(clickEvent)
	*index++
	return menuItem
}

func newScrollBox(parent vcl.IWinControl, width, height int32, align types.TAlign) *vcl.TScrollBox {
	scrollBox := vcl.NewScrollBox(parent)
	scrollBox.SetParent(parent)
	scrollBox.SetWidth(width)
	scrollBox.SetHeight(height)
	scrollBox.SetAlign(align)
	return scrollBox
}

// owner means when the owner is destroyed, all the controls it owns are automatically destroyed as well
// setParent determines where the control is displayed
func newPanel(parent vcl.IWinControl, rect Rect, align types.TAlign) *vcl.TPanel {
	panel := vcl.NewPanel(parent)
	panel.SetParent(parent)
	panel.SetCaption("")
	panel.SetBounds(rect.CenterX, rect.CenterY, rect.Width, rect.Height)
	panel.SetAlign(align)
	return panel
}

func newLabel(parent vcl.IWinControl, text string, rect Rect) *vcl.TLabel {
	label := vcl.NewLabel(parent)
	label.SetParent(parent)
	label.SetCaption(text)
	label.SetBounds(rect.CenterX, rect.CenterY, rect.Width, rect.Height)
	label.SetAutoSize(false)
	return label
}

func newButton(parent vcl.IWinControl, text string, rect Rect, clickEvent func(sender vcl.IObject)) {
	button := vcl.NewButton(parent)
	button.SetParent(parent)
	button.SetCaption(text)
	button.SetBounds(rect.CenterX, rect.CenterY, rect.Width, rect.Height)
	button.SetOnClick(clickEvent)
}

func newEdit(parent vcl.IWinControl, text string, rect Rect) *vcl.TEdit {
	edit := vcl.NewEdit(parent)
	edit.SetParent(parent)
	edit.SetText(text)
	edit.SetBounds(rect.CenterX, rect.CenterY, rect.Width, rect.Height)
	return edit
}

func newComboBox(parent vcl.IWinControl, text string, rect Rect, items []string) *vcl.TComboBox {
	comboBox := vcl.NewComboBox(parent)
	comboBox.SetParent(parent)
	comboBox.SetText(text)
	comboBox.SetBounds(rect.CenterX, rect.CenterY, rect.Width, rect.Height)
	for _, item := range items {
		comboBox.AddItem(item, nil)
	}
	return comboBox
}

// Deprecated
// func newStringGrid(parent vcl.IWinControl, name string, width, height int32, align types.TAlign) *vcl.TStringGrid {
// 	stringGrid := vcl.NewStringGrid(parent)
// 	stringGrid.SetParent(parent)
// 	stringGrid.SetName(name)
// 	stringGrid.SetWidth(width)
// 	stringGrid.SetHeight(height)
// 	stringGrid.SetAlign(align)
// 	stringGrid.SetColCount(1)
// 	stringGrid.SetRowCount(1)
// 	stringGrid.SetColumnClickSorts(false)
// 	stringGrid.SetOptions(stringGrid.Options().Include(
// 		types.GoAlwaysShowEditor,
// 		types.GoEditing,   // cell data can be edit
// 		types.GoTabs,      // Tab can auto next cell
// 		types.GoRowSizing, //Row height can resize by mouse
// 		types.GoColSizing, //Col width can resize by mouse
// 	))
// 	return stringGrid
// }
