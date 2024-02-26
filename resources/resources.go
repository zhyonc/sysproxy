package resources

import (
	_ "embed"
)

//go:embed icon.ico
var IconData []byte

//go:embed abp.js
var AbpData []byte

//go:embed ggfwlist.txt
var GFWListData []byte
