package view

import (
	"sysproxy/controller"

	"github.com/ying32/govcl/vcl"
)

type outboundView struct {
	BoundView
	controller controller.BoundController
	tagEdit    *vcl.TEdit
	srcProto   *vcl.TComboBox
	srcIP      *vcl.TEdit
	srcPort    *vcl.TEdit
	dstProto   *vcl.TComboBox
	dstIP      *vcl.TEdit
	dstPort    *vcl.TEdit
}

func NewOutboundView(controller controller.BoundController) BoundView {
	v := new(outboundView)
	v.BoundView = NewBoundView(v)
	v.controller = controller
	v.BoundView.DrawControls()
	return v
}

// GetFormTitle implements IChild.
func (v *outboundView) GetFormTitle() string {
	return outboundText
}

// GetColRects implements IChild.
func (v *outboundView) GetColRects() []Rect {
	rects := []Rect{
		{CenterX: 0, CenterY: 0, Width: 50, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 150, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 100, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 100, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 50, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 100, Height: 25},
		{CenterX: 0, CenterY: 0, Width: 100, Height: 25},
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
func (v *outboundView) GetBounds() [][]string {
	return v.controller.GetBoundList()
}

// GetHeadText implements IChild.
func (v *outboundView) GetHeadText() []string {
	return []string{
		IDText,
		TagText,
		SrcProtoText,
		SrcIPText,
		SrcPortText,
		DstProtoText,
		DstIPText,
		DstPortText,
	}
}

// NewEditControls implements IChild.
func (v *outboundView) NewEditControls(panel *vcl.TPanel, rects []Rect) {
	v.tagEdit = newEdit(panel, "", rects[1])
	v.srcProto = newComboBox(panel, "", rects[2], v.getInProtoList())
	v.srcIP = newEdit(panel, "", rects[3])
	v.srcPort = newEdit(panel, "", rects[4])
	v.dstProto = newComboBox(panel, "", rects[5], v.getOutProtoList())
	v.dstIP = newEdit(panel, "", rects[6])
	v.dstPort = newEdit(panel, "", rects[7])
}

// NewDataControls implements IChild.
func (v *outboundView) NewDataControls(panel *vcl.TPanel, cols []string, rects []Rect) {
	_ = newEdit(panel, cols[0], rects[1])
	_ = newComboBox(panel, cols[1], rects[2], v.controller.GetInProtoList())
	_ = newEdit(panel, cols[2], rects[3])
	_ = newEdit(panel, cols[3], rects[4])
	_ = newComboBox(panel, cols[4], rects[5], v.controller.GetOutProtoList())
	_ = newEdit(panel, cols[5], rects[6])
	_ = newEdit(panel, cols[6], rects[7])
}

func (v *outboundView) getInProtoList() []string {
	return v.controller.GetInProtoList()
}

func (v *outboundView) getOutProtoList() []string {
	return v.controller.GetOutProtoList()
}

// SetDefaultText implements IChild.
func (v *outboundView) SetDefaultText() {
	v.tagEdit.SetText("pac-8848->http-8080")
	v.srcProto.SetText("pac")
	v.srcIP.SetText("127.0.0.1")
	v.srcPort.SetText("8848")
	v.dstProto.SetText("http")
	v.dstIP.SetText("127.0.0.1")
	v.dstPort.SetText("8080")
}

// ClearEditText implements IChild.
func (v *outboundView) ClearEditText() {
	v.tagEdit.SetText("")
	v.srcProto.SetText("")
	v.srcIP.SetText("")
	v.srcPort.SetText("")
	v.dstProto.SetText("")
	v.dstIP.SetText("")
	v.dstPort.SetText("")
}

// AddBound implements IChild.
func (v *outboundView) AddBound() []string {
	tag := v.tagEdit.Text()
	srcProto := v.srcProto.Text()
	srcIP := v.srcIP.Text()
	srcPort := v.srcPort.Text()
	dstProto := v.dstProto.Text()
	dstIP := v.dstIP.Text()
	dstPort := v.dstPort.Text()
	if tag == "" || srcProto == "" || srcIP == "" || srcPort == "" || dstProto == "" || dstIP == "" || dstPort == "" {
		vcl.ShowMessage("Empty value")
		return nil
	}
	bound := []string{tag, srcProto, srcIP, srcPort, dstProto, dstIP, dstPort}
	v.controller.AddBound(bound)
	return bound
}

// UpdateBound implements IChild.
func (v *outboundView) UpdateBound(rowIndex int, input []string) {
	v.controller.UpdateBound(rowIndex, input)
}

// DeleteBound implements IChild.
func (v *outboundView) DeleteBound(rowIndex int) {
	v.controller.DeleteBound(rowIndex)
}

// UpBound implements IChild.
func (v *outboundView) UpBound(rowIndex int) {
	v.controller.PageUp(rowIndex)
}

// DownBound implements IChild.
func (v *outboundView) DownBound(rowIndex int) {
	v.controller.PageDown(rowIndex)
}

// GetTags implements IChild.
func (v *outboundView) GetTags() (string, []string) {
	return outboundName, v.controller.GetBoundTags()
}

// SaveConfig implements IChild.
func (v *outboundView) SaveConfig() bool {
	return v.controller.SaveConfig()
}
