package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"./translation"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: executable-name --port <port>")
		return
	}

	port := os.Args[2]
	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("Invalid port: " + err.Error())
	} else if portInt < 1 || portInt > 65535 {
		fmt.Println("Invalid port: port must be in the interval (0,65536)")
	}

	runServer(port)
}

func runServer(port string) {
	mux := http.NewServeMux()

	wordHandler := &translation.PostRequestHandler{
		Func: translation.WordTranslationHandler,
	}
	sentenceHandler := &translation.PostRequestHandler{
		Func: translation.SentenceTranslationHandler,
	}
	mux.Handle("/word", wordHandler)
	mux.Handle("/sentence", sentenceHandler)
	mux.HandleFunc("/history", translation.HistoryHandler)

	server := http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	server.ListenAndServe()
}
