package main

import (
	"basiclsp/analysis"
	"basiclsp/lsp"
	"basiclsp/rpc"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fmt.Println("Hi")
	logger := getLogger("/home/juampi/projects/basic-lsp/log.txt")
	logger.Println("LSP Started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	responseWriter := GetResponseWriterFunc(os.Stdout)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}
		handleMessage(logger, responseWriter, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, sendResponse func(msg any), state analysis.State, method string, contents []byte) {
	logger.Printf("Received msg with method: %s", method)

	switch method {

	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error Parsing Initialize request \n%s", err)
		}
		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		// Reply
		msg := lsp.NewInitializeResponse(request.ID)
		sendResponse(msg)

		logger.Println("Sent initialize reply")

	case "textDocument/didOpen":
		var notification lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Error Parsing textDocument/didOpen notification \n%s", err)
		}
		logger.Printf("Opened: %s %s", notification.Params.TextDocument.URI, notification.Params.TextDocument.LanguageId)

		diagnostics := state.OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)
		sendResponse(lsp.PublishDiagnosticNotification{
			Notification: lsp.Notification{RPC: "2.0", Method: "textDocument/publishDiagnostics"},
			Params: lsp.PublishDiagnosticParams{
				URI:         notification.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})

	case "textDocument/didChange":
		var notification lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Error Parsing textDocument/didChange notification \n%s", err)
		}
		logger.Printf("Changed: %s %d", notification.Params.TextDocument.URI, notification.Params.TextDocument.Version)
		for _, change := range notification.Params.ContentChanges {
			diagnostics := state.UpdateDocument(notification.Params.TextDocument.URI, change.Text)
			sendResponse(lsp.PublishDiagnosticNotification{
				Notification: lsp.Notification{RPC: "2.0", Method: "textDocument/publishDiagnostics"},
				Params: lsp.PublishDiagnosticParams{
					URI:         notification.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error Parsing textDocument/hover request\n%s", err)
		}
		msg := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		sendResponse(msg)

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error Parsing textDocument/definition request\n%s", err)
		}
		msg := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		sendResponse(msg)

	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error Parsing textDocument/codeAction request\n%s", err)
		}
		msg := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)
		sendResponse(msg)

	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error Parsing textDocument/completion request\n%s", err)
		}
		msg := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI)
		sendResponse(msg)
	}

}

func GetResponseWriterFunc(writer io.Writer) func(msg any) {
	return func(msg any) {
		reply := rpc.EncodeMessage(msg)
		writer.Write([]byte(reply))
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Error: cannot get Logger")
	}

	l := new(log.Logger)
	l.SetOutput(io.Writer(logFile))
	l.SetPrefix("[basiclsp]")
	l.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return l
}
