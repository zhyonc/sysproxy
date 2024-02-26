package service_test

import (
	"log"
	"sysproxy/config"
	"sysproxy/service"
	"sysproxy/util"
	"testing"
)

func TestSocks5ToHttp(t *testing.T) {
	inbound := config.Inbound{
		DstProto: service.SOCKS5,
		DstIP:    "127.0.0.1",
		DstPort:  "8848",
	}
	outbound := config.Outbound{
		SrcProto: inbound.DstProto,
		SrcIP:    inbound.DstIP,
		SrcPort:  inbound.DstPort,
		DstProto: service.HTTP,
		DstIP:    "127.0.0.1",
		DstPort:  "8080",
	}
	socks5Service := service.NewSocks5Service(outbound)
	if socks5Service == nil {
		log.Println("fail to create socks5Service")
	} else {
		util.EnableSOCKS5(inbound.DstIP, inbound.DstPort)
		socks5Service.StartService()
	}
}

func TestSocks5ToSocks5(t *testing.T) {
	inbound := config.Inbound{
		DstProto: service.SOCKS5,
		DstIP:    "127.0.0.1",
		DstPort:  "8848",
	}
	outbound := config.Outbound{
		SrcProto: inbound.DstProto,
		SrcIP:    inbound.DstIP,
		SrcPort:  inbound.DstPort,
		DstProto: service.SOCKS5,
		DstIP:    "127.0.0.1",
		DstPort:  "8080",
	}
	socks5Service := service.NewSocks5Service(outbound)
	if socks5Service == nil {
		log.Println("fail to create socks5Service")
		return
	} else {
		util.EnableSOCKS5(inbound.DstIP, inbound.DstPort)
		socks5Service.StartService()
	}
}
