package controller

import (
	"sysproxy/config"
)

type MenuController interface {
	GetMenu() *config.Menu
	GetOutboundTags() []string
	GetInboundTags() []string
	SwitchOutbound(index int) error
	SwitchInbound(index int) error
	ToggleAutoStart() error
	OpenAboutURL()
	Exit() string
}

type LogController interface {
	GetLogInfo() string
}

type RuleController interface {
	GetUserRule() []string
	SaveUserRule(rules string) error
}

type BoundController interface {
	GetInProtoList() []string
	GetOutProtoList() []string
	GetBoundList() [][]string
	GetBoundTags() []string
	AddBound(bound []string)
	UpdateBound(index int, bound []string)
	DeleteBound(index int)
	PageUp(index int)
	PageDown(index int)
	SaveConfig() bool
}
