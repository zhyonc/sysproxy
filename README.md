<img src="/resources/Icon.png" alt="[logo]" width="48"/> sysproxy
=======================
### What it is
Windows system proxy forward tool that display as tray menu

### Features
- Support system proxy protocol is pac/http/socks5
- Extend forward protocol is http/socks5
- Customize domain rules to pac file

### Download
Download the compiled version from [release page](https://github.com/zhyonc/sysproxy/releases)

### How to use
#### Inbound
The "inbound" setting represents the configuration of the System Proxy.  
Therefore the changes in this section will affect to all processes that uses the global Internet configuration.  
If you don't want to use system proxy, not need to set "inbound".
- Click Setting menu item to open inbound form
- Input and add inbound info that controls system proxy
- After modified inbound list, don't forget to click update button
- Click Inbound menu item to select inbound tag
- After the tag was selected, call win+R and enter inetcpl.cpl to see the LAN changes
![Inbound](/resources/Inbound.png)
#### Outbound
The "outbound" setting represents the different proxy configurations of this program.  
You can forward the Inbound to any of these proxies, or use directly them from any process configured to use it.
- Click Setting menu item to open outbound form
- Input and add outbound info that forward inbound or other proxy chain
- After modified outbound list, don't forget to click update button
- Click Outbound menu item to select outbound tag
- Click Log menu item to open log form
- If everything is ok, you can see the connection status in the log form
![Outbound](/resources/Outbound.png)
#### Proxy Toggle
Click the tray icon to enable/disable proxy
- <img src="/resources/Icon.png" alt="[Icon]" width="16" height="16"/> Both inbound and outbound is disabled
- <img src="/resources/IconI.png" alt="[IconI]" width="16" height="16"/> Only inbound is enabled
- <img src="/resources/IconO.png" alt="[IconO]" width="16" height="16"/> Only outbound is enabled
- <img src="/resources/IconIO.png" alt="[IconIO]" width="16" height="16"/> Both inbound and outbound is enabled
### Build
- Upgrade go version to above 1.9.2
- Download package tool from [rsrc](https://github.com/akavel/rsrc/releases)
- Package manifest and icon to syso file:
- ```rsrc.exe -manifest sysproxy.manifest -ico ./resources/icon.ico -o resources.syso```
- To get rid of the cmd window:
- ```go build -ldflags="-s -w -H windowsgui"```

### Thanks
[Govcl](https://github.com/ying32/govcl) is Cross-platform Golang GUI library binding [liblcl](https://github.com/ying32/liblcl)

### License
[MIT](LICENSE)
