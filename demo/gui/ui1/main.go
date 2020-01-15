package main

import (
	"fmt"
	"github.com/andlabs/ui"
	"github.com/xs25cn/scanPort/scan"
	_"github.com/andlabs/ui/winmanifest"
)

func main() {
	err := ui.Main(func() {
		ipInput := ui.NewEntry()
		portInput := ui.NewEntry()
		//processInput := ui.NewEntry()

		button := ui.NewButton("开始扫描")
		greeting := ui.NewLabel("")

		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel(" IP地址或域名(如: xs25.cn)"), false)
		box.Append(ipInput, false)
		box.Append(ui.NewLabel(" 端口号范围 (如: 80-1000)"), false)
		box.Append(portInput, false)
		//box.Append(ui.NewLabel("线程数:"), false)
		//box.Append(processInput, false)
		box.Append(ui.NewLabel(" "), false)
		box.Append(button, false)
		box.Append(greeting, false)

		//创建window窗口。并设置长宽。
		window := ui.NewWindow("端口扫描器 By:duzhenxun", 400, 300, true)
		//mac不支持居中。
		//https://github.com/andlabs/ui/issues/162
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			//开始处理
			greeting.SetText("  正在扫描中，请耐心等待......\n")

			ip := ipInput.Text()
			port := portInput.Text()
			//process,_ := strconv.Atoi(processInput.Text())

			scanIp := scan.NewScanIp(200, 200, true)

			ips,err:=scanIp.GetAllIp(ip)
			if err!=nil{
				greeting.SetText(" 输入信息有误：, " + err.Error() + "!")
			}

			for i:=0;i<len(ips);i++{
				ports:=scanIp.GetIpOpenPort(ips[i],port)

				if len(ports)>0{
					greeting.SetText(fmt.Sprintf(" %v 开放的端口：%v\n",ips[i],ports))
				}else{
					greeting.SetText(fmt.Sprintf(" %v 没有找到开放的端口\n",ips[i]))
				}
			}

		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
