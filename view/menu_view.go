package view

import (
	"fmt"
	"sysproxy/controller"
	"sysproxy/resources"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type menuView struct {
	controller   controller.MenuController
	logView      View
	ruleView     View
	inboundView  BoundView
	outboundView BoundView
	popupMenu    *vcl.TPopupMenu
}

func NewMenuView(controller controller.MenuController, logView, ruleView View, inboundView, outboundView BoundView) *menuView {
	v := &menuView{
		controller:   controller,
		logView:      logView,
		ruleView:     ruleView,
		inboundView:  inboundView,
		outboundView: outboundView,
	}
	appName := v.controller.GetMenu().AppName
	v.inboundView.SetRefreshTagsFunc(v.onRefreshTags)
	v.outboundView.SetRefreshTagsFunc(v.onRefreshTags)
	v.drawPopupMenu(newTrayIcon(newMainForm(), appName))
	return v
}

func newMainForm() *vcl.TForm {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(false)
	form := vcl.Application.CreateForm()
	form.SetPosition(types.PoScreenCenter)
	form.SetOnShow(func(sender vcl.IObject) {
		form.Hide()
	})
	return form
}

func newTrayIcon(form *vcl.TForm, appName string) *vcl.TTrayIcon {
	trayIcon := vcl.NewTrayIcon(form)
	trayIcon.SetHint(appName)
	icon := vcl.NewIcon()
	icon.LoadFromBytes(resources.IconData)
	trayIcon.SetIcon(icon)
	return trayIcon
}

func (v *menuView) drawPopupMenu(trayIcon *vcl.TTrayIcon) {
	v.popupMenu = vcl.NewPopupMenu(trayIcon)
	var index int = 0
	v.popupMenu.Items().Add(newMenuItem(v.popupMenu, &index, outboundName, outboundText, nil))
	v.popupMenu.Items().Add(newMenuItem(v.popupMenu, &index, inboundName, inboundText, nil))
	v.popupMenu.Items().Add(newMenuItem(v.popupMenu, &index, logName, logText, v.showForm))
	v.popupMenu.Items().Add(newMenuItem(v.popupMenu, &index, ruleName, ruleText, v.showForm))
	v.popupMenu.Items().Add(v.newSettingMenuItem(v.popupMenu, &index, settingName, settingText, nil))
	v.popupMenu.Items().Add(newMenuItem(v.popupMenu, &index, aboutName, aboutText, v.onAboutMenuItemClick))
	v.popupMenu.Items().Add(newMenuItem(v.popupMenu, &index, exitName, exitText, v.onExitMenuItemClick))
	trayIcon.SetPopupMenu(v.popupMenu)
	trayIcon.SetVisible(true)
}

func (v *menuView) newSettingMenuItem(parent vcl.IComponent, index *int, name, title string, clickEvent vcl.TNotifyEvent) *vcl.TMenuItem {
	settingMenuItem := newMenuItem(parent, index, name, title, clickEvent)
	var subIndex int = 0
	settingMenuItem.Add(newMenuItem(settingMenuItem, &subIndex, settingOutboundName, settingOutboundText, v.showForm))
	settingMenuItem.Add(newMenuItem(settingMenuItem, &subIndex, settingInboundName, settingInboundText, v.showForm))
	menu := v.controller.GetMenu()
	autoProxyMenuItem := newMenuItem(settingMenuItem, &subIndex, settingAutoProxyName, settingAutoProxyText, v.onSettingAutoProxyClick)
	autoProxyMenuItem.SetChecked(menu.AutoProxy)
	settingMenuItem.Add(autoProxyMenuItem)
	autoStartMenuItem := newMenuItem(settingMenuItem, &subIndex, settingAutoStartName, settingAutoStartText, v.onSettingAutoStartClick)
	autoStartMenuItem.SetChecked(menu.AutoStart)
	settingMenuItem.Add(autoStartMenuItem)
	return settingMenuItem
}

func (v *menuView) Run() {
	v.onRefreshTags(outboundName, v.controller.GetOutboundTags())
	v.onRefreshTags(inboundName, v.controller.GetInboundTags())
	vcl.Application.Run()
}

func (v *menuView) showForm(sender vcl.IObject) {
	menuItem := vcl.AsMenuItem(sender)
	switch menuItem.Caption() {
	case logText:
		v.logView.Show()
	case ruleText:
		v.ruleView.Show()
	case settingInboundText:
		v.inboundView.Show()
	case settingOutboundText:
		v.outboundView.Show()
	}
}

func (v *menuView) onSettingAutoProxyClick(sender vcl.IObject) {
	menuItem := vcl.AsMenuItem(sender)
	menu := v.controller.GetMenu()
	checked := !menu.AutoProxy
	menuItem.SetChecked(checked)
	menu.AutoProxy = checked
}

func (v *menuView) onSettingAutoStartClick(sender vcl.IObject) {
	menuItem := vcl.AsMenuItem(sender)
	menu := v.controller.GetMenu()
	checked := !menu.AutoStart
	menuItem.SetChecked(checked)
	menu.AutoStart = checked
	err := v.controller.ToggleAutoStart()
	if err != nil {
		vcl.ShowMessage("Error toggle auto start: " + err.Error())
	}
}

func (v *menuView) onAboutMenuItemClick(sender vcl.IObject) {
	menu := v.controller.GetMenu()
	var msg string = "Version: " + menu.Version + "\n" + "Windows system proxy forward tool"
	res := vcl.MessageDlg(msg, types.MtConfirmation, 0)
	if res == types.MrYes {
		v.controller.OpenAboutURL()
	}
}

func (v *menuView) onExitMenuItemClick(sender vcl.IObject) {
	msg := v.controller.Exit()
	if msg != "" {
		vcl.ShowMessage(msg)
	}
	vcl.Application.Terminate()
}

func (v *menuView) onRefreshTags(name string, tags []string) {
	var boundMenuItem *vcl.TMenuItem
	for i := int32(0); i < v.popupMenu.Items().Count(); i++ {
		item := v.popupMenu.Items().Items(i)
		if item.Name() == name {
			boundMenuItem = item
			break
		}
	}
	if boundMenuItem == nil {
		vcl.ShowMessageFmt("Refresh %s tags failed", name)
		return
	}
	boundMenuItem.Clear()
	var index int = 0
	boundMenuItem.Add(newMenuItem(boundMenuItem, &index, "", "disable", v.onBoundSubMenuItemClick))
	for _, tag := range tags {
		boundMenuItem.Add(newMenuItem(boundMenuItem, &index, "", tag, v.onBoundSubMenuItemClick))
	}
	menu := v.controller.GetMenu()
	if menu.AutoProxy {
		var subMenuItem *vcl.TMenuItem
		if boundMenuItem.Name() == outboundName {
			subMenuItem = boundMenuItem.Items(int32(menu.OutboundCheckedIndex))
		} else if boundMenuItem.Name() == inboundName {
			subMenuItem = boundMenuItem.Items(int32(menu.InboundCheckedIndex))
		} else {
			subMenuItem = nil
		}
		if subMenuItem != nil {
			subMenuItem.Click()
			return
		}
	}
	disableMenuItem := boundMenuItem.Items(0)
	disableMenuItem.Click()
}

func (v *menuView) onBoundSubMenuItemClick(sender vcl.IObject) {
	menuItem := vcl.AsMenuItem(sender)
	if menuItem == nil {
		return
	}
	parentMenuItem := menuItem.Parent()
	var err error
	if parentMenuItem.Name() == outboundName {
		err = v.controller.SwitchOutbound(menuItem.Tag())
	} else if parentMenuItem.Name() == inboundName {
		err = v.controller.SwitchInbound(menuItem.Tag())
	} else {
		err = fmt.Errorf("can't find bound type")
	}
	if err != nil {
		vcl.ShowMessage(menuItem.Caption() + ": " + err.Error())
		return
	}
	for i := 0; i < int(parentMenuItem.Count()); i++ {
		item := parentMenuItem.Items(int32(i))
		item.SetChecked(false)
	}
	menuItem.SetChecked(true)
}
