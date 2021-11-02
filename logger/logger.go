package logger

import (
	"encoding/json"
	"log"
	"time"

	"github.com/telkomdev/go-stash"
)

var Logger *log.Logger
var st *stash.Stash

type LogMessage struct {
	TS      time.Time `json:"timestamp"`
	Message string    `json:"message"`
	LogType string    `json:"type"`
}

func InitLogger() {
	st, err := stash.Connect("localhost", 5000)
	if err != nil {
		log.Fatal(err)
	}
	Logger = log.New(st, "", 0)

}

func Write(infotype string, message string) {
	data := LogMessage{
		TS:      time.Now(),
		Message: message,
		LogType: infotype,
	}
	byteData, _ := json.Marshal(&data)
	Logger.Println(string(byteData))
}

func CloseStash() {
	st.Close()
}
