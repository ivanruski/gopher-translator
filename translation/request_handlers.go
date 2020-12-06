package translation

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type PostRequestHandler struct {
	Func func(http.ResponseWriter, *http.Request)
}

func (handler *PostRequestHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.NotFound(rw, request)
		return
	} else if request.Header.Get("Content-Type") != "application/json" {
		http.Error(rw, "Server accepts only Content-Type: application/json", http.StatusUnsupportedMediaType)
		return
	}

	handler.Func(rw, request)
}

func WordTranslationHandler(rw http.ResponseWriter, request *http.Request) {
	parsedBody := struct {
		Word string `json:"english-word"`
	}{}

	if err := parseRequestBody(request.Body, &parsedBody); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	} else if parsedBody.Word == "" {
		http.Error(rw, "Request body must have 'english-word' field in it", http.StatusBadRequest)
	} else {
		translatedWord, err := TranslateWord(DefaultContext, parsedBody.Word)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		} else {
			result := fmt.Sprintf("{ \"gopher-word\" : \"%s\" }", translatedWord)
			writeOkResponse(rw, result)
		}
	}
}

func SentenceTranslationHandler(rw http.ResponseWriter, request *http.Request) {
	parsedBody := struct {
		Sentence string `json:"english-sentence"`
	}{}

	if err := parseRequestBody(request.Body, &parsedBody); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	} else if parsedBody.Sentence == "" {
		http.Error(rw, "Request body must have 'english-sentence' field in it", http.StatusBadRequest)
	} else {
		translatedSentence, err := TranslateSentence(DefaultContext, parsedBody.Sentence)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		} else {
			result := fmt.Sprintf("{ \"gopher-sentence\" : \"%s\" }", translatedSentence)
			writeOkResponse(rw, result)
		}
	}
}

func HistoryHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.NotFound(rw, request)
		return
	}

	writeOkResponse(rw, ExportCachedItemsAsJson())
}

func parseRequestBody(reader io.Reader, v interface{}) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	parseError := json.Unmarshal(content, v)

	return parseError
}

func writeOkResponse(rw http.ResponseWriter, msg string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	if msg != "" {
		rw.Write([]byte(msg))
	}
}
