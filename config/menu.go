package config

type Menu struct {
	AppName              string
	Version              string
	OutboundCheckedIndex int
	InboundCheckedIndex  int
	AutoStart            bool
	EnableGFWList        bool
}
