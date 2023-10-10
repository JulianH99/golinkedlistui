package gui

import "github.com/rivo/tview"

type Shortcut struct {
	key    rune
	action func(app App)
}

type App struct {
	app     *tview.Application
	pages   *tview.Pages
	frame   *tview.Frame
	form    *tview.Form
	appGrid appGrid
}

func (app App) start() {
	if err := app.app.Run(); err != nil {
		panic(err)
	}
}

func (app App) buildComponents() {
	app.app = tview.NewApplication()
	app.pages = tview.NewPages()
	app.form = tview.NewForm()
	app.appGrid = appGrid{}
	app.appGrid.buildComponents()

	app.frame = tview.NewFrame(app.appGrid.grid)

}
