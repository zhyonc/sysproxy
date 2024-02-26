package view

import (
	"sysproxy/controller"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type RuleView struct {
	controller controller.RuleController
	form       *vcl.TForm
	memo       *vcl.TMemo
}

func NewRuleView(controller controller.RuleController) View {
	v := &RuleView{
		controller: controller,
		form:       newForm("Rule", 500, 300),
	}
	v.form.SetOnClose(v.onFormClose)
	v.DrawControls()
	return v
}

// Show implements View.
func (v *RuleView) Show() {
	if v.memo != nil {
		v.memo.Clear()
		lines := v.controller.GetUserRule()
		for _, line := range lines {
			v.memo.Lines().Add(line)
		}
	}
	v.form.Show()
}

// DrawControls implements View.
func (v *RuleView) DrawControls() {
	v.memo = vcl.NewMemo(v.form)
	v.memo.SetParent(v.form)
	v.memo.SetName("RuleMemo")
	v.memo.SetAlign(types.AlClient)
	v.memo.Clear()
}

func (v *RuleView) onFormClose(sender vcl.IObject, action *types.TCloseAction) {
	if v.memo == nil {
		vcl.ShowMessage("Can't get rule edit result")
		return
	}
	err := v.controller.SaveUserRule(v.memo.Text())
	if err != nil {
		vcl.ShowMessage("Save user rule failed: " + err.Error())
		return
	}
	vcl.ShowMessage("Save user rule successful")
}
