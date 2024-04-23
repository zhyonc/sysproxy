package service

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sysproxy/config"
	"sysproxy/resources"
	"sysproxy/util"
)

const (
	gfwlistURL       string = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"
	abpFilename      string = "abp.js"
	gfwlistFilename  string = "gfwlist.txt"
	userRuleFilename string = "user_rule.txt"
	PACFilename      string = "pac.js"
	PACTempFilename  string = "pac_temp.js"
	proxyWord        string = "__PROXY__"
	ruleWord         string = "__RULES__"
	tipText          string = "! Put user rules line by line in this file.\n! See https://adblockplus.org/en/filter-cheatsheet"
)

var IgnoredLineBegins []string = []string{"!", "["}

var PACService PacService

type pacService struct {
	name          string
	pacText       string
	enableGFWList bool
}

func NewPACService(conf *config.Config) {
	s := &pacService{
		name:          "pacService",
		enableGFWList: conf.Menu.EnableGFWList,
	}
	s.pacText = readPACFile()
	if s.pacText == "" {
		s.pacText = CreatePACFile(s.enableGFWList)
	}
	PACService = s
}

// CreatePACTempFile implements PacService.
func (s *pacService) CreatePACTempFile(outbound config.Outbound) error {
	if s.pacText == "" {
		return fmt.Errorf("can't find %s file", PACFilename)
	}
	temp := replacePACText(s.pacText, outbound.DstProto, net.JoinHostPort(outbound.DstIP, outbound.DstPort))
	return util.CreateFile(PACTempFilename, []byte(temp))
}

// GetUserRule implements PacService.
func (*pacService) GetUserRule() []string {
	return getContent(userRuleFilename)
}

// SaveUserRule implements PacService.
func (s *pacService) SaveUserRule(rule string) error {
	err := util.CreateFile(userRuleFilename, []byte(rule))
	if err != nil {
		return err
	}
	s.pacText = CreatePACFile(s.enableGFWList)
	return nil
}

// SetEnableGFWList implements PacService.
func (s *pacService) SetEnableGFWList(enableGFWList bool) {
	s.enableGFWList = enableGFWList
	s.pacText = CreatePACFile(s.enableGFWList)
}

func readPACFile() string {
	bytes, err := util.ReadFileAll(PACFilename)
	if err != nil {
		log.Printf("Error read file: %v\n", err)
		return ""
	}
	return string(bytes)
}

func replacePACText(pacText string, dstProto string, raddr string) string {
	if pacText == "" {
		return ""
	}
	var replaceWord string = ""
	if dstProto == SOCKS5 {
		replaceWord = "SOCKS5" + " " + raddr
	} else {
		replaceWord = "PROXY" + " " + raddr
	}
	if strings.Contains(pacText, proxyWord) {
		pacText = strings.Replace(pacText, proxyWord, replaceWord, -1)
	} else {
		log.Println("can't replace proxy word " + replaceWord)
	}
	return pacText
}

func CreatePACFile(enableGFWList bool) string {
	// Read file
	abpText := string(resources.AbpData)
	abpContent := strings.Split(abpText, "\n")
	if abpContent == nil {
		return ""
	}
	rules := make([]string, 0)
	if enableGFWList {
		gfwlistBytes, _ := util.ReadFileAll(gfwlistFilename)
		if gfwlistBytes == nil {
			gfwlistBytes = getFileFromURL(gfwlistFilename, gfwlistURL)
		}
		gfwlistContent, err := util.DecodeBase64(gfwlistBytes, "\n")
		if err != nil {
			log.Println("Convert gfwlist to content failed")
		} else {
			log.Println("Convert gfwlist to content successful")
		}
		gfwlistRule := formatRule(gfwlistContent)
		rules = append(rules, gfwlistRule...)
	}
	// Splice gfwlist and user rule
	userRuleContent := getContent(userRuleFilename)
	if userRuleContent == nil {
		err := util.CreateFile(userRuleFilename, []byte(tipText))
		if err != nil {
			log.Printf("Create %s failed:%v\n", userRuleFilename, err)
		}
		log.Printf("Create %s successful\n", userRuleFilename)
	} else {
		userRule := formatRule(userRuleContent)
		rules = append(rules, userRule...)
		log.Printf("Append user rule %d", len(userRule))
	}
	// Relace rule in abp.js
	allRule := fmt.Sprintf(`["%s"]`, strings.Join(rules, `","`))
	for i, s := range abpContent {
		if strings.Contains(s, ruleWord) {
			abpContent[i] = strings.Replace(s, ruleWord, allRule, -1)
			break
		}
	}
	// Crate new pac.js
	pacText := strings.Join(abpContent, "\n")
	pacBytes := []byte(pacText)
	err := util.CreateFile(PACFilename, pacBytes)
	if err != nil {
		log.Printf("Create %s failed:%v\n", PACFilename, err)
		return ""
	}
	log.Printf("Create %s successful\n", PACFilename)
	return pacText
}

func getContent(filename string) []string {
	content, err := util.ReadFileLine(filename)
	if err != nil {
		log.Printf("Read %s failed:%v\n", filename, err)
		return nil
	}
	return content
}

func formatRule(content []string) []string {
	var rules []string
	for i := 0; i < len(content); i++ {
		line := content[i]
		if line == "" {
			continue
		}
		firstChar := string([]rune(line)[0])
		isIgnoredLine := false
		for j := 0; j < len(IgnoredLineBegins); j++ {
			char := IgnoredLineBegins[j]
			if firstChar == char {
				isIgnoredLine = true
				break
			}
		}
		if isIgnoredLine {
			continue
		}
		rules = append(rules, content[i])
	}
	return rules
}

func getFileFromURL(filename string, url string) []byte {
	done := make(chan bool)
	var bytes []byte
	var err error
	go func() {
		log.Printf("Download %s from %s\n", filename, url)
		bytes, err = util.DownloadFile(url)
		if err != nil {
			log.Printf("Download %s failed with %v", filename, err)
			done <- false
			return
		}
		log.Printf("Download %s successful\n", filename)
		done <- true
	}()
	if !<-done {
		log.Println("Load default gfwlist.txt")
		bytes = resources.GFWListData
	}
	err = util.CreateFile(filename, bytes)
	if err != nil {
		log.Printf("Write %s failed:%v\n", filename, err)
	} else {
		log.Printf("Write %s successful\n", filename)
	}
	return bytes
}
