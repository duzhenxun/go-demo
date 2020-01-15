package main

import "github.com/dontpanic92/wxGo/wx"

func main() {

	wx.NewApp()
	f := wx.NewDialog(wx.NullWindow, -1, "Hello World")

	bSizer := wx.NewBoxSizer(wx.VERTICAL)

	checkBox := wx.NewCheckBox(f, wx.ID_ANY, "Check Me!", wx.DefaultPosition, wx.DefaultSize, 0)
	bSizer.Add(checkBox, 0, wx.ALL|wx.EXPAND, 5)

	textCtrl := wx.NewTextCtrl(f, wx.ID_ANY, "", wx.DefaultPosition, wx.DefaultSize, 0)
	bSizer.Add(textCtrl, 0, wx.ALL|wx.EXPAND, 5)

	f.SetSizer(bSizer)
	f.Layout()
	f.ShowModal()
	f.Destroy()
}




























