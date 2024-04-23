package view

import (
	"sysproxy/controller"

	"github.com/ying32/govcl/vcl"
)

type inboundView struct {
	BoundView
	controller controller.BoundController
	tagEdit    *vcl.TEdit
	dstProto   *vcl.TComboBox
	dstIP      *vcl.TEdit
	dstPort    *vcl.TEdit
}

func NewInboundView(controller controller.BoundController) BoundView {
	v := new(inboundView)
	v.BoundView = NewBoundView(v)
	v.controller = controller
	v.BoundView.DrawControls()
	return v
}

// GetFormTitle implements IChild.
func (v *inboundView) GetFormTitle() string {
	return inboundText
}

// GetColRects implements IChild.
func (v *inboundView) GetColRects() []Rect {
	rects := []Rect{
		{CenterX: 0, CenterY: 0, Width: 50, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 100, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 100, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 100, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 50, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 50, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 50, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 50, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 50, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 50, Height: 25},
	}
	for i := 0; i < len(rects)-1; i++ {
		rects[i+1].CenterX = rects[i].CenterX + rects[i].Width
	}
	return rects
}

// GetBounds implements IChild.
func (v *inboundView) GetBounds() [][]string {
	return v.controller.GetBoundList()
}

// GetHeadText implements IChild.
func (v *inboundView) GetHeadText() []string {
	return []string{
		IDText,
		TagText,
		DstProtoText,
		DstIPText,
		DstPortText,
	}
}

// NewEditControls implements IChild.
func (v *inboundView) NewEditControls(panel *vcl.TPanel, rects []Rect) {
	v.tagEdit = newEdit(panel, "", rects[1])
	v.dstProto = newComboBox(panel, "", rects[2], v.getInProtoList())
	v.dstIP = newEdit(panel, "", rects[3])
	v.dstPort = newEdit(panel, "", rects[4])
}

// NewDataControls implements IChild.
func (v *inboundView) NewDataControls(panel *vcl.TPanel, cols []string, rects []Rect) {
	_ = newEdit(panel, cols[0], rects[1])
	_ = newComboBox(panel, cols[1], rects[2], v.controller.GetInProtoList())
	_ = newEdit(panel, cols[2], rects[3])
	_ = newEdit(panel, cols[3], rects[4])
}

func (v *inboundView) getInProtoList() []string {
	return v.controller.GetInProtoList()
}

// SetDefaultText implements IChild.
func (v *inboundView) SetDefaultText() {
	v.tagEdit.SetText("pac-8848")
	v.dstProto.SetText("pac")
	v.dstIP.SetText("127.0.0.1")
	v.dstPort.SetText("8848")
}

// ClearEditText implements IChild.
func (v *inboundView) ClearEditText() {
	v.tagEdit.SetText("")
	v.dstProto.SetText("")
	v.dstIP.SetText("")
	v.dstPort.SetText("")
}

// AddBound implements IChild.
func (v *inboundView) AddBound() []string {
	tag := v.tagEdit.Text()
	dstProto := v.dstProto.Text()
	dstIP := v.dstIP.Text()
	dstPort := v.dstPort.Text()
	if tag == "" || dstProto == "" || dstIP == "" || dstPort == "" {
		vcl.ShowMessage("Empty value")
		return nil
	}
	bound := []string{tag, dstProto, dstIP, dstPort}
	v.controller.AddBound(bound)
	return bound
}

// CopyBound implements IChild.
func (v *inboundView) CopyBound(input []string) {
	v.controller.AddBound(input)
}

// UpdateBound implements IChild.
func (v *inboundView) UpdateBound(rowIndex int, input []string) {
	v.controller.UpdateBound(rowIndex, input)
}

// DeleteBound implements IChild.
func (v *inboundView) DeleteBound(rowIndex int) {
	v.controller.DeleteBound(rowIndex)
}

// UpBound implements IChild.
func (v *inboundView) UpBound(rowIndex int) {
	v.controller.PageUp(rowIndex)
}

// DownBound implements IChild.
func (v *inboundView) DownBound(rowIndex int) {
	v.controller.PageDown(rowIndex)
}

// GetTags implements IChild.
func (v *inboundView) GetTags() (string, []string) {
	return inboundName, v.controller.GetBoundTags()
}

// SaveConfig implements IChild.
func (v *inboundView) SaveConfig() bool {
	return v.controller.SaveConfig()
}
