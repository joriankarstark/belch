package main

import (
	"flag"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
)

var (
	// Global Widgets
	requestHistory *widgets.List
	request        *widgets.Paragraph
	response       *widgets.Paragraph
	footer         *widgets.Paragraph

	termWidth  int
	termHeight int

	defaultRequestHistoryHeight = 10

	requestChan  chan *http.Request
	responseChan chan *http.Response

	//requestHistoryList []string
)

func main() {
	var pemPath string
	var keyPath string
	var proto string

	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	check(err)
	defer f.Close()

	log.SetOutput(f)

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

	requestChan = make(chan *http.Request, 2)
	responseChan = make(chan *http.Response, 2)

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
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "h":
				ui.Clear()
				ui.Render(requestHistory)
			case "l":
				ui.Clear()
				ui.Render(requestHistory)

			case "<Up>":
				requestHistory.ScrollUp()

				ui.Render(requestHistory)
			case "<Down>":
				requestHistory.ScrollDown()

				ui.Render(requestHistory)
			case "<Resize>":
				ui.Clear()
				resetWidgetSizes()
				ui.Render(requestHistory, request, response, footer)
			}

		case r := <-requestChan:
			//requestHistoryList = append(requestHistoryList, r)
			//updateRequestHistory()
			requestHistory.Rows = append(requestHistory.Rows, r.RequestURI)
			ui.Render(requestHistory)

			dump, _ := httputil.DumpRequest(r, true)
			request.Text = string(dump)
			ui.Render(request)

		}

	}
}

func initWidgets() {

	termWidth, termHeight = ui.TerminalDimensions()

	requestHistory = widgets.NewList()
	requestHistory.Rows = make([]string, 0)
	requestHistory.SetRect(1, 1, termWidth-1, defaultRequestHistoryHeight)
	requestHistory.Border = true
	requestHistory.Title = "Request History"
	requestHistory.SelectedRowStyle.Bg = ui.ColorBlue

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

func resetWidgetSizes() {

	termWidth, termHeight = ui.TerminalDimensions()
	footer.SetRect(1, termHeight-4, termWidth-1, termHeight-1)
	request.SetRect(1, defaultRequestHistoryHeight, ((termWidth-1)/2)-1, termHeight-4)
	response.SetRect(((termWidth-1)/2)+1, defaultRequestHistoryHeight, termWidth-1, termHeight-4)
	requestHistory.SetRect(1, 1, termWidth-1, defaultRequestHistoryHeight)
}

func editText(s string) string {
	content := []byte(s)
	tmp, err := ioutil.TempFile("", ".request.tmp")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmp.Name())

	tmp.Write(content)
	tmp.Close()

	cmd := exec.Command("vim", tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
	text, err := ioutil.ReadFile(tmp.Name())
	return (string(text))

}
