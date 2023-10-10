package gui

import "github.com/rivo/tview"

type appGrid struct {
	grid              *tview.Grid
	list              *tview.List
	detailPlaceholder *tview.TextView
}

type ListItem struct {
	title   string
	content string
}

func (grid appGrid) buildComponents() {
	grid.list = tview.NewList()
	grid.detailPlaceholder = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[red]Select a task[white]")
	grid.grid = tview.NewGrid()
}

func (grid appGrid) AddItemToList(listItem ListItem) {
	grid.list.AddItem(listItem.title, listItem.content, 0, nil)
}
