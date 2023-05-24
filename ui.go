package main

import (
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
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

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			rw.WriteString(inputField.GetText())
			rw.Flush()
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

	//ui := UI{
	//	textView,
	//	app,
	//}

	//go streamReader(rw, ui)
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func streamReader(rw *bufio.ReadWriter, ui UI) {
	//tv := ui.tv
	//app := ui.app
	buf := make([]byte, 128)
	for {
		read, err := rw.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println(read)
	}
}
