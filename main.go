package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	fmt.Println("Menu:")

	app := tview.NewApplication()
	list := tview.NewList()
	grid := tview.NewGrid()
	detail_placeholder := tview.NewTextView().SetDynamicColors(true).SetText("[red]select a task[white]")

	detail_placeholder.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			app.SetFocus(list)

		}
		return event
	})

	pages := tview.NewPages()

	task1 := task{title: "title", content: "content", id: '1'}

	list.AddItem(task1.title, task1.content, 0, func() {
		detail := buildDetail(task1)
		detail_placeholder.SetText(detail)
		app.SetFocus(detail_placeholder)
	})
	list.AddItem("This is a list", "With description", 0, nil)
	list.AddItem("This is a list", "With description", 0, nil)
	list.AddItem("This is a list", "With description", 0, nil)
	list.AddItem("This is a list", "With description", 0, nil)
	list.AddItem("This is a list", "With description", 0, nil)
	list.AddItem("This is a list", "With description", 0, nil)

	grid.SetRows(0).
		SetColumns(50, 0).
		SetBorders(true)

	grid.SetGap(0, 1)

	grid.AddItem(list, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(detail_placeholder, 0, 1, 1, 1, 0, 0, false)

	frame := tview.NewFrame(grid)
	frame.AddText("Task list", true, tview.AlignCenter, tcell.ColorGreen)

	pages.AddPage("main", frame, true, true)

	app.SetRoot(pages, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				app.Stop()
			case 'n':
				// show modal
				form := buildForm(func() {
					pages.RemovePage("modal")
				}, func() {
					pages.RemovePage("modal")
				})
				modal := buildModal(form)
				pages.AddPage("modal", modal, true, true)

			}
		}

		return event

	})

	if err := app.Run(); err != nil {
		panic(err)
	}

}

func buildModal(content tview.Primitive) tview.Primitive {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	return modal(content, 60, 20)

}

func buildForm(onSave, onCancel func()) *tview.Form {
	form := tview.NewForm().
		AddInputField("Title", "", 40, nil, nil).
		AddTextArea("Description", "", 40, 10, 300, nil).
		AddButton("Save", onSave).
		AddButton("Quit", onCancel).
		SetButtonsAlign(tview.AlignCenter)

	form.SetTitle("New task").SetBorder(true).SetTitleAlign(tview.AlignCenter)

	return form

}

func buildDetail(task task) string {
	content := fmt.Sprintf("[green]%s[white]\n%s", task.title, task.content)

	return content
}
