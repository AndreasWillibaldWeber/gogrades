package gui

import "fyne.io/fyne/v2"

func (g *GUI) buildMainMenu() *fyne.MainMenu {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open CSV...", g.openCSVDialog),
		fyne.NewMenuItem("Save CSV...", g.saveCSV),
		fyne.NewMenuItemSeparator(),
	)
	return fyne.NewMainMenu(fileMenu)
}
