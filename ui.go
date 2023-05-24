/*
	TviewとStreamの兼ね合いが上手く出来なかったので諦めた名残
*/

package main

import (
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"time"
)

type UI struct {
	tv  *tview.TextView
	app *tview.Application
}

func runUI(rw *bufio.ReadWriter) {
	app := tview.NewApplication()

	textView := tview.NewTextView()
	textView.SetTitle("chat-lobby")
	textView.SetBorder(true)

	inputField := tview.NewInputField()
	inputField.SetLabel("input> ")
	inputField.SetTitle("sendMessage").
		SetBorder(true)

	ui := UI{
		tv:  textView,
		app: app,
	}
	go streamReader(rw, &ui)

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		w := rw.Writer
		switch event.Key() {
		case tcell.KeyEnter:
			w.WriteString(inputField.GetText())
			textView.SetText(textView.GetText(true) + inputField.GetText() + "\n")
			inputField.SetText("")
			return nil
		case tcell.KeyCtrlC:
			os.Exit(1)
		}
		return event
	})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(textView, 0, 2, false).
		AddItem(inputField, 3, 0, true)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func streamReader(rw *bufio.ReadWriter, ui *UI) {
	app := ui.app
	tv := ui.tv
	for {
		time.Sleep(time.Millisecond * 100)
		app.QueueUpdateDraw(func() {
			tv.SetText(tv.GetText(true) + "fuck\n")
		})
		str, err := rw.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if str == "" {
			continue
		}

		if str != "\n" {
			fmt.Fprintln(tv, str)
		}
	}
}
