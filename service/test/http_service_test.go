package service_test

import (
	"log"
	"sysproxy/config"
	"sysproxy/service"
	"sysproxy/util"
	"testing"
)

// Sometimes browser cached pac.js and forward fails when keep the same port
func TestPAC2Http(t *testing.T) {
	inbound := config.Inbound{
		DstProto: service.PAC,
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
	err := service.PACService.CreatePACTempFile(outbound)
	if err != nil {
		log.Printf("fail to create pac temp file")
		return
	}
	httpService := service.NewHttpService(outbound)
	if httpService == nil {
		log.Println("fail to create httpService")
	} else {
		util.EnablePAC(inbound.DstIP, inbound.DstPort)
		httpService.StartService()
	}
}

func TestHttp2Http(t *testing.T) {
	inbound := config.Inbound{
		DstProto: service.HTTP,
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
	httpService := service.NewHttpService(outbound)
	if httpService == nil {
		log.Println("fail to create httpService")
	} else {
		util.EnableHTTP(inbound.DstIP, inbound.DstPort)
		httpService.StartService()
	}
}

func TestHttp2Socks5(t *testing.T) {
	inbound := config.Inbound{
		DstProto: service.HTTP,
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
	httpService := service.NewHttpService(outbound)
	if httpService == nil {
		log.Println("fail to create httpService")
	} else {
		util.EnableHTTP(inbound.DstIP, inbound.DstPort)
		httpService.StartService()
	}
}
