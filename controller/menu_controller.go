package controller

import (
	"fmt"
	"sysproxy/config"
	"sysproxy/service"
	"sysproxy/util"
)

const (
	aboutURL string = "https://github.com/zhyonc/sysproxy"
)

type menuController struct {
	conf    *config.Config
	service service.BaseService
}

func NewMenuController(conf *config.Config) MenuController {
	c := &menuController{
		conf: conf,
	}
	return c
}

// GetMenu implements MenuController.
func (c *menuController) GetMenu() *config.Menu {
	return c.conf.Menu
}

// GetOutboundTags implements MenuController.
func (c *menuController) GetOutboundTags() []string {
	tags := make([]string, 0)
	for _, v := range c.conf.OutboundList {
		tags = append(tags, v.Tag)
	}
	return tags
}

// GetInboundTags implements MenuController.
func (c *menuController) GetInboundTags() []string {
	tags := make([]string, 0)
	for _, v := range c.conf.InboundList {
		tags = append(tags, v.Tag)
	}
	return tags
}

// SwitchOutbound implements MenuController.
func (c *menuController) SwitchOutbound(index int) error {
	if index < 0 || index > len(c.conf.OutboundList) {
		return fmt.Errorf("can't find outbound at %d", index)
	}
	if c.service != nil {
		c.service.StopService()
		c.service = nil
	}
	if index > 0 {
		outbound := c.conf.OutboundList[index-1] // OutboundList don't include disable tag
		if outbound.SrcProto == service.PAC || outbound.SrcProto == service.HTTP {
			err := service.PACService.CreatePACTempFile(outbound)
			if err != nil {
				return err
			}
			c.service = service.NewHttpService(outbound)
		} else if outbound.SrcProto == service.SOCKS5 {
			c.service = service.NewSocks5Service(outbound)
		} else {
			return fmt.Errorf("unknown src protocol at %d", index)
		}
		go c.service.StartService()
	}
	return nil
}

// SwitchInbound implements MenuController.
func (c *menuController) SwitchInbound(index int) error {
	if index < 0 || index > len(c.conf.InboundList) {
		return fmt.Errorf("can't find inbound at %d", index)
	}
	if index == 0 {
		util.ClearProxy()
		return nil
	}
	inbound := c.conf.InboundList[index-1] // InboundList don't include disable tag
	switch inbound.DstProto {
	case service.PAC:
		util.EnablePAC(inbound.DstIP, inbound.DstPort)
	case service.HTTP:
		util.EnableHTTP(inbound.DstIP, inbound.DstPort)
	case service.SOCKS5:
		util.EnableSOCKS5(inbound.DstIP, inbound.DstPort)
	default:
		return fmt.Errorf("unknown dst protocol is %s", inbound.DstProto)
	}
	return nil
}

// ToggleAutoStart implements MenuController.
func (c *menuController) ToggleAutoStart() error {
	return util.HookAutoStart(c.conf.Menu.AutoStart, c.conf.Menu.AppName)
}

// ToggleEnableGFWList implements MenuController.
func (c *menuController) ToggleEnableGFWList() {
	service.PACService.SetEnableGFWList(c.conf.Menu.EnableGFWList)
}

// OpenAboutURL implements MenuController.
func (c *menuController) OpenAboutURL() {
	util.OpenURL(aboutURL)
}

// Exit implements MenuController.
func (c *menuController) Exit() string {
	var msg string = ""
	err := util.ClearProxy()
	if err != nil {
		msg += err.Error() + "\n"
	}
	if c.service != nil {
		c.service.StopService()
	}
	if !c.SaveConfig() {
		msg += "Save config failed\n"
	}
	return msg
}

// SaveConfig implements MenuController.
func (c *menuController) SaveConfig() bool {
	return c.conf.Save()
}
