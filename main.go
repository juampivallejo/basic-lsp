package main

import (
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

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}
		handleMessage(logger, method, contents)
	}
}

func handleMessage(logger *log.Logger, method string, contents []byte) {
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
		reply := rpc.EncodeMessage(msg)
		writer := os.Stdout
		writer.Write([]byte(reply))

		logger.Println("Sent the reply")
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
