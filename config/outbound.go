package config

type Outbound struct {
	Tag      string
	SrcProto string
	SrcIP    string
	SrcPort  string
	DstProto string
	DstIP    string
	DstPort  string
}
