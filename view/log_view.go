package view

import (
	"strings"
	"sysproxy/controller"
	"time"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type logView struct {
	controller controller.LogController
	form       *vcl.TForm
	listBox    *vcl.TListBox
	isHide     bool
}

func NewLogView(controller controller.LogController) View {
	v := &logView{
		controller: controller,
		form:       newForm("Log", 500, 300),
	}
	v.form.SetOnClose(v.onFormClose)
	v.DrawControls()
	return v
}

// Show implements View.
func (v *logView) Show() {
	v.isHide = false
	go func() {
		for {
			if v.isHide {
				break
			}
			logInfo := v.controller.GetLogInfo()
			if logInfo == "" {
				time.Sleep(2 * time.Second)
				continue
			}
			logMessages := strings.Split(logInfo, "\n")
			if v.listBox != nil {
				for _, msg := range logMessages {
					if msg == "" {
						continue
					}
					v.listBox.Items().Add(msg)
				}
			}
			v.listBox.Invalidate()
		}
	}()
	v.form.Show()
}

// DrawControls implements View.
func (v *logView) DrawControls() {
	v.listBox = vcl.NewListBox(v.form)
	v.listBox.SetName("LogListBox")
	v.listBox.SetParent(v.form)
	v.listBox.SetAlign(types.AlClient)
}

func (v *logView) onFormClose(sender vcl.IObject, action *types.TCloseAction) {
	v.isHide = true
}
