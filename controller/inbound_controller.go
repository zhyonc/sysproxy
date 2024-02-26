package controller

import (
	"sysproxy/config"
	"sysproxy/service"
)

type inboundController struct {
	conf *config.Config
}

func NewInboundController(conf *config.Config) BoundController {
	c := &inboundController{
		conf: conf,
	}
	return c
}

// GetInProtoList implements BoundController.
func (c *inboundController) GetInProtoList() []string {
	return service.InProtoList
}

// GetOutProtoList implements BoundController.
func (c *inboundController) GetOutProtoList() []string {
	return nil
}

// GetBoundList implements BoundController.
func (c *inboundController) GetBoundList() [][]string {
	bounds := make([][]string, 0)
	for _, inbound := range c.conf.InboundList {
		cols := []string{
			inbound.Tag,
			inbound.DstProto,
			inbound.DstIP,
			inbound.DstPort,
		}
		bounds = append(bounds, cols)
	}
	return bounds
}

// GetBoundTags implements BoundController.
func (c *inboundController) GetBoundTags() []string {
	tags := make([]string, 0)
	for _, inbound := range c.conf.InboundList {
		tags = append(tags, inbound.Tag)
	}
	return tags
}

// AddBound implements BoundController.
func (c *inboundController) AddBound(bound []string) {
	inbound := config.Inbound{
		Tag:      bound[0],
		DstProto: bound[1],
		DstIP:    bound[2],
		DstPort:  bound[3],
	}
	c.conf.InboundList = append(c.conf.InboundList, inbound)
}

// UpdateBound implements BoundController.
func (c *inboundController) UpdateBound(index int, bound []string) {
	c.conf.InboundList[index].Tag = bound[0]
	c.conf.InboundList[index].DstProto = bound[1]
	c.conf.InboundList[index].DstIP = bound[2]
	c.conf.InboundList[index].DstPort = bound[3]
}

// DeleteBound implements BoundController.
func (c *inboundController) DeleteBound(index int) {
	if index+1 == len(c.conf.InboundList) {
		c.conf.InboundList = c.conf.InboundList[:index]
	} else {
		c.conf.InboundList = append(c.conf.InboundList[:index], c.conf.InboundList[index+1:]...)
	}
}

// PageUp implements BoundController.
func (c *inboundController) PageUp(index int) {
	if index-1 < 0 {
		return
	}
	upper := c.conf.InboundList[index-1]
	current := c.conf.InboundList[index]
	c.conf.InboundList[index-1] = current
	c.conf.InboundList[index] = upper
}

// PageDown implements BoundController.
func (c *inboundController) PageDown(index int) {
	if index >= len(c.conf.InboundList)-1 {
		return
	}
	lower := c.conf.InboundList[index+1]
	current := c.conf.InboundList[index]
	c.conf.InboundList[index+1] = current
	c.conf.InboundList[index] = lower
}

// SaveConfig implements BoundController.
func (c *inboundController) SaveConfig() bool {
	return c.conf.Save()
}
