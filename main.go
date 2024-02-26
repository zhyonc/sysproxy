package main

import (
	"sysproxy/config"
	"sysproxy/controller"
	"sysproxy/util"
	"sysproxy/view"
)

func main() {
	if !CheckSingleton() {
		return
	}
	conf := config.NewConfig()
	menuController := controller.NewMenuController(conf)
	logController := controller.NewLogController()
	logView := view.NewLogView(logController)
	ruleController := controller.NewRuleController(menuController)
	ruleView := view.NewRuleView(ruleController)
	inboundController := controller.NewInboundController(conf)
	inboundView := view.NewInboundView(inboundController)
	outboundController := controller.NewOutboundController(conf)
	outboundView := view.NewOutboundView(outboundController)
	menuView := view.NewMenuView(menuController, logView, ruleView, inboundView, outboundView)
	menuView.Run()
}

func CheckSingleton() bool {
	_, err := util.CheckSingleton()
	if err != nil {
		util.ShowMessage("Another sysproxy is already running")
		return false
	}
	return true
}
