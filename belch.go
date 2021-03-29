package main

import (
	"flag"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	// Global Widgets
	requestHistory *widgets.Paragraph
	request        *widgets.Paragraph
	response       *widgets.Paragraph
	footer         *widgets.Paragraph

	termWidth  int
	termHeight int

	defaultRequestHistoryHeight = 10
)

func main() {
	var pemPath string
	var keyPath string
	var proto string

	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")
	flag.StringVar(&keyPath, "key", "server.key", "path to key file")
	flag.StringVar(&proto, "proto", "http", "Proxy protocol (http or https)")
	flag.Parse()

	if proto != "http" && proto != "https" {
		log.Fatal("Protocol must be either http or https")
	}
	log.Println(pemPath)
	log.Println(keyPath)
	log.Println(proto)
	go StartProxy(pemPath, keyPath, proto)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	initWidgets()

	eventLoop()

}

func eventLoop() {

	ui.Render(requestHistory, request, response, footer)

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "h":
			ui.Clear()
			ui.Render(requestHistory)
		case "l":
			ui.Clear()
			ui.Render(requestHistory)

		case "<Resize>":
			ui.Clear()
			initWidgets()
			ui.Render(requestHistory, request, response, footer)
		}
	}
}

func initWidgets() {

	termWidth, termHeight = ui.TerminalDimensions()

	requestHistory = widgets.NewParagraph()
	requestHistory.Text = "My requests from web proxy will go here....\ntest\ntest\ntest"
	requestHistory.SetRect(1, 1, termWidth-1, defaultRequestHistoryHeight)
	requestHistory.Border = true
	requestHistory.Title = "Request History"

	request = widgets.NewParagraph()
	request.Text = "My current request from web proxy will go here....\ntest\ntest\ntest"
	request.SetRect(1, defaultRequestHistoryHeight, ((termWidth-1)/2)-1, termHeight-4)
	request.Border = true
	request.Title = "Request"

	response = widgets.NewParagraph()
	response.Text = "My current response from web proxy will go here....\ntest\ntest\ntest"
	response.SetRect(((termWidth-1)/2)+1, defaultRequestHistoryHeight, termWidth-1, termHeight-4)
	response.Border = true
	response.Title = "Response"

	footer = widgets.NewParagraph()
	footer.Text = "Shortcuts and info will go here..."
	footer.SetRect(1, termHeight-4, termWidth-1, termHeight-1)
	footer.Border = true
	footer.Title = "info/status bar..."
}
