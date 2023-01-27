package application

import (
	"fmt"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type MessageProcessor struct {
	window *WebviewWindow
}

func NewMessageProcessor(w *WebviewWindow) *MessageProcessor {
	return &MessageProcessor{
		window: w,
	}
}

func (m *MessageProcessor) httpError(rw http.ResponseWriter, message string, args ...any) {
	m.Error(message, args...)
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(message))
}

func (m *MessageProcessor) HandleRuntimeCall(rw http.ResponseWriter, r *http.Request) {
	m.Info("Processing runtime call")
	// Read "method" from query string
	method := r.URL.Query().Get("method")
	if method == "" {
		m.httpError(rw, "No method specified")
		return
	}
	splitMethod := strings.Split(method, ".")
	if len(splitMethod) != 2 {
		m.httpError(rw, "Invalid method format")
		return
	}
	// Get the object
	object := splitMethod[0]
	// Get the method
	method = splitMethod[1]

	switch object {
	case "window":
		m.processWindowMethod(method, rw, r)
	default:
		m.httpError(rw, "Unknown runtime call: %s", object)
	}

}

func (m *MessageProcessor) ProcessMessage(message string) {
	m.Info("ProcessMessage from front end:", message)
}

func (m *MessageProcessor) Error(message string, args ...any) {
	fmt.Printf("[MessageProcessor] Error: "+message, args...)
}

func (m *MessageProcessor) Info(message string, args ...any) {
	fmt.Printf("[MessageProcessor] Info: "+message, args...)
}

func (m *MessageProcessor) json(rw http.ResponseWriter, data any) {
	// convert data to json
	var jsonPayload = []byte("{}")
	var err error
	if data != nil {
		jsonPayload, err = jsoniter.Marshal(data)
		if err != nil {
			m.Error("Unable to convert data to JSON. Please report this to the Wails team! Error: %s", err)
			return
		}
	}
	_, err = rw.Write(jsonPayload)
	if err != nil {
		m.Error("Unable to write json payload. Please report this to the Wails team! Error: %s", err)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
}
