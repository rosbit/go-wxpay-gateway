//+build gateway

/**
 * ResultSaver goroutine
 * Rosbit Xu
 */
package utils

import (
	"time"
	"fmt"
	"os"
	"log"
	"encoding/json"
)

const (
	MAX_Q_LEN = 5
	DATE_FORMAT = "2006-01-02 15:04:05"
)

var (
	resultQ chan interface{}
	exit bool = false
)

func SaveResult(result interface{}) {
	resultQ <- result
}

func StopSaver() {
	exit = true
	resultQ <- nil
}

func saverThread(savePath string) {
	var body []byte
	var err error
	var result interface{}

	log.Printf("[result saver] is working...\n")
	for !exit {
		result = <-resultQ
		if result == nil {
			continue
		}

		switch result.(type) {
		case []byte:
			body = result.([]byte)
		case string:
			body = []byte(result.(string))
		default:
			if body, err = json.Marshal(result); err != nil {
				log.Printf("[result saver] failed to encoding result: %v\n", err)
				continue
			}
		}

		t := time.Now()
		now := t.Format(DATE_FORMAT)
		bl := len(body) + 1

		f, err := os.OpenFile(savePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("[result saver] failed to open file(%s) to save: %v\n", savePath, err)
			continue
		}

		fmt.Fprintf(f, "%s\t%d\n", now, bl)
		f.Write(body)
		f.WriteString("\n")
		f.Close()
	}
	log.Printf("[result saver] will exit.\n")
}

func StartSaver(savePath string) {
	resultQ = make(chan interface{}, MAX_Q_LEN)
	go saverThread(savePath)
}

