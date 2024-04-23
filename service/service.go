package service

import (
	"sysproxy/config"
)

const (
	PAC    string = "pac"
	HTTP   string = "http"
	SOCKS5 string = "socks5"
)

var InProtoList []string = []string{PAC, HTTP, SOCKS5}

var OutProtoList []string = []string{HTTP, SOCKS5}

type BaseService interface {
	StartService()
	StopService()
}

type PacService interface {
	CreatePACTempFile(outbound config.Outbound) error
	GetUserRule() []string
	SaveUserRule(rule string) error
	SetEnableGFWList(enableGFWList bool)
}
