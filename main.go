package main

import (
	"fmt"

	"calendar/config"
	events "calendar/event"
	"calendar/gui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/systray"
	log "github.com/sirupsen/logrus"
)

var (
	calendarApp fyne.App
	editWindow  fyne.Window
)

func main() {
	if err := dbOpenOrCreate(); err != nil {
		log.Error(err)
		return
	}

	calendarApp = app.New()
	config.Config().SetDarkTheme(true)
	editWindow = gui.InitDataWindow(windowWithTitle("Anniversaries"))
	editWindow.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		e := events.NewEvent(events.KeyEvent, map[string]any{"event": event})
		events.Instance().Send(e)
	})

	nativeStart, nativeEnd := systray.RunWithExternalLoop(onReady, onExit)

	nativeStart()
	calendarApp.Run()
	nativeEnd()
}

func windowWithTitle(title string) fyne.Window {
	window := calendarApp.NewWindow(title)
	window.Resize(fyne.NewSize(1000, 600))
	window.CenterOnScreen()
	window.Hide()
	window.SetCloseIntercept(func() {
		window.Hide()
		windowClosed()
	})
	return window
}

func windowClosed() {
	e := events.NewEventWithID(events.WindowClosed)
	events.Instance().Send(e)
}

func windowOpened() {
	e := events.NewEventWithID(events.WindowOpened)
	events.Instance().Send(e)
}

// onReady
// 1. utworzenie menu dla tray'a
// 2. oczekiwanie na event gdy wybrano pozycję menu
// 3. oczekiwanie na eventy gdy otwarto/zamknięto okno
func onReady() {
	systray.SetTitle("Anniversaries")
	mode := systray.AddMenuItemCheckbox("Dark theme", "", config.Config().DarkTheme())
	edit := systray.AddMenuItem("Personal data...", "")
	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "")

	// Oczekiwanie w tle na eventy związane z
	// wyborem pozycji menu.
	go func() {
		for {
			select {
			case <-mode.ClickedCh:
				if mode.Checked() {
					mode.Uncheck()
					config.Config().ChangeTheme()
				} else {
					mode.Check()
					config.Config().ChangeTheme()
				}
			case <-edit.ClickedCh:
				editWindow.Show()
				windowOpened()
			case <-quit.ClickedCh:
				calendarApp.Quit()
				systray.Quit()
				return
			}
		}
	}()

	// Oczekiwanie w tle na eventy związane z
	// otwieraniem/zamykaniem okien.
	go func(stream <-chan *events.Event) {
		for e := range stream {
			switch e.Id() {
			case events.WindowOpened:
				edit.Disable()
			case events.WindowClosed:
				edit.Enable()
			}
		}
	}(events.Instance().Register("Window open/close", events.WindowOpened, events.WindowClosed))
}

// onExit przy wyjściu z programu zamykamy wszystkie okna
// (jeśli zostały utworzone).
func onExit() {
	if editWindow != nil {
		editWindow.Close()
		editWindow = nil
	}
	fmt.Println("Calendar finished")
}
