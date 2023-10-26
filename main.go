package main

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func addTaskToList(list *tview.List, app *tview.Application, detailPlaceholder *tview.TextView, t task) {

	list.AddItem(t.title, t.content, 0, func() {
		detail := buildDetail(t)
		detailPlaceholder.SetText(detail)
		app.SetFocus(detailPlaceholder)
	})
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'd':
			currentTask := list.GetCurrentItem()
			list.RemoveItem(currentTask)
		}

		return event

	})

}

func main() {
	fmt.Println("Menu:")

	app := tview.NewApplication()
	list := tview.NewList()
	grid := tview.NewGrid()
	detail_placeholder := tview.NewTextView().SetDynamicColors(true).SetText("[red]select a task[white]")
	detail_placeholder.SetScrollable(false)
	formShown := false

	detail_placeholder.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			app.SetFocus(list)

		}
		return event
	})

	pages := tview.NewPages()

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

	onSave := func(f *tview.Form) func() {
		return func() {
			title := f.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			description := f.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			task := task{title: title, content: description, id: rand.Intn(100)}
			addTaskToList(list, app, detail_placeholder, task)
			pages.RemovePage("modal")
			formShown = false

			// restore focus
			f.SetFocus(0)
			f.GetFormItemByLabel("Title").(*tview.InputField).SetText("")
			f.GetFormItemByLabel("Description").(*tview.TextArea).SetText("", false)
		}
	}

	onCancel := func() {
		pages.RemovePage("modal")
		formShown = false
	}

	form := buildForm(onSave, onCancel)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				app.Stop()
			case 'n':
				// show modal
				if !formShown {
					modal := buildModal(form)
					pages.AddPage("modal", modal, true, true)
					form.GetFormItemByLabel("Title").(*tview.InputField).SetText("")
					formShown = true
				}
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

func buildForm(onSave func(f *tview.Form) func(), onCancel func()) *tview.Form {
	form := tview.NewForm().
		AddInputField("Title", "", 40, nil, nil).
		AddTextArea("Description", "", 40, 10, 300, nil)
	form.AddButton("Save", onSave(form)).
		AddButton("Quit", onCancel).
		SetButtonsAlign(tview.AlignCenter)
	form.SetTitle("New task").SetBorder(true).SetTitleAlign(tview.AlignCenter)

	return form

}

func buildDetail(task task) string {
	content := fmt.Sprintf("[green]%s[white]\n%s", task.title, task.content)

	return content
}
