package linuxGui

import (
	"Metamorphoun/config"
	"Metamorphoun/server"
	"Metamorphoun/zutil"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func MakeGui() {
	ui := NewUI()

	// This creates a new application window and starts the UI.
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Counter"),
			app.Size(unit.Dp(240), unit.Dp(70)),
		)
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	// This starts Gio main.
	app.Main()
}

type UI struct {
	Theme   *material.Theme
	Button1 widget.Clickable
	Button2 widget.Clickable
	Window  *app.Window
}

func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	ui.Theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	// listen for events happening on the window.
	for {
		// detect the type of the event.
		switch e := w.Event().(type) {
		// this is sent when the application should re-render.
		case app.FrameEvent:
			// gtx is used to pass around rendering and event information.
			gtx := app.NewContext(&ops, e)

			// register a global key listener for the escape key wrapping our entire UI.
			area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			event.Op(gtx.Ops, w)

			// check for presses of the escape key and close the window if we find them.
			for {
				event, ok := gtx.Event(key.Filter{
					Name: key.NameEscape,
				})
				if !ok {
					break
				}
				switch event := event.(type) {
				case key.Event:
					if event.Name == key.NameEscape {
						return nil
					}
				}
			}
			// render and handle UI.
			ui.Layout(gtx)
			area.Pop()
			// render and handle the operations from the UI.
			e.Frame(gtx.Ops)

		// this is sent when the application is closed.
		case app.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			for ui.Button1.Clicked(gtx) {
				urlSettings := "http://" + config.ConfigInstance.ServerAddress + ":" + zutil.AsString(config.ConfigInstance.ServerPort)
				server.OpenFolder("explorer", urlSettings)
			}
			return material.Button(ui.Theme, &ui.Button1, "Settings").Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			for ui.Button2.Clicked(gtx) {
				currPicInfo := "http://" + config.ConfigInstance.ServerAddress + ":" + zutil.AsString(config.ConfigInstance.ServerPort) + "/picInfo.html"
				server.OpenFolder("explorer", currPicInfo)
			}
			return material.Button(ui.Theme, &ui.Button2, "Action").Layout(gtx)
		}),
	)
}
