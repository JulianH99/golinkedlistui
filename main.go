package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Shortcut rune

const (
	ShortcutAdd      Shortcut = 'n'
	ShortcutDelete            = 'd'
	ShortcutComplete          = 'c'
	ShortcutQuit              = 'q'
)

func shortcutToText(shortcut Shortcut) string {
	switch shortcut {
	case ShortcutAdd:
		return fmt.Sprintf("(%c) New task", ShortcutAdd)

	case ShortcutDelete:
		return fmt.Sprintf("(%c) Mark as complete", ShortcutDelete)

	case ShortcutComplete:
		return fmt.Sprintf("(%c) Delete task", ShortcutComplete)

	}
	return ""
}

func addTaskToList(list *tview.List, app *tview.Application, detailPlaceholder *tview.TextView, t *Task) {

	list.AddItem(t.title, t.content, 0, func() {
		detail := buildDetail(t)
		detailPlaceholder.SetText(detail)
		app.SetFocus(detailPlaceholder)
	})
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentTask := list.GetCurrentItem()
		switch event.Rune() {
		case 'd':
			list.RemoveItem(currentTask)
		case 'c':
			list.SetItemText(currentTask, "Completed:[yellow]"+t.title+"[white]", t.content)
		}

		return event

	})

}

func createDetailPlaceholder(app *tview.Application, list *tview.List) *tview.TextView {
	detailPlaceholder := tview.NewTextView().SetDynamicColors(true).SetText("[blue]select a task[white]")
	detailPlaceholder.SetScrollable(false)
	detailPlaceholder.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			app.SetFocus(list)

		}
		return event
	})

	return detailPlaceholder
}

func createCommandsContainer() *tview.Flex {

	commandsContainer := tview.NewFlex().
		AddItem(tview.NewTextView().SetText(shortcutToText(ShortcutAdd)), 0, 1, false).
		AddItem(tview.NewTextView().SetText(shortcutToText(ShortcutComplete)), 0, 1, false).
		AddItem(tview.NewTextView().SetText(shortcutToText(ShortcutDelete)), 0, 1, false)

	return commandsContainer
}

func createGrid(list *tview.List, detailPlaceholder *tview.TextView, commandsContainer *tview.Flex) *tview.Grid {

	grid := tview.NewGrid()
	grid.SetRows(0, 1).
		SetColumns(50, 0).
		SetBorders(true)

	grid.SetGap(0, 1)

	grid.AddItem(list, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(detailPlaceholder, 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(commandsContainer, 1, 0, 1, 2, 0, 0, false)

	return grid
}

func main() {
	app := tview.NewApplication()
	list := tview.NewList()
	detailPlaceholder := createDetailPlaceholder(app, list)
	commandsContainer := createCommandsContainer()
	grid := createGrid(list, detailPlaceholder, commandsContainer)
	formShown := false

	pages := tview.NewPages()

	frame := tview.NewFrame(grid)
	frame.AddText("Task list", true, tview.AlignCenter, tcell.ColorGreen)

	pages.AddPage("main", frame, true, true)

	app.SetRoot(pages, true)

	onSave := func(f *tview.Form) func() {
		return func() {
			title := f.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			content := f.GetFormItemByLabel("Description").(*tview.TextArea).GetText()

			task := createTask(title, content)
			addTaskToList(list, app, detailPlaceholder, task)
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
			case rune(ShortcutQuit):
				app.Stop()
			case rune(ShortcutAdd):
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

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			onCancel()
		}
		return event
	})

	return form

}

func buildDetail(task *Task) string {
	content := fmt.Sprintf("[green]%s[white]\n%s", task.title, task.content)

	return content
}
