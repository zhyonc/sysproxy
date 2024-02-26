package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	appName string = "sysproxy"
	version string = "0.1.1"
)

var path string = "config.toml"

type Config struct {
	OutboundList []Outbound
	InboundList  []Inbound
	Menu         *Menu
}

func NewConfig() *Config {
	conf := &Config{}
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		log.Printf("Error load config from %s\n", path)
		conf = defualtConfig()
		if !conf.Save() {
			return nil
		}
	}
	if conf.Menu.Version < version {
		upgradeConfig(conf)
	}
	return conf
}

func defualtConfig() *Config {
	conf := &Config{
		OutboundList: make([]Outbound, 0),
		InboundList:  make([]Inbound, 0),
		Menu: &Menu{
			AppName:              appName,
			Version:              version,
			InboundCheckedIndex:  0,
			OutboundCheckedIndex: 0,
			AutoStart:            false,
			AutoProxy:            false,
		},
	}
	return conf
}

func upgradeConfig(oldConf *Config) *Config {
	newConf := defualtConfig()
	newConf.OutboundList = oldConf.OutboundList
	newConf.InboundList = oldConf.InboundList
	newConf.Menu.AutoStart = oldConf.Menu.AutoStart
	newConf.Menu.AutoProxy = oldConf.Menu.AutoProxy
	ok := newConf.Save()
	if !ok {
		log.Println("update config failed")
		return oldConf
	}
	log.Println("upgrade config successful")
	return newConf
}

func (c *Config) Save() bool {
	file, err := os.Create(path)
	if err != nil {
		log.Printf("Error create config to %s\n", path)
		return false
	}
	defer file.Close()
	encoder := toml.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		log.Printf("Error encode config: %v\n", err)
		return false
	}
	log.Println("Save config successful")
	return true
}
