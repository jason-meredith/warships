package main

import ui "github.com/gizak/termui"

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	titleTxt := ui.NewPar("\n\n\n[Warships](fg-red,fg-bold)")
	titleTxt.Height = 10
	titleTxt.Width = 10
	titleTxt.Y = 1
	titleTxt.X = 35
	titleTxt.Border = false


	menuTxt := ui.NewPar("\n[[S]](fg-bold)tart New Game\n[[J]](fg-bold)oin Game\n[[Q]](fg-bold)uit")
	menuTxt.Height = 8
	menuTxt.Width = 40
	menuTxt.Y = 4
	menuTxt.X = 20
	menuTxt.Border = false

	authorTxt := ui.NewPar("By Jason Meredith")
	authorTxt.Height = 3
	authorTxt.Width = 30
	authorTxt.X = 25
	authorTxt.Y = 12
	authorTxt.Border = false

	/*
	ui.Render(titleTxt, menuTxt, authorTxt)
*/

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(1, 5, titleTxt)),
		ui.NewRow(
			ui.NewCol(4, 4, menuTxt)),
		ui.NewRow(
			ui.NewCol(4, 4, authorTxt)),
	)


	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)



	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop()



}