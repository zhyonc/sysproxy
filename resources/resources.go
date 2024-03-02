package resources

import (
	_ "embed"
)

//go:embed Icon.ico
var IconData []byte

//go:embed IconI.ico
var IconIData []byte

//go:embed IconO.ico
var IconOData []byte

//go:embed IconIO.ico
var IconIOData []byte

//go:embed abp.js
var AbpData []byte

//go:embed ggfwlist.txt
var GFWListData []byte
