package view

import (
	"strconv"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type Rect struct {
	CenterX int32
	CenterY int32
	Width   int32
	Height  int32
}

type IChild interface {
	GetFormTitle() string
	GetColRects() []Rect
	GetHeadText() []string
	GetBounds() [][]string
	NewEditControls(panel *vcl.TPanel, rects []Rect)
	NewDataControls(panel *vcl.TPanel, data []string, rects []Rect)
	SetDefaultText()
	ClearEditText()
	AddBound() []string
	CopyBound(input []string)
	UpdateBound(rowIndex int, input []string)
	DeleteBound(rowIndex int)
	UpBound(rowIndex int)
	DownBound(rowIndex int)
	GetTags() (string, []string)
	SaveConfig() bool
}

type boundView struct {
	form *vcl.TForm
	IChild
	// Data row
	scrollBox       *vcl.TScrollBox
	dataPanels      []*vcl.TPanel
	sendRefreshTags func(boundType string, tags []string)
}

func NewBoundView(child IChild) BoundView {
	v := &boundView{
		form:   newForm("", 0, 0),
		IChild: child,
	}
	// v.form.SetBorderStyle(types.BsSingle) // disable maximize button
	var width int32 = 0
	for _, rect := range v.IChild.GetColRects() {
		width += rect.Width
	}
	v.form.SetCaption(v.IChild.GetFormTitle())
	v.form.SetWidth(width)
	v.form.SetOnClose(v.onFormClose)
	return v
}

// Show implements View.
func (v *boundView) Show() {
	v.drawDataRows()
	v.form.Show()
}

// DrawControls implements View.
func (v *boundView) DrawControls() {
	h2 := v.drawEditRow()
	h1 := v.drawHeadRow()
	h3 := v.drawDataBox()
	v.form.SetHeight(h1 + h2 + h3)
}

func (v *boundView) drawHeadRow() int32 {
	rects := v.IChild.GetColRects()
	panel := newPanel(v.form, Rect{CenterX: 0, CenterY: 0, Width: v.form.Width(), Height: rects[0].Height}, types.AlTop)
	for index, text := range v.IChild.GetHeadText() {
		rect := rects[index]
		rect.CenterY += 5
		newLabel(panel, text, rect)
		panel.SetWidth(panel.Width() + rect.Width)
	}
	return panel.Height()
}

func (v *boundView) drawEditRow() int32 {
	rects := v.IChild.GetColRects()
	length := len(rects)
	panel := newPanel(v.form, Rect{CenterX: 0, CenterY: 0, Width: v.form.Width(), Height: rects[0].Height}, types.AlTop)
	_ = newLabel(panel, "#", rects[0])
	v.IChild.NewEditControls(panel, rects)
	newButton(panel, AddButtonText, rects[length-4], v.onAddButtonClick)
	newButton(panel, DefaultButtonText, rects[length-3], v.onDefaultButtonClick)
	newButton(panel, ClearButtonText, rects[length-2], v.onClearButtonClick)
	return panel.Height()
}

func (v *boundView) drawDataBox() int32 {
	v.scrollBox = newScrollBox(v.form, v.form.Width(), 300, types.AlClient)
	return v.scrollBox.Height()
}

func (v *boundView) drawDataRows() {
	bounds := v.IChild.GetBounds()
	rects := v.IChild.GetColRects()
	length := len(rects)
	for rowIndex, cols := range bounds {
		v.drawDataCols(rowIndex, cols, rects, length)
	}
	v.scrollBox.Invalidate()
}

func (v *boundView) drawDataCols(rowIndex int, cols []string, rects []Rect, length int) {
	panel := newPanel(v.scrollBox, Rect{CenterX: 0, CenterY: int32(rowIndex) * rects[0].Height, Width: v.form.Width(), Height: rects[0].Height}, types.AlNone)
	panel.SetTag(rowIndex)
	_ = newLabel(panel, strconv.Itoa(rowIndex), rects[0])
	v.IChild.NewDataControls(panel, cols, rects)
	newButton(panel, CopyButtonText, rects[length-5], v.onCopyButtonClick)
	newButton(panel, UpdateButtonText, rects[length-4], v.onUpdateButtonClick)
	newButton(panel, DeleteButtonText, rects[length-3], v.onDeleteButtonClick)
	newButton(panel, PageUpButtonText, rects[length-2], v.onUpButtonClick)
	newButton(panel, PageDownButtonText, rects[length-1], v.onDownButtonClick)
	v.dataPanels = append(v.dataPanels, panel)
}

func (v *boundView) onAddButtonClick(sender vcl.IObject) {
	cols := v.AddBound()
	if cols == nil {
		return
	}
	rects := v.GetColRects()
	length := len(rects)
	v.drawDataCols(len(v.dataPanels), cols, rects, length)
	v.scrollBox.Invalidate()
}

func (v *boundView) onDefaultButtonClick(sender vcl.IObject) {
	v.IChild.SetDefaultText()
}

func (v *boundView) onClearButtonClick(sender vcl.IObject) {
	v.IChild.ClearEditText()
}

func (v *boundView) onCopyButtonClick(sender vcl.IObject) {
	parent := vcl.AsButton(sender).Parent()
	rowIndex := parent.Tag()
	panel := v.dataPanels[rowIndex]
	input := make([]string, 0)
	for i := int32(0); i < panel.ControlCount(); i++ {
		control := panel.Controls(i)
		if control.Is().Edit() {
			edit := vcl.AsEdit(control)
			input = append(input, edit.Text())
		}
		if control.Is().ComboBox() {
			comboBox := vcl.AsComboBox(control)
			input = append(input, comboBox.Text())
		}
	}
	v.IChild.CopyBound(input)
	bounds := v.IChild.GetBounds()
	rowIndex = len(bounds) - 1
	rects := v.IChild.GetColRects()
	length := len(rects)
	v.drawDataCols(rowIndex, input, rects, length)
	v.scrollBox.Invalidate()
}

func (v *boundView) onUpdateButtonClick(sender vcl.IObject) {
	parent := vcl.AsButton(sender).Parent()
	rowIndex := parent.Tag()
	panel := v.dataPanels[rowIndex]
	input := make([]string, 0)
	for i := int32(0); i < panel.ControlCount(); i++ {
		control := panel.Controls(i)
		if control.Is().Edit() {
			edit := vcl.AsEdit(control)
			input = append(input, edit.Text())
		}
		if control.Is().ComboBox() {
			comboBox := vcl.AsComboBox(control)
			input = append(input, comboBox.Text())
		}
	}
	v.IChild.UpdateBound(rowIndex, input)
	vcl.ShowMessage("Update successful")
}

func (v *boundView) onDeleteButtonClick(sender vcl.IObject) {
	deletePanel := vcl.AsButton(sender).Parent()
	rowIndex := deletePanel.Tag()
	v.IChild.DeleteBound(rowIndex)
	dataPanels := v.dataPanels[:rowIndex]
	for i := rowIndex; i < len(v.dataPanels)-1; i++ {
		nextPanel := v.dataPanels[i+1]
		nextPanel.SetTag(i)
		nextPanel.SetTop(int32(i) * nextPanel.Height())
		control := nextPanel.Controls(0)
		label := vcl.AsLabel(control)
		label.SetCaption(strconv.Itoa(i))
		dataPanels = append(dataPanels, nextPanel)
	}
	v.dataPanels = dataPanels
	deletePanel.Free()
	v.scrollBox.Invalidate()
}

func (v *boundView) onUpButtonClick(sender vcl.IObject) {
	parent := vcl.AsButton(sender).Parent()
	rowIndex := parent.Tag()
	if rowIndex-1 < 0 {
		vcl.ShowMessage("current row is on top")
		return
	}
	v.IChild.UpBound(rowIndex)
	srcPanel := vcl.AsPanel(parent)
	dstPanel := v.dataPanels[rowIndex-1]
	v.exchangePosition(dstPanel, srcPanel)
}

func (v *boundView) onDownButtonClick(sender vcl.IObject) {
	parent := vcl.AsButton(sender).Parent()
	rowIndex := parent.Tag()
	if rowIndex == len(v.dataPanels)-1 {
		vcl.ShowMessage("current row is at bottom")
		return
	}
	v.IChild.DownBound(rowIndex)
	srcPanel := vcl.AsPanel(parent)
	dstPanel := v.dataPanels[rowIndex+1]
	v.exchangePosition(dstPanel, srcPanel)
}

func (v *boundView) exchangePosition(dstPanel *vcl.TPanel, srcPanel *vcl.TPanel) {
	dstTag := dstPanel.Tag()
	dstPanel.SetTag(srcPanel.Tag())
	dstPanel.SetTop(int32(dstPanel.Tag()) * dstPanel.Height())
	vcl.AsLabel(dstPanel.Controls(0)).SetCaption(strconv.Itoa(dstPanel.Tag()))
	srcPanel.SetTag(dstTag)
	srcPanel.SetTop(int32(srcPanel.Tag()) * srcPanel.Height())
	vcl.AsLabel(srcPanel.Controls(0)).SetCaption(strconv.Itoa(srcPanel.Tag()))
	v.dataPanels[dstPanel.Tag()] = dstPanel
	v.dataPanels[srcPanel.Tag()] = srcPanel
	v.scrollBox.Invalidate()
}

func (v *boundView) freeDataPanels() {
	for _, panel := range v.dataPanels {
		panel.Free()
	}
	v.dataPanels = make([]*vcl.TPanel, 0)
}

func (v *boundView) onFormClose(sender vcl.IObject, action *types.TCloseAction) {
	ok := v.IChild.SaveConfig()
	if !ok {
		*action = types.CaNone
		vcl.ShowMessage("Save fails")
		return
	}
	v.freeDataPanels()
	v.sendRefreshTags(v.IChild.GetTags())
}

// SetRefreshTagFunc implements BoundView.
func (v *boundView) SetRefreshTagsFunc(onRefreshTags func(boundType string, tag []string)) {
	v.sendRefreshTags = onRefreshTags
}
