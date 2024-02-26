package controller

import (
	"sysproxy/config"
	"sysproxy/service"
)

type outboundController struct {
	conf *config.Config
}

func NewOutboundController(conf *config.Config) BoundController {
	c := &outboundController{
		conf: conf,
	}
	return c
}

// GetInProtoList implements BoundController.
func (c *outboundController) GetInProtoList() []string {
	return service.InProtoList
}

// GetOutProtoList implements BoundController.
func (*outboundController) GetOutProtoList() []string {
	return service.OutProtoList
}

// GetBoundList implements BoundController.
func (c *outboundController) GetBoundList() [][]string {
	bounds := make([][]string, 0)
	for _, outbound := range c.conf.OutboundList {
		cols := []string{
			outbound.Tag,
			outbound.SrcProto,
			outbound.SrcIP,
			outbound.SrcPort,
			outbound.DstProto,
			outbound.DstIP,
			outbound.DstPort,
		}
		bounds = append(bounds, cols)
	}
	return bounds
}

// GetBoundTags implements BoundController.
func (c *outboundController) GetBoundTags() []string {
	tags := make([]string, 0)
	for _, outbound := range c.conf.OutboundList {
		tags = append(tags, outbound.Tag)
	}
	return tags
}

// AddBound implements BoundController.
func (c *outboundController) AddBound(bound []string) {
	outbound := config.Outbound{
		Tag:      bound[0],
		SrcProto: bound[1],
		SrcIP:    bound[2],
		SrcPort:  bound[3],
		DstProto: bound[4],
		DstIP:    bound[5],
		DstPort:  bound[6],
	}
	c.conf.OutboundList = append(c.conf.OutboundList, outbound)
}

// UpdateBound implements BoundController.
func (c *outboundController) UpdateBound(index int, bound []string) {
	c.conf.OutboundList[index].Tag = bound[0]
	c.conf.OutboundList[index].SrcProto = bound[1]
	c.conf.OutboundList[index].SrcIP = bound[2]
	c.conf.OutboundList[index].SrcPort = bound[3]
	c.conf.OutboundList[index].DstProto = bound[4]
	c.conf.OutboundList[index].DstIP = bound[5]
	c.conf.OutboundList[index].DstPort = bound[6]
}

// DeleteBound implements BoundController.
func (c *outboundController) DeleteBound(index int) {
	if index+1 == len(c.conf.OutboundList) {
		c.conf.OutboundList = c.conf.OutboundList[:index]
	} else {
		c.conf.OutboundList = append(c.conf.OutboundList[:index], c.conf.OutboundList[index+1:]...)
	}
}

// PageUp implements BoundController.
func (c *outboundController) PageUp(index int) {
	if index-1 < 0 {
		return
	}
	upper := c.conf.OutboundList[index-1]
	current := c.conf.OutboundList[index]
	c.conf.OutboundList[index-1] = current
	c.conf.OutboundList[index] = upper
}

// PageDown implements BoundController.
func (c *outboundController) PageDown(index int) {
	if index >= len(c.conf.OutboundList)-1 {
		return
	}
	lower := c.conf.OutboundList[index+1]
	current := c.conf.OutboundList[index]
	c.conf.OutboundList[index+1] = current
	c.conf.OutboundList[index] = lower
}

// SaveConfig implements BoundController.
func (c *outboundController) SaveConfig() bool {
	return c.conf.Save()
}
